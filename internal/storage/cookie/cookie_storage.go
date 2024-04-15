package cookie

import (
	"sync"

	"github.com/pkg/errors"
)

type StorageInterface interface {
	GetCookie(string) error
	GetCSRF(string) (string, error)
	Set(string, string)
	SetCSRF(string, string) error
}

type CookieStorage struct {
	Store     []string
	CSRFStore map[string]string
	mu        sync.Mutex
}

func NewCookieStorage() *CookieStorage {
	return &CookieStorage{
		Store:     make([]string, 0),
		CSRFStore: make(map[string]string),
		mu:        sync.Mutex{},
	}
}

func (cs *CookieStorage) GetCSRF(cookie string) (string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	token, ok := cs.CSRFStore[cookie]
	if !ok {
		return "", errors.New("invalid cookie")
	}
	return token, nil
}

func (cs *CookieStorage) GetCookie(cookie string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.Store {
		if c == cookie {
			return nil
		}
	}
	return errors.New("Cookie not found")
}

func (cs *CookieStorage) Set(cookie, token string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.Store = append(cs.Store, cookie)
	cs.CSRFStore[cookie] = token
}

func (cs *CookieStorage) SetCSRF(cookie, token string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	_, ok := cs.CSRFStore[cookie]
	if !ok {
		return errors.New("invalid cookie")
	}
	cs.CSRFStore[cookie] = token
	return nil
}
