package usecase

import (
	"context"
	"image"
	"io"

	"github.com/modasby/futeboxd-backend/core/internal/upload"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

func (uc *uploadUsecase) Upload(ctx context.Context, img upload.File) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	_, ext, err := image.Decode(img.Content)
	if err != nil {
		if err == image.ErrFormat {
			return errors.NewHTTPErr(
				"formato inválido para imagem. Formatos permitidos: png, jpg, gif",
				400,
				"UPLOAD:USECASE:UPLOAD_FILE:INVALID_IMAGE_TYPE",
			)
		}

		return err
	}

	if _, err := img.Content.Seek(0, io.SeekStart); err != nil {
		return err
	}

	img.Name = user.ID + "." + ext

	url, err := uc.uploadRepo.Put(ctx, &img)
	if err != nil {
		return err
	}

	user.ProfilePicture = url

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
