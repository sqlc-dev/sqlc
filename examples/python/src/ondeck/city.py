# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.17.1
# source: city.sql
from typing import AsyncIterator, Optional

import sqlalchemy
import sqlalchemy.ext.asyncio

from ondeck import models


CREATE_CITY = """-- name: create_city \\:one
INSERT INTO city (
    name,
    slug
) VALUES (
    :p1,
    :p2
) RETURNING slug, name
"""


GET_CITY = """-- name: get_city \\:one
SELECT slug, name
FROM city
WHERE slug = :p1
"""


LIST_CITIES = """-- name: list_cities \\:many
SELECT slug, name
FROM city
ORDER BY name
"""


UPDATE_CITY_NAME = """-- name: update_city_name \\:exec
UPDATE city
SET name = :p2
WHERE slug = :p1
"""


class AsyncQuerier:
    def __init__(self, conn: sqlalchemy.ext.asyncio.AsyncConnection):
        self._conn = conn

    async def create_city(self, *, name: str, slug: str) -> Optional[models.City]:
        row = (await self._conn.execute(sqlalchemy.text(CREATE_CITY), {"p1": name, "p2": slug})).first()
        if row is None:
            return None
        return models.City(
            slug=row[0],
            name=row[1],
        )

    async def get_city(self, *, slug: str) -> Optional[models.City]:
        row = (await self._conn.execute(sqlalchemy.text(GET_CITY), {"p1": slug})).first()
        if row is None:
            return None
        return models.City(
            slug=row[0],
            name=row[1],
        )

    async def list_cities(self) -> AsyncIterator[models.City]:
        result = await self._conn.stream(sqlalchemy.text(LIST_CITIES))
        async for row in result:
            yield models.City(
                slug=row[0],
                name=row[1],
            )

    async def update_city_name(self, *, slug: str, name: str) -> None:
        await self._conn.execute(sqlalchemy.text(UPDATE_CITY_NAME), {"p1": slug, "p2": name})
