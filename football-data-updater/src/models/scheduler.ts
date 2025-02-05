import { Job, scheduleJob } from "node-schedule";
import { League } from "./league";

const scheduled: Map<string, Map<League, Job>> = new Map();

export function schedule(executionDate: Date, league: League, func: Function): void {
    executionDate.setMilliseconds(0);
    const executionDateStr = executionDate.toISOString();

    // Verificar se já existe um agendamento para essa data e liga
    if (scheduled.has(executionDateStr) && scheduled.get(executionDateStr)?.has(league)) {
      return;
    }

    const job = scheduleJob(executionDate, async () => {
      console.log(`Executando atualização de partidas para a liga: ${league.name}`);
      try {
        await func();
      } catch (err) {
        console.log(`Erro ao atualizar partidas para a liga ${league.name}: ${(err as Error).message}`);
      }
    });

    // Garantir que a liga seja adicionada no Map correto
    if (!scheduled.has(executionDateStr)) {
      scheduled.set(executionDateStr, new Map());
    }
    scheduled.get(executionDateStr)!.set(league, job);

    console.log(`Agendando atualização de partidas para a liga: ${league.name} para ${executionDateStr}`);
  }