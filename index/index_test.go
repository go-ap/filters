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

func TestIndex_Add(t *testing.T) {
	type fields struct {
		Ref     map[uint32]vocab.IRI
		Indexes map[Type]Indexable
	}
	tests := []struct {
		name    string
		fields  fields
		arg     vocab.LinkOrIRI
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			arg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Index{
				w:       sync.RWMutex{},
				Ref:     tt.fields.Ref,
				Indexes: tt.fields.Indexes,
			}
			if err := i.Add(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIndex_Find(t *testing.T) {
	type fields struct {
		Ref     map[uint32]vocab.IRI
		Indexes map[Type]Indexable
	}
	tests := []struct {
		name    string
		fields  fields
		args    []BasicFilter
		want    []vocab.IRI
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			args:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Index{
				w:       sync.RWMutex{},
				Ref:     tt.fields.Ref,
				Indexes: tt.fields.Indexes,
			}
			got, err := i.Find(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var emptyFullIndex = []byte{
	45, 255, 137, 3, 1, 1, 9, 98, 97, 114, 101, 73, 110, 100, 101, 120, 1, 255,
	138, 0, 1, 2, 1, 3, 82, 101, 102, 1, 255, 140, 0, 1, 7, 73, 110, 100, 101,
	120, 101, 115, 1, 255, 142, 0, 0, 0, 43, 255, 139, 4, 1, 1, 26, 109, 97,
	112, 91, 117, 105, 110, 116, 51, 50, 93, 97, 99, 116, 105, 118, 105, 116,
	121, 112, 117, 98, 46, 73, 82, 73, 1, 255, 140, 0, 1, 6, 1, 255, 134, 0, 0,
	10, 255, 133, 5, 1, 2, 255, 134, 0, 0, 0, 46, 255, 141, 4, 1, 1, 30, 109,
	97, 112, 91, 105, 110, 100, 101, 120, 46, 84, 121, 112, 101, 93, 105, 110,
	100, 101, 120, 46, 73, 110, 100, 101, 120, 97, 98, 108, 101, 1, 255, 142,
	0, 1, 4, 1, 16, 0, 0, 3, 255, 138, 0,
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
				t.Errorf("MarshalBinary() got = %v, want %v", got, tt.want)
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
