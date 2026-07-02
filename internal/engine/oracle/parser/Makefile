# Oracle PL/SQL parser generation.
#
# The grammar (PlSqlLexer.g4, PlSqlParser.g4) is vendored from antlr/grammars-v4
# (sql/plsql) and has already been transformed for the Go target (embedded `this.`
# actions rewritten to `p.`/`l.` via transformGrammar.py). The hand-written base
# classes plsql_lexer_base.go and plsql_parser_base.go supply the grammar
# `superClass` implementations and are committed.
#
# The generated *.go files are committed, so `go build` needs no Java. Run `make`
# here only to regenerate after updating the grammar.

plsql_parser.go: PlSqlLexer.g4 PlSqlParser.g4 antlr-4.13.1-complete.jar
	java -jar antlr-4.13.1-complete.jar -Dlanguage=Go -package parser PlSqlLexer.g4 PlSqlParser.g4

antlr-4.13.1-complete.jar:
	curl -fsSL -o antlr-4.13.1-complete.jar https://repo1.maven.org/maven2/org/antlr/antlr4/4.13.1/antlr4-4.13.1-complete.jar

# Refresh the grammar from antlr/grammars-v4, then re-apply the Go-target transform.
# NOTE: the Go base classes live under sql/plsql/Go in the upstream repo.
.PHONY: fetch-grammar
fetch-grammar:
	curl -O https://raw.githubusercontent.com/antlr/grammars-v4/master/sql/plsql/PlSqlLexer.g4
	curl -O https://raw.githubusercontent.com/antlr/grammars-v4/master/sql/plsql/PlSqlParser.g4
	python3 transformGrammar.py || true

.PHONY: clean
clean:
	rm -f *.interp *.tokens *.bak plsql_lexer.go plsql_parser.go plsqlparser_base_listener.go plsqlparser_listener.go
