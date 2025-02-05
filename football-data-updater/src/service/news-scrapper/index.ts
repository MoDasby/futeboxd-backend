import { News } from "@/models/news";

export default interface NewsScrapper {
    getNews(): Promise<News[]>
}