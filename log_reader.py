#!/usr/bin/env python3
"""
DRRT MLOG READER
Reads the latest MLOG produced by Reassembly and pipes json to CasparCG
"""
import gzip
import json
import os
import subprocess
import sys
import errno
import re

from drrt_common import TOURNAMENT_DIRECTORY, SCRIPT_DIR, print_err
from drrt_datasheet import append_to_sheet, replace_ships, replace_match_schedule, update_sheet

REASSEMBLY_DATA = os.path.join(os.path.expanduser('~'), '.local', 'share', 'Reassembly', 'data')
LATEST_MLOG = os.path.join(REASSEMBLY_DATA, 'match_log_latest.txt')

RED_ALLIANCE_TITLE_COLORS = [0xbaa01e, 0x681818, 0x000000]
BLUE_ALLIANCE_TITLE_COLORS = [0x0aa879, 0x222d84, 0x000000]
SHIPS_DIRECTORY = os.path.abspath( os.path.join( TOURNAMENT_DIRECTORY, 'Ships' ) )

Current_Match_ID = 0

SHIP_NAME_REGEX = re.compile( r'(.*?) \[.*\]' )

Last_Alliance_Name = {'red':None,'blue':None}

PLAYOFFS_SHEET_SELECTION = [
    'DATA_ENTRY_WINNER_BRACKET!A2',
    'DATA_ENTRY_LOSER_BRACKET!A2'
]

ALL_SHIPS = {}

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


# def main(args):
"""
First:
    - starts looking at the Reassembly folder
    - whenever a 1 is passed into a named pipe, do a thing!
"""
MLOG_SIGNAL_PIPE = '/tmp/drrt_mlog_signal_pipe'

def main():
    global ALL_SHIPS

    mlogs = get_mlog_list()
    mlog_initial_count = len(mlogs)

    ALL_SHIPS = get_ship_list()

    try:
        os.mkfifo(MLOG_SIGNAL_PIPE)
    except OSError as oe:
        if oe.errno != errno.EEXIST:
            raise

    is_playoffs = False
    playoffs_losers = False

    # Open script by checking what type of thing we're doing
    while True:
        print("Are we running QUALS or PLAYOFFS?")
        data = pipe_read(MLOG_SIGNAL_PIPE)
        if (data == 'stop'):
            print("Stop recieved!")
            break
        elif (data == 'qualifications'):
            print(bcolors.BOLD + bcolors.OKBLUE + "--- QUALIFICATIONS ---" + bcolors.ENDC)
            print("Hope you have a good Quals!")
            print("Replacing ship lists in sheet.")
            replace_ships(ALL_SHIPS)
            replace_match_schedule()
            break
        elif (data == 'playoffs'):
            print(bcolors.BOLD + bcolors.HEADER + "--- PLAYOFFS ---" + bcolors.ENDC)
            is_playoffs = True
            print("You're nearly there! Get those playoffs done!")
            break
        else:
            print("Please say 'qualifications', 'playoffs', or 'stop' to {MLOG_SIGNAL_PIPE}.")
            print(f"You said: '{data}'")


    while True:
        data = pipe_read(MLOG_SIGNAL_PIPE)
        if (data == 'stop'):
            print("Stop recieved!")
            break
        if data == 'reload':
            ALL_SHIPS = get_ship_list()
            calculate_all_mlogs(get_mlog_list(), False, is_playoffs, playoffs_losers)
        elif data == 'append':
            parse_mlog(read_latest_mlog_symlink(), True)
        elif data == 'review match':
            print("nothing to do!")
        elif data == 'winners':
            playoffs_losers = False
            print(bcolors.OKGREEN + "Entering data into WINNERS BRACKET" + bcolors.ENDC)
        elif data == 'losers':
            playoffs_losers = True
            print(bcolors.OKCYAN + "Entering data into LOSERS BRACKET" + bcolors.ENDC)
        else:
            print("Please say 'reload', 'review match', 'append', or 'stop' to {MLOG_SIGNAL_PIPE}.")
            print(f"You said: '{data}'")
                # print(json.dumps(get_ship_list(), sort_keys=True, indent=4))

