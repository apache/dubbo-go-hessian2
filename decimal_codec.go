package hessian

import (
	big "github.com/dubbogo/gost/math/big"
)

type DecimalCodec struct{}

func init() {
	RegisterPOJO(&big.Decimal{})
	SetCodec("java.math.BigDecimal", DecimalCodec{})
}

func (DecimalCodec) encObject(e *Encoder, v POJO) error {
	decimal, ok := v.(big.Decimal)
	if !ok {
		return e.encObject(v)
	}
	decimal.Value = string(decimal.ToString())
	return e.encObject(decimal)
}

func (DecimalCodec) decObject(d *Decoder) (interface{}, error) {
	dec, err := d.DecodeValue()
	if err != nil {
		return nil, err
	}
	result := dec.(*big.Decimal)
	err = result.FromString([]byte(result.Value))
	if err != nil {
		return nil, err
	}
	return result, nil
}
