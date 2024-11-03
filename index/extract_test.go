package index

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_derefObject(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.Item
		want []vocab.IRI
	}{
		{
			name: "empty",
		},
		{
			name: "item collection",
			arg: vocab.ItemCollection{
				&vocab.Object{ID: "https://example.com"},
				vocab.IRI("https://example.com/1"),
			},
			want: vocab.IRIs{"https://example.com", "https://example.com/1"},
		},
		{
			name: "item",
			arg:  &vocab.Object{ID: "https://example.com/666"},
			want: vocab.IRIs{"https://example.com/666"},
		},
		{
			name: "iri",
			arg:  vocab.IRI("https://example.com/667"),
			want: vocab.IRIs{"https://example.com/667"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := derefObject(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("derefObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
