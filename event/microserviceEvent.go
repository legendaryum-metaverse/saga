package event

type MicroserviceEvent string

type PayloadEvent interface {
	Type() MicroserviceEvent
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	AuthDeletedUserEvent                    MicroserviceEvent = "auth.deleted_user"
	PaymentsNotifyClientEvent               MicroserviceEvent = "payments.notify_client"
	RoomCreatorCreatedRoomEvent             MicroserviceEvent = "room_creator.created_room"
	RoomCreatorUpdatedRoomEvent             MicroserviceEvent = "room_creator.updated_room"
	RoomInventoryUpdateVpBuildingImageEvent MicroserviceEvent = "room_inventory.update_vp_building_image"
	RoomSnapshotBuildingChangeInIslandEvent MicroserviceEvent = "room_snapshot.building_change_in_island"
	RoomSnapshotFirstSnapshotEvent          MicroserviceEvent = "room_snapshot.first_snapshot"
	SocialBlockChatEvent                    MicroserviceEvent = "social.block_chat"
	SocialNewUserEvent                      MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent                  MicroserviceEvent = "social.unblock_chat"
)

func MicroserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImageEvent,
		TestMintEvent,

		AuthDeletedUserEvent,
		PaymentsNotifyClientEvent,
		RoomCreatorCreatedRoomEvent,
		RoomCreatorUpdatedRoomEvent,
		RoomInventoryUpdateVpBuildingImageEvent,
		RoomSnapshotBuildingChangeInIslandEvent,
		RoomSnapshotFirstSnapshotEvent,
		SocialBlockChatEvent,
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

// PaymentsNotifyClientPayload is the payload for the payments.notify_client event.
type PaymentsNotifyClientPayload struct {
	Room    string                 `json:"room"`
	Message map[string]interface{} `json:"message"`
}

func (PaymentsNotifyClientPayload) Type() MicroserviceEvent {
	return PaymentsNotifyClientEvent
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
	Room   Room     `json:"room"`
	Images []string `json:"images"`
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

// SocialNewUserPayload is the payload for the social.new_user event.
type SocialNewUserPayload struct {
	UserID string `json:"userId"`
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
