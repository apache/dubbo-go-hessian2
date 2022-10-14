package hsf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/classpath"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/config"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/mod"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/util"
	"github.com/go-playground/validator/v10"
)

const (
	outTypeConsumer  = "consumer"
	outTypeProvider  = "provider"
	outTypeStruct    = "struct"
	outTypeListClass = "listClass"
	hsfModel         = "gitlab.alibaba-inc.com/amap-go/hsf-go"
	hsfVersion       = "v0.1.23"
	bigModel         = "github.com/dubbogo/gost"
	bigVersion       = "v1.10.1"

	hessianOld     = "github.com/apache/dubbo-go-hessian2"
	hessianNew     = "gitlab.alibaba-inc.com/amap-go/hessian2-go"
	hessianVersion = "v0.0.15"
	//hessianVersion = "replace github.com/apache/dubbo-go-hessian2 => gitlab.alibaba-inc.com/amap-go/hessian2-go v0.0.15"
)

const (
	optDef = "type Options struct {\n" +
		"\tservice        string\n" +
		"\tgroup          string\n" +
		"\tversion        string\n" +
		"\tappName        string\n" +
		"\trequestTimeout time.Duration\n" +
		"\tlogCfg         logger.TraceLoggerConfig\n" +
		"}\n\n"
	optFunc        = "type Option func(*Options)\n\n"
	optServiceFunc = "func WithService(service string) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.service = service\n" +
		"\t}" +
		"\n}\n\n"
	optGroupFunc = "func WithGroup(group string) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.group = group\n" +
		"\t}" +
		"\n}\n\n"
	optVersionFunc = "func WithVersion(version string) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.version = version\n" +
		"\t}" +
		"\n}\n\n"
	optAppNameFunc = "func WithAppName(appName string) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.appName = appName\n" +
		"\t}" +
		"\n}\n\n"
	optRequestTimeoutFunc = "func WithRequestTimeout(requestTimeout time.Duration) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.requestTimeout = requestTimeout\n" +
		"\t}" +
		"\n}\n\n"
	optLoggerFunc = "func WithLogger(logCfg logger.TraceLoggerConfig) Option {\n" +
		"\treturn func(options *Options) {\n" +
		"\t\toptions.logCfg = logCfg\n" +
		"\t}" +
		"\n}\n\n"
	filterFunc = "func filterCtx(ctx context.Context, opt Options) context.Context {\n" +
		"\tctxRequestTimeout := ctx.Value(consumer.RequestTimeoutKey)\n" +
		"\tif ctxRequestTimeout != nil {\n" +
		"\t\t_, ok := ctxRequestTimeout.(time.Duration)\n" +
		"\t\tif !ok {\n" +
		"\t\t\tctx = context.WithValue(ctx,\n" +
		"\t\t\t\tconsumer.RequestTimeoutKey, opt.requestTimeout)\n" +
		"\t\t}\n" +
		"\t} else {\n" +
		"\t\tctx = context.WithValue(ctx,\n" +
		"\t\t\tconsumer.RequestTimeoutKey, opt.requestTimeout)\n" +
		"\t}\n\n" +
		"\tctxAppName := ctx.Value(consumer.RequestAppName)\n" +
		"\tif ctxAppName == nil &&\n" +
		"\t\tlen(opt.appName) != 0 {\n" +
		"\t\tctx = context.WithValue(ctx, consumer.RequestAppName, opt.appName)\n" +
		"\t}\n\n" +
		"\tctxLogger := ctx.Value(consumer.TraceLoggerConfigKey)\n" +
		"\tif ctxLogger == nil {\n" +
		"\t\tctx = context.WithValue(ctx, consumer.TraceLoggerConfigKey, &opt.logCfg)\n" +
		"\t}\n" +
		"\treturn ctx" +
		"\n}\n\n"
)

type Opt struct {
	// in
	InJarPath    string
	InClass      string
	InGroupID    string `validate:"required"`
	InArtifactID string `validate:"required"`
	InVersion    string `validate:"required"`

	// out
	OutType      string `validate:"oneof=consumer provider idl struct listClass clean cleanAll"`
	OutPkg       string
	OutPath      string
	OutQuote     bool
	HaveInternal bool
}

type TransformImp struct {
	opt         Opt
	entry       *classpath.ZipEntry
	goStructMap map[string]struct{} // 已生成的goStruct名称
	structMap   map[string]StructInfo
	consumerMap map[string]StructInfo
	providerMap map[string]StructInfo
	importMap   map[string]struct{}
}

