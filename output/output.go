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

package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	"github.com/apache/dubbo-go-hessian2/output/testfuncs"
)

type outputFunc func() []byte

var (
	funcName = flag.String("func_name", "", "func name")
	funcMap  = make(map[string]outputFunc, 8)
)

// add all output func here
func init() {
	funcMap["HelloWorldString"] = testfuncs.HelloWorldString
	funcMap["Java8TimeYear"] = testfuncs.Java8TimeYear
	funcMap["Java8LocalDate"] = testfuncs.Java8LocalDate
	funcMap["JavaException"] = testfuncs.JavaException
	funcMap["UserArray"] = testfuncs.UserArray
	funcMap["UserList"] = testfuncs.UserList
	funcMap["ByteArray"] = testfuncs.ByteArray
	funcMap["ShortArray"] = testfuncs.ShortArray
	funcMap["IntegerArray"] = testfuncs.IntegerArray
	funcMap["LongArray"] = testfuncs.LongArray
	funcMap["BooleanArray"] = testfuncs.BooleanArray
	funcMap["CharacterArray"] = testfuncs.CharacterArray
	funcMap["FloatArray"] = testfuncs.FloatArray
	funcMap["DoubleArray"] = testfuncs.DoubleArray
}

func main() {
	flag.Parse()

	if *funcName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "func name required")
		os.Exit(1)
	}
	f, exist := funcMap[*funcName]
	if !exist {
		_, _ = fmt.Fprintln(os.Stderr, "func name not exist: ", *funcName)
		os.Exit(1)
	}

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
			os.Exit(1)
		}
	}()
	if _, err := os.Stdout.Write(f()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "call error: ", err)
		os.Exit(1)
	}
	os.Exit(0)
}
