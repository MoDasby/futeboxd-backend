import db from "@/db"
import NewsScrapper from "@/service/news-scrapper"
import logger from "@/util/logger"

export type News = {
    title: string
    description: string
    link: string
    imageLink: string
}

export async function processNews(scrapper: NewsScrapper) {
    logger.info("processando notícias")

    try {
        const news = await scrapper.getNews()

        await Promise.all(news.map(n => insertIfNotExists(n)))
    } catch (err) {
        logger.info("Erro ao processar notícias", {
            err: (err as Error).message
        })
    }
}

export async function deleteOldNews(): Promise<void> {
    const query = `
        DELETE FROM news WHERE created_at >= NOW() - INTERVAL '3 day'
    `

    logger.info("Deletando notícias antigas")

    try {
        await db.query(query)
    } catch (err) {
        logger.error(`Erro ao deletar noticias antigas`, {
            err: (err as Error).message
        })
    }
}

async function insertIfNotExists(news: News) {
    const query = `
        INSERT INTO news (title, description, link, image_link)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT(link) DO NOTHING
    `

    try {
        await db.query(query, [news.title, news.description, news.link, news.imageLink]);
    } catch (err) {
        logger.error(`ocorreu um erro ao inserir notícia`, {
            err: (err as Error).message
        })
    }
}