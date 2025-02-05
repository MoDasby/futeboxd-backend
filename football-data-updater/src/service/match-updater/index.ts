import { League } from "../../models/league"
import { Match, MatchEvent } from "../../models/match"

export default interface MatchDataUpdater {
    getDaySchedule(league: League): Promise<Match[]>
    getMatch(matchId: number): Promise<Match>
    getMatchSummary(matchId: number): Promise<MatchEvent[]>
}