# Timestamp Precision Fix - Implementation Guide for All Libraries

> **Cross-library implementation guide for fixing timestamp precision from seconds to milliseconds in audit events**

---

## Executive Summary

**Problem**: All `legend-saga` libraries (Rust, TypeScript, Go) currently send audit event timestamps as **UNIX timestamps in seconds**, but consuming microservices (specifically `be-audit-eda`) are treating them as **milliseconds**, resulting in incorrect absolute timestamp storage.

**Solution**: Update all libraries to send timestamps as **milliseconds** instead of seconds.

**Impact**: This is ~~~~a **breaking change** that requires coordinated deployment across all libraries.

---

## Background

### The Issue

The `be-audit-eda` microservice uses this code to parse timestamps:

```rust
// Consuming microservice expects milliseconds
let published_at = chrono::DateTime::from_timestamp_millis(payload.published_at as i64)
    .unwrap_or(now);
```

However, all libraries are currently sending timestamps in **seconds**, leading to:

- ❌ Dates showing as **1970** instead of **2025** in the database
- ❌ Loss of sub-second precision (critical for 0.5ms - 2ms latency tracking)
- ❌ Broken cross-system time comparisons

### Why Milliseconds?

1. **Modern standard**: Most APIs and languages use milliseconds (JavaScript, Java, etc.)
2. **Sub-second precision**: Required for accurate performance monitoring
3. **Consumer expectations**: `be-audit-eda` already expects milliseconds
4. **Database format**: ClickHouse uses `DateTime64(3)` (millisecond precision)

---

## Implementation Checklist

### Phase 1: Identify Affected Code

All libraries have **4 audit event payload types** that need updating:

1. **`AuditPublishedPayload`** - Tracks when event is published at source
2. **`AuditReceivedPayload`** - Tracks when event is received before processing
3. **`AuditProcessedPayload`** - Tracks successful event processing
4. **`AuditDeadLetterPayload`** - Tracks when message is rejected/nacked

Each payload has a timestamp field that needs to be changed from seconds → milliseconds.

### Phase 2: Update Documentation

**Action**: Update field documentation for all timestamp fields.

**Change**:

- ❌ OLD: "UNIX timestamp" or "UNIX timestamp in seconds"
- ✅ NEW: "UNIX timestamp in milliseconds"

**Example (Rust)**:

```rust
// BEFORE
pub struct AuditPublishedPayload {
    pub publisher_microservice: String,
    pub published_event: String,
    /// Timestamp when the event was published (UNIX timestamp in seconds) // ❌
    pub published_at: u64,
    pub event_id: String,
}

// AFTER
pub struct AuditPublishedPayload {
    pub publisher_microservice: String,
    pub published_event: String,
    /// Timestamp when the event was published (UNIX timestamp in milliseconds) // ✅
    pub published_at: u64,
    pub event_id: String,
}
```

### Phase 3: Update Timestamp Generation

**Action**: Find all locations where timestamps are generated and change from seconds to milliseconds.

#### Pattern to Search For

Look for code that:

1. Gets current system time
2. Converts to UNIX timestamp
3. Uses **seconds** (not milliseconds)
4. Assigns to audit payload timestamp fields

#### Language-Specific Examples

**Rust** (6 locations in legend-saga):

```rust
// ❌ OLD: Generates seconds
let timestamp = SystemTime::now()
    .duration_since(UNIX_EPOCH)
    .unwrap_or_default()
    .as_secs();  // Returns seconds

// ✅ NEW: Generates milliseconds
let timestamp = SystemTime::now()
    .duration_since(UNIX_EPOCH)
    .unwrap_or_default()
    .as_millis() as u64;  // Returns milliseconds
```

**TypeScript/JavaScript**:

```typescript
// ❌ OLD: Generates seconds
const timestamp = Math.floor(Date.now() / 1000); // Divides by 1000 = seconds

// ✅ NEW: Generates milliseconds
const timestamp = Date.now(); // Already in milliseconds!
```

**Go**:

```go
// ❌ OLD: Generates seconds
timestamp := time.Now().Unix()  // Returns seconds

// ✅ NEW: Generates milliseconds
timestamp := time.Now().UnixMilli()  // Returns milliseconds
// OR (for Go < 1.17)
timestamp := time.Now().UnixNano() / 1e6  // Convert nanoseconds to milliseconds
```

#### Where to Look

Search your codebase for these patterns:

1. **Event publishing functions** - Where `AuditPublishedPayload` is created
2. **Event receiving handlers** - Where `AuditReceivedPayload` is created
3. **Event processing (ACK)** - Where `AuditProcessedPayload` is created
4. **Event rejection (NACK)** - Where `AuditDeadLetterPayload` is created
5. **Test code** - Any tests that create audit payloads

**Rust locations** (as reference):

- `src/publish_event.rs` - Publishing audit.published
- `src/events_consume.rs` (4 locations):
  - `ack()` method - audit.processed
  - `nack_with_delay()` method - audit.dead_letter
  - `nack_with_fibonacci_strategy()` method - audit.dead_letter
  - Event consumption handler - audit.received
