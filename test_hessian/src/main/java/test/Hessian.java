package test;

import com.caucho.hessian.io.Hessian2Output;
import com.caucho.hessian.test.TestHessian2Servlet;

import java.lang.reflect.Method;


public class Hessian {
    public static void main(String[] args) throws Exception {
        Method method = TestHessian2Servlet.class.getMethod(args[0]);

        TestHessian2Servlet servlet = new TestHessian2Servlet();
        Object object = method.invoke(servlet);

        Hessian2Output output = new Hessian2Output(System.out);
        output.writeObject(object);
        output.flush();
    }
}
