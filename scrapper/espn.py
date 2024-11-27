from requests import get

class Espn():
    BASE_URL = "https://site.api.espn.com/apis/site/v2"

    def get_calendar(self, league):
        res = get(f"{self.BASE_URL}/{league}/scoreboard")

        return res.json()["leagues"][0]["calendar"]
    
    def get_day_schedule(self, league):
        res = get(f"{self.BASE_URL}/sports/soccer/{league}/scoreboard?lang=pt").json()

        events = res["events"]

        output = []

        for event in events:
            matchEvents = []

            if len(event["competitions"][0]["details"]) > 0:
                matchEvents = event["competitions"][0]["details"]

            output.append({
            "id": event["id"],
            "venue": event["competitions"][0]["venue"]["fullName"],
            "match_date": event["date"],
            "notes": event["competitions"][0]["notes"],
            "competitors": event["competitions"][0]["competitors"],
            "completed": event["competitions"][0]["status"]["type"]["completed"],
            "status_name": event["competitions"][0]["status"]["type"]["name"],
            "home_away": event["competitions"][0],
            "events": matchEvents
        })
            
        return output


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