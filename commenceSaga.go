package saga

import "fmt"

const (
	CommenceSagaQueue Queue = "commence_saga"
)

type SagaTitle string

const (
	PurchaseResourceFlow                 SagaTitle = "purchase_resource_flow"
	RankingsUsersReward                  SagaTitle = "rankings_users_reward"
	TransferCryptoRewardToMissionWinner  SagaTitle = "transfer_crypto_reward_to_mission_winner"
	TransferCryptoRewardToRankingWinners SagaTitle = "transfer_crypto_reward_to_ranking_winners"
)

type CommencePayload interface {
	Type() SagaTitle
}

// PurchaseResourceFlowPayload is the payload for the purchase_resource_flow event.
type PurchaseResourceFlowPayload struct {
	UserId     string `json:"userId"`
	ResourceId string `json:"resourceId"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
}

func (PurchaseResourceFlowPayload) Type() SagaTitle {
	return PurchaseResourceFlow
}

type UserReward struct {
	UserId string `json:"userId"`
	Reward string `json:"coins"`
}

// RankingsUsersRewardPayload is the payload for the rankings_users_reward event.
type RankingsUsersRewardPayload struct {
	Rewards []UserReward `json:"rewards"`
}

func (RankingsUsersRewardPayload) Type() SagaTitle {
	return RankingsUsersReward
}

type CryptoRankingWinners struct {
	UserID string `json:"userId"`
	Reward string `json:"reward"`
}

type CompletedCryptoRanking struct {
	WalletAddress string                 `json:"walletAddress"`
	Winners       []CryptoRankingWinners `json:"winners"`
}

// TransferCryptoRewardToMissionWinnerPayload is the payload for the transfer_crypto_reward_to_mission_winner event.
type TransferCryptoRewardToMissionWinnerPayload struct {
	// Wallet address from which rewards will be transferred
	WalletAddress string `json:"walletAddress"`
	// ID of the user who completed the mission
	UserID string `json:"userId"`
	// Amount to be transferred
	Reward string `json:"reward"`
}

func (TransferCryptoRewardToMissionWinnerPayload) Type() SagaTitle {
	return TransferCryptoRewardToMissionWinner
}

// TransferCryptoRewardToRankingWinnersPayload is the payload for the transfer_crypto_reward_to_ranking_winners event.
type TransferCryptoRewardToRankingWinnersPayload struct {
	CompletedCryptoRankings []CompletedCryptoRanking `json:"completedCryptoRankings"`
}

func (TransferCryptoRewardToRankingWinnersPayload) Type() SagaTitle {
	return TransferCryptoRewardToRankingWinners
}

type commenceSaga struct {
	Title   SagaTitle   `json:"title"`
	Payload interface{} `json:"payload"`
}

func CommenceSaga(payload CommencePayload) error {
	channel, err := getSendChannel()
	if err != nil {
		return fmt.Errorf("error getting send channel: %w", err)
	}
	title := payload.Type()
	// Para que este micro pueda realizar pasos del saga y realizar commence_saga ops las queue's deben existir, no es responsabilidad
	// de los micros crear estos recursos, el micro "transactional" debe crear estos recursos -> "queue.CommenceSaga" en commenceSagaListener
	// y "queue.ReplyToSaga" en startGlobalSagaStepListener
	err = send(channel, string(CommenceSagaQueue), commenceSaga{
		Title:   title,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	return nil
}
