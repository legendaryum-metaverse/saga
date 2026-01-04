package event

import "time"

type MicroserviceEvent string

type PayloadEvent interface {
	Type() MicroserviceEvent
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	// Audit events - track event lifecycle for monitoring and debugging.
	AuditPublishedEvent  MicroserviceEvent = "audit.published"
	AuditReceivedEvent   MicroserviceEvent = "audit.received"
	AuditProcessedEvent  MicroserviceEvent = "audit.processed"
	AuditDeadLetterEvent MicroserviceEvent = "audit.dead_letter"

	AuthBlockedUserEvent                                     MicroserviceEvent = "auth.blocked_user"
	AuthDeletedUserEvent                                     MicroserviceEvent = "auth.deleted_user"
	AuthLogoutUserEvent                                      MicroserviceEvent = "auth.logout_user"
	AuthNewUserEvent                                         MicroserviceEvent = "auth.new_user"
	LegendMissionsNewMissionCreatedEvent                     MicroserviceEvent = "legend_missions.new_mission_created"
	LegendMissionsOngoingMissionEvent                        MicroserviceEvent = "legend_missions.ongoing_mission"
	LegendMissionsMissionFinishedEvent                       MicroserviceEvent = "legend_missions.mission_finished"
	LegendMissionsSendEmailCryptoMissionCompletedEvent       MicroserviceEvent = "legend_missions.send_email_crypto_mission_completed"
	LegendMissionsSendEmailCodeExchangeMissionCompletedEvent MicroserviceEvent = "legend_missions.send_email_code_exchange_mission_completed"
	LegendMissionsSendEmailNftMissionCompletedEvent          MicroserviceEvent = "legend_missions.send_email_nft_mission_completed"
	LegendRankingsRankingsFinishedEvent                      MicroserviceEvent = "legend_rankings.rankings_finished"
	LegendRankingsNewRankingCreatedEvent                     MicroserviceEvent = "legend_rankings.new_ranking_created"
	LegendRankingsIntermediateRewardEvent                    MicroserviceEvent = "legend_rankings.intermediate_reward"
	LegendRankingsParticipationRewardEvent                   MicroserviceEvent = "legend_rankings.participation_reward"
	LegendShowcaseProductVirtualDeletedEvent                 MicroserviceEvent = "legend_showcase.product_virtual_deleted"
	LegendShowcaseUpdateAllowedMissionSubscriptionIdsEvent   MicroserviceEvent = "legend_showcase.update_allowed_mission_subscription_ids"
	LegendShowcaseUpdateAllowedRankingSubscriptionIdsEvent   MicroserviceEvent = "legend_showcase.update_allowed_ranking_subscription_ids"
	SocialBlockChatEvent                                     MicroserviceEvent = "social.block_chat"
	SocialNewUserEvent                                       MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent                                   MicroserviceEvent = "social.unblock_chat"
	SocialUpdatedUserEvent                                   MicroserviceEvent = "social.updated_user"

	// Billing events - Payment and subscription domain events (No Stripe leakage).
	BillingPaymentCreatedEvent       MicroserviceEvent = "billing.payment.created"
	BillingPaymentSucceededEvent     MicroserviceEvent = "billing.payment.succeeded"
	BillingPaymentFailedEvent        MicroserviceEvent = "billing.payment.failed"
	BillingPaymentRefundedEvent      MicroserviceEvent = "billing.payment.refunded"
	BillingSubscriptionCreatedEvent  MicroserviceEvent = "billing.subscription.created"
	BillingSubscriptionUpdatedEvent  MicroserviceEvent = "billing.subscription.updated"
	BillingSubscriptionRenewedEvent  MicroserviceEvent = "billing.subscription.renewed"
	BillingSubscriptionCanceledEvent MicroserviceEvent = "billing.subscription.canceled"
	BillingSubscriptionExpiredEvent  MicroserviceEvent = "billing.subscription.expired"
)

func MicroserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImageEvent,
		TestMintEvent,

		// Audit events
		AuditPublishedEvent,
		AuditReceivedEvent,
		AuditProcessedEvent,
		AuditDeadLetterEvent,

		AuthBlockedUserEvent,
		AuthDeletedUserEvent,
		AuthLogoutUserEvent,
		AuthNewUserEvent,
		LegendMissionsNewMissionCreatedEvent,
		LegendMissionsOngoingMissionEvent,
		LegendMissionsMissionFinishedEvent,
		LegendMissionsSendEmailCryptoMissionCompletedEvent,
		LegendMissionsSendEmailCodeExchangeMissionCompletedEvent,
		LegendMissionsSendEmailNftMissionCompletedEvent,
		LegendRankingsRankingsFinishedEvent,
		LegendRankingsNewRankingCreatedEvent,
		LegendRankingsIntermediateRewardEvent,
		LegendRankingsParticipationRewardEvent,
		LegendShowcaseProductVirtualDeletedEvent,
		LegendShowcaseUpdateAllowedMissionSubscriptionIdsEvent,
		LegendShowcaseUpdateAllowedRankingSubscriptionIdsEvent,
		SocialBlockChatEvent,
		SocialNewUserEvent,
		SocialUnblockChatEvent,
		SocialUpdatedUserEvent,

		// Billing events
		BillingPaymentCreatedEvent,
		BillingPaymentSucceededEvent,
		BillingPaymentFailedEvent,
		BillingPaymentRefundedEvent,
		BillingSubscriptionCreatedEvent,
		BillingSubscriptionUpdatedEvent,
		BillingSubscriptionRenewedEvent,
		BillingSubscriptionCanceledEvent,
		BillingSubscriptionExpiredEvent,
	}
}

// TestImagePayload is the payload for the test.image event.
type TestImagePayload struct {
	Image string `json:"image"`
}

func (TestImagePayload) Type() MicroserviceEvent {
	return TestImageEvent
}

// TestMintPayload is the payload for the test.mint event.
type TestMintPayload struct {
	Mint string `json:"mint"`
}

func (TestMintPayload) Type() MicroserviceEvent {
	return TestMintEvent
}

// AuthBlockedUserPayload is the payload for the auth.blocked_user event.
type AuthBlockedUserPayload struct {
	UserID               string `json:"userId"`
	BlockType            string `json:"blockType"`
	BlockReason          string `json:"blockReason,omitempty"`
	BlockExpirationHours int    `json:"blockExpirationHours,omitempty"`
}

func (AuthBlockedUserPayload) Type() MicroserviceEvent {
	return AuthBlockedUserEvent
}

// AuthDeletedUserPayload is the payload for the auth.deleted_user event.
type AuthDeletedUserPayload struct {
	UserID string `json:"userId"`
}

func (AuthDeletedUserPayload) Type() MicroserviceEvent {
	return AuthDeletedUserEvent
}

// AuthLogoutUserPayload is the payload for the auth.logout_user event.
type AuthLogoutUserPayload struct {
	UserID string `json:"userId"`
}

func (AuthLogoutUserPayload) Type() MicroserviceEvent {
	return AuthLogoutUserEvent
}

// AuthNewUserPayload is the payload for the auth.new_user event.
type AuthNewUserPayload struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Userlastname string `json:"userlastname"`
}

func (AuthNewUserPayload) Type() MicroserviceEvent {
	return AuthNewUserEvent
}

// LegendMissionsNewMissionCreatedEventPayload is the payload for the legend_missions.new_mission_created.
type LegendMissionsNewMissionCreatedEventPayload struct {
	Title                    string              `json:"title"`
	Author                   string              `json:"author"`
	AuthorEmail              string              `json:"authorEmail"`
	Reward                   int                 `json:"reward"`
	StartDate                string              `json:"startDate"`
	EndDate                  string              `json:"endDate"`
	MaxPlayersClaimingReward int                 `json:"maxPlayersClaimingReward"`
	TimeToReward             int                 `json:"timeToReward"`
	NotificationConfig       *NotificationConfig `json:"notificationConfig,omitempty"`
}

func (LegendMissionsNewMissionCreatedEventPayload) Type() MicroserviceEvent {
	return LegendMissionsNewMissionCreatedEvent
}

// LegendMissionsOngoingMissionEventPayload is the payload for the legend_missions.ongoin_mission.
type LegendMissionsOngoingMissionEventPayload struct {
	RedisKey string `json:"redisKey"`
}

