package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bnf "github.com/lucasscarioca/go-bnf/internal"
)

/*
valid examples:
+12.
-12.3748E38
+347.789E-90
18.
-3878.7918E+400
*/

func main() {
	content, err := os.ReadFile("./bnf.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	fmt.Println("Parsing Grammar file...")

	var parsedGrammar bnf.Grammar
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
