/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 by Bart Kiers
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 * Project : sqlite-parser; an ANTLR4 grammar for SQLite https://github.com/bkiers/sqlite-parser
 * Developed by:
 *     Bart Kiers, bart@big-o.nl
 *     Martin Mirchev, marti_2203@abv.bg
 *     Mike Lische, mike@lischke-online.de
 */

// $antlr-format alignTrailingComments on, columnLimit 130, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments off
// $antlr-format useTab off, allowShortRulesOnASingleLine off, allowShortBlocksOnASingleLine on, alignSemicolons ownLine

parser grammar SQLiteParser;

options {
    tokenVocab = SQLiteLexer;
}

parse: (sql_stmt_list)* EOF
;

sql_stmt_list:
    SCOL* sql_stmt (SCOL+ sql_stmt)* SCOL*
;

sql_stmt: (EXPLAIN_ (QUERY_ PLAN_)?)? (
        alter_table_stmt
        | analyze_stmt
        | attach_stmt
        | begin_stmt
        | commit_stmt
        | create_index_stmt
        | create_table_stmt
        | create_trigger_stmt
        | create_view_stmt
        | create_virtual_table_stmt
        | delete_stmt
        | delete_stmt_limited
        | detach_stmt
        | drop_stmt
        | insert_stmt
        | pragma_stmt
        | reindex_stmt
        | release_stmt
        | rollback_stmt
        | savepoint_stmt
        | select_stmt
        | update_stmt
        | update_stmt_limited
        | vacuum_stmt
    )
;

alter_table_stmt:
    ALTER_ TABLE_ (schema_name DOT)? table_name (
        RENAME_ (
            TO_ new_table_name
            | COLUMN_? old_column_name = column_name TO_ new_column_name = column_name
        )
        | ADD_ COLUMN_? column_def
        | DROP_ COLUMN_? column_name
    )
;

analyze_stmt:
    ANALYZE_ (schema_name | (schema_name DOT)? table_or_index_name)?
;

attach_stmt:
    ATTACH_ DATABASE_? expr AS_ schema_name
;

begin_stmt:
    BEGIN_ (DEFERRED_ | IMMEDIATE_ | EXCLUSIVE_)? (
        TRANSACTION_ transaction_name?
    )?
;

commit_stmt: (COMMIT_ | END_) TRANSACTION_?
;

rollback_stmt:
    ROLLBACK_ TRANSACTION_? (TO_ SAVEPOINT_? savepoint_name)?
;

savepoint_stmt:
    SAVEPOINT_ savepoint_name
;

release_stmt:
    RELEASE_ SAVEPOINT_? savepoint_name
;

create_index_stmt:
    CREATE_ UNIQUE_? INDEX_ (IF_ NOT_ EXISTS_)? (schema_name DOT)? index_name ON_ table_name OPEN_PAR
        indexed_column (COMMA indexed_column)* CLOSE_PAR (WHERE_ expr)?
;

indexed_column: (column_name | expr) (COLLATE_ collation_name)? asc_desc?
;

table_option:
    WITHOUT_ row_ROW_ID = IDENTIFIER
    | STRICT_
;

create_table_stmt:
    CREATE_ (TEMP_ | TEMPORARY_)? TABLE_ (IF_ NOT_ EXISTS_)? (
        schema_name DOT
    )? table_name (
        OPEN_PAR column_def (COMMA column_def)*? (COMMA table_constraint)* CLOSE_PAR (
            table_option (COMMA table_option)*
        )?
        | AS_ select_stmt
    )
;

column_def:
    column_name type_name? column_constraint*
;

type_name:
    name+? (
        OPEN_PAR signed_number CLOSE_PAR
        | OPEN_PAR signed_number COMMA signed_number CLOSE_PAR
    )?
;

column_constraint: (CONSTRAINT_ name)? (
        (PRIMARY_ KEY_ asc_desc? conflict_clause? AUTOINCREMENT_?)
        | (NOT_ NULL_ | UNIQUE_) conflict_clause?
        | CHECK_ OPEN_PAR expr CLOSE_PAR
        | DEFAULT_ (signed_number | literal_value | OPEN_PAR expr CLOSE_PAR)
        | COLLATE_ collation_name
        | foreign_key_clause
        | (GENERATED_ ALWAYS_)? AS_ OPEN_PAR expr CLOSE_PAR (
            STORED_
            | VIRTUAL_
        )?
    )
