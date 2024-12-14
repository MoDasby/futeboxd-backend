from psycopg2.pool import ThreadedConnectionPool
import os
from typing import Any

DB_HOST = os.getenv("DB_HOST", "localhost")
DB_NAME = os.getenv("DB_NAME", "futeboxd-football")
DB_USER = os.getenv("DB_USER", "postgres")
DB_PASSWORD = os.getenv("DB_PASSWORD", "123456")
DB_PORT = os.getenv("DB_PORT", "5433")

class DB:
    _instance = None

    def __new__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super().__new__(cls)
        return cls._instance

    def __init__(self, minconn=2, maxconn=20, **db_params):
        pass

    def execute(self, query: str, values: any):
        pass

    def fetch(self, query: str, values: any) -> Any:
        pass

class PostgresDB(DB):
    def __init__(self, minconn=2, maxconn=20):
        if not hasattr(self, "_pool"):
            DSN = f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
            self._pool = ThreadedConnectionPool(minconn, maxconn, DSN)

    def execute(self, query: str, values: any):
        conn = self._pool.getconn()

        try:
            with conn.cursor() as cursor:
                cursor.execute(query, values)
                conn.commit()
        finally:
            self._pool.putconn(conn)
    
    def fetch(self, query: str, values: any) -> Any:
        conn = self._pool.getconn()

        try:
            with conn.cursor() as cursor:
                cursor.execute(query, values)
                result = cursor.fetchone()

                return result[0]
        finally:
            self._pool.putconn(conn)