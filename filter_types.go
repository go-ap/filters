package filters

func FilterChecks(fns ...Check) Checks {
	if len(fns) == 0 {
		return nil
	}
	c := make([]Check, 0)
	for _, fn := range fns {
		if isFilterFn(fn) {
			switch ff := fn.(type) {
			case checkAny:
				c = append(c, Any(FilterChecks(Checks(ff)...)...))
			case checkAll:
				c = append(c, All(FilterChecks(Checks(ff)...)...))
			default:
				c = append(c, fn)
			}
		}
	}
	return c
}

func ItemChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch check := fn.(type) {
		case objectChecks:
		case actorChecks:
		case targetChecks:
		case tagChecks:
		case checkAny:
			c = append(c, Any(ItemChecks(check...)...))
		case checkAll:
			c = append(c, All(ItemChecks(check...)...))
		default:
			c = append(c, check)
		}
	}
	return c
}

func PaginationChecks(fns ...Check) Checks {
	fn := func(c Check) bool {
		return isCursorFn(c) || isCounterFn(c)
	}
	return filterCheckFns(fn, fns...)
}

func CursorChecks(fns ...Check) Checks {
	return filterCheckFns(isCursorFn, fns...)
}

func filterCheckFns(checkFn func(Check) bool, fns ...Check) Checks {
	c := make([]Check, 0)

	aggCheck := func(c []Check, fil []Check) {
		for _, ff := range fil {
			if !checkFn(ff) {
				continue
			}
			c = append(c, ff)
		}
	}

	for _, fn := range fns {
		if checkFn(fn) {
			c = append(c, fn)
		} else {
			switch fil := fn.(type) {
			case checkAny:
				leftover := make([]Check, 0, len(fil))
				aggCheck(leftover, fil)
				if len(leftover) > 0 {
					c = append(c, checkAny(leftover))
				}
			case checkAll:
				leftover := make([]Check, 0, len(fil))
				aggCheck(leftover, fil)
				if len(leftover) > 0 {
					c = append(c, checkAll(leftover))
				}
			}
		}
	}
	return c
}

func MaxCountCheck(fns ...Check) Check {
	for _, fn := range fns {
		if isCounterFn(fn) {
			return fn
		}
	}
	return nil
}

func MaxCount(fns ...Check) int {
	m := -1
	for _, fn := range fns {
		switch ff := fn.(type) {
		case *counter:
			m = ff.max
		case checkAll:
			a := []Check(ff)
			m = MaxCount(a...)
		case checkAny:
			a := []Check(ff)
			m = MaxCount(a...)
		}
	}
	return m
}

func Counted(fns ...Check) int {
	for _, fn := range fns {
		if f, ok := fn.(*counter); ok {
			return f.cnt
		}
	}
	return -1
}

func AfterChecks(fns ...Check) Checks {
	for _, fn := range fns {
		if f, ok := fn.(*afterCrit); ok {
			return f.fns
		}
	}
	return nil
}

func BeforeChecks(fns ...Check) Checks {
	for _, fn := range fns {
		if f, ok := fn.(*beforeCrit); ok {
			return f.fns
		}
	}
	return nil
}

func ObjectChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case objectChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func ActorChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func TargetChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case targetChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func ActivityChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case targetChecks:
			c = append(c, fns...)
		case objectChecks:
			c = append(c, fns...)
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func IntransitiveActivityChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case targetChecks:
			c = append(c, fns...)
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func TypeChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch t := fn.(type) {
		case withTypes:
			c = append(c, t)
		}
	}
	return c
}

func TagChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch t := fn.(type) {
		case tagChecks:
			c = append(c, t)
		}
	}
	return c
}

func AuthorizedChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch t := fn.(type) {
		case authorized:
			c = append(c, t)
		}
	}
	return c
}

func RecipientsChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch t := fn.(type) {
		case recipients:
			c = append(c, t)
		}
	}
	return c
}
