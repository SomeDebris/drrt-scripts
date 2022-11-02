#!/usr/bin/env python3
"""
#TODO help
"""

import argparse
import csv
import os
import re

from drrt_common import wait_yn, print_err


def main(args):
    #TODO documentation
    #TODO verbose output option
    base_sch_dir = os.path.join(os.getcwd(), 'schedules')
    raw_sch_dir = os.path.join(base_sch_dir, 'raw')
    out_sch_dir = os.path.join(base_sch_dir, 'out')
    if not os.path.exists(raw_sch_dir):
        raise OSError(f'{raw_sch_dir} is not a folder that exists!')

    if not os.path.exists(out_sch_dir):
        os.makedirs(out_sch_dir, exist_ok=True)
    
    for alliance_dir_name in os.listdir(raw_sch_dir):
        alliance_dir_path = os.path.join(raw_sch_dir, alliance_dir_name)
        for txt_fn in os.listdir(alliance_dir_path):
            with open(os.path.join(alliance_dir_path, txt_fn), 'r') as ifh:
                lines = ifh.readlines()
            sch_raw = re.findall(r'\n [0-9 ][0-9]:(.*)\n', '\n'.join(lines))
            schedule = [re.findall(r'(?:\s(\d+\*?))+', s) for s in sch_raw]

            csv_fn = f'{txt_fn[:-4]}.csv'
            out_sch_path = os.path.join(out_sch_dir, csv_fn)
            if not args.no_check and os.path.exists(out_sch_path):
                if not wait_yn(f'{out_sch_path} is a file that already exists. Overwrite?'):
                    print_err(f'{out_sch_path} exists, not overwriting and continuing.', True)
                    continue

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
