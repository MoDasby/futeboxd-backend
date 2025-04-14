import loadEnvs from './config/load-envs';

loadEnvs()

import database from "@/db"
import { listLeagues } from "./models/league"
import { processDaySchedule } from "./models/match"
import { deleteOldNews, processNews } from "./models/news"
import { geScrapper } from "./service/news-scrapper/ge-scrapper"
import { scheduleJob } from "node-schedule"
import { espnMatchUpdater } from "./service/match-updater/espn-match-updater"
import logger from './util/logger';

async function main() {

    await database.waitForDbReady()

    scheduleJob("0 13,18,23 * * *", async () => {
        const leagues = await listLeagues()

        leagues.forEach(async league => {
            await processDaySchedule(league, espnMatchUpdater)
        })
    })

    scheduleJob("0 12,15,19,2 * * *", async () => {
        await deleteOldNews()
        await processNews(geScrapper)
    })

    logger.info("Aplicação inicializada e execuções agendadas")
}

main()