func (LegendMissionsOngoingMissionEventPayload) Type() MicroserviceEvent {
	return LegendMissionsOngoingMissionEvent
}

// MissionFinishedParticipant represents a participant in the mission finished event.
type MissionFinishedParticipant struct {
	UserID   *string `json:"userId,omitempty"`
	Email    *string `json:"email,omitempty"`
	Position *int    `json:"position,omitempty"`
}

// LegendMissionsMissionFinishedEventPayload is the payload for the legend_missions.mission_finished event.
type LegendMissionsMissionFinishedEventPayload struct {
	MissionTitle string                       `json:"missionTitle"`
	Participants []MissionFinishedParticipant `json:"participants"`
}

func (LegendMissionsMissionFinishedEventPayload) Type() MicroserviceEvent {
	return LegendMissionsMissionFinishedEvent
}

type RankingWinners struct {
	UserID string `json:"userId"`
	Reward int    `json:"reward"`
}

type CompletedRanking struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorEmail string `json:"authorEmail"`
	// End date converted to string
	EndsAt string `json:"endsAt"`
	// JSON stringified with each user's rewards
	Reward     string           `json:"reward"`
	RewardType string           `json:"rewardType"`
	Winners    []RankingWinners `json:"winners"`
	// Present only if reward_type is "Nft"
	NftBlockchainNetwork *string `json:"nftBlockchainNetwork,omitempty"`
	NftContractAddress   *string `json:"nftContractAddress,omitempty"`
	// Present only if reward_type is "Crypto"
	WalletCryptoAsset *string `json:"walletCryptoAsset,omitempty"`
	// Optional notification config (dynamic template data)
	NotificationConfig map[string]any `json:"notificationConfig,omitempty"`
}

// LegendMissionsSendEmailCryptoMissionCompletedEventPayload is the payload for the legend_missions.send_email_crypto_mission_completed event.
type LegendMissionsSendEmailCryptoMissionCompletedEventPayload struct {
	UserID            string `json:"userId"`
	MissionTitle      string `json:"missionTitle"`
	Reward            string `json:"reward"`
	BlockchainNetwork string `json:"blockchainNetwork"`
	CryptoAsset       string `json:"cryptoAsset"`
}

func (LegendMissionsSendEmailCryptoMissionCompletedEventPayload) Type() MicroserviceEvent {
	return LegendMissionsSendEmailCryptoMissionCompletedEvent
}

// LegendMissionsSendEmailCodeExchangeMissionCompletedEventPayload is the payload for the legend_missions.send_email_code_exchange_mission_completed event.
type LegendMissionsSendEmailCodeExchangeMissionCompletedEventPayload struct {
	UserID          string `json:"userId"`
	MissionTitle    string `json:"missionTitle"`
	CodeValue       string `json:"codeValue"`
	CodeDescription string `json:"codeDescription"`
}

func (LegendMissionsSendEmailCodeExchangeMissionCompletedEventPayload) Type() MicroserviceEvent {
	return LegendMissionsSendEmailCodeExchangeMissionCompletedEvent
}

// LegendMissionsSendEmailNftMissionCompletedEventPayload is the payload for the legend_missions.send_email_nft_mission_completed event.
type LegendMissionsSendEmailNftMissionCompletedEventPayload struct {
	UserID             string `json:"userId"`
	MissionTitle       string `json:"missionTitle"`
	NftContractAddress string `json:"nftContractAddress"`
	NftTokenID         string `json:"nftTokenId"`
}

func (LegendMissionsSendEmailNftMissionCompletedEventPayload) Type() MicroserviceEvent {
	return LegendMissionsSendEmailNftMissionCompletedEvent
}

// LegendRankingsRankingsFinishedEventPayload is the payload for the legend_rankings.rankings_finished.
type LegendRankingsRankingsFinishedEventPayload struct {
	CompletedRankings []CompletedRanking `json:"completedRankings"`
}

func (LegendRankingsRankingsFinishedEventPayload) Type() MicroserviceEvent {
	return LegendRankingsRankingsFinishedEvent
}

