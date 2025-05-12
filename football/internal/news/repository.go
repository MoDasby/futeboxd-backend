package news

import (
	"context"
	"database/sql"
	"fmt"
)

type NewsRepo interface {
	GetRecentNews(ctx context.Context, pageSize, pageIndex int, keywords []string) ([]News, error)
}

type newsRepo struct {
	db *sql.DB
}

func NewNewsRepo(db *sql.DB) NewsRepo {
	return &newsRepo{db: db}
}

func (repo *newsRepo) GetRecentNews(ctx context.Context, pageSize, pageIndex int, keywords []string) ([]News, error) {
	query := `
		SELECT id, title, description, image_link, created_at, link
		FROM news
    `

	args := make([]any, 0)
	argIndex := 1

	if len(keywords) > 0 {
		query += "WHERE "
		for i, keyword := range keywords {
			if i > 0 {
				query += " OR "
			}

			query += fmt.Sprintf("unaccent(title) ILIKE '%%' || unaccent($%d) || '%%' OR unaccent(description) ILIKE '%%' || unaccent($%d) || '%%'", argIndex, argIndex)
			args = append(args, keyword)
			argIndex++
		}
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (pageIndex-1)*pageSize)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
