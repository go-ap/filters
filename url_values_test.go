package filters

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

var mockURL = url.URL{
	Scheme: "https",
	Host:   "example.com",
	Path:   "/",
}

func withValues(u url.URL, q url.Values) url.URL {
	u.RawQuery = q.Encode()
	return u
}

func withAuth(u url.URL, user *url.Userinfo) url.URL {
	u.User = user
	return u
}

func TestFromIRI(t *testing.T) {
	tests := []struct {
		name    string
		iri     vocab.IRI
		item    vocab.Item
		want    vocab.Item
		wantErr error
	}{
		{name: "empty"},
		{
			name: "no auth, no query values",
			iri:  "https://example.com",
		},
		{
			name: "some query values",
			iri:  "https://example.com?type=Create&type=Follow",
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
			},
		},
		{
			name: "invalid IRI",
			iri:  ":/example-com",
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			wantErr: &url.Error{"parse", ":/example-com", errors.New("missing protocol scheme")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFn, gotErr := FromIRI(tt.iri)
			if (gotErr != nil || tt.wantErr != nil) && !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Error returned FromIRI().Run() = %s, want %s", gotErr, tt.wantErr)
			}
			if got := gotFn.Run(tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromIRI().Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestFromURL(t *testing.T) {
	tests := []struct {
		name string
		u    url.URL
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "no auth, no query values",
			u:    mockURL,
		},
		{
			name: "some query values",
			u:    withValues(mockURL, vals(kv("type", "Create", "Follow"))),
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromURL(tt.u).Run(tt.item)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromURL(%v).Run() = %v, want %v", tt.u, got, tt.want)
			}
		})
	}
}

