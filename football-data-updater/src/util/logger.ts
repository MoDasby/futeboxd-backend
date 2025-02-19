import winston from "winston";
import { WinstonTransport as AxiomTransport } from '@axiomhq/winston';

const logger = winston.createLogger({
    level: 'info',
    format: winston.format.json(),
    defaultMeta: { service: 'football-data-updater' },
    transports: [
        /* new AxiomTransport({
            dataset: "futeboxd-logs",
            token: process.env.AXIOM_TOKEN as string
        }), */
    ],
});

if (process.env.NODE_ENV != 'production') {
    logger.add(
      new winston.transports.Console({
        format: winston.format.simple(),
      }),
    );
  }

export default logger;
