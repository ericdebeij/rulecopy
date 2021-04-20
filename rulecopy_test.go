package rulecopy

import (
	"log"
	"testing"
)

func TestRuleCopy_dryrun(t *testing.T) {
	log.Print("test happy")
	param := RuleCopyParam{
		Rule:   "oidc",
		Var:    "OIDC_*",
		Def:    "rule.json",
		Dryrun: true,
		From: RuleCopyProperty{
			Property: "debeij.lunacooking.com",
			Version:  67,
			Edgerc:   "/workspaces/go/old/ps.edgerc",
			Section:  "default",
			Account:  "",
			Json:     "debug_from.json",
		},
		To: RuleCopyProperty{
			Property: "dsa2_hdebeij",
			Version:  0,
			Edgerc:   "",
			Section:  "",
			Account:  "",
			Backup:   "",
			Json:     "debug_dry.json",
		},
	}
	err := Run(param)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}
func TestRuleCopy_happy(t *testing.T) {
	log.Print("test happy")
	param := RuleCopyParam{
		Rule: "oidc",
		Var:  "OIDC_*",
		Def:  "rule.json",
		From: RuleCopyProperty{
			Property: "debeij.lunacooking.com",
			Version:  67,
			Edgerc:   "/workspaces/go/old/ps.edgerc",
			Section:  "default",
			Account:  "",
			Json:     "debug_from.json",
		},
		To: RuleCopyProperty{
			Property: "dsa2_hdebeij",
			Version:  0,
			Edgerc:   "",
			Section:  "",
			Account:  "",
			Backup:   "debug_backup.json",
			Json:     "debug_to.json",
		},
	}
	err := Run(param)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}

func TestRuleCopy_pbcopy(t *testing.T) {
	log.Print("test pbcopy")
	param := RuleCopyParam{
		Rule:     "authenticated",
		Var:      "AUTH*",
		Def:      "authenticated.json",
		Comments: "auth-1",
		From: RuleCopyProperty{
			Property: "hdebeij4.ps-akamai.nl",
			Version:  126,
			Edgerc:   "/workspaces/go/old/ps.edgerc",
			Section:  "default",
			Account:  "",
			Json:     "debug_from.json",
		},
		To: RuleCopyProperty{
			Property: "",
			Version:  0,
			Edgerc:   "",
			Section:  "",
			Account:  "",
			Json:     "",
		},
	}
	err := Run(param)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}

func TestRuleCopy_pbpaste(t *testing.T) {
	log.Print("test pbcopy")
	param := RuleCopyParam{
		Rule:     "authenticated",
		Var:      "AUTH*",
		Def:      "authenticated.json",
		Comments: "Auth-2",
		From: RuleCopyProperty{
			Property: "",
			Version:  0,
			Edgerc:   "/workspaces/go/old/ps.edgerc",
			Section:  "default",
			Account:  "",
			Json:     "",
		},
		To: RuleCopyProperty{
			Property: "dsa2_hdebeij",
			Version:  0,
			Edgerc:   "",
			Section:  "",
			Account:  "",
			Json:     "",
		},
	}
	err := Run(param)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}

func TestRuleCopy_rollback(t *testing.T) {
	log.Print("test rollback")
	param := RuleCopyParam{
		From: RuleCopyProperty{
			Property: "dsa2_hdebeij",
			Version:  2,
			Edgerc:   "/workspaces/go/old/ps.edgerc",
			Json:     "debug_from.json",
		},
		To: RuleCopyProperty{
			Property: "dsa2_hdebeij",
			Json:     "debug_to.json",
		},
	}
	err := Run(param)
	if err != nil {
		t.Errorf("Error %s", err)
	}
}
