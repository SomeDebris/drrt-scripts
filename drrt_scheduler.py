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

from drrt_common import VERSION, TOURMAMENT_DIRECTORY, SCRIPT_DIR, print_err, wait_yn


ILLEGAL_BLOCK_IDS = [ 854, 863, 838, 833, 273, 927, 928, 929, 930, 931, 932,
                     933, 934, 935, 936, 937, 938, 939, 940, 941, 942, 943,
                     953, 954, 955, 956, 320, 11104, 12130, 15010, 15142,
                     15144, 15146 ] 
FLEET_HEADER = {'name':'Unnamed Alliance',
                'faction':8,
                'currentChild':0,
                'blueprint': {},
                'children': [],
                'blueprints': [],
                'playerprint':{}
                }

RED_ALLIANCE_COLORS =  { 
'color0': '#baa01e',
'color1': '#681818', 
'color2': '#000000' 
}

BLUE_ALLIANCE_COLORS =  { 
'color0': '#0aa879',
'color1': '#222d84', 
'color2': '#000000' 
}


def main( args ):
    # Delete files in quals folder if there are any
    quals_path = os.path.join( TOURMAMENT_DIRECTORY, 'Qualifications' )

    if os.path.exists( quals_path ) and len( os.listdir( quals_path ) ) > 0:
        print( 'Deleting contents of \'Qualifications/\' . . .' )
        shutil.rmtree( quals_path )

    # Create directory structure
    for folder in ( 'Qualifications', 'Playoffs', 'Ships', 'Staging' ):
        path = os.path.join( TOURMAMENT_DIRECTORY, folder )
        if not os.path.exists( path ):
            os.makedirs( path, exist_ok=True )

    # Get list of ship/participant filepaths from ship_index.json
    # Also may check if those files exist
    ships = get_inspected_ship_paths()
    print( f'Found { len( ships ) } ships in ship index.')

    # Checks if there are enough ships to fill both alliances at least once
    if len(ships) < (args.alliances * 2):
        print_err(f'{len(ships)} is lesser than minimum number of ships ({args.alliances * 2}).')

    # Gets paths to input schedule CSV, output schedule file, and output schedule file without asterisks
    sch_in_filename = f'{len(ships)}_{args.alliances}v{args.alliances}.csv'
    sch_in_filepath = _get_script_path(os.path.join('schedules', 'out', sch_in_filename))
    sch_out_filepath = _get_script_path('selected_schedule.csv', False)
    sch_noasterisk_filepath = _get_script_path('.no_asterisks', False)

    with open( sch_in_filepath, 'r' ) as schedule_in, \
            open( sch_out_filepath, 'w' ) as schedule_out, \
            open( sch_noasterisk_filepath, 'w' ) as sch_out_noasterisk:
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

    print(f'Assembling Alliances. . .')
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
    print('Now import selected_schedule.csv into the DRRT Datasheet.')

    # Open the directory containing selected_schedule.csv in the file browser
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



def get_inspected_ship_paths():
    """
    Goal: read all json files from list of arguments or directory and return a
    sorted list of absolute paths based on modification date
    """
    ships_directory = os.path.abspath( os.path.join( TOURMAMENT_DIRECTORY, 'Ships' ) )

    ship_files = []

    for file in os.listdir( ships_directory ):
        if file.endswith( '.json' ):
            ship_files.append( os.path.abspath( file ) )

    # Return a sorted list of files in that directory
    return sorted( ship_files, key = lambda t: os.stat(t).st_mtime )





# Send this a list of ship files
def _assemble(ship_filenames, red_name='Red Alliance', blue_name='Blue Alliance'):
    """Creates a RED ALLIANCE fleet file and a BLUE ALLIANCE fleet file for a specific match."""
    ships = []
    for ship_filename in ship_filenames:
        # Check that each ship file exists
        if not os.path.exists( ship_filename ):
            print_err(f'File {ship_filename} not found!')
        
        with open( ship_filename ) as ship_file:
            ship_file_contents = ship_file.read()
        
        ship_full = json.loads( ship_file_contents )

        # check whether ship is a FLEET FILE or a SHIP FILE, and add the right
        # type of it
        if 'blueprints' in ship_full:
            # This means: THIS IS A FLEET FILE
            # And therefore: it has a 'blueprints' for an array where all the
            # ships are stored
            ships.append( ship_full[ 'blueprints' ][0] )
        else:
            # This means: THIS IS A SHIP FILE
            # and therefore: it is literally just a ship and can be put in
            # without issue as is
            ships.append( ship_full )
    
    # Red is the first half of the schedule, blue is the second half
    length_alliance = len(ships) // 2
    _assemble_alliance(ships[:length_alliance], red_name, RED_ALLIANCE_COLORS)
    _assemble_alliance(ships[length_alliance:], blue_name, BLUE_ALLIANCE_COLORS)


def _assemble_alliance(ships_alliance, name, colors):
    """Creates a match file for one ALLIANCE."""
    # Create output file data/Qualifications/<name>.json

    alliance = dict( colors )
    alliance.update( FLEET_HEADER ) 

    for member in ships_alliance:
        alliance[ 'blueprints' ].append( member )
    
    with gzip.open(os.path.join(TOURMAMENT_DIRECTORY, 'Qualifications', f'{name}.json.gz'), 'wb', encoding='utf-8') as match_file:
        json.dump( alliance, match_file )


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
