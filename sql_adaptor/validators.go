package sql_adaptor

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ahiho/filter/parser"
)

// DefaultMatcherWithValidator wraps the default matcher with validation on the value.
func DefaultMatcherWithValidator(validate ValidatorFunc, comps []string) ParseValidateFunc {
	return func(ex *parser.Expression) (*SqlResponse, error) {
		for _, v := range comps {
			if v == ex.Comparator || v == "*" {
				err := validate(ex.Value)
				if err != nil {
					return nil, errors.New("invalid value")
				}
				return DefaultMatcher(ex), nil
			}
		}
		return nil, errors.New("field is not allowed")
	}
}

// DefaultMatcher takes an expression and spits out the default SqlResponse.
func DefaultMatcher(ex *parser.Expression) *SqlResponse {
	if ex.Comparator == parser.TokenLookup[parser.PERCENT] {
		fmtValue := fmt.Sprintf("%%%s%%", ex.Value)
		sq := SqlResponse{
			Raw:    fmt.Sprintf("%s LIKE ?", ex.Field),
			Values: []string{fmtValue},
		}
		return &sq
	}
	sq := SqlResponse{
		Raw:    fmt.Sprintf("%s%s?", ex.Field, ex.Comparator),
		Values: []string{ex.Value},
	}
	return &sq
}

// NullValidator is a no-op validator on a string, always returns nil error.
func NullValidator(_ string) error {
	return nil
}

// IntegerValidator validates that the input is an integer.
func IntegerValidator(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("value '%s' is not an integer", s)
	}
	return nil
}

// NumericValidator validates that the input is a number.
func NumericValidator(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("value '%s' is not numeric", s)
	}
	return nil
}