def pipe_read(named_pipe=MLOG_SIGNAL_PIPE):
    print(f"Trying to open '{named_pipe}'...")
    with open(named_pipe) as mlog_signal_pipe:
        print(f"'{named_pipe}' opened!")
        while True:
            data = mlog_signal_pipe.read()
            if len(data) == 0:
                print("writer closed!")
                break
            print('\tRead: "{0}"'.format(data))
            return data
    
    

def read_latest_mlog(previous_known_mlogs):
    mlog_initial_count = len(previous_known_mlogs)
    current_known_mlogs = [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]
    mlog_final_count = len(current_known_mlogs) - mlog_initial_count

    return current_known_mlogs

def get_ship_dict( ship_file_path ):
    if not os.path.exists( ship_file_path ):
        print_err(f'File {ship_file_path} not found!')

    if not ship_file_path.endswith( '.json' ):
        print_err( f"'{ship_file_path}' is NOT a json ship file!" )

    with open( ship_file_path ) as ship_file:
        
        ship_dict = json.loads( ship_file.read() )

        # check whether fleet file or ship file
        if 'blueprints' in ship_dict:
            # This means: THIS IS A FLEET FILE
            # And therefore: it has a 'blueprints' for an array where all the
            # ships are stored
            #ships.append( ship_full[ 'blueprints' ][0] )
            return ship_dict[ 'blueprints' ][0]
        else:
            # This means: THIS IS A SHIP FILE
            # and therefore: it is literally just a ship and can be put in
            # without issue as is
            #ships.append( ship_full )
            return ship_dict



def get_ship_list():
    """
    load all ships from ship_index.txt into set of objects
    """
    ship_list = []
    required_fields = ['name', 'author']

    ship_files = []

    for file in os.listdir( SHIPS_DIRECTORY ):
        if file.endswith( '.json' ):
            filepath = os.path.join( SHIPS_DIRECTORY, file )
            ship_files.append( filepath )
        else:
            print( f"Well, THAT ({file}) is not a ship file!")
    
    for file in sorted( ship_files, key = lambda t: os.stat(t).st_mtime ):
        ship = get_ship_dict( file )[ 'data' ]
        
        # Grabs ONLY the 'required_fields' from the ship dict
        ship_required_information = {x:ship[x] for x in required_fields}

        ship_required_information[ 'filename' ] = file

        # The 0th (blank) group is the whole phrase. 
        # The 1st group is the first thing in parentheses. simple!
        ship_required_information[ 'name' ] = \
            SHIP_NAME_REGEX.search( ship_required_information[ 'name' ] ).group(1)

        ship_list.append( ship_required_information )

    return ship_list

def get_mlog_list():
    return [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]


def calculate_all_mlogs(filenames, check_duplicates, in_playoffs=False, playoffs_losers=False):
    # array of ships, data taken from match log and not
    # cognizant of full match record
    all_ship_match_performances = []
    for filename in filenames:
        file_path = os.path.join(REASSEMBLY_DATA, filename)
        if (os.path.exists(file_path)):
            with open(file_path) as mlog:
                alliances = parse_mlog(mlog.read(), check_duplicates, filename)
                if (alliances):
                    all_ship_match_performances += alliances
                else:
                    print_err("calculate_all_mlogs: For some reason, nothing was returned.", True)
        else:
            print_err("calculate_all_mlogs: can't find '{}'!".format(file_path), True)
    sheet_location = None

    if in_playoffs:
        if playoffs_losers:
            sheet_location = PLAYOFFS_SHEET_SELECTION[1]
        else:
            sheet_location = PLAYOFFS_SHEET_SELECTION[0]
    else:
        sheet_location = None

    datasheet_update_ships(all_ship_match_performances, sheet_location)


