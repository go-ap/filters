package index

import (
	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
)

type full roaring.Bitmap

func (f *full) Add(i vocab.LinkOrIRI) uint32 {
	bmp := (*roaring.Bitmap)(f)
	r := HashFn(i)
	bmp.Add(r)
	return r
}
