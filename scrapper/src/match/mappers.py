from .match import Match, Competitor, Event
from espn import Event
from team.models import Team

def from_schedule_to_match(event: Event) -> Match:
    match = Match()

    match.id = event.id
    match.venue = event.venue
    match.date = event.match_date
    match.note = ""

    if len(event.notes) > 0:
        match.note = event.notes[0]["headline"]

    for espnCompetitor in event.competitors:
        competitor = Competitor()
        team = Team()

        team.id = espnCompetitor.team.id
        team.abbreviation = espnCompetitor.team.abbreviation
        team.color = espnCompetitor.team.color
        team.name = espnCompetitor.team.name
        
        if len(espnCompetitor.team.logos) > 0:
            team.logo = espnCompetitor.team.logos[0].href

        competitor.home_away = espnCompetitor.home_away
        competitor.score = espnCompetitor.score
        competitor.team = team
        competitor.winner = espnCompetitor.winner

        if espnCompetitor.home_away == "home":
            match.home_competitor = competitor

        if espnCompetitor.home_away == "away":
            match.away_competitor = competitor
    
    match.completed = event.completed
    match.events = event.events
    match.status_name = event.status_name
    match.competition_name = event.competition_name

    return match