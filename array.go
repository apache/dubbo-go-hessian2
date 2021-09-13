package hessian

func init() {
	SetCollectionSerialize(&IntegerArray{})
}

type IntegerArray struct {
	Values []int32
}

func (ia IntegerArray) Get() []interface{} {
	res := make([]interface{}, len(ia.Values))
	for i, v := range ia.Values {
		res[i] = v
	}
	return res
}

func (ia IntegerArray) Set(vs []interface{}) {
	values := make([]int32, len(vs))
	for i, v := range vs {
		values[i] = v.(int32)
	}
	ia.Values = values
}

func (IntegerArray) JavaClassName() string {
	return "[java.lang.Integer"
}

type ArraySerializer JavaCollectionSerializer
