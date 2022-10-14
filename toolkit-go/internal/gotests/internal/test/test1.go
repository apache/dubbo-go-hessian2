package test

import (
	"fmt"
	"time"
)

type T1 interface {
}

type F1 func()

func name(t1 T1, f1 F1) {
	fmt.Println(time.Now())
}

func name2(f1 F1) {

}
