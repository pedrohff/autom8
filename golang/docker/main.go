package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Input struct {
	SubmoduleModsList string
	SubmoduleList     string
	AppName           string
}

func main() {
	if len(os.Args) < 2 {
		panic("appname is invalid")
	}
	appName := os.Args[1]
	file, err := os.ReadFile("docker/Dockerfile.dockertemplate")
	if err != nil {
		panic(err)
		return
	}

	parsedTemplate, err := template.New("test").Parse(string(file))
	if err != nil {
		panic(err)
		return
	}

	modListBuffer := bytes.NewBufferString("")
	submodulesListBuffer := bytes.NewBufferString("")
	_ = modListBuffer
	goworkFile, err := os.ReadFile("go.work")
	if err != nil {
		panic(err)
		return
	}
	goworkFileLines := strings.Split(string(goworkFile), "\n")
	scan := false
	goworkSubModules := make([]string, 0, 0)
	for _, k := range goworkFileLines {
		if scan && ")" != k {
			goworkSubModules = append(goworkSubModules, strings.TrimSpace(k))
		}
		if "use (" == k {
			scan = true
		}
	}

	for _, subModule := range goworkSubModules {
		modListBuffer.WriteString(fmt.Sprintf("COPY %s/go.mod /app/%s/go.mod\n", subModule, subModule))
		modListBuffer.WriteString(fmt.Sprintf("COPY %s/go.sum /app/%s/go.sum\n", subModule, subModule))
		submodulesListBuffer.WriteString(fmt.Sprintf("COPY %s /app/%s\n", subModule, subModule))
	}

	bufferOutput := bytes.NewBufferString("")
	err = parsedTemplate.Execute(bufferOutput, Input{
		SubmoduleModsList: modListBuffer.String(),
		SubmoduleList:     submodulesListBuffer.String(),
		AppName:           appName,
	})
	if err != nil {
		panic(err)
		return
	}
	err = os.WriteFile("docker/homebridgepooler.Dockerfile", bufferOutput.Bytes(), 0644)
	if err != nil {
		panic(err)
		return
	}
}
