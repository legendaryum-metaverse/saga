package event

import "time"

type MicroserviceEvent string

type PayloadEvent interface {
	Type() MicroserviceEvent
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	AuthDeletedUserEvent                      MicroserviceEvent = "auth.deleted_user"
	AuthLogoutUserEvent                       MicroserviceEvent = "auth.logout_user"
	AuthNewUserEvent                          MicroserviceEvent = "auth.new_user"
	CoinsNotifyClientEvent                    MicroserviceEvent = "coins.notify_client"
	CoinsSendEmail                            MicroserviceEvent = "coins.send_email"
	CoinsUpdateSubscription                   MicroserviceEvent = "coins.update_subscription"
	LegendMissionsCompletedMissionRewardEvent MicroserviceEvent = "legend_missions.completed_mission_reward"
	LegendMissionsOngoingMissionEvent         MicroserviceEvent = "legend_missions.ongoing_mission"
	LegendRankingsRankingsFinishedEvent       MicroserviceEvent = "legend_rankings.rankings_finished"
	RoomCreatorCreatedRoomEvent               MicroserviceEvent = "room_creator.created_room"
	RoomCreatorUpdatedRoomEvent               MicroserviceEvent = "room_creator.updated_room"
	RoomInventoryUpdateVpBuildingImageEvent   MicroserviceEvent = "room_inventory.update_vp_building_image"
	RoomSnapshotBuildingChangeInIslandEvent   MicroserviceEvent = "room_snapshot.building_change_in_island"
	RoomSnapshotFirstSnapshotEvent            MicroserviceEvent = "room_snapshot.first_snapshot"
	SocialBlockChatEvent                      MicroserviceEvent = "social.block_chat"
	SocialMediaRoomsDeleteInBatchEvent        MicroserviceEvent = "social_media_rooms.delete_in_batch"
	SocialNewUserEvent                        MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent                    MicroserviceEvent = "social.unblock_chat"
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
		LegendMissionsOngoingMissionEvent,
		LegendRankingsRankingsFinishedEvent,
		RoomCreatorCreatedRoomEvent,
		RoomCreatorUpdatedRoomEvent,
		RoomInventoryUpdateVpBuildingImageEvent,
		RoomSnapshotBuildingChangeInIslandEvent,
		RoomSnapshotFirstSnapshotEvent,
		SocialBlockChatEvent,
		SocialMediaRoomsDeleteInBatchEvent,
		SocialNewUserEvent,
		SocialUnblockChatEvent,
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

// LegendMissionsOngoingMissionEventPayload is the payload for the legend_missions.ongoin_mission.
type LegendMissionsOngoingMissionEventPayload struct {
	RedisKey string `json:"redisKey"`
}

func (LegendMissionsOngoingMissionEventPayload) Type() MicroserviceEvent {
	return LegendMissionsOngoingMissionEvent
}

type RewardType string

const (
	Legends      RewardType = "Legends"
	CodeExchange RewardType = "CodeExchange"
)

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
	RewardType RewardType       `json:"rewardType"`
	Winners    []RankingWinners `json:"winners"`
}

// LegendRankingsRankingsFinishedEventPayload is the payload for the legend_rankings.rankings_finished.
type LegendRankingsRankingsFinishedEventPayload struct {
	CompletedRankings []CompletedRanking `json:"completedRankings"`
}

func (LegendRankingsRankingsFinishedEventPayload) Type() MicroserviceEvent {
	return LegendRankingsRankingsFinishedEvent
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

// UserLocation represents the user's location.
type UserLocation struct {
	Continent string `json:"continent"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
}

// SocialUser represents the social user model.
type SocialUser struct {
	ID               string             `json:"_id"`
	Username         string             `json:"username"`
	FirstName        *string            `json:"firstName"`
	LastName         *string            `json:"lastName"`
	Gender           Gender             `json:"gender"`
	IsPublicProfile  bool               `json:"isPublicProfile"`
	Followers        []string           `json:"followers"`
	Following        []string           `json:"following"`
	Email            string             `json:"email"`
	Birthday         *time.Time         `json:"birthday"`
	Location         *UserLocation      `json:"location"`
	Avatar           *string            `json:"avatar"`
	AvatarScreenshot *string            `json:"avatarScreenshot"`
	UserImage        *string            `json:"userImage"`
	GlbUrl           *string            `json:"glbUrl"`
	Description      *string            `json:"description"`
	SocialMedia      *map[string]string `json:"socialMedia"`
	Preferences      []string           `json:"preferences"`
	BlockedUsers     []string           `json:"blockedUsers"`
	RPMAvatarID      *string            `json:"RPMAvatarId"`
	RPMUserID        *string            `json:"RPMUserId"`
	PaidPriceID      *string            `json:"paidPriceId"`
	CreatedAt        time.Time          `json:"createdAt"`
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