- `src/start.rs` - Test code

### Phase 4: Add Unit Tests

**Action**: Add tests to verify timestamp precision.

**What to test**:

1. Timestamp values are in reasonable range (year 2020-2030)
2. Timestamps have millisecond precision (not just seconds)
3. Verify milliseconds are ~1000x larger than seconds

**Example (Rust)**:

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use std::time::{SystemTime, UNIX_EPOCH};

    // Year 2020 in milliseconds (Jan 1, 2020 00:00:00 UTC)
    const YEAR_2020_MS: u64 = 1577836800000;
    // Year 2030 in milliseconds (Jan 1, 2030 00:00:00 UTC)
    const YEAR_2030_MS: u64 = 1893456000000;

    #[test]
    fn test_audit_published_payload_timestamp_precision() {
        let current_ms = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_millis() as u64;

        let payload = AuditPublishedPayload {
            publisher_microservice: "test-service".to_string(),
            published_event: "test.event".to_string(),
            published_at: current_ms,
            event_id: "test-uuid".to_string(),
        };

        // Verify timestamp is in reasonable range (year 2020-2030)
        assert!(
            payload.published_at > YEAR_2020_MS,
            "Timestamp {} should be after year 2020",
            payload.published_at
        );
        assert!(
            payload.published_at < YEAR_2030_MS,
            "Timestamp {} should be before year 2030",
            payload.published_at
        );

        // Verify millisecond precision (should have 3+ more digits than seconds)
        let as_seconds = payload.published_at / 1000;
        assert!(
            payload.published_at > as_seconds * 1000,
            "Timestamp should have millisecond precision"
        );
    }

    #[test]
    fn test_millisecond_vs_second_timestamps() {
        // Generate timestamp in milliseconds
        let timestamp_ms = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_millis() as u64;

        // Generate timestamp in seconds (old way)
        let timestamp_s = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();

        // Milliseconds should be ~1000x larger than seconds
        assert!(
            timestamp_ms > timestamp_s * 100,
            "Millisecond timestamp should be much larger than second timestamp"
        );

        // If treated as milliseconds, seconds would show as 1970
        assert!(
            timestamp_s < YEAR_2020_MS,
            "Second timestamp would be before year 2000 if treated as milliseconds"
        );
    }
}
```

**TypeScript/JavaScript example**:

```typescript
describe('Timestamp Precision Tests', () => {
  const YEAR_2020_MS = 1577836800000;
  const YEAR_2030_MS = 1893456000000;

  it('should generate timestamps in milliseconds', () => {
    const timestamp = Date.now();

    // Verify timestamp is in reasonable range
    expect(timestamp).toBeGreaterThan(YEAR_2020_MS);
    expect(timestamp).toBeLessThan(YEAR_2030_MS);

    // Verify millisecond precision
    const asSeconds = Math.floor(timestamp / 1000);
    expect(timestamp).toBeGreaterThan(asSeconds * 1000);
  });

  it('should create audit payload with millisecond timestamp', () => {
    const payload = {
      publisher_microservice: 'test-service',
      published_event: 'test.event',
      published_at: Date.now(),
      event_id: 'test-uuid',
    };

    expect(payload.published_at).toBeGreaterThan(YEAR_2020_MS);
    expect(payload.published_at).toBeLessThan(YEAR_2030_MS);
  });
});
```

**Go example**:

```go
func TestTimestampPrecision(t *testing.T) {
    const YEAR_2020_MS int64 = 1577836800000
    const YEAR_2030_MS int64 = 1893456000000

    t.Run("timestamps should be in milliseconds", func(t *testing.T) {
        timestamp := time.Now().UnixMilli()

        // Verify timestamp is in reasonable range
        if timestamp < YEAR_2020_MS {
            t.Errorf("Timestamp %d should be after year 2020", timestamp)
        }
        if timestamp > YEAR_2030_MS {
            t.Errorf("Timestamp %d should be before year 2030", timestamp)
        }

        // Verify millisecond precision
        asSeconds := timestamp / 1000
        if timestamp <= asSeconds * 1000 {
            t.Error("Timestamp should have millisecond precision")
        }
    })

    t.Run("should create audit payload with millisecond timestamp", func(t *testing.T) {
        payload := AuditPublishedPayload{
            PublisherMicroservice: "test-service",
            PublishedEvent:        "test.event",
            PublishedAt:           time.Now().UnixMilli(),
            EventID:               "test-uuid",
        }

        if payload.PublishedAt < YEAR_2020_MS {
            t.Errorf("Timestamp %d should be after year 2020", payload.PublishedAt)
        }
        if payload.PublishedAt > YEAR_2030_MS {
            t.Errorf("Timestamp %d should be before year 2030", payload.PublishedAt)
        }
    })
}
```

---

## Verification After Implementation

### 1. Run All Tests

Ensure all existing tests pass **AND** new timestamp precision tests pass.

**Rust**:

```bash
make test
# OR
cargo test --lib -- --test-threads=1
```

**TypeScript**:

```bash
npm test
# OR
yarn test
```

**Go**:

```bash
go test ./... -v
```

### 2. Manual Verification

Create a simple test that prints timestamps to verify format:

```rust
// Rust
println!("Timestamp (ms): {}", SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis());
// Should print: Timestamp (ms): 1761368484328 (example)
```

```typescript
// TypeScript
console.log('Timestamp (ms):', Date.now());
// Should print: Timestamp (ms): 1761368484328 (example)
```

```go
// Go
fmt.Printf("Timestamp (ms): %d\n", time.Now().UnixMilli())
// Should print: Timestamp (ms): 1761368484328 (example)
```

**Expected**: Value should be ~1.7 trillion (13 digits) for year 2025.
**Wrong**: Value would be ~1.7 billion (10 digits) if still in seconds.

### 3. Database Verification (After Deployment)

Run this query on the ClickHouse database to verify correct timestamps:

```sql
SELECT
    event_id,
    published_at,
    toUnixTimestamp64Milli(published_at) as stored_ms,
    now64(3) as current_time
