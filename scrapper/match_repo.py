from db import DB
from models.match import Match
import json

class MatchRepository:
    __db: DB

    def __init__(self, db: DB) -> None:
        self.__db = db

    def upsert(self, match: Match) -> None:
        conn = self.__db.get_connection()

        query = """
            INSERT INTO matches (
                id, match_date, venue, home_team_score, home_team_id,
                away_team_score, away_team_id, 
                note, completed, status_name, competition_name, events
            )
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (id) DO UPDATE SET
                match_date = EXCLUDED.match_date,
                venue = EXCLUDED.venue,
                home_team_score = EXCLUDED.home_team_score,
                away_team_score = EXCLUDED.away_team_score,
                note = EXCLUDED.note,
                completed = EXCLUDED.completed,
                status_name = EXCLUDED.status_name,
                competition_name = EXCLUDED.competition_name,
                events = EXCLUDED.events
        """

        events = json.dumps(match.events)

        values = (
            match.id,
            match.date,
            match.venue,
            match.home_competitor.score,
            match.home_competitor.team_id,
            match.away_competitor.score,
            match.away_competitor.team_id,
            match.note,
            match.completed,
            match.status_name,
            match.competition_name,
            events
        )

        with conn.cursor() as cursor:
            cursor.execute(query, values)
        
        self.__db.release(conn)
