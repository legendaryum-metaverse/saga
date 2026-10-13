package test

import (
	"encoding/json"
	"testing"

	"github.com/legendaryum-metaverse/saga/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBillingInvoiceChargeCreatedEvent(t *testing.T) {
	assert.Equal(t, "billing.invoice_charge_created", string(event.BillingInvoiceChargeCreatedEvent))
	assert.Contains(t, event.MicroserviceEventValues(), event.BillingInvoiceChargeCreatedEvent)
}

func TestBillingInvoiceChargeCreatedPayloadWireFormat(t *testing.T) {
	payload := event.BillingInvoiceChargeCreatedPayload{
		ChargeID:   "c-1",
		UserID:     "u-1",
		Amount:     1500,
		Currency:   "usd",
		Metadata:   map[string]string{"product_type": "suggested_game", "product_id": "p-1"},
		OccurredAt: "2026-10-13T00:00:00Z",
	}
	assert.Equal(t, event.BillingInvoiceChargeCreatedEvent, payload.Type())

	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var fields map[string]any
	require.NoError(t, json.Unmarshal(raw, &fields))
	assert.Equal(t, "c-1", fields["chargeId"])
	assert.Equal(t, "u-1", fields["userId"])
	assert.InDelta(t, 1500, fields["amount"], 0)
	assert.Equal(t, "usd", fields["currency"])
	assert.Equal(t, "2026-10-13T00:00:00Z", fields["occurredAt"])
	assert.Len(t, fields["metadata"], 2)

	var back event.BillingInvoiceChargeCreatedPayload
	require.NoError(t, json.Unmarshal(raw, &back))
	assert.Equal(t, payload, back)
}
