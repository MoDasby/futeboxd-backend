package repository

import "database/sql"

type CommentsRepository struct {
	db *sql.DB
}
