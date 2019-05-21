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

import org.apache.dubbo.rpc.RpcInvocation;
import org.apache.dubbo.remoting.exchange.Request;


public class TestDubbo {

    public Request replyRequest() {
        RpcInvocation rpcInvocation = new RpcInvocation();
        rpcInvocation.setMethodName("echo");
        rpcInvocation.setParameterTypes(new Class[]{String.class});
        rpcInvocation.setArguments(new Object[]{"hello world"});
        rpcInvocation.setAttachment("path", "dubbo-x/dubbo.DubboService");
        rpcInvocation.setAttachment("interface", "dubbo.DubboService");

        Request request = new Request(1);
        request.setData(rpcInvocation);

        return request;
    }
}