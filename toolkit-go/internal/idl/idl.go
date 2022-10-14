package idl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/mvn"

	"github.com/go-playground/validator/v10"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/hsf"
	bindata "github.com/apache/dubbo-go-hessian2/toolkit-go/internal/idl/binddata"
)

const (
	outTypeInit     = "init"
	outTypeDeploy   = "deploy"
	outTypeConsumer = "consumer"
	outTypeProvider = "provider"
)

type Opt struct {
	// in
	GroupID    string
	ArtifactID string
	IDLPath    string
	// out
	OutType  string `validate:"oneof=init deploy consumer provider"`
	OutPath  string
	OutQuote bool
}

type printWriter struct {
	buf bytes.Buffer
}

func (pw *printWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return pw.buf.Write(p)
}

// Run 1.16后才支持embed，暂时先用esc包支持
func Run(opt Opt) error {
	va := validator.New()
	err := va.Struct(opt)
	if err != nil {
		return err
	}
	switch opt.OutType {
	case outTypeInit:
		err := createDir(opt)
		if err != nil {
			return err
		}
		return createFile(opt)
	case outTypeDeploy:
		d, err := mvn.Deploy(opt.IDLPath)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] groupId: %s\n", d.GroupID)
		fmt.Printf("[INFO] artifactId: %s\n", d.ArtifactID)
		fmt.Printf("[INFO] version: %s\n", d.Version)
	case outTypeConsumer:
		d, err := mvn.Deploy(opt.IDLPath)
		if err != nil {
			return err
		}
		outPath := path.Join(opt.IDLPath, "../../")
		return hsf.Run(hsf.Opt{
			InClass:      "all",
			InGroupID:    d.GroupID,
			InArtifactID: d.ArtifactID,
			InVersion:    d.Version,
			OutType:      "consumer",
			OutPkg:       "modname",
			OutPath:      outPath,
			OutQuote:     opt.OutQuote,
		})
	case outTypeProvider:
		d, err := mvn.Deploy(opt.IDLPath)
		if err != nil {
			return err
		}
		outPath := path.Join(opt.IDLPath, "../../")
		return hsf.Run(hsf.Opt{
			InClass:      "all",
			InGroupID:    d.GroupID,
			InArtifactID: d.ArtifactID,
			InVersion:    d.Version,
			OutType:      "provider",
			OutPkg:       "modname",
			OutPath:      outPath,
			OutQuote:     opt.OutQuote,
		})
	default:
		return errors.New("outType is error")
	}
	return nil
}

func createDir(opt Opt) error {
	if len(opt.OutPath) == 0 {
		opt.OutPath = "./"
	}
	if len(opt.GroupID) == 0 {
		opt.GroupID = "com.amap.aos"
	}
	mainSrc := mainSrc(opt)
	err := os.MkdirAll(mainSrc, 0755)
	if err != nil {
		return err
	}
	// common
	commonSrc := path.Join(mainSrc, "common")
	err = os.MkdirAll(commonSrc, 0755)
	if err != nil {
		return err
	}
	// request
	requestSrc := path.Join(mainSrc, "request", "demo")
	err = os.MkdirAll(requestSrc, 0755)
	if err != nil {
		return err
	}
	// response
	responseSrc := path.Join(mainSrc, "response", "demo")
	err = os.MkdirAll(responseSrc, 0755)
	if err != nil {
		return err
	}
	// service
	serviceSrc := path.Join(mainSrc, "service")
	err = os.MkdirAll(serviceSrc, 0755)
	return err
}

func pkgName(opt Opt) string {
	tmpD := strings.Split(opt.ArtifactID, "-")
	d := make([]string, 0, len(tmpD))
	for _, s := range tmpD {
		d = append(d, strings.Split(s, ".")...)
	}
	return strings.Join(d, ".")
}

func baseDir(opt Opt) string {
	tmpD := strings.Split(opt.ArtifactID, "-")
	d := make([]string, 0, len(tmpD))
	for _, s := range tmpD {
		d = append(d, strings.Split(s, ".")...)
	}
	return strings.Join(d, "-")
}

func mainSrc(opt Opt) string {
	return path.Join(opt.OutPath, "idl", baseDir(opt),
		"src", "main", "java", pkgName(opt))
}

func createFile(opt Opt) error {
	pom := bindata.FSMustString(false, "/templates/pom.tmpl")
	pom = strings.ReplaceAll(pom, "{{.GroupID}}", opt.GroupID)
	pom = strings.ReplaceAll(pom, "{{.ArtifactID}}", opt.ArtifactID)
	mainSrc := mainSrc(opt)
	err := ioutil.WriteFile(path.Join(opt.OutPath, "idl", baseDir(opt), "pom.xml"),
		[]byte(pom), 0766)
	if err != nil {
		return err
	}
	iml := bindata.FSMustString(false, "/templates/iml.tmpl")
	err = ioutil.WriteFile(path.Join(opt.OutPath, "idl", baseDir(opt),
		fmt.Sprintf("%s.iml", opt.ArtifactID)),
		[]byte(iml), 0766)
	if err != nil {
		return err
	}
	gitignore := bindata.FSMustString(false, "/templates/gitignore.tmpl")
	err = ioutil.WriteFile(path.Join(opt.OutPath, "idl", baseDir(opt), ".gitignore"),
		[]byte(gitignore), 0766)
	if err != nil {
		return err
	}
	common := bindata.FSMustString(false, "/templates/common.tmpl")
	common = strings.ReplaceAll(common, "{{.pkgName}}", pkgName(opt))
	err = ioutil.WriteFile(path.Join(mainSrc, "common", "ResultWrapper.java"),
		[]byte(common), 0766)
	if err != nil {
		return err
	}
	demoRequest := bindata.FSMustString(false, "/templates/demo_request.tmpl")
	demoRequest = strings.ReplaceAll(demoRequest, "{{.pkgName}}", pkgName(opt))
	err = ioutil.WriteFile(path.Join(mainSrc, "request", "demo", "UserName.java"),
		[]byte(demoRequest), 0766)
	if err != nil {
		return err
	}
	demoResponse := bindata.FSMustString(false, "/templates/demo_response.tmpl")
	demoResponse = strings.ReplaceAll(demoResponse, "{{.pkgName}}", pkgName(opt))
	err = ioutil.WriteFile(path.Join(mainSrc, "response", "demo", "Hello.java"),
		[]byte(demoResponse), 0766)
	if err != nil {
		return err
	}
	demoService := bindata.FSMustString(false, "/templates/demo_service.tmpl")
	demoService = strings.ReplaceAll(demoService, "{{.pkgName}}", pkgName(opt))
	err = ioutil.WriteFile(path.Join(mainSrc, "service", "DemoService.java"),
		[]byte(demoService), 0766)
	if err != nil {
		return err
	}
	return nil
}
