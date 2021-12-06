package querytest

var TableName = struct {
	Foo string
}{
	Foo: "foo",
}

var FooColumns = struct {
	Bar  string
	Bars string
}{
	Bar:  "bar",
	Bars: "bars",
}

var FooSelectColumns = struct {
	Bar  string
	Bars string
}{
	Bar:  "foo.bar",
	Bars: "foo.bars",
}
