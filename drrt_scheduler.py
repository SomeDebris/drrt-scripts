#!/usr/bin/env python3
"""
DRRT ALLIANCE SCHEDULER
Generates MATCH SCHEDULE and assembles all QUALIFICATION MATCH ALLIANCES.
"""

import argparse
import csv
import gzip
import json
import os
import re
import shutil
import subprocess
import sys

from drrt_common import VERSION, DATA_DIR, SCRIPT_DIR, print_err, wait_yn


ILLEGAL_SHIP_REGEX = '({854,|{863,|{838,|{833,|{273,|{927,|{928,|{929,|{930,|{931,|{932,|{933,|{934,|{935,|{936,|{937,|{938,|{939,|{940,|{941,|{942,|{943,|{953,|{954,|{955,|{956,|{320,|{11104,|{12130,|{15010,|{15142,|{15144,|{15146,)'
MATCH_TEMPLATE = """{{     -- Created with DRRTscheduler {0}
  color0={3},
  color1={4},
  color2={5},
  name=\"{1}\",
  faction=8,
  currentChild=0,
  blueprint={{}},
  children={{}},
  blueprints={{
    {2}
  }},
  playerprint={{}}
}}
"""
RED_ALLIANCE_COLORS = [0xbaa01e, 0x681818, 0x000000]
BLUE_ALLIANCE_COLORS = [0x0aa879, 0x222d84, 0x000000]


def main(args):
    # Delete files in quals folder if there are any
    quals_path = os.path.join(DATA_DIR, 'Qualifications')
    if os.path.exists(quals_path) and len(os.listdir(quals_path)) > 0:
        print('Deleting contents of \'Qualifications/\' . . .')
        shutil.rmtree(quals_path)

    # Create directory structure
    for folder in ('Qualifications', 'Playoffs', 'Old-Ships'):
        path = os.path.join(DATA_DIR, folder)
        if not os.path.exists(path):
            os.makedirs(path, exist_ok=True)

    # Get list of ship/participant filepaths from ship_index.json
    # Also may check if those files exist
    ships = _get_participants(not args.no_check)
    print(f'Found {len(ships)} ships in ship_index.json.')

    # Checks if there are enough ships to fill both alliances at least once
    if len(ships) < (args.alliances * 2):
        print_err(f'{len(ships)} is lesser than minimum number of ships ({args.alliances * 2}).')

    # Gets paths to input schedule CSV, output schedule file, and output schedule file without asterisks
    sch_in_filename = f'{len(ships)}_{args.alliances}v{args.alliances}.csv'
    sch_in_filepath = _get_script_path(os.path.join('schedules', 'out', sch_in_filename))
    sch_out_filepath = _get_script_path('spreadsheetSCH.txt', False)
    sch_noasterisk_filepath = _get_script_path('.no_asterisks', False)

    with open(sch_in_filepath, 'r') as schedule_in, \
            open(sch_out_filepath, 'w') as schedule_out, \
            open(sch_noasterisk_filepath, 'w') as sch_out_noasterisk:
        # Read all lines of input schedule
        sch_in_lines = schedule_in.readlines()
        # Seek back to beginning of file
        schedule_in.seek(0)
        # Re-read input schedule with CSV reader (to get indexed rows)
        schedule = [row for row in csv.reader(schedule_in)]
        # Number of matches (not rounds) = number of lines in the schedule
        num_matches = len(sch_in_lines)
        # Copy input schedule lines to output schedule file
        schedule_out.writelines(sch_in_lines)
        # Write input schedule lines to output noasterisk file, 
        #   but replace all asterisks with nothing
        sch_out_noasterisk.writelines([line.replace('*', '') for line in sch_in_lines])
    print(f'Schedule has {num_matches} matches.')

    print('Beginning ALLIANCE generation.')
    match_num = 1
    while match_num <= num_matches:
        # Add ship files found in the current match schedule
        # Ship numbers in schedule are 1-indexed, match no here is 1-indexed too
        assemble_ships = [ships[int(idx.replace('*', ''))-1] for idx in schedule[match_num-1]]
        # Check that the correct number of ship files were passed
        if len(assemble_ships) != (2 * args.alliances):
            print('assemble: Not Enough Arguments!')

        # Assemble the red and blue alliance match files
        _assemble(assemble_ships, 
            f'Match {str(match_num).zfill(3)} - ^1The Red Alliance^7',
            f'Match {str(match_num).zfill(3)} - ^4The Blue Alliance^7')
        match_num += 1

    print('Scheduler done.')
    print('Lets get this tournament started!')
    print('Now import spreadsheetSCH.txt into the DRRT Datasheet.')

    # Open the directory containing spreadsheetSCH.txt in the file browser
    # (only if the user requests that it is opened) - cross-platform
    if wait_yn('Open the drrt-scripts directory in the file browser?'):
        if sys.platform == 'win32':
            retval = subprocess.Popen(['start', SCRIPT_DIR], shell=True).wait()
        elif sys.platform == 'darwin':
            retval = subprocess.Popen(['open', SCRIPT_DIR]).wait()
        else:
            retval = subprocess.Popen(['xdg-open', SCRIPT_DIR]).wait()
        if retval:
            print
    else:
        print('Stop.')


