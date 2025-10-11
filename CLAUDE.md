# 📖 CLAUDE.md - legend-saga Golang Library

> **Technical Documentation for Claude Code Assistant**
> Updated: 2025-10-11 | Go 1.25.1 | Production-Ready

---

## Quick Reference

### Essential Patterns

```go
// 1. Initialize
transactional := saga.Config(&saga.Opts{
    RabbitUri:    "amqp://localhost:5672",
    Microservice: micro.Social,
    Events:       []event.MicroserviceEvent{event.AuthNewUserEvent},
})

// 2. Publish Event
saga.PublishEvent(&event.AuthNewUserPayload{ID: "123", Email: "user@example.com"})

// 3. Consume Event
emitter := transactional.ConnectToEvents()
emitter.On(event.AuthNewUserEvent, func(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.AuthNewUserPayload{})
    // Process event...
    handler.Channel.AckMessage() // Success - emits audit.processed
    // OR
    handler.Channel.NackWithFibonacciStrategy(19, 3) // Failure - emits audit.dead_letter
})

// 4. Saga Orchestration
saga.CommenceSaga(&saga.PurchaseResourceFlowPayload{...})
sagaEmitter := transactional.ConnectToSagaCommandEmitter()
sagaEmitter.On(micro.SomeCommand, func(handler saga.CommandHandler) {
    // Handle saga step...
    handler.Channel.AckMessage(saga.NextStepPayload{...})
})
```

---

## Project Overview

**legend-saga** is a Go library for microservice communication via RabbitMQ, supporting event-driven architecture and saga orchestration patterns.

### Key Features

- 🔄 Publish/Subscribe with headers-based routing
- 🎭 Saga orchestration for distributed transactions
- 🔍 **Automatic audit logging** (lifecycle tracking)
- ♻️ Fibonacci backoff retry mechanism
- 🛡️ Durable message delivery
- 🎯 Type-safe Go generics

### Tech Stack

- **Language**: Go 1.25.1
- **Broker**: RabbitMQ (AMQP 0.9.1)
- **Dependencies**: `rabbitmq/amqp091-go v1.10.0`, `go-playground/validator/v10`

---

## Architecture

### System Flow

```
Application
    ↓
Transactional (start.go)
    ↓
┌────────────────┬─────────────────────────┐
↓                ↓                         ↓
ConnectToEvents  ConnectToSagaCommand     (Audit Infrastructure)
    ↓                ↓                         ↓
Event Emitter    Saga Emitter         audit_exchange (direct)
    ↓                ↓                         ↓
matching_exchange commands_exchange    3 Audit Queues
(headers)        (direct)              - audit_received_commands
    ↓                ↓                 - audit_processed_commands
Event Consumers  Saga Consumers        - audit_dead_letter_commands
```

### Core Components

| Component              | Location                   | Purpose                              |
| ---------------------- | -------------------------- | ------------------------------------ |
| `Transactional`        | start.go:24                | Main coordinator, connection manager |
| `Emitter[T, U]`        | emitter.go:7               | Type-safe event dispatcher           |
| `EventHandler`         | consumeCallbackEvent.go    | Event processing handler             |
| `CommandHandler`       | consumeCallbackSagaCommand | Saga step handler                    |
| `EventsConsumeChannel` | consumeChannelEvents.go    | Channel wrapper with audit emission  |

### Directory Structure

```
legend-saga/
├── event/microserviceEvent.go    # 45 events (42 business + 3 audit)
├── micro/sagaCommands.go          # 17 microservices + saga commands
├── test/                          # Integration tests
│
# Core files (~1,400 LOC)
├── start.go                       # Entry point, Transactional, ConnectTo* methods
├── emitter.go                     # Generic type-safe emitter
├── publishEvent.go                # Event publishing (headers exchange)
├── publishAuditEvent.go           # Audit event publishing (direct exchange)
├── consumeCallbackEvent.go        # Event consumer (emits audit.received)
├── consumeChannelEvents.go        # Channel wrapper (emits audit.processed/dead_letter)
├── createAuditInfrastructure.go   # Audit exchange & queues setup
├── createHeaderConsumer.go        # Event subscription setup
├── createConsumers.go             # Saga command consumers
├── commenceSaga.go                # Saga initiation
├── consumeCallbackSagaCommand.go  # Saga step processing
└── resources.go                   # Exchange & queue constants
```