func TestFromValues(t *testing.T) {
	tests := []struct {
		name string
		v    url.Values
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "no auth, no query values",
			v:    vals(),
		},
		{
			name: "type=Create,Follow",
			v:    vals(kv("type", "Create", "Follow")),
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
			},
		},
		{
			name: "ID=https://example.com",
			v:    vals(kv("id", "https://example.com")),
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
			},
		},
		{
			name: "ID=(null)",
			v:    vals(kv("id", "")),
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Follow{Type: "Follow"},
			},
		},
		{
			name: "ID=(not null)",
			v:    vals(kv("id", "!")),
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
		{
			name: "ID=(not nil IRI)",
			v:    vals(kv("id", "!-")),
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
		{
			name: "ID=~example.com",
			v:    vals(kv("id", "~example.com")),
			item: vocab.ItemCollection{
				vocab.Object{ID: "https://activitypub.rocks", Type: "Page"},
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromValues(tt.v).Run(tt.item)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromValues(%v).Run() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func Test_paginationFromValues(t *testing.T) {
	t.Parallel()

	parseQuery := func(s string) url.Values {
		v, _ := url.ParseQuery(s)
		return v
	}

	tests := []struct {
		name string
		q    url.Values
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "after: item is last, nothing after",
			q:    parseQuery("after=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "before-user-2"},
				vocab.Object{ID: "before-user-1"},
				vocab.Object{ID: "user"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "after: some after items",
			q:    parseQuery("after=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "before-2"},
				vocab.Object{ID: "before-1"},
				vocab.Object{ID: "user"},
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: vocab.IRI("after-1")},
				vocab.Object{ID: vocab.IRI("after-2")},
			},
		},
		{
			name: "after: item is first, everything but itself",
			q:    parseQuery("after=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "user"},
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
				vocab.Object{ID: "after-3"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
				vocab.Object{ID: "after-3"},
			},
		},
		{
			name: "before: item is last, everything before",
			q:    parseQuery("before=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "before-2"},
				vocab.Object{ID: "before-1"},
				vocab.Object{ID: "user"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "before-2"},
				vocab.Object{ID: "before-1"},
			},
		},
		{
			name: "before: some before items",
			q:    parseQuery("before=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "before-2"},
				vocab.Object{ID: "before-1"},
				vocab.Object{ID: "user"},
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: vocab.IRI("before-2")},
				vocab.Object{ID: vocab.IRI("before-1")},
			},
		},
		{
			name: "before: item is first, nothing",
			q:    parseQuery("before=user"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "user"},
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
				vocab.Object{ID: "after-3"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "before and after",
			q:    parseQuery("before=stop&after=start"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "before-3"},
				vocab.Object{ID: "before-2"},
				vocab.Object{ID: "before-1"},
				vocab.Object{ID: "start"},
				vocab.Object{ID: "example1"},
				vocab.Object{ID: "example2"},
				vocab.Object{ID: "example3"},
				vocab.Object{ID: "stop"},
				vocab.Object{ID: "after-1"},
				vocab.Object{ID: "after-2"},
				vocab.Object{ID: "after-3"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "example1"},
				vocab.Object{ID: "example2"},
				vocab.Object{ID: "example3"},
			},
		},
		{
			name: "maxItems=0",
			q:    parseQuery("maxItems=0"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "maxItems=0"},
				vocab.Object{ID: "not-1"},
				vocab.Object{ID: "not-2"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "maxItems=2",
			q:    parseQuery("maxItems=2"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
				vocab.Object{ID: "maxItems=2"},
				vocab.Object{ID: "not-1"},
				vocab.Object{ID: "not-2"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
			},
		},
		{
			name: "after=user&maxItems=2",
			q:    parseQuery("after=user&maxItems=2"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "not-1"},
				vocab.Object{ID: "not-2"},
				vocab.Object{ID: "not-3"},
				vocab.Object{ID: "user"},
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
				vocab.Object{ID: "maxItems=2"},
				vocab.Object{ID: "not-5"},
				vocab.Object{ID: "not-6"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
			},
		},
		{
			name: "before=user&maxItems=2",
			q:    parseQuery("before=user&maxItems=2"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
				vocab.Object{ID: "maxItems=2"},
				vocab.Object{ID: "not-1"},
				vocab.Object{ID: "not-2"},
				vocab.Object{ID: "user"},
				vocab.Object{ID: "not-3"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
			},
		},
		{
			name: "after=start&before=end&maxItems=2",
			q:    parseQuery("after=start&before=stop&maxItems=2"),
			item: vocab.ItemCollection{
				vocab.Object{ID: "not-1"},
				vocab.Object{ID: "not-2"},
				vocab.Object{ID: "start"},
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
				vocab.Object{ID: "maxItems=2"},
				vocab.Object{ID: "not-3"},
				vocab.Object{ID: "stop"},
				vocab.Object{ID: "not-4"},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "good1"},
				vocab.Object{ID: "good2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFns := paginationFromValues(tt.q)
			if got := gotFns.Run(tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paginationFromValues().Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromValues(t *testing.T) {
	tests := []struct {
		name string
		arg  url.Values
		want Checks
	}{
		{
			name: "empty",
		},
		// ID
		{
			name: "Empty ID",
			arg:  vals(kv("id", "")),
			want: Checks{NilID},
		},
		{
			name: "Not empty ID",
			arg:  vals(kv("id", "!")),
			want: Checks{NotNilID},
		},
		{
			name: "ID like",
			arg:  vals(kv("id", "~test")),
			want: Checks{IDLike("test")},
		},
		{
			name: "ID equals",
			arg:  vals(kv("id", "https://example.com")),
			want: Checks{SameID("https://example.com")},
		},
		// AttributedTo
		{
			name: "Empty attributedTo",
			arg:  vals(kv("attributedTo", "")),
			want: Checks{NilAttributedTo},
		},
		{
			name: "Not empty attributedTo",
			arg:  vals(kv("attributedTo", "!")),
			want: Checks{Not(NilAttributedTo)},
		},
		{
			name: "attributedTo like",
			arg:  vals(kv("attributedTo", "~test")),
			want: Checks{AttributedToLike("test")},
		},
		{
			name: "attributedTo equals",
			arg:  vals(kv("attributedTo", "https://example.com")),
			want: Checks{SameAttributedTo("https://example.com")},
		},
		// context
		{
			name: "Empty context",
			arg:  vals(kv("context", "")),
			want: Checks{NilContext},
		},
		{
			name: "Not empty context",
			arg:  vals(kv("context", "!")),
			want: Checks{Not(NilContext)},
		},
		{
			name: "context like",
			arg:  vals(kv("context", "~test")),
			want: Checks{ContextLike("test")},
		},
		{
			name: "context equals",
			arg:  vals(kv("context", "https://example.com")),
			want: Checks{SameContext("https://example.com")},
		},
		// URL
		{
			name: "Empty URL",
			arg:  vals(kv("url", "")),
			want: Checks{NilIRI},
		},
		{
			name: "Not empty URL",
			arg:  vals(kv("url", "!")),
			want: Checks{Not(NilIRI)},
		},
		{
			name: "URL like",
			arg:  vals(kv("url", "~test")),
			want: Checks{URLLike("test")},
		},
		{
			name: "URL equals",
			arg:  vals(kv("url", "https://example.com")),
			want: Checks{SameURL("https://example.com")},
		},
		// Name
		{
			name: "Empty Name",
			arg:  vals(kv("name", "")),
			want: Checks{NameEmpty},
		},
		{
			name: "Not empty Name",
			arg:  vals(kv("name", "!")),
			want: Checks{Not(NameEmpty)},
		},
		{
			name: "Name like",
			arg:  vals(kv("name", "~test")),
			want: Checks{NameLike("test")},
		},
		{
			name: "Name equals",
			arg:  vals(kv("name", "john doe")),
			want: Checks{NameIs("john doe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.Contains(tt.name, "Name") {
				t.Skipf("Skipping because the reflect.DeepEquals doesn't work for NaturalLanguageValChecks as they contain function pointers")
			}
			if got := fromValues(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func kv(key string, values ...string) url.Values {
	return url.Values{
		key: values,
	}
}

func vals(v ...url.Values) url.Values {
	res := url.Values{}
	for _, vv := range v {
		for k, vvv := range vv {
			for _, s := range vvv {
				res.Add(k, s)
			}
		}
	}
	return res
}

func Test_urlValue(t *testing.T) {
	tests := []struct {
		name string
		arg  Check
		want url.Values
	}{
		{
			name: "empty",
			arg:  nil,
			want: nil,
		},
		{
			name: "id",
			arg:  idEquals("https://example.com"),
			want: vals(kv(keyID, "https://example.com")),
		},
		{
			name: "iri",
			arg:  iriEquals("https://example.com"),
			want: vals(kv(keyIRI, "https://example.com")),
		},
		{
			name: "nil iri",
			arg:  iriNil{},
			want: vals(kv(keyIRI, "")),
		},
		{
			name: "iri like",
			arg:  iriLike("https://example.com"),
			want: vals(kv(keyIRI, opLike+"https://example.com")),
		},
		{
			name: "maxItems",
			arg:  WithMaxCount(666),
			want: vals(kv(keyMaxItems, "666")),
		},
		{
			name: "after",
			arg:  After(SameID("https://example.com")),
			want: vals(kv(keyAfter, "https://example.com")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlValue(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("urlValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractURLVal(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		arg  Check
		want string
	}{
		{
			name: "empty",
			arg:  nil,
			want: "",
		},
		{
			name: "nil id",
			arg:  idNil{},
			want: "",
		},
		{
			name: "nil iri",
			arg:  iriNil{},
			want: "",
		},
		{
			name: "iri",
			arg:  idEquals("http://example.com"),
			want: "http://example.com",
		},
		{
			name: "iri like",
			arg:  idLike("http://example.com"),
			want: "~http://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractURLVal(tt.arg); got != tt.want {
				t.Errorf("extractURLVal() = %q, want %q", got, tt.want)
			}
		})
	}
}
