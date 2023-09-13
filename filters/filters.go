package filters

import (
	"errors"

	te "github.com/vitaliy-ukiru/telebot-error-handler"
)

func All(filters ...te.Filter) te.Filter {
	return func(err error) bool {
		for _, filter := range filters {
			if !filter(err) {
				return false
			}

		}
		return true
	}
}

func Any(filters ...te.Filter) te.Filter {
	return func(err error) bool {
		for _, filter := range filters {
			if filter(err) {
				return true
			}

		}
		return false
	}
}

func Is(target error) te.Filter {
	return func(err error) bool {
		return errors.Is(err, target)
	}
}

func As[E error](f func(err E) bool) te.Filter {
	return func(err error) bool {
		var e E
		if errors.As(err, &e) {
			// if f is nil return true without additional
			// filter
			return f == nil || f(e)
		}
		return false
	}
}
