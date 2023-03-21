# DRRT Scripts

This repository contains a couple script files that I use to run Debris Regional Reassembly Tournaments.
These scripts should be able to run on any machine, provided that you have python installed.
The main use of these scripts is to select a MATCH SCHEDULE and generate QUALIFICATIONS ALLIANCE fleet files for each MATCH in the selected SCHEDULE.

## Directory Setup

My DRRT directory looks like this before running the scripts, where `#` represents the number of the DRRT being held.

```
DRRT/
| DRRT_#/
| | Qualifications/
| | drrt-scripts/
| | | drrt_scheduler.py*
| | | drrt_common.py*
| | | repack_matchmaker.py*
| | | ship_index.json
| | | old-ships/
| | | ships/
| | | | Reassembly_${SUBMISSION_NAME}_[by_${AUTHOR_NAME}]_DRRT_#_${REV}.lua.gz
| | | schedules/
```
The `Qualifications/` directory is where the script will output all fleet files.

The `drrt-scripts/` directory is this repository.

The `old-ships/` directory is for submission files that participants wanted to switch out for a different ship or a newer version. If the participant switches out an old submission with a new submission with the same name, I increment the `$REV` value. `$REV` starts at `A`, then goes through the alphabet (`B`, `C`, etc.).

The `ships/` directory is the location in which the participating ships must be stored.

The `Reassembly_${SUBMISSION_NAME}_[by_${AUTHOR_NAME}]_DRRT_#_${REV}.lua.gz` file shows the location of all submission files. I name them according to this format. `${SUBMISSION_NAME}` is the submission name given to me by the participant. `${AUTHOR_NAME}` is the author name given to me by the participant. The `${REV}` value is the version of the submission, explained under the description of the `Old-Ships/` directory.

The `schedules/` directory is where all pregenerated match schedules are stored, formatted in csv. Currently, this repository contains schedules for **2v2 tournaments** (4 to 100 ships), **3v3 tournaments** (6 to 100 ships), and **4v4 tournaments** (8 to 100 ships). All schedules were generated with the [MatchMaker v1.5](https://idleloop.com/matchmaker/) with the best quality setting. To generate schedules with different parameters, see the Schedule Generator script (`gen-schedules`). 

## drrt_scheduler.py

The **drrt_scheduler.py** script, referred to simply as **"the Scheduler"**, selects a schedule based on the given ALLIANCE length and creates all ALLIANCE fleet files required for the QUALIFICATIONS.

### Usage:
```
./drrt_scheduler.py [-v | --verbose] [-h | --help] [-a ALLIANCE_LENGTH]
                    [--no-check]
```
The `-a` option sets the number of ships per ALLIANCE. The default ALLIANCE length is 3. This value may be 2, 3, or 4.

The Scheduler uses the file `ship_index.json` so that it can pass the right ship filenames to the Assembler.

The Scheduler ends by adding the file `selected_schedule.csv` to the `drrt-scripts/` directory. It also creates all fleet files for the Red ALLIANCE and Blue ALLIANCE, named according to the match number.

The output RED and BLUE ALLIANCES' filenames are formatted like this:
```
Match ${MATCH_NUMBER} - ${COLOR} Alliance^7.lua.gz
```
where `${MATCH_NUMBER}` is the MATCH number that the ALLIANCE plays in and ${COLOR} is either "Red" or "Blue", depending on the intended color of the ALLIANCE.
This assures that the ALLIANCES are ordered correctly in Reassembly's fleet import screen while keeping the names simple and easy to understand.
The filename uses Reassembly's color escape characters to color the RED ALLIANCE's name red and the BLUE ALLIANCE's name blue. The `^7` changes the color back to white.

To tell the DRRT DATASHEET (google sheets file that keeps track of score) the match schedule, use the file `selected_schedule.csv`. Place the data in the correct spot in the **calc** tab, change the "Change Me" cell's value, and the sheet should recalculate.

The Scheduler will check the block IDs of all ships it reads. If any block IDs match the IDs of ILLEGAL BLOCKS, the script will let you know.

## gen-schedules

If you would like to generate a new set of schedules for the scheduler to use, use the bash script titled **gen-schedules**, referred to from now on as the Schedule Generator. This script generates a bunch of schedules with the MatchMaker program, formats the raw output of the MatchMaker program as \*.csv files, and puts all the output files into the appropriate directory.

This script needs the MatchMaker to be on the $PATH. It does not attempt to download or install the MatchMaker.

### Usage:
```
./gen-schedules ROUNDS MAX_PARTICIPANTS
```
where ROUNDS is the number of MATCHES each participant will play in and MAX_PARTICIPANTS is the number of participants in the schedule generated with the highest participant count.

The Schedule Generator will make the MatchMaker create schedules for 2v2, 3v3, and 4v4 tournaments using the best quality setting (`-b` option).

In the call to MatchMaker, the Schedule Generator will pass ROUNDS to the `-r` option. The script will call the MatchMaker for all participant counts between the minimum and MAX_PARTICIPANTS. The minimum number of participants is the ALLIANCE length multiplied by 2, as anything lower would not allow a single match to be generated. The MatchMaker will throw an error when called with a participant count below the minimum.

To change the schedule generation parameters, you must edit the script itself. Look for the call to `MatchMaker` and edit the arguments.

## Match Schedule Information

- Smaller ALLIANCE length makes a longer schedule, and vice versa.
- In order to run playoffs, a minimum number of participants must be recieved.
    - 2v2 matches
    4 Playoff ALLIANCES: 8 participants
    8 Playoff ALLIANCES: 16 participants
    - 3v3 matches
    4 Playoff ALLIANCES: 12 participants
    8 Playoff ALLIANCES: 24 participants
    - 4v4 matches
    4 Playoff ALLIANCES: 16 participants
    8 Playoff ALLIANCES: 32 participants

