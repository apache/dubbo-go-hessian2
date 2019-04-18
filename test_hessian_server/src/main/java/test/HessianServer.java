package test;

import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.eclipse.jetty.servlet.ServletHolder;
import com.caucho.hessian.test.TestHessian2Servlet;


public class HessianServer {
    public static void main(String[] args) throws Exception {
        ServletContextHandler ctx = new ServletContextHandler(ServletContextHandler.SESSIONS);
        ctx.addServlet(new ServletHolder(new TestHessian2Servlet()), "/");

        Server server = new Server(8080);
        server.setHandler(ctx);
        server.start();
        server.join();
    }
}
