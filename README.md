# gohessian

[![Build Status](https://travis-ci.org/dubbogo/hessian2.png?branch=master)](https://travis-ci.org/dubbogo/hessian2)
[![GoCover](http://gocover.io/_badge/github.com/dubbogo/hessian2)](http://gocover.io/github.com/dubbogo/hessian2)
[![GoDoc](https://godoc.org/github.com/dubbogo/hessian2?status.svg)](https://godoc.org/github.com/dubbogo/hessian2)


---

It's a golang hessian library used by dubbogo.

It was first build in project [viant/gohessian](https://github.com/viant/gohessian), 
and then improved by [AlexStocks](https://github.com/AlexStocks).
Thanks to [viant](https://github.com/viant) and [AlexStocks](https://github.com/AlexStocks) for their great work.

## Basic Usage Examples

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

## Customize Usage Examples

#### Struct filed name encoding
Hessian encoder will convert the filed name of struct to lower camelcase defaulted, but you customize it by `hessian` tag of struct.

Example:
```go
type MyUser struct {
	UserFullName      string   `hessian:"user_full_name"`
	FamilyPhoneNumber string   // default convert to => familyPhoneNumber
}

func (MyUser) JavaClassName() string {
	return "com.company.myuser"
}

user := &MyUser{
    UserFullName:      "username",
    FamilyPhoneNumber: "010-12345678",
}

e := hessian.NewEncoder()
err := e.Encode(user)
if err != nil {
    panic(err)
}
``` 
The encoded bytes of the struct from above example will look like:
```text
 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
 00000010  75 73 65 72 92 0e 75 73  65 72 5f 66 75 6c 6c 5f  |user..user_full_|
 00000020  6e 61 6d 65 11 66 61 6d  69 6c 79 50 68 6f 6e 65  |name.familyPhone|
 00000030  4e 75 6d 62 65 72 60 08  75 73 65 72 6e 61 6d 65  |Number`.username|
 00000040  0c 30 31 30 2d 31 32 33  34 35 36 37 38           |.010-12345678|
```

#### Struct filed name decoding
Hessian decoder will compare all filed's name of struct until it matches the correct name, the order of matching rules is:
```text
tag lookup => lower camelcase => samecase => lowercase
```


##### hessian.SetTagIdentifier

You can use `hessian.SetTagIdentifier` to customize tag-identifier of hessian, it's will effect both encoder and decoder. 

Example:
```go
hessian.SetTagIdentifier("json")

type MyUser struct {
	UserFullName      string   `json:"user_full_name"`
	FamilyPhoneNumber string   // default convert to => familyPhoneNumber
}

func (MyUser) JavaClassName() string {
	return "com.company.myuser"
}

user := &MyUser{
    UserFullName:      "username",
    FamilyPhoneNumber: "010-12345678",
}

e := hessian.NewEncoder()
err := e.Encode(user)
if err != nil {
    panic(err)
}
``` 
The encoded bytes of the struct from above example will look like:
```text
 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
 00000010  75 73 65 72 92 0e 75 73  65 72 5f 66 75 6c 6c 5f  |user..user_full_|
 00000020  6e 61 6d 65 11 66 61 6d  69 6c 79 50 68 6f 6e 65  |name.familyPhone|
 00000030  4e 75 6d 62 65 72 60 08  75 73 65 72 6e 61 6d 65  |Number`.username|
 00000040  0c 30 31 30 2d 31 32 33  34 35 36 37 38           |.010-12345678|
```


## Dubbo Service

TODO