func (c *TransformImp) Opt() Opt {
	return c.opt
}

func (c *TransformImp) CreateStructName(className string) string {
	if len(className) == 0 {
		panic("className is empty")
	}

	if strings.HasPrefix(className, "L") &&
		strings.HasSuffix(className, ";") {
		className = className[1 : len(className)-1]
	}

	if strings.HasSuffix(className, ".class") {
		className = className[:len(className)-6]
	}

	classNameSplit := strings.Split(className, "/")
	if len(classNameSplit) == 0 {
		panic("class name empty")
	}
	idx := len(classNameSplit) - 1
	goName := strings.ReplaceAll(classNameSplit[idx], "$", "")
	goName = upperFirstWord(transSpecialField(goName))
	for {
		if _, ok := c.goStructMap[goName]; ok {
			idx--
			goName = upperFirstWord(transSpecialField(classNameSplit[idx])) + goName
		} else {
			break
		}
	}
	return goName
}

func (c *TransformImp) CreateImportPath(artifactId string) string {
	if c.opt.HaveInternal {
		return fmt.Sprintf("%s/%s/%s/%s",
			c.opt.OutPkg, "internal", "hsf", pkgName(artifactId))
	}
	return fmt.Sprintf("%s/%s/%s",
		c.opt.OutPkg, "hsf", pkgName(artifactId))
}

func (c *TransformImp) ExistStruct(name string) (StructInfo, bool) {
	data, ok := c.structMap[name]
	return data, ok
}

func (c *TransformImp) AppendStruct(name string, data StructInfo) {
	v, ok := c.structMap[name]
	if ok &&
		v.Type != StructTypeInit {
		return
	}
	c.goStructMap[data.Name] = struct{}{}
	c.structMap[name] = data
}

func (c *TransformImp) ExistConsumer(name string) bool {
	_, ok := c.consumerMap[name]
	return ok
}

func (c *TransformImp) AppendConsumer(name string, data StructInfo) {
	_, ok := c.consumerMap[name]
	if ok {
		return
	}
	c.goStructMap[data.Name] = struct{}{}
	c.consumerMap[name] = data
}

func (c *TransformImp) ExistProvider(name string) bool {
	_, ok := c.providerMap[name]
	return ok
}

func (c *TransformImp) AppendProvider(name string, data StructInfo) {
	_, ok := c.providerMap[name]
	if ok {
		return
	}
	c.goStructMap[data.Name] = struct{}{}
	c.providerMap[name] = data
}

func (c *TransformImp) AppendImport(name string) {
	c.importMap[name] = struct{}{}
}

type FileDes struct {
	name       string
	artifactID string
	buf        []byte
}

func (c *TransformImp) Model() []FileDes {
	structMap := make(map[string]map[string]StructInfo)
	for key, val := range c.structMap {
		if subMap, ok := structMap[val.ArtifactID]; ok {
			subMap[key] = val
		} else {
			structMap[val.ArtifactID] = map[string]StructInfo{
				key: val,
			}
		}
	}
	fileList := make([]FileDes, 0)
	for s, m := range structMap {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("package " + pkgName(s) + "\n\n")
		buf.WriteString("//Code generated by \"toolkit-go\"; DO NOT EDIT.\n")
		keys := make([]string, 0)
		for key := range m {
			keys = append(keys, key)
		}
		buf.WriteString("import (\n")
		buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/common\"\n")
		for s := range c.importMap {
			buf.WriteString("\t\"" + s + "\"\n")
		}
		buf.WriteString(")\n\n")
		sort.Strings(keys)
		buf.WriteString("func init()  {\n")
		for _, key := range keys {
			if m[key].Type == StructTypeEnum {
				buf.WriteString("\tcommon.RegisterEnum(&" + m[key].Name + "{})\n")
			} else {
				buf.WriteString("\tcommon.RegisterModel(&" + m[key].Name + "{})\n")
			}
		}
		buf.WriteString("}\n\n")
		for _, key := range keys {
			buf.Write(m[key].buf)
		}
		pomBuf := bytes.NewBuffer(nil)
		pomBuf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
		pomBuf.WriteString(fmt.Sprintf("<dependency>\n"+
			"\t <groupId>%s</groupId>\n"+
			"\t <artifactId>%s</artifactId>\n"+
			"\t <version>%s</version>\n"+
			"</dependency>\n", c.opt.InGroupID, c.opt.InArtifactID, c.opt.InVersion))
		fileList = append(fileList,
			FileDes{
				name:       "dep.xml",
				artifactID: s,
				buf:        pomBuf.Bytes(),
			})
		fileList = append(fileList, FileDes{
			name:       "model.go",
			artifactID: s,
			buf:        buf.Bytes(),
		})
	}

	return fileList
}

