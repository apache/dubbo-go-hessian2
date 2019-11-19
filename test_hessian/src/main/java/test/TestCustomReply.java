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

package test;

import com.alibaba.com.caucho.hessian.io.Hessian2Output;
import com.caucho.hessian.test.A0;
import com.caucho.hessian.test.A1;
import test.model.DateDemo;

import java.io.OutputStream;
import java.io.Serializable;
import java.math.BigDecimal;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

public class TestCustomReply {

    private Hessian2Output output;
    private HashMap<Class<?>, String> typeMap;

    TestCustomReply(OutputStream os) {
        output = new Hessian2Output(os);

        typeMap = new HashMap<>();
        typeMap.put(Void.TYPE, "void");
        typeMap.put(Boolean.class, "boolean");
        typeMap.put(Byte.class, "byte");
        typeMap.put(Short.class, "short");
        typeMap.put(Integer.class, "int");
        typeMap.put(Long.class, "long");
        typeMap.put(Float.class, "float");
        typeMap.put(Double.class, "double");
        typeMap.put(Character.class, "char");
        typeMap.put(String.class, "string");
        typeMap.put(StringBuilder.class, "string");
        typeMap.put(Object.class, "object");
        typeMap.put(Date.class, "date");
        typeMap.put(Boolean.TYPE, "boolean");
        typeMap.put(Byte.TYPE, "byte");
        typeMap.put(Short.TYPE, "short");
        typeMap.put(Integer.TYPE, "int");
        typeMap.put(Long.TYPE, "long");
        typeMap.put(Float.TYPE, "float");
        typeMap.put(Double.TYPE, "double");
        typeMap.put(Character.TYPE, "char");
        typeMap.put(boolean[].class, "[boolean");
        typeMap.put(byte[].class, "[byte");
        typeMap.put(short[].class, "[short");
        typeMap.put(int[].class, "[int");
        typeMap.put(long[].class, "[long");
        typeMap.put(float[].class, "[float");
        typeMap.put(double[].class, "[double");
        typeMap.put(char[].class, "[char");
        typeMap.put(String[].class, "[string");
        typeMap.put(Object[].class, "[object");
        typeMap.put(Date[].class, "[date");
    }

    public void customReplyTypedFixedListHasNull() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableListHasNull() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyUntypedFixedListHasNull() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(o.length, null);
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyUntypedVariableListHasNull() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, null);
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_A0() throws Exception {
        A0[] o = new A0[]{new A0(), new A0(), null};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_A0() throws Exception {
        A0[] o = new A0[]{new A0(), new A0(), null};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, "[com.caucho.hessian.test.A0");
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_int() throws Exception {
        int[] o = new int[]{1, 2, 3};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_int() throws Exception {
        int[] o = new int[]{1, 2, 3};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_long() throws Exception {
        long[] o = new long[]{1, 2, 3};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_long() throws Exception {
        long[] o = new long[]{1, 2, 3};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_float() throws Exception {
        float[] o = new float[]{1, 2, 3};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_float() throws Exception {
        float[] o = new float[]{1, 2, 3};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_double() throws Exception {
        double[] o = new double[]{1, 2, 3};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_double() throws Exception {
        double[] o = new double[]{1, 2, 3};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_short() throws Exception {
        short[] o = new short[]{1, 2, 3};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_short() throws Exception {
        short[] o = new short[]{1, 2, 3};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_char() throws Exception {
        char[] o = new char[]{'1', '2', '3'};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(o.length, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedVariableList_char() throws Exception {
        char[] o = new char[]{'1', '2', '3'};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_boolean() throws Exception {
        boolean[] o = new boolean[]{true, false, true};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_boolean() throws Exception {
        boolean[] o = new boolean[]{true, false, true};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_date() throws Exception {
        Date[] o = new Date[]{new Date(1560864000), new Date(1560864000), new Date(1560864000)};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedVariableList_date() throws Exception {
        Date[] o = new Date[]{new Date(1560864000), new Date(1560864000), new Date(1560864000)};
        if (output.addRef(o)) {
            return;
        }
        boolean hasEnd = output.writeListBegin(-1, typeMap.get(o.getClass()));
        for (Object tmp : o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }

    public void customReplyTypedFixedList_arrays() throws Exception {
        int[][][] o = new int[][][]{{{1, 2, 3}, {4, 5, 6, 7}}, {{8, 9, 10}, {11, 12, 13, 14}}};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedList_A0arrays() throws Exception {
        A0[][][] o = new A0[][][]{{{new A0(), new A0(), new A0()}, {new A0(), new A0(), new A0(), null}}, {{new A0()}, {new A0()}}};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedList_Test() throws Exception {
        TypedListTest o = new TypedListTest();
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedList_Object() throws Exception {
        Object[] o = new Object[]{new A0()};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedDecimal() throws Exception {
        BigDecimal decimal = new BigDecimal("100.256");
        output.writeObject(decimal);
        output.flush();
    }
    public void customReplyTypedFixedDecimalMap() throws Exception {
        BigDecimal decimal = new BigDecimal("100.256");
        HashMap mapDemo = new HashMap<>();
        mapDemo.put("test_BigDecimal",decimal);
        output.writeObject(mapDemo);
        output.flush();
    }

    public void customReplyTypedFixedDateNull() throws Exception {
        DateDemo demo = new DateDemo("zhangshan", null, null);
        output.writeObject(demo);
        output.flush();
    }

    public void customReplyStringEmoji() throws Exception {
        output.writeObject(TestString.getEmojiTestString());
        output.flush();
    }

}

class TypedListTest implements Serializable {
    public A0 a;
    public A0[][] list;
    public A1[][] list1;

    TypedListTest() {
        this.a = new A0();
        this.list = new A0[][]{{new A0(), new A0()}, {new A0(), new A0()}};
        this.list1 = new A1[][]{{new A1(), new A1()}, {new A1(), new A1()}};
    }

}