;

signed_number: (PLUS | MINUS)? NUMERIC_LITERAL
;

table_constraint: (CONSTRAINT_ name)? (
        (PRIMARY_ KEY_ | UNIQUE_) OPEN_PAR indexed_column (
            COMMA indexed_column
        )* CLOSE_PAR conflict_clause?
        | CHECK_ OPEN_PAR expr CLOSE_PAR
        | FOREIGN_ KEY_ OPEN_PAR column_name (COMMA column_name)* CLOSE_PAR foreign_key_clause
    )
;

foreign_key_clause:
    REFERENCES_ foreign_table (
        OPEN_PAR column_name (COMMA column_name)* CLOSE_PAR
    )? (
        ON_ (DELETE_ | UPDATE_) (
            SET_ (NULL_ | DEFAULT_)
            | CASCADE_
            | RESTRICT_
            | NO_ ACTION_
        )
        | MATCH_ name
    )* (NOT_? DEFERRABLE_ (INITIALLY_ (DEFERRED_ | IMMEDIATE_))?)?
;

conflict_clause:
    ON_ CONFLICT_ (
        ROLLBACK_
        | ABORT_
        | FAIL_
        | IGNORE_
        | REPLACE_
    )
;

create_trigger_stmt:
    CREATE_ (TEMP_ | TEMPORARY_)? TRIGGER_ (IF_ NOT_ EXISTS_)? (
        schema_name DOT
    )? trigger_name (BEFORE_ | AFTER_ | INSTEAD_ OF_)? (
        DELETE_
        | INSERT_
        | UPDATE_ (OF_ column_name ( COMMA column_name)*)?
    ) ON_ table_name (FOR_ EACH_ ROW_)? (WHEN_ expr)? BEGIN_ (
        (update_stmt | insert_stmt | delete_stmt | select_stmt) SCOL
    )+ END_
;

create_view_stmt:
    CREATE_ (TEMP_ | TEMPORARY_)? VIEW_ (IF_ NOT_ EXISTS_)? (
        schema_name DOT
    )? view_name (OPEN_PAR column_name (COMMA column_name)* CLOSE_PAR)? AS_ select_stmt
;

create_virtual_table_stmt:
    CREATE_ VIRTUAL_ TABLE_ (IF_ NOT_ EXISTS_)? (schema_name DOT)? table_name USING_ module_name (
        OPEN_PAR module_argument (COMMA module_argument)* CLOSE_PAR
    )?
;

with_clause:
    WITH_ RECURSIVE_? cte_table_name AS_ OPEN_PAR select_stmt CLOSE_PAR (
        COMMA cte_table_name AS_ OPEN_PAR select_stmt CLOSE_PAR
    )*
;

cte_table_name:
    table_name (OPEN_PAR column_name ( COMMA column_name)* CLOSE_PAR)?
;

recursive_cte:
    cte_table_name AS_ OPEN_PAR initial_select UNION_ ALL_? recursive__select CLOSE_PAR
;

common_table_expression:
    table_name (OPEN_PAR column_name ( COMMA column_name)* CLOSE_PAR)? AS_ OPEN_PAR select_stmt CLOSE_PAR
;

returning_clause:
    RETURNING_ (
        (STAR | expr ( AS_? column_alias)?) (
            COMMA (STAR | expr ( AS_? column_alias)?)
        )*
    )
;

delete_stmt:
    with_clause? DELETE_ FROM_ qualified_table_name (WHERE_ expr)? returning_clause?
;

delete_stmt_limited:
    with_clause? DELETE_ FROM_ qualified_table_name (WHERE_ expr)? (
        order_by_stmt? limit_stmt
    )? returning_clause?
;

detach_stmt:
    DETACH_ DATABASE_? schema_name
;

drop_stmt:
    DROP_ object = (INDEX_ | TABLE_ | TRIGGER_ | VIEW_) (
        IF_ EXISTS_
    )? (schema_name DOT)? any_name
;

/*
 SQLite understands the following binary operators, in order from highest to lowest precedence:
    ||
    * / %
    + -
    << >> & |
    < <= > >=
    = == != <> IS IS NOT IN LIKE GLOB MATCH REGEXP
    AND
    OR
 */
