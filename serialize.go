// Copyright 2016-2019 aliiohs
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

func init() {
	RegisterPOJO(&big.Decimal{})
}

type Serializer interface {
	Serialize(*Encoder, POJO) error
	Deserialize(*Decoder) (interface{}, error)
}

var SerializerMap = make(map[string]Serializer, 16)

func SetSerializer(key string, codec Serializer) {
	SerializerMap[key] = codec
}

func GetSerializer(key string) (Serializer, bool) {
	codec, ok := SerializerMap[key]
	return codec, ok
}
