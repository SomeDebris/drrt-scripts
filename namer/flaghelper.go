package namer

import (
	"flag"
	"os"
)
// Source - https://stackoverflow.com/a

// Posted by Brandon, modified by community. See post 'Timeline' for change history

// Retrieved 2025-12-22, License - CC BY-SA 4.0



// ParseFlags works like flag.Parse(), except positional
// args and flag args can be specified in any order.
func ParseFlags() error {
    return ParseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet works like flagset.Parse(), except positional
// args and flag args can be specified in any order.
func ParseFlagSet(flagset *flag.FlagSet, args []string) error {
    var positionalArgs []string
    for {
        if err := flagset.Parse(args); err != nil {
            return err
        }
        // Consume all the flags that were parsed as flags.
        args = args[len(args)-flagset.NArg():]
        if len(args) == 0 {
            break
        }
        // There's at least one flag remaining and it must be a positional arg since
        // we consumed all args that were parsed as flags. Consume just the first
        // one, and retry parsing, since subsequent args may be flags.
        positionalArgs = append(positionalArgs, args[0])
        args = args[1:]
    }
    // Parse just the positional args so that flagset.Args()/flagset.NArgs()
    // return the expected value.
    // Note: This should never return an error.
    return flagset.Parse(positionalArgs)
}

