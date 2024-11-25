from espn import Espn
from models.leagues import League
from match_repo import MatchRepository
from service import upsert_league_schedule
from scheduler import Scheduler
from db import PostgresDB
from utils import get_execution_date

db = PostgresDB()
scheduler = Scheduler()
espn = Espn()
matchRepo = MatchRepository(db)

def init(league: str, espn: Espn) -> None:
    must_continue = upsert_league_schedule(matchRepo, espn, league)
    
    if not must_continue:
        scheduler.schedule(get_execution_date("morning"), upsert_league_schedule, league, espn)
        return
    
    scheduler.schedule(get_execution_date(), upsert_league_schedule, league, espn)

for league in League.list_all():
    init(league, espn)

scheduler.run()