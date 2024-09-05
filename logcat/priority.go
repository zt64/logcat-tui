package logcat

type Priority string

const (
	PriorityVerbose Priority = "V"
	PriorityDebug   Priority = "D"
	PriorityInfo    Priority = "I"
	PriorityWarn    Priority = "W"
	PriorityError   Priority = "E"
	PriorityFatal   Priority = "F"
	PrioritySilent  Priority = "S"
)
