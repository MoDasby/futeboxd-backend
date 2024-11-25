from models.match import Match, Clock, Competitor, Event, EventType

def from_schedule_to_match(schedule_item) -> Match:
    match = Match()

    match.id = schedule_item["id"]
    match.venue = schedule_item["venue"]
    match.date = schedule_item["match_date"]
    match.note = ""

    if len(schedule_item["notes"]) > 0:
        match.note = schedule_item["notes"][0]["headline"]

    for c in schedule_item["competitors"]:
        competitor = Competitor()

        competitor.home_away = c["homeAway"]
        competitor.score = c["score"]
        competitor.team_id = c["team"]["id"]
        competitor.winner = c["winner"]

        if c["homeAway"] == "home":
            match.home_competitor = competitor

        if c["homeAway"] == "away":
            match.away_competitor = competitor
    
    match.completed = schedule_item["status_name"] == "STATUS_FULL_TIME"
    match.events = schedule_item["events"]
    match.status_name = schedule_item["status_name"]
    match.competition_name = ""

    return match