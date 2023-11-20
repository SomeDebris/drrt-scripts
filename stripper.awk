#!/bin/awk -f

# function no_extension(file) {
#     sub(/\.[^.]*$/, "", file)

#     return file
# }

# function basename(file) {
#     sub(".*/", "", file)
#     return file
# }

BEGIN {
    OFS = ","
}

# /^Match Schedule$/      { sched_start = 1; }
/^Schedule Statistics$/ { nextfile; }

/ *[0-9]*: / {
    for (i = 2; i <= NF; i++ ) {
        printf "%s",$i (i == NF ? ORS : OFS)
    }
}        
