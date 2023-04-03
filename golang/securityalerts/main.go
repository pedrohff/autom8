package main

import (
	"bufio"
	"fmt"
	"github.com/nxadm/tail"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	currentFile     os.DirEntry
	currentFileName string
	currentFileLine int64
)

func main() {

	// 1 list files in directory after x timestamp
	dir, err := os.ReadDir(".")
	if err != nil {
		panic(err)
		return
	}
	logFiles := LogFileList{}
	for i, entry := range dir {
		regx, err := regexp.Compile("(.*)(\\.log)$")
		if err != nil {
			panic(err)
			return
		}
		if regx.MatchString(entry.Name()) {

			info, err := entry.Info()
			if err != nil {
				panic(err)
				return
			}
			logFiles = append(logFiles, File{entry: entry, ModTime: info.ModTime()})
		}

		fmt.Printf("%d - %s\n", i, entry.Name())
	}

	if logFiles.Len() <= 0 {
		panic("no log files found")
	}
	sort.Sort(logFiles)
	if currentFile == nil {
		currentFile = logFiles[0].entry
	}

	t, err := tail.TailFile(currentFile.Name(), tail.Config{Follow: true})
	for line := range t.Lines {
		if currentFileLine == 0 || int64(line.Num) > currentFileLine {
			fmt.Printf("[%d] %s\n", line.Num, line.Text)
		}
	}
	//scanner := bufio.NewScanner(cFile)
	//bufio.ScanLines()
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//	bufio.NewReader(os.Stdin).ReadString('\n')
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	log.Fatal(err)
	//}
	// 2 get file currently reading

	// 3 swap files if finished?
}

type LogFileList []File

func (l LogFileList) Len() int {
	return len(l)
}
func (l LogFileList) Less(i, j int) bool {
	return l[i].ModTime.After(l[j].ModTime)
}
func (l LogFileList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type File struct {
	entry   os.DirEntry
	ModTime time.Time
}

type AutheliaLogInput struct {
	Level    string    `json:"level"`
	Method   string    `json:"method"`
	Msg      string    `json:"msg"`
	Path     string    `json:"path"`
	RemoteIP string    `json:"remote_ip"`
	Time     time.Time `json:"time"`
}

type AppTailPointerFile []AppTailPointer

var appTailPointerFileName = ".apptailpointer"

func (a AppTailPointerFile) Open() error {
	_, err := os.Stat(appTailPointerFileName)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(appTailPointerFileName)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	open, err := os.Open(appTailPointerFileName)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		app := strings.Split(scanner.Text(), "|")
		atoi, err := strconv.Atoi(app[2])
		if err != nil {
			return err
		}
		a = append(a, AppTailPointer{
			Name:     app[0],
			LastFile: app[1],
			LastLine: int64(atoi),
		})
	}
	return nil
}
func (a AppTailPointerFile) GetByName(appName string) *AppTailPointer {
	for _, pointer := range a {
		if appName == pointer.Name {
			return &pointer
		}
	}
	return nil
}
func (a AppTailPointerFile) Write() error {
	file, err := os.OpenFile(appTailPointerFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	for _, pointer := range a {
		_, err := file.WriteString(fmt.Sprintf("%s\n", pointer))
		if err != nil {
			return err
		}
	}
	return file.Close()
}

type AppTailPointer struct {
	Name     string
	LastFile string
	LastLine int64
}

func (a AppTailPointer) String() string {
	return fmt.Sprintf("%s|%s|%d", a.Name, a.LastFile, a.LastLine)
}