func (c *TransformImp) WriteModel() error {
	fileList := c.Model()
	for _, val := range fileList {
		outPath := path.Join(c.opt.OutPath, "hsf",
			pkgName(val.artifactID))
		if !util.FileExists(outPath) {
			err := os.MkdirAll(outPath, 0755)
			if err != nil {
				return err
			}
		}
		fileName := path.Join(outPath, val.name)
		err := ioutil.WriteFile(fileName, val.buf, 0766)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}
		cmd := exec.Command("goimports", "-l", "-w", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("[toolkit] 格式化文件报错，请检测goimports是否安装好，" + err.Error())
		}
		cmd = exec.Command("gofmt", "-l", "-w", fileName)
		err = cmd.Run()
		if err != nil {
			fmt.Println("格式化文件报错，请检测gofmt是否安装好，" + err.Error())
		}
	}

	return nil
}

func (c *TransformImp) Consumer(className string) []FileDes {
	pkg := "consumer"

	res := make([]FileDes, 0)
	buf := bytes.NewBuffer(nil)
	buf.WriteString("package " + pkg + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"context\"\n")
	buf.WriteString("\t\"time\"\n")
	buf.WriteString("\t\"errors\"\n")
	buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/codec\"\n")
	buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/consumer\"\n")
	buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/registry\"\n")
	buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/logger\"\n")
	buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/common\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString(optDef)
	buf.WriteString(optFunc)
	buf.WriteString(optServiceFunc)
	buf.WriteString(optGroupFunc)
	buf.WriteString(optVersionFunc)
	buf.WriteString(optAppNameFunc)
	buf.WriteString(optRequestTimeoutFunc)
	buf.WriteString(optLoggerFunc)
	buf.WriteString(filterFunc)
	res = append(res, FileDes{
		name: "options.go",
		buf:  buf.Bytes(),
	})
	if className == "all" {
		isService := func(className string) bool {
			if !strings.HasSuffix(className, ".class") {
				return false
			}
			class, err := getClassFile(c.ZipEntry(), className)
			if err != nil {
				panic(err)
			}
			if class.File.AccessFlags()&0x0200 == 0x0200 {
				return true
			}
			return false
		}
		file := c.entry.FileList()
		for _, s := range file {
			if isService(s) {
				err := appendConsumer(c, s)
				if err != nil {
					panic(s + " " + err.Error())
				}
			}
		}
	} else {
		err := appendConsumer(c, className)
		if err != nil {
			panic(className + " " + err.Error())
		}
	}
	keys := make([]string, 0)
	for key := range c.consumerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("package " + pkg + "\n\n")
		buf.WriteString("import (\n")
		buf.WriteString("\t\"context\"\n")
		buf.WriteString("\t\"time\"\n")
		buf.WriteString("\t\"errors\"\n")
		buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/common\"\n")
		buf.WriteString(")\n\n")
		buf.Write(c.consumerMap[key].buf)
		res = append(res, FileDes{
			name: fmt.Sprintf("%s.go",
				snakeString(c.consumerMap[key].Name)),
			buf: buf.Bytes(),
		})
	}
	return res
}

func (c *TransformImp) WriteConsumer(className string) error {
	pkg := pkgName(c.ZipEntry().ArtifactID())
	outPath := path.Join(c.opt.OutPath, "hsf", pkg, "consumer")

	cList := c.Consumer(className)
	if !util.FileExists(outPath) {
		err := os.MkdirAll(outPath, 0755)
		if err != nil {
			return err
		}
	}

	// 生成当前数据所有model
	file := c.entry.FileList()
	for _, className := range file {
		if !c.isService(className) &&
			strings.HasSuffix(className, ".class") { // todo 后期优化
			transToGoStruct(c, strings.TrimSuffix(className, ".class"))
		}
	}

	err := c.WriteModel()
	if err != nil {
		return err
	}
	for _, service := range cList {
		fileName := path.Join(outPath, service.name)
		err := ioutil.WriteFile(fileName, service.buf, 0766)
		if err != nil {
			return err
		}
		cmd := exec.Command("goimports", "-l", "-w", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("[toolkit] 格式化文件报错，请检测goimports是否安装好，" + err.Error())
		}
		cmd = exec.Command("gofmt", "-l", "-w", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("[toolkit] 格式化文件报错，请检测gofmt是否安装好，" + err.Error())
		}
	}
	return nil
}

