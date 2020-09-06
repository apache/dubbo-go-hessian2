package unit;

import junit.framework.Assert;
import org.junit.Test;
import test.TestCustomReply;

public class TestMap {

    @Test
    public void testHelloWordString() throws Exception {
        new TestCustomReply(System.out).customReplyMapInMap();
    }
}
