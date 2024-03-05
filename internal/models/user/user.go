package user

import "github.com/pkg/errors"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{
	{
		ID:       1,
		Username: "sanboy",
		Password: "12345",
	},
	{
		ID:       2,
		Username: "alex",
		Password: "qwerty",
	},
	{
		ID:       3,
		Username: "fil",
		Password: "12345qwerty",
	},
}

func (h User) Validate() error {
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

func UserValidation(username, password string) bool {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false
	}

	return user.Password == password
}
