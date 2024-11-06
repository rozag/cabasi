package battle

import (
	"slices"
	"testing"
)

func TestPickedTargetsEquals(t *testing.T) {
	tests := []struct {
		name string
		pt1  PickedTargets
		pt2  PickedTargets
		want bool
	}{
		{
			name: "EqualPickedTargets",
			pt1: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			pt2: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			want: true,
		},
		{
			name: "DifferentAttackerID",
			pt1: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			pt2: PickedTargets{
				AttackerID:  CreatureID("attacker2"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			want: false,
		},
		{
			name: "DifferentAttackIdx",
			pt1: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			pt2: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   2,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			want: false,
		},
		{
			name: "DifferentDefenderIDs",
			pt1: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender1", "defender2"},
			},
			pt2: PickedTargets{
				AttackerID:  CreatureID("attacker1"),
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"defender3", "defender4"},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.pt1.Equals(test.pt2); got != test.want {
				t.Errorf("Equals() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPickTargetsFirstAlive(t *testing.T) {
	tests := []struct {
		attackers map[CreatureID]Creature
		defenders map[CreatureID]Creature
		want      []PickedTargets
		name      string
	}{
		{
			name:      "NilAttackers",
			attackers: nil,
			defenders: TODO,
			want:      nil,
		},
		{
			name:      "EmptyAttackers",
			attackers: map[CreatureID]Creature{},
			defenders: TODO,
			want:      nil,
		},
		{
			name:      "NilDefenders",
			attackers: TODO,
			defenders: nil,
			want:      nil,
		},
		{
			name:      "EmptyDefenders",
			attackers: TODO,
			defenders: map[CreatureID]Creature{},
			want:      nil,
		},
		{
			name:      "NilAttackersNilDefenders",
			attackers: nil,
			defenders: nil,
			want:      nil,
		},
		{
			name:      "EmptyAttackersEmptyDefenders",
			attackers: map[CreatureID]Creature{},
			defenders: map[CreatureID]Creature{},
			want:      nil,
		},
		{
			name:      "AllAttackersOut",
			attackers: TODO,
			defenders: TODO,
			want:      nil,
		},
		{
			name:      "AllDefendersOut",
			attackers: TODO,
			defenders: TODO,
			want:      nil,
		},
		{
			name:      "AllAttackersOutAllDefendersOut",
			attackers: TODO,
			defenders: TODO,
			want:      nil,
		},
		{
			name:      TODO,
			attackers: TODO,
			defenders: TODO,
			want:      TODO,
		},
		{
			name:      TODO,
			attackers: TODO,
			defenders: TODO,
			want:      TODO,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := PickTargetsFirstAlive(test.attackers, test.defenders)
			slices.SortFunc(got, func(a, b PickedTargets) int {
				return a.AttackerID.CompareTo(b.AttackerID)
			})
			if !slices.EqualFunc(got, test.want, PickedTargets.Equals) {
				t.Fatalf("PickTargetsFirstAlive() = %v, want %v", got, test.want)
			}
		})
	}
}
