package test.tuple;

public class TupleProviderImpl implements TupleProvider {

    @Override
    public Tuple getTheTuple() {
        Tuple result = new Tuple();
        result.setB((byte) 1);
        result.setByte(Byte.valueOf("1"));
        result.setI(1);
        result.setInteger(Integer.valueOf("1"));
        result.setL(1L);
        result.setLong(Long.valueOf("1"));
        result.setS((short) 1);
        result.setShort(Short.valueOf("1"));
        result.setD(1.23);
        result.setDouble(Double.valueOf("1.23"));
        return result;
    }

}
