import dotenv from "dotenv"

const output = dotenv.config()

if (output.error) throw new Error(output.error.message);

import { Client } from 'pg';

let isConnected = false;

const client = new Client({
    host: process.env.DB_HOST,
    user: process.env.DB_USER,
    password: process.env.POSTGRES_PASSWORD,
    database: process.env.POSTGRES_DB,
    port: Number.parseInt(process.env.PGPORT || "5432"),
});

function getClient(): Client {
    if (!isConnected) {
        throw new Error("Banco de dados não está conectado")
    }

    return client
}

async function waitForDbReady(maxRetries: number = 5, delay: number = 2000) {
    if (isConnected) {
        console.log('Banco de dados pronto!');

        return
    }

    let retries = 0;
    while (retries < maxRetries) {
        try {
            await client.connect(); // Conecta ao banco de dados
            await client.query('SELECT 1'); // Executa uma query simples para verificar a conexão
            console.log('Banco de dados pronto!');
            isConnected = true;
            return;
        } catch (err) {
            console.log(`Banco de dados não está pronto, tentando de novo... (${retries + 1}/${maxRetries})`);
            retries++;
            await new Promise(resolve => setTimeout(resolve, delay));
        }
    }
    throw new Error('Banco de dados não está pronto');
};

export {
    getClient, waitForDbReady
}