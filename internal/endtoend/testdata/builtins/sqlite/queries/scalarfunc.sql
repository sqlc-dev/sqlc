-- name: GetAbs :one
SELECT abs(int_val) FROM test;

-- name: GetChanges :one
SELECT changes();

-- name: GetChar1 :one
SELECT char(65);

-- name: GetChar3 :one
SELECT char(65, 66, 67);

-- name: GetCoalesce :one
SELECT coalesce(NULL, 1, 'test');

-- name: GetFormat :one
SELECT format('Hello %s', 'world');

-- name: GetGlob :one
SELECT glob('a*c', 'abc');

-- name: GetHex :one
SELECT hex(123456);

-- name: GetIfnull :one
SELECT ifnull(1, 2);

-- name: GetIif :one
SELECT iif(1, 2, 3);

-- name: GetLastInsertRowID :one
SELECT last_insert_rowid();

-- name: GetInstr :one
SELECT instr('hello', 'l');

-- name: GetLength :one
SELECT length('12345');

-- name: GetLike2 :one
SELECT like('%bc%', 'abcd');

-- name: GetLike3 :one
SELECT like('$%1%', '%100', '$');

-- name: GetLikelihood :one
SELECT likelihood('12345', 0.5);

-- name: GetLikely :one
SELECT likely('12345');

-- name: GetLower :one
SELECT lower('ABCDE');

-- name: GetLtrim :one
SELECT ltrim(' ABCDE');

-- name: GetLtrim2 :one
SELECT ltrim(':ABCDE', ':');

-- name: GetMax3 :one
SELECT max(1, 3, 2);

-- name: GetMin3 :one
SELECT min(1, 3, 2);

-- name: GetNullif :one
SELECT nullif(1, 2);

-- name: GetPrintf :one
SELECT printf('Hello %s', 'world');

-- name: GetQuote :one
SELECT quote(1);

-- name: GetRandom :one
SELECT random();

-- name: GetRandomBlob :one
SELECT randomblob(16);

-- name: GetRound :one
SELECT round(1.1);

-- name: GetRound2 :one
SELECT round(1.1, 2);

-- name: GetReplace :one
SELECT replace('abc', 'bc', 'df');

-- name: GetRtrim :one
SELECT rtrim('ABCDE ');

-- name: GetRtrim2 :one
SELECT rtrim('ABCDE:', ':');

-- name: GetSign :one
SELECT sign(1);

-- name: GetSoundex :one
SELECT soundex('abc');

-- name: GetSQLiteCompileOptionGet :one
SELECT sqlite_compileoption_get(1);

-- name: GetSQLiteCompileOptionUsed :one
SELECT sqlite_compileoption_used(1);

-- name: GetSQLiteOffset :one
SELECT sqlite_offset(1);

-- name: GetSQLiteSourceID :one
SELECT sqlite_source_id();

-- name: GetSQLiteVersion :one
SELECT sqlite_version();

-- name: GetSubstr3 :one
SELECT substr('abcdef', 1, 2);

-- name: GetSubstr2 :one
SELECT substr('abcdef', 2);

-- name: GetSusbstring3 :one
SELECT substring('abcdef', 1, 2);

-- name: GetSubstring2 :one
SELECT substring('abcdef', 1);

-- name: GetTotalChanges :one
SELECT total_changes();

-- name: GetTrim :one
SELECT trim(' ABCDE ');

-- name: GetTrim2 :one
SELECT trim(':ABCDE:', ':');

-- name: GetTypeof :one
SELECT typeof('ABCDE');

-- name: GetUnicode :one
SELECT unicode('A');

-- name: GetUnlikely :one
SELECT unlikely('12345');

-- name: GetUpper :one
SELECT upper('abcde');

-- name: GetZeroblob :one
SELECT zeroblob(16);
