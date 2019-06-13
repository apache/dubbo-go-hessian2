package exception

import (
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
)
import (
	hessian "github.com/dubbogo/hessian2"
)
const (
	hessianJar = "../test_hessian/target/test_hessian-1.0.0.jar"
)
type _refHolder struct {
	// destinations
	destinations []reflect.Value

	value reflect.Value
}
func isFileExist(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func genHessianJar() {
	existFlag := isFileExist(hessianJar)
	if existFlag {
		return
	}

	cmd := exec.Command("mvn", "clean", "package")
	cmd.Dir = "./test_hessian"
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("after exec command 'mvn clean package', got error:%v, error output:%v",
			err, string(out))
	}
}

func getReply(method string) []byte {
	genHessianJar()
	cmd := exec.Command("java", "-jar", hessianJar, method)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func decodeResponse(method string) (interface{}, error) {
	b := getReply(method)
	d := hessian.NewDecoder(b)
	r, e := d.Decode()
	if e != nil {
		return nil, e
	}
	return r, nil
}

func testDecodeFramework(t *testing.T, method string, expected interface{}) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	if !reflect.DeepEqual(r, expected) {
		t.Errorf("%s: got %v, wanted %v", method, r, expected)
	}
}

func testDecodeFrameworkFunc(t *testing.T, method string, expected func(interface{})) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	expected(r)
}

