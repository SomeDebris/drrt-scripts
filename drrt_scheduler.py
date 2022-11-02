#!/usr/bin/env python3
"""
DRRT ALLIANCE SCHEDULER
Generates MATCH SCHEDULE and assembles all QUALIFICATION MATCH ALLIANCES.
"""

import argparse
import fnmatch
import json
import os
import shutil
import subprocess
import sys
import urllib.request
import zipfile


VERSION = 'v1.3.0'
MM_NAME = 'MatchMaker_1_5_0_b1'

SCRIPT_DIR = os.getcwd()
DATA_DIR = os.path.join(SCRIPT_DIR, '..')
DRRT_ROOT = os.path.join(SCRIPT_DIR, '..', '..')

RED = '\033[0;31m'
YELLOW = '\033[0;33m'
NOCOLOR = '\033[0m'

WARN_ROUNDS = 10


def main(args):
    # Delete files in quals folder if there are any
    quals_path = os.path.join(DATA_DIR, 'Qualifications')
    if len(os.listdir(quals_path)) > 0:
        print('Deleting contents of \'Qualifications/\' . . .')
    shutil.rmtree(quals_path)

    # Create directory structure
    for folder in ('Qualifications', 'Playoffs', 'Old-Ships', 'MatchMaker'):
        path = os.path.join(DATA_DIR, folder)
        if os.path.exists():
            os.makedirs(path, exist_ok=True)

    #TODO make this work cross-platform - DL comes with dmg for mac and exec file for linux + diff arches
    # Find matchmaker executable location if exists (this is messy but it works)
    mm_dir = os.path.join(DATA_DIR, 'MatchMaker')
    mm_exec = '.*MatchMaker.exe' if sys.platform == 'win32' else \
        ('.*MatchMaker.app' if sys.platform == 'darwin' else \
        ('.*/IntelUbuntu/MatchMaker' if platform.machine() in ('AMD64', 'x86_64','x86') 
        else '.*/RaspianArm64/MatchMaker'))
    matchmaker = _find_file(mm_exec, mm_dir)
    # Download matchmaker executable if not
    if not matchmaker:
        print('Attempting to download the MatchMaker (MatchMaker executable). . .')
        mm_zip_path = os.path.join(DATA_DIR, 'MatchMaker', f'{MM_NAME}.zip')
        if not os.path.exists(mm_zip_path):
            urllib.request.urlretrieve(f'https://idleloop.com/matchmaker/{MM_NAME}.zip', mm_zip_path)
            #TODO error handling, L53-58
            # this may do error handling for us?
        with zipfile.ZipFile(mm_zip_path, 'r') as zip_handle:
            zip_handle.extractall(os.path.join(DATA_DIR, 'MatchMaker', MM_NAME))
        
        # Find matchmaker executable, by platform
        if sys.platform == 'win32':
            matchmaker = _find_file('.*MatchMaker.exe', mm_dir)
        elif sys.platform == 'darwin':
            #TODO test to make sure this works
            mm_dmg = _find_file('.*MatchMaker.dmg', mm_dir)
            mm_dmg_mount = os.path.join(mm_dir, 'mount')
            # Mount dmg file to data_dir/MatchMaker/mount/<dmg>
            # make mount folder if it doesn't exist
            os.makedirs(mm_dmg_mount, exist_ok=True)
            subprocess.Popen(['hdiutil', 'attach', '-mountpoint', DATA_DIR, mm_dmg_mount]).wait()
            # Find .app file in mounted dmg
            mm_app_dmg = _find_file('.*MatchMaker.app', mm_dir)
            # Copy app to outside of mount directory
            shutil.copy(mm_app_dmg, mm_dir)
            # Detach mounted dmg and delete mount dir
            subprocess.Popen(['hdiutil', 'detach', mm_dmg_mount]).wait()
            os.removedirs(mm_dmg_mount)
            matchmaker = _find_file('.*MatchMaker.app', mm_dir)
        else:
            # Linx (/UNIX-like other)
            # MM zip include Raspbian, RaspbianArm64, and IntelUbuntu
            # which probably means: x86, ARM64, x64_86
            # TODO test if this works
            matchmaker = _find_file('.*/IntelUbuntu/MatchMaker', mm_dir) \
                if platform.machine() in ('AMD64', 'x86_64','x86') \
                else _find_file('.*/RaspianArm64/MatchMaker', mm_dir)
        
        if not matchmaker:
            _print_err('Could not find MatchMaker executable!')
    else:
        print('MatchMaker executable found.')

    # getopt replaced by argparse, see associated method

    if args.ships < (args.alliances * 2):
        _print_err(f'{args.ships} is lesser than minimum number of ships ({args.alliances * 2}).')

    if args.rounds > WARN_ROUNDS:
        _print_err(f'That\'s a lot of rounds ({args.rounds})! Make sure that all matches were generated!', True)

    ships = get_participants(not args.no_check)
    # Required flag for rounds removes need for error check when rounds set but not ships

    raw_schedule = os.path.join(SCRIPT_DIR, 'rawSchedule.txt')
    # Previous raw schedule file detected
    if os.path.exists(raw_schedule):
        print('A match schedule already exists!')
        # Wait for user to input Y to generate new schedule or n to use current
        regen_resp = _wait_yn('Generate new schedule?')
        if regen_resp:
            os.remove(raw_schedule)
            os.remove(os.path.join(SCRIPT_DIR, 'spreadsheetSCH.txt'))
            run_matchmaker(matchmaker, args.ships, args.rounds, args.alliances)
        else:
            print('Current schedule will be used.')
    else:
        _print_err('rawSchedule.txt not found!', True)
        print(f'{YELLOW}WARNING:{NOCOLOR} ')
        # Do not need to check for ship count here, done in argparse
        print('A new rawSchedule.txt will be generated.')
        run_matchmaker(matchmaker, args.ships, args.rounds, args.alliances)


    print('Beginning ALLIANCE generation.')
    #TODO not sure what this condition should be. $p in original?
    match_num = 1
    while match_num < 1000:
        assemble(ships, args.alliances, 
            f'Match {match_num} - ^1The Red Alliance^7',
            f'Match {match_num} - ^4The Blue Alliance^7')
        match_num += 1

    print('Scheduler done.')
    print('Lets get this tournament started!')
    print('Now import spreadsheetSCH.txt into the DRRT Datasheet.')

    # Open the directory containing spreadsheetSCH.txt in the file browser
    # (only if the user requests that it is opened) - cross-platform
    explorer_resp = _wait_yn('Open the drrt-scripts directory in the file browser?')
    if explorer_resp:
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


