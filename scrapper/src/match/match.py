from dataclasses import dataclass
from typing import Any, List
from team.models import Team

@dataclass
class EventType:
    id: str
    text: str

@dataclass
class Event:
    type: EventType
    text: str
    clock_value: int
    team_id: int
    participant_name: str

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
    events: List[Event]