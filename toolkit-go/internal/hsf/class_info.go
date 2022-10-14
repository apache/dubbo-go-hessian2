package hsf

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/classfile"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/classpath"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/config"
)

const (
	javaTypeChar    = "C"
	javaTypeBoolean = "Z"
	javaTypeByte    = "B"
	javaTypeInt     = "I"
	javaTypeShort   = "S"
	javaTypeLong    = "J"
	javaTypeFloat   = "F"
	javaTypeDouble  = "D"
	javaTypeVoid    = "V"
	javaTypeSlice   = "["

	javaTypeObject           = "Ljava/lang/Object;"
	javaTypeLangString       = "Ljava/lang/String;"
	javaTypeLangNumber       = "Ljava/lang/Number;"
	javaTypeLangBoolean      = "Ljava/lang/Boolean;"
	javaTypeLangByte         = "Ljava/lang/Byte;"
	javaTypeLangShort        = "Ljava/lang/Short;"
	javaTypeLangInteger      = "Ljava/lang/Integer;"
	javaTypeLangLong         = "Ljava/lang/Long;"
	javaTypeLangFloat        = "Ljava/lang/Float;"
	javaTypeLangDouble       = "Ljava/lang/Double;"
	javaTypeUtilDate         = "Ljava/util/Date;"
	javaTypeUtilList         = "Ljava/util/List;"
	javaTypeLinkedList       = "Ljava/util/LinkedList;"
	javaTypeArray            = "Ljava/util/ArrayList;"
	javaTypeUtilMap          = "Ljava/util/Map;"
	javaTypeUtilHashMap      = "Ljava/util/HashMap;"
	javaTypeLinkedHashMap    = "Ljava/util/LinkedHashMap;"
	javaTypeUtilSet          = "Ljava/util/Set;"
	javaTypeUtilHashSet      = "Ljava/util/HashSet;"
	javaTypeBigDecimal       = "Ljava/math/BigDecimal;"
	javaTypeBigInteger       = "Ljava/math/BigInteger;"
	javaTypeEnum             = "Ljava/lang/Enum;"
	javaTypeLangVoid         = "Ljava/lang/Void;"
	javaTypeCollection       = "Ljava/util/Collection;"
	javaTypeRuntimeException = "Ljava/lang/RuntimeException;"
)

var (
	javaTypeToGo = map[string]string{
		javaTypeLangVoid:    "interface{}",
		javaTypeObject:      "interface{}",
		javaTypeChar:        "int16",
		javaTypeBoolean:     "bool",
		javaTypeByte:        "int8",
		javaTypeInt:         "int32",
		javaTypeShort:       "int16",
		javaTypeLong:        "int64",
		javaTypeFloat:       "float32",
		javaTypeDouble:      "float64",
		javaTypeLangString:  "string",
		javaTypeLangBoolean: "bool",
		javaTypeLangInteger: "int32",
		javaTypeLangLong:    "int64",
		javaTypeLangDouble:  "float64",
		javaTypeLangByte:    "*common.Byte",
		javaTypeLangShort:   "*common.Short",
		javaTypeLangFloat:   "*common.Float",
		javaTypeUtilDate:    "time.Time",
		javaTypeUtilSet:     "*common.Set",
		javaTypeUtilHashSet: "*common.HashSet",
		javaTypeBigDecimal:  "*gxbig.Decimal",
		javaTypeBigInteger:  "*gxbig.Integer",
	}
	baseJavaTypeToGoQuote = map[string]string{
		javaTypeLangString:  "*string",
		javaTypeLangBoolean: "*bool",
		javaTypeLangInteger: "*int32",
		javaTypeLangLong:    "*int64",
		javaTypeLangDouble:  "*float64",
	}

	ErrNotService = errors.New("class is not a service")
)

