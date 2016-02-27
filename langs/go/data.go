package gogen

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/samuel/go-thrift/parser"
)

type structsFileData struct {
	TplUtils

	FilePath string

	Package  string
	Includes [][2]string
	Structs  []*structData
}

type structData struct {
	TplUtils

	Name   string
	Fields []*parser.Field
}

func getStructsFileData(pkgName, pkgDir string, includes [][2]string, structs map[string]*parser.Struct) *structsFileData {
	data := &structsFileData{
		FilePath: filepath.Join(pkgDir, "gen_"+pkgName+"_structs.go"),
		Package:  pkgName,
		Includes: includes,
	}

	for structName, parsedStruct := range structs {
		data.Structs = append(data.Structs, &structData{
			Name:   structName,
			Fields: parsedStruct.Fields,
		})
	}

	return data
}

type servicesFileData struct {
	TplUtils

	FilePath string

	Package  string
	Includes [][2]string
	Services []*serviceData
}

type serviceData struct {
	TplUtils

	Name    string
	Methods []*parser.Method
}

func getServicesFileData(pkgName, pkgDir string, includes [][2]string, services map[string]*parser.Service) *servicesFileData {
	data := &servicesFileData{
		FilePath: filepath.Join(pkgDir, "gen_"+pkgName+"_services.go"),
		Package:  pkgName,
		Includes: includes,
	}

	for serviceName, parsedService := range services {
		sData := &serviceData{
			Name: serviceName,
		}

		// sort methods
		methodNames := make([]string, 0, len(parsedService.Methods))

		for methodName, _ := range parsedService.Methods {
			methodNames = append(methodNames, methodName)
		}

		sort.Strings(methodNames)

		for _, name := range methodNames {
			sData.Methods = append(sData.Methods, parsedService.Methods[name])
		}

		data.Services = append(data.Services, sData)
	}

	return data
}

type echoFileData struct {
	TplUtils

	FilePath string

	Package string
	Service *serviceData
}

func getEchoFileData(pkgName, pkgDir string, sData *serviceData) *echoFileData {
	data := &echoFileData{
		FilePath: filepath.Join(pkgDir, fmt.Sprintf("gen_%s_%s_web_apis.go", pkgName, sData.Name)),

		Package: pkgName,
		Service: sData,
	}

	return data
}
