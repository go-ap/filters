package index

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/RoaringBitmap/roaring/roaring64"
	vocab "github.com/go-ap/activitypub"
	"github.com/spaolacci/murmur3"
)

type (
	Type int8

	Tokenizable interface{ ~string | uint32 | uint64 }

	Indexable interface {
		Add(vocab.LinkOrIRI) uint64
	}

	bitmaps[T Tokenizable] interface {
		get(key T) *roaring64.Bitmap
		not(key T) *roaring64.Bitmap
	}

	HashFnType                   func(vocab.LinkOrIRI) uint64
	ExtractFnType[T Tokenizable] func(vocab.LinkOrIRI) []T
)

var HashSeed uint32 = 666

func murmurHash(it vocab.LinkOrIRI) uint64 {
	if it == nil {
		return 0
	}
	h := murmur3.New64WithSeed(HashSeed)
	_, _ = h.Write([]byte(it.GetLink()))
	return h.Sum64()
}

var HashFn HashFnType = murmurHash

type tokenMap[T Tokenizable] struct {
	w               sync.RWMutex
	m               map[T]*roaring64.Bitmap
	refsExtractFn   ExtractFnType[uint64]
	tokensExtractFn ExtractFnType[T]
}

func (i *tokenMap[T]) MarshalBinary() ([]byte, error) {
	buff := bytes.Buffer{}
	err := gob.NewEncoder(&buff).Encode(i.m)
	return buff.Bytes(), err
}

func (i *tokenMap[T]) UnmarshalBinary(data []byte) error {
	if i.m == nil {
		i.m = make(map[T]*roaring64.Bitmap)
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(&i.m)
}

func (i *tokenMap[T]) Add(li vocab.LinkOrIRI) uint64 {
	if i.refsExtractFn == nil {
		i.refsExtractFn = iriRefFn
	}
	refs := i.refsExtractFn(li)
	if len(refs) == 0 {
		return 0
	}
	if i.tokensExtractFn == nil {
		return 0
	}

	tokens := i.tokensExtractFn(li)
	if len(tokens) == 0 {
		return refs[0]
	}

	i.w.Lock()
	defer i.w.Unlock()
	for _, tok := range tokens {
		if _, ok := i.m[tok]; !ok {
			i.m[tok] = roaring64.New()
		}
		i.m[tok].AddMany(refs)
	}
	return refs[0]
}

// get returns the bitmap values corresponding to the key.
func (i *tokenMap[T]) get(key T) *roaring64.Bitmap {
	i.w.RLock()
	defer i.w.RUnlock()
	b, ok := i.m[key]
	if !ok {
		return roaring64.New()
	}
	return b
}

// not returns the OR-ed bitmap values for all token maps not corresponding
// to the key.
func (i *tokenMap[T]) not(key T) *roaring64.Bitmap {
	b := roaring64.New()
	i.w.RLock()
	defer i.w.RUnlock()
	for k, v := range i.m {
		if k == key {
			continue
		}
		b.Or(v)
	}
	return b
}

// NewIndex intializes a new tokenMap index where the function to extract the references that get
// indexed is passed directly.
func NewIndex[T Tokenizable](refsExtractFn ExtractFnType[uint64], tokExtractFn ExtractFnType[T]) Indexable {
	return &tokenMap[T]{
		m:               make(map[T]*roaring64.Bitmap),
		refsExtractFn:   refsExtractFn,
		tokensExtractFn: tokExtractFn,
	}
}

// NewTokenIndex intializes a new tokenMap index where the reference extraction function returns
// the indexed reference based on the ActivityPub Object ID.
func NewTokenIndex[T Tokenizable](extractFn ExtractFnType[T]) Indexable {
	return &tokenMap[T]{
		m:               make(map[T]*roaring64.Bitmap),
		refsExtractFn:   iriRefFn,
		tokensExtractFn: extractFn,
	}
}

// GetBitmaps returns the ORing of the underlying search bitmaps corresponding to the received tokens,
// or to the reverse of the returned tokens if the neg parameter is set.
func GetBitmaps[T Tokenizable](in Indexable, tokens ...T) []*roaring64.Bitmap {
	if in == nil {
		return nil
	}
	if f, ok := in.(*full); ok {
		b := (*roaring64.Bitmap)(f).Clone()
		refs := make([]uint64, len(tokens))
		for i, tok := range tokens {
			if ref, _ := any(tok).(uint64); ref > 0 {
				refs[i] = ref
			}
		}

		if len(refs) > 0 {
			b.And(roaring64.BitmapOf(refs...))
		}
		return []*roaring64.Bitmap{b}
	}

	if bmp, ok := in.(bitmaps[T]); ok {
		getFn := bmp.get

		ors := make([]*roaring64.Bitmap, 0, len(tokens))
		for _, typ := range tokens {
			ti := getFn(typ)
			ors = append(ors, ti)
		}
		return ors
	}
	return nil
}
