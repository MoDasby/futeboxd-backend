import { client } from "@/db"
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

async function insertIfNotExists(news: News) {
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