package com.example.dbtest

import org.junit.jupiter.api.extension.AfterEachCallback
import org.junit.jupiter.api.extension.BeforeEachCallback
import org.junit.jupiter.api.extension.ExtensionContext
import java.nio.file.Files
import java.nio.file.Paths
import java.sql.Connection
import java.sql.DriverManager
import kotlin.streams.toList

class PostgresDbTestExtension(private val migrationsPath: String) : BeforeEachCallback, AfterEachCallback {
    private val schemaConn: Connection
    private val url: String

    companion object {
        const val schema = "dinosql_test"
    }

    init {
        val user = System.getenv("PG_USER") ?: "postgres"
        val pass = System.getenv("PG_PASSWORD") ?: "mysecretpassword"
        val host = System.getenv("PG_HOST") ?: "127.0.0.1"
        val port = System.getenv("PG_PORT") ?: "5432"
        val db = System.getenv("PG_DATABASE") ?: "dinotest"
        url = "jdbc:postgresql://$host:$port/$db?user=$user&password=$pass&sslmode=disable"

        schemaConn = DriverManager.getConnection(url)
    }

    override fun beforeEach(context: ExtensionContext) {
        schemaConn.createStatement().execute("CREATE SCHEMA $schema")
        val path = Paths.get(migrationsPath)
        val migrations = if (Files.isDirectory(path)) {
            Files.list(path).filter{ it.toString().endsWith(".sql")}.sorted().map { Files.readString(it) }.toList()
        } else {
            listOf(Files.readString(path))
        }
        migrations.forEach {
            getConnection().createStatement().execute(it)
        }
    }

    override fun afterEach(context: ExtensionContext) {
        schemaConn.createStatement().execute("DROP SCHEMA $schema CASCADE")
    }

    fun getConnection(): Connection {
        return DriverManager.getConnection("$url&currentSchema=$schema")
    }
}