---

## Audit Feature

### Overview

The audit feature **automatically tracks the lifecycle of every event** without code changes. It emits three types of audit events:

1. **`audit.received`**: When event arrives (before processing)
2. **`audit.processed`**: When event succeeds (on ACK)
3. **`audit.dead_letter`**: When event fails (on NACK)

**Key Design**: Audit uses a **direct exchange** (not headers) for efficient single-consumer routing to the audit microservice.

### Automatic Emission Points

| Trigger         | Location                      | Audit Event Emitted      |
| --------------- | ----------------------------- | ------------------------ |
| Event received  | consumeCallbackEvent.go       | `audit.received`         |
| ACK called      | consumeChannelEvents.go:17    | `audit.processed`        |
| NACK called     | consumeChannelEvents.go:43    | `audit.dead_letter`      |
| Infrastructure  | createAuditInfrastructure.go  | Creates exchanges/queues |
| Publishing      | publishAuditEvent.go          | Sends to audit_exchange  |
| Auto-setup      | start.go:169 (ConnectToEvents)| Called automatically     |

### Audit Payloads

```go
// All audit payloads include:
type AuditReceivedPayload struct {
    Microservice  string  `json:"microservice"`   // Source microservice
    ReceivedEvent string  `json:"received_event"` // Event type
    ReceivedAt    uint64  `json:"received_at"`    // UNIX timestamp
    QueueName     string  `json:"queue_name"`     // Consumer queue
    EventID       *string `json:"event_id,omitempty"` // Optional correlation ID
}

// Similar: AuditProcessedPayload (processed_event, processed_at)
//          AuditDeadLetterPayload (rejected_event, rejection_reason, retry_count)
```

### Infrastructure (Auto-Created)

| Resource                     | Type            | Routing Key           |
| ---------------------------- | --------------- | --------------------- |
| `audit_exchange`             | Direct Exchange | N/A                   |
| `audit_received_commands`    | Queue           | `audit.received`      |
| `audit_processed_commands`   | Queue           | `audit.processed`     |
| `audit_dead_letter_commands` | Queue           | `audit.dead_letter`   |

All resources are **durable** and created automatically when `ConnectToEvents()` is called.

### Usage

**For Regular Microservices**: No changes needed! Audit is automatic.

```go
// This code automatically emits audit events:
emitter.On(event.AuthNewUserEvent, func(handler saga.EventHandler) {
    // audit.received emitted before this runs

    payload := saga.ParsePayload(handler.Payload, &event.AuthNewUserPayload{})
    createUser(payload)

    handler.Channel.AckMessage() // ← audit.processed emitted
})
```

**For Audit Microservice**: Consume from audit queues directly using standard RabbitMQ client (separate service, not using this library).

### Expected Behavior

- **Every event produces exactly 2 audit events**: `received` + (`processed` OR `dead_letter`)
- **Audit events do NOT audit themselves**: No recursive loops
- **Failure tracking**: `audit.dead_letter` includes error reason and retry count
- **Isolated streams**: Each microservice's events are tracked separately

---

## Event System

### Publishing Events

Events use a **headers exchange** (`matching_exchange`) for broadcast routing:

```go
saga.PublishEvent(&event.AuthNewUserPayload{
    ID:       "user-123",
    Email:    "user@example.com",
    Username: "johndoe",
})
```

**Headers created**: `{"AUTH.NEW_USER": "auth.new_user", "all-micro": "yes"}`

### Consuming Events

```go
transactional := saga.Config(&saga.Opts{
    RabbitUri:    "amqp://localhost:5672",
    Microservice: micro.Social,
    Events:       []event.MicroserviceEvent{event.AuthNewUserEvent}, // Subscribe
})

emitter := transactional.ConnectToEvents()
emitter.On(event.AuthNewUserEvent, func(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.AuthNewUserPayload{})

    // Process event
    createUserProfile(payload)

    handler.Channel.AckMessage() // ← Automatically emits audit.processed
})
```

### Error Handling with Fibonacci Backoff

```go
handler.Channel.NackWithFibonacciStrategy(19, 3)
// maxOccurrence: 19 (stops after 19 retries ~ 1.18 hours)
// maxNackRetries: 3 (requeues 3 times, then drops)
```

**Fibonacci Delays**: 1s, 2s, 3s, 5s, 8s, 13s, 21s, 34s, 55s, 89s, 144s...
**Max**: 25 occurrences (~20.84 hours total delay)

