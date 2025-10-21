package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func dnl(v ...string) vocab.NaturalLanguageValues {
	nlv := make(vocab.NaturalLanguageValues)
	for _, vv := range v {
		nlv[vocab.NilLangRef] = vocab.Content(vv)
	}
	return nlv
}

func TestNameIs(t *testing.T) {
	type args struct {
		checkName    string
		toCheckNames vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkName:    "some name",
				toCheckNames: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkName:    "",
				toCheckNames: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkName:    "name",
				toCheckNames: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkName:    "日本語",
				toCheckNames: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkName:    "example",
				toCheckNames: dnl("not example"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Name: tt.args.toCheckNames}
			if got := NameIs(tt.args.checkName).Match(it); got != tt.want {
				t.Errorf("NameIs(%q)(Object.Name=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
			act := vocab.Actor{PreferredUsername: tt.args.toCheckNames}
			if got := NameIs(tt.args.checkName).Match(act); got != tt.want {
				t.Errorf("NameIs(%q)(Actor.PreferredName=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
		})
	}
}

func TestNameLike(t *testing.T) {
	type args struct {
		checkName    string
		toCheckNames vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkName:    "some name",
				toCheckNames: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkName:    "",
				toCheckNames: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkName:    "name",
				toCheckNames: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkName:    "日本語",
				toCheckNames: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring unicode name",
			args: args{
				checkName:    "日本",
				toCheckNames: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring",
			args: args{
				checkName:    "example",
				toCheckNames: dnl("not example"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkName:    "example",
				toCheckNames: dnl("not exampl"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Name: tt.args.toCheckNames}
			if got := NameLike(tt.args.checkName).Match(it); got != tt.want {
				t.Errorf("NameIs(%q)(Object.Name=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
			act := vocab.Actor{PreferredUsername: tt.args.toCheckNames}
			if got := NameLike(tt.args.checkName).Match(act); got != tt.want {
				t.Errorf("NameIs(%q)(Actor.PreferredName=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
		})
	}
}

func TestNameEmpty(t *testing.T) {
	tests := []struct {
		name         string
		toCheckNames vocab.NaturalLanguageValues
		want         bool
	}{
		{
			name:         "nil values",
			toCheckNames: nil,
			want:         true,
		},
		{
			name:         "empty values",
			toCheckNames: vocab.NaturalLanguageValues{},
			want:         true,
		},
		{
			name:         "single value",
			toCheckNames: dnl("not empty"),
			want:         false,
		},
		{
			name:         "multiple values",
			toCheckNames: dnl("not empty", "example"),
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Name: tt.toCheckNames}
			if got := NameEmpty.Match(it); got != tt.want {
				t.Errorf("NameEmpty()(Object.Name=%v) = %v, want %v", tt.toCheckNames, got, tt.want)
			}
			act := vocab.Actor{PreferredUsername: tt.toCheckNames}
			if got := NameEmpty.Match(act); got != tt.want {
				t.Errorf("NameEmpty()(Actor.PreferredName=%v) = %v, want %v", tt.toCheckNames, got, tt.want)
			}
		})
	}
}

func TestContentIs(t *testing.T) {
	type args struct {
		checkContent    string
		toCheckContents vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkContent:    "some name",
				toCheckContents: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkContent:    "",
				toCheckContents: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkContent:    "name",
				toCheckContents: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkContent:    "日本語",
				toCheckContents: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkContent:    "example",
				toCheckContents: dnl("not example"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Content: tt.args.toCheckContents}
			if got := ContentIs(tt.args.checkContent).Match(it); got != tt.want {
				t.Errorf("ContentIs(%q)(Object.Content=%v) = %v, want %v", tt.args.checkContent, tt.args.toCheckContents, got, tt.want)
			}
		})
	}
}

func TestContentLike(t *testing.T) {
	type args struct {
		checkContent    string
		toCheckContents vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkContent:    "some name",
				toCheckContents: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkContent:    "",
				toCheckContents: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkContent:    "name",
				toCheckContents: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkContent:    "日本語",
				toCheckContents: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring unicode name",
			args: args{
				checkContent:    "日本",
				toCheckContents: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring",
			args: args{
				checkContent:    "example",
				toCheckContents: dnl("not example"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkContent:    "example",
				toCheckContents: dnl("not exampl"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Content: tt.args.toCheckContents}
			if got := ContentLike(tt.args.checkContent).Match(it); got != tt.want {
				t.Errorf("ContentIs(%q)(Object.Content=%v) = %v, want %v", tt.args.checkContent, tt.args.toCheckContents, got, tt.want)
			}
		})
	}
}

func TestContentEmpty(t *testing.T) {
	tests := []struct {
		name            string
		toCheckContents vocab.NaturalLanguageValues
		want            bool
	}{
		{
			name:            "nil values",
			toCheckContents: nil,
			want:            true,
		},
		{
			name:            "empty values",
			toCheckContents: vocab.NaturalLanguageValues{},
			want:            true,
		},
		{
			name:            "single value",
			toCheckContents: dnl("not empty"),
			want:            false,
		},
		{
			name:            "multiple values",
			toCheckContents: dnl("not empty", "example"),
			want:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Content: tt.toCheckContents}
			if got := ContentEmpty.Match(it); got != tt.want {
				t.Errorf("ContentEmpty()(Object.Content=%v) = %v, want %v", tt.toCheckContents, got, tt.want)
			}
		})
	}
}

func TestSummaryIs(t *testing.T) {
	type args struct {
		checkSummary    string
		toCheckSummarys vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkSummary:    "some name",
				toCheckSummarys: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkSummary:    "",
				toCheckSummarys: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkSummary:    "name",
				toCheckSummarys: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkSummary:    "日本語",
				toCheckSummarys: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkSummary:    "example",
				toCheckSummarys: dnl("not example"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Summary: tt.args.toCheckSummarys}
			if got := SummaryIs(tt.args.checkSummary).Match(it); got != tt.want {
				t.Errorf("SummaryIs(%q)(Object.Summary=%v) = %v, want %v", tt.args.checkSummary, tt.args.toCheckSummarys, got, tt.want)
			}
		})
	}
}

