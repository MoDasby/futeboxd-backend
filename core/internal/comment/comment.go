package comment

import (
	"context"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
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
	ctx context.Context,
	author *user.User,
	parentID int64,
	content string,
) (*Comment, error) {
	comment := &Comment{
		Author:   author,
		ParentID: parentID,
		Content:  content,
	}

	if err := comment.Validate(ctx); err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *Comment) Validate(ctx context.Context) error {
	if c.ParentID == 0 {
		return &errors.HTTPErr{
			Msg:        "Review inválida",
			Code:       http.StatusBadRequest,
			Context:    "COMMENT:DOMAIN:VALIDATE:INVALID_PARENT_ID",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
	}
	if c.Content == "" {
		return &errors.HTTPErr{
			Msg:        "Não pode fazer um comentário vazio",
			Code:       http.StatusBadRequest,
			Context:    "COMMENT:DOMAIN:VALIDATE:NO_CONTENT",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
	}

	return nil
}
