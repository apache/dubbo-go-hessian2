package classfile

/*
ConstantValue_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 constantvalue_index;
}
*/
type ConstantValueAttribute struct {
	constantValueIndex uint16
}

func (cva *ConstantValueAttribute) readInfo(reader *ClassReader) {
	cva.constantValueIndex = reader.readUint16()
}

func (cva *ConstantValueAttribute) ConstantValueIndex() uint16 {
	return cva.constantValueIndex
}
