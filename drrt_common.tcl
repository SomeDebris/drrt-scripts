# Common commands that we all should use

package require json
package require json::write

# Converts a dict object to a json object
# Taken from the tcl docs
proc dict2json {dict_to_encode} {
    ::json::write object {*}[dict map {k v} $dict_to_encode {
        set v [::json::write string $v]
    }]
}

proc waitYN {prompt_string} {
    while {true} {
        puts "$prompt_string \[y/n\]: "

        gets stdin response

        switch "$response" {
            [yY] { return true }
            [nN] { return false }
        }
    }
}

