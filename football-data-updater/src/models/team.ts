import { client } from "../db";

export type Team = {
    id: string;
    name: string;
    abbreviation: string;
    color: string;
    logo: string;
}

export async function insertTeamIfNotExists(team: Team) { // TODO precisa baixar a logo e colocar no s3
    const query = `
        INSERT INTO teams(id, name, logo, abbreviation, color)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT(id) DO NOTHING
    `

    try {
        await client.query(query, [team.id, team.name, team.logo, team.abbreviation, team.color]);
    } catch (err) {
        console.log(`erro ao inserir time: ${err}`);
    }
}

export async function isAnyMandatory(teams: Team[]): Promise<boolean> {
    const query = `
        SELECT COALESCE(BOOL_OR(mandatory), FALSE) AS has_mandatory
        FROM teams
        WHERE id IN (${teams.map((_, i) => `$${i + 1}`).join(", ")});
    `

    try {
        const result = await client.query<{ has_mandatory: boolean }>(query, teams.map(t => t.id));

        return result.rows[0].has_mandatory
    } catch (err) {
        console.log(`erro ao inserir time: ${err}`);
    }

    return false;
}