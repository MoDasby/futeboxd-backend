import { getClient } from "@/db"
import NewsScrapper from "@/service/news-scrapper"

export type News = {
    title: string
    description: string
    link: string
    imageLink: string
}

export async function processNews(scrapper: NewsScrapper) {
    console.log("processando notícias")

    const news = await scrapper.getNews()

    await Promise.all(news.map(n => insertIfNotExists(n)))
}

export async function deleteOldNews(): Promise<void> {
    const client = getClient()

    const query = `
        DELETE FROM news WHERE created_at >= NOW() - INTERVAL '3 day'
    `

    console.info("Deletando notícias antigas")

    try {
        await client.query(query)
    } catch(err) {
        console.error(`Erro ao deletar noticias antigas: ${err}`)
    }
}

async function insertIfNotExists(news: News) {
    const client = getClient()
    
    const query = `
        INSERT INTO news (title, description, link, image_link)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT(link) DO NOTHING
    `

    try {
        await client.query(query, [news.title, news.description, news.link, news.imageLink]);
    } catch(err) {
        console.error(`ocorreu um erro ao inserir notícia: ${news}`)
    }
}