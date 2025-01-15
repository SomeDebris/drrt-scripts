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
GAME_TEMPLATE_PATH = os.path.join( SCRIPT_DIR, 'html', 'game_TEMPLATE.html')

"""
function that creates two text files when called with 
a match number: 
    RED ALLIANCE NAME
    BLUE ALLIANCE NAME
"""
# def print_ships_at_qualification_match( number_match, ship_list, ranked_ship_list, output_suffix):
#     template = ''
#     with open( GAME_TEMPLATE_PATH, 'r' ) as vic:
#         template = vic.read()

#     filepath_selected_schedule = os.path.join( SCRIPT_DIR, 'selected_schedule.csv' )

#     if ( not os.path.exists( filepath_selected_schedule ) ):
#         print_err( f"drrt_overlay: I can't find the selected_schedule.csv file!" )
    
#     target_match = []

#     with open( filepath_selected_schedule ) as file_schedule:
#         schedule_reader = csv.reader( file_schedule )
        
#         target_match = [ row for idx, row in enumerate( schedule_reader ) if idx == number_match - 1 ][0]
    
#     ships_in_match = []

#     for ship_number in target_match:
#         ship_int = int(ship_number) - 1

#         ships_in_match.append( ship_list[ ship_int ] )
    
#     filename = os.path.join( SCRIPT_DIR, 'html', f"qual{number_match}_game_filled.html" )

#     if os.path.exists( filename ):
#         os.remove(filename)

#     with open( filename, 'w') as f:
#         idx = 0

#         output_string = template.format( rank_red1 = match_info[0]['ships'][0]['rank'],
#                       rank_red1_box = 'rank_neutral_captain' if match_info[0]['ships'][0]['rank'] <= 8 else 'rank_neutral',
#                       name_red1 = strip_author_from_ship_name(match_info[0]['ships'][0]['name']),
#                       author_red1 = match_info[0]['ships'][0]['author'],
#                       rp_red1 = '+' + str(match_info[0]['ships'][0]['RPs']),
#         f.write(output_string)

#     print( "check my work, boss!")

def print_next_html( all_ships, match_number ):
    template = ''
    with open( GAME_TEMPLATE_PATH, 'r' ) as vic:
        template = vic.read()

    filepath_selected_schedule = os.path.join( SCRIPT_DIR, 'selected_schedule.csv' )

    if ( not os.path.exists( filepath_selected_schedule ) ):
        print_err( f"drrt_overlay: I can't find the selected_schedule.csv file!" )
    
    target_match = []

    with open( filepath_selected_schedule ) as file_schedule:
        schedule_reader = csv.reader( file_schedule )
        
        target_match = [ row for idx, row in enumerate( schedule_reader ) if idx == match_number - 1 ][0]
        
    ships_in_match = []

    for ship_number in target_match:
        for rank, ship1 in enumerate(all_ships):
            if (ship1[ 'sub_order' ] + 1 == int(ship_number)):
                ship1[ 'rank' ] = rank + 1
                ships_in_match.append( ship1 )
    
    content = template.format( rank_red1 = ships_in_match[0]['rank'],
                      rank_red1_box = 'rank_neutral_captain' if ships_in_match[0]['rank'] <= 8 else 'rank_neutral',
                      name_red1 = ships_in_match[0]['name'],
                      author_red1 = ships_in_match[0]['author'],

                      rank_red2 = ships_in_match[1]['rank'],
                      rank_red2_box = 'rank_neutral_captain' if ships_in_match[1]['rank'] <= 8 else 'rank_neutral',
                      name_red2 = ships_in_match[1]['name'],
                      author_red2 = ships_in_match[1]['author'],

                      rank_red3 = ships_in_match[2]['rank'],
                      rank_red3_box = 'rank_neutral_captain' if ships_in_match[2]['rank'] <= 8 else 'rank_neutral',
                      name_red3 = ships_in_match[2]['name'],
                      author_red3 = ships_in_match[2]['author'],

                      rank_blue1 = ships_in_match[3]['rank'],
                      rank_blue1_box = 'rank_neutral_captain' if ships_in_match[3]['rank'] <= 8 else 'rank_neutral',
                      name_blue1 = ships_in_match[3]['name'],
                      author_blue1 = ships_in_match[3]['author'],

                      rank_blue2 = ships_in_match[4]['rank'],
                      rank_blue2_box = 'rank_neutral_captain' if ships_in_match[4]['rank'] <= 8 else 'rank_neutral',
                      name_blue2 = ships_in_match[4]['name'],
                      author_blue2 = ships_in_match[4]['author'],

                      rank_blue3 = ships_in_match[5]['rank'],
                      rank_blue3_box = 'rank_neutral_captain' if ships_in_match[5]['rank'] <= 8 else 'rank_neutral',
                      name_blue3 = ships_in_match[5]['name'],
                      author_blue3 = ships_in_match[5]['author'])

    with open( os.path.join(SCRIPT_DIR, 'html', 'game_filled.html'), 'w') as f:
        f.write(content)

