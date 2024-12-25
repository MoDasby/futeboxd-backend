package team

type Repository interface {
	FindOneById(teamID int64) (*Team, error)
	FindAll(name string, pageSize, pageIndex int) ([]Team, error)
}
