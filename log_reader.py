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

from drrt_common import DATA_DIR, SCRIPT_DIR

REASSEMBLY_DATA = os.path.join(os.path.expanduser('~'), '.local', 'share', 'Reassembly', 'data')
LATEST_MLOG = os.path.join(REASSEMBLY_DATA, 'match_log_latest.txt')



# def main(args):
"""
First:
    - starts looking at the Reassembly folder
    - whenever a 1 is passed into a named pipe, do a thing!
"""
MLOG_SIGNAL_PIPE = '/tmp/drrt_mlog_signal_pipe'

def main():
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
                    read_latest_mlog()
                    break
                print('Read: "{0}"'.format(data))
    
def read_latest_mlog():
    # mlogs = [filename for filename in os.listdir(REASSEMBLY_DATA) if filename.startswith('MLOG')]
    # print(mlogs)

def read_latest_mlog_symlink():
    latest_mlog = os.path.join(REASSEMBLY_DATA, 'match_log_latest.txt')
    latest_mlog_file = open(latest_mlog, 'r')

    latest_mlog_content = latest_mlog_file.read()

    latest_mlog_file.close()

    print(latest_mlog_content)
    
if __name__ == '__main__':
    main()
