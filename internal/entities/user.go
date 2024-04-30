package entities

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h UserRequest) Validate() error {
	return nil
}

type User struct {
	ID       int
	Username string
	Password string
	Salt     string
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type ProfileRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

func (h ProfileRequest) Validate() error {
	return nil
}

type UserProfile struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}
