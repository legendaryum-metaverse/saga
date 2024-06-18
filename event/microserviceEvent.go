package event

type MicroserviceEvent string

type PayloadEvent interface {
	Type() MicroserviceEvent
}

const (
	TestImageEvent MicroserviceEvent = "test.image"
	TestMintEvent  MicroserviceEvent = "test.mint"

	PaymentsNotifyClientEvent      MicroserviceEvent = "payments.notify_client"
	RoomCreatorRoomCreatedEvent    MicroserviceEvent = "room_creator.room_created"
	RoomCreatorRoomUpdatedEvent    MicroserviceEvent = "room_creator.room_updated"
	RoomSnapshotFirstSnapshotEvent MicroserviceEvent = "room_snapshot.first_snapshot"
	SocialBlockChatEvent           MicroserviceEvent = "social.block_chat"
	SocialNewUserEvent             MicroserviceEvent = "social.new_user"
	SocialUnblockChatEvent         MicroserviceEvent = "social.unblock_chat"
)

func MicroserviceEventValues() []MicroserviceEvent {
	return []MicroserviceEvent{
		TestImageEvent,
		TestMintEvent,

		PaymentsNotifyClientEvent,
		RoomCreatorRoomCreatedEvent,
		RoomCreatorRoomUpdatedEvent,
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

// PaymentsNotifyClientPayload is the payload for the payments.notify_client event.
type PaymentsNotifyClientPayload struct {
	Room    string                 `json:"room"`
	Message map[string]interface{} `json:"message"`
}

func (PaymentsNotifyClientPayload) Type() MicroserviceEvent {
	return PaymentsNotifyClientEvent
}

// RoomCreatorRoomCreatedPayload is the payload for the room_creator.room_created event.
type RoomCreatorRoomCreatedPayload struct {
	Name    string                 `json:"name"`
	Id      string                 `json:"id"`
	Product map[string]interface{} `json:"product"`
}

func (RoomCreatorRoomCreatedPayload) Type() MicroserviceEvent {
	return RoomCreatorRoomCreatedEvent
}

// RoomCreatorRoomUpdatedPayload is the payload for the room_creator.room_updated event.
type RoomCreatorRoomUpdatedPayload struct {
	Id      string                 `json:"id"`
	Product map[string]interface{} `json:"product"`
}

func (RoomCreatorRoomUpdatedPayload) Type() MicroserviceEvent {
	return RoomCreatorRoomUpdatedEvent
}

// RoomSnapshotFirstSnapshotPayload is the payload for the social.block_chat event.
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