expr:
    literal_value #expr_literal
    | NUMBERED_BIND_PARAMETER #expr_bind
    | NAMED_BIND_PARAMETER #expr_bind
    | ((schema_name DOT)? table_name DOT)? column_name #expr_qualified_column_name
    | unary_operator expr #expr_unary
    | expr PIPE2 expr #expr_binary
    | expr ( PTR | PTR2 ) expr #expr_binary
    | expr ( STAR | DIV | MOD) expr #expr_binary
    | expr ( PLUS | MINUS) expr #expr_binary
    | expr ( LT2 | GT2 | AMP | PIPE) expr #expr_comparison
    | expr ( LT | LT_EQ | GT | GT_EQ) expr #expr_comparison
    | expr (
        ASSIGN
        | EQ
        | NOT_EQ1
        | NOT_EQ2
        | IS_
        | IS_ NOT_
        | NOT_? IN_
        | LIKE_
        | GLOB_
        | MATCH_
        | REGEXP_
    ) expr #expr_comparison
    | expr NOT_? IN_ (
        OPEN_PAR (select_stmt | expr ( COMMA expr)*)? CLOSE_PAR
        | ( schema_name DOT)? table_name
        | (schema_name DOT)? table_function_name OPEN_PAR (expr (COMMA expr)*)? CLOSE_PAR
    ) #expr_in_select
    | expr AND_ expr #expr_bool
    | expr OR_ expr #expr_bool
    | qualified_function_name OPEN_PAR ((DISTINCT_? expr ( COMMA expr)*) | STAR)? CLOSE_PAR filter_clause? over_clause? #expr_function
    | OPEN_PAR expr (COMMA expr)* CLOSE_PAR #expr_list
    | CAST_ OPEN_PAR expr AS_ type_name CLOSE_PAR #expr_cast
    | expr COLLATE_ collation_name #expr_collate
    | expr NOT_? (LIKE_ | GLOB_ | REGEXP_ | MATCH_) expr (
        ESCAPE_ expr
    )? #expr_comparison
    | expr ( ISNULL_ | NOTNULL_ | NOT_ NULL_) #expr_null_comp
    | expr NOT_? BETWEEN_ expr AND_ expr #expr_between
    | ((NOT_)? EXISTS_)? OPEN_PAR select_stmt CLOSE_PAR #expr_in_select
    | CASE_ expr? (WHEN_ expr THEN_ expr)+ (ELSE_ expr)? END_ #expr_case
    | raise_function #expr_raise
;

raise_function:
    RAISE_ OPEN_PAR (
        IGNORE_
        | (ROLLBACK_ | ABORT_ | FAIL_) COMMA error_message
    ) CLOSE_PAR
;

literal_value:
    NUMERIC_LITERAL
    | STRING_LITERAL
    | BLOB_LITERAL
    | NULL_
    | TRUE_
    | FALSE_
    | CURRENT_TIME_
    | CURRENT_DATE_
    | CURRENT_TIMESTAMP_
;

insert_stmt:
    with_clause? (
        INSERT_
        | REPLACE_
        | INSERT_ OR_ (
            REPLACE_
            | ROLLBACK_
            | ABORT_
            | FAIL_
            | IGNORE_
        )
    ) INTO_ (schema_name DOT)? table_name (AS_ table_alias)? (
        OPEN_PAR column_name ( COMMA column_name)* CLOSE_PAR
    )? (
        (
            VALUES_ OPEN_PAR expr (COMMA expr)* CLOSE_PAR (
                COMMA OPEN_PAR expr ( COMMA expr)* CLOSE_PAR
            )*
            | select_stmt
            | DEFAULT_ VALUES_
        ) upsert_clause? returning_clause?
    )
;

upsert_clause:
    ON_ CONFLICT_ (
        OPEN_PAR indexed_column (COMMA indexed_column)* CLOSE_PAR (WHERE_ expr)?
    )? DO_ (
        NOTHING_
        | UPDATE_ SET_ (
            (column_name | column_name_list) ASSIGN expr (
                COMMA (column_name | column_name_list) ASSIGN expr
            )* (WHERE_ expr)?
        )
    )
;

pragma_stmt:
    PRAGMA_ (schema_name DOT)? pragma_name (
        ASSIGN pragma_value
        | OPEN_PAR pragma_value CLOSE_PAR
    )?
;

pragma_value:
    signed_number
    | name
    | STRING_LITERAL
;

reindex_stmt:
    REINDEX_ (collation_name | (schema_name DOT)? (table_name | index_name))?
;