func (c *TransformImp) isService(className string) bool {
	if !strings.HasSuffix(className, ".class") {
		return false
	}
	class, err := getClassFile(c.ZipEntry(), className)
	if err != nil {
		panic(err)
	}
	if class.File.AccessFlags()&0x0200 == 0x0200 {
		return true
	}
	return false
}

func (c *TransformImp) Provider(className string) []FileDes {
	pkg := "provider"

	if className == "all" {
		file := c.entry.FileList()
		for _, s := range file {
			if c.isService(s) {
				err := appendProvider(c, s)
				if err != nil {
					panic(s + " " + err.Error())
				}
			}
		}
	} else {
		err := appendProvider(c, className)
		if err != nil {
			panic(className + " " + err.Error())
		}
	}
	keys := make([]string, 0)
	for key := range c.providerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	res := make([]FileDes, 0)
	for _, key := range keys {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("package " + pkg + "\n\n")
		buf.WriteString("import (\n")
		buf.WriteString("\t\"context\"\n")
		buf.WriteString("\t\"time\"\n")
		buf.WriteString("\t\"errors\"\n")
		buf.WriteString("\t\"gitlab.alibaba-inc.com/amap-go/hsf-go/common\"\n")
		for dpKey := range c.providerMap[key].Dependency {
			buf.WriteString(fmt.Sprintf("\t\"%s\"\n", dpKey))
		}
		buf.WriteString(")\n\n")
		buf.Write(c.providerMap[key].buf)
		res = append(res, FileDes{
			name: fmt.Sprintf("%s.go",
				snakeString(c.providerMap[key].Name)),
			buf: buf.Bytes(),
		})
	}
	return res
}

func (c *TransformImp) WriteProvider(className string) error {
	pkg := pkgName(c.ZipEntry().ArtifactID())
	outPath := path.Join(c.opt.OutPath, "hsf", pkg, "provider")

	pList := c.Provider(className)
	if !util.FileExists(outPath) {
		err := os.MkdirAll(outPath, 0755)
		if err != nil {
			return err
		}
	}
	// 生成当前数据所有model
	file := c.entry.FileList()
	for _, className := range file {
		if !c.isService(className) &&
			strings.HasSuffix(className, ".class") { // todo 后期优化
			transToGoStruct(c, strings.TrimSuffix(className, ".class"))
		}
	}

	err := c.WriteModel()
	if err != nil {
		return err
	}

	for _, service := range pList {
		fileName := path.Join(outPath, service.name)
		err := ioutil.WriteFile(fileName, service.buf, 0766)
		if err != nil {
			return err
		}
		cmd := exec.Command("goimports", "-l", "-w", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("[toolkit] 格式化文件报错，请检测goimports是否安装好，" + err.Error())
		}
		cmd = exec.Command("gofmt", "-l", "-w", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("[toolkit] 格式化文件报错，请检测gofmt是否安装好，" + err.Error())
		}
	}

	return nil
}

func (c *TransformImp) ZipEntry() *classpath.ZipEntry {
	return c.entry
}

func NewTransformImp(ze *classpath.ZipEntry, opt Opt) *TransformImp {
	return &TransformImp{
		opt:         opt,
		entry:       ze,
		goStructMap: make(map[string]struct{}),
		structMap:   make(map[string]StructInfo),
		consumerMap: make(map[string]StructInfo),
		providerMap: make(map[string]StructInfo),
		importMap:   make(map[string]struct{}),
	}
}

