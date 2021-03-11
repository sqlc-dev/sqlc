import os

import pytest
import sqlalchemy.ext.asyncio

from authors import query
from dbtest.migrations import apply_migrations, apply_migrations_async


def test_authors(sqlalchemy_connection: sqlalchemy.engine.Connection):
    apply_migrations(sqlalchemy_connection, [os.path.dirname(__file__) + "/../../../authors/postgresql/schema.sql"])

    db = query.Query(sqlalchemy_connection)

    authors = list(db.list_authors())
    assert authors == []

    author_name = "Brian Kernighan"
    author_bio = "Co-author of The C Programming Language and The Go Programming Language"
    new_author = db.create_author(name=author_name, bio=author_bio)
    assert new_author.id > 0
    assert new_author.name == author_name
    assert new_author.bio == author_bio

    db_author = db.get_author(new_author.id)
    assert db_author == new_author

    author_list = list(db.list_authors())
    assert len(author_list) == 1
    assert author_list[0] == new_author


@pytest.mark.asyncio
async def test_authors_async(async_sqlalchemy_connection: sqlalchemy.ext.asyncio.AsyncConnection):
    await apply_migrations_async(async_sqlalchemy_connection, [os.path.dirname(__file__) + "/../../../authors/postgresql/schema.sql"])

    db = query.AsyncQuery(async_sqlalchemy_connection)

    async for _ in db.list_authors():
        assert False, "No authors should exist"

    author_name = "Brian Kernighan"
    author_bio = "Co-author of The C Programming Language and The Go Programming Language"
    new_author = await db.create_author(name=author_name, bio=author_bio)
    assert new_author.id > 0
    assert new_author.name == author_name
    assert new_author.bio == author_bio

    db_author = await db.get_author(new_author.id)
    assert db_author == new_author

    author_list = []
    async for author in db.list_authors():
        author_list.append(author)
    assert len(author_list) == 1
    assert author_list[0] == new_author
