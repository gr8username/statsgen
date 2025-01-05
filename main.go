package main

import (
	"fmt"
	"runtime"
	"strings"
	"path/filepath"
	"os"
)

var LINUX_DEFAULT_PATH = "~/.minecraft/logs/"
var WINDOWS_DEFAULT_PATH = "%AppData%\\Roaming\\.minecraft\\logs\\"
var DARWIN_DEFAULT_PATH = "~/Library/Application Support/minecraft/logs/"

var PRESET_STATS = ""

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(RED + "Error getting home directory!" + RESET)
		fmt.Println(err)
	}
	ri := getArgsInstructions()
	if ri.help {
		fmt.Println("Example command")
		if runtime.GOOS == "windows" {
			fmt.Println(".\\statsgen " + LOGDIR_ARG + " C:\\Users\\user\\AppData\\Roaming\\.minecraft\\logs\\ " + STATSFILE_ARG + " stats.txt")
		} else {
			fmt.Println("./statsgen " + LOGDIR_ARG + " /home/user/.minecraft/logs/ " + STATSFILE_ARG + " stats.txt")
		}
		fmt.Println("If you run the program with no arguments, it will manually prompt you to enter the information.")
		return
	}
	PRESET_STATS = ri.statsFile
	if len(ri.logLoc) != 0 {
		runScanner(ri.logLoc)
		return
	} 
	LINUX_DEFAULT_PATH = filepath.Join(homeDir, ".minecraft", "logs")
	if runtime.GOOS == "linux" || runtime.GOOS == "bsd" {
		fmt.Print("Enter your minecraft logs directory (default: " + LINUX_DEFAULT_PATH + "): ")
		var logLoc string
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "/") {
			logLoc = LINUX_DEFAULT_PATH
		}
		runScanner([]string{logLoc})
	} else if runtime.GOOS == "darwin" {
		DARWIN_DEFAULT_PATH = filepath.Join(homeDir, "Library", "Application Support", "minecraft", "logs")
		fmt.Print("Enter your minecraft logs directory (default: " + DARWIN_DEFAULT_PATH + "): ")
		var logLoc string
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "/") {
			logLoc = DARWIN_DEFAULT_PATH
		}
		runScanner([]string{logLoc})
	} else if runtime.GOOS == "windows" {
		var logLoc string
		WINDOWS_DEFAULT_PATH = filepath.Join(homeDir, "AppData", "Roaming", ".minecraft", "logs")
		fmt.Print("Enter your minecraft logs directory (default: " + WINDOWS_DEFAULT_PATH + "): ")
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "\\") {
			logLoc = WINDOWS_DEFAULT_PATH
		}
		runScanner([]string{logLoc})
	}
}
