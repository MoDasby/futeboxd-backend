import { client } from "@/db"


export type League = {
    id: number
    name: string
    logo: string
    espn_id: string
}

export async function listLeagues(): Promise<League[]> {
    const result = await client.query<League>(`
        SELECT id, name, logo, espn_id FROM leagues;    
    `);

    return result.rows
}