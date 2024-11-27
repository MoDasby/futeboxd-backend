from match_repo import MatchRepository
from espn import Espn
from scheduler import Scheduler
from models.match import Match
from datetime import datetime, timezone, timedelta
from functools import partial
from mappers import from_schedule_to_match

class MatchService:
    espn: Espn
    match_repo: MatchRepository
    scheduler: Scheduler

    def __init__(self, espn: Espn, match_repo: MatchRepository, scheduler: Scheduler) -> None:
        self.espn = espn
        self.match_repo = match_repo
        self.scheduler = scheduler

    def process_league_schedule(self, league: str) -> None:
        day_schedule = self.espn.get_day_schedule(league)

        for event in day_schedule:
            match = from_schedule_to_match(event)

            self.match_repo.upsert(match)

            self.schedule_if_necessary(match, league)

    def schedule_if_necessary(self, match: Match, league: str) -> None:
        match_date = datetime.strptime(match.date, '%Y-%m-%dT%H:%MZ').replace(tzinfo=timezone.utc, microsecond=0)
        if match_date > datetime.now(timezone.utc):
            self.scheduler.schedule(match_date, partial(self.process_league_schedule, league))
            return

        # checa se a partida está em andamento
        if not match.completed:
            self.scheduler.schedule(datetime.now(timezone.utc) + timedelta(minutes=5, microseconds=0), partial(self.process_league_schedule, league))
            return