const (
	targetFunc = "func (s *%s) WithTarget(targets []string) *%s {\n" +
		"\tif len(targets) == 0 {\n" +
		"\t\treturn s\n" +
		"\t}\n" +
		"\tcloneS := &%s{\n" +
		"\t\topt: s.opt,\n" +
		"\t}\n" +
		"\tcloneS.c = consumer.GetConsumerWithTarget(&registry.ServiceConfig{\n" +
		"\t\tService: s.opt.service,\n" +
		"\t\tGroup:   s.opt.group,\n" +
		"\t\tVersion: s.opt.version,\n" +
		"\t}, targets)\n" +
		"\treturn cloneS\n" +
		"}\n\n"

	consumerReq = "\tctx = filterCtx(ctx, s.opt)\n" +
		"\treq := &codec.HsfRequestPackage{\n" +
		"\t\tType:    codec.RpcRequest,\n" +
		"\t\tGroup:   s.opt.group,\n" +
		"\t\tVersion: s.opt.version,\n" +
		"\t\tService: s.opt.service,\n" +
		"\t\tMethod:  \"%s\",\n" +
		"\t\tArgSigs: []string{\n%s\n},\n" +
		"\t\tArgs:    []interface{}{%s},\n" +
		"\t\tProps:   make(map[string]interface{}),\n" +
		"\t\tTimeout: ctx.Value(consumer.RequestTimeoutKey).(time.Duration),\n" +
		"\t}\n"

	consumerNewFunc = "func New%s(option ...Option) (*%s, error) {\n" +
		"\topt := Options{}\n" +
		"\tfor _, o := range option {\n" +
		"\t\to(&opt)\n" +
		"\t}\n" +
		"\tif len(opt.service) == 0 {\n" +
		"\t\treturn nil, errors.New(\"service is empty\")\n" +
		"\t}\n" +
		"\tif len(opt.group) == 0 {\n" +
		"\t\treturn nil, errors.New(\"group is empty\")\n" +
		"\t}\n" +
		"\tif len(opt.version) == 0 {\n" +
		"\t\treturn nil, errors.New(\"version is empty\")\n" +
		"\t}\n\tif opt.requestTimeout <= 0 {\n" +
		"\t\topt.requestTimeout = 3 * time.Second\n" +
		"\t}\n" +
		"\treturn &%s{\n" +
		"\t\topt: opt,\n" +
		"\t\tc: consumer.GetConsumer(&registry.ServiceConfig{\n" +
		"\t\t\tService: opt.service,\n" +
		"\t\t\tGroup:   opt.group,\n" +
		"\t\t\tVersion: opt.version,\n\t\t}),\n" +
		"\t}, nil\n" +
		"}\n\n"

	providerNewFunc = "func Register%s(p %s) error {\n" +
		"}\n\n"
)

var baseType = map[int32]struct{}{
	'C': {},
	'Z': {},
	'B': {},
	'I': {},
	'S': {},
	'J': {},
	'F': {},
	'D': {},
}

const (
	minInt32 = ^int32(^uint32(0) >> 1)
)

func isEnd(v, next int32) bool {
	if next == minInt32 {
		return true
	}
	if v == ';' {
		return true
	}
	_, vIsBase := baseType[v]
	_, nextIsBase := baseType[next]
	if vIsBase && (nextIsBase || next == 'L') {
		return true
	}
	return false
}
func methodSignature(signature string) (input []string, output []string) {
	if len(signature) == 0 {
		return
	}
	var (
		t    = 1
		cnt  = 0
		name = ""
	)
	for i, v := range signature {
		if v == ' ' {
			continue
		}
		if v == '(' {
			continue
		}
		if v == ')' {
			t = 2
			if len(name) != 0 {
				input = append(input, name)
			}
			name = ""
			continue
		}
		next := minInt32
		if i != len(signature)-1 {
			next = int32(signature[i+1])
		}
		if isEnd(v, next) && cnt == 0 {
			if t == 1 {
				input = append(input, name+string(v))
			} else {
				output = append(output, name+string(v))
			}
			name = ""
			continue
		}
		if v == '<' {
			cnt++
		}
		if v == '>' {
			cnt--
		}
		name += string(v)
	}
	return
}

func goBaseType(tr Transform, name string) (string, bool) {
	if tr.Opt().OutQuote {
		s, ok := baseJavaTypeToGoQuote[name]
		if ok {
			return s, ok
		}
	}
	s, ok := javaTypeToGo[name]
	if !ok {
		return "", false
	}
	switch name {
	case javaTypeBigDecimal, javaTypeBigInteger:
		tr.AppendImport("github.com/dubbogo/gost/math/big")
	case javaTypeUtilDate:
		tr.AppendImport("time")
	}
	return s, true
}

var (
	csStructOnce sync.Once
	csStructMap  map[string]config.CustomStruct
)

func customStruct(className string) (goName, importPath string, have bool) {
	className = formatJavaName(className)
	className = className[1 : len(className)-1]
	className = strings.ReplaceAll(className, "/", ".")
	csStructOnce.Do(func() {
		csStructMap = make(map[string]config.CustomStruct)
		csStructList := config.Config().Hsf.CustomStruct
		for _, cs := range csStructList {
			csStructMap[cs.JavaClassName] = cs
		}
	})
	if v, ok := csStructMap[className]; ok {
		nameArr := strings.Split(v.GoStructName, ".")
		if len(nameArr) == 0 {
			panic(".gokit.toml syntax error, customStruct->[" + className + "]->goStructName")
		}
		if len(nameArr[0]) <= 0 ||
			nameArr[0][0] < 'A' || nameArr[0][0] > 'Z' {
			panic(".gokit.toml syntax error, customStruct->[" + className + "]->goStructName")
		}
		importPathArr := strings.Split(v.GoImportPath, "/")
		if len(importPathArr) == 0 {
			panic(".gokit.toml syntax error, customStruct->[" + className + "]->goImportPath")
		}
		goName = "*" + importPathArr[len(importPathArr)-1] +
			"." + nameArr[len(nameArr)-1]
		importPath = v.GoImportPath
		have = true
	} else {
		have = false
	}
	return
}

