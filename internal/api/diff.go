package api

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"sort"

	"github.com/cubicdaiya/gonp"
)

func writeFiles(ctx context.Context, files map[string]string, stderr io.Writer) error {
	defer trace.StartRegion(ctx, "writefiles").End()
	for filename, source := range files {
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			fmt.Fprintf(stderr, "%s: %s\n", filename, err)
			return err
		}
		if err := os.WriteFile(filename, []byte(source), 0644); err != nil {
			fmt.Fprintf(stderr, "%s: %s\n", filename, err)
			return err
		}
	}
	return nil
}

func diffFiles(ctx context.Context, baseDir string, files map[string]string, stderr io.Writer) error {
	defer trace.StartRegion(ctx, "checkfiles").End()
	var errored bool

	if baseDir == "" {
		baseDir, _ = os.Getwd()
	}

	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, filename := range keys {
		source := files[filename]
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			errored = true
			continue
		}
		existing, err := os.ReadFile(filename)
		if err != nil {
			errored = true
			fmt.Fprintf(stderr, "%s: %s\n", filename, err)
			continue
		}
		d := gonp.New(getLines(existing), getLines([]byte(source)))
		d.Compose()
		uniHunks := filterHunks(d.UnifiedHunks())

		if len(uniHunks) > 0 {
			errored = true
			label := filename
			if baseDir != "" {
				if rel, err := filepath.Rel(baseDir, filename); err == nil {
					label = "/" + rel
				}
			}
			fmt.Fprintf(stderr, "--- a%s\n", label)
			fmt.Fprintf(stderr, "+++ b%s\n", label)
			d.FprintUniHunks(stderr, uniHunks)
		}
	}
	if errored {
		return errors.New("diff found")
	}
	return nil
}

func getLines(f []byte) []string {
	fp := bytes.NewReader(f)
	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func filterHunks[T gonp.Elem](uniHunks []gonp.UniHunk[T]) []gonp.UniHunk[T] {
	var out []gonp.UniHunk[T]
	for i, uniHunk := range uniHunks {
		var changed bool
		for _, e := range uniHunk.GetChanges() {
			switch e.GetType() {
			case gonp.SesDelete:
				changed = true
			case gonp.SesAdd:
				changed = true
			}
		}
		if changed {
			out = append(out, uniHunks[i])
		}
	}
	return out
}
