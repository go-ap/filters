package index

import (
	vocab "github.com/go-ap/activitypub"
	"reflect"
	"testing"
)

func TestFull(t *testing.T) {
	tests := []struct {
		name string
		want *Index
	}{
		{
			name: "some",
			want: &Index{
				Ref: make(map[uint32]vocab.IRI),
				Indexes: map[Type]Indexer{
					ByType:              TokenBitmap(extractType),
					ByName:              TokenBitmap(extractName),
					ByPreferredUsername: TokenBitmap(extractPreferredUsername),
					BySummary:           TokenBitmap(extractSummary),
					ByContent:           TokenBitmap(extractContent),
					ByActor:             TokenBitmap(extractActor),
					ByObject:            TokenBitmap(extractObject),
					ByRecipients:        TokenBitmap(extractRecipients),
					ByAttributedTo:      TokenBitmap(extractAttributedTo),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Full()
			if !reflect.DeepEqual(got.Ref, tt.want.Ref) {
				t.Errorf("Full() = Ref %+v, want %+v", got.Ref, tt.want.Ref)
			}
			for typ, bmp := range tt.want.Indexes {
				gotBmp, ok := got.Indexes[typ]
				if !ok {
					t.Errorf("Full() = Indexes for type [%v] %+v, want %+v", typ, gotBmp, bmp)
					continue
				}
				//if !reflect.DeepEqual(bmp, gotBmp) {
				//	t.Errorf("Full() = Index[%d] %+v, want %+v", typ, gotBmp, bmp)
				//}
			}
		})
	}
}
