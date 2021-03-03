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

package testfuncs

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-hessian2/java8_time"
)

// Java8TimeYear is test java8 java.time.Year
func Java8TimeYear() []byte {
	e := hessian.NewEncoder()
	year := java8_time.Year{Year: 2020}
	e.Encode(year)
	return e.Buffer()
}

// Java8LocalDate is test java8 java.time.LocalDate
func Java8LocalDate() []byte {
	e := hessian.NewEncoder()
	date := java8_time.LocalDate{Year: 2020, Month: 9, Day: 12}
	e.Encode(date)
	return e.Buffer()
}
