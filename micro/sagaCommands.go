package micro

type StepCommand = string

type AvailableMicroservices string

// IsValid checks if the provided value is a valid AvailableMicroservices.
func (m AvailableMicroservices) IsValid() bool {
	switch m {
	case Auth, RapidMessaging, TestImage, TestMint, Payments, RoomInventory, RoomSnapshot, RoomCreator, Showcase, Social, Storage, Missions:
		return true
	}
	return false
}

// image mock microservice.
const (
	TestImage          AvailableMicroservices = "test-image"
	CreateImageCommand StepCommand            = "create_image"
	UpdateTokenCommand StepCommand            = "update_token"
)

// mint mock microservice.
const (
	TestMint         AvailableMicroservices = "test-mint"
	MintImageCommand StepCommand            = "mint_image"
)

const (
	Auth                     AvailableMicroservices = "auth"
	UpdateUserPictureCommand StepCommand            = "update_user:picture"
)

const (
	Missions AvailableMicroservices = "legend-missions"
)

const (
	Payments                            AvailableMicroservices = "payments"
	ResourcePurchasedDeductCoinsCommand StepCommand            = "resource_purchased:deduct_coins"
)

const (
	RapidMessaging AvailableMicroservices = "rapid-messaging"
)

const (
	RoomInventory                AvailableMicroservices = "room-inventory"
	SavePurchasedResourceCommand StepCommand            = "resource_purchased:save_purchased_resource"
)

const (
	RoomSnapshot AvailableMicroservices = "room-snapshot"
)

const (
	RoomCreator                     AvailableMicroservices = "room-creator"
	UpdateIslandRoomTemplateCommand StepCommand            = "update_island_room_template"
)

const (
	Showcase                      AvailableMicroservices = "legend-showcase"
	RandomizeIslandPvImageCommand StepCommand            = "randomize_island_pv_image"
)

const (
	Social                 AvailableMicroservices = "social"
	UpdateUserImageCommand StepCommand            = "update_user:image"
)

const (
	Storage           AvailableMicroservices = "legend-storage"
	UpdateFileCommand StepCommand            = "update_file"
)
