from enum import Enum
from typing import List

class League(Enum):
    # Brasileiros
    BRASILEIRAO_SERIE_A = ("bra.1", "Brasileirão Série A")
    BRASILEIRAO_SERIE_B = ("bra.2", "Brasileirão Série B")
    BRASILEIRAO_SERIE_C = ("bra.3", "Brasileirão Série C")
    COPA_DO_BRASIL = ("bra.copa_do_brazil", "Copa do Brasil")
    SUPERCOPA_DO_BRASIL = ("bra.supercopa_do_brazil", "Supercopa do Brasil")
    
    # Estaduais
    PAULISTAO = ("bra.camp.paulista", "Campeonato Paulista")
    CARIOCA = ("bra.camp.carioca", "Campeonato Carioca")
    GAUCHO = ("bra.camp.gaucho", "Campeonato Gaúcho")
    MINEIRO = ("bra.camp.mineiro", "Campeonato Mineiro")
    COPA_NORDESTE = ("bra.copa_do_nordeste", "Copa do Nordeste")
    
    # UEFA
    CHAMPIONS_LEAGUE = ("uefa.champions", "UEFA Champions League")
    EUROPA_LEAGUE = ("uefa.europa", "UEFA Europa League")
    EURO_QUALIFIERS = ("uefa.euroq", "Euro Qualifiers")
    
    # FIFA/CONMEBOL
    WORLD_CUP = ("fifa.world", "FIFA World Cup")
    WORLD_CUP_CLUBS = ("fifa.intercontinental_cup", "Mundial de clubes")
    LIBERTADORES = ("conmebol.libertadores", "Copa Libertadores")
    SUDAMERICANA = ("conmebol.sudamericana", "Copa Sul-Americana")
    COPA_AMERICA = ("conmebol.america", "Copa América")
    
    # Ligas Europeias
    PREMIER_LEAGUE = ("eng.1", "Premier League")
    CHAMPIONSHIP = ("eng.2", "EFL Championship")
    LA_LIGA = ("esp.1", "La Liga")
    LIGUE_1 = ("fra.1", "Ligue 1")
    BUNDESLIGA = ("ger.1", "Bundesliga")
    SERIE_A = ("ita.1", "Serie A")
    PRIMEIRA_LIGA = ("por.1", "Primeira Liga")

    def __init__(self, id: str, display_name: str):
        self.id = id
        self.display_name = display_name

    @classmethod
    def list_all(self) -> List[str]:
        return [comp.id for comp in self]
