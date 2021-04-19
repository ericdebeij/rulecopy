package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/ericdebeij/rulecopy"
)

var (
	VERSION = "0.1.0"
)

func setlogfile(filename string) (file *os.File, err error) {
	if filename == "" {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		log.SetOutput(file)
	}
	return
}

type args struct {
	From        string `arg:"-f" help:"Source property" placholder:"FROMPROP"`
	Rule        string `arg:"-r" help:"Rule name"`
	Var         string `arg:"-v" help:"Variable names, wildcard * support"`
	Def         string `arg:"-d" help:"Rule definition (json) to store / load the rule"`
	To          string `arg:"-t" help:"Target property"`
	Fromversion string `help:"Versionnumber or Latest, Production Staging"`
	Fromjson    string `help:"Use JSON export as source property (instead of property manager)"`
	Toversion   string `help:"Versionnumber or Latest, Production Staging"`
	Tojson      string `help:"Use JSON export as target property (instead of property manager)"`
	Edgerc      string `help:"[default: ~/.edgerc]" arg:"env" default:"~/.edgerc"`
	Section     string `help:"[default: default]" default:"default"`
	Account     string `help:"Accountswitchkey (partners and Akamai only)"`
	Toedgerc    string `help:"[default EDGERC]"`
	Tosection   string `help:"[default SECTION]"`
	Toaccount   string `help:"[default ACCOUNT]"`
	Log         string `help:"Log file"`
	Silent      bool   `arg:"-s" help:"Quiet mode"`
	Version     string `help:"Version identification"`
}

func (args) Description() string {
	return "copyrule copies a rule (+variables) in property manager to another configuration (version)"
}

func main() {
	//Run(os.Args)
	var err error

	var args args
	argres := arg.MustParse(&args)

	fromVersion, err := rulecopy.VersionConv(args.Fromversion)
	if err != nil {
		argres.Fail(fmt.Sprintf("Fromversion not valid: %s", err))
	}
	toVersion, err := rulecopy.VersionConv(args.Toversion)
	if err != nil {
		argres.Fail(fmt.Sprintf("Toversion not valid: %s", err))
	}

	if args.Log != "" || args.Silent {
		file, err := setlogfile(args.Log)
		if err != nil {
			argres.Fail(fmt.Sprintf("Logfile: %s", err))
		}
		if file != nil {
			defer file.Close()
		}
	}

	param := rulecopy.RuleCopyParam{
		Rule: args.Rule,
		Var:  args.Var,
		Def:  args.Def,
		From: rulecopy.RuleCopyProperty{
			Property: args.From,
			Version:  fromVersion,
			Edgerc:   args.Edgerc,
			Section:  args.Section,
			Account:  args.Account,
			Json:     args.Fromjson,
		},
		To: rulecopy.RuleCopyProperty{
			Property: args.To,
			Version:  toVersion,
			Edgerc:   args.Toedgerc,
			Section:  args.Tosection,
			Account:  args.Toaccount,
			Json:     args.Fromjson,
		},
	}
	err = param.Validate()
	if err != nil {
		argres.Fail(fmt.Sprint(err))
	}

	err = rulecopy.Run(param)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
