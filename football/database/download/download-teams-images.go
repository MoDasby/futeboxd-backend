package main

/* import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/usecase"
)

func main() {
	db := database.InitDatabase()

	teamsRepo := repository.NewTeamRepository(db)

	listTeamsUseCase := usecase.NewListTeamsUseCase(teamsRepo)

	teams, err := listTeamsUseCase.Execute("")
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll("./logos", os.ModePerm); err != nil {
		panic(err)
	}

	for _, team := range teams {
		if team.Logo == "" {
			continue
		}

		file, err := os.Create(fmt.Sprintf("logos/%d.png", team.ID))
		if err != nil {
			panic(err)
		}

		resp, err := http.Get(team.Logo)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(file, resp.Body); err != nil {
			panic(err)
		}

		file.Close()
		resp.Body.Close()
	}
}
*/