def get_participants(check):
    """Read list of all absolute file paths of all participant/ships."""
    # Load json file listing of all ships
    with open('ship_index.json', 'r') as ship_fh:
        participants = json.load(ship_fh)
        if check:
            # If validation, check to make sure each listing is a file that exists
            abs_paths = []
            for ship in participants:
                ship_path = os.path.join(DATA_DIR, ship)
                if not os.path.exists(ship_path):
                    _print_err(f'check_participants: Ship file \'{ship}\' not found!')
                else:
                    abs_paths.append(ship_path)
            print('check_participants: all ship files found!')
            return abs_paths
        else:
            # If no validation, find the absolute path for all given files
            return [os.path.join(DATA_DIR, ship) for ship in participants]


def run_matchmaker(mm_path, num_ships, num_rounds, num_alliances, quality):
    """Runs the MatchMaker, generating a MATCH SCHEDULE."""
    print(f'Creating a schedule with {num_ships} ships each playing in {num_rounds} Rounds.')
    # Set quality to best if passed else fast (should be fast by default)
    quality = '-b' if quality == 'best' else '-f'
    # System call to MatchMaker executable
    subprocess.Popen([mm_path, '-a', num_alliances, '-o', '-t', num_ships, '-r', num_rounds, quality, 
        '>', os.path.join(SCRIPT_DIR, 'rawSchedule.txt')])

    print('rawSchedule.txt generated with MatchMaker executable output.')
    print('Contents of rawSchedule.txt:')
    with open(os.path.join(SCRIPT_DIR, 'rawSchedule.txt'), 'r') as raw_schedule:
        print('\n'.join(raw_schedule.readlines()))


def assemble(ships, ships_per_alliance, red_name='Red Alliance', blue_name='Blue Alliance'):
    """Creates a RED ALLIANCE fleet file and a BLUE ALLIANCE fleet file for a specific match."""
    # Check that the correct number of ship files were passed
    if len(ships) != (2 * ships_per_alliance):
        print('assemble: Not Enough Arguments!')
    
    # Check that each ship file exists
    for ship in ships:
        if not os.path.exists(ship):
            _print_err(f'File {ship} not found!')


    #TODO what does this method actually do???


def _print_err(message, is_warning=False):
    if is_warning:
        print(f'{YELLOW}WARNING:{NOCOLOR} {message}')
    else:
        print(f'{RED}ERROR:{NOCOLOR} {message}')
        print('Stop.')
        sys.exit(1)


def _wait_yn(prompt):
    while True:
        resp = input(f'{prompt} [Y/n]: ')
        if resp == 'n':
            return False
        elif resp == 'Y':
            return True
        else:
            print('Please answer [Y/n].')


def _find_file(pattern, path):
    """Recursively find a file within a path matching a regex pattern."""
    for root, _, files in os.walk(path):
        for name in files:
            if fnmatch.fnmatch(name, pattern):
                return os.path.join(root, name)


def parse_args():
    parser = argparse.ArgumentParser(description=__doc__, 
            formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('-v', '--verbose', 
            action='store_true', 
            help='Enables verbose output.')
    parser.add_argument('-s', '-t', '--ships', '--teams',
            required=True,
            help='Sets number of participating ships/teams.')
    parser.add_argument('-r', '--rounds', 
            default=10, 
            help='Sets minimum number of rounds each ship plays. Defaults to 10.')
    parser.add_argument('-q', 
            nargs='?',
            choices=['fast', 'best'],
            default='fast',
            help='Set schedule generation quality.')
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
