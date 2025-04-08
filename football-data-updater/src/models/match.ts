import { League } from "./league";
import { insertTeamIfNotExists, isAnyMandatory, Team } from "./team";
import { schedule } from "./scheduler";
import MatchDataUpdater from "@/service/match-updater";
import logger from "@/util/logger";
import db from "@/db";

export enum StatusName {
    STATUS_FINAL_PEN = "STATUS_FINAL_PEN",
    STATUS_SCHEDULED = "STATUS_SCHEDULED",
    STATUS_FULL_TIME = "STATUS_FULL_TIME",
    STATUS_SECOND_HALF = "STATUS_SECOND_HALF",
    STATUS_HALFTIME = "STATUS_HALFTIME",
    STATUS_FIRST_HALF = "STATUS_FIRST_HALF",
    STATUS_UNKNOWN = "STATUS_UNKNOWN",
    STATUS_POSTPONED = "STATUS_POSTPONED",
    STATUS_OVERTIME = "STATUS_OVERTIME",
    STATUS_SHOOTOUT = "STATUS_SHOOTOUT"
}

export function isValidStatus(status: string): boolean {
    return Object.values(StatusName).includes(status as StatusName);
}

const activeStatuses = new Set([
    StatusName.STATUS_FIRST_HALF,
    StatusName.STATUS_HALFTIME,
    StatusName.STATUS_SECOND_HALF
]);

function matchStarted(match: Match): boolean {
    return activeStatuses.has(match.statusName)
}


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
    statusName: StatusName;
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
    let processedAnyMatch = false;

    try {
        logger.info("Processando schedule", {
            league: league.espn_id
        })

        const start = performance.now()
        const daySchedule = await gateway.getDaySchedule(league);

        if (daySchedule.length === 0) {
            scheduleProcessNextDay(league, gateway);

            return
        }

        for (const m of daySchedule) {
            if (!(await isAnyMandatory([m.homeCompetitor.team, m.awayCompetitor.team]))) {
                logger.info(`Nenhum time mandátorio, pulando`, {
                    match: m.id,
                    homeTeam: m.homeCompetitor.team.name,
                    awayTeam: m.awayCompetitor.team.name
                });

                continue;
            }

            const matchSummary = await gateway.getMatchSummary(m.id);

            m.events = matchSummary

            logger.info(`Processando partida`, {
                match: m.id
            });

            await insertTeamIfNotExists(m.homeCompetitor.team)
            await insertTeamIfNotExists(m.awayCompetitor.team)
            await upsertMatch(m)

            if (matchStarted(m)) {
                const fiveMinutesLater = new Date();
                fiveMinutesLater.setSeconds(0)
                fiveMinutesLater.setMilliseconds(0)
                fiveMinutesLater.setMinutes(fiveMinutesLater.getMinutes() + 5);
                scheduleNextProcessing(fiveMinutesLater, league, gateway);
                processedAnyMatch = true;
                continue;
            }

            const matchDate = new Date(m.date);
            matchDate.setMilliseconds(0);
            matchDate.setSeconds(0);

            if (matchDate > new Date()) {
                matchDate.setMinutes(matchDate.getMinutes() + 10)
                scheduleNextProcessing(matchDate, league, gateway);
                processedAnyMatch = true;

                continue
            }

            if (!m.completed) {
                const date = new Date();
                date.setMinutes(date.getMinutes() + 5);
                date.setSeconds(0);
                date.setMilliseconds(0);
                scheduleNextProcessing(date, league, gateway);
                processedAnyMatch = true;

                continue;
            }
        }

        logger.info("Schedule Processada", {
            league: league.espn_id,
            duration: (performance.now() - start).toFixed(2)
        })
    } catch (err) {
        const error = err as Error

        logger.error("Erro ao processar schedule", {
            league: league.espn_id,
            originalError: error.message,
            stackTrace: error.stack
        })

        const date = new Date();
        date.setMinutes(date.getMinutes() + 10);
        date.setSeconds(0);
        date.setMilliseconds(0);
        scheduleNextProcessing(date, league, gateway);
    } finally {
        if (!processedAnyMatch) {
            scheduleProcessNextDay(league, gateway)
        }
    }
}

function scheduleProcessNextDay(league: League, gateway: MatchDataUpdater) {
    const nextDayAt10AM = new Date();
    nextDayAt10AM.setHours(10, 0, 0, 0);
    nextDayAt10AM.setDate(nextDayAt10AM.getDate() + 1);
    scheduleNextProcessing(nextDayAt10AM, league, gateway);
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
        await db.query(query, [
            match.id, match.date, match.venue, match.homeCompetitor.score, match.homeCompetitor.team.id,
            match.awayCompetitor.score, match.awayCompetitor.team.id, match.note, match.completed, match.statusName,
            match.league.id, JSON.stringify(match.events)
        ]);
        logger.info("Partida inserida ou atualizada", {
            match: match.id
        })
    } catch (error) {
        const err = error as Error
        logger.error(`Erro ao inserir ou atualizar partida`, {
            match: match.id,
            originalError: err.message,
            stackTrace: err.stack
        });
    }
}
