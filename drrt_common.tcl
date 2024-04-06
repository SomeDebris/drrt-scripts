# Common commands that we all should use

package require json
package require json::write

# Converts a dict object to a json object
# Taken from the tcl docs
# proc dict2json {dict_to_encode} {
#     set accumulator {}
#     ::json::write object {*}[dict map {k v} $dict_to_encode {
#         set v [::json::write string $v]
#     }]
# }

# Lets try to understand this person's code!
proc dict2json {dict_to_encode spec {indent false}} {
    ::json::write::indented $indent

    set accumulator [dict create]

    dict for {field type_info} $spec {
        if {![dict exists $dict_to_encode $field]} {continue}

        lassign $type_info type meta

        set value [dict get $dict_to_encode $field]
        switch $type {
            object {
                set value [dict2json $value $meta $indent]
            }
            array {
                set value [list2jsonarray $value {*}$meta]
            }
            string {
                set value [::json::write string $value]
            }
            bare {}
            default {
                return -code error "Type \"$type\" is not known!"
            }
        }
        dict set accumulator $field $value
    }
    return [::json::write object {*}$accumulator]
}

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

