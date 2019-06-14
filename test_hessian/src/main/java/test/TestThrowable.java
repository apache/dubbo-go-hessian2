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

package test;

import java.lang.invoke.WrongMethodTypeException;
import java.lang.reflect.MalformedParameterizedTypeException;
import java.lang.reflect.MalformedParametersException;
import java.lang.reflect.UndeclaredThrowableException;

public class TestThrowable {
    public static Object throw_exception() {
        return new Exception("exception");
    }

    public static Object throw_throwable() {
        return new Throwable("exception");
    }

    public static Object throw_TypeNotPresentException() {
        return new TypeNotPresentException("exceptiontype1", new Throwable("exception"));
    }

    public static Object throw_UndeclaredThrowableException() {
        return new UndeclaredThrowableException(new Throwable(), "UndeclaredThrowableException");
    }

    public static Object throw_MalformedParametersException() {
        return new MalformedParametersException("MalformedParametersException");
    }

    public static Object throw_WrongMethodTypeException() {
        return new WrongMethodTypeException("WrongMethodTypeException");
    }

    public static Object throw_MalformedParameterizedTypeException() {
        return new MalformedParameterizedTypeException();
    }
}
