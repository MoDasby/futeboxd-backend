from .models import Team
from .team_repo import TeamRepository

class TeamService:
    __teamRepo: TeamRepository

    def __init__(self, teamRepo: TeamRepository) -> None:
        self.__teamRepo = teamRepo

    def createIfNotExists(self, team: Team):
        self.__teamRepo.create(team)