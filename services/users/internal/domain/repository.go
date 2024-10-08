package domain

type UserRepository interface {
	FindOneByID(id string) (*User, error)
	FindOneByCredential(credential string) (*User, error)
	AddUser(user *User) error
	Exists(username, email string) (bool, error)
	FindOneByIdOrUsername(username string) (*User, error)
	FindBatchByID(ids []string) ([]User, error)
}

type SessionRepository interface {
	FindOneByToken(token string) (*Session, error)
	AddSession(session *Session) (*Session, error)
	DeleteSession(id string) error
}
