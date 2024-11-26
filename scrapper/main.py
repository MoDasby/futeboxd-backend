from espn import Espn
from models.leagues import League
from match_repo import MatchRepository
from service import MatchService
from scheduler import Scheduler
from db import PostgresDB

db = PostgresDB()
scheduler = Scheduler()
espn = Espn()
matchRepo = MatchRepository(db)
match_service = MatchService(espn, matchRepo, scheduler)
    
for league in League.list_all():
    match_service.process_league_schedule(league)

scheduler.run()