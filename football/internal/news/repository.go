package news

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-api/services/football/pkg/sq"
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
	query := sq.Query.
		Select("id, title, description, image_link, created_at, link").
		From("news")

	if len(keywords) > 0 {
		for _, keyword := range keywords {
			query = query.Where(
				"unaccent(title) ILIKE '%%' || unaccent(?) || '%%' OR unaccent(description) ILIKE '%%' || unaccent(?) || '%%'",
				keyword, keyword,
			)
		}
	}

	query = query.OrderBy("created_at DESC").Limit(uint64(pageSize)).Offset(uint64(pageIndex) - 1)

	rows, err := query.RunWith(repo.db).QueryContext(ctx)
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
