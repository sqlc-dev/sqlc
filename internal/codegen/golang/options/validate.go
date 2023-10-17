package options

import "fmt"

func Validate(opts *Options) error {
	if opts.EmitMethodsWithDbArgument && opts.EmitPreparedQueries {
		return fmt.Errorf("invalid options: emit_methods_with_db_argument and emit_prepared_queries options are mutually exclusive")
	}

	return nil
}
