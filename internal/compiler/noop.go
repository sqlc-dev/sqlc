// +build !exp

package compiler

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
)

func Run(conf config.SQL, combo config.CombinedSettings) (*Result, error) {
	return nil, fmt.Errorf("unimplemented")
}
