# gen-go-enum

A tool for generate hessian2 java enum define golang code.

```sh
go build -o gen-go-enum tools/gen-go-enum/main.go
```

For example, run command `gen-go-enum -c com.test.enums.TestColorEnum -v RED -v BLUE -v YELLOW` will create a golang code file like this.

```go
package enum

import (
	"strconv"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

const (
	TestColorEnumRed TestColorEnum = iota
	TestColorEnumBlue
	TestColorEnumYellow
)

var _TestColorEnumValues = map[TestColorEnum]string{
	TestColorEnumRed: "RED",
	TestColorEnumBlue: "BLUE",
	TestColorEnumYellow: "YELLOW",
}

var _TestColorEnumEntities = map[string]TestColorEnum{
	"RED": TestColorEnumRed,
	"BLUE": TestColorEnumBlue,
	"YELLOW": TestColorEnumYellow,
}

type TestColorEnum hessian.JavaEnum

func (e TestColorEnum) JavaClassName() string {
	return "com.test.enums.TestColorEnum"
}

func (e TestColorEnum) String() string {
	if v, ok := _TestColorEnumValues[e]; ok {
		return v
	}

	return strconv.Itoa(int(e))
}

func (e TestColorEnum) EnumValue(s string) hessian.JavaEnum {
	if v, ok := _TestColorEnumEntities[s]; ok {
		return hessian.JavaEnum(v)
	}

	return hessian.InvalidJavaEnum
}

func NewTestColorEnum(s string) TestColorEnum {
	if v, ok := _TestColorEnumEntities[s]; ok {
		return v
	}

	return TestColorEnum(hessian.InvalidJavaEnum)
}

func init() {
	for v := range _TestColorEnumValues {
		hessian.RegisterJavaEnum(v)
	}
}
```

You can specify more options, like the usage.

```sh
gen-go-enum can generate golang code for hessian2 java enum.

Usage: gen-go-enum -c java_classname [-p golang_package_name] [-e golang_enum_name] -v java_enum_value [-v java_enum_value] [-o target_file]

Options
  -c	java class name (eg: com.test.enums.TestEnum)
  -p	golang package name, use 'enum' when not specified (eg: test_enum)
  -e	golang enum type name, use java class name when not specified (eg: TestEnum)
  -v	java enum values, can specify multiple (eg: -v TEST1 -v TEST2 -v TEST3)
  -o 	golang code file path, stored in the current directory when not specified

Example
  gen-go-enum -c com.test.enums.TestColorEnum -v RED -v BLUE -v YELLOW
  gen-go-enum -c com.test.enums.TestColorEnum -p test_enums -e ColorEnum -v RED -v BLUE -v YELLOW -o ./color_enum.go
```