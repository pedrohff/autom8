package main

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Pointer struct {
	filesList []AppTailPointer
	sync.Mutex
	logger *zap.Logger
}

var appTailPointerFileName = ".apptailpointer"

func NewPointer(logger *zap.Logger) (*Pointer, error) {
	_, err := os.Stat(appTailPointerFileName)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(appTailPointerFileName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	open, err := os.Open(appTailPointerFileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(open)
	filesList := []AppTailPointer{}
	for scanner.Scan() {
		app := strings.Split(scanner.Text(), "|")
		atoi, err := strconv.Atoi(app[2])
		if err != nil {
			return nil, err
		}
		filesList = append(filesList, AppTailPointer{
			Name:     app[0],
			FileName: app[1],
			LastLine: int64(atoi),
		})
	}
	pointer := &Pointer{filesList: filesList, logger: logger}
	return pointer, nil
}

func (a *Pointer) StartWorker() context.CancelFunc {

	ctx, cancelFunc := context.WithCancel(context.Background())
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				a.logger.Warn("stopping pointer worker due to canceled context")
			case <-ticker.C:
				err := a.Write()
				if err != nil {
					a.logger.Error("failed to write pointer in worker", zap.Error(err))
				}
			}
		}
	}()

	return cancelFunc
}

func (a *Pointer) GetByName(appName string) *AppTailPointer {
	for _, pointer := range a.filesList {
		if appName == pointer.Name {
			return &pointer
		}
	}
	return nil
}

func (a *Pointer) Create(appName, fileName string) {
	a.filesList = append(a.filesList, AppTailPointer{
		Name:     appName,
		FileName: fileName,
		LastLine: 0,
	})
}

func (a *Pointer) UpdateByName(appName string, fileName string, lastLine int64) {
	for i, pointer := range a.filesList {
		if appName == pointer.Name {
			a.filesList[i].FileName = fileName
			a.filesList[i].LastLine = lastLine
		}
	}
}

func (a *Pointer) Write() error {
	a.Lock()
	defer a.Unlock()
	file, err := os.OpenFile(appTailPointerFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	for _, pointer := range a.filesList {
		_, err := file.WriteString(fmt.Sprintf("%s\n", pointer))
		if err != nil {
			return err
		}
	}
	return file.Close()
}

func (a *Pointer) Close() error {
	return a.Write()
}

type AppTailPointer struct {
	Name     string
	FileName string
	LastLine int64
}

func (a AppTailPointer) String() string {
	return fmt.Sprintf("%s|%s|%d", a.Name, a.FileName, a.LastLine)
}
