import db from "@/db";
import logger from "@/util/logger";

export type Team = {
	id: string;
	name: string;
	abbreviation: string;
	color: string;
	logo: string;
};

export async function insertTeamIfNotExists(team: Team) {
	// TODO precisa baixar a logo e colocar no s3
	const query = `
        INSERT INTO teams(id, name, logo, abbreviation, color)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT(id) DO NOTHING
    `;

	try {
		await db.query(query, [
			team.id,
			team.name,
			team.logo,
			team.abbreviation,
			team.color,
		]);
	} catch (error) {
		const err = error as Error;

		logger.error("Erro ao inserir time", {
			originalError: err.message,
			stackTrace: err.stack,
		});
	}
}

export async function isAnyMandatory(teams: Team[]): Promise<boolean> {
	const query = `
        SELECT COALESCE(BOOL_OR(mandatory), FALSE) AS has_mandatory
        FROM teams
        WHERE id IN (${teams.map((_, i) => `$${i + 1}`).join(", ")});
    `;

	try {
		const result = await db.query<{ has_mandatory: boolean }>(
			query,
			teams.map((t) => t.id)
		);

		return result.rows[0].has_mandatory;
	} catch (error) {
		const err = error as Error;

		logger.error("Erro ao verificar se um time é mandatório", {
			originalError: err.message,
			stackTrace: err.stack,
		});
	}

	return false;
}

export async function listTeams({
	mandatory,
}: {
	mandatory?: boolean;
}): Promise<Team[]> {
	const query = `
        SELECT id, name, logo, abbreviation, color
        FROM teams
        ${mandatory ? "WHERE mandatory" : ""}
    `;

	try {
		const result = await db.query<Team>(query);
		return result.rows;
	} catch (error) {
		const err = error as Error;

		logger.error(`Erro ao buscar times mandatórios`, {
			originalError: err.message,
			stackTrace: err.stack,
		});

		return [];
	}
}
