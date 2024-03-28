package entities

import (
	"regexp"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
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
		ID:       1,
		Username: "test_user",
		Password: "TestPassword123",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed creating hash password for test user")
	}

	testUser.Password = string(hashPassword)
	users = append(users, testUser)
}

func GetUserByUsername(username string) (*User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

func GetUserByID(userID int) (*User, error) {
	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

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

func UserExists(username string, password string) error {
	for _, user := range users {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return errors.New("Wrong password")
			}
			return nil
		}
	}
	return errors.New("User with this username doesn't exists")
}

func UserDataVerification(username, password string) error {
	if !ValidateUsername(username) {
		return errors.New("Username doesn't meet requirements")
	}
	if !ValidatePassword(password) {
		return errors.New("Password doesn't meet requirements")
	}

	return nil
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{2,}$`, username)
	return matched
}

func ValidatePassword(password string) bool {
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasMinLength := len(password) >= 8
	hasMaxLength := len(password) <= 40
	return hasDigit && hasUppercase && hasLowercase && hasMinLength && hasMaxLength
}

func ChangeData(userID int, username string, password string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	err = UserDataVerification(username, password)
	if err != nil {
		return nil, err
	}

	if user.Username != username {
		user.Username = username
	}
	if user.Password != password {
		user.Password = string(hashPassword)
	}

	for i := range users {
		if users[i].ID == userID {
			users[i] = *user
			break
		}
	}
	return user, nil
}

func DeleteUser(userID int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return errors.New("User not found")
}
