package index

import (
	"reflect"
	"sync"
	"testing"

	vocab "github.com/go-ap/activitypub"
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
				Indexes: map[Type]Indexable{
					ByType:              TokenBitmap(ExtractType),
					ByName:              TokenBitmap(ExtractName),
					ByPreferredUsername: TokenBitmap(ExtractPreferredUsername),
					BySummary:           TokenBitmap(ExtractSummary),
					ByContent:           TokenBitmap(ExtractContent),
					ByActor:             TokenBitmap(ExtractActor),
					ByObject:            TokenBitmap(ExtractObject),
					ByRecipients:        TokenBitmap(ExtractRecipients),
					ByAttributedTo:      TokenBitmap(ExtractAttributedTo),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Full()
			if !reflect.DeepEqual(got.Ref, tt.want.Ref) {
				t.Errorf("Full() = ref %+v, want %+v", got.Ref, tt.want.Ref)
			}
			for typ, bmp := range tt.want.Indexes {
				gotBmp, ok := got.Indexes[typ]
				if !ok {
					t.Errorf("Full() = indexes for type [%v] %+v, want %+v", typ, gotBmp, bmp)
					continue
				}
			}
		})
	}
}

func TestIndex_Add(t *testing.T) {
	type fields struct {
		Ref     map[uint32]vocab.IRI
		Indexes map[Type]Indexable
	}
	tests := []struct {
		name   string
		fields fields
		arg    vocab.LinkOrIRI
	}{
		{
			name:   "empty",
			fields: fields{},
			arg:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Index{
				w:       sync.RWMutex{},
				Ref:     tt.fields.Ref,
				Indexes: tt.fields.Indexes,
			}
			i.Add(tt.arg)
		})
	}
}

var emptyFullIndex = []byte{
	0x2d, 0xff, 0x87, 0x3, 0x1, 0x1, 0x9, 0x62, 0x61, 0x72, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x1, 0xff, 0x88, 0x0,
	0x1, 0x2, 0x1, 0x3, 0x52, 0x65, 0x66, 0x1, 0xff, 0x8c, 0x0, 0x1, 0x7, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x1, 0xff, 0x8e, 0x0, 0x0, 0x0, 0x2b, 0xff, 0x8b, 0x4, 0x1, 0x1, 0x1a, 0x6d, 0x61, 0x70, 0x5b, 0x75, 0x69, 0x6e,
	0x74, 0x33, 0x32, 0x5d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x70, 0x75, 0x62, 0x2e, 0x49, 0x52, 0x49,
	0x1, 0xff, 0x8c, 0x0, 0x1, 0x6, 0x1, 0xff, 0x8a, 0x0, 0x0, 0xa, 0xff, 0x89, 0x5, 0x1, 0x2, 0xff, 0x8a, 0x0, 0x0,
	0x0, 0x2e, 0xff, 0x8d, 0x4, 0x1, 0x1, 0x1e, 0x6d, 0x61, 0x70, 0x5b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x5d, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x61, 0x62, 0x6c, 0x65, 0x1,
	0xff, 0x8e, 0x0, 0x1, 0x4, 0x1, 0x10, 0x0, 0x0, 0x3, 0xff, 0x88, 0x0,
}

func TestIndex_MarshalBinary(t *testing.T) {
	type fields struct {
		Ref     map[uint32]vocab.IRI
		Indexes map[Type]Indexable
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			want:    emptyFullIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Index{
				Ref:     tt.fields.Ref,
				Indexes: tt.fields.Indexes,
			}
			got, err := i.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBinary() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIndex_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		arg     []byte
		want    Index
		wantErr bool
	}{
		{
			name:    "empty",
			arg:     emptyFullIndex,
			want:    Index{Ref: make(map[uint32]vocab.IRI), Indexes: make(map[Type]Indexable)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Index{}
			if err := i.UnmarshalBinary(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
