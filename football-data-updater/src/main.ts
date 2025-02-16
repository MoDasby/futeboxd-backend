require('module-alias/register');
import { waitForDbReady } from "@/db"
import { listLeagues } from "./models/league"
import { processDaySchedule } from "./models/match"
import { processNews } from "./models/news"
import { geScrapper } from "./service/news-scrapper/ge-scrapper"
import { scheduleJob } from "node-schedule"
import { espnMatchUpdater } from "./service/match-updater/espn-match-updater"

async function main() {
    await waitForDbReady()

    const leagues = await listLeagues()

    leagues.forEach(async league => {
        await processDaySchedule(league, espnMatchUpdater)
    })

    processNews(geScrapper);
    scheduleJob("0 10,15 * * *", () => {
        processNews(geScrapper)
    })
}

main()