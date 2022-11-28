# DRRT Scripts

This repository contains a couple script files that I use to run Debris Regional Reassembly Tournaments.
These scripts should be able to run on any machine, provided that you have python installed.

The main use of these scripts is to select a MATCH SCHEDULE and generate QUALIFICATIONS ALLIANCE fleet files for each MATCH in the selected SCHEDULE.

### Directory Setup

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

The `schedules/` directory is where all pregenerated match schedules are stored, formatted in csv. Currently, this repository contains schedules for **2v2 tournaments** (4 to 100 ships), **3v3 tournaments** (6 to 100 ships), and **4v4 tournaments** (8 to 100 ships). All schedules were generated with the [MatchMaker.](https://idleloop.com/matchmaker/) These schedules were generated with MatchMaker v1.5 using the best quality setting.

### drrt_scheduler.py

The **drrt_scheduler.py** script, referred to simply as **"the Scheduler"**, selects a schedule based on the given ALLIANCE length and creates all ALLIANCE fleet files required for the QUALIFICATIONS.

Basic usage 
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
