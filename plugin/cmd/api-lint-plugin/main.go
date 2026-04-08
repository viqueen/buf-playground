package main

import "buf.build/go/bufplugin/check"

func main() {
	check.Main(&check.Spec{
		Rules: []*check.RuleSpec{},
	})
}
