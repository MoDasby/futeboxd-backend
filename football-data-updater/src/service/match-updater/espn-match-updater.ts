import { League } from "@/models/league";
import { Competitor, Match, MatchEvent } from "@/models/match";
import { Team } from "@/models/team";
import MatchDataUpdater from ".";

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

            const competitor: Competitor = {
                team,
                score: Number.parseInt(espnCompetitor.score),
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

    return matchEvents
}

export const espnMatchUpdater: MatchDataUpdater = {
    getDaySchedule,
    getMatch: () => (new Promise(() => ({} as Match))),
    getMatchSummary
}