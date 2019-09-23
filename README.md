# gohessian

[![Build Status](https://travis-ci.org/apache/dubbo-go-hessian2.png?branch=master)](https://travis-ci.org/apache/dubbo-go-hessian2)
[![GoCover](http://gocover.io/_badge/github.com/apache/dubbo-go-hessian2)](http://gocover.io/github.com/apache/dubbo-go-hessian2)
[![GoDoc](https://godoc.org/github.com/apache/dubbo-go-hessian2?status.svg)](https://godoc.org/github.com/apache/dubbo-go-hessian2)


---

It's a golang hessian library used by [Apache/dubbo-go](https://github.com/apache/dubbo-go).

## Feature List

* [All JDK Exceptions](https://github.com/apache/dubbo-go-hessian2/issues/59)
* [Function Alias By Alias](https://github.com/apache/dubbo-go-hessian2/issues/19)
* [Java Bigdecimal](https://github.com/apache/dubbo-go-hessian2/issues/89)
* [Java Date & Time](https://github.com/apache/dubbo-go-hessian2/issues/90)
* [Java Generic Invokation](https://github.com/apache/dubbo-go-hessian2/issues/84)
* [Dubbo Attachements](https://github.com/apache/dubbo-go-hessian2/issues/49)
* [Skipping unregistered POJO](https://github.com/apache/dubbo-go-hessian2/pull/128)
* [Emoji](https://github.com/apache/dubbo-go-hessian2/issues/129)

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
circular, ok := obj.(*Circular)
// ...
```

## Customize Usage Examples

#### Encoding filed name

Hessian encoder default converts filed names of struct to lower camelcase, but you can customize it using `hessian` tag.

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

The encoded bytes of the struct `MyUser` is as following:
```text
 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
 00000010  75 73 65 72 92 0e 75 73  65 72 5f 66 75 6c 6c 5f  |user..user_full_|
 00000020  6e 61 6d 65 11 66 61 6d  69 6c 79 50 68 6f 6e 65  |name.familyPhone|
 00000030  4e 75 6d 62 65 72 60 08  75 73 65 72 6e 61 6d 65  |Number`.username|
 00000040  0c 30 31 30 2d 31 32 33  34 35 36 37 38           |.010-12345678|
```

#### Decoding filed name

Hessian decoder finds the correct target field though comparing all filed names of struct one by one until matching.

The following example shows the order of the matching rules:
```go
type MyUser struct {
	MobilePhone      string   `hessian:"mobile-phone"`
}

// You must define the tag of struct for lookup filed form encoded binary bytes, in this case：
// 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
// 00000010  75 73 65 72 91 0c 6d 6f  62 69 6c 65 2d 70 68 6f  |user..mobile-pho|
// 00000020  6e 65 60 0b 31 37 36 31  32 33 34 31 32 33 34     |ne`.17612341234|
//
// mobile-phone(tag lookup) => mobilePhone(lowerCameCase) => MobilePhone(SameCase) => mobilephone(lowercase)
// ^ will matched


type MyUser struct {
	MobilePhone      string
}

// The following encoded binary bytes will be hit automatically:
//
// 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
// 00000010  75 73 65 72 91 0b 6d 6f  62 69 6c 65 50 68 6f 6e  |user..mobilePhon|
// 00000020  65 60 0b 31 37 36 31 32  33 34 31 32 33 34        |e`.17612341234|
//
// mobile-phone(tag lookup) => mobilePhone(lowerCameCase) => MobilePhone(SameCase) => mobilephone(lowercase)
//                             ^ will matched
//
// 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
// 00000010  75 73 65 72 91 0b 4d 6f  62 69 6c 65 50 68 6f 6e  |user..MobilePhon|
// 00000020  65 60 0b 31 37 36 31 32  33 34 31 32 33 34        |e`.17612341234|
//
// mobile-phone(tag lookup) => mobilePhone(lowerCameCase) => MobilePhone(SameCase) => mobilephone(lowercase)
//                                                           ^ will matched
//
// 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
// 00000010  75 73 65 72 91 0b 6d 6f  62 69 6c 65 70 68 6f 6e  |user..mobilephon|
// 00000020  65 60 0b 31 37 36 31 32  33 34 31 32 33 34        |e`.17612341234|
//
// mobile-phone(tag lookup) => mobilePhone(lowerCameCase) => MobilePhone(SameCase) => mobilephone(lowercase)
//                                                                                    ^ will matched

```


##### hessian.SetTagIdentifier

You can use `hessian.SetTagIdentifier` to customize tag-identifier of hessian, which takes effect to both encoder and decoder.

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

The encoded bytes of the struct `MyUser` is as following:
```text
 00000000  43 12 63 6f 6d 2e 63 6f  6d 70 61 6e 79 2e 6d 79  |C.com.company.my|
 00000010  75 73 65 72 92 0e 75 73  65 72 5f 66 75 6c 6c 5f  |user..user_full_|
 00000020  6e 61 6d 65 11 66 61 6d  69 6c 79 50 68 6f 6e 65  |name.familyPhone|
 00000030  4e 75 6d 62 65 72 60 08  75 73 65 72 6e 61 6d 65  |Number`.username|
 00000040  0c 30 31 30 2d 31 32 33  34 35 36 37 38           |.010-12345678|
```
