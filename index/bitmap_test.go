package index

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/RoaringBitmap/roaring/roaring64"
	vocab "github.com/go-ap/activitypub"
)

func Test_IRI_NewTokenIndex(t *testing.T) {
	type testCase[T Tokenizable] struct {
		name string
		arg  ExtractFnType[T]
		want tokenMap[T]
	}
	tests := []testCase[uint64]{
		{
			name: "empty",
		},
		{
			name: "iri attributedTo",
			arg:  ExtractAttributedTo,
			want: tokenMap[uint64]{m: make(map[uint64]*roaring64.Bitmap), tokensExtractFn: ExtractAttributedTo},
		},
		{
			name: "iri Actor",
			arg:  ExtractActor,
			want: tokenMap[uint64]{m: make(map[uint64]*roaring64.Bitmap), tokensExtractFn: ExtractActor},
		},
		{
			name: "iri Object",
			arg:  ExtractObject,
			want: tokenMap[uint64]{m: make(map[uint64]*roaring64.Bitmap), tokensExtractFn: ExtractObject},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NewTokenIndex(tt.arg).(*tokenMap[uint64])
			if !ok {
				t.Errorf("NewTokenIndex() = invalid type %T for result", got)
			}
			if got.m != nil && tt.want.m != nil && !reflect.DeepEqual(got.m, tt.want.m) {
				t.Errorf("NewTokenIndex() = invalid token map %+v, expected %+v", got.m, tt.want.m)
			}
			if !sameFunc(got.tokensExtractFn, tt.want.tokensExtractFn) {
				t.Errorf("NewTokenIndex() = invalid tokensExtractFn %p, expected %p", got.tokensExtractFn, tt.want.tokensExtractFn)
			}
		})
	}
}

func Test_Stringy_NewTokenIndex(t *testing.T) {
	type testCase[T Tokenizable] struct {
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
			want: tokenMap[string]{m: make(map[string]*roaring64.Bitmap), tokensExtractFn: ExtractPreferredUsername},
		},
		{
			name: "stringy name",
			arg:  ExtractName,
			want: tokenMap[string]{m: make(map[string]*roaring64.Bitmap), tokensExtractFn: ExtractName},
		},
		{
			name: "stringy content",
			arg:  ExtractContent,
			want: tokenMap[string]{m: make(map[string]*roaring64.Bitmap), tokensExtractFn: ExtractContent},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NewTokenIndex(tt.arg).(*tokenMap[string])
			if !ok {
				t.Errorf("NewTokenIndex() = invalid type %T for result", got)
			}
			if got.m != nil && tt.want.m != nil && !reflect.DeepEqual(got.m, tt.want.m) {
				t.Errorf("NewTokenIndex() = invalid token map %+v, expected %+v", got.m, tt.want.m)
			}
			if !sameFunc(got.tokensExtractFn, tt.want.tokensExtractFn) {
				t.Errorf("NewTokenIndex() = invalid tokensExtractFn %p, expected %p", got.tokensExtractFn, tt.want.tokensExtractFn)
			}
		})
	}
}

func sameFunc(f1, f2 any) bool {
	r1 := reflect.ValueOf(f1)
	r2 := reflect.ValueOf(f2)
	return r1.UnsafePointer() == r2.UnsafePointer()
}

func hashAll(vals ...vocab.LinkOrIRI) []uint64 {
	ints := make([]uint64, 0, len(vals))
	for _, val := range vals {
		ints = append(ints, HashFn(val))
	}
	return ints
}

func tk[T Tokenizable](k T, vals ...vocab.LinkOrIRI) func(mm map[T]*roaring64.Bitmap) {
	return func(mm map[T]*roaring64.Bitmap) {
		mm[k] = roaring64.BitmapOf(hashAll(vals...)...)
	}
}

func tMap[T Tokenizable](fns ...func(map[T]*roaring64.Bitmap)) map[T]*roaring64.Bitmap {
	m := make(map[T]*roaring64.Bitmap)
	for _, fn := range fns {
		fn(m)
	}
	return m
}

func getRef[T ~string](v T) uint64 {
	return HashFn(vocab.IRI(v))
}

