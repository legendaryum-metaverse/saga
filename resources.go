package saga

type (
	Exchange string
	Queue    string
)

const (
	RequeueExchange         Exchange = "requeue_exchange"
	CommandsExchange        Exchange = "commands_exchange"
	MatchingExchange        Exchange = "matching_exchange"
	MatchingRequeueExchange Exchange = "matching_requeue_exchange"
	AuditExchange           Exchange = "audit_exchange"
)

const (
	ReplyToSagaQ             Queue = "reply_to_saga"
	AuditPublishedCommandsQ  Queue = "audit_published_commands"
	AuditReceivedCommandsQ   Queue = "audit_received_commands"
	AuditProcessedCommandsQ  Queue = "audit_processed_commands"
	AuditDeadLetterCommandsQ Queue = "audit_dead_letter_commands"
)