def read_latest_mlog_symlink():
    latest_mlog_content = None
    
    if (os.path.exists(LATEST_MLOG)):
        latest_mlog_file = open(LATEST_MLOG, 'r')
        latest_mlog_content = latest_mlog_file.read()
        latest_mlog_file.close()
        return latest_mlog_content
    else:
        print_err("read_latest_mlog_symlink: can't find \"{}\"!".format(LATEST_MLOG), True)
        return False

    # latest_mlog_content shall be PARSED to heck and back

def parse_mlog(mlog_content, check_duplicates, filename="match_log_latest.txt"):
    """
    What do I do?
        - define list of ships participating in match from [SHIP]
          and [START] statements
        - set ships as destroyed if they don't appear in a [SURVIVAL]
          statement
        - give ranking points to each ship

    Loop through each line. Split at FIRST SPACE:
        - first substring = parse for ID
        - second = parse for fields
    """
    global ALL_SHIPS
    global Current_Match_ID
    global Last_Alliance_Name

    is_qual_match = True

    field_regex = re.compile("(\w+):\{(.+?)\}")
    id_regex = re.compile("\[([A-Z]+)\]")

    match_info = {}
    red_alliance = {} 
    blue_alliance = {}
    red_ships = []
    blue_ships = []

    red_ship_index = {}
    red_ship_index_length = 0
    blue_ship_index = {}
    blue_ship_index_length = 0

    mlog_completion = 0

    #check the mlog to see if its complete


    for line in mlog_content.splitlines():
        if not line:
            continue
        fields = re.findall(field_regex, line)
        fields_dict = dict(fields)
        message_id = re.findall(id_regex, line)

        if len(message_id) < 1:
            continue

        if message_id[0] == 'START':
            if (fields_dict[ 'faction' ] == '100'):
                # It's the alliance on the RIGHT
                red_alliance[ 'name' ] = fields_dict[ 'name' ]
                # TODO game should output fleet colors
            else:
                blue_alliance[ 'name' ] = fields_dict[ 'name' ]

        elif message_id[0] == 'SHIP':
            if (fields_dict[ 'fleet' ] == '0'):
                red_ships.append({ 
                                  'name':fields_dict[ 'ship' ], 
                                  'destroyed':True,
                                  'RPs':0,
                                  'destructions':0
                                  })
                red_ship_index[fields_dict[ 'ship' ]] = red_ship_index_length
                red_ship_index_length += 1
            else:
                blue_ships.append({ 
                                  'name':fields_dict[ 'ship' ], 
                                  'destroyed':True,
                                  'RPs':0,
                                  'destructions':0
                                  })
                blue_ship_index[fields_dict[ 'ship' ]] = blue_ship_index_length
                blue_ship_index_length += 1

        elif message_id[0] == 'DESTRUCTION':
            if (fields_dict[ 'fship' ] == '100'):
                if fields_dict[ 'ship' ] in red_ship_index:
                    red_ships[ red_ship_index[ fields_dict[ 'ship' ] ] ][ 'RPs' ] += 1
                    red_ships[ red_ship_index[ fields_dict[ 'ship' ] ] ][ 'destructions' ] += 1
            else:
                if fields_dict[ 'ship' ] in blue_ship_index:
                    blue_ships[ blue_ship_index[ fields_dict[ 'ship' ] ] ][ 'RPs' ] += 1
                    blue_ships[ blue_ship_index[ fields_dict[ 'ship' ] ] ][ 'destructions' ] += 1

        elif message_id[0] == 'RESULT':
            mlog_completion += 1
            if (fields_dict[ 'fleet' ] == '0'):
                red_alliance[ 'damageTaken' ] = int(fields_dict[ 'DT' ])
                red_alliance[ 'damageInflicted' ] = int(fields_dict[ 'DI' ])
                red_alliance[ 'survivorCount' ] = int(fields_dict[ 'alive' ])
            else:
                blue_alliance[ 'damageTaken' ] = int(fields_dict[ 'DT' ])
                blue_alliance[ 'damageInflicted' ] = int(fields_dict[ 'DI' ])
                blue_alliance[ 'survivorCount' ] = int(fields_dict[ 'alive' ])

        elif message_id[0] == 'SURVIVAL':
            if (fields_dict[ 'fleet' ] == '0'):
                red_ships[ red_ship_index[ fields_dict[ 'ship' ] ] ][ 'destroyed' ] = False
            else:
                blue_ships[ blue_ship_index[ fields_dict[ 'ship' ] ] ][ 'destroyed' ] = False
        else:
            print_err("parse_mlog: {}'s apparently not in my list!".format(message_id[0]), True)


    if (not (mlog_completion >= 1)):
        print_err("parse_mlog: mlog not complete! Cannot continue.", True)
        return
    
    same_name = [
            Last_Alliance_Name[ 'red' ] == red_alliance[ 'name' ],
            Last_Alliance_Name[ 'blue' ] == blue_alliance[ 'name' ]
            ]
    if (same_name[0] and same_name[1]):
        print_err("parse_mlog: Both Alliance names Match!", True)
        print_err("parse_mlog: RED: '{}'".format(Last_Alliance_Name[ 'red' ]), True)
        print_err("parse_mlog: BLUE: '{}'".format(Last_Alliance_Name[ 'blue' ]), True)
        print_err("parse_mlog: mlog filename: '{}'".format(filename), True)
        if (check_duplicates):
            print_err("parse_mlog: Not counting match.", True)
            return
    else:
        Last_Alliance_Name[ 'red' ] = red_alliance[ 'name' ]
        Last_Alliance_Name[ 'blue' ] = blue_alliance[ 'name' ]

    red_score = red_alliance[ 'damageTaken' ]
    blue_score = blue_alliance[ 'damageTaken' ]

    if (red_score >= blue_score):
        match_info[ 'winner' ] = 'red'
        for red_ship in red_ships:
            red_ship[ 'RPs' ] += 2
        for blue_ship in blue_ships:
            blue_ship[ 'deltaL' ] = 1
    else:
        match_info[ 'winner' ] = 'blue'
        for blue_ship in blue_ships:
            blue_ship[ 'RPs' ] += 2
        for red_ship in red_ships:
            red_ship[ 'deltaL' ] = 1

    if (red_score == 0):
        for blue_ship in blue_ships:
            blue_ship[ 'deltaD' ] = 1
    elif (blue_score == 0):
        for red_ship in red_ships:
            red_ship[ 'deltaD' ] = 1

    elif (match_info[ 'winner' ] == 'red'):
        for red_ship in red_ships:
            red_ship[ 'deltaP' ] = 1
    else: 
        for blue_ship in blue_ships:
            blue_ship[ 'deltaP' ] = 1
        
    for ship in blue_ships + red_ships:
        if (not ship[ 'destroyed' ]):
            ship[ 'deltaS' ] = 1

    for ship in red_ships:
        ship[ 'fleet_name' ] = red_alliance[ 'name' ]
        ship[ 'enemy_fleet_name' ] = blue_alliance[ 'name' ]
        ship[ 'mlog_filename' ] = filename
    for ship in blue_ships:
        ship[ 'fleet_name' ] = blue_alliance[ 'name' ]
        ship[ 'enemy_fleet_name' ] = red_alliance[ 'name' ]
        ship[ 'mlog_filename' ] = filename

    # all ship's rank and ranking score is calced
    # distribute_points(red_ships + blue_ships)


    # datasheet_append_ships(red_ships + blue_ships)

    # ALL_SHIPS = recalculated_ranks(ALL_SHIPS)

    # what check do I do to ensure that the ranking score is not freaking duplicated
    # APPEND A THING TO EACH 

    red_alliance[ 'ships' ] = red_ships
    blue_alliance[ 'ships' ] = blue_ships

    alliances = (red_alliance, blue_alliance)
    return red_ships + blue_ships

    # print(json.dumps(red_alliance, indent=4))
    # print(json.dumps(blue_alliance, indent=4))
    
    # print(json.dumps(ALL_SHIPS, indent=2))
    
