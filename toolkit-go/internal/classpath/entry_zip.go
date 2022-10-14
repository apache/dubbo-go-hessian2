package classpath

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/mvn"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/util"
)

var (
	zipEntryMap = make(map[string]struct{})
	// ErrFileNotFound 未找到文件
	ErrFileNotFound = errors.New("file not found")
)

func zipEntryExits(d *ZipEntry) bool {
	_, ok := zipEntryMap[d.String()]
	return ok
}

// NewZipEntry 新建一个
func NewZipEntry(groupID, artifactID, version string) (*ZipEntry, error) {
	ze := &ZipEntry{
		groupID:    groupID,
		artifactID: artifactID,
		version:    version,
	}
	if zipEntryExits(ze) {
		return nil, errors.New("zip exist")
	} else {
		zipEntryMap[ze.String()] = struct{}{}
	}
	jarInfo, err := mvn.JarLocalPath(ze.groupID, ze.artifactID, ze.version)
	if err != nil {
		return nil, err
	}
	r, err := openJar(jarInfo.AbsPath)
	if err != nil {
		return nil, err
	}
	ze.zipRC = r
	for _, v := range jarInfo.DpAbsPath {
		tmp, err := NewZipEntryByJar(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ze.dependencies = append(ze.dependencies, tmp)
	}
	return ze, nil
}

func NewZipEntryByJar(jarFile string) (*ZipEntry, error) {
	r, err := openJar(jarFile)
	if err != nil {
		return nil, err
	}
	ze := &ZipEntry{
		zipRC: r,
	}
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "pom.xml") {
			xmlData, err := readFile(f)
			if err != nil {
				return nil, err
			}
			var x struct {
				XMLName xml.Name `xml:"project"`
				Parent  struct {
					GroupID    string `xml:"groupId"`
					ArtifactID string `xml:"artifactId"`
					Version    string `xml:"version"`
				} `xml:"parent"`
				GroupID    string `xml:"groupId"`
				ArtifactID string `xml:"artifactId"`
				Version    string `xml:"version"`
			}
			err = xml.Unmarshal(xmlData, &x)
			if err != nil {
				return nil, err
			}
			ze.groupID = x.GroupID
			if len(ze.groupID) == 0 {
				ze.groupID = x.Parent.GroupID
			}
			ze.artifactID = x.ArtifactID
			if len(ze.artifactID) == 0 {
				ze.artifactID = x.Parent.ArtifactID
			}
			ze.version = x.Version
			if len(ze.version) == 0 {
				ze.version = x.Parent.Version
			}
			fmt.Println(f.Name, ze.String())
			break
		}
	}

	if zipEntryExits(ze) {
		return nil, errors.New("zip exist")
	} else {
		zipEntryMap[ze.String()] = struct{}{}
	}

	return ze, nil
}

type ClassInfo struct {
	GroupID    string
	ArtifactID string
	Version    string
	File       []byte
}

// ReadClass 获取class信息
func ReadClass(ze *ZipEntry, className string) (*ClassInfo, error) {
	levelList := []*ZipEntry{ze}
	for {
		if len(levelList) == 0 {
			break
		}
		entry := levelList[0]
		if b, err := entry.ReadClass(className); err == nil {
			return &ClassInfo{
				GroupID:    entry.groupID,
				ArtifactID: entry.artifactID,
				Version:    entry.version,
				File:       b,
			}, nil
		}
		levelList = append(levelList[1:], entry.dependencies...)
	}
	return nil, ErrFileNotFound
}

type ZipEntry struct {
	groupID      string
	artifactID   string
	version      string
	zipRC        *zip.ReadCloser
	dependencies []*ZipEntry
}

func (ze *ZipEntry) String() string {
	return fmt.Sprintf("<%s,%s,%s>",
		ze.groupID, ze.artifactID, ze.version)
}

func (ze *ZipEntry) ArtifactID() string {
	return ze.artifactID
}

func (ze *ZipEntry) Close() {
	ze.zipRC.Close()
	for _, dependency := range ze.dependencies {
		dependency.Close()
	}
}

func (ze *ZipEntry) FileList() []string {
	files := make([]string, 0, len(ze.zipRC.File))
	for _, file := range ze.zipRC.File {
		files = append(files, file.Name)
	}
	return files
}

func (ze *ZipEntry) ReadClass(className string) ([]byte, error) {
	findClassName := className
	if !strings.HasSuffix(className, ".class") {
		findClassName = strings.ReplaceAll(className, ".", "/") + ".class"
	}

	classFile, err := ze.findFile(findClassName)
	if err != nil {
		return nil, err
	}

	data, err := readFile(classFile)
	return data, err
}

func (ze *ZipEntry) findFile(className string) (*zip.File, error) {
	for _, f := range ze.zipRC.File {
		if f.Name == className {
			return f, nil
		}
	}
	return nil, ErrFileNotFound
}

func readFile(classFile *zip.File) ([]byte, error) {
	rc, err := classFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func openJar(absPath string) (*zip.ReadCloser, error) {
	if !util.FileExists(absPath) {
		return nil, errors.New("jar path not exist, " + absPath)
	}

	if !strings.HasSuffix(absPath, ".jar") &&
		!strings.HasSuffix(absPath, ".JAR") &&
		!strings.HasSuffix(absPath, ".zip") &&
		!strings.HasSuffix(absPath, ".ZIP") {
		return nil, errors.New("jar path error")
	}

	return zip.OpenReader(absPath)
}
