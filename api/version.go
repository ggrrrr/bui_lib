package api

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/ggrrrr/bui_lib/build"
)

func readBuildInfo() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		apiName = "unset"
		log.Printf("no build info")
		return
	}
	buildInfo = bi
	path := strings.Split(bi.Main.Path, "/")
	apiName = fmt.Sprintf("%s", path[len(path)-1])
	m := ModulesInfo{Path: bi.Main.Path, Version: bi.Main.Version, Sum: bi.Main.Sum}
	modules = append(modules, m)
	log.Printf("%s:main module: %s %+v\n", apiName, bi.Path, bi.Main)
	for _, dep := range bi.Deps {
		m := ModulesInfo{Path: dep.Path, Version: dep.Version, Sum: dep.Sum}
		if dep.Replace != nil {
			m.Replace = fmt.Sprintf("%s", dep.Replace.Path)
		}
		// log.Printf("\tmodule: %d: %+v\n", k, dep)
		modules = append(modules, m)
	}
}

func Version() string {
	return fmt.Sprintf("%s/%s", buildInfo.Main.Version, buildInfo.Main.Sum)
}

func BuildInfo() string {
	return fmt.Sprintf("%s/%s", build.BuildOs, build.BuildTime)
}
