package hessian

import (
	big "github.com/dubbogo/gost/math/big"
)

func init() {
	RegisterPOJO(&big.Decimal{})
}