---

## Saga Orchestration

Sagas use a **direct exchange** (`commands_exchange`) with microservice-specific routing.

### Initiate Saga

```go
saga.CommenceSaga(&saga.PurchaseResourceFlowPayload{
    UserId:     "user-123",
    ResourceId: "resource-456",
    Price:      100,
    Quantity:   1,
})
```

### Handle Saga Step

```go
sagaEmitter := transactional.ConnectToSagaCommandEmitter()
sagaEmitter.On(micro.ResourcePurchasedDeductCoinsCommand, func(handler saga.CommandHandler) {
    payload := saga.ParsePayload(handler.Payload, &DeductCoinsPayload{})

    newBalance, err := deductCoins(payload.UserId, payload.Price)
    if err != nil {
        handler.Channel.NackWithFibonacciStrategy(19, 3)
        return
    }

    // Pass data to next saga step
    handler.Channel.AckMessage(saga.NextStepPayload{
        "__userId":    payload.UserId,
        "newBalance":  newBalance,
    })
})
```

---

## Developer Guidelines

### Adding New Events

1. **Define constant** (event/microserviceEvent.go):
   ```go
   const PaymentsChargeSucceededEvent MicroserviceEvent = "payments.charge_succeeded"
   ```

2. **Add to values array**:
   ```go
   func MicroserviceEventValues() []MicroserviceEvent {
       return []MicroserviceEvent{
           // ...
           PaymentsChargeSucceededEvent,
       }
   }
   ```

3. **Define payload struct**:
   ```go
   type PaymentsChargeSucceededPayload struct {
       ChargeID string  `json:"chargeId"`
       Amount   float64 `json:"amount"`
   }

   func (PaymentsChargeSucceededPayload) Type() MicroserviceEvent {
       return PaymentsChargeSucceededEvent
   }
   ```

4. **Publish and consume** as shown in Quick Reference

### Adding New Microservices

1. **Define constant** (micro/sagaCommands.go):
   ```go
   const Payments AvailableMicroservices = "payments"
   ```

2. **Add to values array** and **update IsValid()** method

3. Define saga commands if needed

### Testing

Integration tests in `test/`:
- `conf_test.go`: Configuration tests
- `pub_test.go`: Publishing tests

---

## Key Design Decisions

### Events vs Saga Commands

| Aspect       | Events (Headers Exchange)              | Saga Commands (Direct Exchange)        |
| ------------ | -------------------------------------- | -------------------------------------- |
| Pattern      | Pub/Sub (broadcast)                    | Point-to-point                         |
| Exchange     | `matching_exchange` (headers)          | `commands_exchange` (direct)           |
| Routing      | Headers: `{"EVENT.NAME": "event.name"}`| Routing key: `microservice_name`       |
| Use Case     | Business events (user created, etc.)   | Orchestrated saga steps                |
| Consumers    | Multiple microservices                 | Single microservice                    |

### Audit Exchange Choice

**Direct Exchange** chosen for audit because:
- Single consumer (audit microservice)
- No broadcast needed
- More efficient than headers exchange
- Routing keys: `audit.received`, `audit.processed`, `audit.dead_letter`

---

## Summary

### Production-Ready Features
- ✅ Event-driven architecture with headers-based routing
- ✅ Saga orchestration for distributed transactions
- ✅ Automatic audit logging (zero-code overhead)
- ✅ Fibonacci backoff retry strategy
- ✅ Type-safe generics for handlers
- ✅ Durable messaging with QoS=1
- ✅ Health checks for RabbitMQ connectivity

### Statistics
- **Events**: 45 total (42 business + 3 audit)
- **Microservices**: 17 supported
- **Exchanges**: 3 (matching, commands, audit)
- **LOC**: ~2,100 (core + events + microservices)

### Critical Concepts
1. **Audit is automatic** - No code changes needed in microservices
2. **Two exchange types** - Headers for events, direct for commands/audit
3. **Fibonacci backoff** - Prevents thundering herd
4. **Type safety** - Generics ensure compile-time correctness
5. **Lifecycle tracking** - Every event produces 2 audit events

---

**For detailed implementation examples**, see existing code in `test/` directory.
**For questions or issues**, see [GitHub Issues](https://github.com/legendaryum-metaverse/saga/issues).