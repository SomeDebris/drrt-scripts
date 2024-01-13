#!/usr/bin/env python3
"""
DRRT overlay
Update the DRRT overlay text files from sheet data.
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

TEMPLATE = """<span size=\"xx-large\">{0}</span>
{1}
"""

"""
function that creates two text files when called with 
a match number: 
    RED ALLIANCE NAME
    BLUE ALLIANCE NAME
"""
def print_ships_at_qualification_match( number_match, ship_list, output_suffix, template=TEMPLATE ):
    filepath_selected_schedule = os.path.join( SCRIPT_DIR, 'selected_schedule.csv' )

    if ( not os.path.exists( filepath_selected_schedule ) ):
        print_err( f"drrt_overlay: I can't find the selected_schedule.csv file!" )
    
    target_match = []

    with open( filepath_selected_schedule ) as file_schedule:
        schedule_reader = csv.reader( file_schedule )
        
        target_match = [ row for idx, row in enumerate( schedule_reader ) if idx == number_match - 1 ][0]
    
    ships_in_match = []

    for ship_number in target_match:
        ship_int = int(ship_number) - 1

        ships_in_match.append( ship_list[ ship_int ] )
    
    red_filename = os.path.join( SCRIPT_DIR, f"red_{output_suffix}.txt" )
    blue_filename = os.path.join( SCRIPT_DIR, f"blue_{output_suffix}.txt" )

    if os.path.exists( red_filename ):
        os.remove(red_filename)
    if os.path.exists( blue_filename ):
        os.remove(blue_filename)

    with open( red_filename, 'a' ) as red_file, open( blue_filename, 'a') as blue_file:
        idx = 0

        for ship in ships_in_match:
            ship_name = ship[ 'name' ]

            if len( ship_name ) > 25:
                ship_name = ship_name[:22] + '...'

            output_string = template.format( ship_name, ship[ 'author' ] )

            if (idx < 3):
                red_file.write(output_string)
            else:
                blue_file.write(output_string)

            idx += 1

    print( "check my work, boss!")


