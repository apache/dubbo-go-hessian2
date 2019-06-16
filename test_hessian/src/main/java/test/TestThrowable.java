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

import java.io.*;
import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;

public class TestThrowable {
    public static Object throw_exception() {
        return new Exception("exception");
    }

    public static Object throw_throwable() {
        return new Throwable("exception");
    }

    public static Object throw_uncheckedIOException() {
        return new java.io.UncheckedIOException("uncheckedIOException", new java.io.IOException("io exception"));
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

    public static Object throw_NullPointerException() {
        return new NullPointerException("nullPointerException");
    }

    public static Object throw_UncheckedIOException() {
        return new UncheckedIOException("uncheckedIOException", new IOException("IOException"));
    }

    public static Object throw_FileNotFoundException() {
        return new FileNotFoundException("fileNotFoundException");
    }

    public static Object throw_EOFException() {
        return new EOFException("EOFException");
    }

    public static Object throw_SyncFailedException() {
        return new SyncFailedException("syncFailedException");
    }

    public static Object throw_ObjectStreamException() {
        return new InvalidObjectException("objectStreamException");
    }

    public static Object throw_WriteAbortedException() {
        return new WriteAbortedException("writeAbortedException", new Exception("detail"));
    }

    public static Object throw_InvalidObjectException() {
        return new InvalidObjectException("invalidObjectException");
    }

    public static Object throw_StreamCorruptedException() {
        return new StreamCorruptedException("streamCorruptedException");
    }

    public static Object throw_InvalidClassException() {
        return new InvalidClassException("invalidClassException");
    }

    public static Object throw_OptionalDataException() throws InvocationTargetException, NoSuchMethodException, IllegalAccessException, InstantiationException {
        Constructor c1= OptionalDataException.class.getDeclaredConstructor(int.class);
        c1.setAccessible(true);
        return c1.newInstance(1);
    }

    public static Object throw_NotActiveException() {
        return new NotActiveException("notActiveException");
    }

    public static Object throw_NotSerializableException() {
        return new NotSerializableException("notSerializableException");
    }

    public static Object throw_UTFDataFormatException() {
        return new UTFDataFormatException("UTFDataFormatException");
    }

    enum TestEnum {
        PASS
    }
}
