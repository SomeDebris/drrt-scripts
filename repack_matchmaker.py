#!/usr/bin/env python3
"""
DRRT MATCHMAKER REPACKER
Repacks TXT output schedules from the MatchMaker execuatable
into CSVs which are used by the DRRT ALLIANCE SCHEDULER.
"""

import argparse
import csv
import os
import re

from drrt_common import wait_yn, print_err


def main(args):
    #TODO verbose output option
    # Paths are: ./schedules/[raw/out]/
    base_sch_dir = os.path.join(os.getcwd(), 'schedules')
    raw_sch_dir = os.path.join(base_sch_dir, 'raw')
    out_sch_dir = os.path.join(base_sch_dir, 'out')

    # Check that input file directory exists - error if true
    if not os.path.exists(raw_sch_dir):
        print_err(f'{raw_sch_dir} is not a folder that exists!')

    # Check that output file directory exists - create if false
    if not os.path.exists(out_sch_dir):
        os.makedirs(out_sch_dir, exist_ok=True)
    
    # Input schedule dir contains one subdirectory per alliance option: 2v2, 3v3, 4v4
    for alliance_dir_name in os.listdir(raw_sch_dir):
        # Iterate through all files (all .txt) in the alliance subdirectory
        alliance_dir_path = os.path.join(raw_sch_dir, alliance_dir_name)
        for txt_fn in os.listdir(alliance_dir_path):
            # Open schedule input file (output from MatchMaker) and read all lines
            with open(os.path.join(alliance_dir_path, txt_fn), 'r') as ifh:
                lines = ifh.readlines()
            # Get schedule lines, which have format "\n [0-9 ][0-9]: <data>\n"
            # Data is of format: "\d \d ..."
            sch_raw = re.findall(r'\n [0-9 ][0-9]:(.*)\n', '\n'.join(lines))
            # Find all space-separated numbers within each line and add to output schedule
            # Some lines may contain an asterisk '*' after them, keep these
            schedule = [re.findall(r'(?:\s(\d+\*?))+', s) for s in sch_raw]

            # Output csv file to "./schedules/out/<name>.csv"
            # Maintain same name convention as the input file
            #   i.e. 'teams_NvN.csv' where N is the number of teams per alliance
            csv_fn = f'{txt_fn[:-4]}.csv'
            out_sch_path = os.path.join(out_sch_dir, csv_fn)
            # If the file already exists and checking is not disabled, ask to overwrite
            if not args.no_check and os.path.exists(out_sch_path):
                if not wait_yn(f'{out_sch_path} is a file that already exists. Overwrite?'):
                    print_err(f'{out_sch_path} exists, not overwriting and continuing.', True)
                    continue

            # Write processed schedule to output CSV file
            with open(out_sch_path, 'w') as ofh:
                csv_writer = csv.writer(ofh, lineterminator='\n')
                csv_writer.writerows(schedule)


def parse_args():
    parser = argparse.ArgumentParser(description=__doc__, 
            formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('-v', '--verbose', 
            action='store_true', 
            help='Enables verbose output.')
    parser.add_argument('--no-check', 
            action='store_true', 
            help='Prevent overwrite checking.')
    return parser.parse_args()


if __name__ == "__main__":
    main(parse_args())
