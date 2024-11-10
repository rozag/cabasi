package dice

import "testing"

type fixedRNG uint

func (f fixedRNG) UintN(uint) uint { return uint(f) }

func TestDiceRollRange(t *testing.T) {
	tests := []struct {
		name     string
		rngValue uint
		dice     Dice
		expected uint8
	}{
		{name: "MinD4", rngValue: 0, dice: D4, expected: 1},
		{name: "MidD4", rngValue: 1, dice: D4, expected: 2},
		{name: "MaxD4", rngValue: 3, dice: D4, expected: 4},

		{name: "MinD6", rngValue: 0, dice: D6, expected: 1},
		{name: "MidD6", rngValue: 2, dice: D6, expected: 3},
		{name: "MaxD6", rngValue: 5, dice: D6, expected: 6},

		{name: "MinD8", rngValue: 0, dice: D8, expected: 1},
		{name: "MidD8", rngValue: 3, dice: D8, expected: 4},
		{name: "MaxD8", rngValue: 7, dice: D8, expected: 8},

		{name: "MinD10", rngValue: 0, dice: D10, expected: 1},
		{name: "MidD10", rngValue: 4, dice: D10, expected: 5},
		{name: "MaxD10", rngValue: 9, dice: D10, expected: 10},

		{name: "MinD12", rngValue: 0, dice: D12, expected: 1},
		{name: "MidD12", rngValue: 5, dice: D12, expected: 6},
		{name: "MaxD12", rngValue: 11, dice: D12, expected: 12},

		{name: "MinD20", rngValue: 0, dice: D20, expected: 1},
		{name: "MidD20", rngValue: 9, dice: D20, expected: 10},
		{name: "MaxD20", rngValue: 19, dice: D20, expected: 20},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rng := fixedRNG(test.rngValue)
			actual := test.dice.Roll(rng)
			if actual != test.expected {
				t.Errorf("Dice.Roll() = %d, want %d", actual, test.expected)
			}
		})
	}
}

type rememberArgRNG struct{ arg uint }

func (r *rememberArgRNG) UintN(n uint) uint { r.arg = n; return 0 }

func TestDiceRollRNGArg(t *testing.T) {
	tests := []struct {
		name     string
		dice     Dice
		expected uint
	}{
		{name: "D4", dice: D4, expected: 4},
		{name: "D6", dice: D6, expected: 6},
		{name: "D8", dice: D8, expected: 8},
		{name: "D10", dice: D10, expected: 10},
		{name: "D12", dice: D12, expected: 12},
		{name: "D20", dice: D20, expected: 20},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &rememberArgRNG{0}
			_ = test.dice.Roll(r)
			if r.arg != test.expected {
				t.Errorf(
					"Dice invoked RNG.UintN() with %d, want %d", r.arg, test.expected,
				)
			}
		})
	}
}
