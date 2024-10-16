package saga

const (
	CommenceSagaQueue Queue = "commence_saga"
)

type SagaTitle string

const (
	PurchaseResourceFlow SagaTitle = "purchase_resource_flow"
	RankingsUsersReward  SagaTitle = "rankings_users_reward"
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

type commenceSaga struct {
	Title   SagaTitle   `json:"title"`
	Payload interface{} `json:"payload"`
}

func (t *Transactional) CommenceSaga(payload CommencePayload) error {
	title := payload.Type()
	err := send(t.sendChannel, string(CommenceSagaQueue), commenceSaga{
		Title:   title,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	return nil
}
