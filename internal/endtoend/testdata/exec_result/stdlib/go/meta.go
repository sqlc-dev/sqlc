package querytest

var TableName = struct {
	Bar string
}{
	Bar: "bar",
}

var BarColumns = struct {
	ID string
}{
	ID: "id",
}

var BarSelectColumns = struct {
	ID string
}{
	ID: "bar.id",
}
