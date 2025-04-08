import { News } from "@/models/news";
import cheerio from "cheerio";
import NewsScrapper from ".";

async function loadHTML(): Promise<cheerio.Root> {
    const headers = {
        "User-Agent":
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
    };

    const res = await fetch("https://ge.globo.com/futebol", {
        headers
    })

    const rawResponse = await res.text()

    return cheerio.load(rawResponse)
}

async function getNews(): Promise<News[]> {
    const $ = await loadHTML()

    const newsList: News[] = [];

    $(".bastian-feed-item").each((_, element) => {
        const titleElement = $(element).find(".feed-post-link");
        const descriptionElement = $(element).find(".feed-post-body-resumo");
        const linkElement = $(element).find(".feed-post-link");
        const imageElement = $(element).find(".bstn-fd-picture-image");

        if (titleElement.length && linkElement.length && descriptionElement.length && imageElement.length) {
            const title = titleElement.text().trim();
            const link = linkElement.attr("href");
            const description = descriptionElement.text().trim();
            const imageLink = imageElement.attr("src");

            if (!link || !imageLink) return

            newsList.push({
                title,
                description,
                link,
                imageLink
            });
        }
    });

    return newsList;
}

export const geScrapper: NewsScrapper = {
    getNews
}