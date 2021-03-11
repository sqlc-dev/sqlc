import asyncio
import os
import random

import asyncpg
import psycopg2
import psycopg2.extensions
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
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    engine = sqlalchemy.create_engine(postgres_uri)
    with engine.connect() as conn:
        conn.execute(f"CREATE SCHEMA {schema_name}")
        conn.execute(f"SET search_path TO {schema_name}")
        yield conn
        conn.execute(f"DROP SCHEMA {schema_name} CASCADE")
        conn.execute("SET search_path TO public")


@pytest.fixture(scope="session")
async def async_sqlalchemy_connection(postgres_uri) -> sqlalchemy.ext.asyncio.AsyncConnection:
    postgres_uri = postgres_uri.replace("postgresql", "postgresql+asyncpg")
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    engine = sqlalchemy.ext.asyncio.create_async_engine(postgres_uri)
    async with engine.connect() as conn:
        await conn.execute(sqlalchemy.text(f"CREATE SCHEMA {schema_name}"))
        await conn.execute(sqlalchemy.text(f"SET search_path TO {schema_name}"))
        await conn.commit()
        yield conn
        await conn.rollback()
        await conn.execute(sqlalchemy.text(f"DROP SCHEMA {schema_name} CASCADE"))
        await conn.execute(sqlalchemy.text("SET search_path TO public"))


@pytest.fixture(scope="session")
def postgres_connection(postgres_uri) -> psycopg2.extensions.connection:
    conn = psycopg2.connect(postgres_uri)
    yield conn
    conn.close()


@pytest.fixture()
def postgres_db(postgres_connection) -> psycopg2.extensions.connection:
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    # schema_name = "sqltest_1"
    cur = postgres_connection.cursor()
    cur.execute(f"CREATE SCHEMA {schema_name}")
    cur.execute(f"SET search_path TO {schema_name}")
    cur.close()
    postgres_connection.commit()
    yield postgres_connection
    postgres_connection.rollback()
    cur = postgres_connection.cursor()
    cur.execute(f"DROP SCHEMA {schema_name} CASCADE")
    cur.execute(f"SET search_path TO public")
    cur.close()
    postgres_connection.commit()


@pytest.fixture(scope="session")
def event_loop():
    """Change event_loop fixture to session level."""
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()


@pytest.fixture(scope="session")
async def async_postgres_connection(postgres_uri: str) -> asyncpg.Connection:
    conn = await asyncpg.connect(postgres_uri)
    yield conn
    await conn.close()


@pytest.fixture()
async def async_postgres_db(async_postgres_connection: asyncpg.Connection) -> asyncpg.Connection:
    conn = async_postgres_connection
    schema_name = f"sqltest_{random.randint(0, 1000)}"
    await conn.execute(f"CREATE SCHEMA {schema_name}")
    await conn.execute(f"SET search_path TO {schema_name}")
    yield conn
    await conn.execute(f"DROP SCHEMA {schema_name} CASCADE")
    await conn.execute(f"SET search_path TO public")
