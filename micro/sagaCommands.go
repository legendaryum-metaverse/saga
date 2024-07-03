package micro

type StepCommand = string

type AvailableMicroservices string

// IsValid checks if the provided value is a valid AvailableMicroservices.
func (m AvailableMicroservices) IsValid() bool {
	switch m {
	case Auth, RapidMessaging, TestImage, TestMint, Payments, RoomInventory, RoomSnapshot, RoomCreator, Showcase, Social, Storage:
		return true
	}
	return false
}

const (
	Auth           AvailableMicroservices = "auth"
	RapidMessaging AvailableMicroservices = "rapid-messaging"
)

// image.
const (
	TestImage          AvailableMicroservices = "test-image"
	CreateImageCommand StepCommand            = "create_image"
	UpdateTokenCommand StepCommand            = "update_token"
)

// mint.
const (
	TestMint         AvailableMicroservices = "test-mint"
	MintImageCommand StepCommand            = "mint_image"
)

// payments.
const (
	Payments                            AvailableMicroservices = "payments"
	ResourcePurchasedDeductCoinsCommand StepCommand            = "resource_purchased:deduct_coins"
)

// room-inventory.
const (
	RoomInventory                AvailableMicroservices = "room-inventory"
	SavePurchasedResourceCommand StepCommand            = "resource_purchased:save_purchased_resource"
)

// room-snapshot.
const (
	RoomSnapshot AvailableMicroservices = "room-snapshot"
)

// room-creator.
const (
	RoomCreator                     AvailableMicroservices = "room-creator"
	UpdateIslandRoomTemplateCommand StepCommand            = "update_island_room_template"
)

// legend-showcase.
const (
	Showcase                      AvailableMicroservices = "legend-showcase"
	RandomizeIslandPvImageCommand StepCommand            = "randomize_island_pv_image"
)

// social.
const (
	Social                 AvailableMicroservices = "social"
	UpdateUserImageCommand StepCommand            = "update_user:image"
	NotifyClientCommand    StepCommand            = "notify_client"
)

// storage.
const (
	Storage           AvailableMicroservices = "legend-storage"
	UpdateFileCommand StepCommand            = "update_file"
)
