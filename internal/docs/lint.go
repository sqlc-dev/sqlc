// Package docs validates the documentation content contract for docs/.
//
// The docs/ directory is plain CommonMark + GitHub-flavored Markdown with no
// rendering toolchain in this repository; a separate repository consumes it
// and builds the documentation site. This package enforces the contract that
// consumer relies on:
//
//   - every page parses as Markdown and starts with exactly one level-1
//     heading (the page title)
//   - every relative link and image resolves to a file inside docs/, and
//     anchor fragments resolve to a heading in the target page
//   - every page appears exactly once in toc.yaml, either in a section or in
//     the unlisted set, and every toc.yaml entry names a real page
//   - no raw HTML or JSX except HTML comments
//   - no MyST/Sphinx directives left over from the old toolchain
package docs

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

// A Problem is a single contract violation, located by file (relative to the
// docs root) and 1-based line number.
type Problem struct {
	File    string
	Line    int
	Message string
}

func (p Problem) String() string {
	return fmt.Sprintf("%s:%d: %s", p.File, p.Line, p.Message)
}

type tocFile struct {
	Index    string `yaml:"index"`
	Sections []struct {
		Title string   `yaml:"title"`
		Pages []string `yaml:"pages"`
	} `yaml:"sections"`
	Unlisted []string `yaml:"unlisted"`
}

type page struct {
	src     []byte
	doc     ast.Node
	anchors map[string]bool
}

// Lint validates every Markdown file under root against the content
// contract and cross-checks the set of files against root/toc.yaml.
func Lint(root string) ([]Problem, error) {
	var problems []Problem

	files, err := markdownFiles(root)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	// First pass: parse everything and collect heading anchors, so link
	// fragments can be validated against any page.
	pages := make(map[string]*page, len(files))
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, err
		}
		doc := md.Parser().Parse(text.NewReader(src))
		pages[rel] = &page{src: src, doc: doc, anchors: headingAnchors(doc, src)}
	}

	for _, rel := range files {
		problems = append(problems, lintPage(rel, pages)...)
	}

	problems = append(problems, lintTOC(root, files)...)

	sort.Slice(problems, func(i, j int) bool {
		if problems[i].File != problems[j].File {
			return problems[i].File < problems[j].File
		}
		return problems[i].Line < problems[j].Line
	})
	return problems, nil
}

// markdownFiles returns the slash-separated paths of all .md files under
// root, relative to root.
func markdownFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func lintPage(rel string, pages map[string]*page) []Problem {
	var problems []Problem
	p := pages[rel]
	report := func(n ast.Node, format string, args ...any) {
		problems = append(problems, Problem{
			File:    rel,
			Line:    nodeLine(n, p.src),
			Message: fmt.Sprintf(format, args...),
		})
	}

	var h1s int
	firstBlock := true
	ast.Walk(p.doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if n.Type() == ast.TypeBlock && n.Parent() == p.doc {
			if firstBlock {
				firstBlock = false
				if h, ok := n.(*ast.Heading); !ok || h.Level != 1 {
					report(n, "page must start with a level-1 heading (its title)")
				}
			}
		}
		switch n := n.(type) {
		case *ast.Heading:
			if n.Level == 1 {
				h1s++
				if h1s == 2 {
					report(n, "page has more than one level-1 heading")
				}
			}
		case *ast.Link:
			problems = append(problems, checkLink(rel, string(n.Destination), false, n, pages)...)
		case *ast.Image:
			problems = append(problems, checkLink(rel, string(n.Destination), true, n, pages)...)
		case *ast.HTMLBlock, *ast.RawHTML:
			if raw := rawHTMLText(n, p.src); !isHTMLComment(raw) {
				report(n, "raw HTML is not allowed: %s", firstLine(raw))
			}
		case *ast.Paragraph:
			if t := nodeText(n, p.src); strings.HasPrefix(t, ":::") {
				report(n, "MyST directive syntax is not allowed; use GitHub alerts (> [!NOTE]) instead")
			}
		case *ast.FencedCodeBlock:
			if info := fenceInfo(n, p.src); strings.HasPrefix(info, "{") {
				report(n, "MyST fenced directive %q is not allowed; use GitHub alerts (> [!NOTE]) instead", info)
			}
		}
		return ast.WalkContinue, nil
	})
	return problems
}

