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

from drrt_common import VERSION, DATA_DIR, SCRIPT_DIR, print_err, wait_yn

from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError

SCOPES = ['https://www.googleapis.com/auth/spreadsheets']

DRRT_DATASHEET_ID = '1rzksRiVxzHi6ukZpWg0OV-U27fCTnkRNALUtAi5iB7w'
DRRT_RANGE_NAME = 'PyTest1!A2:E'

SERVICE = None

# https://stackoverflow.com/questions/68859429/how-to-append-data-in-a-googlesheet-using-python
def main():
    replace_match_schedule()


def get_service():
    credentials = None
    if os.path.exists('token.json'):
        credentials = Credentials.from_authorized_user_file('token.json', SCOPES)

    if not credentials or not credentials.valid:
        if credentials and credentials.expired and credentials.refresh_token:
            credentials.refresh(Request())
        else:
            flow = InstalledAppFlow.from_client_secrets_file(
                    'credentials_drrt.json', SCOPES)
            credentials = flow.run_local_server(port=0)
        with open('token.json', 'w') as token:
            token.write(credentials.to_json())

    return build('sheets', 'v4', credentials=credentials)

def append_to_sheet(values, sheet_range, sheet_id=DRRT_DATASHEET_ID):
    global SERVICE

    try:
        if (not SERVICE):
            SERVICE = get_service()
        body = {'values':values}

        result = SERVICE.spreadsheets().values().append(
                spreadsheetId=sheet_id, range=sheet_range,
                valueInputOption="RAW", body=body).execute()
    except HttpError as error:
        print(f"append_to_sheet: An error occured: {error}")
        return error

def replace_ships(ships, sheet_range='Ships!A2:B', sheet_id=DRRT_DATASHEET_ID):
    global SERVICE

    values = []

    for ship in ships:
        if not 'name' in ship or not 'author' in ship:
            print_err("replace_ship: ship {} doesn't have a name or author.".format(ship))

        values.append([ ship['name'], ship['author'] ])

    body = {
        'values': values
    }

    try:
        if (not SERVICE):
            SERVICE = get_service()

        result = SERVICE.spreadsheets().values().update(
                spreadsheetId=sheet_id, range=sheet_range,
                valueInputOption="RAW", body=body).execute()
        print("replace_ships: Updated Ships List!")
    except HttpError as error:
        print(f"replace_ships: An error occured: {error}")
        return error

def replace_match_schedule(sheet_range='Calc!A1:F', sheet_id = DRRT_DATASHEET_ID):
    global SERVICE
    
    values = []
    deletion = []
    deletionRow = []
    with open('selected_schedule.csv', newline='') as csvfile:
        schedule_reader = csv.reader(csvfile)
        for row in schedule_reader:
            values.append(row)
        deletion_row_local = []
        for item in row:
            deletion_row_local.append("")
        deletionRow = deletion_row_local

    for i in range(0, 200):
        deletion.append(deletionRow)
    
    try:
        if (not SERVICE):
            SERVICE = get_service()
        body_values = {
            'values':values
        }
        body_deletion = {
            'values':deletion
        }
        destroy = SERVICE.spreadsheets().values().update(
                spreadsheetId=sheet_id, range=sheet_range,
                valueInputOption="USER_ENTERED", body=body_deletion).execute()
        create = SERVICE.spreadsheets().values().update(
                spreadsheetId=sheet_id, range=sheet_range,
                valueInputOption="USER_ENTERED", body=body_values).execute()
    except HttpError as error:
        print(f"replace_match_schedule: An error occured: {error}")
        return error



if __name__ == '__main__':
    main()
