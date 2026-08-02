package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func TestApplyQueryComments(t *testing.T) {
	req := &plugin.GenerateRequest{
		Queries: []*plugin.Query{
			{
				Name:     "GetAuthor",
				Cmd:      ":one",
				Filename: "query.sql",
				Text:     "SELECT * FROM authors WHERE id = $1",
			},
		},
	}

	applyQueryComments(req, config.QueryComments{
		Enabled: true,
		Tags:    []string{"name", "cmd", "filename"},
	})

	want := "/*sqlc_name='GetAuthor',sqlc_cmd='%3Aone',sqlc_filename='query.sql'*/ SELECT * FROM authors WHERE id = $1"
	if diff := cmp.Diff(want, req.Queries[0].Text); diff != "" {
		t.Errorf("query text differed (-want +got):\n%s", diff)
	}
}

func TestApplyQueryCommentsMarginalia(t *testing.T) {
	req := &plugin.GenerateRequest{
		Queries: []*plugin.Query{
			{
				Name: "GetAuthor",
				Text: "SELECT * FROM authors WHERE id = $1",
			},
		},
	}

	applyQueryComments(req, config.QueryComments{
		Enabled: true,
		Format:  "marginalia",
	})

	want := "/*sqlc_name:GetAuthor*/ SELECT * FROM authors WHERE id = $1"
	if diff := cmp.Diff(want, req.Queries[0].Text); diff != "" {
		t.Errorf("query text differed (-want +got):\n%s", diff)
	}
}

func TestEscapeQueryCommentValue(t *testing.T) {
	got := escapeQueryCommentValue("a'b,c:d\n*/")
	want := "a%27b%2Cc%3Ad * /"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("escaped value differed (-want +got):\n%s", diff)
	}
}
