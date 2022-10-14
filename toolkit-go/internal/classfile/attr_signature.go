package classfile

/*
Signature_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 signature_index;
}
*/
type SignatureAttribute struct {
	cp             ConstantPool
	signatureIndex uint16
}

func (sa *SignatureAttribute) readInfo(reader *ClassReader) {
	sa.signatureIndex = reader.readUint16()
}

func (sa *SignatureAttribute) Signature() string {
	return sa.cp.getUtf8(sa.signatureIndex)
}
