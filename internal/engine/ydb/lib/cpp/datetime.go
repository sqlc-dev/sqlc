package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func DateTimeFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, dateTimeMakeFuncs()...)
	funcs = append(funcs, dateTimeGetFuncs()...)
	funcs = append(funcs, dateTimeUpdateFuncs()...)
	funcs = append(funcs, dateTimeFromFuncs()...)
	funcs = append(funcs, dateTimeToFuncs()...)
	funcs = append(funcs, dateTimeIntervalFuncs()...)
	funcs = append(funcs, dateTimeStartEndFuncs()...)
	funcs = append(funcs, dateTimeFormatFuncs()...)
	funcs = append(funcs, dateTimeParseFuncs()...)

	return funcs
}

func dateTimeMakeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_makedate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "datetime_makedate32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Date32"},
		},
		{
			Name: "datetime_maketzdate32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDate32"},
		},
		{
			Name: "datetime_makedatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Datetime"},
		},
		{
			Name: "datetime_maketzdatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime"},
		},
		{
			Name: "datetime_makedatetime64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Datetime64"},
		},
		{
			Name: "datetime_maketzdatetime64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime64"},
		},
		{
			Name: "datetime_maketimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "datetime_maketztimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzTimestamp"},
		},
		{
			Name: "datetime_maketimestamp64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "datetime_maketztimestamp64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzTimestamp64"},
		},
	}
}

func dateTimeGetFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_getyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "datetime_getyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Int32"},
		},
		{
			Name: "datetime_getdayofyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "datetime_getmonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getmonthname",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "datetime_getweekofyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getweekofyeariso8601",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getdayofmonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getdayofweek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getdayofweekname",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "datetime_gethour",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getminute",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getsecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "datetime_getmillisecondofsecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "datetime_getmicrosecondofsecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "datetime_gettimezoneid",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "datetime_gettimezonename",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
	}
}

func dateTimeUpdateFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func dateTimeFromFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_fromseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "datetime_fromseconds64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "datetime_frommilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "datetime_frommilliseconds64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "datetime_frommicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "datetime_frommicroseconds64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
	}
}

func dateTimeToFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_toseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tomilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tomicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func dateTimeIntervalFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_todays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tohours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tominutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_toseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tomilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_tomicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_intervalfromdays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64fromdays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "datetime_intervalfromhours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64fromhours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "datetime_intervalfromminutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64fromminutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "datetime_intervalfromseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64fromseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "datetime_intervalfrommilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64frommilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "datetime_intervalfrommicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "datetime_interval64frommicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
	}
}

func dateTimeStartEndFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_startofyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endofyear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_startofquarter",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endofquarter",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_startofmonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endofmonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_startofweek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endofweek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_startofday",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endofday",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_startof",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_endof",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func dateTimeFormatFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_format",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_format_call",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "string"},
		},
	}
}

func dateTimeParseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "datetime_parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_parse64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "datetime_parserfc822",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_parseiso8601",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_parsehttp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "datetime_parsex509",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}
