/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 by Martin Mirchev
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
 * Developed by : Bart Kiers, bart@big-o.nl
 */

// $antlr-format alignTrailingComments on, columnLimit 150, maxEmptyLinesToKeep 1, reflowComments off, useTab off
// $antlr-format allowShortRulesOnASingleLine on, alignSemicolons ownLine

lexer grammar SQLiteLexer;

SCOL:      ';';
DOT:       '.';
OPEN_PAR:  '(';
CLOSE_PAR: ')';
COMMA:     ',';
ASSIGN:    '=';
STAR:      '*';
PLUS:      '+';
PTR2:      '->>';
PTR:       '->';
MINUS:     '-';
TILDE:     '~';
PIPE2:     '||';
DIV:       '/';
MOD:       '%';
LT2:       '<<';
GT2:       '>>';
AMP:       '&';
PIPE:      '|';
LT:        '<';
LT_EQ:     '<=';
GT:        '>';
GT_EQ:     '>=';
EQ:        '==';
NOT_EQ1:   '!=';
NOT_EQ2:   '<>';

// http://www.sqlite.org/lang_keywords.html
ABORT_:             A B O R T;
ACTION_:            A C T I O N;
ADD_:               A D D;
AFTER_:             A F T E R;
ALL_:               A L L;
ALTER_:             A L T E R;
ANALYZE_:           A N A L Y Z E;
AND_:               A N D;
AS_:                A S;
ASC_:               A S C;
ATTACH_:            A T T A C H;
AUTOINCREMENT_:     A U T O I N C R E M E N T;
BEFORE_:            B E F O R E;
BEGIN_:             B E G I N;
BETWEEN_:           B E T W E E N;
BY_:                B Y;
CASCADE_:           C A S C A D E;
CASE_:              C A S E;
CAST_:              C A S T;
CHECK_:             C H E C K;
COLLATE_:           C O L L A T E;
COLUMN_:            C O L U M N;
COMMIT_:            C O M M I T;
CONFLICT_:          C O N F L I C T;
CONSTRAINT_:        C O N S T R A I N T;
CREATE_:            C R E A T E;
CROSS_:             C R O S S;
CURRENT_DATE_:      C U R R E N T '_' D A T E;
CURRENT_TIME_:      C U R R E N T '_' T I M E;
CURRENT_TIMESTAMP_: C U R R E N T '_' T I M E S T A M P;
DATABASE_:          D A T A B A S E;
DEFAULT_:           D E F A U L T;
DEFERRABLE_:        D E F E R R A B L E;
DEFERRED_:          D E F E R R E D;
DELETE_:            D E L E T E;
DESC_:              D E S C;
DETACH_:            D E T A C H;
DISTINCT_:          D I S T I N C T;
DROP_:              D R O P;
EACH_:              E A C H;
ELSE_:              E L S E;
END_:               E N D;
ESCAPE_:            E S C A P E;
EXCEPT_:            E X C E P T;
EXCLUSIVE_:         E X C L U S I V E;
EXISTS_:            E X I S T S;
EXPLAIN_:           E X P L A I N;
FAIL_:              F A I L;
FOR_:               F O R;
FOREIGN_:           F O R E I G N;
FROM_:              F R O M;
FULL_:              F U L L;
GLOB_:              G L O B;
GROUP_:             G R O U P;
HAVING_:            H A V I N G;
IF_:                I F;
IGNORE_:            I G N O R E;
IMMEDIATE_:         I M M E D I A T E;
IN_:                I N;
INDEX_:             I N D E X;
INDEXED_:           I N D E X E D;
INITIALLY_:         I N I T I A L L Y;
INNER_:             I N N E R;
INSERT_:            I N S E R T;
INSTEAD_:           I N S T E A D;
INTERSECT_:         I N T E R S E C T;
INTO_:              I N T O;
IS_:                I S;
ISNULL_:            I S N U L L;
JOIN_:              J O I N;
KEY_:               K E Y;
LEFT_:              L E F T;
LIKE_:              L I K E;
LIMIT_:             L I M I T;
MATCH_:             M A T C H;
NATURAL_:           N A T U R A L;
NO_:                N O;
NOT_:               N O T;
NOTNULL_:           N O T N U L L;
NULL_:              N U L L;
OF_:                O F;
OFFSET_:            O F F S E T;
ON_:                O N;
OR_:                O R;
ORDER_:             O R D E R;
OUTER_:             O U T E R;
PLAN_:              P L A N;
PRAGMA_:            P R A G M A;
PRIMARY_:           P R I M A R Y;
QUERY_:             Q U E R Y;
RAISE_:             R A I S E;
RECURSIVE_:         R E C U R S I V E;
REFERENCES_:        R E F E R E N C E S;
REGEXP_:            R E G E X P;
REINDEX_:           R E I N D E X;
RELEASE_:           R E L E A S E;
RENAME_:            R E N A M E;
REPLACE_:           R E P L A C E;
RESTRICT_:          R E S T R I C T;
RETURNING_:         R E T U R N I N G;
RIGHT_:             R I G H T;
ROLLBACK_:          R O L L B A C K;
ROW_:               R O W;
ROWS_:              R O W S;
SAVEPOINT_:         S A V E P O I N T;
SELECT_:            S E L E C T;
SET_:               S E T;
STRICT_:            S T R I C T;
TABLE_:             T A B L E;
TEMP_:              T E M P;
TEMPORARY_:         T E M P O R A R Y;
THEN_:              T H E N;
TO_:                T O;
TRANSACTION_:       T R A N S A C T I O N;
TRIGGER_:           T R I G G E R;
UNION_:             U N I O N;
UNIQUE_:            U N I Q U E;
UPDATE_:            U P D A T E;
USING_:             U S I N G;
VACUUM_:            V A C U U M;
VALUES_:            V A L U E S;
VIEW_:              V I E W;
VIRTUAL_:           V I R T U A L;
WHEN_:              W H E N;
WHERE_:             W H E R E;
WITH_:              W I T H;
WITHOUT_:           W I T H O U T;
FIRST_VALUE_:       F I R S T '_' V A L U E;
OVER_:              O V E R;
PARTITION_:         P A R T I T I O N;
RANGE_:             R A N G E;
PRECEDING_:         P R E C E D I N G;
UNBOUNDED_:         U N B O U N D E D;
CURRENT_:           C U R R E N T;
FOLLOWING_:         F O L L O W I N G;
CUME_DIST_:         C U M E '_' D I S T;
DENSE_RANK_:        D E N S E '_' R A N K;
LAG_:               L A G;
LAST_VALUE_:        L A S T '_' V A L U E;
LEAD_:              L E A D;
NTH_VALUE_:         N T H '_' V A L U E;
NTILE_:             N T I L E;
PERCENT_RANK_:      P E R C E N T '_' R A N K;
RANK_:              R A N K;
ROW_NUMBER_:        R O W '_' N U M B E R;
GENERATED_:         G E N E R A T E D;
ALWAYS_:            A L W A Y S;
STORED_:            S T O R E D;
TRUE_:              T R U E;
FALSE_:             F A L S E;
WINDOW_:            W I N D O W;
NULLS_:             N U L L S;
FIRST_:             F I R S T;
LAST_:              L A S T;
FILTER_:            F I L T E R;
GROUPS_:            G R O U P S;
EXCLUDE_:           E X C L U D E;
TIES_:              T I E S;
OTHERS_:            O T H E R S;
DO_:                D O;
NOTHING_:           N O T H I N G;

