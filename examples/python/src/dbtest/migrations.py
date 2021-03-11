import os
from typing import List

import sqlalchemy
import sqlalchemy.ext.asyncio


def apply_migrations(conn: sqlalchemy.engine.Connection, paths: List[str]):
    files = _find_sql_files(paths)

    for file in files:
        with open(file, "r") as fd:
            blob = fd.read()
        conn.execute(blob)


async def apply_migrations_async(conn: sqlalchemy.ext.asyncio.AsyncConnection, paths: List[str]):
    files = _find_sql_files(paths)

    for file in files:
        with open(file, "r") as fd:
            blob = fd.read()
        await conn.execute(sqlalchemy.text(blob))


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
