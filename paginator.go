package filters

import (
	"net/url"
	"strconv"

	vocab "github.com/go-ap/activitypub"
)

const MaxItems int = 100

var firstPagePaginator = pagValues(url.Values{
	keyMaxItems: []string{strconv.Itoa(MaxItems)},
})

// FirstPage returns the default url.Values for getting to the first page of a collection.
func FirstPage() pagValues {
	return firstPagePaginator
}

func PrevPage(it vocab.Item) pagValues {
	return pagValues{keyBefore: []string{string(it.GetLink())}}
}

func NextPage(it vocab.Item) pagValues {
	return pagValues{keyAfter: []string{string(it.GetLink())}}
}

type pagValues url.Values

// KeysetPaginator
type KeysetPaginator interface {
	Before() string
	After() string
	Count() int
}

// OffsetPaginator
type OffsetPaginator interface {
	Count() int
	Page() int
}

func PaginatorValues(q url.Values) pagValues {
	return pagValues(q)
}

func (p pagValues) Before() string {
	u := url.Values(p)
	return u.Get(keyBefore)
}

func (p pagValues) After() string {
	u := url.Values(p)
	return u.Get(keyAfter)
}

func (p pagValues) Count() int {
	u := url.Values(p)
	cnt, err := strconv.ParseInt(u.Get(keyMaxItems), 10, 32)
	if err != nil {
		return -1
	}
	return int(cnt)
}

func (p pagValues) Page() int {
	return -1
}
