package repository

import (
	"bytes"
	"context"
	goErrors "errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/internal/upload"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

const (
	BUCKET_NAME string = "futeboxd-profile-pictures"
)

type UploadRepository struct {
	s3Client *s3.Client
	cfg      *config.AWS
}

func NewUploadRepository(client *s3.Client, cfg *config.AWS) upload.Repository {
	return &UploadRepository{s3Client: client, cfg: cfg}
}

func (repo *UploadRepository) Put(ctx context.Context, input *upload.File) (string, error) {
	log.Print(repo.cfg.S3Endpoint != "")
	if _, err := repo.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(BUCKET_NAME),
		Key:         aws.String(input.Name),
		Body:        input.Content,
		ContentType: aws.String(input.ContentType),
		//ACL:         types.ObjectCannedACLPublicRead,
	}); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BUCKET_NAME, repo.cfg.Region, input.Name)

	return url, nil
}

func (repo *UploadRepository) Get(ctx context.Context, filename string) (*upload.File, error) {
	output, err := repo.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if goErrors.As(err, &noSuchKey) {
			return nil, &errors.HTTPErr{
				Msg:        "Imagem não encontrada",
				Code:       http.StatusNotFound,
				Context:    "UPLOAD:REPOSITORY:GET:IMAGE_NOT_FOUND",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  utils.GetTraceIDFromCtx(ctx),
				Timestamp:  time.Now().UTC(),
				Original:   err,
			}
		}

		return nil, err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, output.Body); err != nil {
		return nil, err
	}

	return &upload.File{
		Name:        filename,
		ContentType: *output.ContentType,
		Size:        int64(buf.Len()),
		Content:     bytes.NewReader(buf.Bytes()),
	}, nil
}
