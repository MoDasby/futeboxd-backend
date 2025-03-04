import dotenv from "dotenv"
import fs from "node:fs"

export default function loadEnvs() {
    if (process.env.NODE_ENV == 'production') {
        const { DB_PASSWORD_FILE } = process.env

        const file = fs.readFileSync(DB_PASSWORD_FILE!, { encoding: "utf-8" })

        process.env.POSTGRES_PASSWORD = file.trim()

        return
    }

    const output = dotenv.config()

    if (output.error) throw new Error(output.error.message);
}