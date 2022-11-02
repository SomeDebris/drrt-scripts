#!/usr/bin/env python3
"""
DRRT ALLIANCE SCHEDULER
Generates MATCH SCHEDULE and _assembles all QUALIFICATION MATCH ALLIANCES.
"""

import argparse
import csv
import gzip
import json
import os
import shutil
import subprocess
import sys

from drrt_common import VERSION, DATA_DIR, SCRIPT_DIR, print_err, wait_yn

MATCH_TEMPLATE = """{{     -- Created with DRRTscheduler {0}
  color0=0x0aa879,
  color1=0x222d84,
  color2=0,
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

    ships = _get_participants(not args.no_check)

    if len(ships) < (args.alliances * 2):
        print_err(f'{len(ships)} is lesser than minimum number of ships ({args.alliances * 2}).')

    sch_in_filename = f'{len(ships)}_{args.alliances}v{args.alliances}.csv'
    sch_in_filepath = _get_script_path(os.path.join('schedules', 'out', sch_in_filename))
    sch_out_filepath = _get_script_path('spreadsheetSCH.txt', False)
    sch_noasterisk_filepath = _get_script_path('.no_asterisks', False)

    with open(sch_in_filepath, 'r') as schedule_in, \
            open(sch_out_filepath, 'w') as schedule_out, \
            open(sch_noasterisk_filepath, 'w') as sch_out_noasterisk:
        sch_in_lines = schedule_in.readlines()
        schedule_in.seek(0)
        schedule = [row for row in csv.reader(schedule_in)]
        num_matches = len(sch_in_lines)
        #TODO these may require some modification to line terminators
        schedule_out.writelines(sch_in_lines)
        sch_out_noasterisk.writelines([line.replace('*', '') for line in sch_in_lines])

    print('Beginning ALLIANCE generation.')
    match_num = 1
    while match_num <= num_matches:
        # Add ship files found in the current match schedule
        # Ship numbers in schedule are 1-indexed, match no here is 1-indexed too
        assemble_ships = [ships[int(idx)-1] for idx in schedule[match_num-1]]
        # Check that the correct number of ship files were passed
        if len(assemble_ships) != (2 * args.alliances):
            print('assemble: Not Enough Arguments!')

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
    
    # Check that each ship file exists
    ship_data = []
    for ship in ships:
        if not os.path.exists(ship):
            print_err(f'File {ship} not found!')
        with gzip.open(ship, 'r') as ship_file:
            raw_ship_data = ''.join([b.decode('utf-8') for b in ship_file.readlines()])
        ship_data.append(_parse_ship_data(raw_ship_data))
    
    # Red is the first half of the schedule, blue is the second half
    half_idx = len(ship_data) // 2
    _assemble_alliance(ship_data[:half_idx], red_name)
    _assemble_alliance(ship_data[half_idx:], blue_name)


def _assemble_alliance(ship_data, name):
    with open(os.path.join(DATA_DIR, 'Qualifications', f'{name}.lua'), 'w') as match_file:
        match_file.writelines(MATCH_TEMPLATE.format(VERSION, name, ',\n  '.join(ship_data).replace('\\n', '\n')))


def _parse_ship_data(raw_data):
    data_sense = 0
    start_idx = None
    for idx, char in enumerate(raw_data):
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
    if start_idx is None:
        print_err('Invalid lua ship data file! Cannot find where data starts.')

    end_idx = None
    for idx, char in enumerate(raw_data[start_idx:]):
        if char == '{':
            delim_ctr += 1
        elif char == '}':
            delim_ctr -= 1
            if delim_ctr == 0:
                end_idx = idx + start_idx
                break
    if end_idx is None:
        print_err('Invalid lua ship data file! Cannot find where data ends.')

    return raw_data[start_idx-5:end_idx+1]
    

def _get_script_path(filename, check=True):
    filepath = os.path.join(SCRIPT_DIR, filename)
    if check and not os.path.exists(filepath):
        print_err(f'{filepath} is not a file that exists!')
    return filepath


def parse_args():
    parser = argparse.ArgumentParser(description=__doc__, 
            formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('-v', '--verbose', 
            action='store_true', 
            help='Enables verbose output.')
    parser.add_argument('--no-check', 
            action='store_true', 
            help='Prevent participant checking.')
    parser.add_argument('-a', '--alliances',
            default=3,
            help='') #TODO help for this argument
    #TODO gen-schedule option
    return parser.parse_args()


if __name__ == '__main__':
    main(parse_args())
