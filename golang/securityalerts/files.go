package main

import (
	"fmt"
	"github.com/nxadm/tail"
	"go.uber.org/zap"
	"os"
	"regexp"
	"sort"
	"time"
)

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

type LineProcessor interface {
	Name() string
	LogSuffix() string
	Process(lineNumber int, content string) error
}

type FileReader interface {
	TailDir(directory string) error
}

type fileReader struct {
	app     LineProcessor
	pointer *Pointer
	logger  *zap.Logger
}

func NewFileReader(app LineProcessor, pointer *Pointer, logger *zap.Logger) FileReader {
	return &fileReader{
		app:     app,
		pointer: pointer,
		logger:  logger,
	}
}

func (f *fileReader) TailDir(directory string) error {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	logFiles := LogFileList{}
	for _, entry := range dir {
		regx, err := regexp.Compile(fmt.Sprintf("(.*)(\\.%s)$", f.app.LogSuffix()))
		if err != nil {
			return err
		}
		if regx.MatchString(entry.Name()) {

			info, err := entry.Info()
			if err != nil {
				return err
			}
			logFiles = append(logFiles, File{entry: entry, ModTime: info.ModTime()})
		}

	}

	if logFiles.Len() <= 0 {
		return fmt.Errorf("no log files found")
	}
	sort.Sort(logFiles)
	appName := f.app.Name()
	pointer := f.pointer.GetByName(appName)

	if pointer == nil {
		f.pointer.Create(appName, logFiles[0].entry.Name())
		pointer = f.pointer.GetByName(appName)
	}

	// validando se deve atualizar o arquivo atual para um mais novo
	for i, file := range logFiles {
		isCurrentFile := file.entry.Name() == pointer.FileName
		lastIndex := len(logFiles) - 1
		isNotLastFile := i != lastIndex
		if isCurrentFile && isNotLastFile {
			name := logFiles[lastIndex].entry.Name()
			f.pointer.UpdateByName(appName, name, 0)
			f.logger.Info("rotating current file", zap.String("currentFile", name))
		}
	}

	filePointer := f.pointer.GetByName(appName)
	t, err := tail.TailFile(fmt.Sprintf("%s/%s", directory, filePointer.FileName), tail.Config{Follow: true})

	for {
		lastLine := f.pointer.GetByName(appName).LastLine
		select {
		case <-time.After(10 * time.Second):
			f.logger.Debug("re-reading directory after timeout")
			f.pointer.Write()
			return f.TailDir(directory)

		case line := <-t.Lines:
			if lastLine == 0 || int64(line.Num) > lastLine {
				err := f.app.Process(line.Num, line.Text)
				if err != nil {
					f.logger.Error("failed to process line", zap.Error(err))
					return err
				} else {
					f.logger.Info("line read", zap.String("fileName", filePointer.FileName), zap.Int("line", line.Num))
				}
				f.pointer.UpdateByName(appName, filePointer.FileName, int64(line.Num))
			}
		}
	}
	return nil
}
func (f *fileReader) Close() error {
	return f.pointer.Close()
}
