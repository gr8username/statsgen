package main

import (
	"os"
	"strings"
)

type runInfo struct {
	logLoc []string
	statsFile string
	help bool
}

const STATSFILE_ARG = "--statsfile"
const LOGDIR_ARG = "--logdirs"

func getArgsInstructions() runInfo {
	var ri runInfo
	state := 0
	ri.statsFile = ""
	for _, arg := range os.Args {
		if equalsIgnoreCase(arg, STATSFILE_ARG) {
			state = 1
			continue
		}
		if equalsIgnoreCase(arg, LOGDIR_ARG) {
			state = 2
			continue
		}
		if equalsIgnoreCase(arg, "--help") {
			ri.help = true
			break
		}
		if state == 1 {
			ri.statsFile = arg
			state = 0
		}
		if state == 2 {
			ri.logLoc = append(ri.logLoc, arg)
		}
	}
	return ri
}

func equalsIgnoreCase(str1 string, str2 string) bool {
	return strings.ToLower(str1) == strings.ToLower(str2)
}
