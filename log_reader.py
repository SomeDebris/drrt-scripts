#!/usr/bin/env/ python3
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

from drrt_common import DATA_DIR, SCRIPT_DIR

REASSEMBLY_DATA = os.path.join(os.path.expanduser('~'), '.local', 'share', 'Reassembly', 'data')
LATEST_MLOG = os.path.join(REASSEMBLY_DATA, 'match_log_latest.txt')

RED_ALLIANCE_TITLE_COLORS = [0xbaa01e, 0x681818, 0x000000]
BLUE_ALLIANCE_TITLE_COLORS = [0x0aa879, 0x222d84, 0x000000]

Current_Match_ID = 0

ALL_SHIPS = {}


# def main(args):
"""
First:
    - starts looking at the Reassembly folder
    - whenever a 1 is passed into a named pipe, do a thing!
"""
MLOG_SIGNAL_PIPE = '/tmp/drrt_mlog_signal_pipe'

def main():
    global ALL_SHIPS

    mlogs = [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]
    mlog_initial_count = len(mlogs)

    ALL_SHIPS = get_ship_list()

    try:
        os.mkfifo(MLOG_SIGNAL_PIPE)
    except OSError as oe:
        if oe.errno != errno.EEXIST:
            raise

    while True:
        print("opening MLOG_SIGNAL_PIPE...")
        with open(MLOG_SIGNAL_PIPE) as mlog_signal_pipe:
            print("mlog signal pipe opened!")
            while True:
                data = mlog_signal_pipe.read()
                if len(data) == 0:
                    print("writer closed!")
                    break
                print('Read: "{0}"'.format(data))
                mlogs = read_latest_mlog(mlogs)
                parse_mlog(read_latest_mlog_symlink())
                # print(json.dumps(get_ship_list(), sort_keys=True, indent=4))

def read_latest_mlog(previous_known_mlogs):
    mlog_initial_count = len(previous_known_mlogs)
    current_known_mlogs = [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]
    mlog_final_count = len(current_known_mlogs) - mlog_initial_count

    return current_known_mlogs

def get_ship_list():
    """
    load all ships from ship_index.txt into set of objects
    """
    ship_list = []
    ship_index_file = open(os.path.join(SCRIPT_DIR, "ship_index.txt"), 'r')
    ship_index_content = ship_index_file.read()

    for line in ship_index_content.splitlines():
        if not line:
            continue
        ship_info = line.split("|")
        ship_list.append({ 'name':ship_info[0],'author':ship_info[1],'filename':ship_info[2],
                          'D':0, 'P':0, 'L':0, 'S':0, 'rank':0.0, 'RPs':0})
    
    ship_index_file.close()

    return ship_list

def read_latest_mlog_symlink():
    latest_mlog_content = None
    
    if (os.path.exists(LATEST_MLOG)):
        latest_mlog_file = open(LATEST_MLOG, 'r')
        latest_mlog_content = latest_mlog_file.read()
        latest_mlog_file.close()
        print(latest_mlog_content)
        return latest_mlog_content
    else:
        print("can't find \"{}\"!".format(LATEST_MLOG))
        return False

    # latest_mlog_content shall be PARSED to heck and back
def parse_mlog(mlog_content):
    """
    returns a JSON String of what occured in the match described
    by the input mlog.
    
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
    global Current_Match_ID

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

        print(message_id[0])
        print(fields_dict)
        
        if message_id[0] == 'START':
            if (fields_dict['fleet'] == '0'):
                # It's the alliance on the RIGHT
                red_alliance['name'] = fields_dict['name']
                # TODO game should output fleet colors
            else:
                blue_alliance['name'] = fields_dict['name']

        elif message_id[0] == 'SHIP':
            if (fields_dict['fleet'] == '0'):
                red_ships.append({ 
                                  'name':fields_dict['ship'], 
                                  'destroyed':True,
                                  'RPs':0,
                                  'destructions':0
                                  })
                red_ship_index[fields_dict['ship']] = red_ship_index_length
                red_ship_index_length += 1
            else:
                blue_ships.append({ 
                                  'name':fields_dict['ship'], 
                                  'destroyed':True,
                                  'RPs':0,
                                  'destructions':0
                                  })
                blue_ship_index[fields_dict['ship']] = blue_ship_index_length
                blue_ship_index_length += 1

        elif message_id[0] == 'DESTRUCTION':
            if (fields_dict['fship'] == '100'):
                print(red_ships[ red_ship_index[ fields_dict['ship'] ] ])
                red_ships[ red_ship_index[ fields_dict['ship'] ] ]['RPs'] += 1
                red_ships[ red_ship_index[ fields_dict['ship'] ] ]['destructions'] += 1
                print(red_ships[ red_ship_index[ fields_dict['ship'] ] ])
            else:
                print(blue_ships[ blue_ship_index[ fields_dict['ship'] ] ])
                blue_ships[ blue_ship_index[ fields_dict['ship'] ] ]['RPs'] += 1
                blue_ships[ blue_ship_index[ fields_dict['ship'] ] ]['destructions'] += 1
                print(blue_ships[ blue_ship_index[ fields_dict['ship'] ] ])

        elif message_id[0] == 'RESULT':
            mlog_completion += 1
            if (fields_dict['fleet'] == '0'):
                red_alliance['damageTaken'] = int(fields_dict['DT'])
                red_alliance['damageInflicted'] = int(fields_dict['DI'])
                red_alliance['survivorCount'] = int(fields_dict['alive'])
            else:
                blue_alliance['damageTaken'] = int(fields_dict['DT'])
                blue_alliance['damageInflicted'] = int(fields_dict['DI'])
                blue_alliance['survivorCount'] = int(fields_dict['alive'])

        elif message_id[0] == 'SURVIVAL':
            if (fields_dict['fleet'] == '0'):
                red_ships[ red_ship_index[ fields_dict['ship'] ] ]['destroyed'] = False
            else:
                blue_ships[ blue_ship_index[ fields_dict['ship'] ] ]['destroyed'] = False
        else:
            print("well, {}'s apparently not in my list!".format(message_id[0]))


    if (not (mlog_completion >= 1)):
        print("mlog not complete! Cannot continue.")
        return

    red_score = red_alliance['damageTaken']
    blue_score = blue_alliance['damageTaken']

    if (red_score >= blue_score):
        match_info['winner'] = 'red'
        for blue_ship in blue_ships:
            blue_ship['deltaL'] = 1
    else:
        match_info['winner'] = 'blue'
        for red_ship in red_ships:
            red_ship['deltaL'] = 1

    if (red_score == 0):
        for blue_ship in blue_ships:
            blue_ship['deltaD'] = 1
    elif (blue_score == 0):
        for red_ship in red_ships:
            red_ship['deltaD'] = 1

    elif (match_info['winner'] == 'red'):
        for red_ship in red_ships:
            red_ship['deltaP'] = 1
    else: 
        for blue_ship in blue_ships:
            blue_ship['deltaP'] = 1
        
    for blue_ship in blue_ships:
        if (~blue_ship['destroyed']):
            blue_ship['deltaS'] = 1

    for red_ship in red_ships:
        if (~red_ship['destroyed']):
            red_ship['deltaS'] = 1



    red_alliance['ships'] = red_ships
    blue_alliance['ships'] = blue_ships

    print(json.dumps(red_alliance, indent=4))
    print(json.dumps(blue_alliance, indent=4))
    
if __name__ == '__main__':
    main()
