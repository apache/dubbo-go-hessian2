// Copyright 2016-2019 Alex Stocks, Xinge Gao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hessian

import (
	big "github.com/dubbogo/gost/math/big"
)

type DecimalCodec struct{}

func init() {
	RegisterPOJO(&big.Decimal{})
	SetSerializer("java.math.BigDecimal", DecimalCodec{})
}

func (DecimalCodec) serializeObject(e *Encoder, v POJO) error {
	decimal, ok := v.(big.Decimal)
	if !ok {
		return e.encObject(v)
	}
	decimal.Value = string(decimal.ToString())
	return e.encObject(decimal)
}

func (DecimalCodec) deserializeObject(d *Decoder) (interface{}, error) {
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
