package event

import "time"

type MicroserviceEvent string

type PayloadEvent interface {
	Type() MicroserviceEvent
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	AuthDeletedUserEvent                                     MicroserviceEvent = "auth.deleted_user"
	AuthLogoutUserEvent                                      MicroserviceEvent = "auth.logout_user"
	AuthNewUserEvent                                         MicroserviceEvent = "auth.new_user"
	CoinsNotifyClientEvent                                   MicroserviceEvent = "coins.notify_client"
	CoinsSendEmail                                           MicroserviceEvent = "coins.send_email"
	CoinsUpdateSubscription                                  MicroserviceEvent = "coins.update_subscription"
	LegendMissionsCompletedMissionRewardEvent                MicroserviceEvent = "legend_missions.completed_mission_reward"
	LegendMissionsNewMissionCreatedEvent                     MicroserviceEvent = "legend_missions.new_mission_created"
	LegendMissionsOngoingMissionEvent                        MicroserviceEvent = "legend_missions.ongoing_mission"
	LegendMissionsSendEmailCryptoMissionCompletedEvent       MicroserviceEvent = "legend_missions.send_email_crypto_mission_completed"
	LegendMissionsSendEmailCodeExchangeMissionCompletedEvent MicroserviceEvent = "legend_missions.send_email_code_exchange_mission_completed"
	LegendMissionsSendEmailNftMissionCompletedEvent          MicroserviceEvent = "legend_missions.send_email_nft_mission_completed"
	LegendRankingsRankingsFinishedEvent                      MicroserviceEvent = "legend_rankings.rankings_finished"
	LegendRankingsNewRankingCreatedEvent                     MicroserviceEvent = "legend_rankings.new_ranking_created"
	LegendShowcaseProductVirtualDeletedEvent                 MicroserviceEvent = "legend_showcase.product_virtual_deleted"
	LegendShowcaseUpdateAllowedMissionSubscriptionIdsEvent   MicroserviceEvent = "legend_showcase.update_allowed_mission_subscription_ids"
	LegendShowcaseUpdateAllowedRankingSubscriptionIdsEvent   MicroserviceEvent = "legend_showcase.update_allowed_ranking_subscription_ids"
	RoomCreatorCreatedRoomEvent                              MicroserviceEvent = "room_creator.created_room"
	RoomCreatorUpdatedRoomEvent                              MicroserviceEvent = "room_creator.updated_room"
	RoomInventoryUpdateVpBuildingImageEvent                  MicroserviceEvent = "room_inventory.update_vp_building_image"
	RoomSnapshotBuildingChangeInIslandEvent                  MicroserviceEvent = "room_snapshot.building_change_in_island"
	RoomSnapshotFirstSnapshotEvent                           MicroserviceEvent = "room_snapshot.first_snapshot"
	SocialBlockChatEvent                                     MicroserviceEvent = "social.block_chat"
	SocialMediaRoomsDeleteInBatchEvent                       MicroserviceEvent = "social_media_rooms.delete_in_batch"
	SocialNewUserEvent                                       MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent                                   MicroserviceEvent = "social.unblock_chat"
	SocialUpdatedUserEvent                                   MicroserviceEvent = "social.updated_user"
)

func MicroserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImageEvent,
		TestMintEvent,

		AuthDeletedUserEvent,
		AuthLogoutUserEvent,
		AuthNewUserEvent,
		CoinsUpdateSubscription,
		CoinsNotifyClientEvent,
		CoinsSendEmail,
		LegendMissionsCompletedMissionRewardEvent,
		LegendMissionsNewMissionCreatedEvent,
		LegendMissionsOngoingMissionEvent,
		LegendMissionsSendEmailCryptoMissionCompletedEvent,
		LegendMissionsSendEmailCodeExchangeMissionCompletedEvent,
		LegendMissionsSendEmailNftMissionCompletedEvent,
		LegendRankingsRankingsFinishedEvent,
		LegendRankingsNewRankingCreatedEvent,
		LegendShowcaseProductVirtualDeletedEvent,
		LegendShowcaseUpdateAllowedMissionSubscriptionIdsEvent,
		LegendShowcaseUpdateAllowedRankingSubscriptionIdsEvent,
		RoomCreatorCreatedRoomEvent,
		RoomCreatorUpdatedRoomEvent,
		RoomInventoryUpdateVpBuildingImageEvent,
		RoomSnapshotBuildingChangeInIslandEvent,
		RoomSnapshotFirstSnapshotEvent,
		SocialBlockChatEvent,
		SocialMediaRoomsDeleteInBatchEvent,
		SocialNewUserEvent,
		SocialUnblockChatEvent,
		SocialUpdatedUserEvent,
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

// CoinsUpdateSubscriptionPayload is the payload for the coins.update_subscription event.
type CoinsUpdateSubscriptionPayload struct {
	UserID      string `json:"userId"`
	PaidPriceId string `json:"paidPriceId"`
}

func (CoinsUpdateSubscriptionPayload) Type() MicroserviceEvent {
	return CoinsUpdateSubscription
}

