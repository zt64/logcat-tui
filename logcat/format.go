package logcat

type Format string

const (
	FormatBrief      Format = "brief"
	FormatLong       Format = "long"
	FormatProcess    Format = "process"
	FormatRaw        Format = "raw"
	FormatTag        Format = "tag"
	FormatThread     Format = "thread"
	FormatThreadTime Format = "threadtime"
	FormatTime       Format = "time"
)
