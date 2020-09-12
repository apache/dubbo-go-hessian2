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

/**
 * date 2020/9/12 11:09 <br/>
 * description class <br/>
 * test java8
 *
 * @author zhangyanmingjiayou@163.com
 * @version 1.0
 * @since 1.0
 */
public class GoJavaExceptionTest {

    /**
     * test java java.lang.Exception object and go java_exception Exception struct
     */
    @Test
    public void testException() {
        Exception exception = new Exception("java_exception");
        Object javaException = GoTestUtil.readGoObject("JavaException");
        if (javaException instanceof Exception) {
            Assert.assertEquals(exception.getMessage(), ((Exception) javaException).getMessage());
        }
        // assertEquals don't compare Exception object
        // Exception exception2 = new Exception("java_exception");
        // Assert.assertEquals(exception, exception2);
        // Assert.assertEquals(exception, GoTestUtil.readGoObject("JavaException"));
    }

}
