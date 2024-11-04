package index

type opType uint8

const (
	OPEq opType = iota
	OPLike
	OPNot
)

type BasicFilter struct {
	Values []string
	Op     opType
	Type   Type
}
