package news

import (
	"context"
	"database/sql"
)

type NewsRepo interface {
	GetRecentNews(ctx context.Context, pageSize, pageIndex int) ([]News, error)
}

type newsRepo struct {
	db *sql.DB
}

func NewNewsRepo(db *sql.DB) NewsRepo {
	return &newsRepo{db: db}
}

func (repo *newsRepo) GetRecentNews(ctx context.Context, pageSize, pageIndex int) ([]News, error) {
	query := `
		SELECT id, title, description, image_link, created_at, link
		FROM news
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET ($2 - 1) * $1
	`

	rows, err := repo.db.QueryContext(ctx, query, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}

	output := make([]News, 0)

	for rows.Next() {
		var news News

		if err := rows.Scan(&news.ID, &news.Title, &news.Description, &news.ImageLink, &news.CreatedAt, &news.Link); err != nil {
			return nil, err
		}

		output = append(output, news)
	}

	return output, nil
}
