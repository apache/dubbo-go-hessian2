package util

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type GoEnv struct {
	WorkDir   string
	IsMod     bool
	ModDetail *ModInfo
	GoPath    string
	GoRoot    string
}

type ModInfo struct {
	ModRoot     string
	PackageList []*ModPackage
	CachePath   string
}

type ModPackage struct {
	Path           string
	Version        string
	ReplacePath    string
	ReplaceVersion string
}

func New() *GoEnv {
	goEnv := &GoEnv{}
	goEnv.New()
	return goEnv
}

func (e *GoEnv) New() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("get pwd error")
	}
	e.GoRoot = os.Getenv("GOROOT")
	e.GoPath = os.Getenv("GOPATH")
	e.WorkDir = pwd
	openModule := os.Getenv("GO111MODULE")
	if openModule == "on" {
		e.IsMod = true
		e.ModDetail = &ModInfo{
			CachePath: os.Getenv("GOPATH") + "/pkg/mod",
		}
		e.FormatModFile(pwd + "/go.mod")
	}

}

type RequireList struct {
	Path    string
	Version string
}

type ReplaceList struct {
	Path           string
	Version        string
	ReplacePath    string
	ReplaceVersion string
}

func (e *GoEnv) FormatModFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	mode := 1
	requireList := make([]RequireList, 0)
	replaceMap := make(map[string]ReplaceList, 0)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		line = strings.Replace(line, "\t", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Trim(line, " ")
		mods := strings.Split(line, " ")
		if len(mods) == 0 {
			continue
		}
		switch mode {
		case 1:
			if strings.Contains(line, "require") && strings.Contains(line, "(") {
				mode = 2
			}
			if strings.Contains(line, "replace") && strings.Contains(line, "(") {
				mode = 3
			}

			if len(mods) >= 2 && mods[0] == "module" {
				e.ModDetail.ModRoot = mods[1]
			}
			if len(mods) >= 3 && mods[0] == "require" {
				requireList = append(requireList, RequireList{
					Path:    mods[1],
					Version: mods[2],
				})
			}
			if len(mods) >= 5 && mods[0] == "replace" {
				replaceMap[mods[1]+mods[2]] = ReplaceList{
					Path:           mods[1],
					Version:        mods[2],
					ReplacePath:    mods[3],
					ReplaceVersion: mods[4],
				}
			}

		case 2:
			if len(mods) >= 2 {
				requireList = append(requireList, RequireList{
					Path:    mods[0],
					Version: mods[1],
				})
			}
			if strings.Contains(line, ")") {
				mode = 1
			}
		case 3:
			if len(mods) >= 4 {
				replaceMap[mods[0]+mods[1]] = ReplaceList{
					Path:           mods[0],
					Version:        mods[1],
					ReplacePath:    mods[2],
					ReplaceVersion: mods[3],
				}
			}
			if strings.Contains(line, ")") {
				mode = 1
			}
		}
	}
	e.ModDetail.PackageList = make([]*ModPackage, 0)
	for _, require := range requireList {
		key := require.Path + require.Version
		if replace, ok := replaceMap[key]; ok {
			e.ModDetail.PackageList = append(e.ModDetail.PackageList, &ModPackage{
				Path:           replace.Path,
				Version:        replace.Version,
				ReplacePath:    replace.ReplacePath,
				ReplaceVersion: replace.ReplacePath,
			})
		} else {
			e.ModDetail.PackageList = append(e.ModDetail.PackageList, &ModPackage{
				Path:    require.Path,
				Version: require.Version,
			})
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (e *GoEnv) FindPackageDir(packagePath string) (string, error) {
	// 1. find in vendor

	vendorPath := e.WorkDir + "/vendor/" + packagePath
	if ok, _ := PathExists(vendorPath); ok {
		return vendorPath, nil
	}

	// 2. find in gopath/pkg/mod
	if e.IsMod {
		//项目本身的package
		if strings.HasPrefix(packagePath, e.ModDetail.ModRoot) {
			trimPrefixDir := strings.TrimPrefix(packagePath, e.ModDetail.ModRoot)
			dir := e.WorkDir + trimPrefixDir
			if ok, _ := PathExists(dir); ok {
				return dir, nil
			}
		} else {
			for _, p := range e.ModDetail.PackageList {
				if strings.HasPrefix(packagePath, p.Path) {
					trimPrefixDir := strings.TrimPrefix(packagePath, p.Path)
					if len(p.ReplacePath) == 0 {
						dir := e.ModDetail.CachePath + "/" + p.Path + "@" + p.Version + "/" + trimPrefixDir
						if ok, _ := PathExists(dir); ok {
							return dir, nil
						}
					} else {
						dir := e.ModDetail.CachePath + "/" + p.ReplacePath + "@" + p.ReplaceVersion + "/" + trimPrefixDir
						if ok, _ := PathExists(dir); ok {
							return dir, nil
						}
					}
				}
			}
		}
	}
	// 3.find in goPath/src
	if len(e.GoPath) > 0 {
		dir := e.GoPath + "/src/" + packagePath
		if ok, _ := PathExists(dir); ok {
			return dir, nil
		}
	}
	// 4. find in goRoot/src
	if len(e.GoRoot) > 0 {
		dir := e.GoRoot + "/src/" + packagePath
		if ok, _ := PathExists(dir); ok {
			return dir, nil
		}
	}

	return "", errors.New("pkg not find")
}
