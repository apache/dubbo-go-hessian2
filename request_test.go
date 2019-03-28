// Copyright (c) 2016 ~ 2019, dubbogo.
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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPackRequest(t *testing.T) {
	bytes, err := PackRequest(Service{
		Path:      "/test",
		Interface: "ITest",
		Version:   "v1.0",
		Target:    "test",
		Method:    "test",
		Timeout:   time.Second * 10,
	}, DubboHeader{
		SerialID: 0,
		Type:     Request,
		ID:       123,
	}, []interface{}{1, 2})

	assert.Nil(t, err)
	
	if bytes != nil {
		t.Logf("pack request: %s", string(bytes))
	}
}
