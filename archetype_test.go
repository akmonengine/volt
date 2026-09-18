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
	base := world.entities[e.Index()].archetypeId
	var withC1 archetypeId

	for i := 0; i < 5; i++ {
		if err := AddComponent(world, e, testComponent1{}); err != nil {
			t.Fatalf("add iteration %d: %s", i, err.Error())
		}
		with := world.entities[e.Index()].archetypeId
		if i == 0 {
			withC1 = with
		} else if with != withC1 {
			t.Fatalf("iteration %d: the {c1,c2} archetype id is unstable (%d != %d)", i, with, withC1)
		}

		if err := RemoveComponent[testComponent1](world, e); err != nil {
			t.Fatalf("remove iteration %d: %s", i, err.Error())
		}
		if back := world.entities[e.Index()].archetypeId; back != base {
			t.Fatalf("iteration %d: entity did not return to base archetype (%d != %d)", i, back, base)
		}
	}
}

// TestArchetypeSetIsOrderIndependent: the components of a set, given in any
// order, resolve to the same archetype, whether the entity is created with
// them or receives them afterwards.
func TestArchetypeSetIsOrderIndependent(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	e1, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})
	if err != nil {
		t.Fatal(err)
	}
	e2, err := CreateEntityWithComponents2(world, testComponent2{}, testComponent1{})
	if err != nil {
		t.Fatal(err)
	}
	e3 := world.CreateEntity()
	if err := AddComponents2(world, e3, testComponent2{}, testComponent1{}); err != nil {
		t.Fatal(err)
	}

	a1, a2, a3 := world.entities[e1.Index()].archetypeId, world.entities[e2.Index()].archetypeId, world.entities[e3.Index()].archetypeId
	if a1 != a2 || a1 != a3 {
		t.Fatalf("{c1,c2} resolved to archetypes %d, %d and %d depending on the order", a1, a2, a3)
	}
	if got := len(world.archetypes); got != 2 {
		t.Fatalf("expected the empty archetype and {c1,c2} only, got %d archetypes", got)
	}
	for _, e := range []EntityId{e1, e2, e3} {
		if GetComponent[testComponent1](world, e) == nil || GetComponent[testComponent2](world, e) == nil {
			t.Fatalf("entity %d lost a component", e)
		}
	}
}

// TestArchetypeSetBeyondInlineScratch: a set larger than the inline scratch
// (many tags plus components) still resolves, and resolves once.
func TestArchetypeSetBeyondInlineScratch(t *testing.T) {
	const tags = maxInlineComponents + 1

	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	build := func() EntityId {
		e := world.CreateEntity()
		for i := range tags {
			if err := world.AddTag(TAGS_INDICES+TagId(i), e); err != nil {
				t.Fatal(err)
			}
		}
		if err := AddComponents2(world, e, testComponent1{}, testComponent2{}); err != nil {
			t.Fatal(err)
		}

		return e
	}
	e1, e2 := build(), build()

	a1, a2 := world.entities[e1.Index()].archetypeId, world.entities[e2.Index()].archetypeId
	if a1 != a2 {
		t.Fatalf("the same large set resolved to archetypes %d and %d", a1, a2)
	}
	if got := len(world.archetypes[a1].Type); got != tags+2 {
		t.Fatalf("expected an archetype of %d components, got %d", tags+2, got)
	}
	for i := range tags {
		if !world.HasTag(TAGS_INDICES+TagId(i), e1) {
			t.Fatalf("entity lost tag %d", i)
		}
	}
	if !world.HasComponents(e1, testComponent1Id, testComponent2Id) {
		t.Fatal("entity lost its components")
	}
}

