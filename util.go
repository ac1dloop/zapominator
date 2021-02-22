package main

import (
	"strings"
)

func extractFirst(s *string) string {

	arr := strings.SplitN(*s, " ", 2)

	if len(arr) < 2 {
		return arr[0]
	}

	*s = arr[1]

	return arr[0]
}

func parseInterval(str string) TimeInterval {
	switch str {
	case "once":
		return once
	case "daily":
		return day
	case "monthly":
		return month
	case "weekly":
		return week
	case "annually":
		return year
	}

	return custom
}
