package index

import (
	"reflect"
	"testing"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
)

func Test_IRI_TokenBitmap(t *testing.T) {
	type testCase[T tokener] struct {
		name string
		arg  extractFnType[T]
		want index[T]
	}
	tests := []testCase[vocab.IRI]{
		{
			name: "empty",
		},
		{
			name: "iri attributedTo",
			arg:  extractAttributedTo,
			want: index[vocab.IRI]{tokenMap: make(map[vocab.IRI]*roaring.Bitmap), extractFn: extractAttributedTo},
		},
		{
			name: "iri Actor",
			arg:  extractActor,
			want: index[vocab.IRI]{tokenMap: make(map[vocab.IRI]*roaring.Bitmap), extractFn: extractActor},
		},
		{
			name: "iri Object",
			arg:  extractObject,
			want: index[vocab.IRI]{tokenMap: make(map[vocab.IRI]*roaring.Bitmap), extractFn: extractObject},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TokenBitmap(tt.arg).(*index[vocab.IRI])
			if !ok {
				t.Errorf("TokenBitmap() = invalid type %T for result", got)
			}
			if got.tokenMap != nil && tt.want.tokenMap != nil && !reflect.DeepEqual(got.tokenMap, tt.want.tokenMap) {
				t.Errorf("TokenBitmap() = invalid token map %+v, expected %+v", got.tokenMap, tt.want.tokenMap)
			}
			if !sameFunc(got.extractFn, tt.want.extractFn) {
				t.Errorf("TokenBitmap() = invalid extractFn %p, expected %p", got.extractFn, tt.want.extractFn)
			}
		})
	}
}

func Test_Stringy_TokenBitmap(t *testing.T) {
	type testCase[T tokener] struct {
		name string
		arg  extractFnType[T]
		want index[T]
	}
	tests := []testCase[string]{
		{
			name: "empty",
		},
		{
			name: "stringy preferred username",
			arg:  extractPreferredUsername,
			want: index[string]{tokenMap: make(map[string]*roaring.Bitmap), extractFn: extractPreferredUsername},
		},
		{
			name: "stringy name",
			arg:  extractName,
			want: index[string]{tokenMap: make(map[string]*roaring.Bitmap), extractFn: extractName},
		},
		{
			name: "stringy content",
			arg:  extractContent,
			want: index[string]{tokenMap: make(map[string]*roaring.Bitmap), extractFn: extractContent},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TokenBitmap(tt.arg).(*index[string])
			if !ok {
				t.Errorf("TokenBitmap() = invalid type %T for result", got)
			}
			if got.tokenMap != nil && tt.want.tokenMap != nil && !reflect.DeepEqual(got.tokenMap, tt.want.tokenMap) {
				t.Errorf("TokenBitmap() = invalid token map %+v, expected %+v", got.tokenMap, tt.want.tokenMap)
			}
			if !sameFunc(got.extractFn, tt.want.extractFn) {
				t.Errorf("TokenBitmap() = invalid extractFn %p, expected %p", got.extractFn, tt.want.extractFn)
			}
		})
	}
}

func sameFunc(f1, f2 any) bool {
	r1 := reflect.ValueOf(f1)
	r2 := reflect.ValueOf(f2)
	return r1.UnsafePointer() == r2.UnsafePointer()
}

func hashAll(vals ...vocab.LinkOrIRI) []uint32 {
	ints := make([]uint32, 0, len(vals))
	for _, val := range vals {
		ints = append(ints, hashFn(val))
	}
	return ints
}

func tk[T tokener](k T, vals ...vocab.LinkOrIRI) func(mm map[T]*roaring.Bitmap) {
	return func(mm map[T]*roaring.Bitmap) {
		mm[k] = roaring.BitmapOf(hashAll(vals...)...)
	}
}

func tMap[T tokener](fns ...func(map[T]*roaring.Bitmap)) map[T]*roaring.Bitmap {
	m := make(map[T]*roaring.Bitmap)
	for _, fn := range fns {
		fn(m)
	}
	return m
}

func Test_IRI_index_Add(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       index[T]
		arg     vocab.LinkOrIRI
		want    map[T]*roaring.Bitmap
		wantErr bool
	}
	tests := []testCase[vocab.IRI]{
		{
			name: "empty",
		},
		{
			name: "iri attributedTo",
			i:    index[vocab.IRI]{tokenMap: make(map[vocab.IRI]*roaring.Bitmap), extractFn: extractAttributedTo},
			arg:  &vocab.Object{ID: "https://example.com/1", AttributedTo: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(vocab.IRI("https://example.com/~jane"), vocab.IRI("https://example.com/1"))),
		},
		{
			name: "iri Actor",
			i:    index[vocab.IRI]{tokenMap: make(map[vocab.IRI]*roaring.Bitmap), extractFn: extractActor},
			arg:  &vocab.Activity{ID: "https://example.com/2", Actor: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(vocab.IRI("https://example.com/~jane"), vocab.IRI("https://example.com/2"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Add(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.i
			if got.tokenMap != nil && tt.want != nil && !reflect.DeepEqual(got.tokenMap, tt.want) {
				t.Errorf("Add() = invalid token map %+v, expected %+v", got.tokenMap, tt.want)
			}
		})
	}
}

func Test_Stringy_index_Add(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       index[T]
		arg     vocab.LinkOrIRI
		want    map[T]*roaring.Bitmap
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name: "empty",
		},
		{
			name: "type",
			i:    index[string]{tokenMap: make(map[string]*roaring.Bitmap), extractFn: extractType},
			arg:  &vocab.Object{ID: "https://example.com/1", Type: vocab.NoteType},
			want: tMap(tk("Note", vocab.IRI("https://example.com/1"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.Add(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.i
			if got.tokenMap != nil && tt.want != nil && !reflect.DeepEqual(got.tokenMap, tt.want) {
				t.Errorf("Add() = invalid token map %+v, expected %+v", got.tokenMap, tt.want)
			}
		})
	}
}

func Test_murmurHash(t *testing.T) {
	tests := []struct {
		name string
		seed uint32
		arg  vocab.LinkOrIRI
		want uint32
	}{
		{
			name: "empty",
		},
		{
			name: "http://example.com",
			arg:  vocab.IRI("http://example.com"),
			seed: 666,
			want: 209591596,
		},
		{
			name: "https://localhost:123",
			arg:  vocab.IRI("https://localhost:123"),
			seed: 666,
			want: 3666045539,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.seed > 0 {
				HashSeed = tt.seed
			}

			if got := murmurHash(tt.arg); got != tt.want {
				t.Errorf("murmurHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
