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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestException(t *testing.T) {
	doTestException(t, "throw_throwable", "exception")
	doTestException(t, "throw_exception", "exception")
	doTestException(t, "throw_MalformedParameterizedTypeException", "MalformedParameterizedType")
	doTestException(t, "throw_MalformedParametersException", "MalformedParametersException")
	doTestException(t, "throw_TypeNotPresentException", "Type exceptiontype1 not present")
	doTestException(t, "throw_WrongMethodTypeException", "WrongMethodTypeException")
	doTestException(t, "throw_UndeclaredThrowableException", "UndeclaredThrowableException")
	doTestException(t, "throw_runtimeException", "runtimeException")
	doTestException(t, "throw_arrayStoreException", "arrayStoreException")
	doTestException(t, "throw_classCastException", "classCastException")
	doTestException(t, "throw_enumConstantNotPresentException", "test.TestThrowable$TestEnum.enumConstantNotPresentException")
	doTestException(t, "throw_illegalMonitorStateException", "illegalMonitorStateException")
	doTestException(t, "throw_illegalStateException", "illegalStateException")
	doTestException(t, "throw_IOException", "IOException")
	doTestException(t, "throw_CloneNotSupportedException", "CloneNotSupportedException")
	doTestException(t, "throw_InterruptedException", "InterruptedException")
	doTestException(t, "throw_InterruptedIOException", "InterruptedIOException")
	doTestException(t, "throw_LambdaConversionException", "LambdaConversionException")
	doTestException(t, "throw_UnmodifiableClassException", "UnmodifiableClassException")
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
	doTestException(t, "throw_SecurityException", "SecurityException")
	doTestException(t, "throw_IllegalArgumentException", "IllegalArgumentException")
	doTestException(t, "throw_IllegalThreadStateException", "IllegalThreadStateException")
	doTestException(t, "throw_NumberFormatException", "NumberFormatException")
	doTestException(t, "throw_IndexOutOfBoundsException", "IndexOutOfBoundsException")
	doTestException(t, "throw_ArrayIndexOutOfBoundsException", "ArrayIndexOutOfBoundsException")
	doTestException(t, "throw_StringIndexOutOfBoundsException", "StringIndexOutOfBoundsException")
	doTestException(t, "throw_IllegalFormatWidthException", "1000")
	doTestException(t, "throw_IllegalFormatConversionException", "7 != test.TestThrowable$TestEnum")
	doTestException(t, "throw_DuplicateFormatFlagsException", "flags=DuplicateFormatFlagsException")
	doTestException(t, "throw_MissingResourceException", "MissingResourceException")
	doTestException(t, "throw_ConcurrentModificationException", "ConcurrentModificationException")
	doTestException(t, "throw_RejectedExecutionException", "RejectedExecutionException")
	doTestException(t, "throw_CompletionException", "java.lang.Throwable: exception")
	doTestException(t, "throw_EmptyStackException", "EmptyStackException")
	doTestException(t, "throw_IllformedLocaleException", "IllformedLocaleException")
	doTestException(t, "throw_NoSuchElementException", "NoSuchElementException")
	doTestException(t, "throw_NegativeArraySizeException", "NegativeArraySizeException")
	doTestException(t, "throw_UnsupportedOperationException", "UnsupportedOperationException")
	doTestException(t, "throw_ArithmeticException", "ArithmeticException")
	doTestException(t, "throw_InputMismatchException", "InputMismatchException")
	doTestException(t, "throw_ExecutionException", "ExecutionException")
	doTestException(t, "throw_InvalidPreferencesFormatException", "InvalidPreferencesFormatException")
	doTestException(t, "throw_TimeoutException", "TimeoutException")
	doTestException(t, "throw_BackingStoreException", "BackingStoreException")
	doTestException(t, "throw_DataFormatException", "DataFormatException")
	doTestException(t, "throw_BrokenBarrierException", "BrokenBarrierException")
	doTestException(t, "throw_TooManyListenersException", "TooManyListenersException")
	doTestException(t, "throw_InvalidPropertiesFormatException", "InvalidPropertiesFormatException")
	doTestException(t, "throw_ZipException", "ZipException")
	doTestException(t, "throw_JarException", "JarException")
	doTestException(t, "throw_IllegalClassFormatException", "IllegalClassFormatException")
	doTestException(t, "throw_ReflectiveOperationException", "ReflectiveOperationException")
	doTestException(t, "throw_InvocationTargetException", "InvocationTargetException")
	doTestException(t, "throw_NoSuchMethodException", "NoSuchMethodException")
	doTestException(t, "throw_NoSuchFieldException", "NoSuchFieldException")
	doTestException(t, "throw_IllegalAccessException", "IllegalAccessException")
	doTestException(t, "throw_ClassNotFoundException", "ClassNotFoundException")
	doTestException(t, "throw_InstantiationException", "InstantiationException")
	doTestException(t, "throw_DateTimeException", "DateTimeException")
	doTestException(t, "throw_UnsupportedTemporalTypeException", "UnsupportedTemporalTypeException")
	doTestException(t, "throw_ZoneRulesException", "ZoneRulesException")
	doTestException(t, "throw_DateTimeParseException", "DateTimeParseException")
	doTestException(t, "throw_FormatterClosedException", "null")
	doTestException(t, "throw_CancellationException", "CancellationException")
	doTestException(t, "throw_UnknownFormatConversionException", "Conversion = 'UnknownFormatConversionException'")
	doTestException(t, "throw_UnknownFormatFlagsException", "Flags = UnknownFormatFlagsException")
	doTestException(t, "throw_IllegalFormatFlagsException", "Flags = 'IllegalFormatFlagsException'")
	doTestException(t, "throw_IllegalFormatPrecisionException", "1")
	doTestException(t, "throw_IllegalFormatCodePointException", "Code point = 0x1")
	doTestException(t, "throw_MissingFormatArgumentException", "Format specifier 'MissingFormatArgumentException'")
	doTestException(t, "throw_MissingFormatWidthException", "MissingFormatWidthException")
	doTestException(t, "throw_DubboGenericException", "DubboGenericException")
}

func doTestException(t *testing.T, method, content string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, content, r.(error).Error())
	})
}
