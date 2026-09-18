package volt

import (
	"math"
	"testing"
)

// TestNullHandle: the zero EntityId is the null handle; it never exists, and no
// live entity is ever bound to it (generations start at 1).
func TestNullHandle(t *testing.T) {
	world := CreateWorld(16)

	var null EntityId
	if !null.IsNull() {
		t.Fatal("the zero EntityId should be the null handle")
	}
	if world.Exists(null) {
		t.Fatal("the null handle should not exist")
	}

	e := world.CreateEntity()
	if e.IsNull() {
		t.Fatal("the first created entity should not be the null handle")
	}
	if e.Index() != 0 || e.Generation() != 1 {
		t.Fatalf("the first entity should be slot 0 generation 1, got slot %d generation %d", e.Index(), e.Generation())
	}
}

// TestGenerationalHandles: a handle kept after RemoveEntity is dead, even
// once its slot is reused by a new entity: every API refuses it and never
// answers for the slot's new occupant.
func TestGenerationalHandles(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	e := world.CreateEntity()
	if err := AddComponent[testComponent1](world, e, testComponent1{}); err != nil {
		t.Fatalf("%s", err.Error())
	}
	world.RemoveEntity(e)

	f := world.CreateEntity()
	if f.Index() != e.Index() {
		t.Fatalf("the freed slot should be recycled: got slot %d, want %d", f.Index(), e.Index())
	}
	if f == e {
		t.Fatal("a recycled slot must yield a distinct handle")
	}
	if f.Generation() != e.Generation()+1 {
		t.Fatalf("recycling should bump the generation: got %d, want %d", f.Generation(), e.Generation()+1)
	}
	if err := AddComponent[testComponent1](world, f, testComponent1{}); err != nil {
		t.Fatalf("%s", err.Error())
	}

	// The dead handle answers nothing about the new occupant.
	if world.Exists(e) {
		t.Fatal("a dead handle should not exist")
	}
	if world.HasComponents(e, testComponent1Id) {
		t.Fatal("a dead handle should not report the new occupant's components")
	}
	if GetComponent[testComponent1](world, e) != nil {
		t.Fatal("GetComponent through a dead handle should return nil")
	}
	if _, err := world.GetComponent(e, testComponent1Id); err == nil {
		t.Fatal("World.GetComponent through a dead handle should return an error")
	}
	if err := AddComponent[testComponent1](world, e, testComponent1{}); err == nil {
		t.Fatal("AddComponent through a dead handle should return an error")
	}
	if err := RemoveComponent[testComponent1](world, e); err == nil {
		t.Fatal("RemoveComponent through a dead handle should return an error")
	}
	if err := world.AddTag(TAGS_INDICES, e); err == nil {
		t.Fatal("AddTag through a dead handle should return an error")
	}
	if world.HasTag(TAGS_INDICES, e) {
		t.Fatal("a dead handle should not report tags")
	}
	if err := world.RemoveTag(TAGS_INDICES, e); err == nil {
		t.Fatal("RemoveTag through a dead handle should return an error")
	}

	// Removing a dead handle is a no-op: the new occupant is untouched.
	world.RemoveEntity(e)
	if !world.Exists(f) {
		t.Fatal("removing a dead handle must not remove the slot's new occupant")
	}
	if GetComponent[testComponent1](world, f) == nil {
		t.Fatal("the new occupant should still own its component")
	}
	if world.Count() != 1 {
		t.Fatalf("world should count 1 entity, got %d", world.Count())
	}
}

// TestGenerationCycles: a million remove/create cycles on the same slot
// never hand out a handle colliding with a previous one, and never revive a
// dead handle.
func TestGenerationCycles(t *testing.T) {
	world := CreateWorld(16)

	e := world.CreateEntity()
	for i := range 1_000_000 {
		world.RemoveEntity(e)
		if world.Exists(e) {
			t.Fatalf("cycle %d: the removed handle still exists", i)
		}

		f := world.CreateEntity()
		if f == e || f.Index() != e.Index() || f.Generation() != e.Generation()+1 {
			t.Fatalf("cycle %d: expected slot %d generation %d, got slot %d generation %d", i, e.Index(), e.Generation()+1, f.Index(), f.Generation())
		}
		if world.Exists(e) {
			t.Fatalf("cycle %d: the dead handle was revived by the slot's new occupant", i)
		}

		e = f
	}

	if world.Count() != 1 {
		t.Fatalf("world should count 1 entity, got %d", world.Count())
	}
}

// TestGenerationWrapSkipsNull: bumping the last generation wraps to 1, not 0,
// so the null handle can never be minted.
func TestGenerationWrapSkipsNull(t *testing.T) {
	if got := nextGeneration(1); got != 2 {
		t.Fatalf("nextGeneration(1) = %d, want 2", got)
	}
	if got := nextGeneration(math.MaxUint32); got != 1 {
		t.Fatalf("nextGeneration(MaxUint32) = %d, want 1 (0 is reserved for the null handle)", got)
	}
}

// TestFailedCreationGivesBackTheSlot: a creation that fails leaves no
// half-placed entity behind, returns the null handle, and its slot is reused.
func TestFailedCreationGivesBackTheSlot(t *testing.T) {
	world := CreateWorld(16)
	// testComponent2 is not registered: the creation cannot find its storage.
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	e, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})
	if err == nil {
		t.Fatal("creating with an unregistered component should fail")
	}
	if !e.IsNull() {
		t.Fatalf("a failed creation should return the null handle, got %d", e)
	}
	if world.Count() != 0 {
		t.Fatalf("a failed creation should leave the world empty, got %d", world.Count())
	}

	f := world.CreateEntity()
	if f.Index() != 0 {
		t.Fatalf("the discarded slot should be reused, got slot %d", f.Index())
	}
	if !world.Exists(f) {
		t.Fatal("the entity created in the reused slot should exist")
	}
}
