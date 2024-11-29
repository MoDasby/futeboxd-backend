from dataclasses import dataclass
from typing import Any, List
from team.models import Team

@dataclass
class EventType:
    id: str
    text: str

@dataclass
class Clock:
    value: float
    display_value: str

@dataclass
class Event:
    type: EventType
    text: str
    clock: Clock

class Competitor:
    home_away: str
    winner: bool
    score: int
    team: Team
    roster: str

class Match:
    id: int
    venue: str
    date: str
    note: str
    completed: bool
    status_name: str
    competition_name: str
    home_competitor: Competitor
    away_competitor: Competitor
    events: List[Any]