func Test_IRI_index_Add(t *testing.T) {
	type testCase[T Tokenizable] struct {
		name    string
		i       tokenMap[T]
		arg     vocab.LinkOrIRI
		want    map[T]*roaring64.Bitmap
		wantErr bool
	}
	tests := []testCase[uint64]{
		{
			name: "empty",
		},
		{
			name: "iri attributedTo",
			i:    tokenMap[uint64]{m: make(map[uint64]*roaring64.Bitmap), tokensExtractFn: ExtractAttributedTo},
			arg:  &vocab.Object{ID: "https://example.com/1", AttributedTo: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(getRef("https://example.com/~jane"), vocab.IRI("https://example.com/1"))),
		},
		{
			name: "iri Actor",
			i:    tokenMap[uint64]{m: make(map[uint64]*roaring64.Bitmap), tokensExtractFn: ExtractActor},
			arg:  &vocab.Activity{ID: "https://example.com/2", Actor: vocab.IRI("https://example.com/~jane")},
			want: tMap(tk(getRef("https://example.com/~jane"), vocab.IRI("https://example.com/2"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.i.Add(tt.arg)
			got := tt.i
			if got.m != nil && tt.want != nil && !reflect.DeepEqual(got.m, tt.want) {
				t.Errorf("Add() = invalid token map %+v, expected %+v", got.m, tt.want)
			}
		})
	}
}

func Test_Stringy_index_Add(t *testing.T) {
	type testCase[T Tokenizable] struct {
		name    string
		i       tokenMap[T]
		arg     vocab.LinkOrIRI
		want    map[T]*roaring64.Bitmap
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name: "empty",
		},
		{
			name: "type",
			i:    tokenMap[string]{m: make(map[string]*roaring64.Bitmap), tokensExtractFn: ExtractType},
			arg:  &vocab.Object{ID: "https://example.com/1", Type: vocab.NoteType},
			want: tMap(tk("Note", vocab.IRI("https://example.com/1"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.i.Add(tt.arg)
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
		want uint64
	}{
		{
			name: "empty",
		},
		{
			name: "http://example.com",
			arg:  vocab.IRI("http://example.com"),
			seed: 666,
			want: 11096362743696666034,
		},
		{
			name: "https://localhost:123",
			arg:  vocab.IRI("https://localhost:123"),
			seed: 666,
			want: 17909177958194978167,
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
	0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x2a, 0xff, 0x82, 0x0, 0x1, 0x6, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x1e,
	0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x39, 0xf2, 0x51, 0xb2, 0x3a, 0x30, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x11,
	0x34, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0, 0x27, 0x27,
}

var emptyIRIIndex = []byte{0xf, 0xff, 0x85, 0x4, 0x1, 0x2, 0xff, 0x86, 0x0, 0x1, 0x6, 0x1, 0xff, 0x80, 0x0, 0x0, 0x9,
	0x7f, 0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x4, 0xff, 0x86, 0x0, 0x0,
}

var recipientsIndex = []byte{0xf, 0xff, 0x85, 0x4, 0x1, 0x2, 0xff, 0x86, 0x0, 0x1, 0x6, 0x1, 0xff, 0x80, 0x0, 0x0, 0x9,
	0x7f, 0x6, 0x1, 0x2, 0xff, 0x84, 0x0, 0x0, 0x0, 0x2c, 0xff, 0x86, 0x0, 0x1, 0xf8, 0x31, 0xea, 0x5c, 0x35, 0x27,
	0x2b, 0x86, 0x8f, 0x1e, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x39, 0xf2, 0x51, 0xb2, 0x3a, 0x30, 0x0, 0x0, 0x1,
	0x0, 0x0, 0x0, 0x11, 0x34, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0, 0x27, 0x27,
}

var strIndex = tokenMap[string]{
	m:               make(map[string]*roaring64.Bitmap),
	tokensExtractFn: ExtractType,
}

var iriIndex = tokenMap[uint64]{
	m:               make(map[uint64]*roaring64.Bitmap),
	tokensExtractFn: ExtractRecipients,
}

var Ob = &vocab.Object{
	ID:   "https://example.com/666",
	Type: vocab.CreateType,
	To:   vocab.ItemCollection{vocab.IRI("https://example.com")},
}

func Test_Stringy_index_MarshalBinary(t *testing.T) {
	_ = strIndex.Add(Ob)

	type testCase[T Tokenizable] struct {
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
	_ = iriIndex.Add(Ob)

	type testCase[T Tokenizable] struct {
		name    string
		i       tokenMap[T]
		want    []byte
		wantErr bool
	}
	tests := []testCase[uint64]{
		{
			name:    "empty",
			i:       tokenMap[uint64]{},
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
	type testCase[T Tokenizable] struct {
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
	type testCase[T Tokenizable] struct {
		name    string
		i       tokenMap[T]
		arg     []byte
		wantErr bool
	}
	tests := []testCase[uint64]{
		{
			name:    "empty",
			i:       tokenMap[uint64]{},
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
