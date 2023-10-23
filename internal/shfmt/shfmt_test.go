package shfmt

import "testing"

func TestReplace(t *testing.T) {
	s := "POSTGRES_SQL://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/AUTHORS"
	r := NewReplacer([]string{
		"PG_USER=user",
		"PG_PASSWORD=password",
		"PG_HOST=host",
		"PG_PORT=port",
	})
	e := "POSTGRES_SQL://user:password@host:port/AUTHORS"
	if v := r.Replace(s); v != e {
		t.Errorf("%s != %s", v, e)
	}
}
