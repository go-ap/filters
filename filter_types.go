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

func CursorChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		if isCursorFn(fn) {
			c = append(c, fn)
		}
	}
	return c
}

func MaxCountChecks(fns ...Check) Check {
	for _, fn := range fns {
		if f, ok := fn.(*counter); ok {
			return f
		}
	}
	return nil
}

func MaxCount(fns ...Check) int {
	for _, fn := range fns {
		if f, ok := fn.(*counter); ok {
			return f.max
		}
	}
	return -1
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
