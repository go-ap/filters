package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/google/go-cmp/cmp"
	"github.com/leporo/sqlf"
)

func Test_SQLWhere(t *testing.T) {
	type args struct {
		s *Stmt
		f []Check
	}
	tests := []struct {
		name     string
		args     args
		gotQuery string
		gotArgs  []any
	}{
		{
			name: "empty",
			args: args{},
		},
		{
			name: "one type",
			args: args{
				s: sqlf.New(""),
				f: []Check{HasType("t1")},
			},
			gotQuery: " WHERE type = ?",
			gotArgs:  []any{vocab.ActivityVocabularyType("t1")},
		},
		{
			name: "multiple types",
			args: args{
				s: sqlf.New(""),
				f: []Check{HasType("t1", "t2", "t3")},
			},
			gotQuery: " WHERE type IN (?,?,?)",
			gotArgs:  []any{vocab.ActivityVocabularyType("t1"), vocab.ActivityVocabularyType("t2"), vocab.ActivityVocabularyType("t3")},
		},
		{
			name: "one type with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{HasType("t1", vocab.NilType)},
			},
			gotQuery: " WHERE (type = ? OR type IS NULL)",
			gotArgs:  []any{vocab.ActivityVocabularyType("t1")},
		},
		{
			name: "multiple types with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{HasType("t1", "t2", vocab.NilType)},
			},
			gotQuery: " WHERE (type IN (?,?) OR type IS NULL)",
			gotArgs:  []any{vocab.ActivityVocabularyType("t1"), vocab.ActivityVocabularyType("t2")},
		},
		{
			name: "one ID",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameID("http://example.com")},
			},
			gotQuery: " WHERE iri = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "multiple IDs",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameID("http://example.com"), SameID("http://social.example.com")},
			},
			gotQuery: " WHERE iri IN (?,?)",
			gotArgs:  []any{vocab.IRI("http://example.com"), vocab.IRI("http://social.example.com")},
		},
		{
			name: "one ID with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameID("http://example.com"), NilID},
			},
			gotQuery: " WHERE (iri = ? OR iri IS NULL)",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "multiple IDs with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameID("http://example.com"), SameID("http://social.example.com"), NilID},
			},
			gotQuery: " WHERE (iri IN (?,?) OR iri IS NULL)",
			gotArgs:  []any{vocab.IRI("http://example.com"), vocab.IRI("http://social.example.com")},
		},
		//
		{
			name: "one iri",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameIRI("http://example.com")},
			},
			gotQuery: " WHERE iri = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "iri like",
			args: args{
				s: sqlf.New(""),
				f: []Check{IRILike("http://example.com")},
			},
			gotQuery: " WHERE iri LIKE ?",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "multiple iris",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameIRI("http://example.com"), SameIRI("http://social.example.com")},
			},
			gotQuery: " WHERE iri IN (?,?)",
			gotArgs:  []any{vocab.IRI("http://example.com"), vocab.IRI("http://social.example.com")},
		},
		{
			name: "one iri with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameIRI("http://example.com"), NilIRI},
			},
			gotQuery: " WHERE (iri = ? OR iri IS NULL)",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "iri like with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{IRILike("http://example.com"), NilIRI},
			},
			gotQuery: " WHERE (iri LIKE ? OR iri IS NULL)",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "multiple iri likes with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{IRILike("http://example.com"), IRILike("http://social.example.com"), NilIRI},
			},
			gotQuery: " WHERE (iri LIKE ? OR iri LIKE ? OR iri IS NULL)",
			gotArgs:  []any{"%http://example.com%", "%http://social.example.com%"},
		},
		{
			name: "multiple iris with nil",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameIRI("http://example.com"), SameIRI("http://social.example.com"), NilIRI},
			},
			gotQuery: " WHERE (iri IN (?,?) OR iri IS NULL)",
			gotArgs:  []any{vocab.IRI("http://example.com"), vocab.IRI("http://social.example.com")},
		},
		{
			name: "name empty",
			args: args{
				s: sqlf.New(""),
				f: []Check{NameEmpty},
			},
			gotQuery: " WHERE name IS NULL",
		},
		{
			name: "preferredUsername empty",
			args: args{
				s: sqlf.New(""),
				f: []Check{PreferredUsernameEmpty},
			},
			gotQuery: " WHERE preferred_username IS NULL",
		},
		{
			name: "summary empty",
			args: args{
				s: sqlf.New(""),
				f: []Check{SummaryEmpty},
			},
			gotQuery: " WHERE summary IS NULL",
		},
		{
			name: "content empty",
			args: args{
				s: sqlf.New(""),
				f: []Check{ContentEmpty},
			},
			gotQuery: " WHERE content IS NULL",
		},
		//
		{
			name: "name equals",
			args: args{
				s: sqlf.New(""),
				f: []Check{NameIs("test")},
			},
			gotQuery: " WHERE name = ?",
			gotArgs:  []any{"test"},
		},
		{
			name: "preferredUsername equals",
			args: args{
				s: sqlf.New(""),
				f: []Check{PreferredUsernameIs("test")},
			},
			gotQuery: " WHERE preferred_username = ?",
			gotArgs:  []any{"test"},
		},
		{
			name: "summary equals",
			args: args{
				s: sqlf.New(""),
				f: []Check{SummaryIs("test")},
			},
			gotQuery: " WHERE summary = ?",
			gotArgs:  []any{"test"},
		},
		{
			name: "content equals",
			args: args{
				s: sqlf.New(""),
				f: []Check{ContentIs("test")},
			},
			gotQuery: " WHERE content = ?",
			gotArgs:  []any{"test"},
		},
		//
		{
			name: "name like",
			args: args{
				s: sqlf.New(""),
				f: []Check{NameLike("test")},
			},
			gotQuery: " WHERE name LIKE ?",
			gotArgs:  []any{"%test%"},
		},
		{
			name: "preferredUsername like",
			args: args{
				s: sqlf.New(""),
				f: []Check{PreferredUsernameLike("test")},
			},
			gotQuery: " WHERE preferred_username LIKE ?",
			gotArgs:  []any{"%test%"},
		},
		{
			name: "summary like",
			args: args{
				s: sqlf.New(""),
				f: []Check{SummaryLike("test")},
			},
			gotQuery: " WHERE summary LIKE ?",
			gotArgs:  []any{"%test%"},
		},
		{
			name: "content like",
			args: args{
				s: sqlf.New(""),
				f: []Check{ContentLike("test")},
			},
			gotQuery: " WHERE content LIKE ?",
			gotArgs:  []any{"%test%"},
		},
		//
		{
			name: "multiple rules",
			args: args{
				s: sqlf.New(""),
				f: []Check{ContentLike("test"), NameEmpty, SummaryIs("test1")},
			},
			gotQuery: " WHERE content LIKE ? AND name IS NULL AND summary = ?",
			gotArgs:  []any{"%test%", "test1"},
		},
		{
			name: "inReplyTo nil for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{NilInReplyTo},
			},
			gotQuery: " WHERE json_extract(raw, '$.inReplyTo') IS NULL",
		},
		{
			name: "inReplyTo nil for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{NilInReplyTo},
			},
			gotQuery: " WHERE raw->>'inReplyTo' IS NULL",
		},
		{
			name: "inReplyTo equals for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameInReplyTo("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.inReplyTo') = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "inReplyTo equals for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{SameInReplyTo("http://example.com")},
			},
			gotQuery: " WHERE raw->>'inReplyTo' = $1",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "inReplyTo like for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{InReplyToLike("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.inReplyTo') LIKE ?",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "inReplyTo like for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{InReplyToLike("http://example.com")},
			},
			gotQuery: " WHERE raw->>'inReplyTo' LIKE $1",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "attributedTo nil for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{NilAttributedTo},
			},
			gotQuery: " WHERE json_extract(raw, '$.attributedTo') IS NULL",
		},
		{
			name: "attributedTo nil for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{NilAttributedTo},
			},
			gotQuery: " WHERE raw->>'attributedTo' IS NULL",
		},
		{
			name: "attributedTo equals for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameAttributedTo("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.attributedTo') = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "attributedTo equals for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{SameAttributedTo("http://example.com")},
			},
			gotQuery: " WHERE raw->>'attributedTo' = $1",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "attributedTo like for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{AttributedToLike("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.attributedTo') LIKE ?",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "attributedTo like for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{AttributedToLike("http://example.com")},
			},
			gotQuery: " WHERE raw->>'attributedTo' LIKE $1",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "context nil for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{NilContext},
			},
			gotQuery: " WHERE json_extract(raw, '$.context') IS NULL",
		},
		{
			name: "context nil for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{NilContext},
			},
			gotQuery: " WHERE raw->>'context' IS NULL",
		},
		{
			name: "context equals for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameContext("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.context') = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "context equals for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{SameContext("http://example.com")},
			},
			gotQuery: " WHERE raw->>'context' = $1",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "context like for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{ContextLike("http://example.com")},
			},
			gotQuery: " WHERE json_extract(raw, '$.context') LIKE ?",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "context like for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{ContextLike("http://example.com")},
			},
			gotQuery: " WHERE raw->>'context' LIKE $1",
			gotArgs:  []any{"%http://example.com%"},
		},
		{
			name: "URL nil for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{NilURL},
			},
			gotQuery: " WHERE url IS NULL",
		},
		{
			name: "URL nil for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{NilURL},
			},
			gotQuery: " WHERE url IS NULL",
		},
		{
			name: "URL equals for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{SameURL("http://example.com")},
			},
			gotQuery: " WHERE url = ?",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "URL equals for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{SameURL("http://example.com")},
			},
			gotQuery: " WHERE url = $1",
			gotArgs:  []any{vocab.IRI("http://example.com")},
		},
		{
			name: "URL like for sqlite",
			args: args{
				s: sqlf.New(""),
				f: []Check{URLLike("http://example.com")},
			},
			gotQuery: " WHERE url LIKE ?",
			gotArgs:  []any{vocab.IRI("%http://example.com%")},
		},
		{
			name: "URL like for pgsql",
			args: args{
				s: sqlf.PostgreSQL.New(""),
				f: []Check{URLLike("http://example.com")},
			},
			gotQuery: " WHERE url LIKE $1",
			gotArgs:  []any{vocab.IRI("%http://example.com%")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = SQLWhere(tt.args.s, tt.args.f...)

			var gotQuery string
			var gotArgs []any
			if tt.args.s != nil {
				gotQuery = tt.args.s.String()
				gotArgs = tt.args.s.Args()
			}
			if gotQuery != tt.gotQuery {
				t.Errorf("SQLWhere() query %s does not match expected: %s", gotQuery, tt.gotQuery)
			}

			if !cmp.Equal(gotArgs, tt.gotArgs) {
				t.Errorf("SQLWhere() query args are different: %s", cmp.Diff(tt.gotArgs, gotArgs))
			}
		})
	}
}

func TestSQLLimit(t *testing.T) {
	type args struct {
		st *Stmt
		f  []Check
	}
	tests := []struct {
		name     string
		args     args
		gotQuery string
		gotArgs  []any
	}{
		{
			name: "empty",
			args: args{},
		},
		{
			name: "no limit",
			args: args{
				st: sqlf.New(""),
				f:  []Check{HasType("t1"), SameInReplyTo("http://example.com")},
			},
			gotQuery: "",
			gotArgs:  []any{},
		},
		{
			name: "limit 1",
			args: args{
				st: sqlf.New(""),
				f:  []Check{WithMaxCount(1)},
			},
			gotQuery: "LIMIT ?",
			gotArgs:  []any{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SQLLimit(tt.args.st, tt.args.f...)
		})
	}
}
