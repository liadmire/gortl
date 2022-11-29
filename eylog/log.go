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

func Debug(format string, v ...interface{}) {
	logInst().Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	logInst().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	logInst().Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	logInst().Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	logInst().Fatal(format, v...)
}
