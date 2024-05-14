package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/redis/csrf"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CSRFInterface interface {
	CreateToken(userID int) (string, error)
	UpdateToken(userID int) (string, error)
	SetToken(token string, w http.ResponseWriter)
	CompareToken(encoded string, userID int) error
	ClearToken(userID int) error
}

type CSRFUseCase struct {
	CSRFStorage *csrf.CSRFStorage
}

func NewCSRFUseCase(storage *csrf.CSRFStorage) *CSRFUseCase {
	return &CSRFUseCase{
		CSRFStorage: storage,
	}
}

func (a *CSRFUseCase) CreateToken(userID int) (string, error) {
	token := uuid.New().String()
	key, err := generateRandomKey()
	if err != nil {
		return "", err
	}
	err = a.CSRFStorage.SaveToken(token, key, userID)

	encryptedToken, err := encrypt(key, []byte(token))
	if err != nil {
		return "", fmt.Errorf("error encrypting token: %w", err)
	}
	base64Token := base64.StdEncoding.EncodeToString(encryptedToken)
	return base64Token, nil
}

func (a *CSRFUseCase) UpdateToken(userID int) (string, error) {
	token := uuid.New().String()
	key, err := a.CSRFStorage.GetKey(userID)
	if err != nil {
		return "", fmt.Errorf("error getting key: %w", err)
	}
	err = a.CSRFStorage.SaveToken(token, key, userID)
	encryptedToken, err := encrypt(key, []byte(token))
	if err != nil {
		return "", fmt.Errorf("error encrypting token: %w", err)
	}
	base64Token := base64.StdEncoding.EncodeToString(encryptedToken)
	return base64Token, nil
}

func (a *CSRFUseCase) SetToken(token string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "X-CSRF-Token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
}

func (a *CSRFUseCase) CompareToken(base64Token string, userID int) error {
	key, err := a.CSRFStorage.GetKey(userID)
	if err != nil {
		return fmt.Errorf("error getting key: %w", err)
	}

	encoded, err := base64.StdEncoding.DecodeString(base64Token)
	if err != nil {
		return fmt.Errorf("error decoding token: %w", err)
	}
	decryptedToken, err := decrypt(key, encoded)
	if err != nil {
		return fmt.Errorf("error encrypting token: %w", err)
	}

	currentToken, err := a.CSRFStorage.GetToken(key)
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}

	if currentToken != string(decryptedToken) {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func (a *CSRFUseCase) ClearToken(userID int) error {
	key, err := a.CSRFStorage.GetKey(userID)
	if err != nil {
		return fmt.Errorf("error getting key: %w", err)
	}
	return a.CSRFStorage.DeleteToken(key)
}

// generateRandomKey создает случайный ключ шифрования AES.
func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func encrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Создаем вектор инициализации (IV) для шифрования.
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Шифруем данные.
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Проверяем, что текст шифрования достаточно велик для хранения IV.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Расшифровываем данные.
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
