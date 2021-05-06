import datetime
import os

import pytest
import sqlalchemy.ext.asyncio

from booktest import query, models
from dbtest.migrations import apply_migrations_async


@pytest.mark.asyncio
async def test_books(async_db: sqlalchemy.ext.asyncio.AsyncConnection):
    await apply_migrations_async(async_db, [os.path.dirname(__file__) + "/../../../booktest/postgresql/schema.sql"])

    querier = query.AsyncQuerier(async_db)

    author = await querier.create_author(name="Unknown Master")
    assert author is not None

    now = datetime.datetime.now()
    await querier.create_book(query.CreateBookParams(
        author_id=author.author_id,
        isbn="1",
        title="my book title",
        book_type=models.BookType.FICTION,
        year=2016,
        available=now,
        tags=[],
    ))

    b1 = await querier.create_book(query.CreateBookParams(
        author_id=author.author_id,
        isbn="2",
        title="the second book",
        book_type=models.BookType.FICTION,
        year=2016,
        available=now,
        tags=["cool", "unique"],
    ))

    await querier.update_book(book_id=b1.book_id, title="changed second title", tags=["cool", "disastor"])

    b3 = await querier.create_book(query.CreateBookParams(
        author_id=author.author_id,
        isbn="3",
        title="the third book",
        book_type=models.BookType.FICTION,
        year=2001,
        available=now,
        tags=["cool"],
    ))

    b4 = await querier.create_book(query.CreateBookParams(
        author_id=author.author_id,
        isbn="4",
        title="4th place finisher",
        book_type=models.BookType.NONFICTION,
        year=2011,
        available=now,
        tags=["other"],
    ))

    await querier.update_book_isbn(book_id=b4.book_id, isbn="NEW ISBN", title="never ever gonna finish, a quatrain", tags=["someother"])

    books0 = querier.books_by_title_year(title="my book title", year=2016)
    expected_titles = {"my book title"}
    async for book in books0:
        expected_titles.remove(book.title)  # raises a key error if the title does not exist
        assert len(book.tags) == 0

        author = await querier.get_author(author_id=book.author_id)
        assert author.name == "Unknown Master"
    assert len(expected_titles) == 0

    books = querier.books_by_tags(dollar_1=["cool", "other", "someother"])
    expected_titles = {"changed second title", "the third book", "never ever gonna finish, a quatrain"}
    async for book in books:
        expected_titles.remove(book.title)
    assert len(expected_titles) == 0

    b5 = await querier.get_book(book_id=b3.book_id)
    assert b5 is not None
    await querier.delete_book(book_id=b5.book_id)
    b6 = await querier.get_book(book_id=b5.book_id)
    assert b6 is None
