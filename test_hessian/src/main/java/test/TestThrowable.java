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

public class TestThrowable {
    public static Object throw_exception() {
        return new Exception("exception");
    }

    public static Object throw_throwable() {
        return new Throwable("exception");
    }

    public static Object throw_uncheckedIOException() {
        return new java.io.UncheckedIOException("uncheckedIOException", getIOException());
    }

    public static Object throw_runtimeException() {
        return new RuntimeException("runtimeException");
    }

    public static Object throw_illegalStateException() {
        return new IllegalStateException("illegalStateException");
    }

    public static Object throw_illegalMonitorStateException() {
        return new IllegalMonitorStateException("illegalMonitorStateException");
    }

    public static Object throw_enumConstantNotPresentException() {
        return new EnumConstantNotPresentException(TestEnum.class, "enumConstantNotPresentException");
    }

    public static Object throw_classCastException() {
        return new ClassCastException("classCastException");
    }

    public static Object throw_arrayStoreException() {
        return new ArrayStoreException("arrayStoreException");
    }

    public static Object throw_IOException() {
        return new ArrayStoreException("IOException");
    }

    private static java.io.IOException getIOException() {
        try {
            Class clz = Class.forName("java.io.IOException");
            return (java.io.IOException)clz.newInstance();
        } catch (java.lang.Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    enum TestEnum {
        PASS
    }
}
