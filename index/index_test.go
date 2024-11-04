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

func TestIndex_Add(t *testing.T) {
	type fields struct {
		Ref     map[uint32]vocab.IRI
		Indexes map[Type]Indexer
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
		Indexes map[Type]Indexer
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
