## TNT Wizards stats generator

This is a utility written in go designed to generated a text file of statistics for kills and deaths against specific
players. It works by analyzing Minecraft log files, most commonly found in ~/.minecraft/logs/\*.gz.

If you use a custom Minecraft client (such as Lunar), or have set a different Minecraft directory, you may need to change the default filepath shown when running the utility.

## Current Status
* Tested on Windows and Linux, not tested on Mac OS/Darwin/BSD.

## Known Issues
* If a player changes their username, they will show up twice in the statistics as separate players.

## Install steps (Linux)
1. [Install golang](https://go.dev/doc/install).
2. Open a terminal.
3. run `git clone https://github.com/gr8username/statsgen/`.
4. run `cd statsgen`.
5. run `go build`.
6. run `./statsgen`.
7. To view stats.txt (or the filename you selected) type ```less stats.txt``` Or you can open it with a graphical text editor.

## Install steps (Windows)
1. [Install golang](https://go.dev/doc/install).
2. [Download repo as ZIP file](https://github.com/gr8username/statsgen/archive/refs/heads/main.zip).
3. **Extract the files**.
4. Open `Windows Terminal`
5. Run `cd <PATH TO EXTRACTED FILES>`.
6. Run `go build`.
7. When finished (when the prompt shows back up) run `.\statsgen.exe`).
8. Follow the prompts (leave blank if defaults okay, press enter after typing response).
9. Run `notepad stats.txt`.

Note, you may find the path to the extracted files by dragging a folder into the Terminal.


Example stats.txt output.

```
Recorded Kills: 9112
Recorded Deaths: 4152
Note, the above numbers are not necessarily the actual amount of kills and deaths on your account.
This script can only possibly read data from logs, so if you, for example, play Minecraft on another computer, any logs of Wizards kills will be missing.
This file is most likely thousands of lines long, it is recommended to use a text editor with search capability
Additionally, the class section of the tables define what class you were using when you died or got a kill, it does not and cannot determine which kit the other player was using.

                     -----STATS AGAINST PLAYER1-----                     
You have killed PLAYER1 252 times
PLAYER1 has killed you 390 times
Individual K/D: 0.65
                              CLASS BREAKDOWN                              
Class         Killed by PLAYER1             Killed PLAYER1              
Kinetic       21                            12                            
Blood         119                           78                            
Arcane        3                             2                             
Toxic         102                           93                            
Fire          1                             0                             
Ancient       0                             0                             
Ice           137                           66                            
Wither        6                             0                             
Storm         1                             0                             
Hydro         0                             1                             


                     -----STATS AGAINST PLAYER2-----                     
You have killed PLAYER2 170 times
PLAYER2 has killed you 389 times
Individual K/D: 0.44
                              CLASS BREAKDOWN                              
Class         Killed by PLAYER2             Killed PLAYER2              
Kinetic       9                             1                             
Blood         170                           66                            
Arcane        0                             0                             
Toxic         116                           64                            
Fire          12                            7                             
Ancient       0                             0                             
Ice           70                            31                            
Wither        0                             0                             
Storm         0                             0                             
Hydro         12                            1                             
...                             
```
