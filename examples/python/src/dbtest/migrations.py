import os
from typing import List

import asyncpg
import psycopg2.extensions


def apply_migrations(db: psycopg2.extensions.connection, paths: List[str]):
    files = _find_sql_files(paths)

    for file in files:
        with open(file, "r") as fd:
            blob = fd.read()
        cur = db.cursor()
        cur.execute(blob)
        cur.close()
        db.commit()


async def apply_migrations_async(db: asyncpg.Connection, paths: List[str]):
    files = _find_sql_files(paths)

    for file in files:
        with open(file, "r") as fd:
            blob = fd.read()
        await db.execute(blob)


def _find_sql_files(paths: List[str]) -> List[str]:
    files = []
    for path in paths:
        if not os.path.exists(path):
            raise FileNotFoundError(f"{path} does not exist")
        if os.path.isdir(path):
            for file in os.listdir(path):
                if file.endswith(".sql"):
                    files.append(os.path.join(path, file))
        else:
            files.append(path)
    files.sort()
    return files
