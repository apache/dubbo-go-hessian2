# gohessian

[![Build Status](https://travis-ci.org/dubbogo/hessian2.png?branch=master)](https://travis-ci.org/dubbogo/hessian2)
[![GoCover](http://gocover.io/_badge/github.com/dubbogo/hessian2)](http://gocover.io/github.com/dubbogo/hessian2)
[![GoDoc](https://godoc.org/github.com/dubbogo/hessian2?status.svg)](https://godoc.org/github.com/dubbogo/hessian2)


---

It's a golang hessian library used by dubbogo.

It was first build in project [viant/gohessian](https://github.com/viant/gohessian), 
and then improved by [AlexStocks](https://github.com/AlexStocks).
Thanks to [viant](https://github.com/viant) and [AlexStocks](https://github.com/AlexStocks) for their great work.

## Usage Examples

### Encode To Bytes

```go
type Circular struct {
	Num      int
	Previous *Circular
	Next     *Circular
}

func (Circular) JavaClassName() string {
	return "com.company.Circular"
}

c := &Circular{}
c.Num = 12345
c.Previous = c
c.Next = c

e := NewEncoder()
err := e.Encode(c)
if err != nil {
    panic(err)
}

bytes := e.Buffer()
```

### Decode From Bytes

```go
decodedObject, err := NewDecoder(bytes).Decode()
if err != nil {
    panic(err)
}
obj = decodedObject.(reflect.Value).Interface()
circular, ok := obj.(*Circular)
// ...
```

### Dubbo Service

TODO