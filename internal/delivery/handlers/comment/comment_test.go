package comment_test

import (
	"context"
	"testing"

	comment "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/handlers/sight"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	"github.com/stretchr/testify/assert"
)

// Get comments
func TestGetComments(t *testing.T) {
	handler := &comment.SightsHandler{}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	resp, err := handler.GetSightByID(ctx, entities.Sight{})

	assert.NotEmpty(t, resp.Comms)

	assert.Nil(t, err)
}

// Create comment
func TestCreateComment(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.CreateComment(ctx, comment)

	assert.Nil(t, err)
}

func TestCreateCommentIDNotInt(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["id"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.CreateComment(ctx, comment)

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Edit comment
func TestEditComment(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["cid"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.EditComment(ctx, comment)

	assert.Nil(t, err)
}

func TestEditCommentIDNotInt(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["cid"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.EditComment(ctx, comment)

	assert.EqualError(t, err, "cannot parsing not integer")
}

// Delete comment
func TestDeleteComment(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["cid"] = "1"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteComment(ctx, comment)

	assert.Nil(t, err)
}

func TestDeleteCommentIDNotInt(t *testing.T) {
	handler := &CommentHandler{}
	comment := entities.Comment{
		Rating:   4,
		Feedback: "Siiiuuu",
		UserID:   1,
	}

	ctx := context.Background()
	param := make(map[string]string)
	param["cid"] = "ok"
	ctx = wrapper.SetPathParamsToCtx(ctx, param)

	_, err := handler.DeleteComment(ctx, comment)

	assert.EqualError(t, err, "cannot parsing not integer")
}
