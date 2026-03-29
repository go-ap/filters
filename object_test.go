package filters

import (
	"reflect"
	"runtime"
	"strings"
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
			want: true,
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
			if got := HasType(tt.args.checkTypes...).Match(&ob); got != tt.want {
				t.Errorf("Type(%v).Match(Object.Type=%v) = %v, want %v", tt.args.checkTypes, tt.args.toCheckType, got, tt.want)
			}
		})
	}
}

func TestID(t *testing.T) {
	tests := []struct {
		name    string
		arg     vocab.IRI
		matchTo vocab.Item
		want    bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name:    "empty check iri",
			matchTo: &vocab.Object{ID: "http://example.com"},
			want:    false,
		},
		{
			name: "empty iri",
			arg:  "http://example.com",
			want: false,
		},
		{
			name:    "matching iris",
			arg:     "http://example.com",
			matchTo: &vocab.Object{ID: "http://example.com"},
			want:    true,
		},
		{
			name:    "non matching iris - different scheme",
			arg:     "https://example.com",
			matchTo: &vocab.Object{ID: "http://example.com"},
			want:    true,
		},
		{
			name:    "non matching iris - different domain",
			arg:     "http://example.com",
			matchTo: &vocab.Object{ID: "http://example.org"},
			want:    false,
		},
		{
			name:    "non matching iris - different path",
			arg:     "http://example.com/index",
			matchTo: &vocab.Object{ID: "http://example.com"},
			want:    false,
		},
		{
			name:    "matching iris - path is root vs no root",
			arg:     "http://example.com/",
			matchTo: &vocab.Object{ID: "http://example.com"},
			want:    true,
		},
		{
			name:    "non matching iris - different query params",
			arg:     "http://example.com",
			matchTo: &vocab.Object{ID: "http://example.com/?ana=are"},
			want:    false,
		},
		{
			name:    "iri in list matches",
			arg:     "http://example.com",
			matchTo: vocab.ItemCollection{&vocab.Object{ID: "http://example.com"}},
			want:    true,
		},
		{
			name:    "iri in list doesn't match",
			arg:     "http://example.com",
			matchTo: vocab.ItemCollection{&vocab.Object{ID: "http://example.com/ana"}, vocab.IRI("http://no.example.com")},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameID(tt.arg).Match(tt.matchTo); got != tt.want {
				t.Errorf("SameID(%s).Match(Object.ID=%s) = %v, want %v", tt.arg, tt.matchTo, got, tt.want)
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
			if got := NilID.Match(ob); got != tt.want {
				t.Errorf("NilID(Object.ID=%s) = %v, want %v", tt.args.toCheckIRI, got, tt.want)
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
			if got := SameIRI(tt.args.checkIRI).Match(tt.args.toCheckIRI); got != tt.want {
				t.Errorf("IRI(%s).Match(%s) = %v, want %v", tt.args.checkIRI, tt.args.toCheckIRI, got, tt.want)
			}
			ob := vocab.Object{ID: tt.args.toCheckIRI}
			if got := SameIRI(tt.args.checkIRI).Match(ob); got != tt.want {
				t.Errorf("IRI(%s).Match(Object.ID=%s) = %v, want %v", tt.args.checkIRI, tt.args.toCheckIRI, got, tt.want)
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
			if got := IDLike(tt.args.checkFrag).Match(ob); got != tt.want {
				t.Errorf("IDLike(%s).Match(Object.ID=%v) = %v, want %v", tt.args.checkFrag, tt.args.toCheckIRI, got, tt.want)
			}
		})
	}
}

func testAccumFn(accumFn func(vocab.Item) vocab.IRIs) func(*testing.T) {
	accumFnName := runtime.FuncForPC(reflect.ValueOf(accumFn).Pointer()).Name()
	if idx := strings.LastIndex(accumFnName, ".") + 1; idx < len(accumFnName) {
		accumFnName = accumFnName[idx:]
	}
	return func(t *testing.T) {
		tests := []struct {
			name string
			item vocab.Item
			want vocab.IRIs
		}{
			{
				name: "empty Item",
			},
			{
				name: "empty InReplyTo",
				item: &vocab.Object{},
			},
			{
				name: "one InReplyTo IRI",
				item: &vocab.Object{InReplyTo: vocab.IRI("https://example.com")},
				want: vocab.IRIs{"https://example.com"},
			},
			{
				name: "two InReplyTo IRIs",
				item: &vocab.Object{InReplyTo: vocab.IRIs{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/one")}},
				want: vocab.IRIs{"https://example.com", "https://example.com/one"},
			},
			{
				name: "two InReplyTo IRIs as Items",
				item: &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/one")}},
				want: vocab.IRIs{"https://example.com", "https://example.com/one"},
			},
			{
				name: "two InReplyTo Items",
				item: &vocab.Object{InReplyTo: vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}, &vocab.Profile{ID: "https://example.com/one"}}},
				want: vocab.IRIs{"https://example.com", "https://example.com/one"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := accumFn(tt.item)
				if len(got) != len(tt.want) {
					t.Errorf("%s() has %d items, wanted %d items", accumFnName, len(got), len(tt.want))
					return
				}
				for i, it := range tt.want {
					git := got[i]
					if !git.Equal(it) {
						t.Errorf("%s() at pos %d = %v, want %v", accumFnName, i, git.GetLink(), it)
					}
				}
			})
		}
	}
}

func Test_accumInReplyTos(t *testing.T) {
	testAccumFn(accumInReplyTos)
}

func Test_accumContexts(t *testing.T) {
	testAccumFn(accumContexts)
}

func Test_accumAttributedTos(t *testing.T) {
	testAccumFn(accumAttributedTos)
}

func Test_accumURLs(t *testing.T) {
	testAccumFn(accumURLs)
}

func Test_inReplyToEquals_Match(t *testing.T) {
	tests := []struct {
		name string
		i    inReplyToEquals
		it   vocab.Item
		want bool
	}{
		{
			name: "empty equals empty",
			i:    "",
			want: true,
		},
		{
			name: "https://example.com does not equal empty",
			i:    "https://example.com",
			want: false,
		},
		{
			name: "https://example.com w/ empty inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{},
			want: false,
		},
		{
			name: "https://example.com not equal w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRI("https://example.com/not-equal")},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.Object{ID: "https://example.com/not-equal"}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRIs{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRIs{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}}},
			want: true,
		},
		//
		{
			name: "https://example.com not equal w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRIs{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.IRIs{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}, vocab.Object{ID: "https://example.com/still-not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object inReplyTo",
			i:    "https://example.com",
			it:   &vocab.Object{InReplyTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}, vocab.Object{ID: "https://example.com/not-equal"}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Match(tt.it); got != tt.want {
				t.Errorf("inReplyTo.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attributedToEquals_Match(t *testing.T) {
	tests := []struct {
		name string
		i    attributedToEquals
		it   vocab.Item
		want bool
	}{
		{
			name: "empty equals empty",
			i:    "",
			want: true,
		},
		{
			name: "https://example.com does not equal empty",
			i:    "https://example.com",
			want: false,
		},
		{
			name: "https://example.com w/ empty attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{},
			want: false,
		},
		{
			name: "https://example.com not equal w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRI("https://example.com/not-equal")},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.Object{ID: "https://example.com/not-equal"}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRIs{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRIs{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}}},
			want: true,
		},
		//
		{
			name: "https://example.com not equal w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRIs{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.IRIs{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}, vocab.Object{ID: "https://example.com/still-not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object attributedTo",
			i:    "https://example.com",
			it:   &vocab.Object{AttributedTo: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}, vocab.Object{ID: "https://example.com/not-equal"}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Match(tt.it); got != tt.want {
				t.Errorf("attributedTo.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextEquals_Match(t *testing.T) {
	tests := []struct {
		name string
		i    contextEquals
		it   vocab.Item
		want bool
	}{
		{
			name: "empty equals empty",
			i:    "",
			want: true,
		},
		{
			name: "https://example.com does not equal empty",
			i:    "https://example.com",
			want: false,
		},
		{
			name: "https://example.com w/ empty context",
			i:    "https://example.com",
			it:   &vocab.Object{},
			want: false,
		},
		{
			name: "https://example.com not equal w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRI("https://example.com/not-equal")},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.Object{ID: "https://example.com/not-equal"}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRIs{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRIs{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}}},
			want: true,
		},
		//
		{
			name: "https://example.com not equal w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRIs{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.IRIs{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}, vocab.Object{ID: "https://example.com/still-not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object context",
			i:    "https://example.com",
			it:   &vocab.Object{Context: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}, vocab.Object{ID: "https://example.com/not-equal"}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Match(tt.it); got != tt.want {
				t.Errorf("context.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_URLEquals_Match(t *testing.T) {
	tests := []struct {
		name string
		i    urlEquals
		it   vocab.Item
		want bool
	}{
		{
			name: "empty equals empty",
			i:    "",
			want: true,
		},
		{
			name: "https://example.com does not equal empty",
			i:    "https://example.com",
			want: false,
		},
		{
			name: "https://example.com w/ empty URL",
			i:    "https://example.com",
			it:   &vocab.Object{},
			want: false,
		},
		{
			name: "https://example.com not equal w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRI("https://example.com/not-equal")},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.Object{ID: "https://example.com/not-equal"}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRIs{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRIs{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}}},
			want: true,
		},
		//
		{
			name: "https://example.com not equal w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRIs{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.IRIs{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.IRI("https://example.com/not-equal"), vocab.IRI("https://example.com/still-not-equal")}},
			want: false,
		},
		{
			name: "https://example.com equals w/ IRI URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.IRI("https://example.com"), vocab.IRI("https://example.com/not-equal")}},
			want: true,
		},
		{
			name: "https://example.com not equal w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.Object{ID: "https://example.com/not-equal"}, vocab.Object{ID: "https://example.com/still-not-equal"}}},
			want: false,
		},
		{
			name: "https://example.com equals w/ Object URL",
			i:    "https://example.com",
			it:   &vocab.Object{URL: vocab.ItemCollection{vocab.Object{ID: "https://example.com"}, vocab.Object{ID: "https://example.com/not-equal"}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Match(tt.it); got != tt.want {
				t.Errorf("URL.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_idLike_GoString(t *testing.T) {
	tests := []struct {
		name string
		l    idLike
		want string
	}{
		{
			name: "empty",
			l:    "",
			want: "id=~",
		},
		{
			name: "example.com",
			l:    "example.com",
			want: "id=~example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GoString(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_idEquals_GoString(t *testing.T) {
	tests := []struct {
		name string
		i    idEquals
		want string
	}{
		{
			name: "empty",
			i:    "",
			want: "id=",
		},
		{
			name: "example.com",
			i:    "example.com",
			want: "id=example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GoString(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inReplyToLike_Match(t *testing.T) {
	tests := []struct {
		name string
		a    inReplyToLike
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "IRI does not match",
			it:   vocab.IRI("https://example.com"),
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
		{
			name: "IRIs don't match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/o inReplyTo does not match",
			it:   &vocab.Object{ID: "https://example.com"},
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/ same inReplyTo does not match",
			it:   &vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://example.com")},
			a:    inReplyToLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ like inReplyTo does not match",
			it:   &vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://example.com/~jdoe")},
			a:    inReplyToLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ different inReplyTo does not match",
			it:   &vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://social.example.com")},
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} w/ same inReplyTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://example.com")}},
			a:    inReplyToLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ like inReplyTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://example.com/~jdoe")}},
			a:    inReplyToLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ wrong inReplyTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("https://social.example.com")}},
			a:    inReplyToLike("https://example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inReplyToNil_Match(t *testing.T) {
	tests := []struct {
		name string
		c    inReplyToNil
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "nil IRI",
			it:   vocab.NilIRI,
			want: true,
		},
		{
			name: "empty IRI",
			it:   vocab.EmptyIRI,
			want: true,
		},
		{
			name: "IRI matches",
			it:   vocab.IRI("https://example.com"),
			want: true,
		},
		{
			name: "IRIs match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			want: true,
		},
		{
			name: "not nil Items{IRI} match",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "Object w/o inReplyTo matches",
			it:   &vocab.Object{ID: "https://example.com"},
			want: true,
		},
		{
			name: "Object w/ inReplyTo does not match",
			it:   &vocab.Object{ID: "https://example.com/one", InReplyTo: vocab.IRI("http://example.com/zero")},
			want: false,
		},
		{
			name: "Items{Object} w/ inReplyTo matches",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "Items{Object} w/ inReplyTo does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", InReplyTo: vocab.IRI("http://example.com/t")}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attributedToNil_Match(t *testing.T) {
	tests := []struct {
		name string
		a    attributedToNil
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "nil IRI",
			it:   vocab.NilIRI,
			want: true,
		},
		{
			name: "empty IRI",
			it:   vocab.EmptyIRI,
			want: true,
		},
		{
			name: "IRI matches",
			it:   vocab.IRI("https://example.com"),
			want: true,
		},
		{
			name: "IRIs match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			want: true,
		},
		{
			name: "not nil Items{IRI} match",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "Object w/o attributedTo matches",
			it:   &vocab.Object{ID: "https://example.com"},
			want: true,
		},
		{
			name: "Object w/ attributedTo does not match",
			it:   &vocab.Object{ID: "https://example.com/one", AttributedTo: vocab.IRI("http://example.com/zero")},
			want: false,
		},
		{
			name: "Items{Object} w/ attributedTo matches",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "Items{Object} w/ attributedTo does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("http://example.com/t")}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attributedToLike_Match(t *testing.T) {
	tests := []struct {
		name string
		a    attributedToLike
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "IRI does not match",
			it:   vocab.IRI("https://example.com"),
			a:    attributedToLike("https://example.com"),
			want: false,
		},
		{
			name: "IRIs don't match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			a:    attributedToLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/o attributedTo does not match",
			it:   &vocab.Object{ID: "https://example.com"},
			a:    attributedToLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			a:    attributedToLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/ same attributedTo does not match",
			it:   &vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://example.com")},
			a:    attributedToLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ like attributedTo does not match",
			it:   &vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://example.com/~jdoe")},
			a:    attributedToLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ different attributedTo does not match",
			it:   &vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://social.example.com")},
			a:    attributedToLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} w/ same attributedTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://example.com")}},
			a:    attributedToLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ like attributedTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://example.com/~jdoe")}},
			a:    attributedToLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ wrong attributedTo match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", AttributedTo: vocab.IRI("https://social.example.com")}},
			a:    attributedToLike("https://example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextLike_Match(t *testing.T) {
	tests := []struct {
		name string
		c    contextLike
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "IRI does not match",
			it:   vocab.IRI("https://example.com"),
			c:    contextLike("https://example.com"),
			want: false,
		},
		{
			name: "IRIs don't match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			c:    contextLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/o context does not match",
			it:   &vocab.Object{ID: "https://example.com"},
			c:    contextLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			c:    contextLike("https://example.com"),
			want: false,
		},
		{
			name: "Object w/ same context does not match",
			it:   &vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://example.com")},
			c:    contextLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ like context does not match",
			it:   &vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://example.com/~jdoe")},
			c:    contextLike("https://example.com"),
			want: true,
		},
		{
			name: "Object w/ different context does not match",
			it:   &vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://social.example.com")},
			c:    contextLike("https://example.com"),
			want: false,
		},
		{
			name: "Items{Object} w/ same context match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://example.com")}},
			c:    contextLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ like context match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://example.com/~jdoe")}},
			c:    contextLike("https://example.com"),
			want: true,
		},
		{
			name: "Items{Object} w/ wrong context match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", Context: vocab.IRI("https://social.example.com")}},
			c:    contextLike("https://example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextNil_Match(t *testing.T) {
	tests := []struct {
		name string
		c    contextNil
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "nil IRI",
			it:   vocab.NilIRI,
			want: true,
		},
		{
			name: "empty IRI",
			it:   vocab.EmptyIRI,
			want: true,
		},
		{
			name: "IRI matches",
			it:   vocab.IRI("https://example.com"),
			want: true,
		},
		{
			name: "IRIs match",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			want: true,
		},
		{
			name: "not nil Items{IRI} match",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			want: true,
		},
		{
			name: "Object w/o context matches",
			it:   &vocab.Object{ID: "https://example.com"},
			want: true,
		},
		{
			name: "Object w/ context does not match",
			it:   &vocab.Object{ID: "https://example.com/one", Context: vocab.IRI("http://example.com/zero")},
			want: false,
		},
		{
			name: "Items{Object} w/ context matches",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			want: true,
		},
		{
			name: "Items{Object} w/ context does not match",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com", Context: vocab.IRI("http://example.com/t")}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
