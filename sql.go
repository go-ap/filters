package filters

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"github.com/leporo/sqlf"
)

func SQLLimit(st *Stmt, f ...Check) {
	lim := MaxItems
	for _, check := range f {
		switch c := check.(type) {
		case *counter:
			lim = c.max
			break
		}
	}
	st.Limit(lim)
}

type Stmt = sqlf.Stmt

func SQLBuild(s *Stmt, ff ...Check) error {
	getWhereClauses(s, ff...)
	return nil
}

func getWhereClauses(s *Stmt, f ...Check) {
	addTypeWheres(s, f...)
	getIRIWheres(s, f...)
	addNLVWheres(s, f...)
	addInReplyToWheres(s, f...)
	addAttributedToWheres(s, f...)
	addURLWheres(s, f...)
	addContextWheres(s, f...)
}

func getIRIWheres(s *Stmt, f ...Check) (string, []any) {
	inVal := make([]any, 0)

	for _, check := range f {
		switch i := check.(type) {
		case idEquals:
			inVal = append(inVal, i)
		case iriEquals:
			inVal = append(inVal, i)
		case iriLike:
			s.Where("iri LIKE ?", "%"+i+"%")
		case idLike:
			s.Where("iri LIKE ?", "%"+i+"%")
		case iriNil:
			s.Where("iri IS NULL")
		case idNil:
			s.Where("iri IS NULL")
		}
	}

	if len(inVal) > 0 {
		if len(inVal) == 1 {
			s.Where("iri = ?", inVal[0])
		} else {
			s.Where("iri").In(inVal...)
		}
	}
	return "", nil
}

func addTypeWheres(s *Stmt, f ...Check) {
	inVal := make([]any, 0)

	for _, check := range f {
		c, ok := check.(withTypes)
		if !ok {
			continue
		}
		for _, typ := range c {
			// TODO(marius): add support for type being NULL
			if typ != vocab.NilType {
				inVal = append(inVal, typ)
			}
		}
	}

	if len(inVal) > 0 {
		if len(inVal) == 1 {
			s.Where("type = ?", inVal[0])
		} else {
			s.Where("type").In(inVal...)
		}
	}
}

func stmtIsPostgres(s *Stmt) bool {
	sc := s.Clone()
	sc.Where("t = ?", 1)
	isPg := strings.Contains(sc.String(), "$1")
	return isPg
}

func addContextWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case contextEquals:
			jsonEquals(s, "context", string(c))
		case contextLike:
			jsonLike(s, "context", string(c))
		case contextNil:
			jsonIsNull(s, "context")
		}
	}
}

func addURLWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case urlEquals:
			s.Where("url = ?", c)
		case urlLike:
			s.Where("url LIKE ?", "%"+c+"%")
		case urlNil:
			s.Where("url IS NULL")
		}
	}
}

func sameFns(f1, f2 any) bool {
	p1 := reflect.ValueOf(f1).Pointer()
	p2 := reflect.ValueOf(f2).Pointer()
	if p1 == p2 {
		return true
	}
	if p1 == 0 || p2 == 0 {
		return false
	}
	s1, l1 := runtime.FuncForPC(p1).FileLine(p1)
	s2, l2 := runtime.FuncForPC(p2).FileLine(p2)
	return s1 == s2 && l1 == l2
}

func addNLVWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case naturalLanguageValCheck:
			var field string
			switch c.typ {
			case byName:
				field = keyName
			case byPreferredUsername:
				field = "preferred_username"
			case bySummary:
				field = keySummary
			case byContent:
				field = keyContent
			}
			switch {
			case sameFns(c.checkFn, naturalLanguageEmpty):
				s.Where(field + " IS NULL")
			case sameFns(c.checkFn, naturalLanguageValuesLike):
				s.Where(field+" LIKE ?", "%"+c.checkValue+"%")
			case sameFns(c.checkFn, naturalLanguageValuesEquals):
				s.Where(field+" = ?", c.checkValue)
			}
		}
	}
}

func jsonLike(s *Stmt, prop, val string) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' LIKE ?`, prop), "%"+val+"%")
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') LIKE ?`, prop), "%"+val+"%")
	}
}

func jsonIsNull(s *Stmt, prop string) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' IS NULL`, prop))
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') IS NULL`, prop))
	}
}

func jsonIsNotNull(s *Stmt, prop string) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' IS NOT NULL`, prop))
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') IS NOT NULL`, prop))
	}
}

func jsonEquals(s *Stmt, prop, val string) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' = ?`, prop), val)
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') = ?`, prop), val)
	}
}

func addInReplyToWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case iriNil:
			jsonIsNull(s, "inReplyTo")
		case iriEquals:
			jsonEquals(s, "inReplyTo", string(c))
		}
	}
}

func addAttributedToWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case iriNil:
			jsonIsNull(s, "attributedTo")
		case iriEquals:
			jsonEquals(s, "attributedTo", string(c))
		}
	}
}
