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


# def main(args):
"""
First:
    - starts looking at the Reassembly folder
    - whenever a 1 is passed into a named pipe, do a thing!
"""
MLOG_SIGNAL_PIPE = '/tmp/drrt_mlog_signal_pipe'

def main():

    mlogs = [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]
    mlog_initial_count = len(mlogs)

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
        ship_info = line.split("|")
        ship_list.append({ 'name':ship_info[0],'author':ship_info[1],'filename':ship_info[2] })
    
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
    field_regex = re.compile("(\w+):\{(.+?)\}")
    id_regex = re.compile("\[([A-Z]+)\]")
    for line in mlog_content.splitlines():
        if not line:
            continue
        fields = re.findall(field_regex, line)
        print(fields)
        message_id = re.findall(id_regex, line)
        print(fields)
     

    

    
if __name__ == '__main__':
    main()
