from datetime import timezone, datetime
from time import sleep
from typing import Any, Dict, List
from util import logger


class Scheduler():
    scheduled: Dict[datetime, List[Any]]

    def __init__(self) -> None:
        self.scheduled = {}

    def schedule(self, executionDate: datetime, func: Any) -> None:
        executionDate = executionDate.replace(microsecond=0)

        # Verifica se `executionDate` está em `self.scheduled`
        tasks = self.scheduled.get(executionDate)

        if not tasks:
            tasks = []

        # Verifica se `func` tem os atributos esperados
        if not hasattr(func, "args") or not hasattr(func, "func"):
            raise TypeError("O objeto `func` deve ter os atributos `args` e `func`.")

        # Verifica se a função já foi agendada
        for f in tasks:
            if isinstance(f, type(func)) and f.args[0] == func.args[0] and f.func.__name__ == func.func.__name__:
                return
        
        # Agenda a execução e faz o log
        logger.info(f"agendando execucao de {func.func.__name__}, args: {', '.join(repr(arg) for arg in func.args)} para {executionDate}")
        tasks.append(func)

        self.scheduled[executionDate] = tasks

    
    def _start(self) -> None:
        while True:
            currentDate = datetime.now(timezone.utc).replace(microsecond=0)
            for executionDate, tasks in list(self.scheduled.items()):
                if currentDate >= executionDate:
                    for func in tasks:
                        logger.info(f"executando {func.func.__name__}, args: {', '.join(repr(arg) for arg in func.args)}")
                        func()
                    
                    self.scheduled.pop(executionDate, [])

            sleep(1)

    def run(self) -> None:
        self._start()
