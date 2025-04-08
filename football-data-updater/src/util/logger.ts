import winston from "winston";
import { WinstonTransport as AxiomTransport } from '@axiomhq/winston';

let transports: winston.transport[] = [];

if (process.env.NODE_ENV === "production") {
  transports.push(new AxiomTransport({
    dataset: "futeboxd-logs",
    token: process.env.AXIOM_TOKEN as string
  }))
} else {
  transports.push(new winston.transports.Console({
    format: winston.format.simple(),
  }))
}

const logger = winston.createLogger({
  format: winston.format.json(),
  defaultMeta: { service: 'football-data-updater' },
  transports
});

export default logger;
