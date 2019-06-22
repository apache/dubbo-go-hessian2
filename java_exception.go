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

import "github.com/dubbogo/hessian2/java_exception"

////////////////////////////
// Throwable interface
////////////////////////////

type Throwabler interface {
	Error() string
	JavaClassName() string
}

func init() {
	RegisterPOJO(&java_exception.Class{})
	RegisterPOJO(&java_exception.Throwable{})
	RegisterPOJO(&java_exception.Exception{})
	RegisterPOJO(&java_exception.IOException{})
	RegisterPOJO(&java_exception.RuntimeException{})
	RegisterPOJO(&java_exception.StackTraceElement{})
	RegisterPOJO(&java_exception.ClassCastException{})
	RegisterPOJO(&java_exception.ArrayStoreException{})
	RegisterPOJO(&java_exception.IllegalStateException{})
	RegisterPOJO(&java_exception.IllegalMonitorStateException{})
	RegisterPOJO(&java_exception.EnumConstantNotPresentException{})
	RegisterPOJO(&java_exception.MalformedParameterizedTypeException{})
	RegisterPOJO(&java_exception.MalformedParametersException{})
	RegisterPOJO(&java_exception.TypeNotPresentException{})
	RegisterPOJO(&java_exception.UndeclaredThrowableException{})
	RegisterPOJO(&java_exception.WrongMethodTypeException{})
	RegisterPOJO(&java_exception.NullPointerException{})
	RegisterPOJO(&java_exception.UncheckedIOException{})
	RegisterPOJO(&java_exception.FileNotFoundException{})
	RegisterPOJO(&java_exception.EOFException{})
	RegisterPOJO(&java_exception.SyncFailedException{})
	RegisterPOJO(&java_exception.ObjectStreamException{})
	RegisterPOJO(&java_exception.WriteAbortedException{})
	RegisterPOJO(&java_exception.InvalidObjectException{})
	RegisterPOJO(&java_exception.StreamCorruptedException{})
	RegisterPOJO(&java_exception.InvalidClassException{})
	RegisterPOJO(&java_exception.OptionalDataException{})
	RegisterPOJO(&java_exception.NotActiveException{})
	RegisterPOJO(&java_exception.NotSerializableException{})
	RegisterPOJO(&java_exception.UTFDataFormatException{})
	RegisterPOJO(&java_exception.IllegalArgumentException{})
}