def distribute_points(alliance):
    global ALL_SHIPS

    for participant in ALL_SHIPS:
        for ship in alliance:
            if (participant[ 'name' ] == ship[ 'name' ]):
                participant[ 'RPs' ] += ship[ 'RPs' ]
                if ('deltaD' in ship):
                    participant[ 'D' ] += ship[ 'deltaD' ]
                if ('deltaP' in ship):
                    participant[ 'P' ] += ship[ 'deltaP' ]
                if ('deltaL' in ship):
                    participant[ 'L' ] += ship[ 'deltaL' ]
                if ('deltaS' in ship):
                    participant[ 'S' ] += ship[ 'deltaS' ]
                matches_played = participant[ 'D' ] + participant[ 'P' ] + participant[ 'L' ]
                participant[ 'ranking_score' ] = participant[ 'RPs' ] / matches_played

def recalculated_ranks(ship_array):
    return sorted(ship_array, key=lambda d: d[ 'ranking_score' ], reverse=True) 

def datasheet_update_ships(ships, sheet_range):
    values = []

    for ship in ships:
        if (not 'deltaD' in ship):
            ship[ 'deltaD' ] = 0
        if (not 'deltaP' in ship):
            ship[ 'deltaP' ] = 0
        if (not 'deltaL' in ship):
            ship[ 'deltaL' ] = 0
        if (not 'deltaS' in ship):
            ship[ 'deltaS' ] = 0
        if (not 'mlog_filename' in ship):
            ship[ 'mlog_filename' ] = 'match_log_latest.txt'
        if (not 'fleet_name' in ship):
            ship[ 'fleet_name' ] = 'NONE'
        if (not 'enemy_fleet_name' in ship):
            ship[ 'enemy_fleet_name' ] = 'NONE'
        values.append([ ship[ 'name' ], ship[ 'destructions' ], ship[ 'RPs' ], ship[ 'deltaD' ], ship[ 'deltaP' ], ship[ 'deltaL' ], ship[ 'deltaS' ], ship[ 'fleet_name' ], ship[ 'enemy_fleet_name' ], ship[ 'mlog_filename' ]])

    update_sheet(values, sheet_range)

def datasheet_append_ships(ships):
    values = []

    for ship in ships:
        if (not 'deltaD' in ship):
            ship[ 'deltaD' ] = 0
        if (not 'deltaP' in ship):
            ship[ 'deltaP' ] = 0
        if (not 'deltaL' in ship):
            ship[ 'deltaL' ] = 0
        if (not 'deltaS' in ship):
            ship[ 'deltaS' ] = 0
        if (not 'mlog_filename' in ship):
            ship[ 'mlog_filename' ] = 'match_log_latest.txt'
        if (not 'fleet_name' in ship):
            ship[ 'fleet_name' ] = 'NONE'
        if (not 'enemy_fleet_name' in ship):
            ship[ 'enemy_fleet_name' ] = 'NONE'
        values.append([ SHIP_NAME_REGEX.search( ship[ 'name' ] ).group(1), ship[ 'destructions' ], ship[ 'RPs' ], ship[ 'deltaD' ], ship[ 'deltaP' ], ship[ 'deltaL' ], ship[ 'deltaS' ], ship[ 'fleet_name' ], ship[ 'enemy_fleet_name' ], ship[ 'mlog_filename' ]])

    append_to_sheet(values, 'DATA_ENTRY!A1')
    
    

if __name__ == '__main__':
    main()
