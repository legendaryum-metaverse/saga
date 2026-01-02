package micro

type StepCommand = string

type AvailableMicroservices string

// IsValid checks if the provided value is a valid AvailableMicroservices.
func (m AvailableMicroservices) IsValid() bool {
	switch m {
	case TestImage, TestMint,
		Auth,
		Blockchain,
		Missions,
		Rankings,
		SendEmail, Showcase, Social, Storage,
		AuditEda, Billing:
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
	Blockchain                    AvailableMicroservices = "blockchain"
	TransferMissionRewardToWinner StepCommand            = "crypto_reward:transfer_mission_reward_to_winner"
	TransferRewardToWinners       StepCommand            = "crypto_reward:transfer_reward_to_winners"
)

const (
	Missions AvailableMicroservices = "legend-missions"
)

const (
	Rankings AvailableMicroservices = "rankings"
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
	Storage           AvailableMicroservices = "legend-storage"
	UploadFileCommand StepCommand            = "upload_file"
)

const (
	AuditEda AvailableMicroservices = "audit-eda"
)

const (
	Billing                           AvailableMicroservices = "billing"
	RefundPaymentCommand              StepCommand            = "refund_payment"
	CancelSubscriptionCommand         StepCommand            = "cancel_subscription"
	CreateSubscriptionScheduleCommand StepCommand            = "create_subscription_schedule"
)