const (
	// StructTypeInit 默认值
	StructTypeInit = 0
	// StructTypeDefault struct普通类型
	StructTypeDefault = 1
	// StructTypeEnum struct枚举类型
	StructTypeEnum = 2
)

type StructInfo struct {
	GroupID    string
	ArtifactID string
	Version    string
	Dependency map[string]struct{}
	Name       string
	Type       int
	buf        []byte
}

// Transform 生成器
type Transform interface {
	Opt() Opt
	ZipEntry() *classpath.ZipEntry
	CreateImportPath(artifactId string) string
	CreateStructName(name string) string
	ExistStruct(name string) (StructInfo, bool)
	AppendStruct(name string, data StructInfo)
	ExistConsumer(name string) bool
	AppendConsumer(name string, data StructInfo)
	ExistProvider(name string) bool
	AppendProvider(name string, data StructInfo)
	AppendImport(name string)
}

func transP(c multiCtx, p []P) string {
	f := func(tmpP P) string {
		switch tmpP.Name {
		case javaTypeLangNumber:
			panic("not support number")
		case javaTypeUtilList, javaTypeLinkedList, javaTypeSlice, javaTypeCollection, javaTypeArray:
			if len(tmpP.Sub) != 1 {
				panic("list sub must 1")
			}
			return "[]" + transP(c, tmpP.Sub)
		case javaTypeUtilMap, javaTypeUtilHashMap, javaTypeLinkedHashMap:
			if len(tmpP.Sub) < 1 {
				panic("map sub must >= 1")
			}
			return "map" + transP(c, tmpP.Sub)
		default:
			name, ok := goBaseType(c.tr, tmpP.Name)
			if ok {
				return name
			}
			// 处理泛型
			if len(tmpP.Sub) != 0 {
				_ = transP(c, tmpP.Sub)
			}
			goName, artifactID := transToGoStruct(c.tr, tmpP.Name)
			if len(artifactID) != 0 {
				c.dependency[artifactID] = struct{}{}
			}
			return goName
		}
	}
	if len(p) == 1 {
		return f(p[0])
	}
	return "[" + f(p[0]) + "]" + f(p[1])
}

type multiCtx struct {
	tr         Transform
	dependency map[string]struct{}
}

func transMultiToGoStruct(c multiCtx, className string) (name string, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("%+v", pErr)
		}
	}()
	name = transP(c, parseParam(className))
	return
}

