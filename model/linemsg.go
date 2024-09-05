package model

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/zt64/logcat-tui/logcat"
)

var re = regexp.MustCompile(`(?<Timestamp>\d+\.\d+) +(?<PID>\d+) +(?<TID>\d+) (?<Priority>\w) (?<Tag>.+?): (?<Message>.+)`)

type lineMsg struct {
	Timestamp float64
	Pid       int
	Tid       int
	Priority  logcat.Priority
	Tag       string
	Message   string
}

func parseLine(s string) lineMsg {
	match := re.FindStringSubmatch(s)

	paramsMap := make(map[string]string)
	for i, name := range match {
		paramsMap[re.SubexpNames()[i]] = name
	}

	timestamp, err := strconv.ParseFloat(paramsMap["Timestamp"], 64)
	if err != nil {
		panic(err)
	}

	pid, err := strconv.Atoi(paramsMap["PID"])
	if err != nil {
		fmt.Println("Error parsing PID: ", err)
		fmt.Println("PID: ", paramsMap)
		// print the line that caused the error
		fmt.Println("Line: ", s)
		panic(err)
	}

	tid, err := strconv.Atoi(paramsMap["TID"])
	if err != nil {
		panic(err)
	}

	return lineMsg{
		Timestamp: timestamp,
		Pid:       pid,
		Tid:       tid,
		Priority:  logcat.Priority(paramsMap["Priority"]),
		Tag:       paramsMap["Tag"],
		Message:   paramsMap["Message"],
	}
}
