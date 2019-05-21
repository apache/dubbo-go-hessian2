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

import org.apache.dubbo.remoting.Codec2;
import org.apache.dubbo.rpc.protocol.dubbo.DubboCodec;
import org.apache.dubbo.remoting.Channel;
import org.apache.dubbo.remoting.buffer.ChannelBuffer;
import org.apache.dubbo.remoting.buffer.DynamicChannelBuffer;

import java.lang.reflect.Method;


public class Dubbo {

    private static Object getReply(String methodString) throws Exception {
        Method method = TestDubbo.class.getMethod(methodString);
        TestDubbo testDubbo = new TestDubbo();
        return method.invoke(testDubbo);
    }

    public static void main(String[] args) throws Exception {
        Channel channel = new SimpleChannel();
        ChannelBuffer buffer = new DynamicChannelBuffer(4096);
        Object object = getReply(args[0]);

        Codec2 codec = new DubboCodec();
        codec.encode(channel, buffer, object);

        System.out.write(buffer.array());
        System.out.flush();
    }
}