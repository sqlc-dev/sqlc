package com.example.authors

import com.example.dbtest.DbTestExtension
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import org.junit.jupiter.api.extension.RegisterExtension
import java.time.Duration
import java.util.concurrent.TimeoutException

class QueryRuntimeTest {
    @Test
    fun testTimeout() {
        val db = QueriesImpl(QueriesImplTest.dbtest.getConnection())
        assertThrows<TimeoutException> {
            db.deleteAuthor(1).apply {
                timeout = Duration.ZERO
            }.execute()
        }
    }
}