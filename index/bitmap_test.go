package index

import (
	"bytes"
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

var emptyStrIndex = []byte{
	15, 255, 129, 4, 1, 2, 255, 130, 0, 1, 12, 1, 255, 128, 0, 0, 9, 127, 6, 1, 2, 255, 132, 0, 0, 0,
	4, 255, 130, 0, 0,
}

var typeIndex = []byte{
	15, 255, 129, 4, 1, 2, 255, 130, 0, 1, 12, 1, 255, 128, 0, 0, 9, 127, 6, 1, 2, 255, 132, 0, 0, 0,
	30, 255, 130, 0, 1, 6, 67, 114, 101, 97, 116, 101, 18, 58, 48, 0, 0, 1, 0, 0, 0, 40, 122, 0, 0, 16, 0, 0, 0, 163, 148,
}

var emptyIRIIndex = []byte{
	16, 255, 135, 4, 1, 2, 255, 136, 0, 1, 255, 134, 1, 255, 128, 0, 0, 10, 255, 133, 5, 1, 2, 255,
	134, 0, 0, 0, 9, 127, 6, 1, 2, 255, 132, 0, 0, 0,
	4, 255, 136, 0, 0,
}

var recipientsIndex = []byte{
	16, 255, 135, 4, 1, 2, 255, 136, 0, 1, 255, 134, 1, 255, 128, 0, 0, 10, 255, 133, 5, 1, 2, 255,
	134, 0, 0, 0, 9, 127, 6, 1, 2, 255, 132, 0, 0, 0,
	43, 255, 136, 0, 1, 19, 104, 116, 116, 112, 115, 58, 47, 47, 101, 120, 97, 109, 112, 108, 101, 46, 99, 111, 109, 18, 58, 48, 0, 0, 1, 0, 0, 0, 40,
	122, 0, 0, 16, 0, 0, 0, 163, 148,
}

var strIndex = index[string]{
	tokenMap:  make(map[string]*roaring.Bitmap),
	extractFn: extractType,
}

var iriIndex = index[vocab.IRI]{
	tokenMap:  make(map[vocab.IRI]*roaring.Bitmap),
	extractFn: extractRecipients,
}

var Ob = &vocab.Object{
	ID:   "https://example.com/666",
	Type: vocab.CreateType,
	To:   vocab.ItemCollection{vocab.IRI("https://example.com")},
}

func Test_Stringy_index_MarshalBinary(t *testing.T) {
	_ = strIndex.Add(Ob)

	type testCase[T tokener] struct {
		name    string
		i       index[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "empty",
			i:       index[string]{},
			want:    emptyStrIndex,
			wantErr: false,
		},
		{
			name:    "index with type",
			i:       strIndex,
			want:    typeIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("MarshalBinary() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IRI_index_MarshalBinary(t *testing.T) {
	_ = iriIndex.Add(Ob)

	type testCase[T tokener] struct {
		name    string
		i       index[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[vocab.IRI]{
		{
			name:    "empty",
			i:       index[vocab.IRI]{},
			want:    emptyIRIIndex,
			wantErr: false,
		},
		{
			name:    "index with recipients",
			i:       iriIndex,
			want:    recipientsIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("MarshalBinary() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Stringy_index_UnmarshalBinary(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       index[T]
		arg     []byte
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "empty",
			i:       index[string]{},
			arg:     emptyStrIndex,
			wantErr: false,
		},
		{
			name:    "index with type",
			i:       strIndex,
			arg:     typeIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalBinary(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_IRI_index_UnmarshalBinary(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       index[T]
		arg     []byte
		wantErr bool
	}
	tests := []testCase[vocab.IRI]{
		{
			name:    "empty",
			i:       index[vocab.IRI]{},
			arg:     emptyIRIIndex,
			wantErr: false,
		},
		{
			name:    "index with recipients",
			i:       iriIndex,
			arg:     recipientsIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.UnmarshalBinary(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
