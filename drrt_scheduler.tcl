#!/usr/bin/env tclsh

source drrt_common.tcl

puts "Heyy this worked!"

set Block_Keys_To_Keep [list ident offset angle bindingId faction command]

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

    set filecontents {}

    switch $shipfiletype {
        json {
            set filehandle [open "$filename" r]
            set filecontents [read $filehandle]
            close $filehandle
            # I can do this because we operate DIRECTLY on json (pretty cool)
        }
        jsongz {
            set filehandle [open "$filename" r]
            zlib push gunzip $filehandle
            set filecontents [read $filehandle]
            close $filehandle
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

    return $filecontents
}

proc isShipJSON {ship_json_varname} {
    upvar 1 $ship_json_varname ship_json

    return [::rl_json::json exists $ship_json "data"]
}


proc saveFleetToFile {filename fleet} {
}
    

proc makeShipsIntoFleet {ships} {

}

proc removeBlockDataFromShip {ship_json_varname \
    {keys_to_keep {ident offset angle bindingId faction command}} } {
    upvar 1 $ship_json_varname ship_json

    set new_blocks_array [::rl_json::json array]

    ::rl_json::json foreach block [::rl_json::json extract $ship_json "blocks"] {
        ::rl_json::json foreach {key value} $block {
            if {[lsearch -exact $keys_to_keep $key] < 0} {
                ::rl_json::json unset block $key
            }
        }
        puts $block
    }
}

proc getShipsFromFleetJSON {fleet_json_varname} {
    upvar 1 $fleet_json_varname fleet_json

    if {[::rl_json::json exists $fleet_json "data"]} {
        error "Ship [::rl_json::json get $fleet_json "data" "name"] is a ship [
                ]file, but was passed to a proc expecting fleet files!"
    }

    set fleet_overview [::rl_json::json get $fleet_json "name"]
    
    set new_ships_array [::rl_json::json array]
}
    

    

