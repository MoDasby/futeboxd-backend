from match_repo import MatchRepository
from espn import Espn
from mappers import from_schedule_to_match

"""
cria ou atualiza a partida na tabela e retorna se deve tentar daqui cinco minutos
se todas as partidas foram finalizadas não há necessidade
"""
def upsert_league_schedule(repo: MatchRepository, espn: Espn, league: str) -> bool:
    day_schedule = espn.get_day_schedule(league)
    endedMatches = 0

    for event in day_schedule:
        match = from_schedule_to_match(event)

        repo.upsert(match)
    
    return endedMatches == len(day_schedule)