package saga

import (
	"fmt"
	"slices"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/legendaryum-metaverse/saga/event"
)

func getEventKey(event event.MicroserviceEvent) string {
	return strings.ToUpper(string(event))
}

func getEventObject(event event.MicroserviceEvent) amqp.Table {
	key := getEventKey(event)
	return amqp.Table{key: string(event)}
}

// createHeaderConsumer creates the exchanges, queues, and bindings for the given microservice and events.
func (t *Transactional) createHeaderConsumer(queueName string, events []event.MicroserviceEvent) error {
	requeueQueue := fmt.Sprintf("%s_matching_requeue", queueName)

	for _, exchange := range []Exchange{MatchingExchange, MatchingRequeueExchange} {
		err := t.eventsChannel.ExchangeDeclare(string(exchange), "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", string(exchange), err)
		}
	}

	_, err := t.eventsChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}
	_, err = t.eventsChannel.QueueDeclare(requeueQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": string(MatchingExchange),
	})
	if err != nil {
		return fmt.Errorf("failed to declare requeue queue %s: %w", requeueQueue, err)
	}

	// Handle individual events
	for _, ev := range event.MicroserviceEventValues() {
		headerEvent := getEventObject(ev)

		// Declare and bind event-specific exchanges
		err = t.eventsChannel.ExchangeDeclare(string(ev), "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", ev, err)
		}
		headersArgs := amqp.Table{
			"all-micro": "yes",
			"x-match":   "all",
		}
		for k, v := range headerEvent {
			headersArgs[k] = v
		}
		err = t.eventsChannel.ExchangeBind(string(ev), "", string(MatchingExchange), false, headersArgs)
		if err != nil {
			return fmt.Errorf("failed to bind exchange %s to %s: %w", ev, MatchingExchange, err)
		}

		requeueExchange := fmt.Sprintf("%s_requeue", ev)

		err = t.eventsChannel.ExchangeDeclare(requeueExchange, "headers", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare requeue exchange %s: %w", requeueExchange, err)
		}
		err = t.eventsChannel.ExchangeBind(requeueExchange, "", string(MatchingRequeueExchange), false, headerEvent)
		if err != nil {
			return fmt.Errorf("failed to bind requeue exchange %s to %s: %w", requeueExchange, MatchingRequeueExchange, err)
		}

		headersArgs = amqp.Table{
			"micro":   queueName,
			"x-match": "all",
		}
		for k, v := range headerEvent {
			headersArgs[k] = v
		}
		if slices.Contains(events, ev) {
			// Bindings for included events
			err = t.eventsChannel.QueueBind(queueName, "", string(ev), false, headerEvent)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, ev, err)
			}

			err = t.eventsChannel.QueueBind(requeueQueue, "", fmt.Sprintf("%s_requeue", ev), false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind requeue queue %s to exchange %s_requeue: %w", requeueQueue, ev, err)
			}

			microExchange := fmt.Sprintf("%s_%s", ev, queueName)
			err = t.eventsChannel.ExchangeDeclare(microExchange, "headers", true, false, false, false, nil)
			if err != nil {
				return fmt.Errorf("failed to declare exchange %s: %w", microExchange, err)
			}
			err = t.eventsChannel.ExchangeBind(microExchange, "", string(MatchingExchange), false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind exchange %s to %s: %w", microExchange, MatchingExchange, err)
			}

			err = t.eventsChannel.QueueBind(queueName, "", microExchange, false, headersArgs)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, microExchange, err)
			}

		} else {
			// Attempt to unbind the queue, ignoring errors if it's already unbound
			err = t.eventsChannel.QueueUnbind(queueName, "", string(ev), headerEvent)
			if err != nil {
				return fmt.Errorf("failed to unbind queue %s from exchange %s: %w", queueName, ev, err)
			}
			err = t.eventsChannel.QueueUnbind(requeueQueue, "", fmt.Sprintf("%s_requeue", ev), headersArgs)
			if err != nil {
				return fmt.Errorf("failed to unbind requeue queue %s from exchange %s_requeue: %w", requeueQueue, ev, err)
			}
			err = t.eventsChannel.ExchangeDelete(fmt.Sprintf("%s_%s", ev, queueName), false, false)
			if err != nil {
				return fmt.Errorf("failed to delete exchange %s_%s: %w", ev, queueName, err)
			}
		}
	}
	t.healthCheckQueue = queueName
	return nil
}
