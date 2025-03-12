import { League } from "@/models/league";
import { Competitor, isValidStatus, Match, MatchEvent, StatusName } from "@/models/match";
import { Team } from "@/models/team";
import MatchDataUpdater from ".";
import logger from "@/util/logger";

const BASE_URL = "https://site.api.espn.com/apis/site/v2";

export async function getDaySchedule(league: League): Promise<Match[]> {
    const res = await fetch(`${BASE_URL}/sports/soccer/${league.espn_id}/scoreboard?lang=pt`);
    const data = await res.json();

    const dayMatches = data.events;
    const output: Match[] = []

    for (const espnMatch of dayMatches) {
        const matchEvents: MatchEvent[] = [];

        let homeCompetitor: Competitor | undefined = undefined
        let awayCompetitor: Competitor | undefined = undefined

        for (const espnCompetitor of espnMatch.competitions[0].competitors) {
            const team: Team = {
                id: espnCompetitor.team.id,
                name: espnCompetitor.team.name,
                color: espnCompetitor.team.color || "000000",
                abbreviation: espnCompetitor.team.abbreviation,
                logo: espnCompetitor.team.logo || ""
            };

            const parsedScore = Number.parseInt(espnCompetitor.score);
            const parsedValue = Number.parseInt(espnCompetitor.score?.value);

            const competitor: Competitor = {
                team,
                score: !Number.isNaN(parsedScore) ? parsedScore
                    : !Number.isNaN(parsedValue) ? parsedValue
                        : 0,
                winner: espnCompetitor.winner,
            };

            if (espnCompetitor.homeAway === "home") homeCompetitor = competitor
            if (espnCompetitor.homeAway === "away") awayCompetitor = competitor
        }

        if (!homeCompetitor || !awayCompetitor) throw new Error("Home competitor ou away competitor está nulo");

        const matchEvent: Match = {
            id: espnMatch.id,
            venue: espnMatch.competitions[0].venue?.fullName || null,
            date: espnMatch.date,
            note: espnMatch.competitions[0].notes[0]?.headline || null,
            homeCompetitor: homeCompetitor,
            awayCompetitor: awayCompetitor,
            completed: espnMatch.competitions[0].status.type.completed,
            statusName: espnMatch.competitions[0].status.type.name,
            events: matchEvents,
            league: league,
        };

        output.push(matchEvent);
    }

    return output;
}

async function getMatchSummary(matchId: number): Promise<MatchEvent[]> {
    const res = await fetch(`${BASE_URL}/sports/soccer/all/summary?lang=pt&event=${matchId}`);
    const data = await res.json();

    if (!data.keyEvents) return [];

    const matchEvents: MatchEvent[] = []

    try {
        for (const event of data.keyEvents) {
            if (/gol|goal/i.test(event.type.text as string)) {
                matchEvents.push({
                    clock_value: Number.parseInt(event.clock.value),
                    participant_name: event.participants[0].athlete.displayName,
                    team_id: Number.parseInt(event.team.id),
                    type: {
                        id: event.type.id,
                        text: event.type.text
                    }
                })
            }
        }
    } catch {
        return []
    }

    return matchEvents
}

async function getTeamSchedule(teamId: number, season: number, leagues: League[]): Promise<Match[]> {
    const res = await fetch(`${BASE_URL}/sports/soccer/all/teams/${teamId}/schedule?lang=pt&season=${season}`)
    const data = await res.json();

    const output: Match[] = []

    for (const espnMatch of data.events) {
        const league = leagues.find(l => l.espn_id === espnMatch.league.slug)

        if (league) {
            const matchEvents: MatchEvent[] = [];

            let homeCompetitor: Competitor | undefined = undefined
            let awayCompetitor: Competitor | undefined = undefined

            for (const espnCompetitor of espnMatch.competitions[0].competitors) {
                const team: Team = {
                    id: espnCompetitor.team.id,
                    name: espnCompetitor.team.displayName,
                    color: espnCompetitor.team.color || "000000",
                    abbreviation: espnCompetitor.team.abbreviation || (espnCompetitor.team.displayName as string).substring(0, 3),
                    logo: espnCompetitor.team.logo || ""
                };

                const parsedScore = Number.parseInt(espnCompetitor.score);
                const parsedValue = Number.parseInt(espnCompetitor.score?.value);

                const competitor: Competitor = {
                    team,
                    score: !Number.isNaN(parsedScore) ? parsedScore
                        : !Number.isNaN(parsedValue) ? parsedValue
                            : 0,
                    winner: espnCompetitor.winner,
                };

                if (espnCompetitor.homeAway === "home") homeCompetitor = competitor
                if (espnCompetitor.homeAway === "away") awayCompetitor = competitor
            }

            if (!homeCompetitor || !awayCompetitor) throw new Error("Home competitor ou away competitor está nulo");

            const matchEvent: Match = {
                id: espnMatch.id,
                venue: espnMatch.competitions[0].venue?.fullName || "",
                date: espnMatch.date,
                note: espnMatch.competitions[0].notes[0]?.headline || null,
                homeCompetitor: homeCompetitor,
                awayCompetitor: awayCompetitor,
                completed: espnMatch.competitions[0].status.type.completed,
                statusName: normalizeStatusName(espnMatch.id, espnMatch.competitions[0].status.type.name),
                events: matchEvents,
                league: league,
            };

            output.push(matchEvent);
        }
    }

    return output
}

function normalizeStatusName(matchId: string, statusName: string): StatusName {
    const endMatchStatusNames = ["STATUS_END_OF_EXTRATIME", "STATUS_FINAL_PEN", "STATUS_FINAL_AET"]
    
    if (endMatchStatusNames.includes(statusName)) {
        return StatusName.STATUS_FULL_TIME
    }

    if (isValidStatus(statusName)) {
        return statusName as StatusName
    }

    logger.warn(`Status desconhecido encontrado, substituindo por STATUS_UNKNOWN`, {
        match: matchId,
        originalStatus: statusName
    });

    return StatusName.STATUS_UNKNOWN
}

export const espnMatchUpdater: MatchDataUpdater = {
    getDaySchedule,
    getMatch: () => (new Promise(() => ({} as Match))),
    getMatchSummary,
    getTeamSchedule,
    normalizeStatusName
}