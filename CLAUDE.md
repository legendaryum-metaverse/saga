# 📖 CLAUDE.md - legend-saga Golang Library

> **Technical Documentation for Claude Code Assistant**
> Updated: 2025-10-05 | Go 1.25.0 | Production-Ready

---

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [Codebase Structure](#codebase-structure)
3. [Core Architecture](#core-architecture)
4. [Audit Feature (NEW)](#audit-feature)
5. [Event System](#event-system)
6. [Saga Orchestration](#saga-orchestration)
7. [Developer Guidelines](#developer-guidelines)
8. [Code Metrics](#code-metrics)

---

## Project Overview

**legend-saga** is a Go library for microservice communication via RabbitMQ, supporting event-driven architecture and saga patterns.

### Key Features

- 🔄 Publish/Subscribe with headers-based routing
- 🎭 Saga orchestration for distributed transactions
- ♻️ Fibonacci backoff retry mechanism
- 🛡️ Durable message delivery
- 🎯 Type-safe Go generics
- **🔍 Automatic audit logging (NEW)**

### Tech Stack

- **Language**: Go 1.25.0
- **Broker**: RabbitMQ
- **AMQP**: `rabbitmq/amqp091-go v1.10.0`
- **Validation**: `go-playground/validator/v10`

---

## Codebase Structure

### Directory Layout

```
legend-saga/
├── event/microserviceEvent.go          # Event types & payloads (590+ LOC)
├── micro/sagaCommands.go                # Microservices & commands (106 LOC)
├── test/                                # Integration tests
│
# Core files (flat structure ~1,400 LOC)
├── start.go                             # Entry point, Transactional, ConnectToAudit
├── emitter.go                           # Generic event emitter
├── publishEvent.go                      # Event publishing
├── publishAuditEvent.go                 # Audit event publishing (NEW)
├── consumeCallbackEvent.go              # Event handler (emits audit.received)
├── consumeCallbackAudit.go              # AuditHandler (non-recursive) (NEW)
├── consumeChannelEvents.go              # Event channel (emits audit.processed/dead_letter)
├── createAuditInfrastructure.go         # Audit exchange & queues setup (NEW)
├── resources.go                         # Exchange & queue constants
└── ...
```

### Package Structure

- **`saga`** (root): Core functionality
- **`event`**: Event definitions, payloads, **audit payloads**
- **`micro`**: Microservice identifiers, **AuditEda**

---

## Core Architecture

### High-Level Flow

```
Application → Transactional → [ConnectToEvents | ConnectToSagaCommandEmitter | ConnectToAudit]
                                      ↓                    ↓                           ↓
                              Event Emitter        Saga Emitter                 Audit Emitter
                                      ↓                    ↓                           ↓
                             RabbitMQ Exchanges   RabbitMQ Exchanges           RabbitMQ Audit Exchange
                                      ↓                    ↓                           ↓
                              Event Consumers      Saga Consumers               Audit Consumers
                                      ↓                    ↓                           ↓
                              Business Logic        Saga Steps                  Audit Processing
```

### Key Components

| Component              | File                       | Purpose                               |
| ---------------------- | -------------------------- | ------------------------------------- |
| `Transactional`        | start.go:46                | Main coordinator, connection manager  |
| `Emitter[T, U]`        | emitter.go:7               | Type-safe event dispatcher (generics) |
| `EventHandler`         | consumeCallbackEvent.go:13 | Event processing handler              |
| `AuditHandler`         | consumeCallbackAudit.go:13 | **Non-recursive audit handler (NEW)** |
| `EventsConsumeChannel` | consumeChannelEvents.go:11 | **Auto-emits audit events (NEW)**     |

---

## Audit Feature

### Overview

The audit feature automatically tracks the **lifecycle of every event** for monitoring, debugging, and compliance. It emits three types of audit events:

1. **`audit.received`**: When an event is received (before processing)
2. **`audit.processed`**: When an event is successfully processed (on ACK)
3. **`audit.dead_letter`**: When an event fails and is NACKed

### Architecture

#### Direct Exchange Pattern

```
All Microservices           Audit Microservice
       ↓                            ↓
[Event Processing]      [Audit Event Consumers]
       ↓                            ↓
Emit Audit Events → audit_exchange (direct) → 3 Queues
                           ↓
                    ┌──────┴──────┬──────────┐
                    ↓             ↓          ↓
        audit_received_commands  audit_processed_commands  audit_dead_letter_commands
```

**Why Direct Exchange?**

- Efficient single-consumer delivery to audit microservice
- Routing keys: `audit.received`, `audit.processed`, `audit.dead_letter`
- No broadcast overhead (unlike headers exchange for events)

#### Non-Recursive Design

**Critical**: Audit events themselves must NOT emit more audits (infinite loop prevention).

```go
// Standard EventHandler (emits audits)
eventCallback() → emits audit.received → processes event → handler.AckMessage() → emits audit.processed

// AuditHandler (NO recursive audits)
auditEventCallback() → processes audit event → handler.AuditAck() → NO audit emission
```

### Implementation Details

#### 1. Event Definitions (event/microserviceEvent.go)

```go
const (
    AuditReceivedEvent   MicroserviceEvent = "audit.received"
    AuditProcessedEvent  MicroserviceEvent = "audit.processed"
    AuditDeadLetterEvent MicroserviceEvent = "audit.dead_letter"
)

type AuditReceivedPayload struct {
    Microservice  string  `json:"microservice"`
    ReceivedEvent string  `json:"received_event"`
    ReceivedAt    uint64  `json:"received_at"`    // UNIX timestamp
    QueueName     string  `json:"queue_name"`
    EventID       *string `json:"event_id,omitempty"`
}
// Similar: AuditProcessedPayload, AuditDeadLetterPayload
```

#### 2. Infrastructure (createAuditInfrastructure.go)

```go
func (t *Transactional) createAuditLoggingResources() error {
    // 1. Declare direct exchange
    t.eventsChannel.ExchangeDeclare("audit_exchange", "direct", ...)

    // 2. Declare 3 separate queues
    t.eventsChannel.QueueDeclare("audit_received_commands", ...)
    t.eventsChannel.QueueDeclare("audit_processed_commands", ...)
    t.eventsChannel.QueueDeclare("audit_dead_letter_commands", ...)

    // 3. Bind queues to exchange with routing keys
    t.eventsChannel.QueueBind("audit_received_commands", "audit.received", "audit_exchange", ...)
    // ... similar for other queues
}
```

#### 3. Automatic Emission (consumeChannelEvents.go)

**On ACK** (consumeChannelEvents.go:17):

```go
func (m *EventsConsumeChannel) AckMessage() {
    m.channel.Ack(m.msg.DeliveryTag, false)

    // Auto-emit audit.processed
    auditPayload := &event.AuditProcessedPayload{
        Microservice:   m.microservice,
        ProcessedEvent: m.eventType,
        ProcessedAt:    uint64(time.Now().Unix()),
        QueueName:      m.queueName,
    }
    PublishAuditProcessedEvent(auditPayload)
}
```

**On NACK** (consumeChannelEvents.go:43):

```go
func (m *EventsConsumeChannel) NackWithFibonacciStrategy(...) {
    m.ConsumeChannel.NackWithFibonacciStrategy(...)

    // Auto-emit audit.dead_letter
    auditPayload := &event.AuditDeadLetterPayload{
        Microservice:    m.microservice,
        RejectedEvent:   m.eventType,
        RejectionReason: "fibonacci_strategy",
        RetryCount:      &retryCount,
        ...
    }
    PublishAuditDeadLetterEvent(auditPayload)
}
```

**On Receive** (consumeCallbackEvent.go:87):

```go
func eventCallback(..., microservice string) {
    // Parse event from message
    eventType := string(eventKey[0])

    // Auto-emit audit.received BEFORE processing
    auditPayload := &event.AuditReceivedPayload{
        Microservice:  microservice,
        ReceivedEvent: eventType,
        ReceivedAt:    uint64(time.Now().Unix()),
        QueueName:     queueName,
    }
    PublishAuditReceivedEvent(auditPayload)

    // Then process event normally
    emitter.Emit(eventKey[0], EventHandler{...})
}
```

#### 4. Audit Consumer (start.go:181)

```go
func (t *Transactional) ConnectToAudit() *Emitter[AuditHandler, event.MicroserviceEvent] {
    // Create audit infrastructure
    t.createAuditLoggingResources()

    e := newEmitter[AuditHandler, event.MicroserviceEvent]()

    // Spawn 3 consumers (one per audit event type)
    go consume(e, "audit_received_commands", t.eventsChannel, auditEventCallback, ...)
    go consume(e, "audit_processed_commands", t.eventsChannel, auditEventCallback, ...)
    go consume(e, "audit_dead_letter_commands", t.eventsChannel, auditEventCallback, ...)

    return e
}
```

#### 5. Non-Recursive Handler (consumeCallbackAudit.go)

```go
type AuditHandler struct {
    Channel *AuditConsumeChannel
    Payload map[string]interface{}
}

type AuditConsumeChannel struct {
    *ConsumeChannel
}

// AuditAck does NOT emit audit events (breaks recursion)
func (a *AuditConsumeChannel) AuditAck() {
    a.channel.Ack(a.msg.DeliveryTag, false)
    // NO PublishAuditProcessedEvent() call here!
}

// auditEventCallback does NOT emit audit.received (breaks recursion)
func auditEventCallback(...) {
    // Determine event type from queue name (not headers)
    var eventType event.MicroserviceEvent
    switch queueName {
    case "audit_received_commands":
        eventType = event.AuditReceivedEvent
    // ...
    }

    // NO PublishAuditReceivedEvent() call here!
    emitter.Emit(eventType, AuditHandler{...})
}
```

### Usage

#### For Regular Microservices (Automatic)

No code changes needed! Audit events are emitted automatically:

```go
// Existing code works as-is
transactional := saga.Config(&saga.Opts{
    RabbitUri:    "amqp://...",
    Microservice: micro.Social,
    Events:       []event.MicroserviceEvent{event.AuthNewUserEvent},
})

emitter := transactional.ConnectToEvents()
emitter.On(event.AuthNewUserEvent, func(handler saga.EventHandler) {
    // Process event
    // ...
    handler.Channel.AckMessage() // ← Automatically emits audit.processed
})
```

**Audit events emitted**:

1. `audit.received` when event arrives (before handler)
2. `audit.processed` when `AckMessage()` called
3. `audit.dead_letter` if `NackWithFibonacciStrategy()` called

#### For Audit Microservice

```go
transactional := saga.Config(&saga.Opts{
    RabbitUri:    "amqp://...",
    Microservice: micro.AuditEda, // Special audit microservice
    Events:       []event.MicroserviceEvent{},
})

// Use ConnectToAudit() instead of ConnectToEvents()
auditEmitter := transactional.ConnectToAudit()

// Handle audit events without recursion
auditEmitter.On(event.AuditReceivedEvent, func(handler saga.AuditHandler) {
    payload := ParseAuditPayload(handler.Payload, &event.AuditReceivedPayload{})

    // Store in database, send to monitoring, etc.
    storeAuditEvent(payload)

    // Use AuditAck() to avoid recursive audit emission
    handler.Channel.AuditAck()
})

auditEmitter.On(event.AuditProcessedEvent, func(handler saga.AuditHandler) {
    payload := ParseAuditPayload(handler.Payload, &event.AuditProcessedPayload{})
    storeAuditEvent(payload)
    handler.Channel.AuditAck()
})

auditEmitter.On(event.AuditDeadLetterEvent, func(handler saga.AuditHandler) {
    payload := ParseAuditPayload(handler.Payload, &event.AuditDeadLetterPayload{})
    alertOnFailure(payload)
    handler.Channel.AuditAck()
})
```

### Expected Behavior (from Rust specification)

1. **Every event produces exactly 2 audit events**:

   - `audit.received` (when received)
   - `audit.processed` OR `audit.dead_letter` (never both)

2. **Audit events do NOT audit themselves**:

   - Processing `audit.received` does NOT emit another `audit.received`
   - `AuditHandler` prevents recursion

3. **Failure scenarios**:

   - Failed events emit `audit.dead_letter` instead of `audit.processed`
   - Error reason and retry count included in payload

4. **Isolated streams**:
   - Each microservice has separate audit stream
   - Audit microservice processes all audit events

### Resources

| Resource                     | Type            | Purpose                        |
| ---------------------------- | --------------- | ------------------------------ |
| `audit_exchange`             | Direct Exchange | Routes audit events            |
| `audit_received_commands`    | Queue           | Holds audit.received events    |
| `audit_processed_commands`   | Queue           | Holds audit.processed events   |
| `audit_dead_letter_commands` | Queue           | Holds audit.dead_letter events |

**Constants** (resources.go):

```go
const (
    AuditExchange           Exchange = "audit_exchange"
    AuditReceivedCommandsQ  Queue    = "audit_received_commands"
    AuditProcessedCommandsQ Queue    = "audit_processed_commands"
    AuditDeadLetterCommandsQ Queue   = "audit_dead_letter_commands"
)
```

### Key Rust→Go Patterns

| Rust Pattern                | Go Equivalent            | Notes                       |
| --------------------------- | ------------------------ | --------------------------- |
| `Arc<Mutex<Channel>>`       | `*amqp.Channel`          | Shared channel via pointer  |
| `async fn` + `tokio::spawn` | `go func()`              | Goroutine per consumer      |
| `Result<T, Error>`          | `(T, error)`             | Explicit error returns      |
| Trait `PayloadEvent`        | Interface `PayloadEvent` | `Type() MicroserviceEvent`  |
| `u64` timestamp             | `uint64`                 | UNIX timestamp in seconds   |
| `Option<String>`            | `*string`                | Pointer for optional fields |

---

## Event System

### Event Publishing

```go
saga.PublishEvent(&event.AuthNewUserPayload{
    ID:       "user-123",
    Email:    "user@example.com",
    Username: "johndoe",
})
```

**Flow**:

1. `PublishEvent()` → `getSendChannel()` → creates/reuses channel
2. Headers created: `{"AUTH.NEW_USER": "auth.new_user", "all-micro": "yes"}`
3. Publishes to `matching_exchange` (headers exchange)
4. RabbitMQ routes to subscribed microservices

### Event Consumption

```go
emitter := transactional.ConnectToEvents()
emitter.On(event.AuthNewUserEvent, func(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.AuthNewUserPayload{})

    createUserProfile(payload.ID, payload.Email)

    handler.Channel.AckMessage() // ← Emits audit.processed
})
```

**Flow**:

1. Consumer goroutine receives message
2. **Emits `audit.received`** (NEW)
3. `eventCallback()` extracts event type from headers
4. Parses JSON payload
5. Emits to registered handler
6. Handler processes and ACKs/NACKs
7. **Emits `audit.processed` or `audit.dead_letter`** (NEW)

### Retry Mechanism

```go
handler.Channel.NackWithFibonacciStrategy(MAX_OCCURRENCE, MAX_NACK_RETRIES)
```

**Fibonacci Backoff**:

- Occurrence: 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144...
- Delay (seconds): 1s, 2s, 3s, 5s, 8s, 13s, 21s, 34s, 55s, 89s, 144s...
- After 19 occurrences: ~1.18 hours delay
- Max: 25 occurrences (~20.84 hours)

---

## Saga Orchestration

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
emitter := transactional.ConnectToSagaCommandEmitter()
emitter.On(micro.ResourcePurchasedDeductCoinsCommand, func(handler saga.CommandHandler) {
    payload := saga.ParsePayload(handler.Payload, &DeductCoinsPayload{})

    newBalance, err := deductCoins(payload.UserId, payload.Price)
    if err != nil {
        handler.Channel.NackWithFibonacciStrategy(19, 3)
        return
    }

    handler.Channel.AckMessage(saga.NextStepPayload{
        "__userId":    payload.UserId,
        "__price":     payload.Price,
        "newBalance":  newBalance,
    })
})
```

---

## Developer Guidelines

### Adding Audit Support to Existing Code

**No changes needed!** Audit is automatic. Just ensure:

1. Use `handler.Channel.AckMessage()` for success
2. Use `handler.Channel.NackWithFibonacciStrategy()` for failure
3. Audit events will be emitted automatically

### Creating Audit Microservice

```go
// 1. Configure with AuditEda microservice
transactional := saga.Config(&saga.Opts{
    RabbitUri:    "amqp://...",
    Microservice: micro.AuditEda,
    Events:       []event.MicroserviceEvent{},
})

// 2. Use ConnectToAudit() instead of ConnectToEvents()
auditEmitter := transactional.ConnectToAudit()

// 3. Handle audit events with AuditHandler
auditEmitter.On(event.AuditReceivedEvent, func(handler saga.AuditHandler) {
    // Process audit event
    handler.Channel.AuditAck() // ← Use AuditAck(), not AckMessage()
})
```

### Adding New Events

1. **Define event** (event/microserviceEvent.go):

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

3. **Define payload**:

```go
type PaymentsChargeSucceededPayload struct {
    ChargeID string  `json:"chargeId"`
    Amount   float64 `json:"amount"`
}
func (PaymentsChargeSucceededPayload) Type() MicroserviceEvent {
    return PaymentsChargeSucceededEvent
}
```

4. **Publish**:

```go
saga.PublishEvent(&event.PaymentsChargeSucceededPayload{
    ChargeID: "ch_123",
    Amount:   99.99,
})
```

5. **Subscribe**:

```go
emitter.On(event.PaymentsChargeSucceededEvent, func(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.PaymentsChargeSucceededPayload{})
    // Process
    handler.Channel.AckMessage()
})
```

---

## Code Metrics

### Lines of Code

| Category                 | LOC        | Notes                      |
| ------------------------ | ---------- | -------------------------- |
| Event Definitions        | 590        | +70 LOC for audit payloads |
| Core Logic               | 1,400      | +260 LOC for audit feature |
| Microservice Definitions | 106        | +4 LOC for AuditEda        |
| **Total**                | **~2,100** | **+334 LOC for audit**     |

### New Files (Audit Feature)

| File                           | LOC | Purpose                     |
| ------------------------------ | --- | --------------------------- |
| `publishAuditEvent.go`         | 76  | Audit event publishing      |
| `consumeCallbackAudit.go`      | 94  | Non-recursive audit handler |
| `createAuditInfrastructure.go` | 100 | Audit infrastructure setup  |

### Modified Files (Audit Feature)

| File                         | Changes | Purpose                          |
| ---------------------------- | ------- | -------------------------------- |
| `event/microserviceEvent.go` | +70 LOC | Audit event types & payloads     |
| `consumeChannelEvents.go`    | +75 LOC | Auto-emit audit on ACK/NACK      |
| `consumeCallbackEvent.go`    | +25 LOC | Auto-emit audit.received         |
| `start.go`                   | +55 LOC | ConnectToAudit() method          |
| `resources.go`               | +4 LOC  | Audit exchange & queue constants |
| `micro/sagaCommands.go`      | +4 LOC  | AuditEda microservice            |

### Supported Events

- **Standard Events**: 42 (auth, coins, missions, rankings, rooms, social, etc.)
- **Audit Events**: 3 (audit.received, audit.processed, audit.dead_letter)
- **Total**: 45 events

### Supported Microservices

- **Business Microservices**: 16 (auth, blockchain, coins, etc.)
- **Audit Microservice**: 1 (audit-eda)
- **Total**: 17 microservices

---

## Summary

### What's New in Audit Feature

1. **Automatic Lifecycle Tracking**: Every event emits 2 audit events automatically
2. **Non-Recursive Design**: `AuditHandler` prevents infinite audit loops
3. **Direct Exchange**: Efficient routing to single audit microservice
4. **Transparent Integration**: No changes needed in existing code
5. **Rust Parity**: 100% functional parity with production Rust implementation

### Key Takeaways

- **For Regular Microservices**: Audit is automatic, zero code changes required
- **For Audit Microservice**: Use `ConnectToAudit()` and `AuditHandler`
- **Payload Format**: JSON with snake_case, UNIX timestamps (uint64)
- **Event Flow**: received → [processing] → processed OR dead_letter
- **Infrastructure**: 1 direct exchange, 3 queues, 3 routing keys

---

**End of Documentation**
For questions, see [GitHub Issues](https://github.com/legendaryum-metaverse/saga/issues).
