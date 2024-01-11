#!/usr/bin/env python3
"""
DRRT DATASHEET CONNECTION
Connects to the DRRT Datasheet and does stuff.
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

from drrt_common import VERSION, TOURNAMENT_DIRECTORY, SCRIPT_DIR, print_err, wait_yn

"""
function that creates two text files when called with 
a match number: 
    RED ALLIANCE NAME
    BLUE ALLIANCE NAME
"""
def print_ships_at_match( number_match ):
    filepath_selected_schedule = os.path.join( SCRIPT_DIR, 'selected_schedule.csv' )

    if ( not os.path.exists( filepath_selected_schedule ) ):
        print_err( f"drrt_overlay: I can't find the selected_schedule.csv file!" )
    
    target_match = []

    with open( filepath_selected_schedule ) as file_schedule:
        schedule_reader = csv.reader( file_schedule )
        
        target_match = [ row for idx, row in enumerate( schedule_reader ) if idx in (number_match - 1) ]

    print(target_match)
        



