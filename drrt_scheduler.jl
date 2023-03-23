# The whole reason that I'm doing this is just so I can learn this language
# a little better. It might be fun!

import JSON

global const RED_ALLIANCE_COLORS    = [ 0xbaa01e, 0x681818, 0x000000 ]
global const BLUE_ALLIANCE_COLORS   = [ 0xbaa01e, 0x681818, 0x000000 ]

global const QUALS_DIR = ""


function ship_at(index::Int)
    # return the ship file at this index
end

function ship_at(index::String)
    # return the ship file at this index.
    # surrogate ships will be strings, as they have an asterisk.
    # straight up call the ship_at(index::Int) function once int has been parsed

    regex_noast = r"\*"
    return ship_at( parse(Int, replace(index, regex_noast => "")) )
end


function ship_dict(ship_index::String)
    open(ship_index) do file

    end
end


function parse_alliance(partners, Name::String)
    # creates an alliance with a certain name.
    # doesn't color alliance
end

function parse_alliance(partners, is_red::Bool)
    # creates an alliance that is either red or blue
    # colors alliance properly
end


# Goal: Split participants into Red and Blue alliances.
# calls alliance parser function on red and blue halves
function parse_match(participants)
    # participant_count = size(participants, 1)
end


function main(args)
end
     
