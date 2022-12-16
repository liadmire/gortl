package eylog

func SetLevel(lt LevelType) {
	logInst().SetLevel(lt)
}

func SetMaxSize(size int64) {
	logInst().SetMaxSize(size)
}

func SetName(name string) {
	logInst().SetName(name)
}

func SetBackupCount(cnt int) {
	logInst().SetBackupCount(cnt)
}

func SetPath(path string) {
	logInst().SetPath(path)
}

func Debug(v ...any) {
	logInst().Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	logInst().Debugf(format, v...)
}

func Info(v ...any) {
	logInst().Info(v...)
}

func Infof(format string, v ...any) {
	logInst().Infof(format, v...)
}

func Warnf(format string, v ...any) {
	logInst().Warnf(format, v...)
}

func Error(v ...any) {
	logInst().Error(v...)
}

func Errorf(format string, v ...any) {
	logInst().Errorf(format, v...)
}

func Fatal(v ...any) {
	logInst().Fatal(v...)
}

func Fatalf(format string, v ...any) {
	logInst().Fatalf(format, v...)
}
