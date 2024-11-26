from multiprocessing import Process
from datetime import timezone, datetime
from time import sleep
from typing import Any, Dict
import logger

class Scheduler():
    scheduled: Dict[datetime, Any]

    def __init__(self) -> None:
        self.scheduled = {}

    def schedule(self, executionDate: datetime, func: Any) -> None:
        logger.info(f"agendando execucao de {func.func.__name__}, args: {", ".join(repr(arg) for arg in func.args)} para {executionDate}")
        self.scheduled[executionDate] = func
    
    def _start(self) -> None:
        while True:
            currentDate = datetime.now(timezone.utc)

            for executionDate, func in list(self.scheduled.items()):
                if currentDate >= executionDate:
                    logger.info(f"executando {func}")
                    proc = Process(target=func)
                    proc.start()
                    self.scheduled.pop(executionDate, None)

            sleep(1)

    def run(self) -> None:
        p1 = Process(target=self._start)
        p1.start()
