package bnf

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type Grammar struct {
	NonTerminals  []string
	Rules         []string
	RulesMap      map[string][]string
	Terminals     []string
	InitialSymbol string
}

func (g *Grammar) ValidateGrammar() error {
	for _, nonTerminal := range g.NonTerminals {
		if !strings.HasPrefix(nonTerminal, "<") && !strings.HasSuffix(nonTerminal, ">") {
			return errors.New("invalid nonTerminal")
		}
	}
	for _, rule := range g.Rules {
		ruleset := strings.Split(rule, " ")

		ruleSymbol, ruleAssignment, ruleExpression := ruleset[0], ruleset[1], ruleset[2:]
		rulesMap := make(map[string][]string)
		rulesMap[ruleSymbol] = ruleExpression

		if g.RulesMap == nil {
			g.RulesMap = rulesMap
		}
		for k, v := range rulesMap {
			g.RulesMap[k] = v
		}

		if !slices.Contains(g.NonTerminals, ruleSymbol) {
			return errors.New("unmapped nonTerminal in rules")
		}
		if ruleAssignment != "::=" {
			return errors.New("missing rule assignment")
		}
		for i := 0; i < len(ruleExpression); i++ {
			if ruleExpression[i] != "|" && !slices.Contains(g.NonTerminals, ruleExpression[i]) && !slices.Contains(g.Terminals, ruleExpression[i]) {
				return errors.New("invalid rule")
			}
		}
	}
	return nil
}

func (g *Grammar) ValidateInput(input string) bool {
	inputSlice := strings.Split(strings.Trim(input, " "), "")
	return g.validateRule(inputSlice, g.InitialSymbol, g.RulesMap[g.InitialSymbol])
}

// TODO separate validation for rule: if input does not exist in rule expression
func (g *Grammar) validateRule(inputSlice []string, rule string, expressions []string) bool {
	var valid bool
	if slices.Contains(expressions, "|") {
		concatExpression := strings.Join(expressions, "<*separator*>")
		splitExpressions := strings.Split(concatExpression, "|")
		for _, exp := range splitExpressions {
			exp = strings.TrimPrefix(exp, "<*separator*>")
			exp = strings.TrimSuffix(exp, "<*separator*>")
			possibleExpressions := strings.Split(strings.Trim(exp, " "), "<*separator*>")
			res := g.validateRule(inputSlice, rule, possibleExpressions)
			if res && rule == g.InitialSymbol {
				return true
			}
		}
		if rule == g.InitialSymbol {
			valid = false
		}
		return valid
	}

	if len(inputSlice) == 0 && len(expressions) > 0 {
		return false
	}

	if !slices.Contains(g.NonTerminals, inputSlice[0]) && !slices.Contains(g.Terminals, inputSlice[0]) {
		return false
	}

	if len(expressions) == 0 {
		return true
	}

	if inputSlice[0] == expressions[0] {
		g.validateRule(inputSlice[1:], rule, expressions[1:])
	}

	if expressions[0] == "null" {
		g.validateRule(inputSlice, rule, expressions[1:])
	}

	if slices.Contains(g.NonTerminals, expressions[0]) {
		if rule == expressions[0] {
			if slices.Contains(g.NonTerminals, expressions[1]) && slices.Contains(g.RulesMap[expressions[1]], inputSlice[0]) {
				g.validateRule(inputSlice[1:], rule, expressions)
			}

			g.validateRule(inputSlice, rule, expressions[1:])
		}

		g.validateRule(inputSlice, expressions[0], g.RulesMap[expressions[0]])
	}
	return valid
}