select_stmt:
    common_table_stmt? select_core (compound_operator select_core)* order_by_stmt? limit_stmt?
;

join_clause:
    table_or_subquery (join_operator table_or_subquery join_constraint)*
;

select_core:
    (
        SELECT_ (DISTINCT_ | ALL_)? result_column (COMMA result_column)* (
            FROM_ (table_or_subquery (COMMA table_or_subquery)* | join_clause)
        )? (WHERE_ expr)? (GROUP_ BY_ expr (COMMA expr)* (HAVING_ expr)?)? (
            WINDOW_ window_name AS_ window_defn (
                COMMA window_name AS_ window_defn
            )*
        )?
    )
    | VALUES_ OPEN_PAR expr (COMMA expr)* CLOSE_PAR (
        COMMA OPEN_PAR expr ( COMMA expr)* CLOSE_PAR
    )*
;

factored_select_stmt:
    select_stmt
;

simple_select_stmt:
    common_table_stmt? select_core order_by_stmt? limit_stmt?
;

compound_select_stmt:
    common_table_stmt? select_core (
        (UNION_ ALL_? | INTERSECT_ | EXCEPT_) select_core
    )+ order_by_stmt? limit_stmt?
;

table_or_subquery:
    (schema_name DOT)? table_name (AS_? table_alias)? (INDEXED_ BY_ index_name | NOT_ INDEXED_)?
    | (schema_name DOT)? table_function_name OPEN_PAR expr (COMMA expr)* CLOSE_PAR (AS_? table_alias)?
    | OPEN_PAR (table_or_subquery (COMMA table_or_subquery)* | join_clause) CLOSE_PAR
    | OPEN_PAR select_stmt CLOSE_PAR (AS_? table_alias)?
    | (schema_name DOT)? table_name (AS_? table_alias_fallback)? (INDEXED_ BY_ index_name | NOT_ INDEXED_)?
    | (schema_name DOT)? table_function_name OPEN_PAR expr (COMMA expr)* CLOSE_PAR (AS_? table_alias_fallback)?
    | OPEN_PAR (table_or_subquery (COMMA table_or_subquery)* | join_clause) CLOSE_PAR
    | OPEN_PAR select_stmt CLOSE_PAR (AS_? table_alias_fallback)?
;

result_column:
    STAR
    | table_name DOT STAR
    | expr ( AS_? column_alias)?
;

join_operator:
    COMMA
    | NATURAL_? (((LEFT_ | RIGHT_ | FULL_) OUTER_?) | INNER_)? JOIN_
    | CROSS_ JOIN_
;

join_constraint:
    (ON_ expr
    | USING_ OPEN_PAR column_name ( COMMA column_name)* CLOSE_PAR)?
;

compound_operator:
    UNION_ ALL_?
    | INTERSECT_
    | EXCEPT_
;

update_stmt:
    with_clause? UPDATE_ (
        OR_ (ROLLBACK_ | ABORT_ | REPLACE_ | FAIL_ | IGNORE_)
    )? qualified_table_name SET_ (column_name | column_name_list) ASSIGN expr (
        COMMA (column_name | column_name_list) ASSIGN expr
    )* (WHERE_ expr)? returning_clause?
;

column_name_list:
    OPEN_PAR column_name (COMMA column_name)* CLOSE_PAR
;

update_stmt_limited:
    with_clause? UPDATE_ (
        OR_ (ROLLBACK_ | ABORT_ | REPLACE_ | FAIL_ | IGNORE_)
    )? qualified_table_name SET_ (column_name | column_name_list) ASSIGN expr (
        COMMA (column_name | column_name_list) ASSIGN expr
    )* (WHERE_ expr)? (order_by_stmt? limit_stmt)?
;

qualified_table_name: (schema_name DOT)? table_name (AS_ alias)? (
        INDEXED_ BY_ index_name
        | NOT_ INDEXED_
    )?
;

vacuum_stmt:
    VACUUM_ schema_name? (INTO_ filename)?
;

filter_clause:
    FILTER_ OPEN_PAR WHERE_ expr CLOSE_PAR
;

window_defn:
    OPEN_PAR base_window_name? (PARTITION_ BY_ expr (COMMA expr)*)? (
        ORDER_ BY_ ordering_term (COMMA ordering_term)*
    ) frame_spec? CLOSE_PAR
;

