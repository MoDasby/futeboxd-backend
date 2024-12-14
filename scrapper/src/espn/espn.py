from requests import get
from .models import LeagueSchedule, Event, Competitor, Team, Logo, MatchSummary, Roster, RosterItem, Athlete, HeadShot, Position
from typing import List

class Espn():
    BASE_URL = "https://site.api.espn.com/apis/site/v2"

    def get_calendar(self, league):
        res = get(f"{self.BASE_URL}/{league}/scoreboard")

        return res.json()["leagues"][0]["calendar"]
    
    def get_day_schedule(self, league) -> LeagueSchedule:
        res = get(f"{self.BASE_URL}/sports/soccer/{league}/scoreboard?lang=pt").json()

        events = res["events"]
        competition_name = f"{res["leagues"][0]["name"]} {res["leagues"][0]["season"]["year"]}"

        output = LeagueSchedule()
        output.events = []

        for event in events:
            matchEvents = []

            if "details" in event["competitions"][0] and len(event["competitions"][0]["details"]) > 0:
                for detail in event["competitions"][0]["details"]:
                    if "gol" in detail["type"]["text"].lower() or "goal" in detail["type"]["text"].lower():
                        matchEvents.append(detail)
            
            matchCompetitors = []

            for espnCompetitor in event["competitions"][0]["competitors"]:
                competitor = Competitor()
                team = Team()

                team.name = espnCompetitor["team"]["name"]
                team.id = espnCompetitor["team"]["id"]
                team.color = espnCompetitor["team"].get("color", "000000")
                team.abbreviation = espnCompetitor["team"]["abbreviation"]

                logo = Logo()
                logo.href = espnCompetitor["team"]["logo"]
                team.logos = [logo]

                competitor.team = team

                competitor.home_away = espnCompetitor["homeAway"]
                competitor.score = espnCompetitor["score"]
                competitor.winner = espnCompetitor["winner"]

                matchCompetitors.append(competitor)

            event = {
                "id": event["id"],
                "venue": event["competitions"][0]["venue"]["fullName"],
                "match_date": event["date"],
                "notes": event["competitions"][0]["notes"],
                "competitors": matchCompetitors,
                "completed": event["competitions"][0]["status"]["type"]["completed"],
                "status_name": event["competitions"][0]["status"]["type"]["name"],
                "events": matchEvents,
                "competition_name": competition_name
            }

            output.events.append(Event(**event))
            
        return output

    def get_match_summary(self, match_id: int) -> MatchSummary:
        url = f"{self.BASE_URL}/sports/soccer/all/summary?lang=pt&event={match_id}"
        print(f"pegando summary da partida: {match_id}")

        data = get(url).json()

        rosters = data["rosters"]

        matchSummary = MatchSummary()
        matchSummary.rosters = []
        matchSummary.events = data.get("keyEvents", [])


        for espnRoster in rosters:
            roster = Roster()
            rosterItems: List[RosterItem] = []

            roster.homeAway = espnRoster["homeAway"]

            if not espnRoster.get("roster"):
                continue
            
            for espnRosterItem in espnRoster["roster"]:
                rosterItem = RosterItem()
                position = Position()
                athlete = Athlete()

                athlete.displayName = espnRosterItem["athlete"]["displayName"]
                athlete.fullName = espnRosterItem["athlete"]["fullName"]
                athlete.id = espnRosterItem["athlete"]["id"]

                if espnRosterItem["athlete"].get("headshot"):
                    headshot = HeadShot()

                    headshot.alt = espnRosterItem["athlete"]["headshot"]["alt"]
                    headshot.href = espnRosterItem["athlete"]["headshot"]["href"]

                    athlete.headshot = headshot
                else:
                    athlete.headshot = None

                espnPosition = espnRosterItem.get("position")

                if espnPosition:
                    position.abbreviation = espnPosition.get("abbreviation", "")
                    position.displayName = espnPosition.get("displayName", "")
                else:
                    position.abbreviation = ""
                    position.displayName = ""

                rosterItem.jersey = espnRosterItem["jersey"]
                rosterItem.starter = espnRosterItem["starter"]
                rosterItem.subbedIn = espnRosterItem["subbedIn"]
                rosterItem.subbedOut = espnRosterItem["subbedOut"]
                rosterItem.position = position
                rosterItem.athlete = athlete

                rosterItems.append(rosterItem)

            roster.roster = rosterItems
            matchSummary.rosters.append(roster)

        return matchSummary




    def get_match(self, match_id):
        url = f"{self.BASE_URL}/sports/soccer/all/summary?lang=pt&event={match_id}"
        event_data = get(url).json()
        
        match = match.new_empty_match()

        if len(event_data["header"]["competitions"][0]["notes"]) > 0:
            match["Note"] = event_data["header"]["competitions"][0]["notes"]

        for _, espn_competitor in enumerate(event_data['header']['competitions'][0]['competitors']):
            competitor = {
                'HomeAway': espn_competitor['homeAway'],
                'Winner': espn_competitor['winner'],
                'Score': int(espn_competitor.get('score', espn_competitor.get('order', 0))),
                'Team': {
                    'ID': int(espn_competitor['team']['id']),
                    'Name': espn_competitor['team']['name'],
                    'Abbreviation': espn_competitor['team']['abbreviation'],
                    'Color': espn_competitor['team']['color'],
                    'Logo': espn_competitor['team']['logos'][0]['href'] if espn_competitor['team']['logos'] else ''
                }
            }

            if competitor['HomeAway'] == 'home':
                match['HomeCompetitor'] = competitor
            else:
                match['AwayCompetitor'] = competitor

        return match