package hysteria2

import (
	"reflect"
	"testing"
)

// Regression for cedar2025/Xboard-Node#54: IDs must be stable across hot
// updates and must never be list indices — removing users must neither
// misattribute surviving connections nor panic lookups.
func TestUserIDAllocator_StableAcrossUpdates(t *testing.T) {
	a := newUserIDAllocator()

	ids1 := a.reset([]string{"alice", "bob", "carol"})
	if len(ids1) != 3 {
		t.Fatalf("len(ids1) = %d", len(ids1))
	}
	alice1 := ids1[0]

	// bob removed, dave added; alice keeps her ID
	ids2 := a.reset([]string{"alice", "carol", "dave"})
	if ids2[0] != alice1 {
		t.Errorf("alice id changed across update: %d -> %d", alice1, ids2[0])
	}
	if ids2[1] == ids1[1] || ids2[2] == ids1[1] {
		t.Errorf("removed user bob's id %d was reused for another user", ids1[1])
	}

	// surviving long connection still resolves to alice
	if name, ok := a.name(alice1); !ok || name != "alice" {
		t.Errorf("name(alice1) = %q, %v", name, ok)
	}

	// stale ID from a removed user resolves to nothing — no panic
	if name, ok := a.name(ids1[1]); ok {
		t.Errorf("removed bob's id still resolves to %q", name)
	}
	_, _ = a.name(9999) // out-of-range must not panic

	// re-adding a removed name keeps traffic continuity on the same ID
	ids3 := a.reset([]string{"alice", "bob"})
	if ids3[1] != ids1[1] {
		t.Errorf("re-added bob got new id %d, want %d", ids3[1], ids1[1])
	}
	if !reflect.DeepEqual(ids3, []int{alice1, ids1[1]}) {
		t.Errorf("ids3 = %v", ids3)
	}
}
