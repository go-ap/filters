package index

import (
	"encoding"

	"github.com/RoaringBitmap/roaring/roaring64"
	vocab "github.com/go-ap/activitypub"
)

type full roaring64.Bitmap

func All() Indexable {
	return (*full)(new(roaring64.Bitmap))
}

func (f *full) Add(i vocab.LinkOrIRI) uint64 {
	bmp := (*roaring64.Bitmap)(f)
	r := HashFn(i)
	bmp.Add(r)
	return r
}

func (f *full) UnmarshalBinary(data []byte) error {
	return (*roaring64.Bitmap)(f).UnmarshalBinary(data)
}

func (f *full) MarshalBinary() (data []byte, err error) {
	return (*roaring64.Bitmap)(f).MarshalBinary()
}

var _ encoding.BinaryMarshaler = new(full)
var _ encoding.BinaryUnmarshaler = new(full)
