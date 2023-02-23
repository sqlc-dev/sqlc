# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.17.2
# source: query.sql
import sqlalchemy
import sqlalchemy.ext.asyncio

from querytest import models


DELETE_BAR_BY_ID = """-- name: delete_bar_by_id \\:execresult
DELETE FROM bar WHERE id = :p1
"""


class Querier:
    def __init__(self, conn: sqlalchemy.engine.Connection):
        self._conn = conn

    def delete_bar_by_id(self, *, id: int) -> sqlalchemy.engine.Result:
        return self._conn.execute(sqlalchemy.text(DELETE_BAR_BY_ID), {"p1": id})


class AsyncQuerier:
    def __init__(self, conn: sqlalchemy.ext.asyncio.AsyncConnection):
        self._conn = conn

    async def delete_bar_by_id(self, *, id: int) -> sqlalchemy.engine.Result:
        return await self._conn.execute(sqlalchemy.text(DELETE_BAR_BY_ID), {"p1": id})
