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

```
Options:
  --from FROM, -f FROM   Source property
  --rule RULE, -r RULE   Rule name
  --var VAR, -v VAR      Variable names, wildcard * support
  --def DEF, -d DEF      Rule definition (json) to store / load the rule
  --to TO, -t TO         Target property
  --comments COMMENTS, -m COMMENTS
                        Overrule default commit message / version note
  --fromversion FROMVERSION
                        Versionnumber or Latest, Production Staging
  --fromjson FROMJSON    Use JSON export as source property (instead of property manager)
  --toversion TOVERSION
                        Versionnumber or Latest, Production Staging
  --tojson TOJSON        Use JSON export as target property (instead of property manager)
  --edgerc EDGERC        [default: ~/.edgerc] [default: ~/.edgerc]
  --section SECTION      [default: default] [default: default]
  --account ACCOUNT      Accountswitchkey (partners and Akamai only)
  --toedgerc TOEDGERC    [default EDGERC]
  --tosection TOSECTION
                        [default SECTION]
  --toaccount TOACCOUNT
                        [default ACCOUNT]
  --log LOG              Log file
  --silent, -s           Quiet mode
  --dryrun               additional validation and supress actual update
  --backup BACKUP        Backup of the to-property
  --help, -h             display this help and exit
  --version              display version and exit
```
## Installation
### Using akamai CLI
    $ akamai install https://github.com/ericdebeij/rulecopy.git

### Without akamai CLI
Download 
[latest release binary](https://github.com/ericdebeij/rulecopy/releases)
for your system, or by cloning this repository and compiling it yourself.

## TODO
- --values - copy variable values from definition
- subcommands for quick usage:
  - akamai rcp COPY from_prop:version rule vars (into pasteboard)
  - akamai rcp PASTE to_prop:version (from pasteboard)
- --new BASEVERSION - create a new property version instead of updating latest
