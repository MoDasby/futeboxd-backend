package team

type Repository interface {
	FindOneById(teamID int64) (*Team, error)
	FindByName(name string, pageSize, pageIndex int) ([]Team, error)
}
