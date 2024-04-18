package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (h User) Validate() error {
	return nil
}

type UserProfile struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

func (h UserProfile) Validate() error {
	return nil
}
