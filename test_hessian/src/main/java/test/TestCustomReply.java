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
import com.caucho.hessian.test.A0;
import com.caucho.hessian.test.A1;

import java.io.OutputStream;
import java.util.Date;
import java.util.HashMap;


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
        for (Object tmp: o) {
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
        for (Object tmp: o) {
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
        for (Object tmp: o) {
            output.writeObject(tmp);
        }
        if (hasEnd) {
            output.writeListEnd();
        }
        output.flush();
    }
}