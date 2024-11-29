from database.db import DB
from .models import Team

class TeamRepository:
    __db: DB

    def __init__(self, db: DB) -> None:
        self.__db = db

    def existsById(self, id: int) -> bool:
        query = "SELECT EXISTS(SELECT 1 FROM teams WHERE id = %s)"
        
        result = self.__db.fetch(query, (id, ))

        return result
    
    def create(self, team: Team):
        query = """
            INSERT INTO teams (id, name, abbreviation, color, logo)
		    VALUES (%s, %s, %s, %s, %s)
            ON CONFLICT (id) DO NOTHING
        """

        self.__db.execute(query, (team.id, team.name, team.abbreviation, team.color, team.logo))
