# akamai-rulecopy
**DRAFT see draft note at the end**

Rulecopy is a sample utility / package for Akamai configurations to copy a rule and variable from one property manager configuration to another one.

Sample usage:

    $ akamai-rulecopy -f from_config -r some_rule -v some_var -t to_config

Some features:
* the utility will search in the source property for the rule with the given name and copies the content together with the variable definitions
* the utility will search in the target property for the rule with the given name and overrides this with the rule as found in the source. Variables will be merged, if the variable exists the initial-value from the target property will be re-used
* the rule definitions can be stored in a configuration file or read from a configuration file
* the variable selection does support the wildcard character *

## Installation
### Using akamai CLI
    $ akamai install https://github.com/ericdebeij/rulecopy.git

### Without akamai CLI
Download 
[latest release binary](https://github.com/ericdebeij/rulecopy/releases)
for your system, or by cloning this repository and compiling it yourself.

## STILL IN DRAFT
TODO:
- //Done --json/-j => --def/-d
- //rule not in target => place at the end of the Default Section
- --dryrun - perform a dryrun
- --values - copy variable values from definition
- subcommands for quick usage:
  - //DONE akamai rcp
  - akamai rcp COPY from_prop:version rule vars (into pasteboard)
  - akamai rcp PASTE to_prop:version (from pasteboard)
- --new BASEVERSION - create a new property version instead of updating latest
- --comment / -m - alternative note instead of pregenerated message
