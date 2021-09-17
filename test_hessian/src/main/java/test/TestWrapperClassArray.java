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

public class TestWrapperClassArray {

    public static Byte[] byteArray() {
        return new Byte[]{1, 100, -56};
    }

    public static Short[] shortArray() {
        return new Short[]{1, 100, 10000};
    }

    public static Integer[] integerArray() {
        return new Integer[]{1, 100, 10000};
    }

    public static Long[] longArray() {
        return new Long[]{1L, 100L, 10000L};
    }

    public static Boolean[] booleanArray() {
        return new Boolean[]{true, false, true};
    }

    public static Float[] floatArray() {
        return new Float[]{1.0f, 100.0f, 10000.1f};
    }

    public static Double[] doubleArray() {
        return new Double[]{1.0, 100.0, 10000.1};
    }

    public static Character[] characterArray() {
        return new Character[]{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'};
    }
}
