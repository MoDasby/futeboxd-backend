import { Pool, PoolClient, PoolConfig, QueryResultRow } from 'pg';
import logger from "./util/logger";

let isConnected = false;

const config: PoolConfig = {
    host: process.env.DB_HOST,
    user: process.env.DB_USER,
    password: process.env.POSTGRES_PASSWORD,
    database: process.env.POSTGRES_DB,
    port: Number.parseInt(process.env.PGPORT || "5432"),
    max: 3,
    min: 1,
    idleTimeoutMillis: 20000,
    allowExitOnIdle: true
}

const pool = new Pool(config)

async function query<T extends QueryResultRow>(query: string, params?: any[]) {
    const client = await pool.connect()

    try {
        const result = client.query<T>(query, params)

        return result
    } finally {
        client.release()
    }
}

async function getClient(): Promise<PoolClient> {
    if (!isConnected) {
        throw new Error("Banco de dados não está conectado")
    }

    const client = await pool.connect()

    return client
}

async function waitForDbReady(maxRetries: number = 5, delay: number = 2000) {
    
    let retries = 0;
    while (retries < maxRetries) {
        try {
            const client = await pool.connect()
            logger.info('Banco de dados pronto!');
            isConnected = true;
            return;
        } catch (err) {
            logger.info(`Banco de dados não está pronto, tentando de novo... (${retries + 1}/${maxRetries})`);
            retries++;
            await new Promise(resolve => setTimeout(resolve, delay));
        }
    }
    throw new Error('Banco de dados não está pronto');
};

export default {
    getClient, query, waitForDbReady
}