package com.example.authors

import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import java.nio.file.Files
import java.nio.file.Paths
import java.sql.Connection
import java.sql.DriverManager

const val schema = "dinosql_test"

class QueriesImplTest {
    lateinit var schemaConn: Connection
    lateinit var conn: Connection

    @BeforeEach
    fun setup() {
        val user = System.getenv("PG_USER") ?: "postgres"
        val pass = System.getenv("PG_PASSWORD") ?: "mysecretpassword"
        val host = System.getenv("PG_HOST") ?: "127.0.0.1"
        val port = System.getenv("PG_PORT") ?: "5432"
        val db = System.getenv("PG_DATABASE") ?: "dinotest"
        val url = "jdbc:postgresql://$host:$port/$db?user=$user&password=$pass&sslmode=disable"
        println("db: $url")

        schemaConn = DriverManager.getConnection(url)
        schemaConn.createStatement().execute("CREATE SCHEMA $schema")

        conn = DriverManager.getConnection("$url&currentSchema=$schema")
        val stmt = Files.readString(Paths.get("src/main/resources/schema.sql"))
        conn.createStatement().execute(stmt)
    }

    @AfterEach
    fun teardown() {
        schemaConn.createStatement().execute("DROP SCHEMA $schema CASCADE")
    }

    @Test
    fun testCreateAuthor() {
        val db = QueriesImpl(conn)

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val params = CreateAuthorParams(
            name = "Brian Kernighan",
            bio = "Co-author of The C Programming Language and The Go Programming Language"
        )
        val insertedAuthor = db.createAuthor(params)
        val expectedAuthor = Author(insertedAuthor.id, params.name, params.bio)
        assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])
    }

    @Test
    fun testNull() {
        val db = QueriesImpl(conn)

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val params = CreateAuthorParams(
            name = "Brian Kernighan",
            bio = null
        )
        val insertedAuthor = db.createAuthor(params)
        val expectedAuthor = Author(insertedAuthor.id, params.name, params.bio)
        assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])
    }
}
