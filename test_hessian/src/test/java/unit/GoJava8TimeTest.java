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

import junit.framework.Assert;
import org.junit.Test;

import java.time.LocalDate;
import java.time.Year;

/**
 * date 2020/9/12 11:09 <br/>
 * description class <br/>
 * test java8
 *
 * @author zhangyanmingjiayou@163.com
 * @version 1.0
 * @since 1.0
 */
public class GoJava8TimeTest {

    /**
     * test java8 java.time.* object and go java8_time/* struct
     */
    @Test
    public void testJava8Year() {
        Year year = Year.of(2020);
        Assert.assertEquals(year
                , GoTestUtil.readGoObject("Java8TimeYear"));
        LocalDate localDate = LocalDate.of(2020, 9, 12);
        Assert.assertEquals(localDate, GoTestUtil.readGoObject("Java8LocalDate"));
    }

}
