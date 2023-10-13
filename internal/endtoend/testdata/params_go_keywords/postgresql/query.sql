-- name: KeywordBreak :exec
SELECT sqlc.arg('break')::text;

-- name: KeywordDefault :exec
SELECT sqlc.arg('default')::text;

-- name: KeywordFunc :exec
SELECT sqlc.arg('func')::text;

-- name: KeywordInterface :exec
SELECT sqlc.arg('interface')::text;

-- name: KeywordSelect :exec
SELECT sqlc.arg('select')::text;

-- name: KeywordCase :exec
SELECT sqlc.arg('case')::text;

-- name: KeywordDefer :exec
SELECT sqlc.arg('defer')::text;

-- name: KeywordGo :exec
SELECT sqlc.arg('go')::text;

-- name: KeywordMap :exec
SELECT sqlc.arg('map')::text;

-- name: KeywordStruct :exec
SELECT sqlc.arg('struct')::text;

-- name: KeywordChan :exec
SELECT sqlc.arg('chan')::text;

-- name: KeywordElse :exec
SELECT sqlc.arg('else')::text;

-- name: KeywordGoto :exec
SELECT sqlc.arg('goto')::text;

-- name: KeywordPackage :exec
SELECT sqlc.arg('package')::text;

-- name: KeywordSwitch :exec
SELECT sqlc.arg('switch')::text;

-- name: KeywordConst :exec
SELECT sqlc.arg('const')::text;

-- name: KeywordFallthrough :exec
SELECT sqlc.arg('fallthrough')::text;

-- name: KeywordIf :exec
SELECT sqlc.arg('if')::text;

-- name: KeywordRange :exec
SELECT sqlc.arg('range')::text;

-- name: KeywordType :exec
SELECT sqlc.arg('type')::text;

-- name: KeywordContinue :exec
SELECT sqlc.arg('continue')::text;

-- name: KeywordFor :exec
SELECT sqlc.arg('for')::text;

-- name: KeywordImport :exec
SELECT sqlc.arg('import')::text;

-- name: KeywordReturn :exec
SELECT sqlc.arg('return')::text;

-- name: KeywordVar :exec
SELECT sqlc.arg('var')::text;

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
