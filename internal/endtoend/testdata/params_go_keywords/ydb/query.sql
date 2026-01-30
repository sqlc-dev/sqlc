-- name: KeywordBreak :exec
SELECT CAST(sqlc.arg('break') AS Text);

-- name: KeywordDefault :exec
SELECT CAST(sqlc.arg('default') AS Text);

-- name: KeywordFunc :exec
SELECT CAST(sqlc.arg('func') AS Text);

-- name: KeywordInterface :exec
SELECT CAST(sqlc.arg('interface') AS Text);

-- name: KeywordSelect :exec
SELECT CAST(sqlc.arg('select') AS Text);

-- name: KeywordCase :exec
SELECT CAST(sqlc.arg('case') AS Text);

-- name: KeywordDefer :exec
SELECT CAST(sqlc.arg('defer') AS Text);

-- name: KeywordGo :exec
SELECT CAST(sqlc.arg('go') AS Text);

-- name: KeywordMap :exec
SELECT CAST(sqlc.arg('map') AS Text);

-- name: KeywordStruct :exec
SELECT CAST(sqlc.arg('struct') AS Text);

-- name: KeywordChan :exec
SELECT CAST(sqlc.arg('chan') AS Text);

-- name: KeywordElse :exec
SELECT CAST(sqlc.arg('else') AS Text);

-- name: KeywordGoto :exec
SELECT CAST(sqlc.arg('goto') AS Text);

-- name: KeywordPackage :exec
SELECT CAST(sqlc.arg('package') AS Text);

-- name: KeywordSwitch :exec
SELECT CAST(sqlc.arg('switch') AS Text);

-- name: KeywordConst :exec
SELECT CAST(sqlc.arg('const') AS Text);

-- name: KeywordFallthrough :exec
SELECT CAST(sqlc.arg('fallthrough') AS Text);

-- name: KeywordIf :exec
SELECT CAST(sqlc.arg('if') AS Text);

-- name: KeywordRange :exec
SELECT CAST(sqlc.arg('range') AS Text);

-- name: KeywordType :exec
SELECT CAST(sqlc.arg('type') AS Text);

-- name: KeywordContinue :exec
SELECT CAST(sqlc.arg('continue') AS Text);

-- name: KeywordFor :exec
SELECT CAST(sqlc.arg('for') AS Text);

-- name: KeywordImport :exec
SELECT CAST(sqlc.arg('import') AS Text);

-- name: KeywordReturn :exec
SELECT CAST(sqlc.arg('return') AS Text);

-- name: KeywordVar :exec
SELECT CAST(sqlc.arg('var') AS Text);

-- name: KeywordQ :exec
SELECT CAST(sqlc.arg('q') AS Text);

-- name: SelectBreak :one
SELECT "break" FROM go_keywords;

-- name: SelectDefault :one
SELECT "default" FROM go_keywords;

-- name: SelectFunc :one
SELECT "func" FROM go_keywords;

-- name: SelectInterface :one
SELECT "interface" FROM go_keywords;

-- name: SelectSelect :one
SELECT "select" FROM go_keywords;

-- name: SelectCase :one
SELECT "case" FROM go_keywords;

-- name: SelectDefer :one
SELECT "defer" FROM go_keywords;

-- name: SelectGo :one
SELECT "go" FROM go_keywords;

-- name: SelectMap :one
SELECT "map" FROM go_keywords;

-- name: SelectStruct :one
SELECT "struct" FROM go_keywords;

-- name: SelectChan :one
SELECT "chan" FROM go_keywords;

-- name: SelectElse :one
SELECT "else" FROM go_keywords;

-- name: SelectGoto :one
SELECT "goto" FROM go_keywords;

-- name: SelectPackage :one
SELECT "package" FROM go_keywords;

-- name: SelectSwitch :one
SELECT "switch" FROM go_keywords;

-- name: SelectConst :one
SELECT "const" FROM go_keywords;

-- name: SelectFallthrough :one
SELECT "fallthrough" FROM go_keywords;

-- name: SelectIf :one
SELECT "if" FROM go_keywords;

-- name: SelectRange :one
SELECT "range" FROM go_keywords;

-- name: SelectType :one
SELECT "type" FROM go_keywords;

-- name: SelectContinue :one
SELECT "continue" FROM go_keywords;

-- name: SelectFor :one
SELECT "for" FROM go_keywords;

-- name: SelectImport :one
SELECT "import" FROM go_keywords;

-- name: SelectReturn :one
SELECT "return" FROM go_keywords;

-- name: SelectVar :one
SELECT "var" FROM go_keywords;

-- name: SelectQ :one
SELECT "q" FROM go_keywords;
