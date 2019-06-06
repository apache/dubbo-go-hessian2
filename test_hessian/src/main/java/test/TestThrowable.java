package test;

public class TestThrowable {
    public static Object throw_exception()  {
        return new Exception("exception");
    }
    public static Object throw_throwable()  {
        return new Throwable("exception");
    }
}
