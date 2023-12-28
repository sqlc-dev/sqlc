package quickdb

import (
	"fmt"
	"net/url"
)

// The database URI returned by the QuickDB service isn't understood by the
// go-mysql-driver
func MySQLReformatURI(original string) (string, error) {
	u, err := url.Parse(original)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@tcp(%s)%s?multiStatements=true&parseTime=true&tls=true", u.User, u.Host, u.Path), nil
}
