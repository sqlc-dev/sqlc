package sqlc.runtime

import java.sql.Statement
import java.time.Duration

abstract class Query {
    @Volatile
    var timeout: Duration? = null

    private var _statement: Statement? = null
    protected var statement: Statement
        @Synchronized get() {
            return this._statement ?: throw Exception("Cannot get Statement before Query is executed")
        }
        @Synchronized set(value) {
            val timeout = this.timeout
            if (timeout != null) {
                value.queryTimeout = timeout.seconds.toInt()
            }
            this._statement = value
        }

    fun cancel() {
        this.statement.cancel()
    }
}

abstract class RowQuery<T> : Query() {
    abstract fun execute(): T
}

abstract class ListQuery<T> : Query() {
    abstract fun execute(): List<T>
}

abstract class ExecuteQuery : Query() {
    abstract fun execute()
}

abstract class ExecuteUpdateQuery : Query() {
    abstract fun execute(): Int
}