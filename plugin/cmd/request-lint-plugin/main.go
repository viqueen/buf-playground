package main

import (
	"buf.build/go/bufplugin/check"
	"github.com/viqueen/buf-playground/plugin/internal/request"
)

func main() {
	check.Main(&check.Spec{
		Rules: []*check.RuleSpec{
			request.RepeatedFieldValidationRule(),
		},
	})
}
