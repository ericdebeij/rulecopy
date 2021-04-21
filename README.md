# akamai-rulecopy
Rulecopy is a sample utility / package for Akamai configurations to copy a rule and variable from one property manager configuration to another one.

Sample usage to copy the rule cors and the related variables CORS_*:

    $ akamai-rulecopy -f fromprop -r cors -v 'CORS_*' -t toprop

As a best practice the utility can also be used as an akamai cli:

    $ akamai rcp -f fromprop -r cors -v 'CORS_*' -t toprop

Rules+variables can also be stored in a definition file and use them as text objects:

    $ akamai rcp -d cors.json -t toprop

Some features:
* the utility will search in the source property for the rule with the given name and copies the content together with the variable definitions
* the utility will search in the target property for the rule with the given name and overrides this with the rule as found in the source, if the rule is not found it will be added at the end. Variables will be merged, if the variable exists the initial-value from the target property will be re-used
* the rule definitions can be stored in a configuration file or read from a configuration file
* the variable selection does support the wildcard character *

```
Command line paramters:
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
  --toversion TOVERSION  Versionnumber or Latest, Production Staging
  --tojson TOJSON        Use JSON export as target property (instead of property manager)
  --edgerc EDGERC        [default: ~/.edgerc] [default: ~/.edgerc]
  --section SECTION      [default: default] [default: default]
  --account ACCOUNT      Accountswitchkey (partners and Akamai only)
  --toedgerc TOEDGERC    [default EDGERC]
  --tosection TOSECTION  [default SECTION]
  --toaccount TOACCOUNT  [default ACCOUNT]
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
- --values - copy variable values from definition; currently variables values are not overwritten (by design) as they can contain configurable values)
- --new BASEVERSION - create a new property version instead of updating latest

## Not implemented
- Considered to have the utility work like the pasteboard (pbcopy, pbpaste). Didn't do that but you can simply mimik that by using a definition file.
- Considered exporting a json diff of the changes that will be / are made. Didn't do that as you can simply compare the backup and target export.
