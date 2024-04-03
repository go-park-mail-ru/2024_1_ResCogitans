package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserResponse struct {
	Username string `json:"username"`
}

func (h User) Validate() error {
	return nil
}
