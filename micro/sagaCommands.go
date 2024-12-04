package micro

type StepCommand = string

type AvailableMicroservices string

// IsValid checks if the provided value is a valid AvailableMicroservices.
func (m AvailableMicroservices) IsValid() bool {
	switch m {
	case TestImage, TestMint,
		Auth,
		Blockchain,
		Coins,
		Missions,
		Rankings,
		RapidMessaging, RoomCreator, RoomInventory, RoomSnapshot,
		SendEmail, Showcase, Social, SocialMediaRooms, Storage:
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
	Auth              AvailableMicroservices = "auth"
	CreateUserCommand StepCommand            = "create_user"
)

const (
	Blockchain              AvailableMicroservices = "blockchain"
	TransferRewardToWinners StepCommand            = "crypto_reward:transfer_reward_to_winners"
)

const (
	Coins                               AvailableMicroservices = "coins"
	ResourcePurchasedDeductCoinsCommand StepCommand            = "resource_purchased:deduct_coins"
	RankingsRewardCoinsCommand          StepCommand            = "rankings_users_reward:reward_coins"
)

const (
	Missions AvailableMicroservices = "legend-missions"
)

const (
	Rankings AvailableMicroservices = "rankings"
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
	SendEmail AvailableMicroservices = "legend-send-email"
)

const (
	Showcase                      AvailableMicroservices = "legend-showcase"
	RandomizeIslandPvImageCommand StepCommand            = "randomize_island_pv_image"
)

const (
	Social                  AvailableMicroservices = "social"
	UpdateUserImageCommand  StepCommand            = "update_user:image"
	CreateSocialUserCommand StepCommand            = "create_social_user"
)

const (
	SocialMediaRooms AvailableMicroservices = "social-media-rooms"
)

const (
	Storage           AvailableMicroservices = "legend-storage"
	UploadFileCommand StepCommand            = "upload_file"
)
