package volt

import "testing"

// TestArchetypeGraph_ReallocationSafety drives the number of archetypes
// well past the 1024 preallocated capacity, so world.archetypes reallocates
// several times. It guards against a stale-pointer hazard: add/remove resolve a
// destination archetype (which may grow world.archetypes) while operating on the
// source archetype. Component values and tag membership must survive intact.
func TestArchetypeGraph_ReallocationSafety(t *testing.T) {
	const n = 1500 // exceeds the 1024 archetype preallocation

	world := CreateWorld(n)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	entities := make([]EntityId, n)
	for i := 0; i < n; i++ {
		e := world.CreateEntity()
		entities[i] = e

		c := testComponent1{}
		c.x = i
		if err := AddComponent(world, e, c); err != nil {
			t.Fatalf("AddComponent: %s", err.Error())
		}
		// A distinct tag per entity forces a distinct archetype {c1, tag_i}.
		if err := world.AddTag(TAGS_INDICES+TagId(i), e); err != nil {
			t.Fatalf("AddTag: %s", err.Error())
		}
	}

	// Every component value must have survived the archetype reallocations.
	for i, e := range entities {
		c := GetComponent[testComponent1](world, e)
		if c == nil {
			t.Fatalf("entity %d lost its component", e)
		}
		if c.x != i {
			t.Fatalf("entity %d: expected component x=%d, got %d", e, i, c.x)
		}
	}

	// Removing the component drives ~n more archetype creations (well past 1024)
	// exactly while removeComponent holds the source archetype — the hazard.
	for _, e := range entities {
		if err := RemoveComponent[testComponent1](world, e); err != nil {
			t.Fatalf("RemoveComponent: %s", err.Error())
		}
	}

	for i, e := range entities {
		if world.HasComponents(e, testComponent1Id) {
			t.Fatalf("entity %d still owns the component after removal", e)
		}
		if !world.HasTag(TAGS_INDICES+TagId(i), e) {
			t.Fatalf("entity %d lost its tag after component removal", e)
		}
	}
}

// TestArchetypeGraph_EdgesAreReused checks that repeated identical transitions
// resolve to the same archetype (the graph stays consistent across many hops),
// by cycling a component on and off and confirming the entity returns to the
// exact same archetype each time.
func TestArchetypeGraph_EdgesAreReused(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	e := world.CreateEntity()
	if err := AddComponent(world, e, testComponent2{}); err != nil {
		t.Fatalf("%s", err.Error())
	}
	base := world.entities[e].archetypeId
	var withC1 archetypeId

	for i := 0; i < 5; i++ {
		if err := AddComponent(world, e, testComponent1{}); err != nil {
			t.Fatalf("add iteration %d: %s", i, err.Error())
		}
		with := world.entities[e].archetypeId
		if i == 0 {
			withC1 = with
		} else if with != withC1 {
			t.Fatalf("iteration %d: the {c1,c2} archetype id is unstable (%d != %d)", i, with, withC1)
		}

		if err := RemoveComponent[testComponent1](world, e); err != nil {
			t.Fatalf("remove iteration %d: %s", i, err.Error())
		}
		if back := world.entities[e].archetypeId; back != base {
			t.Fatalf("iteration %d: entity did not return to base archetype (%d != %d)", i, back, base)
		}
	}
}
