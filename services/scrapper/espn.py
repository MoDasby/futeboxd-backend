from requests import get

class Espn():
    BASE_URL = "https://site.api.espn.com/apis/site/v2/sports/soccer"

    def get_calendar(self, league):
        res = get(f"{self.BASE_URL}/{league}/scoreboard")

        return res.json()["leagues"][0]["calendar"]
    
    def get_match(self, match_id):
        res = get(f"{self.BASE_URL}/all/summary?lang=pt&event={match_id}").json()

        output = {
            "id": res["header"]["id"],
            "venue": {
                "name": res["gameInfo"]["venue"]["fullName"],
                "city": res["gameInfo"]["venue"]["address"]["city"]
            },
            "date": res["header"]["competitions"][0]["date"],
            "notes": res["header"]["competitions"][0]["notes"],
            "competitors": res["header"]["competitions"][0]["competitors"],
            "completed": res["header"]["competitions"][0]["status"]["type"]["completed"],
            "status_name": res["header"]["competitions"][0]["status"]["type"]["name"],
            "competition_name": res["header"]["league"]["name"],
            "events": res["keyEvents"]
        }

        return output