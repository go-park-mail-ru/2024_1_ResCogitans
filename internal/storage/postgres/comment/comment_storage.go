package comment

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentStorage struct {
	db *pgxpool.Pool
}

func NewCommentStorage(db *pgxpool.Pool) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (cs *CommentStorage) GetCommentsBySightID(ctx context.Context, id int) ([]entities.Comment, error) {
	var comments []*entities.Comment

	err := pgxscan.Select(ctx, cs.db, &comments,
		`SELECT f.id, f.user_id, p.username, p.avatar, f.sight_id, f.rating, f.feedback FROM feedback AS f 
    			INNER JOIN profile_data AS p 
    			ON f.user_id = p.user_id 
    			WHERE sight_id =  $1`, id)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var commentsList []entities.Comment
	for _, s := range comments {
		commentsList = append(commentsList, *s)
	}
	return commentsList, nil
}

func (cs *CommentStorage) GetCommentsByUserID(ctx context.Context, userID int) ([]entities.Comment, error) {
	var comments []*entities.Comment

	err := pgxscan.Select(ctx, cs.db, &comments,
		`SELECT f.id, f.user_id, f.sight_id, f.rating, f.feedback FROM feedback AS f 
                WHERE user_id =  $1 `, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}

	var commentsList []entities.Comment
	for _, s := range comments {
		commentsList = append(commentsList, *s)
	}
	return commentsList, nil
}

func (cs *CommentStorage) CreateCommentBySightID(ctx context.Context, sightID int, comment entities.Comment) error {
	_, err := cs.db.Exec(ctx, `INSERT INTO feedback(user_id, sight_id, rating, feedback) VALUES($1, $2, $3, $4)`,
		comment.UserID, sightID, comment.Rating, comment.Feedback)
	return err
}

func (cs *CommentStorage) EditComment(ctx context.Context, commentID int, comment entities.Comment) error {
	_, err := cs.db.Exec(ctx, `UPDATE feedback SET rating = $1, feedback = $2 
                					WHERE id = $3`, comment.Rating, comment.Feedback, commentID)
	return err
}

func (cs *CommentStorage) DeleteComment(ctx context.Context, commentID int) error {
	_, err := cs.db.Exec(ctx, `DELETE FROM feedback WHERE id = $1`, commentID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return err
	}
	return nil
}
