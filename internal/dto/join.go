package dto

type JoinRequest struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
}

type JoinResponse struct {
	CanJoin bool `json:"can_join"`
}