// LegendRankingsNewRankingCreatedEventPayload is the payload for the legend_rankings.new_ranking_created event.
type LegendRankingsNewRankingCreatedEventPayload struct {
	Title                string              `json:"title"`
	Description          string              `json:"description"`
	AuthorEmail          string              `json:"authorEmail"`
	RewardType           string              `json:"rewardType"`
	StartAt              string              `json:"startAt"`
	EndsAt               string              `json:"endsAt"`
	NftBlockchainNetwork *string             `json:"nftBlockchainNetwork,omitempty"`
	NftContractAddress   *string             `json:"nftContractAddress,omitempty"`
	WalletCryptoAsset    *string             `json:"walletCryptoAsset,omitempty"`
	NotificationConfig   *NotificationConfig `json:"notificationConfig,omitempty"`
}

// NotificationConfig represents the notification configuration.
type NotificationConfig struct {
	CustomEmails *[]string `json:"customEmails,omitempty"`
	TemplateName string    `json:"templateName"`
}

func (LegendRankingsNewRankingCreatedEventPayload) Type() MicroserviceEvent {
	return LegendRankingsNewRankingCreatedEvent
}

// LegendRankingsIntermediateRewardEventPayload is the payload for the legend_rankings.intermediate_reward event.
type LegendRankingsIntermediateRewardEventPayload struct {
	UserID                 string                 `json:"userId"`
	RankingID              int                    `json:"rankingId"`
	IntermediateRewardType string                 `json:"intermediateRewardType"`
	RewardConfig           map[string]interface{} `json:"rewardConfig"`
	TemplateName           string                 `json:"templateName"`
	TemplateData           map[string]interface{} `json:"templateData"`
}

func (LegendRankingsIntermediateRewardEventPayload) Type() MicroserviceEvent {
	return LegendRankingsIntermediateRewardEvent
}

// LegendRankingsParticipationRewardEventPayload is the payload for the legend_rankings.participation_reward event.
type LegendRankingsParticipationRewardEventPayload struct {
	UserID                  string                 `json:"userId"`
	RankingID               int                    `json:"rankingId"`
	ParticipationRewardType string                 `json:"participationRewardType"`
	RewardConfig            map[string]interface{} `json:"rewardConfig"`
	TemplateName            string                 `json:"templateName"`
	TemplateData            map[string]interface{} `json:"templateData"`
}

func (LegendRankingsParticipationRewardEventPayload) Type() MicroserviceEvent {
	return LegendRankingsParticipationRewardEvent
}

// LegendShowcaseProductVirtualDeletedEventPayload is the payload for the legend_showcase.product_virtual_deleted event.
type LegendShowcaseProductVirtualDeletedEventPayload struct {
	ProductVirtualID   string `json:"productVirtualId"`
	ProductVirtualSlug string `json:"productVirtualSlug"`
}

func (LegendShowcaseProductVirtualDeletedEventPayload) Type() MicroserviceEvent {
	return LegendShowcaseProductVirtualDeletedEvent
}

// LegendShowcaseUpdateAllowedMissionSubscriptionIdsEventPayload is the payload for the legend_showcase.update_allowed_mission_subscription_ids.
type LegendShowcaseUpdateAllowedMissionSubscriptionIdsEventPayload struct {
	ProductVirtualSlug     string   `json:"productVirtualSlug"`
	AllowedSubscriptionIds []string `json:"allowedSubscriptionIds"`
}

func (LegendShowcaseUpdateAllowedMissionSubscriptionIdsEventPayload) Type() MicroserviceEvent {
	return LegendShowcaseUpdateAllowedMissionSubscriptionIdsEvent
}

// LegendShowcaseUpdateAllowedRankingSubscriptionIdsEventPayload is the payload for the legend_showcase.update_allowed_ranking_subscription_ids.
type LegendShowcaseUpdateAllowedRankingSubscriptionIdsEventPayload struct {
	ProductVirtualID       string   `json:"productVirtualId"`
	AllowedSubscriptionIds []string `json:"allowedSubscriptionIds"`
}

func (LegendShowcaseUpdateAllowedRankingSubscriptionIdsEventPayload) Type() MicroserviceEvent {
	return LegendShowcaseUpdateAllowedRankingSubscriptionIdsEvent
}

// SocialBlockChatPayload is the payload for the social.block_chat event.
type SocialBlockChatPayload struct {
	UserID        string `json:"userId"`
	UserToBlockID string `json:"userToBlockId"`
}

