from espn import Espn
from .match_repo import MatchRepository
from .scheduler import Scheduler
from .match import Match
from .mappers import from_schedule_to_match
from typing import List
from datetime import datetime, timezone, timedelta
from functools import partial
from team import TeamService
import json

class MatchService:
    __espn: Espn
    __match_repo: MatchRepository
    __team_service: TeamService
    __scheduler: Scheduler

    def __init__(self, espn: Espn, match_repo: MatchRepository, team_service: TeamService, scheduler: Scheduler) -> None:
        self.__espn = espn
        self.__match_repo = match_repo
        self.__team_service = team_service
        self.__scheduler = scheduler

    def process_league_schedule(self, league: str) -> None:
        day_schedule = self.__espn.get_day_schedule(league)
        matches: List[Match] = []

        for event in day_schedule.events:
            match = from_schedule_to_match(event)
            match.home_competitor.roster = "[]"
            match.away_competitor.roster = "[]"

            self.__team_service.createIfNotExists(match.home_competitor.team)

            self.__team_service.createIfNotExists(match.away_competitor.team)

            if match.status_name != "STATUS_SCHEDULED":
                matchSummary = self.__espn.get_match_summary(match.id)

                for roster in matchSummary.rosters:
                    if roster.homeAway == "home":
                        match.home_competitor.roster = json.dumps([item.to_dict() for item in roster.roster])

                    if roster.homeAway == "away":
                        match.away_competitor.roster = json.dumps([item.to_dict() for item in roster.roster])

            self.__match_repo.upsert(match)
            matches.append(match)

        
        self.schedule_match_if_necessary(matches, league)

    def schedule_match_if_necessary(self, matches: List[Match], league: str) -> None:
        for match in matches:
            match_date = datetime.strptime(match.date, '%Y-%m-%dT%H:%MZ').replace(tzinfo=timezone.utc, microsecond=0)
            if match_date > datetime.now(timezone.utc):
                self.__scheduler.schedule(match_date, partial(self.process_league_schedule, league))
                return

            # checa se a partida está em andamento
            if not match.completed:
                self.__scheduler.schedule(datetime.now(timezone.utc) + timedelta(minutes=5, microseconds=0), partial(self.process_league_schedule, league))
                return
        
        self.__scheduler.schedule(datetime.now(timezone.utc).replace(microsecond=0, hour=10, minute=0) + timedelta(days=1, microseconds=0), partial(self.process_league_schedule, league))