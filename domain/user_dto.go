package domain

// User represents a record from users table
type User struct {
	Id        UUID        `json:"id"`
	CreatedAt Timestamptz `json:"createdAt"`
}

type UserProfile struct {
	UserId   UUID   `json:"userId"`
	FullName string `json:"fullName"`
}

type EmailCredential struct {
	UserId         UUID        `json:"userId"`
	Email          string      `json:"email"`
	HashedPassword string      `json:"-"`
	CreatedAt      Timestamptz `json:"createdAt"`
}

type InputCreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type InputCreatTokenByEmailCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OutputToken struct {
	Token string `json:"token"`
}
