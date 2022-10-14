package mvn

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	bindata "github.com/apache/dubbo-go-hessian2/toolkit-go/internal/mvn/binddata"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/util"
)

const (
	dockerName = "reg.docker.alibaba-inc.com/tgh/maven-deploy"
)

func m2Dir() (*string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	resDir := path.Join(home, ".toolkit", ".m2", "repository")
	err = util.CreateDir(resDir)
	if err != nil {
		return nil, err
	}
	return &resDir, nil
}

type printWriter struct {
	buf bytes.Buffer
}

func (pw *printWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return pw.buf.Write(p)
}

func cmdMvn(resource string, cmds ...string) ([]byte, error) {
	m2Path, err := m2Dir()
	if err != nil {
		return nil, err
	}
	arg := []string{"run", "-t",
		"-v", fmt.Sprintf("%s:/usr/src/.m2/repository", *m2Path)}
	if len(resource) != 0 {
		arg = append(arg, "-v", fmt.Sprintf("%s:/usr/src/resource", resource))
	}
	arg = append(arg, dockerName)

	arg = append(arg, "/usr/src/apache-maven/bin/mvn")
	arg = append(arg, "-s")
	arg = append(arg, "/usr/src/apache-maven/conf/settings.xml")
	arg = append(arg, cmds...)
	cmd := exec.Command("docker", arg...)
	pw := &printWriter{}
	cmd.Stdout = pw
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	d := pw.buf.Bytes()
	if !strings.Contains(string(d),
		"BUILD SUCCESS") {
		return nil, errors.New("build fail")
	}
	return d, nil
}

func DownLoadPackage(groupID, artifactID, version string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pomPath := path.Join(home, ".toolkit", ".pom")
	err = util.CreateDir(pomPath)
	if err != nil {
		return nil, err
	}
	pomFile := path.Join(pomPath, artifactID+"-pom.xml")
	pomTmpl := bindata.FSMustString(false,
		"/templates/pom.tmpl")

	tmpl, err := template.New("pom.tmpl").Parse(pomTmpl)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(pomFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(f, struct {
		GroupID    string
		ArtifactID string
		Version    string
	}{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Version:    version,
	})
	if err != nil {
		return nil, err
	}

	downByte, err := cmdMvn(pomPath,
		"-U", "-f", "/usr/src/resource/"+artifactID+"-pom.xml", "install")
	if err != nil {
		return nil, err
	}
	result, err := parseDependencies(downByte)
	if err != nil {
		return nil, err
	}
	vMap := make(map[string]struct{})
	for _, v1 := range result {
		vMap[v1] = struct{}{}
	}
	result = make([]string, 0)
	for v1 := range vMap {
		result = append(result, v1)
	}
	m2Path, err := m2Dir()
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(path.Join(*m2Path, "toolkit"))
	if err != nil {
		return nil, err
	}

	return result, os.RemoveAll(pomPath)
}

func parseDependencies(downByte []byte) ([]string, error) {
	m2Path, err := m2Dir()
	if err != nil {
		return nil, err
	}

	dList, err := os.ReadDir(path.Join(*m2Path))
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile("(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
	result := reg.FindAll(downByte, -1)

	dp := make([]string, 0)
	for _, v := range result {
		for _, entry := range dList {
			i := strings.Index(string(v), entry.Name()+"/")
			if i == -1 {
				continue
			}
			//maven-metadata.xml
			fileDir := path.Join(*m2Path, string(v[i:len(v)-19]))
			sDir := strings.Split(fileDir, "/")
			jarFile := path.Join(fileDir, sDir[len(sDir)-2]+"-"+sDir[len(sDir)-1]+".jar")
			if !util.FileExists(jarFile) {
				continue
			}
			dp = append(dp, jarFile)
		}
	}
	return dp, nil
}

type JarInfo struct {
	AbsPath                      string
	GroupID, ArtifactID, Version string
	DpAbsPath                    []string
}

// JarLocalPath 获取jar包信息
func JarLocalPath(groupID, artifactID, version string) (*JarInfo, error) {
	m2Path, err := m2Dir()
	if err != nil {
		return nil, err
	}
	gPath := path.Join(strings.Split(groupID, ".")...)

	if strings.HasPrefix(version, "${") ||
		len(version) == 0 {
		dList, err := os.ReadDir(path.Join(*m2Path, gPath, artifactID))
		if err != nil {
			return nil, err
		}
		vList := make([]string, 0, len(dList))
		for _, v := range dList {
			vList = append(vList, v.Name())
		}
		sort.Strings(vList)
		if len(vList) != 0 {
			version = vList[len(vList)-1]
		}
	}
	jarPath := path.Join(*m2Path, gPath, artifactID, version,
		artifactID+"-"+version+".jar")
	fmt.Printf("下载依赖：%s %s %s\n", groupID, artifactID, version)
	dp, err := DownLoadPackage(groupID, artifactID, version)
	if err != nil {
		return nil, err
	}
	if !util.FileExists(jarPath) {
		return nil, errors.New("not found, " + jarPath)
	}
	return &JarInfo{
		AbsPath:    jarPath,
		GroupID:    groupID,
		ArtifactID: artifactID,
		Version:    version,
		DpAbsPath:  dp,
	}, nil
}

type DeployData struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// Deploy 生成本地数据
func Deploy(idlPath string) (*DeployData, error) {
	pomPath := path.Join(idlPath, "pom.xml")
	if !util.FileExists(pomPath) {
		return nil, errors.New("pom.xml not found")
	}

	_, err := cmdMvn(idlPath, "deploy")
	if err != nil {
		return nil, err
	}

	xmlData, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return nil, err
	}
	var x struct {
		XMLName    xml.Name `xml:"project"`
		GroupID    string   `xml:"groupId"`
		ArtifactID string   `xml:"artifactId"`
		Version    string   `xml:"version"`
	}
	err = xml.Unmarshal(xmlData, &x)
	if err != nil {
		return nil, err
	}
	targetPath := path.Join(idlPath, "target")
	if util.IsDir(targetPath) {
		err := os.RemoveAll(targetPath)
		if err != nil {
			return nil, err
		}
	}
	return &DeployData{
		GroupID:    x.GroupID,
		ArtifactID: x.ArtifactID,
		Version:    x.Version,
	}, nil
}
