// Copyright 2019 Xinge Gao
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

import com.alibaba.com.caucho.hessian.io.Hessian2Output;
import com.caucho.hessian.test.TestHessian2Servlet;

import java.lang.reflect.Method;


public class Hessian {
    public static void main(String[] args) throws Exception {
        if (args[0].startsWith("reply")) {
            Method method = TestHessian2Servlet.class.getMethod(args[0]);
            TestHessian2Servlet servlet = new TestHessian2Servlet();
            Object object = method.invoke(servlet);

            Hessian2Output output = new Hessian2Output(System.out);
            output.writeObject(object);
            output.flush();
        } else {
            Method method = TestHessian.class.getMethod(args[0]);
            TestHessian testHessian = new TestHessian(System.out);
            method.invoke(testHessian);
        }
    }
}