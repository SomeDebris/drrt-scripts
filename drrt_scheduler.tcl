#!/usr/bin/env tclsh

source drrt_common.tcl

puts "Heyy this worked!"

# What does this script need to accomplish in order to truly function as the 
# DRRT Scheduler?
# - Create correct folder structure. This should likely be in a specified 
#   location
# - Reference a match schedule file and index the array based on it
# - alliance assembler function
proc checkShipFileExtension {filename} {
    if {![file exists $filename]} {
        error "File \"$filename\" cannot be found!"
    }

    switch -glob -nocase -- "$filename" {
        *.json {
            return {json}
        }
        *.json.gz {
            return {jsongz}
        }
        *.lua {
            return {lua}
        }
        *.lua.gz {
            return {luagz}
        }
        default {
            error "File extension of $filename should be that of a ship file."
        }
    }
}

proc makeShipDictFromFile {filename} {
    set shipfiletype [checkShipFileExtension $filename]

    switch $shipfiletype {
        json {
            set filehandle [open "$filename" r]
            set filecontents [read $filehandle]
            close $filehandle
            # I can do this because we operate DIRECTLY on json (pretty cool)
            return $filecontents
        }
        jsongz {
            set filehandle [open "$filename" rb]
            set filecontents [read $filehandle]
            close $filehandle

            set filecontents_uncompressed [zlib inflate $filecontents]

            return $filecontents_uncompressed
        }
        lua {
            error "Lua ship files cannot be parsed yet. Please re-export your ships as JSON."
        }
        luagz {
            error "gzipped Lua ship files cannot be parsed yet. Please re-export your ships as JSON."
        }
        default {
            error "Ship file \"$filename\" type found to be \"$shipfiletype\", but \"$shipfiletype\" has no dict creation procedure defined."
        }
    }
}

proc saveFleetToFile {filename fleet} {
}
    

proc makeShipsIntoFleet {ships} {

}

proc removeBlockDataFromShip {ship_dict} {
}
