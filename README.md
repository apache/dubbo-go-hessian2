# dubbo-go-hessian2

[![Build Status](https://travis-ci.org/apache/dubbo-go-hessian2.png?branch=master)](https://travis-ci.org/apache/dubbo-go-hessian2)
[![codecov](https://codecov.io/gh/apache/dubbo-go-hessian2/branch/master/graph/badge.svg)](https://codecov.io/gh/apache/dubbo-go-hessian2)
[![GoDoc](https://godoc.org/github.com/apache/dubbo-go-hessian2?status.svg)](https://godoc.org/github.com/apache/dubbo-go-hessian2)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/dubbo-go-hessian2)](https://goreportcard.com/report/github.com/apache/dubbo-go-hessian2)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

---

A Go implementation of the [Hessian 2.0 serialization](http://hessian.caucho.com/doc/hessian-serialization.html) protocol, enabling seamless cross-language communication between Go and Java applications. Primarily used by [Apache Dubbo-Go](https://github.com/apache/dubbo-go) for RPC serialization.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
  - [The POJO Interface](#the-pojo-interface)
  - [Registering Types](#registering-types)
- [Type Mapping (Java <-> Go)](#type-mapping-java---go)
  - [Primitive Types](#primitive-types)
  - [Java Wrapper Types](#java-wrapper-types)
  - [Extended Types](#extended-types)
- [API Reference](#api-reference)
  - [Encoder](#encoder)
  - [Decoder](#decoder)
  - [Registration Functions](#registration-functions)
  - [Configuration Functions](#configuration-functions)
- [Usage Examples](#usage-examples)
  - [Encoding and Decoding Custom Objects](#encoding-and-decoding-custom-objects)
  - [Customizing Field Names with Tags](#customizing-field-names-with-tags)
  - [Decoding Field Name Matching Rules](#decoding-field-name-matching-rules)
  - [Using a Custom Tag Identifier](#using-a-custom-tag-identifier)
  - [Specifying Java Parameter Types (Inheritance)](#specifying-java-parameter-types-inheritance)
  - [Working with Java Collections](#working-with-java-collections)
  - [Working with Java Enums](#working-with-java-enums)
  - [Custom Serializer](#custom-serializer)
  - [Strict Mode](#strict-mode)
  - [Reusing Encoder/Decoder (Object Pool)](#reusing-encoderdecoder-object-pool)
- [Struct Inheritance](#struct-inheritance)
- [Tools](#tools)
- [Reference](#reference)

## Features

- Full Hessian 2.0 protocol implementation (encode/decode)
- Circular reference and object graph support
- [All JDK exceptions](https://github.com/apache/dubbo-go-hessian2/issues/59) (88+ exception types)
- [Java wrapper types](https://github.com/apache/dubbo-go-hessian2/issues/349) (Integer, Long, Boolean, etc.)
- [Java BigDecimal / BigInteger](https://github.com/apache/dubbo-go-hessian2/issues/89)
- [Java Date & Time](https://github.com/apache/dubbo-go-hessian2/issues/90) (java.util.Date, java.sql.Date/Time)
- [Java 8 Time API](https://github.com/apache/dubbo-go-hessian2/pull/212) (LocalDate, LocalDateTime, ZonedDateTime, Instant, Duration, etc.)
- [Java UUID](https://github.com/apache/dubbo-go-hessian2/pull/256)
- [Java collections](https://github.com/apache/dubbo-go-hessian2/issues/84) (HashSet, HashMap, etc.)
- [Java inheritance / extends](https://github.com/apache/dubbo-go-hessian2/issues/157)
- [Field alias via struct tags](https://github.com/apache/dubbo-go-hessian2/issues/19)
- [Generic invocation](https://github.com/apache/dubbo-go-hessian2/issues/84)
- [Dubbo attachments](https://github.com/apache/dubbo-go-hessian2/issues/49)
- [Skipping unregistered POJOs](https://github.com/apache/dubbo-go-hessian2/pull/128)
- [Emoji support in strings](https://github.com/apache/dubbo-go-hessian2/issues/129)
- Custom serializer support
- Strict mode for decoding validation

> **Note:** From v1.6.0+, the decoder skips non-existent fields (matching Java hessian behavior). Versions before v1.6.0 returned errors for non-existent fields.

## Installation

```bash
go get github.com/apache/dubbo-go-hessian2
```

Requires **Go 1.21+**.

## Quick Start

```go
// Define a struct and implement the POJO interface
type User struct {
    Name string
    Age  int32
}

func (User) JavaClassName() string {
    return "com.example.User"
}

// Register the type before encoding/decoding
hessian.RegisterPOJO(&User{})

// Encode
user := &User{Name: "Alice", Age: 30}
encoder := hessian.NewEncoder()
encoder.Encode(user)
data := encoder.Buffer()

// Decode
obj, _ := hessian.NewDecoder(data).Decode()
decoded := obj.(*User)
```

## Core Concepts

### The POJO Interface

Any Go struct that needs to be serialized as a Java object must implement the `POJO` interface:

```go
type POJO interface {
    JavaClassName() string  // Returns the fully qualified Java class name
}
```

This establishes the mapping between your Go struct and the corresponding Java class.

### Registering Types

Before encoding or decoding custom types, you must register them. This is typically done in an `init()` function:

```go
func init() {
    // Register a single type
    hessian.RegisterPOJO(&MyStruct{})

    // Register multiple types at once
    hessian.RegisterPOJOs(&TypeA{}, &TypeB{}, &TypeC{})

    // Register with a custom Java class name (overrides JavaClassName())
    hessian.RegisterPOJOMapping("com.example.CustomName", &MyStruct{})
}
```

If a type is not registered, the decoder will:
- **Default mode:** Decode unknown objects as `map[interface{}]interface{}`
- **Strict mode:** Return an error

## Type Mapping (Java <-> Go)

### Primitive Types

| Hessian Type | Java Type          | Go Type     |
|--------------|--------------------|-------------|
| null         | null               | nil         |
| binary       | byte[]             | []byte      |
| boolean      | boolean            | bool        |
| date         | java.util.Date     | time.Time   |
| double       | double             | float64     |
| int          | int                | int32       |
| long         | long               | int64       |
| string       | java.lang.String   | string      |
| list         | java.util.List     | slice       |
| map          | java.util.Map      | map         |
| object       | custom class       | struct      |

### Java Wrapper Types

| Java Type             | Go Type   |
|-----------------------|-----------|
| java.lang.Integer     | \*int32   |
| java.lang.Long        | \*int64   |
| java.lang.Boolean     | \*bool    |
| java.lang.Short       | \*int16   |
| java.lang.Byte        | \*uint8   |
| java.lang.Float       | \*float32 |
| java.lang.Double      | \*float64 |
| java.lang.Character   | \*hessian.Rune |

### Extended Types

| Java Type                   | Go Type / Package                                          |
|-----------------------------|------------------------------------------------------------|
| java.math.BigDecimal        | `github.com/dubbogo/gost/math/big` Decimal                |
| java.math.BigInteger        | `github.com/dubbogo/gost/math/big` Integer                |
| java.sql.Date               | `github.com/apache/dubbo-go-hessian2/java_sql_time` Date  |
| java.sql.Time               | `github.com/apache/dubbo-go-hessian2/java_sql_time` Time  |
| java.util.UUID              | `github.com/apache/dubbo-go-hessian2/java_util` UUID      |
| java.util.Locale            | `github.com/apache/dubbo-go-hessian2/java_util` Locale    |
| Java 8 time types (LocalDate, LocalTime, LocalDateTime, ZonedDateTime, Instant, Duration, Period, etc.) | `github.com/apache/dubbo-go-hessian2/java8_time` |

> **Tip:** Avoid defining objects that only exist in one language. Use error codes/messages instead of Java exceptions for cross-language communication.

## API Reference

### Encoder

```go
// Create a new encoder
encoder := hessian.NewEncoder()

// Encode a value (primitives, slices, maps, structs, etc.)
err := encoder.Encode(value)

// Get the encoded bytes
data := encoder.Buffer()

// Reset encoder state for reuse
encoder.Clean()

// Reset encoder state but reuse the underlying buffer
encoder.ReuseBufferClean()
```

**Special encoding methods:**

```go
// Encode a map as a typed Java object using "_class" key in the map
encoder.EncodeMapClass(map[string]interface{}{"_class": "com.example.Foo", "name": "bar"})

// Encode a map as a specific Java class
encoder.EncodeMapAsClass("com.example.Foo", map[string]interface{}{"name": "bar"})
```

### Decoder

```go
// Standard decoder
decoder := hessian.NewDecoder(data)

// Strict mode - returns error for unregistered types
decoder := hessian.NewStrictDecoder(data)

// Decode the next value
obj, err := decoder.Decode()

// Reset decoder with new data (for reuse with object pools)
decoder.Reset(newData)
```

**Decoder modes:**

| Constructor                      | Behavior |
|----------------------------------|----------|
| `NewDecoder(data)`               | Decodes unknown objects as maps |
| `NewStrictDecoder(data)`         | Returns error for unregistered objects |
| `NewDecoderWithSkip(data)`       | Skips non-existent fields |
| `NewCheapDecoderWithSkip(data)`  | Poolable decoder, use with `Reset()` |

### Registration Functions

```go
hessian.RegisterPOJO(&MyStruct{})                         // Register a POJO type
hessian.RegisterPOJOs(&A{}, &B{})                         // Register multiple POJOs
hessian.RegisterPOJOMapping("com.example.Name", &Struct{}) // Register with custom Java class name
hessian.RegisterJavaEnum(&MyEnum{})                        // Register a Java enum type
hessian.UnRegisterPOJOs(&A{}, &B{})                        // Unregister POJOs
hessian.SetCollectionSerialize(&MyHashSet{})               // Register a Java collection type
hessian.SetSerializer("com.example.Foo", &FooSerializer{}) // Register a custom serializer
```

### Configuration Functions

```go
// Change the struct tag used for field name mapping (default: "hessian")
hessian.SetTagIdentifier("json")

// Look up a registered custom serializer
serializer, ok := hessian.GetSerializer("com.example.Foo")

// Find class info in the decoder
classInfo := hessian.FindClassInfo("com.example.Foo")
```

## Usage Examples

### Encoding and Decoding Custom Objects

```go
type Circular struct {
    Value
    Previous *Circular
    Next     *Circular
}

type Value struct {
    Num int
}

func (Circular) JavaClassName() string {
    return "com.company.Circular"
}

func init() {
    hessian.RegisterPOJO(&Circular{})
}

// Encode
c := &Circular{}
c.Num = 12345
c.Previous = c  // circular reference - handled automatically
c.Next = c

e := hessian.NewEncoder()
if err := e.Encode(c); err != nil {
    panic(err)
}
data := e.Buffer()

// Decode
obj, err := hessian.NewDecoder(data).Decode()
if err != nil {
    panic(err)
}
circular := obj.(*Circular)
fmt.Println(circular.Num) // 12345
```

### Customizing Field Names with Tags

The encoder converts Go field names to **lowerCamelCase** by default. Use the `hessian` tag to override:

```go
type MyUser struct {
    UserFullName      string `hessian:"user_full_name"` // encoded as "user_full_name"
    FamilyPhoneNumber string                            // encoded as "familyPhoneNumber" (default)
}

func (MyUser) JavaClassName() string {
    return "com.company.myuser"
}
```

### Decoding Field Name Matching Rules

When decoding, fields are matched in the following order:

1. **Tag match** - matches the `hessian` tag value
2. **lowerCamelCase** - e.g., `mobilePhone` matches `MobilePhone`
3. **Exact case** - e.g., `MobilePhone` matches `MobilePhone`
4. **Lowercase** - e.g., `mobilephone` matches `MobilePhone`

```go
type MyUser struct {
    MobilePhone string `hessian:"mobile-phone"`
}
// Incoming field "mobile-phone" -> matched via tag (rule 1)
// Incoming field "mobilePhone"  -> matched via lowerCamelCase (rule 2)
// Incoming field "MobilePhone"  -> matched via exact case (rule 3)
// Incoming field "mobilephone"  -> matched via lowercase (rule 4)
```

### Using a Custom Tag Identifier

Use `SetTagIdentifier` to read field names from a different struct tag (e.g., `json`):

```go
hessian.SetTagIdentifier("json")

type MyUser struct {
    UserFullName      string `json:"user_full_name"`
    FamilyPhoneNumber string
}

func (MyUser) JavaClassName() string {
    return "com.company.myuser"
}
```

### Specifying Java Parameter Types (Inheritance)

When a Java method expects a parent class but you send a subclass, implement the `Param` interface:

**Java side:**
```java
public abstract class User {}

public class MyUser extends User implements Serializable {
    private String userFullName;
    private String familyPhoneNumber;
}

public interface UserProvider {
    String GetUser(User user);  // accepts parent type
}
```

**Go side:**
```go
type MyUser struct {
    UserFullName      string `hessian:"userFullName"`
    FamilyPhoneNumber string
}

func (m *MyUser) JavaClassName() string {
    return "com.company.MyUser"
}

// JavaParamName tells the encoder to use the parent class name in the method signature
func (m *MyUser) JavaParamName() string {
    return "com.company.User"
}
```

### Working with Java Collections

Map a Java collection class (e.g., `HashSet`) to a Go struct:

```go
type JavaHashSet struct {
    value []interface{}
}

func (j *JavaHashSet) Get() []interface{}       { return j.value }
func (j *JavaHashSet) Set(v []interface{})      { j.value = v }
func (j *JavaHashSet) JavaClassName() string    { return "java.util.HashSet" }

func init() {
    hessian.SetCollectionSerialize(&JavaHashSet{})
}
```

Without this registration, Java collections are decoded as `[]interface{}`.

### Working with Java Enums

```go
type Color int32

const (
    RED   Color = 0
    GREEN Color = 1
    BLUE  Color = 2
)

var colorNames = map[Color]string{
    RED: "RED", GREEN: "GREEN", BLUE: "BLUE",
}
var colorValues = map[string]Color{
    "RED": RED, "GREEN": GREEN, "BLUE": BLUE,
}

func (c Color) JavaClassName() string   { return "com.example.Color" }
func (c Color) String() string          { return colorNames[c] }
func (c Color) EnumValue(s string) hessian.JavaEnum {
    return hessian.JavaEnum(colorValues[s])
}

func init() {
    hessian.RegisterJavaEnum(RED)
}
```

### Custom Serializer

Implement the `Serializer` interface for full control over encoding/decoding:

```go
type MySerializer struct{}

func (s *MySerializer) EncObject(encoder *hessian.Encoder, obj hessian.POJO) error {
    // Custom encoding logic
    return nil
}

func (s *MySerializer) DecObject(decoder *hessian.Decoder, typ reflect.Type, cls *hessian.ClassInfo) (interface{}, error) {
    // Custom decoding logic
    return nil, nil
}

func init() {
    hessian.SetSerializer("com.example.MyClass", &MySerializer{})
}
```

### Strict Mode

By default, unregistered objects are decoded as maps. Use strict mode to get errors instead:

```go
decoder := hessian.NewDecoder(data)
decoder.Strict = true  // returns error for unregistered types

// Or use the convenience constructor:
decoder = hessian.NewStrictDecoder(data)
```

### Reusing Encoder/Decoder (Object Pool)

For high-performance scenarios, reuse encoder/decoder instances:

```go
// Encoder reuse
encoder := hessian.NewEncoder()
encoder.Encode(obj1)
data1 := encoder.Buffer()

encoder.Clean()  // or encoder.ReuseBufferClean() to keep the buffer
encoder.Encode(obj2)
data2 := encoder.Buffer()

// Decoder reuse (poolable decoder)
decoder := hessian.NewCheapDecoderWithSkip(data1)
obj1, _ := decoder.Decode()

decoder.Reset(data2)  // reuse with new data
obj2, _ := decoder.Decode()
```

## Struct Inheritance

Go struct embedding is supported for modeling Java inheritance:

```go
type Animal struct {
    Name string
}

func (Animal) JavaClassName() string { return "com.example.Animal" }

type Dog struct {
    Animal          // embedded parent struct
    Breed  string
}

func (Dog) JavaClassName() string { return "com.example.Dog" }
```

**Avoid these patterns:**

1. **Duplicate field names across parents** - ambiguous field resolution:
    ```go
    type A struct { Name string }
    type B struct { Name string }
    type C struct { A; B }  // which Name?
    ```

2. **Pointer embedding** - nil parent at initialization, not supported:
    ```go
    type Dog struct {
        *Animal  // will be nil in Dog{}, not supported
    }
    ```

## Tools

### gen-go-enum

A code generation tool for creating Go enum types compatible with Java enums. See [tools/gen-go-enum/README.md](tools/gen-go-enum/README.md) for details.

## Reference

- [Hessian 2.0 Serialization Protocol Spec](http://hessian.caucho.com/doc/hessian-serialization.html)
- [Apache Dubbo-Go](https://github.com/apache/dubbo-go)
- [GoDoc API Reference](https://godoc.org/github.com/apache/dubbo-go-hessian2)