func checkLink(rel, dest string, isImage bool, n ast.Node, pages map[string]*page) []Problem {
	p := pages[rel]
	problem := func(format string, args ...any) []Problem {
		return []Problem{{
			File:    rel,
			Line:    nodeLine(n, p.src),
			Message: fmt.Sprintf(format, args...),
		}}
	}

	switch {
	case dest == "":
		return problem("empty link destination")
	case strings.HasPrefix(dest, "http://"), strings.HasPrefix(dest, "https://"), strings.HasPrefix(dest, "mailto:"):
		return nil
	case strings.HasPrefix(dest, "#"):
		if !p.anchors[strings.TrimPrefix(dest, "#")] {
			return problem("anchor %s not found in this page", dest)
		}
		return nil
	case strings.HasPrefix(dest, "/"):
		return problem("absolute link %q; use a path relative to this file", dest)
	case strings.Contains(dest, "://"):
		return problem("unsupported link scheme in %q", dest)
	}

	pathPart, frag, _ := strings.Cut(dest, "#")
	target := path.Join(path.Dir(rel), pathPart)
	if target == ".." || strings.HasPrefix(target, "../") {
		return problem("link %q points outside docs/", dest)
	}

	tp, ok := pages[target]
	if !ok {
		if isImage || !strings.HasSuffix(target, ".md") {
			// Non-page assets (images, other files) just need to exist on
			// disk; they were not parsed in the first pass.
			return problem("link target %q does not exist", target)
		}
		return problem("link target %q does not exist", target)
	}
	if frag != "" && !tp.anchors[frag] {
		return problem("anchor #%s not found in %s", frag, target)
	}
	return nil
}

func lintTOC(root string, files []string) []Problem {
	const tocName = "toc.yaml"
	problem := func(format string, args ...any) Problem {
		return Problem{File: tocName, Line: 1, Message: fmt.Sprintf(format, args...)}
	}

	data, err := os.ReadFile(filepath.Join(root, tocName))
	if err != nil {
		return []Problem{problem("%v", err)}
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var toc tocFile
	if err := dec.Decode(&toc); err != nil {
		return []Problem{problem("%v", err)}
	}

	var problems []Problem
	listed := make(map[string]int)
	add := func(entry string) {
		listed[entry]++
	}
	add(toc.Index)
	for _, s := range toc.Sections {
		for _, pg := range s.Pages {
			add(pg)
		}
	}
	for _, pg := range toc.Unlisted {
		add(pg)
	}

	onDisk := make(map[string]bool, len(files))
	for _, f := range files {
		onDisk[f] = true
	}
	for entry, count := range listed {
		if count > 1 {
			problems = append(problems, problem("%s is listed %d times", entry, count))
		}
		if !onDisk[entry] {
			problems = append(problems, problem("%s is listed but does not exist", entry))
		}
	}
	for _, f := range files {
		if listed[f] == 0 {
			problems = append(problems, problem("%s is not listed; add it to a section or to unlisted", f))
		}
	}
	return problems
}

// headingAnchors returns the set of GitHub-style anchor slugs for the
// headings in doc, applying GitHub's -1, -2 suffixes for duplicates.
func headingAnchors(doc ast.Node, src []byte) map[string]bool {
	anchors := make(map[string]bool)
	seen := make(map[string]int)
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if h, ok := n.(*ast.Heading); ok && entering {
			slug := slugify(nodeText(h, src))
			if c := seen[slug]; c > 0 {
				anchors[fmt.Sprintf("%s-%d", slug, c)] = true
			} else {
				anchors[slug] = true
			}
			seen[slug]++
		}
		return ast.WalkContinue, nil
	})
	return anchors
}

// slugify converts heading text to a GitHub-style anchor: lowercase, spaces
// become hyphens, and everything except letters, digits, hyphens and
// underscores is dropped.
func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}

func nodeText(n ast.Node, src []byte) string {
	var b strings.Builder
	ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch c := c.(type) {
			case *ast.Text:
				b.Write(c.Segment.Value(src))
			case *ast.String:
				b.Write(c.Value)
			}
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

func rawHTMLText(n ast.Node, src []byte) string {
	var b strings.Builder
	switch n := n.(type) {
	case *ast.HTMLBlock:
		for i := 0; i < n.Lines().Len(); i++ {
			seg := n.Lines().At(i)
			b.Write(seg.Value(src))
		}
		if n.HasClosure() {
			b.Write(n.ClosureLine.Value(src))
		}
	case *ast.RawHTML:
		for i := 0; i < n.Segments.Len(); i++ {
			seg := n.Segments.At(i)
			b.Write(seg.Value(src))
		}
	}
	return strings.TrimSpace(b.String())
}

func isHTMLComment(raw string) bool {
	return strings.HasPrefix(raw, "<!--") && strings.HasSuffix(raw, "-->")
}

func fenceInfo(n *ast.FencedCodeBlock, src []byte) string {
	if n.Info == nil {
		return ""
	}
	return strings.TrimSpace(string(n.Info.Segment.Value(src)))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// nodeLine returns the 1-based line number of n, using the nearest
// enclosing block that carries source positions.
func nodeLine(n ast.Node, src []byte) int {
	if rh, ok := n.(*ast.RawHTML); ok && rh.Segments.Len() > 0 {
		return 1 + bytes.Count(src[:rh.Segments.At(0).Start], []byte("\n"))
	}
	for cur := n; cur != nil; cur = cur.Parent() {
		if cur.Type() == ast.TypeBlock && cur.Lines().Len() > 0 {
			return 1 + bytes.Count(src[:cur.Lines().At(0).Start], []byte("\n"))
		}
	}
	return 1
}
