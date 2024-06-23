#!/usr/bin/env tclsh

package require rl_json

# Converts a dict object to a json object
# Taken from the tcl docs
# proc dict2json {dict_to_encode} {
#     set accumulator {}
#     ::json::write object {*}[dict map {k v} $dict_to_encode {
#         set v [::json::write string $v]
#     }]
# }
set Block_Keys_To_Keep [list ident offset angle bindingId faction command]

proc waitYN {prompt_string} {
    while {true} {
        puts "$prompt_string \[y/n\]: "

        gets stdin response

        switch -nocase "$response" {
            y { return 1 }
            n { return 0 }
        }
    }
}

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

proc makeShipJSONFromFile {filename} {
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

proc sanitizeShipJSON {ship_json_varname {keys_to_keep_varname Block_Keys_To_Keep}} {
    upvar 1 $ship_json_varname ship_json
    upvar 1 $keys_to_keep_varname keys_to_keep

    set new_blocks_array [::rl_json::json amap block [::rl_json::json extract $ship_json "blocks"] {
        ::rl_json::json foreach {key value} $block {
            if {[lsearch -exact $keys_to_keep $key] < 0} {
                ::rl_json::json unset block $key
            }
        }
        set block
    }]

    return [::rl_json::json set ship_json blocks $new_blocks_array]
}

proc getSanitizedShipJSON {ship_json_varname {keys_to_keep_varname Block_Keys_To_Keep}} {
    upvar 1 $ship_json_varname ship_json
    upvar 1 $keys_to_keep_varname keys_to_keep

    set modified_json $ship_json

    return [sanitizeShipJSON modified_json keys_to_keep]
}

# Returns all ships from the json file.
# If this is a single ship file, it will return just the ship.
# If this is a fleet file, it will return an array of each ship.
proc getShipsArrayFromJSON {ships_json_varname} {
    global Block_Keys_To_Keep

    upvar 1 $ships_json_varname ships_json

    set output [::rl_json::json array]

    switch [isShipJSON ships_json] {
        1 {
            ::rl_json::json set output 0 [getSanitizedShipJSON ships_json]
        }
        0 {
            set output [::rl_json::json amap ship [::rl_json::json extract $ships_json "blueprints"] {
                sanitizeShipJSON ship
            }]
        }
    }

    return $output
}

proc shipName {ship_json_varname {default_name "Unnamed Spaceship"}} {
    upvar 1 $ship_json_varname ship_json

    return [::rl_json::json get -default $default_name $ship_json_varname "data" "name"]
}

proc shipAuthor {ship_json_varname {default_name "Unknown Author"}} {
    upvar 1 $ship_json_varname ship_json

    return [::rl_json::json get -default $default_name $ship_json_varname "data" "author"]
}
