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

package test.util;

import com.alibaba.com.caucho.hessian.io.Hessian2Input;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.util.ArrayList;

/**
 * @author tiltwind
 */
public class JavaHessianUtil {

    private static volatile boolean isRemove = false;

    public static void removeTestClassLimitFromAllowList() {
        if (isRemove) {
            return;
        }

        synchronized (JavaHessianUtil.class) {
            if (isRemove) {
                return;
            }

            isRemove = true;

            try {
                Object classFactory = new Hessian2Input(null).findSerializerFactory().getClassFactory();
                Class classFactoryClass = Class.forName("com.alibaba.com.caucho.hessian.io.ClassFactory");
                Field allowListField = classFactoryClass.getDeclaredField("_staticAllowList");
                allowListField.setAccessible(true);
                Object allowListObject = allowListField.get(classFactory);
                ArrayList allowList = (ArrayList) allowListObject;

                Class allowClass = Class.forName("com.alibaba.com.caucho.hessian.io.ClassFactory$Allow");
                Method allowMethod = allowClass.getDeclaredMethod("allow", String.class);
                allowMethod.setAccessible(true);

                for (int i = 0; i < allowList.size(); i++) {
                    Object item = allowList.get(i);
                    Object o = allowMethod.invoke(item, "com.caucho.hessian.test.A0");
                    if (o == null) {
                        continue;
                    }
                    Boolean allow = (Boolean) o;

                    if (!allow) {
                        System.err.println("remove hessian inner test object deny limit");
                        allowList.remove(i);
                        return;
                    }
                }

            } catch (Exception e) {
                System.err.println(e.getMessage());
            }
        }
    }
}
