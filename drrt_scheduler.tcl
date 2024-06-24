#!/usr/bin/env tclsh

package require rl_json
package require Tk

# Converts a dict object to a json object
# Taken from the tcl docs
# proc dict2json {dict_to_encode} {
#     set accumulator {}
#     ::json::write object {*}[dict map {k v} $dict_to_encode {
#         set v [::json::write string $v]
#     }]
# }

### User interface
set padthick 5

wm title . "DRRT Alliance Assembler"

ttk::frame .c -padding "5 5 5 5"

set past_export_directories [list]
set fleets_in_export_list [list]

ttk::combobox .c.directory_entry -textvariable export_directory
ttk::button .c.export_button -text "Export Fleets" -command exportFleetsFromGui
ttk::label .c.directory_label -text "Export directory:"
ttk::label .c.fleetlist_label -text "Fleets:"

ttk::treeview .c.fleetlist_tree -columns {author points}
.c.fleetlist_tree heading author -text "Author Name"
.c.fleetlist_tree heading points -text "P Total"

ttk::button .c.browse_button -text "Browse..." -command browseForExportDirectoryGui

ttk::button .c.new_fleet -text "New" -command createNewFleetGui
ttk::button .c.edit_fleet -text "Edit" -command editSelectedFleetGui
ttk::button .c.remove_fleet -text "Remove" -command removeSelectedFleetGui

grid .c -column 0 -row 0 -sticky nsew

grid .c.directory_entry -column 0 -row 5 -sticky we 
grid .c.export_button -column 1 -row 6 -pady "$padthick 0"
grid .c.directory_label -column 0 -row 4 -sticky swe -pady "$padthick 0"
grid .c.fleetlist_label -column 0 -row 0 -sticky swe
grid .c.fleetlist_tree  -column 0 -row 1 -rowspan 3 -sticky nsew

grid .c.browse_button -column 1 -row 5

grid .c.new_fleet  -column 1 -row 1
grid .c.edit_fleet -column 1 -row 2
grid .c.remove_fleet -column 1 -row 3 -sticky s

grid columnconfigure . 0 -weight 1
grid rowconfigure . 0 -weight 1

grid columnconfigure .c 0 -weight 3 -pad $padthick

grid rowconfigure .c 1 -weight 0
grid rowconfigure .c 2 -weight 0
grid rowconfigure .c 3 -weight 1

wm geometry . "400x200"

### Constants
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

proc browseForExportDirectoryGui {} {
    set dir [tk_chooseDirectory]

    if {$dir ne ""} {
        set ::export_directory $dir

        appendToExportCombobox $::export_directory
    }
}

proc appendToExportCombobox {input} {
    lappend ::past_export_directories $input

    set ::past_export_directories [lsearch -all -inline -not -exact $::past_export_directories {}]

    ::.c.directory_entry configure -values $::past_export_directories
}

proc createNewFleetGui {} {
    set id [.c.fleetlist_tree insert {} end -text "Unnamed Fleet" \
        -values [list "Unknown Author" 0]]

    # TODO: This list gets appended to every time we add a fleet; this may want
    # to be a dictionary. It may also just be possible to use the list of stuff
    # in the tree as the actual variable!
    lappend ::fleets_in_export_list $id

    puts $id
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

<<<<<<< HEAD
    set new_blocks_array [::rl_json::json amap block [::rl_json::json extract $ship_json "blocks"] {
=======
    if {![info exists ship_json]} {
        error "Variable with name $ship_json_varname does not exist!"
    }

    set new_blocks_array [::rl_json::json array]

    ::rl_json::json foreach block [::rl_json::json extract $ship_json "blocks"] {
>>>>>>> 688db09 (add check for if variable ship_json actually exists)
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

proc shipName {ship_json_varname} {
    upvar 1 $ship_json_varname ship_json

    return [::rl_json::json get $ship_json "data" "name"]
}

proc shipAuthor {ship_json_varname} {
    upvar 1 $ship_json_varname ship_json

    return [::rl_json::json get $ship_json "data" "author"]
}

