package model

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/zt64/logcat-tui/logcat"
)

var re = regexp.MustCompile(`(?<Timestamp>\d+\.\d+) +(?<PID>\d+) +(?<TID>\d+) (?<Priority>\w) (?<Tag>.+?): (?<Message>.+)`)

// lineMsg represents a single logcat message
type lineMsg struct {
	Timestamp float64         // Timestamp is the time in seconds since the start of the logcat process
	Pid       int             // Pid is the process ID of the logcat message
	Tid       int             // Tid is the thread ID of the logcat message
	Priority  logcat.Priority // Priority is the logcat priority of the message
	Tag       string          // Tag is the logcat tag of the message
	Message   string          // Message is the logcat message itself
}

// parseLine parses a single line of logcat output into a lineMsg struct and returns it along with an error if one occurred
func parseLine(s string) (lineMsg, error) {
	match := re.FindStringSubmatch(s)

	paramsMap := make(map[string]string)
	for i, name := range match {
		paramsMap[re.SubexpNames()[i]] = name
	}

	timestamp, err := strconv.ParseFloat(paramsMap["Timestamp"], 64)
	if err != nil {
		return lineMsg{}, fmt.Errorf("error parsing timestamp: %w", err)
	}

	pid, err := strconv.Atoi(paramsMap["PID"])
	if err != nil {
		return lineMsg{}, fmt.Errorf("error parsing PID: %w\nPID: %v\nLine: %s", err, paramsMap, s)
	}

	tid, err := strconv.Atoi(paramsMap["TID"])
	if err != nil {
		return lineMsg{}, fmt.Errorf("error parsing TID: %w", err)
	}

	return lineMsg{
		Timestamp: timestamp,
		Pid:       pid,
		Tid:       tid,
		Priority:  logcat.Priority(paramsMap["Priority"]),
		Tag:       paramsMap["Tag"],
		// Message:   paramsMap["Message"],
		Message: s,
	}, nil
}
