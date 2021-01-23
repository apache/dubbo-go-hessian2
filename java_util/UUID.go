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

package java_util

import (
	"fmt"
)

//java.util.UUID
type UUID struct {
	MostSigBits  int64 `hessian:"mostSigBits"`
	LeastSigBits int64 `hessian:"leastSigBits"`
}

func (UUID) JavaClassName() string {
	return "java.util.UUID"
}

func (UUID) Error() string {
	return "encode java.util.UUID error"
}

//ToString
//Returns a string object representing this UUID.
//The UUID string representation is as described by this BNF:
//
// UUID                   = <time_low> "-" <time_mid> "-"
//                          <time_high_and_version> "-"
//                          <variant_and_sequence> "-"
//                          <node>
// time_low               = 4*<hexOctet>
// time_mid               = 2*<hexOctet>
// time_high_and_version  = 2*<hexOctet>
// variant_and_sequence   = 2*<hexOctet>
// node                   = 6*<hexOctet>
// hexOctet               = <hexDigit><hexDigit>
// hexDigit               =
//       "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
//       | "a" | "b" | "c" | "d" | "e" | "f"
//       | "A" | "B" | "C" | "D" | "E" | "F"
//
//Returns:
//A string representation of this UUID
func (uuid UUID) ToString() string {
	uuidStr := fmt.Sprintf("%v-%v-%v-%v-%v",
		digits((uuid.MostSigBits>>32), 8),
		digits((uuid.MostSigBits>>16), 4),
		digits(uuid.MostSigBits, 4),
		digits((uuid.LeastSigBits>>48), 4),
		digits(uuid.LeastSigBits, 12),
	)
	return uuidStr
}

// digits Returns arg represented by the specified number of hex digits.
func digits(arg int64, digits int) string {
	var hi = int64(1) << (digits * 4)
	i := hi | (arg & (hi - 1))
	return fmt.Sprintf("%x", i)[1:]
}