// TestArchetypeKeyCollisionFallsBackToScan: when two sets share a key, the
// second one is still resolved (by scan) to a single archetype of its own, and
// the first one keeps the key.
func TestArchetypeKeyCollisionFallsBackToScan(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	single := world.CreateEntity()
	if err := AddComponent(world, single, testComponent1{}); err != nil {
		t.Fatal(err)
	}
	owner := world.entities[single.Index()].archetypeId

	// Forge a collision: the key of {c1,c2} points at the {c1} archetype.
	pair := []ComponentId{testComponent1Id, testComponent2Id}
	world.archetypesByKey[archetypeKey(pair)] = owner

	e1, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})
	if err != nil {
		t.Fatal(err)
	}
	e2, err := CreateEntityWithComponents2(world, testComponent2{}, testComponent1{})
	if err != nil {
		t.Fatal(err)
	}

	a1, a2 := world.entities[e1.Index()].archetypeId, world.entities[e2.Index()].archetypeId
	if a1 == owner || a1 != a2 {
		t.Fatalf("colliding set resolved to archetypes %d and %d (key owner %d)", a1, a2, owner)
	}
	if got := len(world.archetypes); got != 3 {
		t.Fatalf("expected the empty archetype, {c1} and {c1,c2}, got %d archetypes", got)
	}
	if world.archetypesByKey[archetypeKey(pair)] != owner {
		t.Fatal("the colliding archetype stole the key from its first owner")
	}
	if GetComponent[testComponent2](world, e1) == nil || GetComponent[testComponent1](world, single) == nil {
		t.Fatal("a component was lost across the collision")
	}
}

// TestTransitionCache: a multi-component transition is cached after its first
// resolution; a slot whose key differs is never trusted and gets refilled; the
// same set in another order is another entry, but the same archetype.
func TestTransitionCache(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})
	archetypeOf := func(e EntityId) archetypeId { return world.entities[e.Index()].archetypeId }

	e1, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})
	if err != nil {
		t.Fatal(err)
	}
	slot := transitionSlot(0, []ComponentId{testComponent1Id, testComponent2Id})
	if entry := world.transitions[slot]; entry.from != 0 || entry.n != 2 || entry.dest != archetypeOf(e1) {
		t.Fatalf("the transition should be cached after its first resolution, got %+v", entry)
	}

	// Corrupt the key of the slot and point it at the wrong archetype: the
	// lookup must miss, resolve correctly, and refill the slot.
	world.transitions[slot].ids[0] = testComponent2Id
	world.transitions[slot].dest = 0
	e2, err := CreateEntityWithComponents2(world, testComponent1{}, testComponent2{})
	if err != nil {
		t.Fatal(err)
	}
	if archetypeOf(e2) != archetypeOf(e1) {
		t.Fatalf("a slot whose key differs must not be trusted: got archetype %d, want %d", archetypeOf(e2), archetypeOf(e1))
	}
	if world.transitions[slot].dest != archetypeOf(e1) {
		t.Fatal("the miss should have refilled the slot")
	}

	e3, err := CreateEntityWithComponents2(world, testComponent2{}, testComponent1{})
	if err != nil {
		t.Fatal(err)
	}
	if archetypeOf(e3) != archetypeOf(e1) {
		t.Fatal("the same set in another order should reach the same archetype")
	}
	if got := len(world.archetypes); got != 2 {
		t.Fatalf("expected the empty archetype and {c1,c2} only, got %d", got)
	}
}

// TestTransitionCacheFromNonEmptyArchetype: the transition is keyed by the
// archetype the entity starts from, so adding the same components to entities
// living in different archetypes leads to different, correct archetypes.
func TestTransitionCacheFromNonEmptyArchetype(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})
	RegisterComponent[testComponent2](world, &ComponentConfig[testComponent2]{})

	bare := world.CreateEntity()
	tagged := world.CreateEntity()
	if err := world.AddTag(TAGS_INDICES, tagged); err != nil {
		t.Fatal(err)
	}
	for _, e := range []EntityId{bare, tagged} {
		if err := AddComponents2(world, e, testComponent1{}, testComponent2{}); err != nil {
			t.Fatal(err)
		}
	}

	if world.entities[bare.Index()].archetypeId == world.entities[tagged.Index()].archetypeId {
		t.Fatal("entities starting from different archetypes must not share the destination")
	}
	if !world.HasComponents(bare, testComponent1Id, testComponent2Id) || world.HasTag(TAGS_INDICES, bare) {
		t.Fatal("the bare entity should own the two components and no tag")
	}
	if !world.HasComponents(tagged, testComponent1Id, testComponent2Id) || !world.HasTag(TAGS_INDICES, tagged) {
		t.Fatal("the tagged entity should own the two components and keep its tag")
	}
}
