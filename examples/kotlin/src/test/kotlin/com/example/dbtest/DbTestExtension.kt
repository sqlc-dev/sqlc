package com.example.dbtest

import org.junit.jupiter.api.extension.AfterEachCallback
import org.junit.jupiter.api.extension.BeforeEachCallback
import org.junit.jupiter.api.extension.ExtensionContext
import org.junit.jupiter.api.extension.ParameterContext
import org.junit.jupiter.api.extension.ParameterResolver
import java.nio.file.Files
import java.nio.file.Paths
import java.sql.Connection
import java.sql.DriverManager

const val schema = "dinosql_test"

class DbTestExtension(private val migrationsPath: String) : BeforeEachCallback, AfterEachCallback, ParameterResolver {
    private val schemaConn: Connection
    private val url: String

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
        val stmt = Files.readString(Paths.get(migrationsPath))
        getConnection().createStatement().execute(stmt)
    }

    override fun afterEach(context: ExtensionContext) {
        schemaConn.createStatement().execute("DROP SCHEMA $schema CASCADE")
    }

    private fun getConnection(): Connection {
        return DriverManager.getConnection("$url&currentSchema=$schema")
    }

    override fun supportsParameter(parameterContext: ParameterContext, extensionContext: ExtensionContext): Boolean {
        return parameterContext.parameter.type == Connection::class.java
    }

    override fun resolveParameter(parameterContext: ParameterContext, extensionContext: ExtensionContext): Any {
        return getConnection()
    }
}