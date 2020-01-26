package com.example.authors

import com.example.dbtest.DbTestExtension
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension
import java.sql.Connection

class QueriesImplTest(private val conn: Connection) {

    companion object {
        @JvmField @RegisterExtension val db = DbTestExtension("src/main/resources/schema.sql")
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
