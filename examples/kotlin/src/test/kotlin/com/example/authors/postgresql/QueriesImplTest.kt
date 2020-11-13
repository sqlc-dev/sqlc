package com.example.authors.postgresql

import com.example.dbtest.PostgresDbTestExtension
import org.junit.jupiter.api.Assertions
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension

class QueriesImplTest() {

    companion object {
        @JvmField
        @RegisterExtension
        val dbtest = PostgresDbTestExtension("src/main/resources/authors/postgresql/schema.sql")
    }

    @Test
    fun testCreateAuthor() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = "Co-author of The C Programming Language and The Go Programming Language"
        val insertedAuthor = db.createAuthor(
            name = name,
            bio = bio
        )!!
        val expectedAuthor = Author(insertedAuthor.id, name, bio)
        Assertions.assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        Assertions.assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        Assertions.assertEquals(1, listedAuthors.size)
        Assertions.assertEquals(expectedAuthor, listedAuthors[0])
    }

    @Test
    fun testNull() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = null
        val insertedAuthor = db.createAuthor(name, bio)!!
        val expectedAuthor = Author(insertedAuthor.id, name, bio)
        Assertions.assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        Assertions.assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        Assertions.assertEquals(1, listedAuthors.size)
        Assertions.assertEquals(expectedAuthor, listedAuthors[0])
    }
}