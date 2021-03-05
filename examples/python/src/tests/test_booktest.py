import datetime
import os

import asyncpg
import pytest
from sqlc_runtime.asyncpg import build_asyncpg_connection

from booktest import query, models
from dbtest.migrations import apply_migrations_async


@pytest.mark.asyncio
async def test_books(async_postgres_db: asyncpg.Connection):
    await apply_migrations_async(async_postgres_db, [os.path.dirname(__file__) + "/../../../booktest/postgresql/schema.sql"])

    db = build_asyncpg_connection(async_postgres_db)

    author = await query.create_author(db, "Unknown Master")
    assert author is not None

    async with async_postgres_db.transaction():
        now = datetime.datetime.now()
        await query.create_book(db, query.CreateBookParams(
            author_id=author.author_id,
            isbn="1",
            title="my book title",
            book_type=models.BookType.FICTION,
            year=2016,
            available=now,
            tags=[],
        ))

        b1 = await query.create_book(db, query.CreateBookParams(
            author_id=author.author_id,
            isbn="2",
            title="the second book",
            book_type=models.BookType.FICTION,
            year=2016,
            available=now,
            tags=["cool", "unique"],
        ))

        await query.update_book(db, book_id=b1.book_id, title="changed second title", tags=["cool", "disastor"])

        b3 = await query.create_book(db, query.CreateBookParams(
            author_id=author.author_id,
            isbn="3",
            title="the third book",
            book_type=models.BookType.FICTION,
            year=2001,
            available=now,
            tags=["cool"],
        ))

        b4 = await query.create_book(db, query.CreateBookParams(
            author_id=author.author_id,
            isbn="4",
            title="4th place finisher",
            book_type=models.BookType.NONFICTION,
            year=2011,
            available=now,
            tags=["other"],
        ))

    await query.update_book_isbn(db, book_id=b4.book_id, isbn="NEW ISBN", title="never ever gonna finish, a quatrain", tags=["someother"])

    books0 = query.books_by_title_year(db, title="my book title", year=2016)
    expected_titles = {"my book title"}
    async for book in books0:
        expected_titles.remove(book.title)  # raises a key error if the title does not exist
        assert len(book.tags) == 0

        author = await query.get_author(db, author_id=book.author_id)
        assert author.name == "Unknown Master"
    assert len(expected_titles) == 0

    books = query.books_by_tags(db, ["cool", "other", "someother"])
    expected_titles = {"changed second title", "the third book", "never ever gonna finish, a quatrain"}
    async for book in books:
        expected_titles.remove(book.title)
    assert len(expected_titles) == 0

    b5 = await query.get_book(db, b3.book_id)
    assert b5 is not None
    await query.delete_book(db, book_id=b5.book_id)
    b6 = await query.get_book(db, b5.book_id)
    assert b6 is None
