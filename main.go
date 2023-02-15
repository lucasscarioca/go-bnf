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
	// for _, c := range inputSlice {
	// 	if !g.validateRule(c, g.InitialSymbol, g.RulesMap[g.InitialSymbol]) {
	// 		return false
	// 	}
	// }
	return g.validateFormat(inputSlice, g.InitialSymbol, g.RulesMap[g.InitialSymbol])
}

// func (g *Grammar) validateRule(c string, rule string, expression []string) bool {

// 	if slices.Contains(expression, "|") {
// 		concatExpression := strings.Join(expression, "<*separator*>")
// 		splitExpressions := strings.Split(concatExpression, "|")
// 		for _, e := range splitExpressions {
// 			possibleExpressions := strings.Split(e, "<*separator*>")
// 			res := g.validateRule(c, rule, possibleExpressions)
// 			if res {
// 				return true
// 			}
// 		}
// 	}

// 	for _, expressionSymbol := range expression {
// 		if expressionSymbol != rule {
// 			if slices.Contains(g.Terminals, expressionSymbol) {
// 				if c == expressionSymbol {
// 					return true
// 				}
// 			}

// 			exp, isNonTerminal := g.RulesMap[expressionSymbol]
// 			if isNonTerminal {
// 				res := g.validateRule(c, expressionSymbol, exp)
// 				if res {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

func (g *Grammar) validateFormat(input []string, rule string, expression []string) bool {

	if slices.Contains(expression, "|") {
		concatExpression := strings.Join(expression, "<*separator*>")
		splitExpressions := strings.Split(concatExpression, "|")
		var isValidVariant bool
		for _, exp := range splitExpressions {
			possibleExpressions := strings.Split(exp, "<*separator*>")
			isValidVariant = g.validateFormat(input, rule, possibleExpressions)
		}
		return isValidVariant
	}
	for _, symbol := range expression {
		if symbol == "null" {
			continue
		}

		if slices.Contains(g.Terminals, symbol) && !slices.Contains(input, symbol) {
			if rule == g.InitialSymbol {
				return false
			}
			continue
		}

		if slices.Contains(g.NonTerminals, symbol) {
			if slices.Contains(g.RulesMap[symbol], symbol) {
				continue
			} else {
				g.validateFormat(input, symbol, g.RulesMap[symbol])
			}
		}
	}
	return true
}

func main() {
	content, err := os.ReadFile("./bnf.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	fmt.Println("Parsing Grammar file...")

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

	for {
		fmt.Print("\n(type :quit to close)\nExpression to validate: ")

		var input string
		fmt.Scanln(&input)

		if input == ":quit" {
			fmt.Println("Closing program...")
			break
		}

		res := parsedGrammar.ValidateInput(input)
		if !res {
			fmt.Println("\nInvalid expression")
		} else {
			fmt.Println("\nValid expression")
		}
	}

}
