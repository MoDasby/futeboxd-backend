import loadEnvs from "@/config/load-envs";

loadEnvs()

import db from "@/db";
import { League, listLeagues } from "@/models/league";
import { upsertMatch } from "@/models/match";
import { getMandatoryTeams, insertTeamIfNotExists } from "@/models/team";
import { espnMatchUpdater } from "@/service/match-updater/espn-match-updater";
import { Command } from "commander";

const program = new Command();

program
    .name("get team schedule")
    .option("--from <year>", "Ano de inicio para a busca, ex: 2020")
    .option("--to <year>", "Ano de fim para a busca, ex: 2023")
    .option("--source <source>", "Fonte para buscar as partidas, por enquanto apenas espn")
    .argument("<id>", "id do time no banco de dados")
    .action(async (id: string = "all") => {
        const options = program.opts();

        const from = Number(options.from)
        const to = Number(options.to) || new Date().getFullYear()

        await db.waitForDbReady()

        const leagues = await listLeagues()

        if (id === "all") {
            const teams = await getMandatoryTeams()

            for (const team of teams) {
                for (let season = from; season <= to; season++) {
                    await saveSeasonSchedule(Number.parseInt(team.id), leagues, season)
                }
            }

            return
        }

        for (let season = from; season <= to; season++) {
            await saveSeasonSchedule(Number.parseInt(id), leagues, season)
        }
    })
    .parse(process.argv);

async function saveSeasonSchedule(teamId: number, leagues: League[], season: number) {
    const matches = await espnMatchUpdater.getTeamSchedule(teamId, season, leagues)

    matches.forEach(async match => {
        const summary = await espnMatchUpdater.getMatchSummary(match.id)

        match.events = summary;

        await insertTeamIfNotExists(match.homeCompetitor.team)
        await insertTeamIfNotExists(match.awayCompetitor.team)
        await upsertMatch(match);

        console.log(`inserido partida: liga: ${match.league.name} ${match.id}: ${match.homeCompetitor.team.name} ${match.homeCompetitor.score} x ${match.awayCompetitor.score} ${match.awayCompetitor.team.name}`)
    })
}