#!/bin/bash

# I tell you whos fighting who and you write it!

while true; do
    read -p "Left alliance?  : " LEFT
    read -p "Right alliance? : " RIGHT

    LEFT_LARGE="./playoffs_hacked/${LEFT}_large.txt"
    LEFT_SMALL="./playoffs_hacked/${LEFT}_small.txt"
    RIGHT_LARGE="./playoffs_hacked/${RIGHT}_large.txt"
    RIGHT_SMALL="./playoffs_hacked/${RIGHT}_small.txt"

    echo "--- LEFT SPAWNING ALLIANCE ---"
    cat "$LEFT_LARGE" > ./red_NEXT.txt
    cat "$LEFT_SMALL" | tee ./red_NEXT_SMALL.txt
    echo "--- RIGHT SPAWNING ALLIANCE ---"
    cat "$RIGHT_LARGE" > ./blue_NEXT.txt
    cat "$RIGHT_SMALL" | tee ./blue_NEXT_SMALL.txt
    echo ""
done