over_clause:
    OVER_ (
        window_name
        | OPEN_PAR base_window_name? (PARTITION_ BY_ expr (COMMA expr)*)? (
            ORDER_ BY_ ordering_term (COMMA ordering_term)*
        )? frame_spec? CLOSE_PAR
    )
;

frame_spec:
    frame_clause (
        EXCLUDE_ (NO_ OTHERS_)
        | CURRENT_ ROW_
        | GROUP_
        | TIES_
    )?
;

frame_clause: (RANGE_ | ROWS_ | GROUPS_) (
        frame_single
        | BETWEEN_ frame_left AND_ frame_right
    )
;

simple_function_invocation:
    simple_func OPEN_PAR (expr (COMMA expr)* | STAR) CLOSE_PAR
;

aggregate_function_invocation:
    aggregate_func OPEN_PAR (DISTINCT_? expr (COMMA expr)* | STAR)? CLOSE_PAR filter_clause?
;

window_function_invocation:
    window_function OPEN_PAR (expr (COMMA expr)* | STAR)? CLOSE_PAR filter_clause? OVER_ (
        window_defn
        | window_name
    )
;

common_table_stmt: //additional structures
    WITH_ RECURSIVE_? common_table_expression (COMMA common_table_expression)*
;

order_by_stmt:
    ORDER_ BY_ ordering_term (COMMA ordering_term)*
;

limit_stmt:
    LIMIT_ expr ((OFFSET_ | COMMA) expr)?
;

ordering_term:
    expr (COLLATE_ collation_name)? asc_desc? (NULLS_ (FIRST_ | LAST_))?
;

asc_desc:
    ASC_
    | DESC_
;

frame_left:
    expr PRECEDING_
    | expr FOLLOWING_
    | CURRENT_ ROW_
    | UNBOUNDED_ PRECEDING_
;

frame_right:
    expr PRECEDING_
    | expr FOLLOWING_
    | CURRENT_ ROW_
    | UNBOUNDED_ FOLLOWING_
;

frame_single:
    expr PRECEDING_
    | UNBOUNDED_ PRECEDING_
    | CURRENT_ ROW_
;

// unknown

window_function:
    (FIRST_VALUE_ | LAST_VALUE_) OPEN_PAR expr CLOSE_PAR OVER_ OPEN_PAR partition_by? order_by_expr_asc_desc frame_clause
        ? CLOSE_PAR
    | (CUME_DIST_ | PERCENT_RANK_) OPEN_PAR CLOSE_PAR OVER_ OPEN_PAR partition_by? order_by_expr? CLOSE_PAR
    | (DENSE_RANK_ | RANK_ | ROW_NUMBER_) OPEN_PAR CLOSE_PAR OVER_ OPEN_PAR partition_by? order_by_expr_asc_desc
        CLOSE_PAR
    | (LAG_ | LEAD_) OPEN_PAR expr of_OF_fset? default_DEFAULT__value? CLOSE_PAR OVER_ OPEN_PAR partition_by?
        order_by_expr_asc_desc CLOSE_PAR
    | NTH_VALUE_ OPEN_PAR expr COMMA signed_number CLOSE_PAR OVER_ OPEN_PAR partition_by? order_by_expr_asc_desc
        frame_clause? CLOSE_PAR
    | NTILE_ OPEN_PAR expr CLOSE_PAR OVER_ OPEN_PAR partition_by? order_by_expr_asc_desc CLOSE_PAR
;

of_OF_fset:
    COMMA signed_number
;

default_DEFAULT__value:
    COMMA signed_number
;

partition_by:
    PARTITION_ BY_ expr+
;

order_by_expr:
    ORDER_ BY_ expr+
;

order_by_expr_asc_desc:
    ORDER_ BY_ order_by_expr_asc_desc
;

expr_asc_desc:
    expr asc_desc? (COMMA expr asc_desc?)*
;

//TODO BOTH OF THESE HAVE TO BE REWORKED TO FOLLOW THE SPEC
initial_select:
    select_stmt
;

recursive__select:
    select_stmt
;

unary_operator:
    MINUS
    | PLUS
    | TILDE
    | NOT_
;

error_message:
    STRING_LITERAL
;

module_argument: // TODO check what exactly is permitted here
    expr
    | column_def
;

column_alias:
    IDENTIFIER
    | STRING_LITERAL
;

