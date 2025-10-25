package test

import (
	"testing"
	"time"

	"github.com/legendaryum-metaverse/saga/event"
)

// Constants for timestamp validation
// Year 2020 in milliseconds (Jan 1, 2020 00:00:00 UTC)
const YEAR_2020_MS int64 = 1577836800000

// Year 2030 in milliseconds (Jan 1, 2030 00:00:00 UTC)
const YEAR_2030_MS int64 = 1893456000000

func TestTimestampPrecision(t *testing.T) {
	t.Run("timestamps should be in milliseconds", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()

		// Verify timestamp is in reasonable range (year 2020-2030)
		if timestamp < YEAR_2020_MS {
			t.Errorf("Timestamp %d should be after year 2020 (expected > %d)", timestamp, YEAR_2020_MS)
		}
		if timestamp > YEAR_2030_MS {
			t.Errorf("Timestamp %d should be before year 2030 (expected < %d)", timestamp, YEAR_2030_MS)
		}

		// Verify millisecond precision (should have 3+ more digits than seconds)
		asSeconds := timestamp / 1000
		if timestamp <= asSeconds*1000 {
			t.Error("Timestamp should have millisecond precision")
		}
	})

	t.Run("should create AuditPublishedPayload with millisecond timestamp", func(t *testing.T) {
		payload := event.AuditPublishedPayload{
			PublisherMicroservice: "test-service",
			PublishedEvent:        "test.event",
			PublishedAt:           uint64(time.Now().UnixMilli()),
			EventID:               "test-uuid",
		}

		// Verify timestamp is in reasonable range
		if int64(payload.PublishedAt) < YEAR_2020_MS {
			t.Errorf("PublishedAt %d should be after year 2020", payload.PublishedAt)
		}
		if int64(payload.PublishedAt) > YEAR_2030_MS {
			t.Errorf("PublishedAt %d should be before year 2030", payload.PublishedAt)
		}

		// Verify millisecond precision
		asSeconds := payload.PublishedAt / 1000
		if payload.PublishedAt <= asSeconds*1000 {
			t.Error("PublishedAt should have millisecond precision")
		}
	})

	t.Run("should create AuditReceivedPayload with millisecond timestamp", func(t *testing.T) {
		payload := event.AuditReceivedPayload{
			PublisherMicroservice: "test-publisher",
			ReceiverMicroservice:  "test-receiver",
			ReceivedEvent:         "test.event",
			ReceivedAt:            uint64(time.Now().UnixMilli()),
			QueueName:             "test-queue",
			EventID:               "test-uuid",
		}

		// Verify timestamp is in reasonable range
		if int64(payload.ReceivedAt) < YEAR_2020_MS {
			t.Errorf("ReceivedAt %d should be after year 2020", payload.ReceivedAt)
		}
		if int64(payload.ReceivedAt) > YEAR_2030_MS {
			t.Errorf("ReceivedAt %d should be before year 2030", payload.ReceivedAt)
		}

		// Verify millisecond precision
		asSeconds := payload.ReceivedAt / 1000
		if payload.ReceivedAt <= asSeconds*1000 {
			t.Error("ReceivedAt should have millisecond precision")
		}
	})

	t.Run("should create AuditProcessedPayload with millisecond timestamp", func(t *testing.T) {
		payload := event.AuditProcessedPayload{
			PublisherMicroservice: "test-publisher",
			ProcessorMicroservice: "test-processor",
			ProcessedEvent:        "test.event",
			ProcessedAt:           uint64(time.Now().UnixMilli()),
			QueueName:             "test-queue",
			EventID:               "test-uuid",
		}

		// Verify timestamp is in reasonable range
		if int64(payload.ProcessedAt) < YEAR_2020_MS {
			t.Errorf("ProcessedAt %d should be after year 2020", payload.ProcessedAt)
		}
		if int64(payload.ProcessedAt) > YEAR_2030_MS {
			t.Errorf("ProcessedAt %d should be before year 2030", payload.ProcessedAt)
		}

		// Verify millisecond precision
		asSeconds := payload.ProcessedAt / 1000
		if payload.ProcessedAt <= asSeconds*1000 {
			t.Error("ProcessedAt should have millisecond precision")
		}
	})

	t.Run("should create AuditDeadLetterPayload with millisecond timestamp", func(t *testing.T) {
		retryCount := uint32(3)
		payload := event.AuditDeadLetterPayload{
			PublisherMicroservice: "test-publisher",
			RejectorMicroservice:  "test-rejector",
			RejectedEvent:         "test.event",
			RejectedAt:            uint64(time.Now().UnixMilli()),
			QueueName:             "test-queue",
			RejectionReason:       "fibonacci_strategy",
			RetryCount:            &retryCount,
			EventID:               "test-uuid",
		}

		// Verify timestamp is in reasonable range
		if int64(payload.RejectedAt) < YEAR_2020_MS {
			t.Errorf("RejectedAt %d should be after year 2020", payload.RejectedAt)
		}
		if int64(payload.RejectedAt) > YEAR_2030_MS {
			t.Errorf("RejectedAt %d should be before year 2030", payload.RejectedAt)
		}

		// Verify millisecond precision
		asSeconds := payload.RejectedAt / 1000
		if payload.RejectedAt <= asSeconds*1000 {
			t.Error("RejectedAt should have millisecond precision")
		}
	})
}

func TestMillisecondVsSecondTimestamps(t *testing.T) {
	t.Run("millisecond timestamp should be ~1000x larger than second timestamp", func(t *testing.T) {
		// Generate timestamp in milliseconds (new way)
		timestampMs := time.Now().UnixMilli()

		// Generate timestamp in seconds (old way)
		timestampS := time.Now().Unix()

		// Milliseconds should be ~1000x larger than seconds
		if timestampMs < timestampS*100 {
			t.Errorf("Millisecond timestamp (%d) should be much larger than second timestamp (%d)", timestampMs, timestampS)
		}

		// If seconds timestamp was treated as milliseconds, it would show as 1970
		if timestampS > YEAR_2020_MS {
			t.Errorf("Second timestamp (%d) would be invalid if treated as milliseconds (should be < %d)", timestampS, YEAR_2020_MS)
		}
	})

	t.Run("millisecond timestamps should have 13 digits, seconds should have 10 digits", func(t *testing.T) {
		timestampMs := time.Now().UnixMilli()
		timestampS := time.Now().Unix()

		// Count digits in millisecond timestamp (should be 13 for year 2025)
		msDigits := len(toString(timestampMs))
		if msDigits != 13 {
			t.Errorf("Millisecond timestamp should have 13 digits, got %d", msDigits)
		}

		// Count digits in second timestamp (should be 10 for year 2025)
		sDigits := len(toString(timestampS))
		if sDigits != 10 {
			t.Errorf("Second timestamp should have 10 digits, got %d", sDigits)
		}
	})
}

// Helper function to convert int64 to string for digit counting
func toString(n int64) string {
	if n == 0 {
		return "0"
	}

	if n < 0 {
		n = -n
	}

	s := ""
	for n > 0 {
		s = string(rune('0'+(n%10))) + s
		n /= 10
	}
	return s
}
