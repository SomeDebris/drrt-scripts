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

from drrt_common import VERSION, TOURNAMENT_DIRECTORY, SCRIPT_DIR, print_err, wait_yn, strip_author_from_ship_name

TEMPLATE = """<span size=\"xx-large\">{0}</span>
{1}
"""

VICTORY_TEMPLATE_PATH = os.path.join( SCRIPT_DIR, 'html', 'victory_TEMPLATE.html')

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
    
    red_filename = os.path.join( SCRIPT_DIR, 'html', f"red_{output_suffix}.html" )
    blue_filename = os.path.join( SCRIPT_DIR, 'html', f"blue_{output_suffix}.html" )

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

def print_victory_html( match_info, all_ships ):
    temp = ''
    with open( VICTORY_TEMPLATE_PATH, 'r' ) as vic:
        temp = vic.read()
    
    print(match_info)
    winner = 0
    if match_info[0]['damageTaken'] < match_info[1]['damageTaken']:
        winner = 1
    
    red_vic_txt  = 'VICTORY!'
    blue_vic_txt = 'LOSES!'
    if winner == 1:
        red_vic_txt = 'LOSES!'
        blue_vic_txt = 'VICTORY!'
   
    # for alliance in match_info:
    #     for ship in alliance['ships']:
    #         if ( not 'author' in ship ):
    #             for participant in all_ships:
    #                 if (participant['name'] == ship['name']):
    #                     ship['author'] = participant['authro
    
    # set the author name
    for ship1 in all_ships:
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][0]['name']):
            match_info[0]['ships'][0]['author'] = ship1['author']
            match_info[0]['ships'][0]['ranking_score'] = ship1['ranking_score']
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][1]['name']):
            match_info[0]['ships'][1]['author'] = ship1['author']
            match_info[0]['ships'][1]['ranking_score'] = ship1['ranking_score']
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][2]['name']):
            match_info[0]['ships'][2]['author'] = ship1['author']
            match_info[0]['ships'][2]['ranking_score'] = ship1['ranking_score']

        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][0]['name']):
            match_info[1]['ships'][0]['author'] = ship1['author']
            match_info[1]['ships'][0]['ranking_score'] = ship1['ranking_score']
        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][1]['name']):
            match_info[1]['ships'][1]['author'] = ship1['author']
            match_info[1]['ships'][1]['ranking_score'] = ship1['ranking_score']
        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][2]['name']):
            match_info[1]['ships'][2]['author'] = ship1['author']
            match_info[1]['ships'][2]['ranking_score'] = ship1['ranking_score']
            
    
    print(match_info)
    print(match_info[0]['ships'][0])
    print(temp.format(red_victory_txt = red_vic_txt,
                      blue_victory_txt = blue_vic_txt,
                      rank_red1 = 1,
                      rank_red1_box = 'rank_neutral',
                      name_red1 = strip_author_from_ship_name(match_info[0]['ships'][0]['name']),
                      author_red1 = match_info[0]['ships'][0]['author'],
                      rp_red1 = '+' + str(match_info[0]['ships'][0]['RPs']),

                      rank_red2 = 1,
                      rank_red2_box = 'rank_neutral',
                      name_red2 = strip_author_from_ship_name(match_info[0]['ships'][1]['name']),
                      author_red2 = match_info[0]['ships'][1]['author'],
                      rp_red2 = '+' + str(match_info[0]['ships'][1]['RPs']),

                      rank_red3 = 1,
                      rank_red3_box = 'rank_neutral',
                      name_red3 = strip_author_from_ship_name(match_info[0]['ships'][2]['name']),
                      author_red3 = match_info[0]['ships'][2]['author'],
                      rp_red3 = '+' + str(match_info[0]['ships'][2]['RPs']),

                      rank_blue1 = 1,
                      rank_blue1_box = 'rank_neutral',
                      name_blue1 = strip_author_from_ship_name(match_info[1]['ships'][0]['name']),
                      author_blue1 = match_info[1]['ships'][0]['author'],
                      rp_blue1 = '+' + str(match_info[1]['ships'][0]['RPs']),

                      rank_blue2 = 1,
                      rank_blue2_box = 'rank_neutral',
                      name_blue2 = strip_author_from_ship_name(match_info[1]['ships'][1]['name']),
                      author_blue2 = match_info[1]['ships'][1]['author'],
                      rp_blue2 = '+' + str(match_info[1]['ships'][1]['RPs']),

                      rank_blue3 = 1,
                      rank_blue3_box = 'rank_neutral',
                      name_blue3 = strip_author_from_ship_name(match_info[1]['ships'][2]['name']),
                      author_blue3 = match_info[1]['ships'][2]['author'],
                      rp_blue3 = '+' + str(match_info[1]['ships'][2]['RPs'])))
                      