keyword:
    ABORT_
    | ACTION_
    | ADD_
    | AFTER_
    | ALL_
    | ALTER_
    | ANALYZE_
    | AND_
    | AS_
    | ASC_
    | ATTACH_
    | AUTOINCREMENT_
    | BEFORE_
    | BEGIN_
    | BETWEEN_
    | BY_
    | CASCADE_
    | CASE_
    | CAST_
    | CHECK_
    | COLLATE_
    | COLUMN_
    | COMMIT_
    | CONFLICT_
    | CONSTRAINT_
    | CREATE_
    | CROSS_
    | CURRENT_DATE_
    | CURRENT_TIME_
    | CURRENT_TIMESTAMP_
    | DATABASE_
    | DEFAULT_
    | DEFERRABLE_
    | DEFERRED_
    | DELETE_
    | DESC_
    | DETACH_
    | DISTINCT_
    | DROP_
    | EACH_
    | ELSE_
    | END_
    | ESCAPE_
    | EXCEPT_
    | EXCLUSIVE_
    | EXISTS_
    | EXPLAIN_
    | FAIL_
    | FOR_
    | FOREIGN_
    | FROM_
    | FULL_
    | GLOB_
    | GROUP_
    | HAVING_
    | IF_
    | IGNORE_
    | IMMEDIATE_
    | IN_
    | INDEX_
    | INDEXED_
    | INITIALLY_
    | INNER_
    | INSERT_
    | INSTEAD_
    | INTERSECT_
    | INTO_
    | IS_
    | ISNULL_
    | JOIN_
    | KEY_
    | LEFT_
    | LIKE_
    | LIMIT_
    | MATCH_
    | NATURAL_
    | NO_
    | NOT_
    | NOTNULL_
    | NULL_
    | OF_
    | OFFSET_
    | ON_
    | OR_
    | ORDER_
    | OUTER_
    | PLAN_
    | PRAGMA_
    | PRIMARY_
    | QUERY_
    | RAISE_
    | RECURSIVE_
    | REFERENCES_
    | REGEXP_
    | REINDEX_
    | RELEASE_
    | RENAME_
    | REPLACE_
    | RESTRICT_
    | RETURNING_
    | RIGHT_
    | ROLLBACK_
    | ROW_
    | ROWS_
    | SAVEPOINT_
    | SELECT_
    | SET_
    | STRICT_
    | TABLE_
    | TEMP_
    | TEMPORARY_
    | THEN_
    | TO_
    | TRANSACTION_
    | TRIGGER_
    | UNION_
    | UNIQUE_
    | UPDATE_
    | USING_
    | VACUUM_
    | VALUES_
    | VIEW_
    | VIRTUAL_
    | WHEN_
    | WHERE_
    | WITH_
    | WITHOUT_
    | FIRST_VALUE_
    | OVER_
    | PARTITION_
    | RANGE_
    | PRECEDING_
    | UNBOUNDED_
    | CURRENT_
    | FOLLOWING_
    | CUME_DIST_
    | DENSE_RANK_
    | LAG_
    | LAST_VALUE_
    | LEAD_
    | NTH_VALUE_
    | NTILE_
    | PERCENT_RANK_
    | RANK_
    | ROW_NUMBER_
    | GENERATED_
    | ALWAYS_
    | STORED_
    | TRUE_
    | FALSE_
    | WINDOW_
    | NULLS_
    | FIRST_
    | LAST_
    | FILTER_
    | GROUPS_
    | EXCLUDE_
;

// TODO: check all names below

name:
    any_name
;

function_name:
    any_name
;

qualified_function_name:
    (schema_name DOT)? function_name
;

schema_name:
    any_name
;

table_name:
    any_name
;

table_or_index_name:
    any_name
;

new_table_name:
    any_name
;

column_name:
    any_name
;

collation_name:
    any_name
;

foreign_table:
    any_name
;

index_name:
    any_name
;

trigger_name:
    any_name
;

view_name:
    any_name
;

module_name:
    any_name
;

pragma_name:
    any_name
;

savepoint_name:
    any_name
;

table_alias: IDENTIFIER | STRING_LITERAL;

table_alias_fallback: any_name;

transaction_name:
    any_name
;

window_name:
    any_name
;

alias:
    any_name
;

filename:
    any_name
;

base_window_name:
    any_name
;

simple_func:
    any_name
;

aggregate_func:
    any_name
;

table_function_name:
    any_name
;

any_name:
    IDENTIFIER
    | keyword
    | STRING_LITERAL
    | OPEN_PAR any_name CLOSE_PAR
;
