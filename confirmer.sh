#!/bin/bash

# Exit on failure
set -e

DRRT_directory="$HOME/Documents/reassembly_ships/tournaments/DRRT/2025 Winter DRRT"

if [ ! -n "$DRRT_directory" ]; then
    echo "Variable \$DRRT_directory is empty: '$DRRT_directory'"
    exit 1
fi


Ships_directory="$DRRT_directory/Ships"
Staging_directory="$DRRT_directory/Staging"


function name_reformat() {
    initial="$1"
    indicator="$2"
    
    new_name=$(basename "$initial" | sed -f "$Staging_directory/name_transformer.sed" )

    echo "$new_name"
}

# we're getting piped filenames
# if [ -p /dev/stdin ]; then


for input_file in $@; do
    # That's right, you need GNU Parallel for this one
    # sem 
    new_filename=$(name_reformat "$input_file" "2025W")

    commit_message=''

    if [ ! -f "${Ships_directory}/${new_filename}" ]; then
        commit_message="New ship file: $new_filename"
    else
        commit_message="Updated ship file: $new_filename"
    fi

    mv "$input_file" "${Ships_directory}/${new_filename}"
    echo "$input_file renamed to $new_filename"

    cd "$Ships_directory"
    git add "$new_filename"

    git commit -m "$commit_message"

    cd "$Staging_directory"
done




