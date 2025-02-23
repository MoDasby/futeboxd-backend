package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/modasby/futeboxd-backend/core/internal/upload"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type uploadUsecase struct {
	uploadRepo upload.Repository
	userRepo   user.Repository
}

func NewUploadUsecase(
	uploadRepo upload.Repository,
	userRepo user.Repository,
) upload.Usecase {
	return &uploadUsecase{uploadRepo: uploadRepo, userRepo: userRepo}
}

var allowedTypes = map[string]bool{"image/jpeg": true, "image/png": true}

func calculateSize(file io.Seeker) (int64, error) {
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	return size, nil
}

func validateImage(ctx context.Context, file io.ReadSeeker) (*mimetype.MIME, error) {
	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, err
	}

	if ok := allowedTypes[mimeType.String()]; !ok {
		return nil, &errors.HTTPErr{
			Msg:        "Arquivo não é uma imagem válida, imagens suportadas: png, jpg",
			Code:       http.StatusBadRequest,
			StackTrace: errors.CaptureStackTrace(),
			Context:    "UPLOAD:USECASE:UPLOAD_FILE:INVALID_IMAGE_TYPE",
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   nil,
		}
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return mimeType, nil
}

func (uc *uploadUsecase) Upload(ctx context.Context, fileRaw io.ReadSeeker) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	mime, err := validateImage(ctx, fileRaw)
	if err != nil {
		return err
	}

	size, err := calculateSize(fileRaw)
	if err != nil {
		return err
	}

	file := &upload.File{
		Name:        fmt.Sprintf("%s%s", user.ID, mime.Extension()),
		Size:        size,
		Content:     fileRaw,
		ContentType: mime.String(),
	}

	url, err := uc.uploadRepo.Put(ctx, file)
	if err != nil {
		return err
	}

	user.ProfilePicture = url

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
