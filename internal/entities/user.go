package entities

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

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

func CreateUser(username, password string) (*User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Failed creating hash password")
	}
	var lastUser User

	if len(users) == 0 {
		lastUser = User{}
	} else {
		lastUser = users[len(users)-1]
	}

	newUser := User{
		lastUser.ID + 1,
		username,
		string(hashPassword),
	}

	users = append(users, newUser)

	fmt.Printf("Created new user: ID=%d, Username=%s\n", newUser.ID, newUser.Username)

	return &newUser, nil
}

func UserDataVerification(username, password string) error {
	if username == "" || password == "" {
		return errors.New("Data must not be empty")
	}
	return nil
}
