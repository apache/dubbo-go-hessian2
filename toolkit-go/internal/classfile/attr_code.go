package classfile

/*
Code_attribute {
    u2 attribute_name_index;
    u4 attribute_length;
    u2 max_stack;
    u2 max_locals;
    u4 code_length;
    u1 code[code_length];
    u2 exception_table_length;
    {   u2 start_pc;
        u2 end_pc;
        u2 handler_pc;
        u2 catch_type;
    } exception_table[exception_table_length];
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
*/
type CodeAttribute struct {
	cp             ConstantPool
	maxStack       uint16
	maxLocals      uint16
	code           []byte
	exceptionTable []*ExceptionTableEntry
	attributes     []AttributeInfo
}

func (ca *CodeAttribute) readInfo(reader *ClassReader) {
	ca.maxStack = reader.readUint16()
	ca.maxLocals = reader.readUint16()
	codeLength := reader.readUint32()
	ca.code = reader.readBytes(codeLength)
	ca.exceptionTable = readExceptionTable(reader)
	ca.attributes = readAttributes(reader, ca.cp)
}

func (ca *CodeAttribute) MaxStack() uint {
	return uint(ca.maxStack)
}
func (ca *CodeAttribute) MaxLocals() uint {
	return uint(ca.maxLocals)
}
func (ca *CodeAttribute) Code() []byte {
	return ca.code
}
func (ca *CodeAttribute) ExceptionTable() []*ExceptionTableEntry {
	return ca.exceptionTable
}

type ExceptionTableEntry struct {
	startPc   uint16
	endPc     uint16
	handlerPc uint16
	catchType uint16
}

func readExceptionTable(reader *ClassReader) []*ExceptionTableEntry {
	exceptionTableLength := reader.readUint16()
	exceptionTable := make([]*ExceptionTableEntry, exceptionTableLength)
	for i := range exceptionTable {
		exceptionTable[i] = &ExceptionTableEntry{
			startPc:   reader.readUint16(),
			endPc:     reader.readUint16(),
			handlerPc: reader.readUint16(),
			catchType: reader.readUint16(),
		}
	}
	return exceptionTable
}

func (et *ExceptionTableEntry) StartPc() uint16 {
	return et.startPc
}
func (et *ExceptionTableEntry) EndPc() uint16 {
	return et.endPc
}
func (et *ExceptionTableEntry) HandlerPc() uint16 {
	return et.handlerPc
}
func (et *ExceptionTableEntry) CatchType() uint16 {
	return et.catchType
}
