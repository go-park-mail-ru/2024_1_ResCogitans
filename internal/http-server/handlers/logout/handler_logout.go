package logout

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"net/http"
)

type Logout struct{}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *Logout) LogOut(ctx context.Context, _ entities.User) (Response, error) {
	w, ok := ctx.Value("responseWriter").(http.ResponseWriter)
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: "failed getting responseWriter"}, nil
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return Response{Status: http.StatusOK}, nil
}
