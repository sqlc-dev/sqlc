# Tables

```sql
CREATE TABLE authors (
  id   SERIAL PRIMARY KEY,
  name text   NOT NULL
);
```

```go
package db

// Struct names use the singular form of table names
type Author struct {
	ID   int
	Name string
}
```
