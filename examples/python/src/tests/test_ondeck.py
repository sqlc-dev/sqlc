import os

import asyncpg
import pytest
from sqlc_runtime.asyncpg import build_asyncpg_connection

from ondeck import models
from ondeck import city as city_queries
from ondeck import venue as venue_queries
from dbtest.migrations import apply_migrations_async


@pytest.mark.asyncio
async def test_ondeck(async_postgres_db: asyncpg.Connection):
    await apply_migrations_async(async_postgres_db, [os.path.dirname(__file__) + "/../../../ondeck/postgresql/schema"])

    db = build_asyncpg_connection(async_postgres_db)

    city = await city_queries.create_city(db, slug="san-francisco", name="San Francisco")
    assert city is not None

    venue_id = await venue_queries.create_venue(db, venue_queries.CreateVenueParams(
        slug="the-fillmore",
        name="The Fillmore",
        city=city.slug,
        spotify_playlist="spotify:uri",
        status=models.Status.OPEN,
        statuses=[models.Status.OPEN, models.Status.CLOSED],
        tags=["rock", "punk"],
    ))
    assert venue_id is not None

    venue = await venue_queries.get_venue(db, slug="the-fillmore", city=city.slug)
    assert venue is not None
    assert venue.id == venue_id

    assert city == await city_queries.get_city(db, city.slug)
    assert [venue_queries.VenueCountByCityRow(city=city.slug, count=1)] == await _to_list(venue_queries.venue_count_by_city(db))
    assert [city] == await _to_list(city_queries.list_cities(db))
    assert [venue] == await _to_list(venue_queries.list_venues(db, city=city.slug))

    await city_queries.update_city_name(db, slug=city.slug, name="SF")
    _id = await venue_queries.update_venue_name(db, slug=venue.slug, name="Fillmore")
    assert _id == venue_id

    await venue_queries.delete_venue(db, slug=venue.slug)


async def _to_list(it):
    out = []
    async for i in it:
        out.append(i)
    return out
