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

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(RED + "Error getting home directory!" + RESET)
		fmt.Println(err)
	}
	LINUX_DEFAULT_PATH = filepath.Join(homeDir, ".minecraft", "logs")
	if runtime.GOOS == "linux" || runtime.GOOS == "bsd" {
		fmt.Print("Enter your minecraft logs directory (default: " + LINUX_DEFAULT_PATH + "): ")
		var logLoc string
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "/") {
			logLoc = LINUX_DEFAULT_PATH
		}
		runScanner(logLoc)
	} else if runtime.GOOS == "darwin" {
		DARWIN_DEFAULT_PATH = filepath.Join(homeDir, "Library", "Application Support", "minecraft", "logs")
		fmt.Print("Enter your minecraft logs directory (default: " + DARWIN_DEFAULT_PATH + "): ")
		var logLoc string
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "/") {
			logLoc = DARWIN_DEFAULT_PATH
		}
		runScanner(logLoc)
	} else if runtime.GOOS == "windows" {
		var logLoc string
		WINDOWS_DEFAULT_PATH = filepath.Join(homeDir, "AppData", "Roaming", ".minecraft", "logs")
		fmt.Print("Enter your minecraft logs directory (default: " + WINDOWS_DEFAULT_PATH + "): ")
		fmt.Scanln(&logLoc)
		if !strings.Contains(logLoc, "\\") {
			logLoc = WINDOWS_DEFAULT_PATH
		}
		runScanner(logLoc)
	}
}
