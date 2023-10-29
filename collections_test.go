package filters

import (
	"fmt"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func TestBefore(t *testing.T) {
	tests := []struct {
		name     string
		checkIRI vocab.IRI
		with     vocab.ItemCollection
		want     []bool
	}{
		{
			name:     "empty",
			checkIRI: "https://example.com",
			want:     []bool{true},
		},
		{
			name:     "one iri",
			checkIRI: "https://example.com",
			with:     vocab.ItemCollection{vocab.IRI("http://example.com")},
			want:     []bool{true},
		},
		{
			name:     "two iris - at the end",
			checkIRI: "https://example.com",
			with:     vocab.ItemCollection{vocab.IRI("https://example1.com"), vocab.IRI("http://example.com")},
			want:     []bool{true, true},
		},
		{
			name:     "two iris - at start",
			checkIRI: "https://example1.com",
			with:     vocab.ItemCollection{vocab.IRI("https://example1.com"), vocab.IRI("http://example.com")},
			want:     []bool{false, false},
		},
		{
			name:     "three iris - in the middle",
			checkIRI: "https://example1.com",
			with: vocab.ItemCollection{
				vocab.IRI("https://example.dev"),
				vocab.IRI("https://example1.com"),
				vocab.IRI("http://example.com"),
			},
			want: []bool{true, false, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeFn := Before(ID(tt.checkIRI))
			for i, it := range tt.with {
				t.Run(fmt.Sprintf("it(%d)", i), func(t *testing.T) {
					if got := beforeFn(it); got != tt.want[i] {
						t.Errorf("Before() = %t, want %t", got, tt.want[i])
					}
				})
			}
		})
	}
}

func TestAfter(t *testing.T) {
	tests := []struct {
		name     string
		checkIRI vocab.IRI
		with     vocab.ItemCollection
		want     []bool
	}{
		{
			name:     "empty",
			checkIRI: "https://example.com",
			want:     []bool{false},
		},
		{
			name:     "one iri",
			checkIRI: "https://example.com",
			with:     vocab.ItemCollection{vocab.IRI("http://example.com")},
			want:     []bool{false},
		},
		{
			name:     "two iris - at the end",
			checkIRI: "https://example.com",
			with:     vocab.ItemCollection{vocab.IRI("https://example1.com"), vocab.IRI("http://example.com")},
			want:     []bool{false, false},
		},
		{
			name:     "two iris - at start",
			checkIRI: "https://example1.com",
			with:     vocab.ItemCollection{vocab.IRI("https://example1.com"), vocab.IRI("http://example.com")},
			want:     []bool{false, true},
		},
		{
			name:     "three iris - in the middle",
			checkIRI: "https://example1.com",
			with: vocab.ItemCollection{
				vocab.IRI("https://example.dev"),
				vocab.IRI("https://example1.com"),
				vocab.IRI("http://example.com"),
			},
			want: []bool{false, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			afterFn := After(ID(tt.checkIRI))
			for i, it := range tt.with {
				t.Run(fmt.Sprintf("it(%d)", i), func(t *testing.T) {
					if got := afterFn(it); got != tt.want[i] {
						t.Errorf("After() = %t, want %t", got, tt.want[i])
					}
				})
			}
		})
	}
}
