package dto

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

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type CreatTokenByEmailCredentialInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
