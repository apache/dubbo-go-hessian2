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

func init() {
	RegisterPOJO(&Throwable{})
	RegisterPOJO(&Exception{})
	RegisterPOJO(&StackTraceElement{})
}

////////////////////////////
// Throwable interface
////////////////////////////

type Throwabler interface {
	Error() string
	JavaClassName() string
}

////////////////////////////
// Throwable
////////////////////////////

type Throwable struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwable
	StackTrace           []StackTraceElement
	Cause                *Throwable
}

func NewThrowable(detailMessage string) *Throwable {
	return &Throwable{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e Throwable) Error() string {
	return e.DetailMessage
}

func (Throwable) JavaClassName() string {
	return "java.lang.Throwable"
}

////////////////////////////
// Exception
////////////////////////////

type Exception struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Exception
	StackTrace           []StackTraceElement
	Cause                *Exception
}

func NewException(detailMessage string) *Exception {
	return &Exception{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e Exception) Error() string {
	return e.DetailMessage
}

func (Exception) JavaClassName() string {
	return "java.lang.Exception"
}

type StackTraceElement struct {
	DeclaringClass string
	MethodName     string
	FileName       string
	LineNumber     int
}

func (StackTraceElement) JavaClassName() string {
	return "java.lang.StackTraceElement"
}

type Class struct {
	Name string
}

func (Class) JavaClassName() string {
	return "java.lang.Class"
}


