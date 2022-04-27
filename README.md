# DRRT Scripts

This repository contains a couple script files that I use to run Debris Regional Reassembly Tournaments.
These scripts were made to work on my machine only, so don't expect them to work instantly without modifications.

#### Usage

These scripts were meant to run **only on Windows in WSL2.** The [MatchMaker](https://idleloop.com/matchmaker/) program, used by the scripts to generate a MATCH SCHEDULE, does have a MacOS and some linux versions, but I haven't tested it whatsoever. Neither the `DRRTassembler` nor the `DRRTscheduler` have been run on any machine other than my own.

My DRRT directory looks like this before running the scripts, where `#` represents the number of the DRRT being held.

```
DRRT/
| DRRT_#/
| | oldships/
| | drrt-scripts/
| | | DRRTassembler*
| | | DRRTscheduler*
| | | shipIndex.conf
| | Reassembly_${SUBMISSION_NAME}_[by_${AUTHOR_NAME}]_DRRT_#_${REV}.lua.gz
| Matchmaker/
| | MatchMaker.exe*
```
The `drrt-scripts/` directory is this repository.

The `oldships/` directory is for submission files that participants wanted to switch out for a different ship or a newer version. If the participant switches out an old submission with a new submission with the same name, I increment the `$REV` value. `$REV` starts at `A`, then goes through the alphabet (`B`, `C`, etc.).

The `Reassembly_${SUBMISSION_NAME}_[by_${AUTHOR_NAME}]_DRRT_#_${REV}.lua.gz` file shows the location of all submission files. They are named according to this format. `${SUBMISSION_NAME}` is the submission name given to me by the participant. `${AUTHOR_NAME}` is the author name given to me by the participant. The `${REV}` value is the version of the submission, explained under the description of the `oldships/` directory.

The `Matchmaker/` directory is where the MatchMaker program is stored. The MatchMaker's functionality is explained [here.](https://idleloop.com/matchmaker/) The version of the MatchMaker program doesn't seem to matter. I am currently using Version 1.4 for Windows. DRRT 6, the most recent DRRT as of writing this, used this version.

#### DRRTscheduler

The **DRRTscheduler** script, referred to as simply **"the Scheduler"**, runs the MatchMaker and generates ALLIANCES according to its output.

Basic usage 
```
./DRRTscheduler [-r|--rounds NUMBER_OF_ROUNDS] [-t|-s|--ships|--teams NUMBER_OF_SHIPS] [-v] [-h] [-b] [-f]
```
Run the MatchMaker and check out the options and what they mean. This will show you what arguments you need to pass to the Scheduler.

The Scheduler ends by adding the files `rawSchedule.txt` and `spreadsheetSCH.txt` to the `drrt-scripts/` directory. It also creates all fleet files for the Red ALLIANCE and Blue ALLIANCE, named according to the match number. It creates those files with the Assembler.

The Scheduler sources the file `shipIndex.conf` so that it can pass the right ship filenames to the Assembler.

If the file `rawSchedule.txt` already exists, the Scheduler will ask if you would like to use it to make ALLIANCES instead of calling the MatchMaker and creating a brand new schedule.
If the file `rawSchedule.txt` already exists AND you pass no arguments to the Scheduler, the script will use it to generate ALLIANCES without prompting you.

Now that I've mentioned the Assembler:

#### DRRTassembler

The **DRRTassembler** script, referred to as simply **"the Assembler"**, takes 6 individual Reassembly Fleet files (`*.lua.gz` extension) and creates a RED ALLIANCE fleet file and a BLUE ALLIANCE fleet file, according to the order of arguments.

It requires 6 arguments:
```
./DRRTassembler $RED1 $RED2 $RED3 $BLUE1 $BLUE2 $BLUE3
```

If you would like to name the generated fleet files, pass two more arguments:
```
./DRRTassembler $RED1 $RED2 $RED3 \
                $BLUE1 $BLUE2 $BLUE3 \
                $RED_ALLIANCE_FILENAME $BLUE_ALLIANCE_FILENAME
```
To save the generated fleet files to a certain directory, just change the `$RED_ALLIANCE_FILENAME` and the `$BLUE_ALLIANCE_FILENAME`. (e.g. the Scheduler passes an absolute path to the Assembler for the filenames)

Note: The Assembler isn't meant to be run by itself. Its only purpose was to be called by the Scheduler repeatedly. I might make it a function in the Scheduler eventually.

When the Assembler is called by the Scheduler, the `$RED_ALLIANCE_FILENAME` and the `$BLUE_ALLIANCE_FILENAME` are formatted like this:
```
Match ${MATCH_NUMBER} - ${COLOR} Alliance.lua.gz
```
where `${MATCH_NUMBER}` is the MATCH number that the ALLIANCE plays in and ${COLOR} is either "Red" or "Blue", depending on the intended color of the ALLIANCE.
This assures that the ALLIANCES are ordered correctly in Reassembly's fleet import screen while keeping the names simple and easy to understand.
