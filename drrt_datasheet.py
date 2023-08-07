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

SCOPES = ['https://www.googleapis.com/auth/spreadsheets.readonly']

DRRT_DATASHEET_ID = '1rzksRiVxzHi6ukZpWg0OV-U27fCTnkRNALUtAi5iB7w'
DRRT_RANGE_NAME = 'PyTest1!A2:E

# https://stackoverflow.com/questions/68859429/how-to-append-data-in-a-googlesheet-using-python
def main():
    append_to_sheet([['test1','test2']], 'PyTest1!A1')


def get_service()
    credentials = None
    if os.path.exists('token.json'):
        credentials = Credentials.from_authorized_user_file('token.json', SCOPES)

    if not credentials or not credentials.valid:
        if credentials and credentials.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            flow = InstalledAppFlow.from_client_secrets_file(
                    'credentials_drrt.json', SCOPES)
            credentials = flow.run_local_server(port=0)
        with open('token.json', 'w') as token:
            token.write(credentials.to_json())

    return build('sheets', 'v4', credentials=credentials)

def append_to_sheet(values, sheet_range, sheet_id=DRRT_DATASHEET_ID):
    service = get_service()
    body = {'values':values}

    result = service.spreadsheets().values().append(
            spreadsheetId=sheet_id, range=sheet_range,
            valueInputOption="RAW", body=body).execute()


if __name__ == '__main__':
    main()
