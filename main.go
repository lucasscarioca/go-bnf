package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

/*
valid examples:
+12.
-12.3748E38
+347.789E-90
18.
-3878.7918E+400
*/

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
	for _, c := range inputSlice {
		if !g.validateRule(c, g.InitialSymbol, g.RulesMap[g.InitialSymbol]) {
			return false
		}
	}
	return true
}

func (g *Grammar) validateRule(c string, rule string, expression []string) bool {
	for i, expressionSymbol := range expression {
		exp, isNonTerminal := g.RulesMap[expressionSymbol]
		if isNonTerminal {
			g.validateRule(c, expressionSymbol, exp)
		}
		if expressionSymbol == "|" {
			return g.validateRule(c, rule, expression[:i]) || g.validateRule(c, rule, expression[i:])
		}

	}
	return true
}

func main() {
	content, err := os.ReadFile("./bnf.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var parsedGrammar Grammar
	err = json.Unmarshal(content, &parsedGrammar)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	err = parsedGrammar.ValidateGrammar()
	if err != nil {
		log.Println("[INVALID GRAMMAR] reason: ", err)
		return
	}
	fmt.Println("Grammar Accepted")

	fmt.Print("Input: ")

	var input string
	fmt.Scanln(&input)

	res := parsedGrammar.ValidateInput(input)
	if !res {
		fmt.Println("Invalid expression")
	} else {
		fmt.Println("Valid expression")
	}
}
