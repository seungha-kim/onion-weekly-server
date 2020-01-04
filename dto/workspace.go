package dto

type CreateWorkspaceInput struct {
	Name string `json:"name"`
}

type Workspace struct {
	ID        UUID        `json:"id"`
	Name      string      `json:"name"`
	CreatedBy UUID        `json:"-"`
	CreatedAt Timestamptz `json:"createdAt"`
}