func (SocialBlockChatPayload) Type() MicroserviceEvent {
	return SocialBlockChatEvent
}

// Gender represents the possible genders a social user can have.
type Gender string

const (
	GenderMale      Gender = "MALE"
	GenderFemale    Gender = "FEMALE"
	GenderUndefined Gender = "UNDEFINED"
)

// UserLocation represents the geographical location of a user.
type UserLocation struct {
	Continent string `json:"continent" bson:"continent"`
	Country   string `json:"country" bson:"country"`
	Region    string `json:"region" bson:"region"`
	City      string `json:"city" bson:"city"`
}

// SocialMedia represents social media links as a map.
type SocialMedia map[string]string

// SocialUser represents the main user model.
type SocialUser struct {
	ID               string        `json:"_id" bson:"_id"`
	Username         string        `json:"username" bson:"username"`
	FirstName        *string       `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName         *string       `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Gender           Gender        `json:"gender" bson:"gender"`
	IsPublicProfile  bool          `json:"isPublicProfile,omitempty" bson:"isPublicProfile,omitempty"`
	Followers        []string      `json:"followers" bson:"followers"`
	Following        []string      `json:"following" bson:"following"`
	Email            string        `json:"email" bson:"email"`
	Birthday         *time.Time    `json:"birthday,omitempty" bson:"birthday,omitempty"`
	Location         *UserLocation `json:"location,omitempty" bson:"location,omitempty"`
	Avatar           *string       `json:"avatar,omitempty" bson:"avatar,omitempty"`
	AvatarScreenshot *string       `json:"avatarScreenshot,omitempty" bson:"avatarScreenshot,omitempty"`
	UserImage        *string       `json:"userImage,omitempty" bson:"userImage,omitempty"`
	GlbURL           *string       `json:"glbUrl,omitempty" bson:"glbUrl,omitempty"`
	Description      *string       `json:"description,omitempty" bson:"description,omitempty"`
	SocialMedia      *SocialMedia  `json:"socialMedia,omitempty" bson:"socialMedia,omitempty"`
	Preferences      []string      `json:"preferences" bson:"preferences"`
	BlockedUsers     []string      `json:"blockedUsers" bson:"blockedUsers"`
	RPMAvatarID      *string       `json:"rpmAvatarId,omitempty" bson:"rpmAvatarId,omitempty"`
	RPMUserID        *string       `json:"rpmUserId,omitempty" bson:"rpmUserId,omitempty"`
	PaidPriceID      *string       `json:"paidPriceId,omitempty" bson:"paidPriceId,omitempty"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
}

// SocialNewUserPayload is the payload for the social.new_user event.
type SocialNewUserPayload struct {
	SocialUser SocialUser `json:"socialUser"`
}

func (SocialNewUserPayload) Type() MicroserviceEvent {
	return SocialNewUserEvent
}

// SocialUnblockChatPayload is the payload for the social.unblock_chat event.
type SocialUnblockChatPayload struct {
	UserID          string `json:"userId"`
	UserToUnblockID string `json:"userToUnblockId"`
}

func (SocialUnblockChatPayload) Type() MicroserviceEvent {
	return SocialUnblockChatEvent
}

// SocialUpdatedUserPayload is the payload for the social.updated_user event.
type SocialUpdatedUserPayload struct {
	SocialUser SocialUser `json:"socialUser"`
}

func (SocialUpdatedUserPayload) Type() MicroserviceEvent {
	return SocialUpdatedUserEvent
}

// ********** AUDIT PAYLOADS ************** //

// AuditReceivedPayload is the payload for audit.received event - tracks when event is received before processing.
type AuditReceivedPayload struct {
	// The microservice that published the original event
	PublisherMicroservice string `json:"publisher_microservice"`
	// The microservice that received the event
	ReceiverMicroservice string `json:"receiver_microservice"`
	// The event that was received
	ReceivedEvent string `json:"received_event"`
	// Timestamp when the event was received (UNIX timestamp in milliseconds)
	ReceivedAt uint64 `json:"received_at"`
	// The queue name from which the event was consumed
	QueueName string `json:"queue_name"`
	// Event identifier for tracking across the event lifecycle
	EventID string `json:"event_id"`
}

func (AuditReceivedPayload) Type() MicroserviceEvent {
	return AuditReceivedEvent
}

// AuditProcessedPayload is the payload for audit.processed event - tracks successful event processing.
type AuditProcessedPayload struct {
	// The microservice that published the original event
	PublisherMicroservice string `json:"publisher_microservice"`
	// The microservice that processed the event
	ProcessorMicroservice string `json:"processor_microservice"`
	// The original event that was processed
	ProcessedEvent string `json:"processed_event"`
	// Timestamp when the event was processed (UNIX timestamp in milliseconds)
	ProcessedAt uint64 `json:"processed_at"`
	// The queue name where the event was consumed
	QueueName string `json:"queue_name"`
	// Event identifier for tracking across the event lifecycle
	EventID string `json:"event_id"`
}

func (AuditProcessedPayload) Type() MicroserviceEvent {
	return AuditProcessedEvent
}

// AuditDeadLetterPayload is the payload for audit.dead_letter event - tracks when message is rejected/nacked.
type AuditDeadLetterPayload struct {
	// The microservice that published the original event
	PublisherMicroservice string `json:"publisher_microservice"`
	// The microservice that rejected the event
	RejectorMicroservice string `json:"rejector_microservice"`
	// The original event that was rejected
	RejectedEvent string `json:"rejected_event"`
	// Timestamp when the event was rejected (UNIX timestamp in milliseconds)
	RejectedAt uint64 `json:"rejected_at"`
	// The queue name where the event was rejected from
	QueueName string `json:"queue_name"`
	// Reason for rejection (delay, fibonacci_strategy, etc.)
	RejectionReason string `json:"rejection_reason"`
	// Optional retry count
	RetryCount *uint32 `json:"retry_count,omitempty"`
	// Event identifier for tracking across the event lifecycle
	EventID string `json:"event_id"`
}

func (AuditDeadLetterPayload) Type() MicroserviceEvent {
	return AuditDeadLetterEvent
}

// AuditPublishedPayload is the payload for audit.published event - tracks when event is published by a microservice.
type AuditPublishedPayload struct {
	// The microservice that published the event
	PublisherMicroservice string `json:"publisher_microservice"`
	// The event that was published
	PublishedEvent string `json:"published_event"`
	// Timestamp when the event was published (UNIX timestamp in milliseconds)
	PublishedAt uint64 `json:"published_at"`
	// Event identifier for tracking across the event lifecycle
	EventID string `json:"event_id"`
}

func (AuditPublishedPayload) Type() MicroserviceEvent {
	return AuditPublishedEvent
}

// ======================================================================================================
// BILLING PAYLOADS - Payment and subscription domain events (No Stripe leakage - only internal IDs)
// ======================================================================================================

// BillingPaymentCreatedPayload is the payload for billing.payment.created event.
type BillingPaymentCreatedPayload struct {
	PaymentID  string            `json:"paymentId"`
	UserID     string            `json:"userId"`
	Amount     int64             `json:"amount"`
	Currency   string            `json:"currency"`
	Status     string            `json:"status"` // "pending" | "processing"
	Metadata   map[string]string `json:"metadata"`
	OccurredAt string            `json:"occurredAt"`
}

func (BillingPaymentCreatedPayload) Type() MicroserviceEvent {
	return BillingPaymentCreatedEvent
}

// BillingPaymentSucceededPayload is the payload for billing.payment.succeeded event.
type BillingPaymentSucceededPayload struct {
	PaymentID  string            `json:"paymentId"`
	UserID     string            `json:"userId"`
	Amount     int64             `json:"amount"`
	Currency   string            `json:"currency"`
	Metadata   map[string]string `json:"metadata"`
	OccurredAt string            `json:"occurredAt"`
}

func (BillingPaymentSucceededPayload) Type() MicroserviceEvent {
	return BillingPaymentSucceededEvent
}

// BillingPaymentFailedPayload is the payload for billing.payment.failed event.
type BillingPaymentFailedPayload struct {
	PaymentID     string            `json:"paymentId"`
	UserID        string            `json:"userId"`
	Amount        int64             `json:"amount"`
	Currency      string            `json:"currency"`
	FailureReason *string           `json:"failureReason,omitempty"`
	Metadata      map[string]string `json:"metadata"`
	OccurredAt    string            `json:"occurredAt"`
}

func (BillingPaymentFailedPayload) Type() MicroserviceEvent {
	return BillingPaymentFailedEvent
}

// BillingPaymentRefundedPayload is the payload for billing.payment.refunded event.
type BillingPaymentRefundedPayload struct {
	PaymentID      string            `json:"paymentId"`
	UserID         string            `json:"userId"`
	Amount         int64             `json:"amount"`
	RefundedAmount int64             `json:"refundedAmount"`
	Currency       string            `json:"currency"`
	Metadata       map[string]string `json:"metadata"`
	OccurredAt     string            `json:"occurredAt"`
}

func (BillingPaymentRefundedPayload) Type() MicroserviceEvent {
	return BillingPaymentRefundedEvent
}

// BillingSubscriptionCreatedPayload is the payload for billing.subscription.created event.
type BillingSubscriptionCreatedPayload struct {
	SubscriptionID string `json:"subscriptionId"`
	UserID         string `json:"userId"`
	PlanID         string `json:"planId"`
	PlanSlug       string `json:"planSlug"`
	Status         string `json:"status"` // "pending" | "active" | "trialing"
	PeriodStart    string `json:"periodStart"`
	PeriodEnd      string `json:"periodEnd"`
	OccurredAt     string `json:"occurredAt"`
}

func (BillingSubscriptionCreatedPayload) Type() MicroserviceEvent {
	return BillingSubscriptionCreatedEvent
}

// BillingSubscriptionUpdatedPayload is the payload for billing.subscription.updated event.
type BillingSubscriptionUpdatedPayload struct {
	SubscriptionID    string `json:"subscriptionId"`
	UserID            string `json:"userId"`
	PlanID            string `json:"planId"`
	PlanSlug          string `json:"planSlug"`
	Status            string `json:"status"` // "active" | "past_due" | "unpaid" | "paused" | "trialing"
	CancelAtPeriodEnd bool   `json:"cancelAtPeriodEnd"`
	PeriodStart       string `json:"periodStart"`
	PeriodEnd         string `json:"periodEnd"`
	OccurredAt        string `json:"occurredAt"`
}

func (BillingSubscriptionUpdatedPayload) Type() MicroserviceEvent {
	return BillingSubscriptionUpdatedEvent
}

// BillingSubscriptionRenewedPayload is the payload for billing.subscription.renewed event.
type BillingSubscriptionRenewedPayload struct {
	SubscriptionID string `json:"subscriptionId"`
	UserID         string `json:"userId"`
	PlanID         string `json:"planId"`
	PlanSlug       string `json:"planSlug"`
	PeriodStart    string `json:"periodStart"`
	PeriodEnd      string `json:"periodEnd"`
	OccurredAt     string `json:"occurredAt"`
}

func (BillingSubscriptionRenewedPayload) Type() MicroserviceEvent {
	return BillingSubscriptionRenewedEvent
}

// BillingSubscriptionCanceledPayload is the payload for billing.subscription.canceled event.
type BillingSubscriptionCanceledPayload struct {
	SubscriptionID string `json:"subscriptionId"`
	UserID         string `json:"userId"`
	PlanID         string `json:"planId"`
	PlanSlug       string `json:"planSlug"`
	CanceledAt     string `json:"canceledAt"`
	OccurredAt     string `json:"occurredAt"`
}

func (BillingSubscriptionCanceledPayload) Type() MicroserviceEvent {
	return BillingSubscriptionCanceledEvent
}

// BillingSubscriptionExpiredPayload is the payload for billing.subscription.expired event.
type BillingSubscriptionExpiredPayload struct {
	SubscriptionID string `json:"subscriptionId"`
	UserID         string `json:"userId"`
	PlanID         string `json:"planId"`
	PlanSlug       string `json:"planSlug"`
	ExpiredAt      string `json:"expiredAt"`
	OccurredAt     string `json:"occurredAt"`
}

func (BillingSubscriptionExpiredPayload) Type() MicroserviceEvent {
	return BillingSubscriptionExpiredEvent
}
