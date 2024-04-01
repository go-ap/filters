package filters

import (
	vocab "github.com/go-ap/activitypub"
	"reflect"
	"testing"
)

func TestIRILike(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Check
	}{
		{
			name: "empty",
			want: iriLike(""),
		},
		{
			name: "example.com",
			arg:  "https://example.com",
			want: iriLike("https://example.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IRILike(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IRILike() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSameIRI(t *testing.T) {
	tests := []struct {
		name string
		iri  vocab.IRI
		want Check
	}{
		{
			name: "empty",
			want: iriEquals(""),
		},
		{
			name: "example.com",
			iri:  "https://example.com",
			want: iriEquals("https://example.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameIRI(tt.iri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SameIRI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iriEquals_Apply(t *testing.T) {
	tests := []struct {
		name string
		i    iriEquals
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "example.com equals IRI",
			it:   vocab.IRI("https://example.com"),
			i:    iriEquals("https://example.com"),
			want: true,
		},
		{
			name: "example.com not equals IRI",
			it:   vocab.IRI("https://example.com/one"),
			i:    iriEquals("https://example.com"),
			want: false,
		},
		{
			name: "example.com not equals IRIs",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			i:    iriEquals("https://example.com"),
			want: false,
		},
		{
			name: "example.com not equals Items{IRI}",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			i:    iriEquals("https://example.com"),
			want: false,
		},
		{
			name: "example.com equals Object",
			it:   &vocab.Object{ID: "https://example.com"},
			i:    iriEquals("https://example.com"),
			want: true,
		},
		{
			name: "example.com not equals Object",
			it:   &vocab.Object{ID: "https://example.com/one"},
			i:    iriEquals("https://example.com"),
			want: false,
		},
		{
			name: "example.com not equals Items{Object}",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			i:    iriEquals("https://example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Apply(tt.it); got != tt.want {
				t.Errorf("Apply(%v // %v) = %v, want %v", tt.i, tt.it, got, tt.want)
			}
		})
	}
}

func Test_iriLike_Apply(t *testing.T) {
	tests := []struct {
		name string
		i    iriLike
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
		},
		{
			name: "example.com equals IRI",
			it:   vocab.IRI("https://example.com"),
			i:    iriLike("https://example.com"),
			want: true,
		},
		{
			name: "example.com like IRI",
			it:   vocab.IRI("https://example.com/one"),
			i:    iriLike("https://example.com"),
			want: true,
		},
		{
			name: "example.com not like IRI",
			it:   vocab.IRI("https://not.example.com"),
			i:    iriLike("https://example.com"),
			want: false,
		},
		{
			name: "example.com not like IRIs",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			i:    iriLike("https://example.com"),
			want: false,
		},
		{
			name: "example.com not like Items{IRI}",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			i:    iriLike("https://example.com"),
			want: false,
		},
		{
			name: "example.com equals Object",
			it:   &vocab.Object{ID: "https://example.com"},
			i:    iriLike("https://example.com"),
			want: true,
		},
		{
			name: "example.com like Object",
			it:   &vocab.Object{ID: "https://example.com/one"},
			i:    iriLike("https://example.com"),
			want: true,
		},
		{
			name: "example.com not like Object",
			it:   &vocab.Object{ID: "https://not.example.com"},
			i:    iriLike("https://example.com"),
			want: false,
		},
		{
			name: "example.com not like Items{Object}",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			i:    iriLike("https://example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Apply(tt.it); got != tt.want {
				t.Errorf("Apply(%v // %v) = %v, want %v", tt.i, tt.it, got, tt.want)
			}
		})
	}
}

func Test_iriNil_Apply(t *testing.T) {
	tests := []struct {
		name string
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
			name: "not nil IRI",
			it:   vocab.IRI("https://example.com"),
			want: false,
		},
		{
			name: "not nil IRIs",
			it:   vocab.IRIs{vocab.IRI("https://example.com/one")},
			want: false,
		},
		{
			name: "not nil Items{IRI}",
			it:   vocab.ItemCollection{vocab.IRI("https://example.com")},
			want: false,
		},
		{
			name: "not nil Object",
			it:   &vocab.Object{ID: "https://example.com"},
			want: false,
		},
		{
			name: "not nil Object",
			it:   &vocab.Object{ID: "https://example.com/one"},
			want: false,
		},
		{
			name: "not nil Items{Object}",
			it:   vocab.ItemCollection{&vocab.Object{ID: "https://example.com"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := iriNil{}
			if got := n.Apply(tt.it); got != tt.want {
				t.Errorf("iriNil.Apply(%v) = %v, want %v", tt.it, got, tt.want)
			}
		})
	}
}
