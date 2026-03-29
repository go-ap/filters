package filters

import (
	"net/url"
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/google/go-cmp/cmp"
)

func TestFirstPage(t *testing.T) {
	tests := []struct {
		name string
		want pagValues
	}{
		{
			name: "empty",
			want: firstPagePaginator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstPage(); !cmp.Equal(got, tt.want) {
				t.Errorf("FirstPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextPage(t *testing.T) {
	tests := []struct {
		name string
		it   vocab.Item
		want pagValues
	}{
		{
			name: "empty",
			want: pagValues{},
		},
		{
			name: "object w/o id",
			it:   &vocab.Object{},
			want: pagValues{},
		},
		{
			name: "object w/ id",
			it:   &vocab.Object{ID: "http://example.com"},
			want: pagValues{keyAfter: []string{"http://example.com"}},
		},
		{
			name: "iri",
			it:   vocab.IRI("http://social.example.com"),
			want: pagValues{keyAfter: []string{"http://social.example.com"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextPage(tt.it); !cmp.Equal(got, tt.want) {
				t.Errorf("NextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginatorValues(t *testing.T) {
	tests := []struct {
		name string
		q    url.Values
		want pagValues
	}{
		{
			name: "nil",
			want: nil,
		},
		{
			name: "empty",
			q:    url.Values{},
			want: pagValues{},
		},
		{
			name: "after",
			q: url.Values{
				keyAfter: []string{"http://example.com"},
			},
			want: pagValues{
				keyAfter: []string{"http://example.com"},
			},
		},
		{
			name: "before",
			q: url.Values{
				keyBefore: []string{"http://example.com"},
			},
			want: pagValues{
				keyBefore: []string{"http://example.com"},
			},
		},
		{
			name: "id",
			q: url.Values{
				keyID: []string{"http://example.com"},
			},
			want: pagValues{
				keyID: []string{"http://example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PaginatorValues(tt.q); !cmp.Equal(got, tt.want) {
				t.Errorf("PaginatorValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrevPage(t *testing.T) {
	tests := []struct {
		name string
		it   vocab.Item
		want pagValues
	}{
		{
			name: "empty",
			want: pagValues{},
		},
		{
			name: "object w/o id",
			it:   &vocab.Object{},
			want: pagValues{},
		},
		{
			name: "object w/ id",
			it:   &vocab.Object{ID: "http://example.com"},
			want: pagValues{keyBefore: []string{"http://example.com"}},
		},
		{
			name: "iri",
			it:   vocab.IRI("http://social.example.com"),
			want: pagValues{keyBefore: []string{"http://social.example.com"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrevPage(tt.it); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrevPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pagValues_After(t *testing.T) {
	tests := []struct {
		name string
		p    pagValues
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "with id no after",
			p: pagValues{
				keyID: []string{"http://example.com"},
			},
			want: "",
		},
		{
			name: "with id with after",
			p: pagValues{
				keyID:    []string{"http://example.com"},
				keyAfter: []string{"http://example.com"},
			},
			want: "http://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.After(); got != tt.want {
				t.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pagValues_Before(t *testing.T) {
	tests := []struct {
		name string
		p    pagValues
		want string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "with id no before",
			p: pagValues{
				keyID: []string{"http://example.com"},
			},
			want: "",
		},
		{
			name: "with id with before",
			p: pagValues{
				keyID:     []string{"http://example.com"},
				keyBefore: []string{"http://example.com"},
			},
			want: "http://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Before(); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pagValues_Count(t *testing.T) {
	tests := []struct {
		name string
		p    pagValues
		want int
	}{
		{
			name: "empty",
			want: -1,
		},
		{
			name: "not empty, no count",
			p: pagValues{
				keyID: []string{"http://example.com"},
			},
			want: -1,
		},
		{
			name: "not empty, with count 100",
			p: pagValues{
				keyMaxItems: []string{"100"},
			},
			want: 100,
		},
		{
			name: "not empty, with count 10",
			p: pagValues{
				keyMaxItems: []string{"10"},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Count(); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_pagValues_Page(t *testing.T) {
	tests := []struct {
		name string
		p    pagValues
		want int
	}{
		{
			name: "meh",
			p:    pagValues{},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Page(); got != tt.want {
				t.Errorf("Page() = %v, want %v", got, tt.want)
			}
		})
	}
}
