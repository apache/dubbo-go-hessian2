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
	d := v.(big.Decimal)
	d.Value = string(d.ToString())
	return e.encObject(d)
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
