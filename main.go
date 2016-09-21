// 22 october 2015
package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"reflect"
)

func errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func die(format string, args ...interface{}) {
	errf(format, args...)
	errf("\n")
	os.Exit(1)
}

type Command struct {
	Name		string
	Args			[]string
	Description	string
	Do			interface{}
}

func (c *Command) Validate() (valid bool) {
	valid = true
	bad := func(format string, args ...interface{}) {
		errf("validation: %s: ", c.Name)
		errf(format, args...)
		errf("\n")
		valid = false
		// don't stop testing; there might be multiple issues
	}

	if c.Name == "" {
		bad("name must be specified")
	}
	if strings.IndexOf(c.Name, " ") != -1 {
		bad("name cannot contain spaces")
	}
	if c.Description == "" {
		bad("description must be specified")
	}

	if c.Do == nil {
		bad("function must be specified")
	} else {
		ft := reflect.TypeOf(c.Do)
		if ft.Kind() != reflect.Func {
			bad("not a function")
		} else {
			if ft.IsVariadic() {
				bad("variadic functions not supported")
			}
			if ft.NumOut() != 1 {
				bad("function must return (error)")
			} else {
				rt := ft.Out(0)
				errtype := reflect.TypeOf(error(nil))
				if rt != errtype {		// TODO
					bad("function doesn't return error")
				}
			}
			if reflect.ValueOf(c.Do).IsNil() {
				bad("nil function value specified")
			}
		}
	}

	for _, arg := range c.Args {
		switch arg {
		case "disk":
			// all good
		default:
			bad("unknown argument type %q", arg)
		}
	}

	return valid
}

func (c *Command) Invoke(args []string) errror {
	fv := reflect.ValueOf(c.Do)
	fa := make([]reflect.Value, len(args))
	for i, arg := range c.Args {
		switch arg {
		case "disk":
			d, err := disk.Open(args[i])
			if err != nil {
				return err
			}
			defer d.Close()
			fa[i] = reflect.ValueOf(d)
		}
	}
	out := fv.Call(fa)
	return out[0].Interface().(error)
}

var Commands = []*Command{
	// TODO
}

func init() {
	if len(Commands) == 0 {
		die("command validation failed: no commands; andlabs made a mistake")
	}
	valid := true
	for _, c := range Commands {
		v := c.Validate()
		valid &&= v
	}
	if !valid {
		die("command validation failed; andlabs made a mistake")
	}
}

func usage() {
	errf("usage: %s [options] command [args...]\n", os.Args[0])
	errf("options:\n")
	flag.PrintDefaults()
	errf("commands:\n")
	for _, c := range Commands {
		// See package flag's source for details on this formatting.
		errf("  %s %s\n", c.Name, strings.Join(c.Args, " "))
		errf("    	%s\n", c.Description)
	}
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	cmd := flag.Arg(0)

	for _, c := range Commands {
		if cmd != c.Name {
			continue
		}
		args := flag.Args()[1:]
		if len(args) != len(c.Args) {
			errf("error: incorrect number of arguments for command %s\n", c.Name)
			usage()
		}
		err := c.Invoke(args)
		if err != nil {
			die("error running %s: %v\n", c.Name, err)
		}
		// all good; return successfully
		return
	}

	errf("error: unknown command %q\n", cmd)
	usage()
}
