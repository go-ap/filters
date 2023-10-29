package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

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
		toCheckIRI vocab.IRI
		checkIRI   vocab.IRI
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

func TestIDLike(t *testing.T) {
	type args struct {
		toCheckIRI vocab.IRI
		checkFrag  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty is contained",
			args: args{
				checkFrag:  "",
				toCheckIRI: "https://example.com",
			},
			want: true,
		},
		{
			name: "empty is contained by empty",
			args: args{
				checkFrag:  "",
				toCheckIRI: "",
			},
			want: true,
		},
		{
			name: "something is not contained by empty",
			args: args{
				checkFrag:  "something",
				toCheckIRI: "",
			},
			want: false,
		},
		{
			name: "something is not contained by https://example",
			args: args{
				checkFrag:  "something",
				toCheckIRI: "https://example",
			},
			want: false,
		},
		{
			name: "http://example is not contained by https://example",
			args: args{
				checkFrag:  "http://example",
				toCheckIRI: "https://example",
			},
			want: false,
		},
		{
			name: "https://example/ is not contained by https://example",
			args: args{
				checkFrag:  "https://example/",
				toCheckIRI: "https://example",
			},
			want: false,
		},
		{
			name: "https://example is contained by https://example/test",
			args: args{
				checkFrag:  "https://example",
				toCheckIRI: "https://example/test",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := IDLike(tt.args.checkFrag)(ob); got != tt.want {
				t.Errorf("IDLike(%s)(Object.ID=%v) = %v, want %v", tt.args.checkFrag, tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}
