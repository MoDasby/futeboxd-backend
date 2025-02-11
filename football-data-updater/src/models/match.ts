
import { League } from "./league";
import { insertTeamIfNotExists, isAnyMandatory, Team } from "./team";
import { schedule } from "./scheduler";
import { client } from "@/db";
import MatchDataUpdater from "@/service/match-updater";


const statusNameValues = ["STATUS_FINAL_PEN", "STATUS_SCHEDULED", "STATUS_FULL_TIME", "STATUS_SECOND_HALF", "STATUS_HALFTIME", "STATUS_FIRST_HALF", "STATUS_UNKNOWN", "STATUS_POSTPONED"]

export type Competitor = {
    winner: boolean;
    score: number;
    team: Team;
}

export type Match = {
    id: number;
    venue: string;
    date: string;
    note: string;
    completed: boolean;
    statusName: string;
    league: League;
    homeCompetitor: Competitor;
    awayCompetitor: Competitor;
    events: MatchEvent[]
}

export type MatchEvent = {
    type: {
        id: string
        text: string
    },
    clock_value: number
    team_id: number
    participant_name: string
}

export async function processDaySchedule(league: League, gateway: MatchDataUpdater): Promise<void> {
    const daySchedule = await gateway.getDaySchedule(league);

    if (daySchedule.length === 0) {
        scheduleProcessNextDay(league, gateway);

        return
    }

    daySchedule.forEach(async m => {
        if (!(await isAnyMandatory([m.homeCompetitor.team, m.awayCompetitor.team]))) {
            console.info(`Nenhum time mandátorio em ${getFormattedMatchName(m)}, pulando`)
            return;
        }

        if (!statusNameValues.includes(m.statusName)) {
            m.statusName = "STATUS_UNKNOWN"
        }

        const matchSummary = await gateway.getMatchSummary(m.id);

        m.events = matchSummary

        console.info(`Processando partida ${getFormattedMatchName(m)}`);

        await insertTeamIfNotExists(m.homeCompetitor.team)
        await insertTeamIfNotExists(m.awayCompetitor.team)
        await upsertMatch(m)

        if (matchStarted(m)) {
            const fiveMinutesLater = new Date();
            fiveMinutesLater.setSeconds(0)
            fiveMinutesLater.setMilliseconds(0)
            fiveMinutesLater.setMinutes(fiveMinutesLater.getMinutes() + 5);
            scheduleNextProcessing(fiveMinutesLater, league, gateway);
            return;
        }

        const matchDate = new Date(m.date);
        matchDate.setMilliseconds(0);
        matchDate.setSeconds(0);

        if (matchDate > new Date()) {
            matchDate.setMinutes(matchDate.getMinutes() + 10)
            scheduleNextProcessing(matchDate, league, gateway);

            return
        }

        if (!matchStarted(m)) {
            const fiveMinutesLater = new Date();
            fiveMinutesLater.setMinutes(fiveMinutesLater.getMinutes() + 10);
            scheduleNextProcessing(fiveMinutesLater, league, gateway);
            return;
        }
    })
}

export function scheduleProcessNextDay(league: League, gateway: MatchDataUpdater) {
    const nextDayAt10AM = new Date();
    nextDayAt10AM.setHours(10, 0, 0, 0);
    nextDayAt10AM.setDate(nextDayAt10AM.getDate() + 1);
    schedule(nextDayAt10AM, league, () => processDaySchedule(league, gateway));
}

function scheduleNextProcessing(time: Date, league: League, gateway: MatchDataUpdater) {
    schedule(time, league, () => processDaySchedule(league, gateway));
}

export async function upsertMatch(match: Match) {


    const query = `
        INSERT INTO matches (
		id, match_date, venue, home_team_score, home_team_id,
		away_team_score, away_team_id,
		note, completed, status_name, league_id, events
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        ON CONFLICT (id) DO UPDATE SET
            match_date = EXCLUDED.match_date,
            venue = EXCLUDED.venue,
            home_team_score = EXCLUDED.home_team_score,
            home_team_id = EXCLUDED.home_team_id,
            away_team_score = EXCLUDED.away_team_score,
            away_team_id = EXCLUDED.away_team_id,
            note = EXCLUDED.note,
            completed = EXCLUDED.completed,
            status_name = EXCLUDED.status_name,
            league_id = EXCLUDED.league_id,
            events = EXCLUDED.events;
    `

    try {
        await client.query(query, [
            match.id, match.date, match.venue, match.homeCompetitor.score, match.homeCompetitor.team.id,
            match.awayCompetitor.score, match.awayCompetitor.team.id, match.note, match.completed, match.statusName,
            match.league.id, JSON.stringify(match.events)
        ]);
    } catch (err) {
        console.error(`erro ao inserir partida: ${err}`);
    }
}

function getFormattedMatchName(match: Match): string {
    return `liga: ${match.league.name} ${match.id}: ${match.homeCompetitor.team.name} ${match.homeCompetitor.score} x ${match.awayCompetitor.score} ${match.awayCompetitor.team.name}`;
}

function matchStarted(match: Match): boolean {
    const activeStatuses = new Set(["STATUS_FIRST_HALF", "STATUS_HALFTIME", "STATUS_SECOND_HALF"]);
    return activeStatuses.has(match.statusName)
}
