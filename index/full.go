package index

import (
	"encoding"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
)

type full roaring.Bitmap

func All() Indexable {
	return (*full)(new(roaring.Bitmap))
}

func (f *full) Add(i vocab.LinkOrIRI) uint32 {
	bmp := (*roaring.Bitmap)(f)
	r := HashFn(i)
	bmp.Add(r)
	return r
}

func (f *full) UnmarshalBinary(data []byte) error {
	return (*roaring.Bitmap)(f).UnmarshalBinary(data)
}

func (f *full) MarshalBinary() (data []byte, err error) {
	return (*roaring.Bitmap)(f).MarshalBinary()
}

var _ encoding.BinaryMarshaler = new(full)
var _ encoding.BinaryUnmarshaler = new(full)
