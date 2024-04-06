#!/usr/bin/env tclsh

source drrt_common.tcl

puts "Heyy this worked!"

# What does this script need to accomplish in order to truly function as the 
# DRRT Scheduler?
# - Create correct folder structure. This should likely be in a specified 
#   location
# - Reference a match schedule file and index the array based on it
# - alliance assembler function
proc makeShipDictFromFile {filename} {
    set filehandle [open "$filename" r]
    set filecontents [read $filehandle]
    close $filehandle

    switch -glob -nocase -- "$filename" {
        *.json {
            makeShipDictFromJson 
    }
}

proc makeShipDictFromJson {varname} {
    upvar $varname jsonstring

    set shipdict {[::json::json2dict $jsonstring]}

    return $shipdict
}

proc saveFleetToFile {filename fleet} {
}

proc makeShipsIntoFleet {ships} {

}

proc removeBlockDataFromShip {ship_dict} {
}
