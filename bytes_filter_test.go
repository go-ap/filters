package filters

import (
	"bytes"
	"testing"
)

func Test_quaminaPattern(t *testing.T) {
	tests := []struct {
		name    string
		checks  Checks
		want    []byte
		wantErr error
	}{
		{
			name:   "empty",
			checks: nil,
			want:   nil,
		},
		{
			name:   "id",
			checks: Checks{SameID("http://example.com")},
			want:   []byte(`{"id":["http://example.com"]}`),
		},
		{
			name:   "prefix id",
			checks: Checks{IDLike("http://example.com")},
			want:   []byte(`{"id":[{"prefix":"http://example.com"}]}`),
		},
		{
			name:   "id nil",
			checks: Checks{NilID},
			want:   []byte(`{"id":[{"exists":false}]}`),
		},
		{
			name:   "id",
			checks: Checks{SameIRI("http://example.com")},
			want:   []byte(`{"id":["http://example.com"]}`),
		},
		{
			name:   "prefix iri",
			checks: Checks{IRILike("http://example.com")},
			want:   []byte(`{"id":[{"prefix":"http://example.com"}]}`),
		},
		{
			name:   "IRI nil",
			checks: Checks{NilIRI},
			want:   []byte(`{"id":[{"exists":false}]}`),
		},
		{
			name:    "not idNil",
			checks:  Checks{Not(NilID)},
			want:    []byte(`{"id":[{"exists":true}]}`),
			wantErr: nil,
		},
		{
			name:    "not iriNil",
			checks:  Checks{Not(NilIRI)},
			want:    []byte(`{"id":[{"exists":true}]}`),
			wantErr: nil,
		},
		{
			name:    "not same ID",
			checks:  Checks{Not(SameID("https://example.com"))},
			want:    []byte(`{"id":[{"anything-but":"https://example.com"}]}`),
			wantErr: nil,
		},
		{
			name:    "not same IRI",
			checks:  Checks{Not(SameIRI("http://example.com"))},
			want:    []byte(`{"id":[{"anything-but":"http://example.com"}]}`),
			wantErr: nil,
		},
		{
			name:    "one type",
			checks:  Checks{HasType("Note")},
			want:    []byte(`{"type":["Note"]}`),
			wantErr: nil,
		},
		{
			name:    "multiple types",
			checks:  Checks{HasType("Note", "Article", "Image")},
			want:    []byte(`{"type":["Note","Article","Image"]}`),
			wantErr: nil,
		},
		{
			name:    "object with one type filter",
			checks:  Checks{Object(HasType("Note"))},
			want:    []byte(`{"object":{"type":["Note"]}}`),
			wantErr: nil,
		},
		{
			name:    "object with multiple filters",
			checks:  Checks{Object(HasType("Note"), NilID)},
			want:    []byte(`{"object":{"type":["Note"],"id":[{"exists":false}]}}`),
			wantErr: nil,
		},
		{
			name:    "name eq",
			checks:  Checks{NameIs("jdoe")},
			want:    []byte(`{"name":["jdoe"]}`),
			wantErr: nil,
		},
		{
			name:    "name empty",
			checks:  Checks{NameEmpty},
			want:    []byte(`{"name":[{"exists":false}]}`),
			wantErr: nil,
		},
		{
			name:    "name like",
			checks:  Checks{NameLike("test")},
			want:    []byte(`{"name":[{"prefix":"test"}]}`),
			wantErr: nil,
		},

		{
			name:    "summary eq",
			checks:  Checks{SummaryIs("jdoe")},
			want:    []byte(`{"summary":["jdoe"]}`),
			wantErr: nil,
		},
		{
			name:    "summary empty",
			checks:  Checks{SummaryEmpty},
			want:    []byte(`{"summary":[{"exists":false}]}`),
			wantErr: nil,
		},
		{
			name:    "summary like",
			checks:  Checks{SummaryLike("test")},
			want:    []byte(`{"summary":[{"prefix":"test"}]}`),
			wantErr: nil,
		},

		{
			name:    "content eq",
			checks:  Checks{ContentIs("jdoe")},
			want:    []byte(`{"content":["jdoe"]}`),
			wantErr: nil,
		},
		{
			name:    "content empty",
			checks:  Checks{ContentEmpty},
			want:    []byte(`{"content":[{"exists":false}]}`),
			wantErr: nil,
		},
		{
			name:    "content like",
			checks:  Checks{ContentLike("test")},
			want:    []byte(`{"content":[{"prefix":"test"}]}`),
			wantErr: nil,
		},

		{
			name:    "tag equal",
			checks:  Checks{Tag(NameIs("example"))},
			want:    []byte(`{"tag":{"name":["example"]}}`),
			wantErr: nil,
		},
		{
			name:    "tag like",
			checks:  Checks{Tag(NameLike("example"))},
			want:    []byte(`{"tag":{"name":[{"prefix":"example"}]}}`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quaminaPattern(tt.checks); !bytes.Equal(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s wanted %s", got, tt.want)
			}
		})
	}
}

func TestMatchRaw(t *testing.T) {
	type args struct {
		filters Checks
		raw     []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: true,
		},
		{
			name: "nil IRI",
			args: args{
				filters: Checks{NilIRI},
				raw:     []byte(`{"type":"Note"}`),
			},
			want: true,
		},
		{
			name: "nil IRI no match",
			args: args{
				filters: Checks{NilIRI},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: false,
		},
		{
			name: "single iri",
			args: args{
				filters: Checks{SameIRI("http://example.com")},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: true,
		},
		{
			name: "prefix iri",
			args: args{
				filters: Checks{IRILike("http://example.com")},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: true,
		},
		{
			name: "nil ID",
			args: args{
				filters: Checks{NilID},
				raw:     []byte(`{"type":"Note"}`),
			},
			want: true,
		},
		{
			name: "nil ID no match",
			args: args{
				filters: Checks{NilID},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: false,
		},
		{
			name: "single ID",
			args: args{
				filters: Checks{SameID("http://example.com")},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: true,
		},
		{
			name: "prefix ID",
			args: args{
				filters: Checks{IDLike("http://example.com")},
				raw:     []byte(`{"id":"http://example.com"}`),
			},
			want: true,
		},
		{
			name: "object with type and nil id matches",
			args: args{
				filters: Checks{Object(HasType("Note"), NilID)},
				raw:     []byte(`{"object":{"type":"Note"}}`),
			},
			want: true,
		},
		{
			name: "object with type and nil id does not match",
			args: args{
				filters: Checks{Object(HasType("Note"), NilID)},
				raw:     []byte(`{"object":{"type":"Note","id":"http://example.com"}}`),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchRaw(tt.args.filters, tt.args.raw); got != tt.want {
				t.Errorf("MatchRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}
