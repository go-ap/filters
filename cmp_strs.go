package filters

import "github.com/mariusor/qstring"

type CompStr = qstring.ComparativeString
type CompStrs []CompStr

func StringEquals(s string) CompStr {
	return CompStr{Str: s}
}
func StringLike(s string) CompStr {
	return CompStr{Operator: "~", Str: s}
}
func StringDifferent(s string) CompStr {
	return CompStr{Operator: "!", Str: s}
}

func (cs CompStrs) Contains(f CompStr) bool {
	for _, c := range cs {
		if c.Str == f.Str && c.Operator == f.Operator {
			return true
		}
	}
	return false
}
