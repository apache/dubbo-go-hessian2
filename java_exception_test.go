// Copyright 2016-2019 Yincheng Fang
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestException(t *testing.T) {
	doTestException(t, "throw_throwable", "exception")
	doTestException(t, "throw_exception", "exception")
	doTestException(t, "throw_runtimeException", "runtimeException")
	doTestException(t, "throw_arrayStoreException", "arrayStoreException")
	doTestException(t, "throw_classCastException", "classCastException")
	doTestException(t, "throw_enumConstantNotPresentException", "test.TestThrowable$TestEnum.enumConstantNotPresentException")
	doTestException(t, "throw_illegalMonitorStateException", "illegalMonitorStateException")
	doTestException(t, "throw_illegalStateException", "illegalStateException")
	doTestException(t, "throw_IOException", "IOException")
	doTestException(t, "throw_NullPointerException", "nullPointerException")
	doTestException(t, "throw_UncheckedIOException", "uncheckedIOException")
	doTestException(t, "throw_FileNotFoundException", "fileNotFoundException")
	doTestException(t, "throw_EOFException", "EOFException")
	doTestException(t, "throw_SyncFailedException", "syncFailedException")
	doTestException(t, "throw_ObjectStreamException", "objectStreamException")
	doTestException(t, "throw_WriteAbortedException", "writeAbortedException")
	doTestException(t, "throw_InvalidObjectException", "invalidObjectException")
	doTestException(t, "throw_StreamCorruptedException", "streamCorruptedException")
	doTestException(t, "throw_InvalidClassException", "null; invalidClassException")
	doTestException(t, "throw_OptionalDataException", "null")
	doTestException(t, "throw_NotActiveException", "notActiveException")
	doTestException(t, "throw_NotSerializableException", "notSerializableException")
	doTestException(t, "throw_UTFDataFormatException", "UTFDataFormatException")
}

func doTestException(t *testing.T, method, content string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		assert.Equal(t, content, r.(error).Error())
	})
}
