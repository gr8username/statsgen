package main

import (
	"path/filepath"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const RED = "\033[31m"
const RESET = "\033[0m"

const (
	KINETIC = iota
	BLOOD
	ARCANE
	TOXIC
	FIRE
	ANCIENT
	ICE
	WITHER
	STORM
	HYDRO
)

var CLASS_NAMES = []string{"Kinetic", "Blood", "Arcane", "Toxic", "Fire", "Ancient", "Ice", "Wither", "Storm", "Hydro"}

type ScannerState struct {
	totalKills int
	totalDeaths int
	currentClass int
	scannedFile bool
	kills map[string]PlayerRelation
}

func (state *ScannerState) logKill(playerName string) {
	// playerName is the player by whom user was killed
	// fmt.Println("Logged kill of " + playerName)
	state.initPlayer(playerName)
	state.totalKills++
	state.kills[playerName].killAs[state.currentClass]++
}

func (state *ScannerState) initPlayer(playerName string) {
	_, ok := state.kills[playerName]
	if state.kills == nil {
		state.kills = make(map[string]PlayerRelation)
	}
	if !ok {
		var placeHolder PlayerRelation
		placeHolder.killAs = make(map[int]int)
		placeHolder.deathAs = make(map[int]int)
		placeHolder.name = playerName
		state.kills[playerName] = placeHolder
	}
}

func (state *ScannerState) logDeath(playerName string) {
	state.initPlayer(playerName)
	state.totalDeaths++
	// fmt.Println("Death by " + playerName + " logged!")
	state.kills[playerName].deathAs[state.currentClass]++
}

func (this *PlayerRelation) deathsByPlayer() int {
	sum := 0
	for _, k := range this.deathAs {
		sum += k
	}
	return sum
}

func (this *PlayerRelation) killsByUser() int {
	sum := 0
	for _, k := range this.killAs {
		sum += k
	}
	return sum
}

func (this *PlayerRelation) deathsAndKills() int {
	return this.killsByUser()+this.deathsByPlayer()
}

func centerText(text string, length int) string {
	length -= len(text)
	retVal := text
	for ; length > 0 ; length -= 2 {
		retVal = " " + retVal
		retVal += " "
	}
	retVal += "\n"
	return retVal
}

func (this *PlayerRelation) ToString() string {
	returnValue := "\n\n"
	returnValue += centerText("-----STATS AGAINST " + this.name + "-----", 74)
	returnValue += fmt.Sprintf("You have killed %s %d times\n", this.name, this.killsByUser())
	returnValue += fmt.Sprintf("%s has killed you %d times\n", this.name, this.deathsByPlayer())
	returnValue += fmt.Sprintf("Individual K/D: %.2f\n", this.getKD())
	returnValue += centerText("CLASS BREAKDOWN", 74)
	returnValue += fmt.Sprintf("%-14s%-30s%-30s\n", "Class", "Killed by " + this.name, "Killed " + this.name)
	for i := KINETIC; i <= HYDRO ; i++ {
		returnValue += fmt.Sprintf("%-14s%-30d%-30d\n", CLASS_NAMES[i], this.deathAs[i], this.killAs[i])
	}
	return returnValue
}

type PlayerRelation struct {
	killAs map[int]int  // number of times user killed player as X class
	deathAs map[int]int  // number of times user was killed by player as X class
	name string
}

func (this *PlayerRelation) getKD() float64 {
	return safeDivide(float64(this.killsByUser()), float64(this.deathsByPlayer()))
} 

func safeDivide(dividend float64, divisor float64) float64 {
	if dividend == float64(0) || divisor == float64(0) {
		return float64(0)
	}
	return dividend/divisor
}

func getClassIndex(className string) int {
	for i := 0 ; i < len(CLASS_NAMES) ; i++ {
		if CLASS_NAMES[i] == className {
			return i
		}
	}
	return 0
}

func isLogfile(file os.DirEntry) (bool, error) {
	// fmt.Println(file.Name())
	return regexp.MatchString("[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9]?\\.log\\.gz", file.Name())
}

func runScanner(logDirs []string) {
	var state ScannerState
	state.currentClass = KINETIC
	for _, dir := range logDirs {
		runScannerOnDir(dir, &state)
	}
	rets := sortPlayerRelations(state.kills)
	// fmt.Println(rets[8].ToString())
	genReportFile(state, rets)
}

func runScannerOnDir(logDir string, state *ScannerState) {
	fmt.Println("Starting scan of " + logDir + " for log files.")
	
	dirData, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Println(RED +"ERROR: DIRECTORY MOST LIKELY DOES NOT EXIST OR PERMISSION DENIED" + RESET)
		fmt.Println(err)
		return
	}
	// sortEntries(dirData)
	if len(dirData) == 0 {
		fmt.Println(RED + "ERROR: No logfiles available!" + RESET)  // Print red error message
		return
	}
	for i := 0 ; i < len(dirData) ; i++ {
		logF, _ := isLogfile(dirData[i])
		fullPath := filepath.Join(logDir + string(os.PathSeparator), dirData[i].Name())
		if !logF {
			// fmt.Printf("Not logfile: %s\n", fullPath)
			if dirData[i].IsDir() {
				// recursion is fun to write, and absolutely nothing can ever go wrong, right?
				newPath := filepath.Join(logDir + string(os.PathSeparator), dirData[i].Name())
				runScannerOnDir(newPath, state)
			}
			continue
		}	
		// fmt.Printf("Found logfile: %s\n", fullPath)
		scanLogFile(fullPath, state)
		state.scannedFile = true
	}
}

