/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Strings []string

func (t *Strings) String() string {
	return fmt.Sprint(*t)
}

func (t *Strings) Set(value string) error {
	*t = append(*t, value)
	return nil
}

type viewData struct {
	ClassName   string
	PackageName string
	EnumName    string
	EnumValues  []string
}

func snakeCase(name string) string {
	if name == "" {
		return ""
	}

	var s string

	for k, v := range []rune(name) {
		if v >= 'A' && v <= 'Z' {
			if k > 0 {
				s += "_"
			}

			v += 32
		}

		s += string(v)
	}

	return s
}

func pascalCase(name string) string {
	if name == "" {
		return ""
	}

	var s string

	for _, v := range strings.Split(name, "_") {
		runes := []rune(v)

		for k, vv := range runes {
			if k == 0 {
				if vv >= 'a' && vv <= 'z' {
					vv -= 32
				}
			} else if vv >= 'A' && vv <= 'Z' {
				vv += 32
			}

			s += string(vv)
		}
	}

	return s
}

const (
	usage = `gen-go-enum can generate golang code for hessian2 java enum.

Usage: gen-go-enum -c java_classname [-p golang_package_name] [-e golang_enum_name] -v java_enum_value [-v java_enum_value] [-o target_file]

Options
  -c	java class name (eg: com.test.enums.TestEnum)
  -p	golang package name, use 'enum' when not specified (eg: test_enums)
  -e	golang enum type name, use java class name when not specified (eg: TestEnum)
  -v	java enum values, can specify multiple (eg: -v TEST1 -v TEST2 -v TEST3)
  -o 	golang code file path, stored in the current directory when not specified

Example
  gen-go-enum -c com.test.enums.TestColorEnum -v RED -v BLUE -v YELLOW
  gen-go-enum -c com.test.enums.TestColorEnum -p test_enums -e ColorEnum -v RED -v BLUE -v YELLOW -o ./color_enum.go
`

	enumTempl = `package {{.PackageName}}

import (
	"strconv"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

const (
{{- range $index, $value := .EnumValues }}
	{{$.EnumName}}{{ pascalCase $value}}{{ if eq $index 0 }} {{$.EnumName}} = iota{{ end }}
{{- end }}
)

var _{{.EnumName}}Values = map[{{.EnumName}}]string{
{{- range $index, $value := .EnumValues }}
	{{$.EnumName}}{{ pascalCase $value}}: "{{$value}}",
{{- end }}
}

var _{{.EnumName}}Entities = map[string]{{.EnumName}}{
{{- range $index, $value := .EnumValues }}
	"{{$value}}": {{$.EnumName}}{{ pascalCase $value}},
{{- end }}
}

type {{.EnumName}} hessian.JavaEnum

func (e {{.EnumName}}) JavaClassName() string {
	return "{{.ClassName}}"
}

func (e {{.EnumName}}) String() string {
	if v, ok := _{{.EnumName}}Values[e]; ok {
		return v
	}

	return strconv.Itoa(int(e))
}

func (e {{.EnumName}}) EnumValue(s string) hessian.JavaEnum {
	if v, ok := _{{.EnumName}}Entities[s]; ok {
		return hessian.JavaEnum(v)
	}

	return hessian.InvalidJavaEnum
}

func New{{.EnumName}}(s string) {{.EnumName}} {
	if v, ok := _{{.EnumName}}Entities[s]; ok {
		return v
	}

	return {{.EnumName}}(hessian.InvalidJavaEnum)
}

func init() {
	for v := range _{{.EnumName}}Values {
		hessian.RegisterJavaEnum(v)
	}
}`
)

var (
	className   string
	packageName string
	enumName    string
	enumValues  Strings
	targetFile  string
)

func init() {
	flag.StringVar(&className, "c", "", "")
	flag.StringVar(&packageName, "p", "enum", "")
	flag.StringVar(&enumName, "e", "", "")
	flag.Var(&enumValues, "v", "")
	flag.StringVar(&targetFile, "o", "", "")

	flag.Usage = func() {
		fmt.Print(usage)
	}
}

func main() {
	flag.Parse()

	if className == "" || len(enumValues) == 0 {
		flag.Usage()
		return
	}

	if packageName == "" {
		packageName = "enum"
	}

	if enumName == "" {
		classItems := strings.Split(className, ".")
		enumName = classItems[len(classItems)-1]
	}

	if targetFile == "" {
		targetFile = snakeCase(enumName) + ".go"
	}

	tmpl, err := template.New("enum").Funcs(template.FuncMap{"pascalCase": pascalCase}).Parse(enumTempl)

	if err != nil {
		log.Fatalln("Error: can't parse code template!!!", err)
	}

	fd, err := os.Create(targetFile)

	if err != nil {
		log.Fatalln("Error: can't open target file!!!", err)
	}

	defer func() { _ = fd.Close() }()

	err = tmpl.Execute(fd, &viewData{
		ClassName:   className,
		PackageName: packageName,
		EnumName:    enumName,
		EnumValues:  enumValues,
	})

	if err != nil {
		log.Fatalln("Error: can't write target file!!!", err)
	}

	fmt.Printf("Create file '%s'!\n", targetFile)
}
