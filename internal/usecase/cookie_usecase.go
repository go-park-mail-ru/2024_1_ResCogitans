package usecase

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage/cookie"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CookieInterface interface {
	GetCookie(*http.Request) error
	CompareCSRF(*http.Request) error
	Set(http.ResponseWriter) error
	ChangeCSRF(*http.Request) (string, error)
}

type CookieUseCase struct {
	CookieStorage cookie.StorageInterface
}

func NewCookieUseCase(storage cookie.StorageInterface) CookieInterface {
	return &CookieUseCase{
		CookieStorage: storage,
	}
}

func (cu *CookieUseCase) GetCookie(r *http.Request) error {
	cookie, err := r.Cookie("cookie_id")
	if err != nil {
		return err
	}
	var cookieID string
	if err = CookieHandler.Decode("cookie_id", cookie.Value, &cookieID); err == nil {
		return cu.CookieStorage.GetCookie(cookieID)
	}
	return err
}

func (cu *CookieUseCase) CompareCSRF(r *http.Request) error {
	cookie, err := r.Cookie("cookie_id")
	if err != nil {
		return err
	}

	var cookieID string
	if err = CookieHandler.Decode("cookie_id", cookie.Value, &cookieID); err != nil {
		return err
	}
	token, err := cu.CookieStorage.GetCSRF(cookieID)
	if err != nil {
		return err
	}

	tokenFromRequest := r.Header.Get("X-CSRF-Token")
	if token != tokenFromRequest {
		return errors.New("invalid token")
	}
	return nil
}

func (cu *CookieUseCase) Set(w http.ResponseWriter) error {
	cookieID := uuid.New().String()
	csrfToken := uuid.New().String()
	cu.CookieStorage.Set(cookieID, csrfToken)
	encoded, err := CookieHandler.Encode("cookie_id", cookieID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "cookie_id",
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	w.Header().Set("X-CSRF-Token", csrfToken)
	return nil
}

func (cu *CookieUseCase) ChangeCSRF(r *http.Request) (string, error) {
	cookie, err := r.Cookie("cookie_id")
	if err != nil {
		return "", err
	}
	var cookieID string
	if err = CookieHandler.Decode("cookie_id", cookie.Value, &cookieID); err != nil {
		return "", err
	}
	csrfToken := uuid.New().String()
	err = cu.CookieStorage.SetCSRF(cookieID, csrfToken)
	if err != nil {
		return "", err
	}
	return csrfToken, nil
}
