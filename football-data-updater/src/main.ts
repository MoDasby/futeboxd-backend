import loadEnvs from './config/load-envs';

loadEnvs()

import database from "@/db"
import { listLeagues } from "./models/league"
import { processDaySchedule } from "./models/match"
import { deleteOldNews, processNews } from "./models/news"
import { geScrapper } from "./service/news-scrapper/ge-scrapper"
import { scheduleJob } from "node-schedule"
import { espnMatchUpdater } from "./service/match-updater/espn-match-updater"

async function main() {

    await database.waitForDbReady()

    const leagues = await listLeagues()

    leagues.forEach(async league => {
        await processDaySchedule(league, espnMatchUpdater)
    })

    await processNews(geScrapper);

    scheduleJob("0 10,15 * * *", async () => {
        await deleteOldNews()
        await processNews(geScrapper)
    })
}

main()