func TestSummaryLike(t *testing.T) {
	type args struct {
		checkSummary    string
		toCheckSummarys vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				checkSummary:    "some name",
				toCheckSummarys: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				checkSummary:    "",
				toCheckSummarys: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				checkSummary:    "name",
				toCheckSummarys: dnl("name"),
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkSummary:    "日本語",
				toCheckSummarys: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring unicode name",
			args: args{
				checkSummary:    "日本",
				toCheckSummarys: dnl("日本語"),
			},
			want: true,
		},
		{
			name: "matching substring",
			args: args{
				checkSummary:    "example",
				toCheckSummarys: dnl("not example"),
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkSummary:    "example",
				toCheckSummarys: dnl("not exampl"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Summary: tt.args.toCheckSummarys}
			if got := SummaryLike(tt.args.checkSummary).Match(it); got != tt.want {
				t.Errorf("SummaryIs(%q)(Object.Summary=%v) = %v, want %v", tt.args.checkSummary, tt.args.toCheckSummarys, got, tt.want)
			}
		})
	}
}

func TestSummaryEmpty(t *testing.T) {
	tests := []struct {
		name            string
		toCheckSummarys vocab.NaturalLanguageValues
		want            bool
	}{
		{
			name:            "nil values",
			toCheckSummarys: nil,
			want:            true,
		},
		{
			name:            "empty values",
			toCheckSummarys: vocab.NaturalLanguageValues{},
			want:            true,
		},
		{
			name:            "single value",
			toCheckSummarys: dnl("not empty"),
			want:            false,
		},
		{
			name:            "multiple values",
			toCheckSummarys: dnl("not empty", "example"),
			want:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Summary: tt.toCheckSummarys}
			if got := SummaryEmpty.Match(it); got != tt.want {
				t.Errorf("SummaryEmpty()(Object.Summary=%v) = %v, want %v", tt.toCheckSummarys, got, tt.want)
			}
		})
	}
}
