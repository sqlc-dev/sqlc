package sqlite

import "github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"

// signature is what SQLite does not record about a function: the type it
// returns, whether that can be NULL when its arguments are not, and the
// types its arguments are meant to hold. Args types the leading arguments
// in order; a position it does not cover is "any", which the analyzer
// resolves to the argument's own type. Variadic types the arguments an
// overload of variable arity repeats, "any" when empty. How many arguments
// an overload takes, and how many of them it requires, comes from the
// shell, not from here.
type signature struct {
	Args     []string
	Variadic string
	Returns  string
	Nullable bool
}

// args builds an overload's parameters from the shell's argument count. A
// fixed count is that many parameters; a negative count is the required
// leading parameters, any further ones the signature types marked as
// having a default, and a variadic tail.
func (s signature) args(narg int) []dialect.Arg {
	if narg >= 0 {
		args := make([]dialect.Arg, narg)
		for i := range args {
			args[i] = dialect.Arg{Type: s.argType(i)}
		}
		return args
	}
	required := minArgs(narg)
	n := max(required, len(s.Args))
	args := make([]dialect.Arg, 0, n+1)
	for i := 0; i < n; i++ {
		args = append(args, dialect.Arg{Type: s.argType(i), HasDefault: i >= required})
	}
	tail := s.Variadic
	if tail == "" {
		tail = "any"
	}
	return append(args, dialect.Arg{Type: tail, Mode: "v"})
}

func (s signature) argType(i int) string {
	if i < len(s.Args) {
		return s.Args[i]
	}
	return "any"
}

// omitted are built-in functions the dialect leaves out: ones that exist
// for their side effect and return nothing a query can use, and one the
// shell's build adds that the library does not have.
var omitted = map[string]bool{
	// Loads a shared library and returns NULL.
	"load_extension": true,
	// Writes to the error log and returns NULL.
	"sqlite_log": true,
	// SQLITE_ENABLE_UNKNOWN_SQL_FUNCTION's stand-in for any function the
	// shell does not know, so that EXPLAIN works on queries that use one.
	"unknown": true,
	// FTS3's and FTS5's ways of passing pointers to their virtual tables,
	// not functions a query calls.
	"fts3_tokenizer": true,
	"fts5":           true,
	// A debugging aid of GEOPOLY's.
	"geopoly_debug": true,
}

