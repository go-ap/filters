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
	if st == nil || len(f) == 0 {
		return
	}
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

func SQLWhere(s *Stmt, ff ...Check) error {
	getWhereClauses(s, ff...)
	return nil
}

func getWhereClauses(s *Stmt, f ...Check) {
	//addNotClauses(s, f...)
	addTypeWheres(s, f...)
	addIRIWheres(s, f...)
	addNLVWheres(s, f...)
	addInReplyToWheres(s, f...)
	addAttributedToWheres(s, f...)
	addURLWheres(s, f...)
	addContextWheres(s, f...)
}

func addNotClauses(s *Stmt, f ...Check) {
	if s == nil || len(f) == 0 {
		return
	}

	nots := make([]any, 0)
	var os *Stmt
	hasNots := false
	/*
		for _, check := range f {
			switch c := check.(type) {
			case notCrit:
				nots = append()
			}
		}
	*/
	if len(nots) > 0 {
		if len(nots) == 1 {
			s.Where("iri = ?", nots[0])
		} else {
			s.Where("iri").In(nots...)
		}
	}
	if hasNots {
		tsql := strings.TrimPrefix(s.String(), " WHERE ")
		os.Where("("+tsql+" OR iri IS NULL)", s.Args()...)
	}
}

func addIRIWheres(s *Stmt, f ...Check) {
	if s == nil || len(f) == 0 {
		return
	}

	inVal := make([]any, 0)
	likeVal := make([]any, 0)

	var os *Stmt
	andNil := false
	for _, check := range f {
		switch i := check.(type) {
		case idEquals:
			inVal = append(inVal, vocab.IRI(i))
		case iriEquals:
			inVal = append(inVal, vocab.IRI(i))
		case iriLike:
			likeVal = append(likeVal, "%"+string(i)+"%")
		case idLike:
			likeVal = append(likeVal, "%"+string(i)+"%")
		case iriNil:
			andNil = true
			os = s
			s = sqlf.New("")
		case idNil:
			andNil = true
			os = s
			s = sqlf.New("")
		}
	}

	if len(inVal) > 0 {
		if len(inVal) == 1 {
			s.Where("iri = ?", inVal[0])
		} else {
			s.Where("iri").In(inVal...)
		}
	}
	if len(likeVal) > 0 {
		if len(likeVal) == 1 {
			s.Where("iri LIKE ?", likeVal[0])
		} else {
			lors := sqlf.New("")
			for _, like := range likeVal {
				lors.Where("iri LIKE ?", like)
			}
			orsq := strings.TrimPrefix(lors.String(), " WHERE ")
			s.Where(strings.ReplaceAll(orsq, " AND ", " OR "), lors.Args()...)
		}
	}
	if andNil {
		tsql := strings.TrimPrefix(s.String(), " WHERE ")
		os.Where("("+tsql+" OR iri IS NULL)", s.Args()...)
	}
}

func addTypeWheres(s *Stmt, f ...Check) {
	if s == nil || len(f) == 0 {
		return
	}

	inVal := make([]any, 0)

	var os *Stmt
	andNil := false
	for _, check := range f {
		c, ok := check.(withTypes)
		if !ok {
			continue
		}
		for _, typ := range c {
			if typ == vocab.NilType {
				andNil = true
				os = s
				s = sqlf.New("")
				continue
			}
			inVal = append(inVal, typ)
		}
	}

	if len(inVal) > 0 {
		if len(inVal) == 1 {
			s.Where("type = ?", inVal[0])
		} else {
			s.Where("type").In(inVal...)
		}
	}
	if andNil {
		tsql := strings.TrimPrefix(s.String(), " WHERE ")
		os.Where("("+tsql+" OR type IS NULL)", s.Args()...)
	}
}

func addContextWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case contextNil:
			jsonIsNull(s, "context")
		case contextEquals:
			jsonEquals(s, "context", vocab.IRI(c))
		case contextLike:
			jsonLike(s, "context", vocab.IRI(c))
		}
	}
}

func addURLWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case urlNil:
			s.Where("url IS NULL")
		case urlEquals:
			s.Where("url = ?", vocab.IRI(c))
		case urlLike:
			s.Where("url LIKE ?", "%"+vocab.IRI(c)+"%")
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

func addInReplyToWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case inReplyToNil:
			jsonIsNull(s, "inReplyTo")
		case inReplyToEquals:
			jsonEquals(s, "inReplyTo", vocab.IRI(c))
		case inReplyToLike:
			jsonLike(s, "inReplyTo", vocab.IRI(c))
		}
	}
}

func addAttributedToWheres(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case attributedToNil:
			jsonIsNull(s, "attributedTo")
		case attributedToEquals:
			jsonEquals(s, "attributedTo", vocab.IRI(c))
		case attributedToLike:
			jsonLike(s, "attributedTo", vocab.IRI(c))
		}
	}
}

func jsonLike(s *Stmt, prop string, val fmt.Stringer) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' LIKE ?`, prop), "%"+val.String()+"%")
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') LIKE ?`, prop), "%"+val.String()+"%")
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

func jsonEquals(s *Stmt, prop string, val any) {
	isPg := stmtIsPostgres(s)
	if isPg {
		s.Where(fmt.Sprintf(`raw->>'%s' = ?`, prop), val)
	} else {
		s.Where(fmt.Sprintf(`json_extract(raw, '$.%s') = ?`, prop), val)
	}
}

func stmtIsPostgres(s *Stmt) bool {
	sc := s.Clone()
	sc.Where("t = ?", 1)
	isPg := strings.Contains(sc.String(), "$1")
	return isPg
}
