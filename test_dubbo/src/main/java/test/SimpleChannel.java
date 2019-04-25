package test;

import com.alibaba.dubbo.common.URL;
import com.alibaba.dubbo.remoting.Channel;
import com.alibaba.dubbo.remoting.ChannelHandler;
import com.alibaba.dubbo.remoting.RemotingException;

import java.net.InetSocketAddress;


public class SimpleChannel implements Channel {
    // Override the method in Channel interface.

    @Override
    public InetSocketAddress getRemoteAddress() {
        return null;
    }

    @Override
    public boolean isConnected() {
        return true;
    }

    @Override
    public boolean hasAttribute(String var1) {
        return true;
    }

    @Override
    public Object getAttribute(String var1) {
        return null;
    }

    @Override
    public void setAttribute(String var1, Object var2) {

    }

    @Override
    public void removeAttribute(String var1) {

    }

    // Override the method in Endpoint interface.

    @Override
    public URL getUrl() {
        URL url = new URL("dubbo", "127.0.0.1", 8080);
        url.addParameter("serialization", "hessian2");
        return url;
    }

    @Override
    public ChannelHandler getChannelHandler() {
        return null;
    }

    @Override
    public InetSocketAddress getLocalAddress() {
        return null;
    }

    @Override
    public void send(Object var1) throws RemotingException {

    }

    @Override
    public void send(Object var1, boolean var2) throws RemotingException {

    }

    @Override
    public void close() {

    }

    @Override
    public void close(int var1) {

    }

    @Override
    public void startClose() {

    }

    @Override
    public boolean isClosed() {
        return true;
    }
}