package main

import (
	"buf.build/go/bufplugin/check"
	"github.com/viqueen/buf-playground/plugin/internal/api"
)

func main() {
	check.Main(&check.Spec{
		Rules: []*check.RuleSpec{
			api.FileNameConventionRule(),
			api.UpdateRequestFieldMaskRule(),
		},
	})
}