// transToGoStruct 转换为go类型
func transToGoStruct(tr Transform, className string) (string, string) {
	defer func() {
		if pErr := recover(); pErr != nil {
			fmt.Printf("\033[31m[toolkit] className: %s error， %v \033[0m\n", className, pErr)
		}
	}()
	// 处理成标准数据
	className = formatJavaName(className)
	goName, ok := goBaseType(tr, className)
	if ok {
		return goName, ""
	}
	goName, importPath, ok := customStruct(className)
	if ok {
		return goName, importPath
	}
	if data, ok := tr.ExistStruct(className); ok {
		return "*" + pkgName(data.ArtifactID) +
			"." + data.Name, tr.CreateImportPath(data.ArtifactID)
	}

	goStructName := tr.CreateStructName(className)

	class, err := getClassFile(tr.ZipEntry(), className)
	if err != nil {
		panic(err)
	}

	// 处理循环
	tr.AppendStruct(className, StructInfo{
		GroupID:    class.GroupID,
		ArtifactID: class.ArtifactID,
		Version:    class.Version,
		Type:       StructTypeInit,
		Name:       goStructName,
	})

	superClass := class.File.SuperClassName()
	superClassList := make([]*classfile.ClassFile, 0)
	for {
		if !isSupperObject(superClass) {
			classInfo, err := getClassFile(tr.ZipEntry(), superClass)
			if err != nil {
				panic(err)
			}
			superClassList = append([]*classfile.ClassFile{classInfo.File},
				superClassList...)
			superClass = classInfo.File.SuperClassName()
		} else {
			break
		}
	}

	var (
		idx           = 0
		fieldKeyIndex = make(map[string]int)
		fieldMap      = make(map[string]*classfile.MemberInfo)
	)

	for _, info := range superClassList {
		for _, v := range info.Fields() {
			if v.AccessFlags()&0x8 == 0x8 { // final, static
				continue
			}
			fieldKeyIndex[v.Name()] = idx
			fieldMap[v.Name()] = v
			idx++
		}
	}
	for _, v := range class.File.Fields() {
		if v.AccessFlags()&0x8 == 0x8 { // final, static
			continue
		}
		fieldKeyIndex[v.Name()] = idx
		fieldMap[v.Name()] = v
		idx++
	}

	var (
		idxList     = make([]int, 0)
		fieldKeyMap = make(map[int]string)
	)
	for fieldName, idx := range fieldKeyIndex {
		fieldKeyMap[idx] = fieldName
		idxList = append(idxList, idx)
	}
	sort.Ints(idxList)
	keys := make([]string, 0, len(fieldMap))
	for _, idx := range idxList {
		keys = append(keys, fieldKeyMap[idx])
	}

	sByte := bytes.NewBuffer(nil)
	sByte.WriteString(fmt.Sprintf("type %s struct {\n", goStructName))

	fullClassName := strings.ReplaceAll(class.File.ClassName(), "/", ".")

	if class.File.AccessFlags()&0x4000 == 0x4000 { // 处理枚举
		sByte.WriteString("\tName string")
		sByte.WriteString("}\n\n")
		sByte.WriteString("func (*" + goStructName + ") JavaClassName() string {\n")
		sByte.WriteString("\treturn \"" + fullClassName + "\"\n")
		sByte.WriteString("}\n\n")
		sByte.WriteString("func (e *" + goStructName + ") String() string {\n")
		sByte.WriteString("\tif e == nil {\n")
		sByte.WriteString("\t\treturn \"\"\n")
		sByte.WriteString("\t}\n")
		sByte.WriteString("\treturn e.Name\n")
		sByte.WriteString("}\n\n")
		sByte.WriteString("func (e *" + goStructName + ") EnumValue(s string) hessian.JavaEnum {\n")
		sByte.WriteString("\treturn map[string]hessian.JavaEnum{\n")
		var i int
		for _, info := range class.File.Fields() {
			if info.Name() == "$VALUES" {
				continue
			}
			if !strings.Contains(info.Descriptor(), class.File.ClassName()) {
				continue
			}
			i++
			sByte.WriteString(fmt.Sprintf("\t\t\"%s\": %d,\n", info.Name(), i))
		}
		sByte.WriteString("\t}[s]\n")
		sByte.WriteString("}\n\n")
		// 写入枚举类型
		sByte.WriteString("var (\n")
		for _, info := range class.File.Fields() {
			if info.Name() == "$VALUES" {
				continue
			}
			if !strings.Contains(info.Descriptor(), class.File.ClassName()) {
				continue
			}
			i++
			sByte.WriteString(fmt.Sprintf("\t%s%s = &%s{\"%s\"}\n", goStructName,
				upperFirstWord(transSpecialField(info.Name())), goStructName, info.Name()))
		}
		sByte.WriteString(")\n")
		sByte.WriteString(fmt.Sprintf("func New%s(s string) *%s {\n",
			goStructName, goStructName))
		sByte.WriteString(fmt.Sprintf("\treturn &%s{\n\t\tName: s,\n\t}\n", goStructName))
		sByte.WriteString("}\n")
		tr.AppendStruct(className, StructInfo{
			GroupID:    class.GroupID,
			ArtifactID: class.ArtifactID,
			Version:    class.Version,
			Type:       StructTypeEnum,
			Name:       goStructName,
			buf:        sByte.Bytes(),
		})

		return "*" + pkgName(class.ArtifactID) +
			"." + goStructName, tr.CreateImportPath(class.ArtifactID)
	}
	dependency := make(map[string]struct{})
	for _, key := range keys {
		field := fieldMap[key]
		fieldDes := field.Descriptor()
		m, ok := goBaseType(tr, fieldDes)
		if !ok {
			if field.SignatureAttribute() != nil {
				fieldDes = field.SignatureAttribute().Signature()
			}
			c := multiCtx{
				tr:         tr,
				dependency: make(map[string]struct{}),
			}
			m, err = transMultiToGoStruct(c, fieldDes)
			if err != nil {
				fmt.Printf("\033[31m[toolkit] className: %s, fieldName: %s, fieldType: %s continue\033[0m\n",
					className, field.Name(), fieldDes)
				continue
			}
			m = strings.ReplaceAll(m,
				"*"+pkgName(class.ArtifactID)+".", "*")
			for dpKey := range c.dependency {
				dependency[dpKey] = struct{}{}
			}
		}
		notes := ""
		if field.AccessFlags()&0x0010 == 0x0010 { // final
			notes = "// final"
		}

		var (
			hessianTag = ""
			fieldName  = upperFirstWord(field.Name())
		)
		fieldName = lowerFirstWord(transSpecialField(field.Name()))
		if fieldName != field.Name() {
			hessianTag = fmt.Sprintf("`json:\"%s\" hessian:\"%s\"`",
				field.Name(), field.Name())
		} else {
			hessianTag = fmt.Sprintf("`json:\"%s\"`",
				field.Name())
		}
		fieldName = upperFirstWord(fieldName)
		sByte.WriteString(fmt.Sprintf("\t%s %s %s %s\n",
			fieldName, m, hessianTag, notes))
	}
	sByte.WriteString("}\n\n")
	sByte.WriteString("func (*" + goStructName + ") JavaClassName() string {\n")
	sByte.WriteString("\treturn \"" + fullClassName + "\"\n")
	sByte.WriteString("}\n\n")

	tr.AppendStruct(className, StructInfo{
		GroupID:    class.GroupID,
		ArtifactID: class.ArtifactID,
		Version:    class.Version,
		Dependency: dependency,
		Type:       StructTypeDefault,
		Name:       goStructName,
		buf:        sByte.Bytes(),
	})

	return "*" + pkgName(class.ArtifactID) +
		"." + goStructName, tr.CreateImportPath(class.ArtifactID)
}

