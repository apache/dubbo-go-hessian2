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

import com.caucho.hessian.io.Hessian2Input;
import com.caucho.hessian.test.A0;
import com.caucho.hessian.test.A1;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;


public class TestCustomDecode {

    private Hessian2Input input;

    TestCustomDecode(InputStream is) {
        input = new Hessian2Input(is);
    }

    public Object customArgUntypedFixedListHasNull() throws Exception {
        List list = new ArrayList();
        list.add(new A0());
        list.add(new A1());
        list.add(null);

        Object o = input.readObject();
        return list.equals(o);
    }
}