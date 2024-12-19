package filters

func FilterChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		if isFilterFn(fn) {
			c = append(c, fn)
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

func CursorChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		if isCursorFn(fn) {
			c = append(c, fn)
		}
	}
	return c
}

func MaxCountCheck(fns ...Check) Check {
	for _, fn := range fns {
		if f, ok := fn.(*counter); ok {
			return f
		}
	}
	return nil
}

func MaxCount(fns ...Check) int {
	for _, fn := range fns {
		switch ff := fn.(type) {
		case *counter:
			return ff.max
		case checkAll:
			a := []Check(ff)
			return MaxCount(a...)
		case checkAny:
			a := []Check(ff)
			return MaxCount(a...)
		}
	}
	return -1
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