type ParamDefine struct {
	GoClassName   string
	JavaClassName string
}

type MethodDefine struct {
	GoName   string
	JavaName string
	In       []ParamDefine
	Out      []ParamDefine
}

// appendConsumer 插入一个service
func appendConsumer(tr Transform, className string) error {
	// 处理成标准数据
	className = formatJavaName(className)
	if tr.ExistConsumer(className) {
		return nil
	}

	fmt.Printf("[toolkit] generate consumer %s starting...\n", className)

	class, err := getClassFile(tr.ZipEntry(), className)
	if err != nil {
		panic(err)
	}

	if class.File.AccessFlags()&0x0200 != 0x0200 {
		return ErrNotService
	}
	var (
		goNameMap = make(map[string]struct{})
		methods   = make([]MethodDefine, 0)
	)
	dependency := make(map[string]struct{})
	for _, v := range class.File.Methods() {
		if v.AccessFlags()&0x8 == 0x8 {
			continue
		}
		var (
			i      int
			goName = v.Name()
		)
		for {
			if _, ok := goNameMap[goName]; ok {
				i++
				goName = fmt.Sprintf("%s%d", v.Name(), i)
			} else {
				goNameMap[goName] = struct{}{}
				break
			}
		}
		method := MethodDefine{
			GoName:   goName,
			JavaName: v.Name(),
			In:       make([]ParamDefine, 0),
			Out:      make([]ParamDefine, 0),
		}
		signature := v.Descriptor()
		if v.SignatureAttribute() != nil {
			signature = v.SignatureAttribute().Signature()
		}
		in, out := methodSignature(signature)
		for _, javaName := range in {
			c := multiCtx{
				tr:         tr,
				dependency: make(map[string]struct{}),
			}
			goName, err := transMultiToGoStruct(c, javaName)
			if err != nil {
				fmt.Printf("\033[31m[toolkit] className: %s, methodName: %s, fieldType: %s in continue\033[0m\n",
					className, v.Name(), javaName)
				continue
			}
			for dpKey := range c.dependency {
				dependency[dpKey] = struct{}{}
			}
			method.In = append(method.In, ParamDefine{
				goName,
				javaName,
			})
		}
		for _, javaName := range out {
			if javaName == javaTypeVoid {
				continue
			}
			c := multiCtx{
				tr:         tr,
				dependency: make(map[string]struct{}),
			}
			goName, err := transMultiToGoStruct(c, javaName)
			if err != nil {
				fmt.Printf("\033[31m[toolkit] className: %s, methodName: %s, fieldType: %s out continue\033[0m\n",
					className, v.Name(), javaName)
				continue
			}
			for dpKey := range c.dependency {
				dependency[dpKey] = struct{}{}
			}
			method.Out = append(method.Out, ParamDefine{
				goName,
				javaName,
			})
		}
		methods = append(methods, method)
	}
	if len(methods) == 0 {
		fmt.Printf("[toolkit] generate consumer %s end, is empty...\n", className)
		return nil
	}

	goStructName := tr.CreateStructName(className)

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("var _ %sInterface = &%s{}\n\n",
		goStructName, goStructName))
	buf.WriteString(fmt.Sprintf("type %sInterface interface {\n", goStructName))
	for _, method := range methods {
		buf.WriteString("\t")
		buf.WriteString(upperFirstWord(method.GoName))
		buf.WriteString("(ctx context.Context")
		for inK, inDefine := range method.In {
			buf.WriteString(" ,p" + strconv.Itoa(inK+1) + " " + inDefine.GoClassName)
		}
		buf.WriteString(")")
		if len(method.Out) == 0 {
			buf.WriteString("error\n")
		} else {
			buf.WriteString(" (")
			for _, outDefine := range method.Out {
				buf.WriteString(outDefine.GoClassName + ", ")
			}
			buf.WriteString("error)\n")
		}
	}
	buf.WriteString("}\n\n")
	buf.WriteString(fmt.Sprintf("// %s %s\n",
		goStructName, toJavaClassName(className)))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", goStructName))
	buf.WriteString("\topt Options\n")
	buf.WriteString("\tc   *consumer.Consumer\n")
	buf.WriteString("}\n\n")
	buf.WriteString(fmt.Sprintf("// New%s 创建客户端实例，使用时必须是单例的\n", goStructName))
	buf.WriteString(fmt.Sprintf(consumerNewFunc, goStructName, goStructName, goStructName))
	buf.WriteString("// WithTarget 自定义ip访问，用于本地调试\n")
	buf.WriteString(fmt.Sprintf(targetFunc, goStructName, goStructName, goStructName))
	for _, method := range methods {
		buf.WriteString(fmt.Sprintf("func (s *%s) %s(ctx context.Context",
			goStructName, upperFirstWord(method.GoName)))
		for inK, inDefine := range method.In {
			buf.WriteString(" ,p" + strconv.Itoa(inK+1) + " " + inDefine.GoClassName)
		}
		buf.WriteString(")")
		if len(method.Out) == 0 {
			buf.WriteString("error {\n")
		} else {
			buf.WriteString(" (")
			for _, outDefine := range method.Out {
				buf.WriteString(outDefine.GoClassName + ", ")
			}
			buf.WriteString("error) {\n")
		}
		var (
			args    = make([]string, 0)
			argSigs = make([]string, 0)
		)
		for inK, inDefine := range method.In {
			args = append(args, "p"+strconv.Itoa(inK+1))
			javaSign := toJavaClassName(inDefine.JavaClassName)
			javaSign = removeGenericName(javaSign)
			argSigs = append(argSigs, "\""+javaSign+"\",\n")
		}
		buf.WriteString(fmt.Sprintf(consumerReq,
			method.JavaName, strings.Join(argSigs, ""), strings.Join(args, ", ")))
		if len(method.Out) == 0 {
			buf.WriteString("\t_, err := s.c.CallSync(ctx, req)\n")
			buf.WriteString("\treturn err")
		} else {
			// java没有多返回值，暂时这么处理
			goOut := method.Out[0].GoClassName
			buf.WriteString(fmt.Sprintf("\tres, err := s.c.CallSync(ctx, req)\n"+
				"\tif err != nil {\n\t\treturn %s, err\n\t}\n"+
				"\tresT, ok := res.(%s)\n"+
				"\tif !ok {\n\t\treturn %s, errors.New(\"res type error\")\n\t}\n"+
				"\treturn resT, nil\n", goDefaultVal(goOut), goOut, goDefaultVal(goOut)))
		}
		buf.WriteString("}\n\n")
	}

	tr.AppendConsumer(className, StructInfo{
		GroupID:    class.GroupID,
		ArtifactID: class.ArtifactID,
		Version:    class.Version,
		Dependency: dependency,
		Name:       goStructName,
		buf:        buf.Bytes(),
	})
	fmt.Printf("[toolkit] generate consumer %s end...\n", className)
	return nil
}

