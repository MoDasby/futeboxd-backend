from typing import List, Dict, Any
from dataclasses import field, dataclass

class HeadShot:
    href: str
    alt: str

    def to_dict(self):
        return {
            "href": self.href,
            "alt": self.alt
        }

class Athlete:
    id: str
    fullName: str
    displayName: str
    headshot: HeadShot

    def to_dict(self):
        return {
            "id": self.id,
            "fullName": self.fullName,
            "displayName": self.displayName,
            "headshot": self.headshot.to_dict() if self.headshot else None
        }

class Position:
    displayName: str
    abbreviation: str

    def to_dict(self):
        return {
            "displayName": self.displayName,
            "abbreviation": self.abbreviation
        }

class RosterItem:
    starter: bool
    jersey: str
    subbedIn: bool
    subbedOut: bool
    athlete: Athlete
    position: Position

    def to_dict(self):
        return {
            "starter": self.starter,
            "jersey": self.jersey,
            "subbedIn": self.subbedIn,
            "subbedOut": self.subbedOut,
            "athlete": self.athlete.to_dict(),
            "position": self.position.to_dict() if self.position else None
        }

class Roster:
    homeAway: str
    roster: List[RosterItem]

class MatchSummary:
    rosters: List[Roster]

class Logo:
    href: str
    alt: str
    rel: List[str]
    width: int
    height: int

class Team:
    id: str
    name: str
    abbreviation: str
    color: str
    logos: List[Logo]

class Competitor:
    home_away: str
    winner: bool
    score: int
    team: Team

@dataclass
class Event:
    id: str
    venue: str
    match_date: str
    notes: List[Any]
    competitors: List[Competitor]
    completed: bool
    status_name: str
    events: List[Dict[str, Any]]
    competition_name: str

class LeagueSchedule:
    events: List[Event] = field(default_factory=list)