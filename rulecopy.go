package rulecopy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ryanuber/go-glob"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type RuleCopyProperty struct {
	Property string
	Version  int
	Edgerc   string
	Section  string
	Account  string
	Json     string
}
type RuleCopyParam struct {
	From RuleCopyProperty
	To   RuleCopyProperty
	Rule string
	Var  string
	Json string
}

func (p *RuleCopyParam) Validate() (err error) {
	if p.From.Section == "" {
		p.From.Section = "default"
	}
	if p.To.Edgerc == "" {
		p.To.Edgerc = p.From.Edgerc
	}
	if p.To.Section == "" {
		p.To.Section = p.From.Section
	}
	if p.To.Account == "" {
		p.To.Account = p.From.Account
	}
	if p.From.Property == "" && p.From.Json == "" && p.Json == "" {
		err = fmt.Errorf("no source defined, either json, property or json-property needs to be provided")
	}
	if p.To.Property == "" && p.To.Json == "" && p.Json == "" {
		err = fmt.Errorf("no target defined, either json, property or json-property needs to be provided")
	}
	if p.From == p.To {
		err = fmt.Errorf("source and target required and should not be identical")
	}
	return
}
func (p RuleCopyParam) SameEdgerc() bool {
	p.Validate()
	return p.From.Edgerc == p.To.Edgerc && p.From.Section == p.To.Section && p.From.Account == p.To.Account
}

const LATEST int = 0
const PRODUCTION int = -1
const STAGING int = -2

func VersionConv(s string) (i int, err error) {
	switch strings.ToUpper(s) {
	case "PRODUCTION":
		i = PRODUCTION
	case "STAGING":
		i = STAGING
	case "LATEST":
		i = LATEST
	case "":
		i = LATEST
	default:
		i, err = strconv.Atoi(s)
	}
	return
}

// Rules contains Rule object
type CopyRule struct {
	Name      string              `json:"name"`
	Comments  string              `json:"comment"`
	Rules     []papi.Rules        `json:"rules,omitempty"`
	Variables []papi.RuleVariable `json:"variables,omitempty"`
}

func FetchRules(p papi.PAPI, name string, version int) (treeResponse *papi.GetRuleTreeResponse, err error) {
	log.Printf("searching property %s:%d", name, version)
	ps := papi.SearchRequest{Key: papi.SearchKeyPropertyName, Value: name}
	k, err := p.SearchProperties(context.Background(), ps)

	if err != nil {
		return
	}
	if len(k.Versions.Items) == 0 {
		err = fmt.Errorf("property %s not found", name)
		return
	}
	u := k.Versions.Items[0]
	if version <= 0 {
		for _, v := range k.Versions.Items {
			if v.PropertyVersion > u.PropertyVersion {
				u = v
			}
			if version == PRODUCTION && v.ProductionStatus == "ACTIVE" {
				version = v.PropertyVersion
			}
			if version == STAGING && v.ProductionStatus == "ACTIVE" {
				version = v.PropertyVersion
			}
		}
		if version == LATEST {
			version = u.PropertyVersion
		}
	}

	treeRequest := papi.GetRuleTreeRequest{
		PropertyID:      u.PropertyID,
		PropertyVersion: version,
		ContractID:      u.ContractID,
		GroupID:         u.GroupID,
	}

	treeResponse, err = p.GetRuleTree(context.Background(), treeRequest)
	if err != nil {
		return
	}
	log.Printf("property %s:%d loaded", name, version)
	return
}

func StoreRules(p papi.PAPI, name string, g *papi.GetRuleTreeResponse) (res *papi.UpdateRulesResponse, err error) {
	ptr := papi.UpdateRulesRequest{
		PropertyID:      g.PropertyID,
		PropertyVersion: g.PropertyVersion,
		ContractID:      g.ContractID,
		DryRun:          false,
		GroupID:         g.GroupID,
		ValidateMode:    "fast",
		ValidateRules:   false,
		Rules: papi.RulesUpdate{
			Comments: g.Comments,
			Rules:    g.Rules,
		},
	}
	res, err = p.UpdateRuleTree(context.Background(), ptr)
	return
}

func BuildCopyRule(rulename string, varname string, rules *papi.Rules) *CopyRule {
	c := &CopyRule{Name: rulename}

	log.Printf("searching for rule %s\n", rulename)
	walkrule(0, rules, func(r *papi.Rules) (stop bool) {
		if c.Name == r.Name {
			c.Rules = append(c.Rules, *r)
			log.Printf("-found source rule %s\n", c.Name)
			stop = true
		}
		return
	})

	log.Printf("searching for variables %s\n", varname)
	for _, v := range rules.Variables {
		if glob.Glob("PMUSER_"+varname, v.Name) {
			c.Variables = append(c.Variables, v)
			log.Printf("-found source variable %s\n", v.Name)
		}
	}
	return c
}

func replaceRule(c *CopyRule, r *papi.Rules) (found bool) {
	for i := range r.Children {
		for j := range c.Rules {
			if r.Children[i].Name == c.Rules[j].Name {
				r.Children[i] = c.Rules[j]
				log.Printf("rule %s found and replaced", c.Rules[j].Name)
				found = true
			}
		}
		if !found {
			found = replaceRule(c, &r.Children[i])
		}
	}
	return
}