IDENTIFIER:
    '"' (~'"' | '""')* '"'
    | '`' (~'`' | '``')* '`'
    | '[' ~']'* ']'
    | [a-zA-Z_] [a-zA-Z_0-9]*
; // TODO check: needs more chars in set

NUMERIC_LITERAL: ((DIGIT+ ('.' DIGIT*)?) | ('.' DIGIT+)) (E [-+]? DIGIT+)? | '0x' HEX_DIGIT+;

NUMBERED_BIND_PARAMETER: '?' DIGIT*;

NAMED_BIND_PARAMETER: [:@$] IDENTIFIER;

STRING_LITERAL: '\'' ( ~'\'' | '\'\'')* '\'';

BLOB_LITERAL: X STRING_LITERAL;

SINGLE_LINE_COMMENT: '--' ~[\r\n]* (('\r'? '\n') | EOF) -> channel(HIDDEN);

MULTILINE_COMMENT: '/*' .*? '*/' -> channel(HIDDEN);

SPACES: [ \u000B\t\r\n] -> channel(HIDDEN);

UNEXPECTED_CHAR: .;

fragment HEX_DIGIT: [0-9a-fA-F];
fragment DIGIT:     [0-9];

fragment A: [aA];
fragment B: [bB];
fragment C: [cC];
fragment D: [dD];
fragment E: [eE];
fragment F: [fF];
fragment G: [gG];
fragment H: [hH];
fragment I: [iI];
fragment J: [jJ];
fragment K: [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment Q: [qQ];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];
fragment Y: [yY];
fragment Z: [zZ];
