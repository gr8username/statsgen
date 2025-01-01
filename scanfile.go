package main

import (
	"compress/gzip"
	"strings"
	"os"
	"fmt"
	"bufio"
	"regexp"
)

func scanLogFile(path string, state *ScannerState) {
	startedWizGame := false
	// it is difficult to determine with certainty whether kill messages actually
	// come from wizards, as the syntax is very non-unique ("You killed <PlayerName>!")
	// and the script will confuse it with a Wizards kill. Thus, we check if the log file shows joining wizards or selecting a kit, and in that case, we can have a greater degree of certainty that the kill actually comes from Wizards.
	// fmt.Println("Opening " + path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	usingPattern := regexp.MustCompile("You are using the (Kinetic|Blood|Arcane|Toxic|Fire|Ancient|Ice|Wither|Storm|Hydro) Wizard kit")
	// In theory, someone could write "You are using the Blood Wizard kit" in chat, and it would screw up the program.
	// A fix would require parsing the starting portion of the logs across multiple versions of MC. Though this is absolutely incomplete parsing
	// there are no security issues with doing this, the most someone could do is annoy you by screwing up your statistics a little, so for now, we'll leave this without
	// checking beginning of the line
	switchedPattern := regexp.MustCompile("You will respawn as (Kinetic|Blood|Arcane|Toxic|Fire|Ancient|Ice|Wither|Storm|Hydro) Wizard next time!")
	killedPat := regexp.MustCompile("You killed [A-Za-z0-9_]{3,16}!")
	killedByPat := regexp.MustCompile("You were killed by [A-Za-z0-9_]{3,16}!")
	defer reader.Close()
	// now, we go line by line
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if !isChatLine(line) {
			continue
		}
		if usingPattern.MatchString(line) {
			startedWizGame = true
			state.currentClass = getClass(line)
			// fmt.Println("Class started as " + CLASS_NAMES[state.currentClass])
		}
		if switchedPattern.MatchString(line) {
			startedWizGame = true
			state.currentClass = getClass(line)
			// fmt.Println("Class changed to " + CLASS_NAMES[state.currentClass])
		}
		if startedWizGame {
			if killedPat.MatchString(line) {
				name := extractName(line)
				state.logKill(name)
			}
			if killedByPat.MatchString(line) {
				name := extractName(line)
				state.logDeath(name)
			}
		}
	}
}

func extractName(line string) string {
	// for You killed and You were killed by
	user := ""
	for i := len(line)-1 ; i >= 0 && line[i] != ' '; i-- {
		if line[i] != '!' {
			user = string(line[i]) + user
		}
	}
	return user 
}

func getClass(line string) int {
	// first, parse the name out
	pat := regexp.MustCompile("(Kinetic|Blood|Arcane|Toxic|Fire|Ancient|Ice|Wither|Storm|Hydro)")
	classTypeStr := pat.FindString(line)
	return getClassIndex(classTypeStr)
}

func isChatLine(line string) bool {
	return strings.Contains(line, "[CHAT]")
}
