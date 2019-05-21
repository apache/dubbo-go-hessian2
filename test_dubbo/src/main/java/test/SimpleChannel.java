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

import org.apache.dubbo.common.URL;
import org.apache.dubbo.remoting.Channel;
import org.apache.dubbo.remoting.ChannelHandler;
import org.apache.dubbo.remoting.RemotingException;

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