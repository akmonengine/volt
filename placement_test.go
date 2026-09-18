package volt

import "testing"

// TestEmptyArchetypeReleasesEntityOnAddComponent: an entity created bare is
// listed in the empty archetype; adding a component must move it out through a
// swap-remove, not leave a stale entry behind. A stale entry leaks and is
// visible to any query whose filter matches the empty archetype (e.g. a query
// with only optional components), which then reports the entity twice.
func TestEmptyArchetypeReleasesEntityOnAddComponent(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	e := world.CreateEntity()
	if got := len(world.archetypes[0].entities); got != 1 {
		t.Fatalf("bare entity should be listed once in the empty archetype, got %d", got)
	}

	if err := AddComponent[testComponent1](world, e, testComponent1{}); err != nil {
		t.Fatalf("%s", err.Error())
	}

	if got := len(world.archetypes[0].entities); got != 0 {
		t.Fatalf("empty archetype still lists %d entity(ies) after the component was added: stale entry", got)
	}

	query := CreateQuery1[testComponent1](world, QueryConfiguration{OptionalComponents: []OptionalComponent{testComponent1Id}})
	if n := query.Count(); n != 1 {
		t.Fatalf("optional-only query counted %d results for a single entity (stale entry in the empty archetype)", n)
	}
}

// TestStaleEntryCorruptsLiveEntity: the stale entry is worse than a leak. When
// the empty archetype later swap-removes one of its entities, the stale entry
// may be the one moved, and the swap rewrites the key of the entity it names —
// an entity that now lives in another archetype, whose row then aliases a
// neighbour's.
func TestStaleEntryCorruptsLiveEntity(t *testing.T) {
	world := CreateWorld(16)
	RegisterComponent[testComponent1](world, &ComponentConfig[testComponent1]{})

	a := world.CreateEntity()
	b := world.CreateEntity()
	if err := AddComponent[testComponent1](world, b, testComponent1{}); err != nil {
		t.Fatal(err)
	}
	c := world.CreateEntity()
	if err := AddComponent[testComponent1](world, c, testComponent1{}); err != nil {
		t.Fatal(err)
	}
	// {C1}: [b, c]. With a stale entry, the empty archetype would still list b and c
	// behind a, and removing a would move the stale c onto a's row, rewriting c's key.
	world.RemoveEntity(a)

	if got := world.entities[c.Index()].key; got != 1 {
		t.Fatalf("live entity c should still be at key 1 in its archetype, got key %d", got)
	}
	GetComponent[testComponent1](world, c).x = 42
	if GetComponent[testComponent1](world, b).x == 42 {
		t.Fatal("writing c's component wrote into b's: rows are aliased")
	}
}