func Run(opt Opt) error {
	va := validator.New()
	err := va.Struct(opt)
	if err != nil {
		return err
	}
	ze, err := classpath.NewZipEntry(opt.InGroupID, opt.InArtifactID, opt.InVersion)
	if err != nil {
		panic(err)
	}
	defer ze.Close()

	switch opt.OutType {
	case outTypeListClass:
		file := ze.FileList()
		for _, s := range file {
			fmt.Println(s)
		}
	case outTypeConsumer:
		if len(opt.OutPkg) == 0 {
			opt.OutPath = "./"
		}
		md, err := checkEnv()
		if err != nil {
			return err
		}
		if md != nil {
			opt.OutPath = "./"
			opt.OutPkg = md.Module.Path
		}
		internalPath := path.Join(opt.OutPath, "internal")
		if util.FileExists(internalPath) &&
			!config.Config().Hsf.CloseInternal {
			opt.OutPath = internalPath
			opt.HaveInternal = true
		}
		c := NewTransformImp(ze, opt)
		return c.WriteConsumer(opt.InClass)
	case outTypeProvider:
		if len(opt.OutPkg) == 0 {
			opt.OutPath = "./"
		}
		md, err := checkEnv()
		if err != nil {
			return err
		}
		if md != nil {
			opt.OutPath = "./"
			opt.OutPkg = md.Module.Path
		}
		internalPath := path.Join(opt.OutPath, "internal")
		if util.FileExists(internalPath) &&
			!config.Config().Hsf.CloseInternal {
			opt.OutPath = internalPath
			opt.HaveInternal = true
		}
		c := NewTransformImp(ze, opt)
		return c.WriteProvider(opt.InClass)
	case outTypeStruct:
		if len(opt.OutPkg) == 0 {
			return errors.New("--outPkg is require")
		}
		md, err := checkEnv()
		if err != nil {
			return err
		}
		if md != nil {
			opt.OutPath = "./"
			opt.OutPkg = md.Module.Path
		}
		internalPath := path.Join(opt.OutPath, "internal")
		if util.FileExists(internalPath) &&
			!config.Config().Hsf.CloseInternal {
			opt.OutPath = internalPath
			opt.HaveInternal = true
		}
		tr := NewTransformImp(ze, opt)
		goName, _ := transToGoStruct(tr, opt.InClass)
		if len(goName) == 0 {
			return errors.New(opt.InClass + " not found")
		}
		return tr.WriteModel()
	default:
		return errors.New("outType error")
	}
	return nil
}

func checkEnv() (*mod.Mod, error) {
	fmt.Println("依赖版本:")
	fmt.Println(hsfModel + " " + hsfVersion)
	fmt.Println(bigModel + " " + bigVersion)
	fmt.Println("replace " + hessianOld + " => " + hessianNew + " " + hessianVersion)

	dS, err := mod.ReadMod("./")
	if err == nil {
		requireMap := make(map[string]string)
		for _, v := range dS.Require {
			requireMap[v.Path] = v.Version
		}
		replaceMap := make(map[string]struct {
			Old mod.Version
			New mod.Version
		})
		for _, v := range dS.Replace {
			replaceMap[v.Old.Path] = v
		}
		var updateHsf, updateHessian, updateBig bool
		if v, ok := requireMap[hsfModel]; !ok || versionLess(v, hsfVersion) {
			updateHsf = true
		}
		if v, ok := requireMap[bigModel]; ok && versionLess(v, bigVersion) {
			updateBig = true
		}
		if v, ok := replaceMap[hessianOld]; !ok || versionLess(v.New.Version, hessianVersion) {
			updateHessian = true
		}
		if !updateHsf && !updateBig && !updateHessian {
			return dS, nil
		}

		fmt.Print("检测到go.mod，是否需要更新依赖[Y/N]:")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		text = strings.Trim(text, "\t\n ")
		if strings.ToUpper(text) == "Y" {
			if updateHsf {
				cmd := exec.Command("go", "get", "-v", hsfModel+"@"+hsfVersion)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					fmt.Println("[toolkit] 更新hsf版本出错，请检测go环境是否安装好，" + err.Error())
				}
			}

			if updateHessian {
				cmd := exec.Command("go", "mod", "edit",
					fmt.Sprintf("-replace=%s=%s@%s", hessianOld, hessianNew, hessianVersion))
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					fmt.Println("[toolkit] 更新hessian版本出错，请检测go环境是否安装好，" + err.Error())
				}
			}

			if updateBig {
				cmd := exec.Command("go", "get", "-v", bigModel+"@"+bigVersion)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					fmt.Println("[toolkit] 更新big版本出错，请检测go环境是否安装好，" + err.Error())
				}
			}
		}
		return dS, nil
	}
	return nil, nil
}

// versionLess v1 < v2
func versionLess(v1, v2 string) bool {
	f := func(s string) ([]int, error) {
		ts := ""
		for _, v := range s {
			if (v >= '0' && v <= '9') || v == '.' {
				ts += string(v)
			}
		}

		tArr := strings.Split(ts, ".")
		rArr := make([]int, 0, len(tArr))

		for _, v := range tArr {
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			rArr = append(rArr, tmp)
		}
		if len(rArr) < 3 {
			rArr = append(rArr, 0)
		}
		return rArr, nil
	}
	v1Arr, err := f(v1)
	if err != nil {
		return true
	}
	v2Arr, err := f(v2)
	if err != nil {
		return true
	}

	for k, v := range v1Arr {
		if len(v2Arr) == k {
			break
		}
		if v < v2Arr[k] {
			return true
		}
	}
	return false
}
