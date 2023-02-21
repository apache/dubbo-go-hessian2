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

package unit;

import com.caucho.hessian.test.A0;
import org.junit.Assert;
import org.junit.Test;
import test.util.JavaHessianUtil;

public class GoWrapperClassArrayTest {

    public GoWrapperClassArrayTest(){
        JavaHessianUtil.removeTestClassLimitFromAllowList();
    }

    @Test
    public void testArrays() {
        Assert.assertArrayEquals(new Byte[]{1, 100, -56}
                , (Byte[])GoTestUtil.readGoObject("ByteArray"));
        Assert.assertArrayEquals(new Short[]{1, 100, 10000}
                , (Short[])GoTestUtil.readGoObject("ShortArray"));
        Assert.assertArrayEquals(new Integer[]{1, 100, 10000}
                , (Integer[])GoTestUtil.readGoObject("IntegerArray"));
        Assert.assertArrayEquals(new Long[]{1L, 100L, 10000L}
                , (Long[])GoTestUtil.readGoObject("LongArray"));
        Assert.assertArrayEquals(new Boolean[]{true, false, true}
                , (Boolean[])GoTestUtil.readGoObject("BooleanArray"));
        Assert.assertArrayEquals(new Float[]{1.0f, 100.0f, 10000.1f}
                , (Float[])GoTestUtil.readGoObject("FloatArray"));
        Assert.assertArrayEquals(new Double[]{1.0, 100.0, 10000.1}
                , (Double[])GoTestUtil.readGoObject("DoubleArray"));
        Assert.assertArrayEquals(new Character[]{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
                , (Character[])GoTestUtil.readGoObject("CharacterArray"));
    }

    @Test
    public void testMultipleLevelArrays() {
        A0[][][] multipleLevelA0Array = new A0[][][]{{{new A0(), new A0(), new A0()}, {new A0(), new A0(), new A0(), null}}, {{new A0()}, {new A0()}}};

        Assert.assertArrayEquals(multipleLevelA0Array, (Object[]) GoTestUtil.readGoObject("MultipleLevelA0Array"));
    }
}
