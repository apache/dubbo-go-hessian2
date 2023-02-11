package testfuncs

import hessian "github.com/apache/dubbo-go-hessian2"

type A0 struct{}

// JavaClassName  java fully qualified path
func (*A0) JavaClassName() string {
	return "com.caucho.hessian.test.A0"
}

func ObjectA0() []byte {
	e := hessian.NewEncoder()
	var a = A0{}
	e.Encode(&a)
	return e.Buffer()
}
