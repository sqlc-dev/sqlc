package sqltest

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"
)

var containerNameRe = regexp.MustCompile(`[[:^alpha:]]`)

func containerName(t *testing.T, driver string) string {
	testName := strings.ToLower(containerNameRe.ReplaceAllString(t.Name(), "_"))
	suffix, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		t.Fatalf("failed to generate random suffix: %v", err)
	}

	return fmt.Sprintf("sqlc-%s-%s-%s", driver, testName, suffix)
}
