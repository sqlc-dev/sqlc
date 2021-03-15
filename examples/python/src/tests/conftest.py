import asyncio
import os
import random

import pytest
import sqlalchemy
import sqlalchemy.ext.asyncio


@pytest.fixture(scope="session")
def postgres_uri() -> str:
    pg_host = os.environ.get("PG_HOST", "postgres")
    pg_port = os.environ.get("PG_PORT", 5432)
    pg_user = os.environ.get("PG_USER", "postgres")
    pg_password = os.environ.get("PG_PASSWORD", "mysecretpassword")
    pg_db = os.environ.get("PG_DATABASE", "dinotest")

    return f"postgresql://{pg_user}:{pg_password}@{pg_host}:{pg_port}/{pg_db}"


@pytest.fixture(scope="session")
def sqlalchemy_connection(postgres_uri) -> sqlalchemy.engine.Connection:
    engine = sqlalchemy.create_engine(postgres_uri, future=True)
    with engine.connect() as conn:
        yield conn


@pytest.fixture(scope="function")
def db(sqlalchemy_connection: sqlalchemy.engine.Connection) -> sqlalchemy.engine.Connection:
    conn = sqlalchemy_connection
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    conn.execute(sqlalchemy.text(f"CREATE SCHEMA {schema_name}"))
    conn.execute(sqlalchemy.text(f"SET search_path TO {schema_name}"))
    conn.commit()
    yield conn
    conn.rollback()
    conn.execute(sqlalchemy.text(f"DROP SCHEMA {schema_name} CASCADE"))
    conn.execute(sqlalchemy.text("SET search_path TO public"))


@pytest.fixture(scope="session")
async def async_sqlalchemy_connection(postgres_uri) -> sqlalchemy.ext.asyncio.AsyncConnection:
    postgres_uri = postgres_uri.replace("postgresql", "postgresql+asyncpg")
    engine = sqlalchemy.ext.asyncio.create_async_engine(postgres_uri)
    async with engine.connect() as conn:
        yield conn


@pytest.fixture(scope="function")
async def async_db(async_sqlalchemy_connection: sqlalchemy.ext.asyncio.AsyncConnection) -> sqlalchemy.ext.asyncio.AsyncConnection:
    conn = async_sqlalchemy_connection
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    await conn.execute(sqlalchemy.text(f"CREATE SCHEMA {schema_name}"))
    await conn.execute(sqlalchemy.text(f"SET search_path TO {schema_name}"))
    await conn.commit()
    yield conn
    await conn.rollback()
    await conn.execute(sqlalchemy.text(f"DROP SCHEMA {schema_name} CASCADE"))
    await conn.execute(sqlalchemy.text("SET search_path TO public"))


@pytest.fixture(scope="session")
def event_loop():
    """Change event_loop fixture to session level."""
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()