// CoinsNotifyClientPayload is the payload for the coins.notify_client event.
type CoinsNotifyClientPayload struct {
	Room    string                 `json:"room"`
	Message map[string]interface{} `json:"message"`
}

func (CoinsNotifyClientPayload) Type() MicroserviceEvent {
	return CoinsNotifyClientEvent
}

// CoinsSendEmailPayload is the payload for the coins.send_email event.
type CoinsSendEmailPayload struct {
	UserId    string `json:"userId"`
	EmailType string `json:"emailType"`
	Coins     int    `json:"coins"`
	Email     string `json:"email"`
}

func (CoinsSendEmailPayload) Type() MicroserviceEvent {
	return CoinsSendEmail
}

// LegendMissionsCompletedMissionRewardEventPayload is the payload for the legend_missions.completed_mission_reward event.
type LegendMissionsCompletedMissionRewardEventPayload struct {
	UserID string `json:"userId"`
	Coins  int    `json:"coins"`
}

func (LegendMissionsCompletedMissionRewardEventPayload) Type() MicroserviceEvent {
	return LegendMissionsCompletedMissionRewardEvent
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
	EndsAt               string              `json:"endsAt"`
	NftBlockchainNetwork *string             `json:"nftBlockchainNetwork,omitempty"`
	NftContractAddress   *string             `json:"nftContractAddress,omitempty"`
	WalletCryptoAsset    *string             `json:"walletCryptoAsset,omitempty"`
	NotificationConfig   *NotificationConfig `json:"notificationConfig,omitempty"`
}

// NotificationConfig represents the notification configuration
type NotificationConfig struct {
	CustomEmails *[]string `json:"customEmails,omitempty"`
	TemplateName string    `json:"templateName"`
}

func (LegendRankingsNewRankingCreatedEventPayload) Type() MicroserviceEvent {
	return LegendRankingsNewRankingCreatedEvent
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

type Room struct {
	Id         string `json:"Id"`
	CreateAt   string `json:"CreateAt"`
	UpdateAt   string `json:"UpdateAt"`
	RoomType   string `json:"type"`
	Name       string `json:"name"`
	OwnerId    string `json:"ownerId"`
	OwnerEmail string `json:"ownerEmail"`
	MaxPlayers int    `json:"maxPlayers"`
	MaxLayers  int    `json:"maxLayers"`
	TemplateId string `json:"templateId"`
	HaveEditor bool   `json:"haveEditor"`
}

// RoomCreatorCreatedRoomPayload is the payload for the room_creator.created_room event.
type RoomCreatorCreatedRoomPayload struct {
	Room Room `json:"room"`
}

func (RoomCreatorCreatedRoomPayload) Type() MicroserviceEvent {
	return RoomCreatorCreatedRoomEvent
}

// RoomCreatorUpdatedRoomPayload is the payload for the room_creator.updated_room event.
type RoomCreatorUpdatedRoomPayload struct {
	Room Room `json:"room"`
}

func (RoomCreatorUpdatedRoomPayload) Type() MicroserviceEvent {
	return RoomCreatorUpdatedRoomEvent
}

// RoomInventoryUpdateVpBuildingImagePayload is the payload for the room_snapshot.room_inventory.update_vp_building_image event.
type RoomInventoryUpdateVpBuildingImagePayload struct {
	Images   []string `json:"images"`
	RoomType string   `json:"roomType"`
	UserID   string   `json:"userId"`
}

func (RoomInventoryUpdateVpBuildingImagePayload) Type() MicroserviceEvent {
	return RoomInventoryUpdateVpBuildingImageEvent
}

// RoomSnapshotBuildingChangeInIslandPayload is the payload for the room_snapshot.building_change_in_island event.
type RoomSnapshotBuildingChangeInIslandPayload struct {
	Building string `json:"building"`
	UserID   string `json:"userId"`
}

func (RoomSnapshotBuildingChangeInIslandPayload) Type() MicroserviceEvent {
	return RoomSnapshotBuildingChangeInIslandEvent
}

// RoomSnapshotFirstSnapshotPayload is the payload for the room_snapshot.first_snapshot event.
type RoomSnapshotFirstSnapshotPayload struct {
	Slug string `json:"slug"`
}

func (RoomSnapshotFirstSnapshotPayload) Type() MicroserviceEvent {
	return RoomSnapshotFirstSnapshotEvent
}

// SocialBlockChatPayload is the payload for the social.block_chat event.
type SocialBlockChatPayload struct {
	UserID        string `json:"userId"`
	UserToBlockID string `json:"userToBlockId"`
}

func (SocialBlockChatPayload) Type() MicroserviceEvent {
	return SocialBlockChatEvent
}

// SocialMediaRoomsDeleteInBatchPayload is the payload for the social_media_rooms.delete_in_batch event.
type SocialMediaRoomsDeleteInBatchPayload struct {
	BucketName string   `json:"bucketName"`
	FilePaths  []string `json:"filePaths"`
}

func (SocialMediaRoomsDeleteInBatchPayload) Type() MicroserviceEvent {
	return SocialMediaRoomsDeleteInBatchEvent
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
