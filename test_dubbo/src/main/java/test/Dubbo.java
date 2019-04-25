package test;

import com.alibaba.dubbo.remoting.Codec2;
import com.alibaba.dubbo.rpc.protocol.dubbo.DubboCodec;
import com.alibaba.dubbo.remoting.Channel;
import com.alibaba.dubbo.remoting.buffer.ChannelBuffer;
import com.alibaba.dubbo.remoting.buffer.DynamicChannelBuffer;

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