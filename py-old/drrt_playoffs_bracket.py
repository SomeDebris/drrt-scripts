#!/usr/bin/env python3
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
DRRT Playoffs Bracket
This code runs the bracket locally on my pc.
"""

class Alliance:
    self.name = ""
    self.members = []
    self.colors = []
    self.match_record = []




