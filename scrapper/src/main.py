from match import MatchRepository, League, Scheduler, MatchService
from espn import Espn
from database.db import PostgresDB
from team import TeamService
from team.team_repo import TeamRepository

db = PostgresDB()

team_repo = TeamRepository(db)
team_service = TeamService(team_repo)

scheduler = Scheduler()

espn = Espn()

matchRepo = MatchRepository(db)
match_service = MatchService(espn, matchRepo, team_service, scheduler)
    
for league in League.list_all():
    match_service.process_league_schedule(league)

scheduler.run()
