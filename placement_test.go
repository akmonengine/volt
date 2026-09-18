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
