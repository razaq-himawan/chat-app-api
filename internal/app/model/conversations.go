package model

type Conversation struct {
	ID          string `json:"id"`
	MemberOneID string `json:"member_one_id"`
	MemberTwoID string `json:"member_two_id"`
}
