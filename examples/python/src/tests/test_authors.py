import os

import asyncpg
import psycopg2.extensions
import pytest
from sqlc_runtime.psycopg2 import build_psycopg2_connection
from sqlc_runtime.asyncpg import build_asyncpg_connection

from authors import query
from dbtest.migrations import apply_migrations, apply_migrations_async


def test_authors(postgres_db: psycopg2.extensions.connection):
    apply_migrations(postgres_db, [os.path.dirname(__file__) + "/../../../authors/postgresql/schema.sql"])

    db = build_psycopg2_connection(postgres_db)

    authors = list(query.list_authors(db))
    assert authors == []

    author_name = "Brian Kernighan"
    author_bio = "Co-author of The C Programming Language and The Go Programming Language"
    new_author = query.create_author(db, name=author_name, bio=author_bio)
    assert new_author.id > 0
    assert new_author.name == author_name
    assert new_author.bio == author_bio

    db_author = query.get_author(db, new_author.id)
    assert db_author == new_author

    author_list = list(query.list_authors(db))
    assert len(author_list) == 1
    assert author_list[0] == new_author


@pytest.mark.asyncio
async def test_authors_async(async_postgres_db: asyncpg.Connection):
    await apply_migrations_async(async_postgres_db, [os.path.dirname(__file__) + "/../../../authors/postgresql/schema.sql"])

    db = build_asyncpg_connection(async_postgres_db)

    async for _ in query.list_authors(db):
        assert False, "No authors should exist"

    author_name = "Brian Kernighan"
    author_bio = "Co-author of The C Programming Language and The Go Programming Language"
    new_author = await query.create_author(db, name=author_name, bio=author_bio)
    assert new_author.id > 0
    assert new_author.name == author_name
    assert new_author.bio == author_bio

    db_author = await query.get_author(db, new_author.id)
    assert db_author == new_author

    author_list = []
    async for author in query.list_authors(db):
        author_list.append(author)
    assert len(author_list) == 1
    assert author_list[0] == new_author
