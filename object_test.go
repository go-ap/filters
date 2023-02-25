package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func dnl(v string) vocab.LangRefValue {
	return nl("-", v)
}
func nl(ref string, v string) vocab.LangRefValue {
	return vocab.LangRefValue{Ref: vocab.LangRef(ref), Value: vocab.Content(v)}
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
				toCheckNames: vocab.NaturalLanguageValues{dnl("name")},
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				checkName:    "日本語",
				toCheckNames: vocab.NaturalLanguageValues{dnl("日本語")},
			},
			want: true,
		},
		{
			name: "not matching",
			args: args{
				checkName:    "example",
				toCheckNames: vocab.NaturalLanguageValues{dnl("not example")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Name: tt.args.toCheckNames}
			if got := NameIs(tt.args.checkName)(it); tt.want != got {
				t.Errorf("NameIs(%q)(Object.Name=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
			act := vocab.Actor{PreferredUsername: tt.args.toCheckNames}
			if got := NameIs(tt.args.checkName)(act); tt.want != got {
				t.Errorf("NameIs(%q)(Actor.PreferredName=%v) = %v, want %v", tt.args.checkName, tt.args.toCheckNames, got, tt.want)
			}
		})
	}
}

func TestType(t *testing.T) {
	type args struct {
		checkTypes  vocab.ActivityVocabularyTypes
		toCheckType vocab.ActivityVocabularyType
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
			name: "empty object",
			args: args{
				checkTypes: vocab.ObjectTypes,
			},
			want: false,
		},
		{
			name: "empty check types",
			args: args{
				toCheckType: vocab.CreateType,
			},
			want: false,
		},
		{
			name: "matching types",
			args: args{
				checkTypes:  vocab.ActivityVocabularyTypes{vocab.CreateType, vocab.UpdateType},
				toCheckType: vocab.CreateType,
			},
			want: true,
		},
		{
			name: "non matching types",
			args: args{
				checkTypes:  vocab.ActivityVocabularyTypes{vocab.CreateType, vocab.UpdateType},
				toCheckType: vocab.FollowType,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{Type: tt.args.toCheckType}
			if got := HasType(tt.args.checkTypes...)(&ob); got != tt.want {
				t.Errorf("Type(%v)(Object.Type=%v) = %v, want %v", tt.args.checkTypes, tt.args.toCheckType, got, tt.want)
			}
		})
	}
}

func TestID(t *testing.T) {
	type args struct {
		checkIRI   vocab.IRI
		toCheckIRI vocab.IRI
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
			name: "empty check iri",
			args: args{
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "empty iri",
			args: args{
				checkIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "matching iris",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.com",
			},
			want: true,
		},
		{
			name: "non matching iris - different scheme",
			args: args{
				checkIRI:   "https://example.com",
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "non matching iris - different domain",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.org",
			},
			want: false,
		},
		{
			name: "non matching iris - different path",
			args: args{
				checkIRI:   "http://example.com/index",
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "matching iris - path is root vs no root",
			args: args{
				checkIRI:   "http://example.com/",
				toCheckIRI: "http://example.com",
			},
			want: true,
		},
		{
			name: "non matching iris - different query params",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.com/?ana=are",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := ID(tt.args.checkIRI)(ob); got != tt.want {
				t.Errorf("ID(%s)(Object.ID=%s) = %v, want %v", tt.args.checkIRI, tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}

func TestNilID(t *testing.T) {
	type args struct {
		toCheckIRI vocab.IRI
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
			name: "non nil iri",
			args: args{
				toCheckIRI: "http://example.org",
			},
			want: false,
		},
		{
			name: "empty IRI",
			args: args{
				toCheckIRI: "",
			},
			want: true,
		},
		{
			name: "nil IRI",
			args: args{
				toCheckIRI: "-",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := NilID(ob); got != tt.want {
				t.Errorf("NilID(Object.ID=%s) = %v, want %v", tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}

func TestNotNilID(t *testing.T) {
	type args struct {
		toCheckIRI vocab.IRI
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
			name: "non nil iri",
			args: args{
				toCheckIRI: "http://example.org",
			},
			want: true,
		},
		{
			name: "empty IRI",
			args: args{
				toCheckIRI: "",
			},
			want: false,
		},
		{
			name: "nil IRI",
			args: args{
				toCheckIRI: "-",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := NotNilID(ob); got != tt.want {
				t.Errorf("NotNilID(Object.ID=%s) = %v, want %v", tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}

func TestIRI(t *testing.T) {
	type args struct {
		checkIRI   vocab.IRI
		toCheckIRI vocab.IRI
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
			name: "empty check iri",
			args: args{
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "empty iri",
			args: args{
				checkIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "matching iris",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.com",
			},
			want: true,
		},
		{
			name: "non matching iris - different scheme",
			args: args{
				checkIRI:   "https://example.com",
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "non matching iris - different domain",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.org",
			},
			want: false,
		},
		{
			name: "non matching iris - different path",
			args: args{
				checkIRI:   "http://example.com/index",
				toCheckIRI: "http://example.com",
			},
			want: false,
		},
		{
			name: "matching iris - path is root vs no root",
			args: args{
				checkIRI:   "http://example.com/",
				toCheckIRI: "http://example.com",
			},
			want: true,
		},
		{
			name: "non matching iris - different query params",
			args: args{
				checkIRI:   "http://example.com",
				toCheckIRI: "http://example.com/?ana=are",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameIRI(tt.args.checkIRI)(tt.args.toCheckIRI); got != tt.want {
				t.Errorf("IRI(%s)(%s) = %v, want %v", tt.args.checkIRI, tt.args.toCheckIRI, got, tt.want)
			}
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := SameIRI(tt.args.checkIRI)(ob); got != tt.want {
				t.Errorf("IRI(%s)(Object.ID=%s) = %v, want %v", tt.args.checkIRI, tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}