def print_victory_html( match_info, all_ships ):
    temp = ''
    with open( VICTORY_TEMPLATE_PATH, 'r' ) as vic:
        temp = vic.read()
    
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
    for rank, ship1 in enumerate(all_ships):
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][0]['name']):
            match_info[0]['ships'][0]['author'] = ship1['author']
            match_info[0]['ships'][0]['ranking_score'] = ship1['ranking_score']
            match_info[0]['ships'][0]['rank'] = rank + 1
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][1]['name']):
            match_info[0]['ships'][1]['author'] = ship1['author']
            match_info[0]['ships'][1]['ranking_score'] = ship1['ranking_score']
            match_info[0]['ships'][1]['rank'] = rank + 1
        if ship1['name'] == strip_author_from_ship_name(match_info[0]['ships'][2]['name']):
            match_info[0]['ships'][2]['author'] = ship1['author']
            match_info[0]['ships'][2]['ranking_score'] = ship1['ranking_score']
            match_info[0]['ships'][2]['rank'] = rank + 1

        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][0]['name']):
            match_info[1]['ships'][0]['author'] = ship1['author']
            match_info[1]['ships'][0]['ranking_score'] = ship1['ranking_score']
            match_info[1]['ships'][0]['rank'] = rank + 1
        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][1]['name']):
            match_info[1]['ships'][1]['author'] = ship1['author']
            match_info[1]['ships'][1]['ranking_score'] = ship1['ranking_score']
            match_info[1]['ships'][1]['rank'] = rank + 1
        if ship1['name'] == strip_author_from_ship_name(match_info[1]['ships'][2]['name']):
            match_info[1]['ships'][2]['author'] = ship1['author']
            match_info[1]['ships'][2]['ranking_score'] = ship1['ranking_score']
            match_info[1]['ships'][2]['rank'] = rank + 1
    
    content = temp.format(red_victory_txt = red_vic_txt,
                      blue_victory_txt = blue_vic_txt,
                      rank_red1 = match_info[0]['ships'][0]['rank'],
                      rank_red1_box = 'rank_neutral_captain' if match_info[0]['ships'][0]['rank'] <= 8 else 'rank_neutral',
                      name_red1 = strip_author_from_ship_name(match_info[0]['ships'][0]['name']),
                      author_red1 = match_info[0]['ships'][0]['author'],
                      rp_red1 = '+' + str(match_info[0]['ships'][0]['RPs']),

                      rank_red2 = match_info[0]['ships'][1]['rank'],
                      rank_red2_box = 'rank_neutral_captain' if match_info[0]['ships'][1]['rank'] <= 8 else 'rank_neutral',
                      name_red2 = strip_author_from_ship_name(match_info[0]['ships'][1]['name']),
                      author_red2 = match_info[0]['ships'][1]['author'],
                      rp_red2 = '+' + str(match_info[0]['ships'][1]['RPs']),

                      rank_red3 = match_info[0]['ships'][2]['rank'],
                      rank_red3_box = 'rank_neutral_captain' if match_info[0]['ships'][2]['rank'] <= 8 else 'rank_neutral',
                      name_red3 = strip_author_from_ship_name(match_info[0]['ships'][2]['name']),
                      author_red3 = match_info[0]['ships'][2]['author'],
                      rp_red3 = '+' + str(match_info[0]['ships'][2]['RPs']),

                      rank_blue1 = match_info[1]['ships'][0]['rank'],
                      rank_blue1_box = 'rank_neutral_captain' if match_info[1]['ships'][0]['rank'] <= 8 else 'rank_neutral',
                      name_blue1 = strip_author_from_ship_name(match_info[1]['ships'][0]['name']),
                      author_blue1 = match_info[1]['ships'][0]['author'],
                      rp_blue1 = '+' + str(match_info[1]['ships'][0]['RPs']),

                      rank_blue2 = match_info[1]['ships'][1]['rank'],
                      rank_blue2_box = 'rank_neutral_captain' if match_info[1]['ships'][1]['rank'] <= 8 else 'rank_neutral',
                      name_blue2 = strip_author_from_ship_name(match_info[1]['ships'][1]['name']),
                      author_blue2 = match_info[1]['ships'][1]['author'],
                      rp_blue2 = '+' + str(match_info[1]['ships'][1]['RPs']),

                      rank_blue3 = match_info[1]['ships'][2]['rank'],
                      rank_blue3_box = 'rank_neutral_captain' if match_info[1]['ships'][2]['rank'] <= 8 else 'rank_neutral',
                      name_blue3 = strip_author_from_ship_name(match_info[1]['ships'][2]['name']),
                      author_blue3 = match_info[1]['ships'][2]['author'],
                      rp_blue3 = '+' + str(match_info[1]['ships'][2]['RPs']))
                      
    with open( os.path.join(SCRIPT_DIR, 'html', 'victory_filled.html'), 'w') as f:
        f.write(content)



