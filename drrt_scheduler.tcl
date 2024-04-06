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
            error "File extension of $filename should be legal for a ship file, but wasn't!"
        }
    }
}

proc makeShipDictFromFile {filename} {
    set shipfiletype [checkShipFileExtension $filename]

    set filehandle [open "$filename" r]
    set filecontents [read $filehandle]
    close $filehandle

    switch $shipfiletype {
        json {
            return -code 0 [makeShipDictFromJson filecontents]
        }
        jsongz {
            error "gzipped JSON files cannot be parsed yet. Please gunzip them."
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

proc makeShipDictFromJson {varname} {
    upvar 1 $varname jsonstring

    set shipdict {[::json::json2dict $jsonstring]}

    return $shipdict
}

proc saveFleetToFile {filename fleet} {
}

proc makeShipsIntoFleet {ships} {

}

proc removeBlockDataFromShip {ship_dict} {
}