// appendProvider 插入一个provider service
func appendProvider(tr Transform, className string) error {
	// 处理成标准数据
	className = formatJavaName(className)
	if tr.ExistProvider(className) {
		return nil
	}

	fmt.Printf("[toolkit] generate provider %s starting...\n", className)

	class, err := getClassFile(tr.ZipEntry(), className)
	if err != nil {
		panic(err)
	}

	if class.File.AccessFlags()&0x0200 != 0x0200 {
		return ErrNotService
	}
	var (
		goNameMap = make(map[string]struct{})
		methods   = make([]MethodDefine, 0)
	)
	dependency := make(map[string]struct{})
	for _, v := range class.File.Methods() {
		if v.AccessFlags()&0x8 == 0x8 {
			continue
		}
		var (
			i      int
			goName = v.Name()
		)
		for {
			if _, ok := goNameMap[goName]; ok {
				i++
				goName = fmt.Sprintf("%s%d", v.Name(), i)
			} else {
				goNameMap[goName] = struct{}{}
				break
			}
		}
		method := MethodDefine{
			GoName:   goName,
			JavaName: v.Name(),
			In:       make([]ParamDefine, 0),
			Out:      make([]ParamDefine, 0),
		}
		signature := v.Descriptor()
		if v.SignatureAttribute() != nil {
			signature = v.SignatureAttribute().Signature()
		}
		in, out := methodSignature(signature)
		for _, javaName := range in {
			c := multiCtx{
				tr:         tr,
				dependency: make(map[string]struct{}),
			}
			goName, err := transMultiToGoStruct(c, javaName)
			if err != nil {
				fmt.Printf("\033[31m[toolkit] className: %s, methodName: %s, fieldType: %s in params continue\033[0m\n",
					className, v.Name(), javaName)
				continue
			}
			for dpKey := range c.dependency {
				dependency[dpKey] = struct{}{}
			}
			method.In = append(method.In, ParamDefine{
				goName,
				javaName,
			})
		}
		for _, javaName := range out {
			if javaName == javaTypeVoid {
				continue
			}
			c := multiCtx{
				tr:         tr,
				dependency: make(map[string]struct{}),
			}
			goName, err := transMultiToGoStruct(c, javaName)
			if err != nil {
				fmt.Printf("\033[31m[toolkit] className: %s, methodName: %s, fieldType: %s out params continue\033[0m\n",
					className, v.Name(), javaName)
				continue
			}
			for dpKey := range c.dependency {
				dependency[dpKey] = struct{}{}
			}
			method.Out = append(method.Out, ParamDefine{
				goName,
				javaName,
			})
		}
		methods = append(methods, method)
	}
	if len(methods) == 0 {
		fmt.Printf("[toolkit] generate provider %s end, is empty...\n", className)
		return nil
	}

	goStructName := tr.CreateStructName(className)

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("// %sProvider %s\n",
		goStructName, toJavaClassName(className)))
	buf.WriteString(fmt.Sprintf("type %sProvider interface {\n", goStructName))
	for _, method := range methods {
		buf.WriteString("\t")
		buf.WriteString(upperFirstWord(method.GoName))
		buf.WriteString("(ctx context.Context")
		for inK, inDefine := range method.In {
			buf.WriteString(" ,p" + strconv.Itoa(inK+1) + " " + inDefine.GoClassName)
		}
		buf.WriteString(")")
		if len(method.Out) == 0 {
			buf.WriteString("error\n")
		} else {
			buf.WriteString(" (")
			for _, outDefine := range method.Out {
				buf.WriteString(outDefine.GoClassName + ", ")
			}
			buf.WriteString("error)\n")
		}
	}
	buf.WriteString("\tInterfaceName() string \n")
	buf.WriteString("\tVersion() string \n")
	buf.WriteString("\tGroup() string \n")
	buf.WriteString("}\n\n")
	buf.WriteString("")
	tr.AppendProvider(className, StructInfo{
		GroupID:    class.GroupID,
		ArtifactID: class.ArtifactID,
		Version:    class.Version,
		Dependency: dependency,
		Name:       goStructName,
		buf:        buf.Bytes(),
	})
	fmt.Printf("[toolkit] generate provider %s end...\n", className)
	return nil
}

