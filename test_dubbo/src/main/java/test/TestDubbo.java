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