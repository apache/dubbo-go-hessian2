/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_exception

// ObjectStreamException represents an exception of the same name in java
type ObjectStreamException struct {
	SerialVersionUID     int64
	DetailMessage        string
	StackTrace           []StackTraceElement
	SuppressedExceptions []Throwabler
	Cause                Throwabler
}

// NewObjectStreamException is the constructor
func NewObjectStreamException(detailMessage string) *ObjectStreamException {
	return &ObjectStreamException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

// Error output error message
func (e ObjectStreamException) Error() string {
	return e.DetailMessage
}

// JavaClassName  java fully qualified path
func (ObjectStreamException) JavaClassName() string {
	return "java.io.ObjectStreamException"
}

// equals to getStackTrace in java
func (e ObjectStreamException) GetStackTrace() []StackTraceElement {
	return e.StackTrace
}