FROM audit_published
ORDER BY timestamp DESC
LIMIT 5;
```

**Expected**: `published_at` should show dates in 2025, not 1970.

---

## Common Pitfalls

### ❌ Don't: Mix Seconds and Milliseconds

```typescript
// ❌ WRONG: This creates seconds then treats as milliseconds
const timestampSeconds = Math.floor(Date.now() / 1000);
const payload = {
  published_at: timestampSeconds, // Will show as 1970!
};

// ✅ CORRECT: Use milliseconds directly
const payload = {
  published_at: Date.now(),
};
```

### ❌ Don't: Forget Test Code

```rust
// ❌ WRONG: Test still uses seconds
#[test]
fn test_something() {
    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs();  // Old code in test!
}

// ✅ CORRECT: Update tests too
#[test]
fn test_something() {
    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis() as u64;
}
```

### ❌ Don't: Assume Type Changes Are Needed

The data type remains the same (usually `u64`, `number`, `int64`). Only the **value** changes:

- ❌ OLD: `1761368484` (10 digits, seconds)
- ✅ NEW: `1761368484328` (13 digits, milliseconds)

---

## Language-Specific Quick Reference

### Rust

| Operation     | Seconds (OLD ❌) | Milliseconds (NEW ✅) |
| ------------- | ---------------- | --------------------- |
| Get timestamp | `.as_secs()`     | `.as_millis() as u64` |
| Type          | `u64`            | `u64`                 |
| Example value | `1761368484`     | `1761368484328`       |

### TypeScript/JavaScript

| Operation     | Seconds (OLD ❌)                | Milliseconds (NEW ✅) |
| ------------- | ------------------------------- | --------------------- |
| Get timestamp | `Math.floor(Date.now() / 1000)` | `Date.now()`          |
| Type          | `number`                        | `number`              |
| Example value | `1761368484`                    | `1761368484328`       |

### Go

| Operation     | Seconds (OLD ❌)    | Milliseconds (NEW ✅)    |
| ------------- | ------------------- | ------------------------ |
| Get timestamp | `time.Now().Unix()` | `time.Now().UnixMilli()` |
| Type          | `int64`             | `int64`                  |
| Example value | `1761368484`        | `1761368484328`          |

---

## Deployment Coordination

### Breaking Change

This is a **breaking change** because:

- Old libraries send seconds (10 digits)
- New libraries send milliseconds (13 digits)
- Consumer expects one format consistently

### Deployment Strategy

**Option 1: Big Bang Deployment** (Recommended if automated)

1. Deploy all libraries simultaneously via CI/CD
2. All microservices update at once
3. Clean cutover with no mixed state

**Option 2: Gradual Deployment** (If needed)

1. Deploy consumer validation first (optional, not in this guide)
2. Deploy libraries one by one
3. Monitor logs for issues

### Success Criteria

After deployment, verify:

- ✅ Database queries show dates in 2025 (not 1970)
- ✅ Latency calculations remain accurate with millisecond precision
- ✅ All library tests pass
- ✅ No errors in microservice logs
- ✅ Sub-second latency metrics preserved (0.5ms, 1.45ms, etc.)

---

## Summary

**Changes Required**:

1. ✅ Update documentation: "UNIX timestamp in seconds" → "UNIX timestamp in milliseconds"
2. ✅ Update timestamp generation: seconds → milliseconds
3. ✅ Add unit tests: verify millisecond precision
4. ✅ Run tests: ensure everything passes

**Files to Search**:

- Event payload definitions (4 structs)
- Event publishing code
- Event consuming code (ACK/NACK handlers)
- Test files

**Key Insight**: The fix is simple but requires systematic changes across all timestamp generation points. Use search functionality to find all occurrences of second-based timestamp generation and update them to milliseconds.

---

## Questions?

If you encounter issues:

1. Check that ALL timestamp generation uses milliseconds
2. Verify tests include reasonable range checks (2020-2030)
3. Ensure no test code was missed
4. Confirm timestamp values are ~13 digits (not 10)

**Reference Implementation**: See the Rust library (`legend-saga` v0.1.0+) for a complete working example of this fix.
