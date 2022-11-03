#!/usr/bin/env python3
"""
DRRT COMMON FUNCTIONS
General functions which are useful for multiple script files.
"""

import fnmatch
import os
import sys


VERSION = 'v1.3.0'
MM_NAME = 'MatchMaker_1_5_0_b1'

SCRIPT_DIR = os.path.dirname(__file__)
DATA_DIR = os.path.abspath(os.path.join(SCRIPT_DIR, os.pardir))
DRRT_ROOT = os.path.abspath(os.path.join(SCRIPT_DIR, os.pardir, os.pardir))

RED = '\033[0;31m'
YELLOW = '\033[0;33m'
NOCOLOR = '\033[0m'


def print_err(message, is_warning=False):
    """Print an error message and exit OR print a warning."""
    if is_warning:
        print(f'{YELLOW}WARNING:{NOCOLOR} {message}')
    else:
        print(f'{RED}ERROR:{NOCOLOR} {message}')
        print('Stop.')
        sys.exit(1)


def wait_yn(prompt):
    """Wait for the user to enter Y or n to a given prompt."""
    while True:
        resp = input(f'{prompt} [Y/n]: ')
        if resp == 'n':
            return False
        elif resp == 'N':
            return False
        elif resp == 'Y':
            return True
        elif resp == 'y':
            return True
        else:
            print('Please answer [Y/n].')


def find_file(pattern, path):
    """Recursively find a file within a path matching a regex pattern."""
    for root, _, files in os.walk(path):
        for name in files:
            if fnmatch.fnmatch(name, pattern):
                return os.path.join(root, name)


if __name__ == "__main__":
    print_err('This is not a script file. Do not run individually.')
