package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	errorsTypes "github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindBatchByID(ids []string) ([]domain.User, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)) // slice de argumentos

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	placeholderStr := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(`
		SELECT u.id::text, u.username, u.email
		FROM users u 
		WHERE u.id IN (%s)
	`, placeholderStr)

	rows, err := r.db.Query(query, args...) // Use args com o operador variadic
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}

		fmt.Println(user)

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) FindOneByID(id string) (*domain.User, error) {
	query := `
		SELECT u.id::text, u.username, u.email FROM users u
		WHERE u.id = $1
	`

	rows := r.db.QueryRow(query, id)

	var user *domain.User

	var userId, username, email sql.NullString

	if err := rows.Scan(&userId, &username, &email); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewErrNotFound("usuário não encontrado")
		}

		return nil, err
	}

	user = &domain.User{ID: userId.String, Username: username.String, Email: email.String}

	return user, nil
}

func (r *userRepository) FindOneByIdOrUsername(identificator string) (*domain.User, error) {
	query := `
		SELECT u.id::text, u.username, u.email, u.favorite_team FROM users u
		WHERE LOWER(u.username) = LOWER($1) OR u.id::text = $1
	`

	rows := r.db.QueryRow(query, identificator)

	var user *domain.User

	var userId, username, email sql.NullString
	var favoriteTeamID sql.NullInt64

	if err := rows.Scan(&userId, &username, &email, &favoriteTeamID); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewErrNotFound("usuário não encontrado")
		}

		return nil, err
	}

	user = &domain.User{
		ID:             userId.String,
		Username:       username.String,
		Email:          email.String,
		FavoriteTeamID: favoriteTeamID.Int64,
	}

	return user, nil
}

func (r *userRepository) FindOneByCredential(credential string) (*domain.User, error) {
	query := `
		SELECT u.id::text, u.username, u.email, u.password FROM users u
		WHERE LOWER(u.email) = LOWER($1) OR LOWER(u.username) = LOWER($1)
	`
	fmt.Printf("credential: %s", credential)
	rows := r.db.QueryRow(query, credential)

	var user *domain.User

	var id, username, email, password sql.NullString

	if err := rows.Scan(&id, &username, &email, &password); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewErrNotFound("usuario não encontrado")
		}

		return nil, err
	}

	user = &domain.User{ID: id.String, Username: username.String, Email: email.String, Password: password.String}

	return user, nil
}

func (r *userRepository) AddUser(user *domain.User) error {
	query := `
		INSERT INTO users(username, email, password, favorite_team)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.FavoriteTeamID)

	return err
}

func (r *userRepository) Exists(username, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($2))
	`

	row := r.db.QueryRow(query, username, email)

	var exists *sql.NullBool

	if err := row.Scan(&exists); err != nil {

		return false, err
	}

	return exists.Bool, nil
}
