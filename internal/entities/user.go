package entities

import (
	"fmt"
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
		Username: "testuser",
		Password: "Testpassword123",
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
	return nil, fmt.Errorf("can't find user by ID")
}

func IsAuthenticated(username, password string) bool {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
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

func UserExists(username string) error {
	for _, user := range users {
		if user.Username == username {
			return errors.New("username already exists")
		}
	}

	return nil
}

func UserDataVerification(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password must not be empty")
	}

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
