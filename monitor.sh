#!/bin/bash

WATCHME=$1
REPLAY_DIR=$2
KERCHUNK=$3

inotifywait -m "$WATCHME" -e create |
    while read -r directory action file; do
        echo "New shit just got made: $file"

        case "$file" in
            (*.mkv)
            # if [ "$KERCHUNK" -eq "yes" ]; then
            #     KERCHUNK=''

            #     sleep 1

            #     continue
            # fi
            if [ 3500 -gt $(du "$file" | awk '{ print $1 }') ]; then
                mv "$file" "$REPLAY_DIR"
                
                move_error=$?

                echo "Moved '$file', error code $move_error"
            else
                echo "I don't think thats a replay."
            fi

            ;;
        esac
    done
