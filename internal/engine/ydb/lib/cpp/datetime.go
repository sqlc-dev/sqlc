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
			Name: "DateTime::MakeDate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "DateTime::MakeDate32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Date32"},
		},
		{
			Name: "DateTime::MakeTzDate32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDate32"},
		},
		{
			Name: "DateTime::MakeDatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Datetime"},
		},
		{
			Name: "DateTime::MakeTzDatetime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime"},
		},
		{
			Name: "DateTime::MakeDatetime64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Datetime64"},
		},
		{
			Name: "DateTime::MakeTzDatetime64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzDatetime64"},
		},
		{
			Name: "DateTime::MakeTimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "DateTime::MakeTzTimestamp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "TzTimestamp"},
		},
		{
			Name: "DateTime::MakeTimestamp64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "DateTime::MakeTzTimestamp64",
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
			Name: "DateTime::GetYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "DateTime::GetYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Int32"},
		},
		{
			Name: "DateTime::GetDayOfYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "DateTime::GetMonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetMonthName",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "DateTime::GetWeekOfYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetWeekOfYearIso8601",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetDayOfMonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetDayOfWeek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetDayOfWeekName",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "DateTime::GetHour",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetMinute",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetSecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint8"},
		},
		{
			Name: "DateTime::GetMillisecondOfSecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "DateTime::GetMicrosecondOfSecond",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint32"},
		},
		{
			Name: "DateTime::GetTimezoneId",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint16"},
		},
		{
			Name: "DateTime::GetTimezoneName",
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
			Name: "DateTime::Update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::Update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::Update",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::Update",
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
			Name: "DateTime::Update",
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
			Name: "DateTime::Update",
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
			Name: "DateTime::Update",
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
			Name: "DateTime::Update",
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
			Name: "DateTime::FromSeconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "DateTime::FromSeconds64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "DateTime::FromMilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "DateTime::FromMilliseconds64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp64"},
		},
		{
			Name: "DateTime::FromMicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Timestamp"},
		},
		{
			Name: "DateTime::FromMicroseconds64",
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
			Name: "DateTime::ToSeconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToMilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToMicroseconds",
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
			Name: "DateTime::ToDays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToHours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToMinutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToSeconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToMilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ToMicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::IntervalFromDays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromDays",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "DateTime::IntervalFromHours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromHours",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "DateTime::IntervalFromMinutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromMinutes",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "DateTime::IntervalFromSeconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromSeconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "DateTime::IntervalFromMilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromMilliseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval64"},
		},
		{
			Name: "DateTime::IntervalFromMicroseconds",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType: &ast.TypeName{Name: "Interval"},
		},
		{
			Name: "DateTime::Interval64FromMicroseconds",
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
			Name: "DateTime::StartOfYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOfYear",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::StartOfQuarter",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOfQuarter",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::StartOfMonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOfMonth",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::StartOfWeek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOfWeek",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::StartOfDay",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOfDay",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::StartOf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::EndOf",
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
			Name: "DateTime::Format",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func dateTimeParseFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "DateTime::Parse",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::Parse64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "DateTime::ParseRfc822",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::ParseIso8601",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::ParseHttp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "DateTime::ParseX509",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}
