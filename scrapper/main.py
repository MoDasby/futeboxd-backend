from espn import Espn

espn = Espn()

calendar = espn.get_calendar("bra.1")

match = espn.get_match("348167")

print(calendar)