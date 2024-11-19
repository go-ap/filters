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
		arg  ExtractFnType[T]
		want tokenMap[T]
	}
	tests := []testCase[vocab.IRI]{
		{
			name: "empty",
		},
		{
			name: "iri attributedTo",
			arg:  ExtractAttributedTo,
			want: tokenMap[vocab.IRI]{m: make(map[vocab.IRI]*roaring.Bitmap), extractFn: ExtractAttributedTo},
		},
		{
			name: "iri Actor",
			arg:  ExtractActor,
			want: tokenMap[vocab.IRI]{m: make(map[vocab.IRI]*roaring.Bitmap), extractFn: ExtractActor},
		},
		{
			name: "iri Object",
			arg:  ExtractObject,
			want: tokenMap[vocab.IRI]{m: make(map[vocab.IRI]*roaring.Bitmap), extractFn: ExtractObject},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TokenBitmap(tt.arg).(*tokenMap[vocab.IRI])
			if !ok {
				t.Errorf("TokenBitmap() = invalid type %T for result", got)
			}
			if got.m != nil && tt.want.m != nil && !reflect.DeepEqual(got.m, tt.want.m) {
				t.Errorf("TokenBitmap() = invalid token map %+v, expected %+v", got.m, tt.want.m)
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
		arg  ExtractFnType[T]
		want tokenMap[T]
	}
	tests := []testCase[string]{
		{
			name: "empty",
		},
		{
			name: "stringy preferred username",
			arg:  ExtractPreferredUsername,
			want: tokenMap[string]{m: make(map[string]*roaring.Bitmap), extractFn: ExtractPreferredUsername},
		},
		{
			name: "stringy name",
			arg:  ExtractName,
			want: tokenMap[string]{m: make(map[string]*roaring.Bitmap), extractFn: ExtractName},
		},
		{
			name: "stringy content",
			arg:  ExtractContent,
			want: tokenMap[string]{m: make(map[string]*roaring.Bitmap), extractFn: ExtractContent},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TokenBitmap(tt.arg).(*tokenMap[string])
			if !ok {
				t.Errorf("TokenBitmap() = invalid type %T for result", got)
			}
			if got.m != nil && tt.want.m != nil && !reflect.DeepEqual(got.m, tt.want.m) {
				t.Errorf("TokenBitmap() = invalid token map %+v, expected %+v", got.m, tt.want.m)
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
		ints = append(ints, HashFn(val))
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
		i       tokenMap[T]
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
			i:    tokenMap[vocab.IRI]{m: make(map[vocab.IRI]*roaring.Bitmap), extractFn: ExtractAttributedTo},
			arg:  &vocab.Object{ID: "https://example.com/1", AttributedTo: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(vocab.IRI("https://example.com/~jane"), vocab.IRI("https://example.com/1"))),
		},
		{
			name: "iri Actor",
			i:    tokenMap[vocab.IRI]{m: make(map[vocab.IRI]*roaring.Bitmap), extractFn: ExtractActor},
			arg:  &vocab.Activity{ID: "https://example.com/2", Actor: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(vocab.IRI("https://example.com/~jane"), vocab.IRI("https://example.com/2"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.i.Add(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.i
			if got.m != nil && tt.want != nil && !reflect.DeepEqual(got.m, tt.want) {
				t.Errorf("Add() = invalid token map %+v, expected %+v", got.m, tt.want)
			}
		})
	}
}

func Test_Stringy_index_Add(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       tokenMap[T]
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
			i:    tokenMap[string]{m: make(map[string]*roaring.Bitmap), extractFn: ExtractType},
			arg:  &vocab.Object{ID: "https://example.com/1", Type: vocab.NoteType},
			want: tMap(tk("Note", vocab.IRI("https://example.com/1"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.i.Add(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.i
			if got.m != nil && tt.want != nil && !reflect.DeepEqual(got.m, tt.want) {
				t.Errorf("Add() = invalid token map %+v, expected %+v", got.m, tt.want)
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

var emptyStrIndex = []byte{0xf, 0xff, 0x81, 0x4, 0x1, 0x2, 0xff, 0x82, 0x0, 0x1, 0xc, 0x1, 0xff, 0x80, 0x0, 0x0, 0x9,
	0x7f, 0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x4, 0xff, 0x82, 0x0, 0x0,
}

var typeIndex = []byte{0xf, 0xff, 0x81, 0x4, 0x1, 0x2, 0xff, 0x82, 0x0, 0x1, 0xc, 0x1, 0xff, 0x80, 0x0, 0x0, 0x9, 0x7f,
	0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x1e, 0xff, 0x82, 0x0, 0x1, 0x6, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x12, 0x3a, 0x30, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x28, 0x7a, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0, 0xa3, 0x94,
}

var emptyIRIIndex = []byte{0x10, 0xff, 0x87, 0x4, 0x1, 0x2, 0xff, 0x88, 0x0, 0x1, 0xff, 0x86, 0x1, 0xff, 0x80, 0x0, 0x0,
	0xa, 0xff, 0x85, 0x5, 0x1, 0x2, 0xff, 0x86, 0x0, 0x0, 0x0, 0x9, 0x7f, 0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x4,
	0xff, 0x88, 0x0, 0x0,
}

var recipientsIndex = []byte{0x10, 0xff, 0x87, 0x4, 0x1, 0x2, 0xff, 0x88, 0x0, 0x1, 0xff, 0x86, 0x1, 0xff, 0x80, 0x0,
	0x0, 0xa, 0xff, 0x85, 0x5, 0x1, 0x2, 0xff, 0x86, 0x0, 0x0, 0x0, 0x9, 0x7f, 0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0,
	0x2b, 0xff, 0x88, 0x0, 0x1, 0x13, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x12, 0x3a, 0x30, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x28, 0x7a, 0x0, 0x0, 0x10,
	0x0, 0x0, 0x0, 0xa3, 0x94,
}

var strIndex = tokenMap[string]{
	m:         make(map[string]*roaring.Bitmap),
	extractFn: ExtractType,
}

var iriIndex = tokenMap[vocab.IRI]{
	m:         make(map[vocab.IRI]*roaring.Bitmap),
	extractFn: ExtractRecipients,
}

var Ob = &vocab.Object{
	ID:   "https://example.com/666",
	Type: vocab.CreateType,
	To:   vocab.ItemCollection{vocab.IRI("https://example.com")},
}

func Test_Stringy_index_MarshalBinary(t *testing.T) {
	_, _ = strIndex.Add(Ob)

	type testCase[T tokener] struct {
		name    string
		i       tokenMap[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "empty",
			i:       tokenMap[string]{},
			want:    emptyStrIndex,
			wantErr: false,
		},
		{
			name:    "tokenMap with type",
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
				t.Errorf("MarshalBinary() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_IRI_index_MarshalBinary(t *testing.T) {
	_, _ = iriIndex.Add(Ob)

	type testCase[T tokener] struct {
		name    string
		i       tokenMap[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[vocab.IRI]{
		{
			name:    "empty",
			i:       tokenMap[vocab.IRI]{},
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
				t.Errorf("MarshalBinary() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_Stringy_index_UnmarshalBinary(t *testing.T) {
	type testCase[T tokener] struct {
		name    string
		i       tokenMap[T]
		arg     []byte
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "empty",
			i:       tokenMap[string]{},
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
		i       tokenMap[T]
		arg     []byte
		wantErr bool
	}
	tests := []testCase[vocab.IRI]{
		{
			name:    "empty",
			i:       tokenMap[vocab.IRI]{},
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