// signatures covers every function of the default build and of each
// extension's, by the SQLite documentation of each — lang_corefunc,
// lang_mathfunc, lang_datefunc, lang_aggfunc, windowfunctions, json1,
// fts3, fts5, rtree and geopoly. A function marked Nullable can return NULL
// for arguments that are not: an aggregate over no rows, a lookup that
// finds nothing, an input that does not parse.
//
// SQLite's values carry their own types, so a function's result type is
// what it typically produces. abs of an integer is an integer, but abs
// returns real here, as sum does, because the analyzer wants one answer;
// the polymorphic "any" is for functions that hand back one of their
// arguments.
var signatures = map[string]signature{
	// Aggregates.
	"avg":          {Returns: "real", Nullable: true},
	"count":        {Returns: "integer"},
	"group_concat": {Args: []string{"any", "text"}, Returns: "text", Nullable: true},
	"max":          {Returns: "any", Nullable: true},
	"min":          {Returns: "any", Nullable: true},
	"string_agg":   {Args: []string{"any", "text"}, Returns: "text", Nullable: true},
	"sum":          {Returns: "real", Nullable: true},
	"total":        {Returns: "real"},

	// Percentiles, added by SQLITE_ENABLE_PERCENTILE.
	"median":          {Returns: "real", Nullable: true},
	"percentile":      {Args: []string{"any", "real"}, Returns: "real", Nullable: true},
	"percentile_cont": {Args: []string{"any", "real"}, Returns: "real", Nullable: true},
	"percentile_disc": {Args: []string{"any", "real"}, Returns: "any", Nullable: true},

	// Window functions.
	"cume_dist":    {Returns: "real"},
	"dense_rank":   {Returns: "integer"},
	"first_value":  {Returns: "any", Nullable: true},
	"lag":          {Args: []string{"any", "integer", "any"}, Returns: "any", Nullable: true},
	"last_value":   {Returns: "any", Nullable: true},
	"lead":         {Args: []string{"any", "integer", "any"}, Returns: "any", Nullable: true},
	"nth_value":    {Args: []string{"any", "integer"}, Returns: "any", Nullable: true},
	"ntile":        {Args: []string{"integer"}, Returns: "integer"},
	"percent_rank": {Returns: "real"},
	"rank":         {Returns: "integer"},
	"row_number":   {Returns: "integer"},

	// Core functions.
	"changes":                   {Returns: "integer"},
	"char":                      {Variadic: "integer", Returns: "text"},
	"coalesce":                  {Returns: "any", Nullable: true},
	"concat":                    {Returns: "text"},
	"concat_ws":                 {Args: []string{"text"}, Returns: "text"},
	"format":                    {Args: []string{"text"}, Returns: "text", Nullable: true},
	"glob":                      {Args: []string{"text", "text"}, Returns: "integer"},
	"hex":                       {Returns: "text"},
	"if":                        {Returns: "any", Nullable: true},
	"ifnull":                    {Returns: "any", Nullable: true},
	"iif":                       {Returns: "any", Nullable: true},
	"instr":                     {Args: []string{"text", "text"}, Returns: "integer", Nullable: true},
	"last_insert_rowid":         {Returns: "integer"},
	"length":                    {Returns: "integer", Nullable: true},
	"like":                      {Args: []string{"text", "text", "text"}, Returns: "integer"},
	"likelihood":                {Args: []string{"any", "real"}, Returns: "any", Nullable: true},
	"likely":                    {Returns: "any", Nullable: true},
	"lower":                     {Args: []string{"text"}, Returns: "text"},
	"ltrim":                     {Args: []string{"text", "text"}, Returns: "text"},
	"nullif":                    {Returns: "any", Nullable: true},
	"octet_length":              {Returns: "integer", Nullable: true},
	"printf":                    {Args: []string{"text"}, Returns: "text", Nullable: true},
	"quote":                     {Returns: "text"},
	"random":                    {Returns: "integer"},
	"randomblob":                {Args: []string{"integer"}, Returns: "blob"},
	"replace":                   {Args: []string{"text", "text", "text"}, Returns: "text"},
	"round":                     {Args: []string{"real", "real"}, Returns: "real"},
	"rtrim":                     {Args: []string{"text", "text"}, Returns: "text"},
	"sign":                      {Returns: "integer", Nullable: true},
	"sqlite_compileoption_get":  {Args: []string{"integer"}, Returns: "text", Nullable: true},
	"sqlite_compileoption_used": {Args: []string{"text"}, Returns: "integer"},
	"sqlite_source_id":          {Returns: "text"},
	"sqlite_version":            {Returns: "text"},
	"substr":                    {Args: []string{"any", "integer", "integer"}, Returns: "text"},
	"substring":                 {Args: []string{"any", "integer", "integer"}, Returns: "text"},
	"subtype":                   {Returns: "integer"},
	"total_changes":             {Returns: "integer"},
	"trim":                      {Args: []string{"text", "text"}, Returns: "text"},
	"typeof":                    {Returns: "text"},
	"unhex":                     {Args: []string{"text", "text"}, Returns: "blob", Nullable: true},
	"unicode":                   {Args: []string{"text"}, Returns: "integer"},
	"unistr":                    {Args: []string{"text"}, Returns: "text"},
	"unistr_quote":              {Args: []string{"text"}, Returns: "text"},
	"unlikely":                  {Returns: "any", Nullable: true},
	"upper":                     {Args: []string{"text"}, Returns: "text"},
	"zeroblob":                  {Args: []string{"integer"}, Returns: "blob"},

	// Math functions.
	"abs":     {Returns: "real"},
	"acos":    {Returns: "real"},
	"acosh":   {Returns: "real"},
	"asin":    {Returns: "real"},
	"asinh":   {Returns: "real"},
	"atan":    {Returns: "real"},
	"atan2":   {Returns: "real"},
	"atanh":   {Returns: "real"},
	"ceil":    {Returns: "integer"},
	"ceiling": {Returns: "integer"},
	"cos":     {Returns: "real"},
	"cosh":    {Returns: "real"},
	"degrees": {Returns: "real"},
	"exp":     {Returns: "real"},
	"floor":   {Returns: "integer"},
	"ln":      {Returns: "real"},
	"log":     {Returns: "real"},
	"log10":   {Returns: "real"},
	"log2":    {Returns: "real"},
	"mod":     {Returns: "real"},
	"pi":      {Returns: "real"},
	"pow":     {Returns: "real"},
	"power":   {Returns: "real"},
	"radians": {Returns: "real"},
	"sin":     {Returns: "real"},
	"sinh":    {Returns: "real"},
	"sqrt":    {Returns: "real"},
	"tan":     {Returns: "real"},
	"tanh":    {Returns: "real"},
	"trunc":   {Returns: "integer"},

	// Date and time functions, which return text in ISO-8601 form and NULL
	// for a time value they cannot parse.
	"current_date":      {Returns: "text"},
	"current_time":      {Returns: "text"},
	"current_timestamp": {Returns: "text"},
	"date":              {Returns: "text", Nullable: true},
	"datetime":          {Returns: "text", Nullable: true},
	"julianday":         {Returns: "real", Nullable: true},
	"strftime":          {Args: []string{"text"}, Returns: "text", Nullable: true},
	"time":              {Returns: "text", Nullable: true},
	"timediff":          {Returns: "text", Nullable: true},
	"unixepoch":         {Returns: "integer", Nullable: true},

	// JSON functions. A JSON argument is text or a JSONB blob, so it is
	// "any"; a path is text. The json_ forms return JSON text and the
	// jsonb_ forms JSONB blobs, and the extracting forms return whatever
	// the path leads to, or NULL when it leads nowhere.
	"->":                  {Args: []string{"any", "text"}, Returns: "text", Nullable: true},
	"->>":                 {Args: []string{"any", "text"}, Returns: "any", Nullable: true},
	"json":                {Returns: "text"},
	"json_array":          {Returns: "text"},
	"json_array_insert":   {Args: []string{"any"}, Returns: "text"},
	"json_array_length":   {Args: []string{"any", "text"}, Returns: "integer", Nullable: true},
	"json_error_position": {Returns: "integer"},
	"json_extract":        {Args: []string{"any", "text"}, Returns: "any", Nullable: true},
	"json_group_array":    {Returns: "text"},
	"json_group_object":   {Args: []string{"text", "any"}, Returns: "text"},
	"json_insert":         {Args: []string{"any"}, Returns: "text"},
	"json_object":         {Returns: "text"},
	"json_patch":          {Returns: "text"},
	"json_pretty":         {Args: []string{"any", "text"}, Returns: "text"},
	"json_quote":          {Returns: "text"},
	"json_remove":         {Args: []string{"any"}, Returns: "text"},
	"json_replace":        {Args: []string{"any"}, Returns: "text"},
	"json_set":            {Args: []string{"any"}, Returns: "text"},
	"json_type":           {Args: []string{"any", "text"}, Returns: "text", Nullable: true},
	"json_valid":          {Args: []string{"any", "integer"}, Returns: "integer"},
	"jsonb":               {Returns: "blob"},
	"jsonb_array":         {Returns: "blob"},
	"jsonb_array_insert":  {Args: []string{"any"}, Returns: "blob"},
	"jsonb_extract":       {Args: []string{"any", "text"}, Returns: "any", Nullable: true},
	"jsonb_group_array":   {Returns: "blob"},
	"jsonb_group_object":  {Args: []string{"text", "any"}, Returns: "blob"},
	"jsonb_insert":        {Args: []string{"any"}, Returns: "blob"},
	"jsonb_object":        {Returns: "blob"},
	"jsonb_patch":         {Returns: "blob"},
	"jsonb_remove":        {Args: []string{"any"}, Returns: "blob"},
	"jsonb_replace":       {Args: []string{"any"}, Returns: "blob"},
	"jsonb_set":           {Args: []string{"any"}, Returns: "blob"},

	// Added by SQLITE_SOUNDEX.
	"soundex": {Args: []string{"text"}, Returns: "text"},

	// Added by SQLITE_ENABLE_OFFSET_SQL_FUNC.
	"sqlite_offset": {Returns: "integer", Nullable: true},

	// FTS3's auxiliary functions, whose first argument names the table.
	"matchinfo": {Args: []string{"any", "text"}, Returns: "blob"},
	"offsets":   {Args: []string{"any"}, Returns: "text"},
	"optimize":  {Args: []string{"any"}, Returns: "text"},

	// FTS5's auxiliary and locale functions. snippet is FTS3's too, with the
	// same result.
	"bm25":            {Args: []string{"text"}, Variadic: "real", Returns: "real"},
	"fts5_get_locale": {Args: []string{"any", "any"}, Returns: "text", Nullable: true},
	"fts5_insttoken":  {Args: []string{"text"}, Returns: "text"},
	"fts5_locale":     {Args: []string{"text", "text"}, Returns: "text"},
	"fts5_source_id":  {Returns: "text"},
	"highlight":       {Args: []string{"text", "integer", "text", "text"}, Returns: "text"},
	"snippet":         {Args: []string{"text", "integer", "text", "text", "text", "integer"}, Returns: "text"},

	// R*Tree's functions for inspecting an index's nodes.
	"rtreecheck": {Args: []string{"text", "text"}, Returns: "text"},
	"rtreedepth": {Args: []string{"blob"}, Returns: "integer"},
	"rtreenode":  {Args: []string{"integer", "blob"}, Returns: "text"},

	// GEOPOLY's functions. A polygon argument is GeoJSON text or the binary
	// form, so it is "any"; the functions that return a polygon return the
	// binary form.
	"geopoly_area":           {Returns: "real", Nullable: true},
	"geopoly_bbox":           {Returns: "blob", Nullable: true},
	"geopoly_blob":           {Returns: "blob", Nullable: true},
	"geopoly_ccw":            {Returns: "blob", Nullable: true},
	"geopoly_contains_point": {Args: []string{"any", "real", "real"}, Returns: "integer", Nullable: true},
	"geopoly_group_bbox":     {Returns: "blob", Nullable: true},
	"geopoly_json":           {Returns: "text", Nullable: true},
	"geopoly_overlap":        {Returns: "integer", Nullable: true},
	"geopoly_regular":        {Args: []string{"real", "real", "real", "integer"}, Returns: "blob"},
	"geopoly_svg":            {Args: []string{"any", "text"}, Returns: "text", Nullable: true},
	"geopoly_within":         {Returns: "integer", Nullable: true},
	"geopoly_xform":          {Args: []string{"any", "real", "real", "real", "real", "real", "real"}, Returns: "blob", Nullable: true},
}