def _get_participants(check):
    """Read list of all absolute file paths of all participant/ships."""
    # Load json file listing of all ships
    with open('ship_index.json', 'r') as ship_fh:
        participants = json.load(ship_fh)['participants']
        if check:
            # Check that there are one or more ships in the config json
            if len(participants) <= 0:
                print_err(f'check_participants: No ship files found!')

            # If validation, check to make sure each listing is a file that exists
            abs_paths = []
            for ship in participants:
                ship_path = os.path.join(SCRIPT_DIR, 'ships', ship)
                if not os.path.exists(ship_path):
                    print_err(f'check_participants: Ship file \'{ship}\' not found!')
                else:
                    abs_paths.append(ship_path)
            print('check_participants: all ship files found!')
            return abs_paths
        else:
            # If no validation, find the absolute path for all given files
            return [os.path.join(DATA_DIR, ship) for ship in participants]


def _assemble(ships, red_name='Red Alliance', blue_name='Blue Alliance'):
    """Creates a RED ALLIANCE fleet file and a BLUE ALLIANCE fleet file for a specific match."""
    
    ship_data = []
    for ship in ships:
        # Check that each ship file exists
        if not os.path.exists(ship):
            print_err(f'File {ship} not found!')
        # Open gzipped lua ship file (.lua.gz)
        # Read file and decode to single string
        with gzip.open(ship, 'r') as ship_file:
            raw_ship_data = ''.join([b.decode('utf-8') for b in ship_file.readlines()])
            if re.match(ILLEGAL_SHIP_REGEX, raw_ship_data):
                print_err(f'Ship {ship} contains ILLEGAL BLOCKS')
        # Parse ship data out of file content and append to list
        ship_data.append(_parse_ship_data(raw_ship_data))
    
    # Red is the first half of the schedule, blue is the second half
    half_idx = len(ship_data) // 2
    _assemble_alliance(ship_data[:half_idx], red_name)
    _assemble_alliance(ship_data[half_idx:], blue_name)


def _assemble_alliance(ship_data, name):
    """Creates a match file for one ALLIANCE."""
    # Create output file data/Qualifications/<name>.lua
    with open(os.path.join(DATA_DIR, 'Qualifications', f'{name}.lua'), 'w') as match_file:
        # Write match template to file filled out with version, name, and ship data
        # Ship data has escaped \\n in it, replace with \n for newlines
        #   Also join each ship (data field) in the match together with a comma and newline
        match_file.writelines(MATCH_TEMPLATE.format(VERSION, name, ',\n  '.join(ship_data).replace('\\n', '\n')))


def _parse_ship_data(raw_data):
    """Parse and return a ship .lua file for its data field."""
    # raw_data is a string of the ship's data
    data_sense = 0
    start_idx = None
    # Loop through all characters in the ship's data and search for "data="
    for idx, char in enumerate(raw_data):
        # This checks for the characters 'd', 'a', 't', 'a' sequentially
        # Works because data_sense is incremented each letter
        # If a character not in the sequence is found, data_sense is reset to 0
        # So the sequence needs to be found and *then* an '=' character,
        #   which denotes the start of the data block
        if (char == 'd' and data_sense == 0) or \
                (char == 'a' and data_sense == 1) or \
                (char == 't' and data_sense == 2) or \
                (char == 'a' and data_sense == 3):
            data_sense += 1
        elif char == '=' and data_sense == 4:
            delim_ctr = 1
            start_idx = idx
            break
        else:
            data_sense = 0
    # If "data=" is not found, error
    if start_idx is None:
        print_err('Invalid lua ship data file! Cannot find where data starts.')

    end_idx = None
    # If "data=" is found, start at the found index and search for the closing } delimeter
    for idx, char in enumerate(raw_data[start_idx:]):
        # If a { is found, add one to the count
        # If a } is found, subtract one
        # When the count reaches 0, we've closed the data block
        #  (this is because we start at 1 opening brace)
        if char == '{':
            delim_ctr += 1
        elif char == '}':
            delim_ctr -= 1
            if delim_ctr == 0:
                end_idx = idx + start_idx
                break
    # If the delimeter counter never reaches 0, the file is invalid - error
    if end_idx is None:
        print_err('Invalid lua ship data file! Cannot find where data ends.')

    # Return the data found between the two delimiter indices found
    return raw_data[start_idx-len('data='):end_idx+1]
    

def _get_script_path(filename, check=True):
    """Get an OS filepath within the script directory from a filename."""
    filepath = os.path.join(SCRIPT_DIR, filename)
    if check and not os.path.exists(filepath):
        print_err(f'{filepath} is not a file that exists!')
    return filepath


def parse_args():
    #TODO verbose mode
    parser = argparse.ArgumentParser(description=__doc__, 
            formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('-v', '--verbose', 
            action='store_true', 
            help='Enables verbose output.')
    parser.add_argument('--no-check', 
            action='store_true', 
            help='Prevent participant checking.')
    parser.add_argument('-a', '--alliances',
            type=int,
            choices=range(2, 5),
            default=3,
            help='Sets number of ships per alliance.')
    return parser.parse_args()


if __name__ == '__main__':
    main(parse_args())
