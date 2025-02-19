import db from "@/db";

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
        await db.query(query, [team.id, team.name, team.logo, team.abbreviation, team.color]);
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
        const result = await db.query<{ has_mandatory: boolean }>(query, teams.map(t => t.id));

        return result.rows[0].has_mandatory
    } catch (err) {
        console.log(`erro ao inserir time: ${err}`);
    }

    return false;
}

export async function getMandatoryTeams(): Promise<Team[]> {
    const query = `
        SELECT id, name, logo, abbreviation, color
        FROM teams
        WHERE mandatory
    `;

    try {
        const result = await db.query<Team>(query);
        return result.rows; // Retorna os times mandatórios encontrados
    } catch (err) {
        console.log(`Erro ao buscar times mandatórios: ${err}`);
        return []; // Retorna um array vazio em caso de erro
    }
}
