import os

import pytest
import sqlalchemy.ext.asyncio

from ondeck import models
from ondeck import city as city_queries
from ondeck import venue as venue_queries
from dbtest.migrations import apply_migrations_async


@pytest.mark.asyncio
async def test_ondeck(async_db: sqlalchemy.ext.asyncio.AsyncConnection):
    await apply_migrations_async(async_db, [os.path.dirname(__file__) + "/../../../ondeck/postgresql/schema"])

    city_querier = city_queries.AsyncQuerier(async_db)
    venue_querier = venue_queries.AsyncQuerier(async_db)

    city = await city_querier.create_city(slug="san-francisco", name="San Francisco")
    assert city is not None

    venue_id = await venue_querier.create_venue(venue_queries.CreateVenueParams(
        slug="the-fillmore",
        name="The Fillmore",
        city=city.slug,
        spotify_playlist="spotify:uri",
        status=models.Status.OPEN,
        statuses=[models.Status.OPEN, models.Status.CLOSED],
        tags=["rock", "punk"],
    ))
    assert venue_id is not None

    venue = await venue_querier.get_venue(slug="the-fillmore", city=city.slug)
    assert venue is not None
    assert venue.id == venue_id

    assert city == await city_querier.get_city(slug=city.slug)
    assert [venue_queries.VenueCountByCityRow(city=city.slug, count=1)] == await _to_list(venue_querier.venue_count_by_city())
    assert [city] == await _to_list(city_querier.list_cities())
    assert [venue] == await _to_list(venue_querier.list_venues(city=city.slug))

    await city_querier.update_city_name(slug=city.slug, name="SF")
    _id = await venue_querier.update_venue_name(slug=venue.slug, name="Fillmore")
    assert _id == venue_id

    await venue_querier.delete_venue(slug=venue.slug)


async def _to_list(it):
    out = []
    async for i in it:
        out.append(i)
    return out
