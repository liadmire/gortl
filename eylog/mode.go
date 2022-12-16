package eylog

const (
	debugMode = iota
	releaseMode
)

var mode int = debugMode

func ReleaseMode() {
	mode = releaseMode
}