func MergeCopyRule(c *CopyRule, torules *papi.Rules) (err error) {
	foundrule := replaceRule(c, torules)
	if !foundrule {
		err = fmt.Errorf("rule %s not found in target", c.Name)
	}

	for _, tv := range c.Variables {
		found := false
		for si, sv := range torules.Variables {
			if !found && sv.Name == tv.Name {
				log.Printf("variable %s found and synced", sv.Name)
				torules.Variables[si].Description = tv.Description
				torules.Variables[si].Hidden = tv.Hidden
				torules.Variables[si].Sensitive = tv.Sensitive
				found = true
			}
		}
		if !found {
			log.Printf("variable %s added\n", tv.Name)
			torules.Variables = append(torules.Variables, tv)
		}
	}
	return
}

func papiClient(param RuleCopyProperty) (p papi.PAPI, err error) {
	var e *edgegrid.Config

	e, err = edgegrid.New(edgegrid.WithFile(param.Edgerc), edgegrid.WithSection(param.Section))
	if err != nil {
		return
	}

	s, err := session.New(session.WithSigner(e))
	if err != nil {
		return
	}
	p = papi.Client(s)
	return
}

func Run(param RuleCopyParam) (err error) {
	var papiFrom, papiTo papi.PAPI

	err = param.Validate()
	if err != nil {
		return
	}
	papiFrom, err = papiClient(param.From)
	if err != nil {
		return
	}

	if param.SameEdgerc() {
		papiTo = papiFrom
	} else {
		papiTo, err = papiClient(param.To)
		if err != nil {
			return
		}
	}

	// Source
	var fromPropertyRules *papi.GetRuleTreeResponse
	if param.From.Property != "" {
		fromPropertyRules, err = FetchRules(papiFrom, param.From.Property, param.From.Version)
		if err != nil {
			return
		}
		if param.From.Json != "" {
			err = exportJson(param.From.Json, fromPropertyRules)
			if err != nil {
				return
			}
		}
	} else {
		if param.From.Json != "" {
			var ff papi.GetRuleTreeResponse
			err = importJson(param.From.Json, &fromPropertyRules)
			if err != nil {
				return
			}
			fromPropertyRules = &ff
			log.Printf("source property loaded from %s", param.From.Json)
		}
	}

	var copyRule *CopyRule

	if fromPropertyRules == nil && param.Json != "" {
		var cp CopyRule
		err = importJson(param.Json, &cp)
		if err != nil {
			return
		}
		copyRule = &cp
		log.Printf("%s loaded from %s", copyRule.Comments, param.Json)
	}

	if param.Rule != "" || param.Var != "" {
		if fromPropertyRules != nil {
			copyRule = BuildCopyRule(param.Rule, param.Var, &fromPropertyRules.Rules)
			copyRule.Comments = fmt.Sprintf("rule %s, vars %s (%s:%d)",
				param.Rule, param.Var,
				param.From.Property, fromPropertyRules.PropertyVersion)
			if param.Json != "" {
				err = exportJson(param.Json, copyRule)
				if err != nil {
					return
				}
				log.Printf("%s stored in %s", copyRule.Comments, param.Json)
			}
		}
	}

	var toPropertyRules *papi.GetRuleTreeResponse
	if param.To.Property != "" {
		// Target
		toPropertyRules, err = FetchRules(papiTo, param.To.Property, param.To.Version)
		if err != nil {
			return
		}
	}
	//exportJson("debug_unchanged.json", toPropertyRules)

	if toPropertyRules != nil {
		if copyRule != nil {
			err = MergeCopyRule(copyRule, &toPropertyRules.Rules)
			if err != nil {
				return
			}
			toPropertyRules.Comments = strings.Trim(fmt.Sprintf("%s\n%s", toPropertyRules.Comments, copyRule.Comments), "\n ")
		} else {
			// Copy the entire ruletree as no subrule name nor variables are provided
			toPropertyRules.Rules = fromPropertyRules.Rules
			toPropertyRules.Comments = strings.Trim(fmt.Sprintf("%s\nContent copied from %s:%d",
				fromPropertyRules.Comments, param.From.Property, param.From.Version), "\n ")
		}

		if param.To.Json != "" {
			err = exportJson(param.To.Json, toPropertyRules)
			if err != nil {
				return
			}
			log.Printf("property %s:%d exported to %s", param.To.Property, toPropertyRules.PropertyVersion, param.To.Json)
		}

		if param.To.Property != "" {
			_, err = StoreRules(papiTo, param.To.Property, toPropertyRules)
			if err != nil {
				return
			}
			log.Printf("property %s:%d stored", param.To.Property, toPropertyRules.PropertyVersion)
		}
	}
	return
}

func exportJson(filename string, content interface{}) error {
	file, err := json.MarshalIndent(content, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, file, 0644)
	return err
}

func importJson(filename string, content interface{}) (err error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, content)
	return
}

func walkrule(d int, r *papi.Rules, f func(r *papi.Rules) bool) {
	if f(r) {
		return
	}
	// Note: don't use for _, c as that will create a copy and we need to work by-reference onlt
	for i := range r.Children {
		walkrule(d+1, &r.Children[i], f)
	}
}
