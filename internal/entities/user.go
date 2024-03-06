package entities

import (
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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
		Password: "testpassword",
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

func UserValidation(username, password string) bool {
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

func UserDataVerification(username, password string) (int, error) {
	if username == "" || password == "" {
		return http.StatusBadRequest, errors.New("username and password must not be empty")
	}

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return http.StatusBadRequest, errors.New("username and password must not contain only whitespace")
	}

	for _, user := range users {
		if user.Username == username {
			return http.StatusBadRequest, errors.New("username already exists")
		}
	}

	// if !isPasswordComplex(password) {
	// 	return http.StatusBadRequest, errors.New("Password is not complex enough")
	// }

	return http.StatusOK, nil
}

func isPasswordComplex(password string) bool {
	// Соответствует паролю, содержащему как минимум одну цифру, одну заглавную букву, одну строчную букву и имеет длину не менее 8 символов
	complexityRegex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	match, _ := regexp.MatchString(complexityRegex, password)
	return match
}
