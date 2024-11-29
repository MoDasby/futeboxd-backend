from database.db import DB
from .match import Match
import json

class MatchRepository:
    __db: DB

    def __init__(self, db: DB) -> None:
        self.__db = db

    def upsert(self, match: Match) -> None:
        query = """
            INSERT INTO matches (
                id, match_date, venue, home_team_score, home_team_id,
                away_team_score, away_team_id, 
                note, completed, status_name, competition_name, events,
                home_team_rosters, away_team_rosters
            )
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (id) DO UPDATE SET
                match_date = EXCLUDED.match_date,
                venue = EXCLUDED.venue,
                home_team_score = EXCLUDED.home_team_score,
                away_team_score = EXCLUDED.away_team_score,
                note = EXCLUDED.note,
                completed = EXCLUDED.completed,
                status_name = EXCLUDED.status_name,
                competition_name = EXCLUDED.competition_name,
                events = EXCLUDED.events,
                home_team_rosters = EXCLUDED.home_team_rosters,
                away_team_rosters = EXCLUDED.away_team_rosters
        """

        events = json.dumps(match.events)

        values = (
            match.id,
            match.date,
            match.venue,
            match.home_competitor.score,
            match.home_competitor.team.id,
            match.away_competitor.score,
            match.away_competitor.team.id,
            match.note,
            match.completed,
            match.status_name,
            match.competition_name,
            events,
            match.home_competitor.roster,
            match.away_competitor.roster
        )

        self.__db.execute(query, values)
