package logcat

type Buffer string

const (
	BufferMain    Buffer = "main"
	BufferSystem  Buffer = "system"
	BufferRadio   Buffer = "radio"
	BufferEvents  Buffer = "events"
	BufferCrash   Buffer = "crash"
	BufferDefault Buffer = "default"
	BufferAll     Buffer = "all"
)
