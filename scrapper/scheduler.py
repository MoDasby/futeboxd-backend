from multiprocessing import Process
from datetime import timezone, datetime
from time import sleep
from functools import partial
from typing import Callable, Any, Dict

class Scheduler():
    scheduled: Dict[datetime, Callable[[Any], Any]]

    def __init__(self) -> None:
        self.scheduled = {}

    def schedule(self, executionDate: datetime, func: Callable[[Any], Any], *args, **kwargs) -> None:
        print(f"agendando execucao de {func} para {executionDate}, {args}")
        self.scheduled[executionDate] = partial(func, *args, **kwargs)
    
    def _start(self) -> None:
        while True:
            currentDate = datetime.now(timezone.utc)

            for executionDate, func in list(self.scheduled.items()):
                if currentDate >= executionDate:
                    print(f"executando {func}")
                    proc = Process(target=func)
                    proc.start()
                    self.scheduled.pop(executionDate, None)

            sleep(1)

    def run(self) -> None:
        p1 = Process(target=self._start)
        p1.start()