// isSupperObject 是否是基础类
func isSupperObject(className string) bool {
	className = formatJavaName(className)
	switch className {
	case javaTypeObject, javaTypeEnum, javaTypeRuntimeException,
		javaTypeArray, javaTypeCollection, javaTypeUtilHashSet:
		return true
	default:
		return false
	}
}

// upperFirstWord 首字母大写
func upperFirstWord(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func lowerFirstWord(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToLower(s[:1]) + s[1:]
}

var specialMap = map[string]string{
	"ACL":   "ACL",
	"API":   "API",
	"ASCII": "ASCII",
	"CPU":   "CPU",
	"CSS":   "CSS",
	"DNS":   "DNS",
	"EOF":   "EOF",
	"GUID":  "GUID",
	"HTML":  "HTML",
	"HTTP":  "HTTP",
	"HTTPS": "HTTPS",
	"ID":    "ID",
	"IP":    "IP",
	"JSON":  "JSON",
	"JSONP": "JSONP",
	"LHS":   "LHS",
	"QPS":   "QPS",
	"RAM":   "RAM",
	"RHS":   "RHS",
	"RPC":   "RPC",
	"SLA":   "SLA",
	"SMTP":  "SMTP",
	"SQL":   "SQL",
	"SSH":   "SSH",
	"TCP":   "TCP",
	"TLS":   "TLS",
	"TTL":   "TTL",
	"UDP":   "UDP",
	"UI":    "UI",
	"UID":   "UID",
	"UUID":  "UUID",
	"URI":   "URI",
	"URL":   "URL",
	"UTF8":  "UTF8",
	"VM":    "VM",
	"XML":   "XML",
	"XMPP":  "XMPP",
	"XSRF":  "XSRF",
	"XSS":   "XSS",

	"ECS":  "ECS",
	"VPC":  "VPC",
	"OSS":  "OSS",
	"SLB":  "SLB",
	"RDS":  "RDS",
	"CDN":  "CDN",
	"HPC":  "HPC",
	"NAT":  "NAT",
	"ICMP": "ICMP",
	"GW":   "GW",
	"EDAS": "EDAS",
	"DRDS": "DRDS",
	"ARMS": "ARMS",
	"MQ":   "MQ",
	"CSB":  "CSB",
	"TS":   "TS",
}

func splitName(name string) []string {
	nameItem := ""
	nameSlice := make([]string, 0)
	for i, r := range name {
		if (unicode.IsUpper(r) || r == '_') && i != 0 {
			nameSlice = append(nameSlice, upperFirstWord(nameItem))
			if r != '_' {
				nameItem = string(r)
			} else {
				nameItem = ""
			}
			continue
		}
		nameItem += string(r)
	}
	if len(nameItem) != 0 {
		nameSlice = append(nameSlice, upperFirstWord(nameItem))
	}
	return nameSlice
}

func transSpecialField(name string) string {
	nameSlice := splitName(name)
	tmpSlice := make([]string, 0, len(nameSlice))
	for _, s := range nameSlice {
		if v, ok := specialMap[strings.ToUpper(s)]; ok {
			tmpSlice = append(tmpSlice, v)
			continue
		}
		tmpSlice = append(tmpSlice, s)
	}
	return strings.Join(tmpSlice, "")
}

func toJavaClassName(className string) string {
	if v, ok := map[string]string{
		javaTypeChar:    "char",
		javaTypeBoolean: "boolean",
		javaTypeByte:    "byte",
		javaTypeInt:     "int",
		javaTypeShort:   "short",
		javaTypeLong:    "long",
		javaTypeFloat:   "float",
		javaTypeDouble:  "double",
	}[className]; ok {
		return v
	}
	if strings.HasPrefix(className, "L") &&
		strings.HasSuffix(className, ";") {
		className = className[1 : len(className)-1]
	}
	if strings.HasSuffix(className, ".class") {
		className = className[:len(className)-6]
	}
	return strings.ReplaceAll(className, "/", ".")
}

type P struct {
	Name string
	Sub  []P
}

func parseParam(paramName string) []P {
	p := P{}
	reP := make([]P, 0)
	for k, v := range paramName {
		if v == '<' {
			p.Sub = parseParam(paramName[k+1 : len(paramName)-2])
			break
		}
		if v == '[' {
			p.Name = "["
			p.Sub = parseParam(paramName[1:])
			break
		}
		if v == ';' {
			p.Name += string(v)
			reP = append(reP, p)
			p = P{}
			continue
		}
		p.Name += string(v)
	}
	if len(p.Name) != 0 {
		p.Name = formatJavaName(p.Name)
		reP = append(reP, p)
	}
	return reP
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		d := s[i]
		if (i > 0 && d >= 'A' && d <= 'Z') &&
			(!(i+1 < len(s) && s[i+1] >= 'A' && s[i+1] <= 'Z') ||
				!(i-1 > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z')) {
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func formatJavaName(name string) string {
	if _, ok := map[string]bool{
		javaTypeChar:    true,
		javaTypeBoolean: true,
		javaTypeByte:    true,
		javaTypeInt:     true,
		javaTypeShort:   true,
		javaTypeLong:    true,
		javaTypeFloat:   true,
		javaTypeDouble:  true,
		javaTypeSlice:   true,
	}[name]; ok {
		return name
	}
	if !strings.HasPrefix(name, "L") {
		name = "L" + name
	}
	if !strings.HasSuffix(name, ";") {
		name += ";"
	}
	return name
}

type ClassFile struct {
	GroupID    string
	ArtifactID string
	Version    string
	File       *classfile.ClassFile
}

var (
	classFileMap = make(map[string]*ClassFile)
)

// getClassFile 获取类信息
func getClassFile(ze *classpath.ZipEntry, className string) (*ClassFile, error) {
	// 处理成标准数据
	if strings.HasPrefix(className, "L") &&
		strings.HasSuffix(className, ";") {
		className = className[1 : len(className)-1]
	}
	if strings.HasSuffix(className, ".class") {
		className = className[:len(className)-6]
	}

	if d, ok := classFileMap[className]; ok {
		return d, nil
	}
	classData, err := classpath.ReadClass(ze, className)
	if err != nil {
		return nil, errors.New(className + " " + err.Error())
	}
	cf, err := classfile.Parse(classData.File)
	if err != nil {
		return nil, err
	}
	classFileMap[className] = &ClassFile{
		GroupID:    classData.GroupID,
		ArtifactID: classData.ArtifactID,
		Version:    classData.Version,
		File:       cf,
	}
	return classFileMap[className], nil
}

func pkgName(artifactID string) string {
	var d string
	for _, v := range artifactID {
		if (v >= 'a' && v <= 'z') ||
			(v >= 'A' && v <= 'Z') {
			d += string(v)
		}
	}
	d = strings.TrimSuffix(d, "api")
	d = strings.TrimPrefix(d, "amap")
	d = strings.TrimPrefix(d, "aos")
	return d
}

func goDefaultVal(name string) string {
	switch name {
	case "int8", "int16", "int32", "int64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "string":
		return "\"\""
	default:
		return "nil"
	}
}

func removeGenericName(name string) string {
	i := strings.Index(name, "<")
	if i < 0 {
		return name
	}
	return name[:i]
}
