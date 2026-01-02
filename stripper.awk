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
    code_exit = 0;
    OFS=",";
}

# /^Match Schedule$/      { sched_start = 1; }
/^Schedule Statistics$/ { 
    exit code_exit;
}

/ *[0-9]+: / {
    for (i = 2; i <= NF; i++) {
        printf "%s",$i (i == NF ? ORS : OFS);

        if (match($i, /\*/)) {
            code_exit = 1;
        }
    }
}
