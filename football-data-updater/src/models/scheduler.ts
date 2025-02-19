import { Job, scheduleJob } from "node-schedule";
import { League } from "./league";
import logger from "@/util/logger";

const scheduled: Map<string, Map<League, Job>> = new Map();

export function schedule(executionDate: Date, league: League, func: Function): void {
  executionDate.setMilliseconds(0);
  const executionDateStr = executionDate.toISOString();

  // Verificar se já existe um agendamento para essa data e liga
  if (scheduled.has(executionDateStr) && scheduled.get(executionDateStr)?.has(league)) {
    return;
  }

  const job = scheduleJob(executionDate, async () => {
    logger.info(`Executando atualização de partidas para a liga`, {
      espn_id: league.espn_id
    });
    try {
      await func();
    } catch (err) {
      logger.error(`Erro ao atualizar partidas para a liga`, {
        espn_id: league.espn_id,
        err
      });
    }
  });

  // Garantir que a liga seja adicionada no Map correto
  if (!scheduled.has(executionDateStr)) {
    scheduled.set(executionDateStr, new Map());
  }
  scheduled.get(executionDateStr)!.set(league, job);

  logger.info(`Agendando processamento`, {
    espn_id: league.espn_id,
    scheduledTo: executionDateStr
  })
}