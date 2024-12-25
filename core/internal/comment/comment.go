package comment

import (
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type Comment struct {
	ID        int64
	Author    *user.User
	ParentID  int64
	Content   string
	LikeCount int
	IsLiked   bool
	CreatedAt time.Time
}

func NewComment(
	author *user.User,
	parentID int64,
	content string,
) (*Comment, error) {
	comment := &Comment{
		Author:   author,
		ParentID: parentID,
		Content:  content,
	}

	if err := comment.Validate(); err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *Comment) Validate() error {
	if c.ParentID == 0 {
		return errors.NewHTTPErr("review inválida", 400, "DOMAIN:COMMENT:VALIDATE:INVALID_PARENT_ID")
	}

	if c.Content == "" {
		return errors.NewHTTPErr("não pode fazer um comentário vazio", 400, "DOMAIN:COMMENT:VALIDATE:NO_CONTENT")
	}

	return nil
}
