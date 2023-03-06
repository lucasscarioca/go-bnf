package bnf

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidation(t *testing.T) {
	content, err := os.ReadFile("../bnf.json")
	if err != nil {
		t.Errorf("Got %q", err)
		return
	}

	var parsedGrammar Grammar
	err = json.Unmarshal(content, &parsedGrammar)
	if err != nil {
		t.Errorf("Got %q", err)
		return
	}

	err = parsedGrammar.ValidateGrammar()

	if err != nil {
		t.Errorf("Got %q", err)
		return
	}

	validInputs := []string{"+12.", "-12.3748E38", "+347.789E-90", "18.", "-3878.7918E+400"}

	for _, input := range validInputs {
		got := parsedGrammar.ValidateInput(input)
		want := true
		if got != want {
			t.Errorf("For input: %s Got %t, want %t", input, got, want)
			return
		}
	}

	invalidInput := []string{"asdasdas", "", "12", "123123123", "1213.12312.3123.12.31.2312.3...."}

	for _, input := range invalidInput {
		got := parsedGrammar.ValidateInput(input)
		want := false
		if got != want {
			t.Errorf("For invalid input: %s Got %t, want %t", input, got, want)
			return
		}
	}
}
