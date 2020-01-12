package dto

type Recurring struct {
	Id          UUID        `json:"id"`
	WorkspaceId UUID        `json:"workspaceId"`
	Title       string      `json:"title"`
	Interval    uint64      `json:"interval"`
	CreatedAt   Timestamptz `json:"createdAt"`
}

type RecurringRecord struct {
	Id               UUID        `json:"id"`
	Description      string      `json:"description"`
	ActorId          UUID        `json:"actorId"`
	IntervalSnapshot uint64      `json:"intervalSnapshot"`
	CreatedAt        Timestamptz `json:"createdAt"`
	RecurringId      UUID        `json:"recurringId"`
}

type CreateRecurringInput struct {
	Title    string `json:"title"`
	Interval uint64 `json:"interval"`
}

type UpdateRecurringInput struct {
	Title    string `json:"title"`
	Interval uint64 `json:"interval"`
}

type CreateRecurringRecordInput struct {
	Description string `json:"description"`
}