func genReportFile(state ScannerState, rets []PlayerRelation) {
	var loc string
	if PRESET_STATS == "" {
		fmt.Print("Enter the filename under which you would like to save stats (default: stats.txt): ")
		fmt.Scanln(&loc)
	} else {
		loc = PRESET_STATS
	}
	if loc == "" {
		loc = "stats.txt"
	}
	if !strings.HasSuffix(loc, ".txt") {
		loc += ".txt"
	}
	repFile, err := os.Create(loc)
	if err != nil {
		panic(err)
	}
	defer repFile.Close()
	introductionStr := fmt.Sprintf("Recorded Kills: %d\nRecorded Deaths: %d\nNote, the above numbers are not necessarily the actual amount of kills and deaths on your account.\nThis script can only possibly read data from logs, so if you, for example, play Minecraft on another computer, any logs of Wizards kills will be missing.\nThis file is most likely thousands of lines long, it is recommended to use a text editor with search capability\nAdditionally, the class section of the tables define what class you were using when you died or got a kill, it does not and cannot determine which kit the other player was using.", state.totalKills, state.totalDeaths)
	repFile.WriteString(introductionStr)
	for i := 0 ; i < len(rets) ; i++ {
		repFile.WriteString(rets[i].ToString())
	}

}


func sortPlayerRelations(data map[string]PlayerRelation) []PlayerRelation {
	// sort in DESCENDING order
	slice := make([]PlayerRelation, 0, len(data))
	for _, item := range data {
		slice = append(slice, item)
	}
	for i := 0 ; i < len(data)-1; i++ {
		maximum := i
		for j := i+1 ; j < len(data) ; j++ {
			if slice[j].deathsAndKills() > slice[maximum].deathsAndKills() {
				maximum = j
			} 
		}
		temp := slice[maximum]
		slice[maximum] = slice[i]
		slice[i] = temp
	}
	return slice
}

func sortEntries(data []os.DirEntry) {
	// sorts through in ASCENDING order
	// yes, I know there are better sorting algorithms than selection sort,
	// this is just quick and easy to code
	for i := 0 ; i < len(data)-1 ; i++ {
		minimum := i
		for j := i+1; j < len(data) ; j++ {
			if data[j].Name() < data[i].Name() {
				minimum = j
			}
		}
		swapDirSlice(data, i, minimum)
	}
}

func swapDirSlice(data []os.DirEntry, ind1 int, ind2 int) {
	temp := data[ind1]
	data[ind1] = data[ind2]
	data[ind2] = temp
}
