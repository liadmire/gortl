package eylog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/liadmire/gortl/eysys"
)

type LevelType int

const (
	DEBUG LevelType = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelPrefix = [FATAL + 1]string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[FATAL]"}

type logger struct {
	sync.RWMutex
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	fileWriter    *os.File  `json:"-"`
	consoleWriter io.Writer `json:"-"`
	Level         LevelType `json:"leveltype"`
	MaxSize       int64     `json:"maxsize"`
	currentSize   int64     `json:"-"`
	Rotate        bool      `json:"rotate"`
	Perm          string    `json:"perm"`
	RotatePerm    string    `json:"rotateperm"`
	BackupCount   int       `json:"backupcount"`
}

var (
	gLogInst *logger
	once     sync.Once
)

func logInst() *logger {
	once.Do(func() {
		gLogInst = &logger{
			Name:          fmt.Sprintf("%s.txt", eysys.SelfNameWithoutExt()),
			Path:          path.Join(eysys.SelfDir(), "log"),
			consoleWriter: os.Stdout,
			Level:         DEBUG,
			MaxSize:       16 * 1024 * 1024,
			Rotate:        true,
			Perm:          "0660",
			RotatePerm:    "0440",
			BackupCount:   32,
		}
	})
	return gLogInst
}

func (l *logger) logFile() string {
	return path.Join(l.Path, l.Name)
}

func (l *logger) doClean() {
	files, err := ioutil.ReadDir(l.Path)
	if err != nil {
		return
	}

	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == l.Name {
			continue
		}

		fileList = append(fileList, file.Name())
	}

	if len(fileList) <= l.BackupCount {
		return
	}

	sort.Strings(fileList)

	for i := 0; i < len(fileList)-l.BackupCount; i++ {
		os.Remove(path.Join(l.Path, fileList[i]))
	}
}

func (l *logger) doRotate() error {
	defer l.doClean()

	var fName string
	var t time.Time
	var err error
	_, err = os.Lstat(l.logFile())
	if err != nil {
		goto RESTART_logger
	}

	if !(l.MaxSize > 0 && l.currentSize >= l.MaxSize) {
		return fmt.Errorf("%s", l.Name)
	}

	t = time.Now()
	fName = fmt.Sprintf("%s.txt", t.Format("20060102150405"))
	fName = path.Join(l.Path, fName)
	_, err = os.Lstat(fName)
	if err == nil {
		return err
	}

	err = l.fileWriter.Close()
	if err != nil {
		return err
	}

	err = os.Rename(l.logFile(), fName)
	if err != nil {
		goto RESTART_logger
	}

RESTART_logger:
	l.startlogger()

	return nil
}

func (l *logger) writeMsg(level LevelType, msg string, v ...interface{}) error {
	l.Lock()
	defer l.Unlock()

	if l.fileWriter == nil {
		if err := l.startlogger(); err != nil {
			return err
		}
	}

	if l.Rotate {
		l.doRotate()
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}

	when := time.Now().Format("2006-01-02 15:04:05")
	prefix := levelPrefix[level]

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	_, fileName := path.Split(file)
	msg = fmt.Sprintf("%s%s %s:%d %s\n", when, prefix, fileName, line, msg)

	var err error
	if level == DEBUG {
		_, err = l.consoleWriter.Write([]byte(msg))
	} else {
		_, err = l.fileWriter.Write([]byte(msg))
		if err == nil {
			l.currentSize += int64(len(msg))
		}
	}
	return err
}

func (l *logger) startlogger() error {
	perm, err := strconv.ParseInt(l.Perm, 8, 64)
	if err != nil {
		return err
	}

	if !eysys.FileExists(l.Path) {
		os.MkdirAll(l.Path, os.ModePerm)
	}

	fd, err := os.OpenFile(l.logFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return err
	}

	if err := os.Chmod(l.logFile(), os.FileMode(perm)); err != nil {
		return err
	}

	l.fileWriter = fd

	fInfo, err := fd.Stat()
	if err != nil {
		return err
	}

	l.currentSize = fInfo.Size()
	return nil
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.Level > DEBUG {
		return
	}

	l.writeMsg(DEBUG, format, v...)
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.Level > INFO {
		return
	}
	l.writeMsg(INFO, format, v...)
}

func (l *logger) Warn(format string, v ...interface{}) {
	if l.Level > WARN {
		return
	}
	l.writeMsg(WARN, format, v...)
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.Level > ERROR {
		return
	}
	l.writeMsg(ERROR, format, v...)
}

func (l *logger) Fatal(format string, v ...interface{}) {
	if l.Level > FATAL {
		return
	}
	l.writeMsg(FATAL, format, v...)
}

func (l *logger) SetLevel(lt LevelType) {
	l.Level = lt
}

func (l *logger) GetLevel() LevelType {
	return l.Level
}

func (l *logger) SetMaxSize(size int64) {
	l.MaxSize = size
}

func (l *logger) SetName(name string) {
	l.Name = name
}

func (l *logger) SetPath(path string) {
	l.Path = path
}

func (l *logger) SetBackupCount(cnt int) {
	l.BackupCount = cnt
}
