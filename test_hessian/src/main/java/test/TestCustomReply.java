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
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.caucho.hessian.test.A0;
import com.caucho.hessian.test.A1;
import test.generic.BusinessData;
import test.generic.Response;
import test.model.CustomMap;
import test.model.DateDemo;
import test.model.JavaLangObjectHolder;
import test.model.User;

import java.io.OutputStream;
import java.io.Serializable;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Date;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Set;
import java.util.UUID;


public class TestCustomReply {

    private Hessian2Output output;
    private HashMap<Class<?>, String> typeMap;

    public TestCustomReply(OutputStream os) {
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

    public void customReplyJsonString() throws Exception {
        String s = "{\"params\":{\"fromAccid\":\"23495382\",\"msgType\":100,\"msgId\":\"148ef1b2-808d-48f2-b268-7a1018a27bdb\",\"attach\":\"{\\\"accid\\\":\\\"23495382\\\",\\\"classRoomFlag\\\":50685,\\\"msgId\\\":\\\"599645021431398400\\\",\\\"msgType\\\":\\\"100\\\",\\\"nickname\\\":\\\"橙子������\\\"}\",\"roomid\":413256699},\"url\":\"https://api.netease.im/nimserver/chatroom/sendMsg.action\"}";
        output.writeObject(s);
        output.flush();
    }

    public void customReplyTypedIntegerHasNull() throws Exception {
        User user = new User();
        user.setId(null);
        output.writeObject(user);
        output.flush();
    }

    public void customReplyTypedListIntegerHasNull() throws Exception {
        User user = new User();
        user.setId(null);
        List<Integer> list = new ArrayList<>();
        list.add(1);
        list.add(null);
        user.setList(list);
        output.writeObject(user);
        output.flush();
    }

    public void customReplyTypedFixedListHasNull() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedListRefSelf() throws Exception {
        Object[] o = new Object[]{new A0(), new A1(), null};
        o[2] = o;
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

    public void customReplySingleShort() throws Exception {
        Short o = Short.valueOf((short) 123);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaShortArray() throws Exception {
        Short[] arr = new Short[]{Short.valueOf((short) 123), null, Short.valueOf((short) -456)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleInteger() throws Exception {
        Integer o = new Integer(123);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaIntegerArray() throws Exception {
        Integer[] arr = new Integer[]{new Integer(123), null, new Integer(-456)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleLong() throws Exception {
        Long o = new Long(12345);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaLongArray() throws Exception {
        Long[] arr = new Long[]{new Long(12345), null, new Long(-67890)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleBoolean() throws Exception {
        Boolean o = new Boolean(true);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaBooleanArray() throws Exception {
        Boolean[] arr = new Boolean[]{new Boolean(true), null, new Boolean(false)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleByte() throws Exception {
        Byte o = Byte.valueOf((byte) 'A');
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaByteArray() throws Exception {
        Byte[] arr = new Byte[]{Byte.valueOf((byte) 'A'), null, Byte.valueOf((byte) 'C')};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleFloat() throws Exception {
        Float o = Float.valueOf((float) 1.23);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaFloatArray() throws Exception {
        Float[] arr = new Float[]{Float.valueOf((float) 1.23), null, Float.valueOf((float) 4.56)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleDouble() throws Exception {
        Double o = Double.valueOf(1.23);
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaDoubleArray() throws Exception {
        Double[] arr = new Double[]{Double.valueOf(1.23), null, Double.valueOf(4.56)};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplySingleCharacter() throws Exception {
        Character o = new Character('A');
        output.writeObject(o);
        output.flush();
    }

    public void customReplyJavaCharacterArray() throws Exception {
        Character[] arr = new Character[]{new Character('A'), null, new Character('C')};
        output.writeObject(arr);
        output.flush();
    }

    public void customReplyTypedFixedList_Object() throws Exception {
        Object[] o = new Object[]{new A0()};
        output.writeObject(o);
        output.flush();
    }

    public void customReplyTypedFixedInteger() throws Exception {
        BigInteger integer = new BigInteger("4294967298");
        output.writeObject(integer);
        output.flush();
    }

    public void customReplyTypedFixedList_BigInteger() throws Exception {
        BigInteger[] integers = new BigInteger[]{
                new BigInteger("1234"),
                new BigInteger("12347890"),
                new BigInteger("123478901234"),
                new BigInteger("1234789012345678"),
                new BigInteger("123478901234567890"),
                new BigInteger("1234789012345678901234"),
                new BigInteger("12347890123456789012345678"),
                new BigInteger("123478901234567890123456781234"),
                new BigInteger("1234789012345678901234567812345678"),
                new BigInteger("12347890123456789012345678123456781234"),
                new BigInteger("-12347890123456789012345678123456781234"),
                new BigInteger("0"),
        };
        output.writeObject(integers);
        output.flush();
    }

    public void customReplyTypedFixedList_CustomObject() throws Exception {
        Object[] objects = new Object[]{
                new BigInteger("1234"),
                new BigInteger("-12347890"),
                new BigInteger("0"),
                new BigDecimal("123.4"),
                new BigDecimal("-123.45"),
                new BigDecimal("0"),
        };
        output.writeObject(objects);
        output.flush();
    }

    public void customReplyTypedFixedIntegerZero() throws Exception {
        BigInteger integer = new BigInteger("0");
        output.writeObject(integer);
        output.flush();
    }

    public void customReplyTypedFixedIntegerSigned() throws Exception {
        BigInteger integer = new BigInteger("-4294967298");
        output.writeObject(integer);
        output.flush();
    }

    public void customReplyTypedFixedDecimal() throws Exception {
        BigDecimal decimal = new BigDecimal("100.256");
        output.writeObject(decimal);
        output.flush();
    }

    public void customReplyTypedFixedList_BigDecimal() throws Exception {
        BigDecimal[] decimals = new BigDecimal[]{
                new BigDecimal("123.4"),
                new BigDecimal("123.45"),
                new BigDecimal("123.456"),
        };
        output.writeObject(decimals);
        output.flush();
    }

    public void customReplyObjectJsonObjectBigDecimal() throws Exception {
        JSONObject t = new JSONObject();
        BigDecimal decimal = new BigDecimal("100");
        t.put("test_BigDecimal", decimal);
        output.writeObject(t);
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

    public void customReplyStringEmoji2() throws Exception {
        output.writeObject(TestString.getEmojiTestString2());
        output.flush();
    }

    public void customReplyPerson183() throws Exception {
        Person183 p = new Person183();
        p.name = "pname";
        p.age = 13;
        InnerPerson innerPerson = new InnerPerson();
        innerPerson.name = "pname2";
        innerPerson.age = 132;
        p.innerPerson = innerPerson;
        output.writeObject(p);
        output.flush();
    }

    public void customReplyComplexString() throws Exception {
        output.writeObject(TestString.getComplexString());
        output.flush();
    }

    public void customReplyExtendClass() throws Exception {
        Dog dog = new Dog();
        dog.name = "a dog";
        dog.gender = "male";
        output.writeObject(dog);
        output.flush();
    }

    public void customReplyExtendClassToSingleStruct() throws Exception {
        Dog dog = new DogAll();
        dog.name = "a dog";
        dog.gender = "male";
        output.writeObject(dog);
        output.flush();
    }

    public void customReplyTypedFixedList_HashSet() throws Exception {
        Set<Integer> set = new HashSet<>();
        set.add(0);
        set.add(1);
        output.writeObject(set);
        output.flush();
    }

    public void customReplyTypedFixedList_HashSetCustomObject() throws Exception {
        Set<Object> set = new HashSet<>();
        set.add(new BigInteger("1234"));
        set.add(new BigDecimal("123.4"));
        output.writeObject(set);
        output.flush();
    }

    public void customReplyMap() throws Exception {
        Map<String, Object> map = new HashMap<String, Object>(4);
        map.put("a", 1);
        map.put("b", 2);
        output.writeObject(map);
        output.flush();
    }

    public void customReplyMapRefMap() throws Exception {
        Map<String, Object> map = new HashMap<String, Object>(4);
        map.put("a", 1);
        map.put("b", 2);
        map.put("self", map);

        output.writeObject(map);
        output.flush();
    }

    public void customReplyMultipleTypeMap() throws Exception {
        Map<String, Integer> map1 = new HashMap<String, Integer>(4);
        map1.put("a", 1);
        map1.put("b", 2);
        Map<Long, String> map2 = new HashMap<Long, String>(4);
        map2.put(3L, "c");
        map2.put(4L, "d");
        Map<Integer, BigDecimal> map3 = new HashMap<Integer, BigDecimal>(4);
        map3.put(5, new BigDecimal("55.55"));
        map3.put(3, new BigDecimal("33.33"));
        Map<String, Object> map = new HashMap<String, Object>(4);
        map.put("m1", map1);
        map.put("m2", map2);
        map.put("m3", map3);

        output.writeObject(map);
        output.flush();
    }

    public void customReplyListMapListMap() throws Exception {
        List<Object> list = new ArrayList<>();

        Map<String, Object> listMap1 = new HashMap<String, Object>(4);
        listMap1.put("a", 1);
        listMap1.put("b", 2);

        List<Object> items = new ArrayList<>();
        items.add(new BigDecimal("55.55"));
        items.add("hello");
        items.add(123);

        CustomMap<String, Object> innerMap = new CustomMap<String, Object>();
        innerMap.put("Int", 456);
        innerMap.put("S", "string");
        items.add(innerMap);

        listMap1.put("items", items);

        list.add(listMap1);

        CustomMap<String, Object> listMap2 = new CustomMap<String, Object>();
        listMap2.put("Int", 789);
        listMap2.put("S", "string2");

        list.add(listMap2);

        output.writeObject(list);
        output.flush();
    }

    public Map<String, Object> mapInMap() throws Exception {
        Map<String, Object> map1 = new HashMap<String, Object>();
        map1.put("a", 1);
        Map<String, Object> map2 = new HashMap<String, Object>();
        map2.put("b", 2);

        Map<String, Object> map = new HashMap<String, Object>();
        map.put("obj1", map1);
        map.put("obj2", map2);
        return map;
    }

    public void customReplyMapInMap() throws Exception {
        output.writeObject(mapInMap());
        output.flush();
    }

    public void customReplyMapInMapJsonObject() throws Exception {
        JSONObject json = JSON.parseObject(JSON.toJSONString(mapInMap()));
        output.writeObject(json);
        output.flush();
    }

    public void customReplyGenericResponseLong() throws Exception {
        Response<Long> response = new Response<>(200, 123L);
        output.writeObject(response);
        output.flush();
    }

    public void customReplyGenericResponseBusinessData() throws Exception {
        Response<BusinessData> response = new Response<>(201, new BusinessData("apple", 5));
        output.writeObject(response);
        output.flush();
    }

    public void customReplyGenericResponseList() throws Exception {
        List<BusinessData> list = new ArrayList<>();
        list.add(new BusinessData("apple", 5));
        list.add(new BusinessData("banana", 6));
        Response<List<BusinessData>> response = new Response<>(202, list);
        output.writeObject(response);
        output.flush();
    }

    public void customReplyUUID() throws Exception {
        Map<String, Object> map = new HashMap<String, Object>();
        UUID uuid1 = new UUID(459021424248441700L, -7160773830801198154L);
        UUID uuid2 = UUID.randomUUID();
        map.put("uuid1", uuid1);
        map.put("uuid1_string", uuid1.toString());
        map.put("uuid2", uuid2);
        map.put("uuid2_string", uuid2.toString());
        output.writeObject(map);
        output.flush();
    }

    public void customReplyLocale() throws Exception {
        Map<String, Object> map = new HashMap<>();
        map.put("english", Locale.ENGLISH);
        map.put("french", Locale.FRENCH);
        map.put("german", Locale.GERMAN);
        map.put("italian", Locale.ITALIAN);
        map.put("japanese", Locale.JAPANESE);
        map.put("korean", Locale.KOREAN);
        map.put("chinese", Locale.CHINESE);
        map.put("simplified_chinese", Locale.SIMPLIFIED_CHINESE);
        map.put("traditional_chinese", Locale.TRADITIONAL_CHINESE);
        map.put("france", Locale.FRANCE);
        map.put("germany", Locale.GERMANY);
        map.put("japan", Locale.JAPAN);
        map.put("korea", Locale.KOREA);
        map.put("china", Locale.CHINA);
        map.put("prc", Locale.PRC);
        map.put("taiwan", Locale.TAIWAN);
        map.put("uk", Locale.UK);
        map.put("us", Locale.US);
        map.put("canada", Locale.CANADA);
        map.put("root", Locale.ROOT);
        // The two objects below is java hessian bug
        // map.put("italy", Locale.ITALY);
        // map.put("canada_french", Locale.CANADA_FRENCH);
        // LocaleHandle
        output.writeObject(map);
        output.flush();
    }

    public void customReplyEnumSet() throws Exception {
        Map<String, Object> map = new HashMap<>();
        EnumSet<Locale.Category> enumSet = EnumSet.allOf(Locale.Category.class);
        map.put("enumset", enumSet);
        System.out.println(map);
        output.writeObject(map);
        output.flush();
    }

    public void customReplyEnumVariableList() throws Exception {
        List<Locale.Category> enumList = new ArrayList<>();
        enumList.add(Locale.Category.DISPLAY);
        enumList.add(null);
        enumList.add(Locale.Category.FORMAT);
        output.writeObject(enumList.toArray(new Locale.Category[enumList.size()]));
        output.flush();
    }

    public void customReplyJavaLangObjectHolder() throws Exception {
        JavaLangObjectHolder holder = new JavaLangObjectHolder();

        holder.setFieldInteger(123);
        holder.setFieldLong(456L);
        holder.setFieldBoolean(true);
        holder.setFieldShort((short) 789);
        holder.setFieldByte((byte) 12);
        holder.setFieldFloat(3.45f);
        holder.setFieldDouble(6.78);
        holder.setFieldCharacter('A');

        output.writeObject(holder);
        output.flush();
    }

    public void customReplyJavaLangObjectHolderForNull() throws Exception {
        // all fields are default null.
        JavaLangObjectHolder holder = new JavaLangObjectHolder();
        output.writeObject(holder);
        output.flush();
    }
}

interface Leg {
    public int legConnt = 4;
}

class Animal {
    public String name;
}

class Dog extends Animal implements Serializable, Leg {
    public String gender;
}

class DogAll extends Dog {
    public boolean all = true;
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

class Person183 implements Serializable {
    public String name;
    public Integer age;
    public InnerPerson innerPerson;
}

class InnerPerson implements Serializable {
    public String name;
    public Integer age;
}

