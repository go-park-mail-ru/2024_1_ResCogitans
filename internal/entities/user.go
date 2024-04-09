package entities

import (
	"regexp"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID      int    `json:"id"`
	Email   string `json:"username"`
	Passwrd string `json:"password"`
}

var (
	users []User
	mu    sync.Mutex
)

func (h User) Validate() error {
	return nil
}

func init() {
	testUser := User{
		ID:      1,
		Email:   "testuser",
		Passwrd: "testpassword",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Passwrd), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed creating hash password for test user")
	}

	testUser.Passwrd = string(hashPassword)
	users = append(users, testUser)
}

// func GetUserByUsername(username string) (*User, error) {
// 	for _, user := range users {
// 		if user.Username == username {
// 			return &user, nil
// 		}
// 	}
// 	return nil, errors.New("User not found")
// }

// func IsAuthenticated(username, password string) bool {
// 	user, err := GetUserByUsername(username)
// 	if err != nil {
// 		return false
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	return err == nil
// }

func CreateUser(username, password string) (User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("Failed creating hash password")
	}
	var lastUser User

	mu.Lock()
	defer mu.Unlock()

	if len(users) != 0 {
		lastUser = users[len(users)-1]
	}

	newUser := User{
		lastUser.ID + 1,
		username,
		string(hashPassword),
	}

	users = append(users, newUser)

	return newUser, nil
}

func UserDataVerification(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password must not be empty")
	}

	// for _, user := range users {
	// 	if user.Username == username {
	// 		return errors.New("username already exists")
	// 	}
	// }

	if !ValidatePassword(password) {
		return errors.New("password is not complex")
	}

	return nil
}

func ValidatePassword(password string) bool {
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasMinLength := len(password) >= 8

	return hasDigit && hasUppercase && hasLowercase && hasMinLength
}
