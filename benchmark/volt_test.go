package benchmark

import (
	"math/rand/v2"
	"testing"

	"github.com/akmonengine/volt"
)

func BenchmarkCreateEntityVolt(b *testing.B) {
	for b.Loop() {
		world := volt.CreateWorld(ENTITIES_COUNT)
		volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
		volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

		for range ENTITIES_COUNT {
			volt.CreateEntityWithComponents2(world,
				testTransform{
					x: rand.Float64() * 100,
					y: rand.Float64() * 100,
					z: rand.Float64() * 100,
				},
				testTag{},
			)
		}
	}

	b.ReportAllocs()
}

func BenchmarkIterateVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.CreateEntity()
		volt.AddComponent[testTransform](world, id, testTransform{})
		volt.AddComponent[testTag](world, id, testTag{})
	}

	for b.Loop() {
		query := volt.CreateQuery2[testTransform, testTag](world, volt.QueryConfiguration{})
		for result := range query.Foreach(nil) {
			transformData(result.A)
		}
	}

	b.ReportAllocs()
}

func BenchmarkIterateConcurrentlyVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.CreateEntity()
		volt.AddComponent[testTransform](world, id, testTransform{})
		volt.AddComponent[testTag](world, id, testTag{})
	}

	for b.Loop() {
		query := volt.CreateQuery2[testTransform, testTag](world, volt.QueryConfiguration{})
		queryChannel := query.ForeachChannel(ENTITIES_COUNT/WORKERS, nil)

		runWorkers(WORKERS, func(workerId int) {
			for results := range queryChannel {
				for result := range results {
					transformData(result.A)
				}
			}
		})
	}

	b.ReportAllocs()
}

func BenchmarkTaskVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.CreateEntity()
		volt.AddComponent[testTransform](world, id, testTransform{})
		volt.AddComponent[testTag](world, id, testTag{})
	}

	for b.Loop() {
		query := volt.CreateQuery2[testTransform, testTag](world, volt.QueryConfiguration{})
		query.Task(WORKERS, nil, func(result volt.QueryResult2[testTransform, testTag]) {
			transformData(result.A)
		})
	}

	b.ReportAllocs()
}

func BenchmarkAddVolt(b *testing.B) {
	b.StopTimer()

	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	entities := make([]volt.EntityId, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		entityId := world.CreateEntity()
		volt.AddComponent(world, entityId, testTag{})
		entities = append(entities, entityId)
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, entityId := range entities {
			volt.AddComponent(world, entityId, testTransform{})
		}

		b.StopTimer()
		for _, entityId := range entities {
			volt.RemoveComponent[testTransform](world, entityId)
		}
	}

	b.ReportAllocs()
}

func BenchmarkRemoveVolt(b *testing.B) {
	b.StopTimer()

	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	entities := make([]volt.EntityId, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		entityId := world.CreateEntity()
		volt.AddComponent(world, entityId, testTag{})
		entities = append(entities, entityId)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, entityId := range entities {
			volt.AddComponent(world, entityId, testTransform{})
		}

		b.StartTimer()
		for _, entityId := range entities {
			volt.RemoveComponent[testTransform](world, entityId)
		}
	}

	b.ReportAllocs()
}

// BenchmarkCreateRemoveVolt measures a full recycle cycle: every entity is
// removed then recreated with the same components, so the id pool is exercised
// on both sides (free then reuse) on every iteration.
func BenchmarkCreateRemoveVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	entities := make([]volt.EntityId, ENTITIES_COUNT)
	for i := range entities {
		entities[i], _ = volt.CreateEntityWithComponents2(world, testTransform{}, testTag{})
	}

	for b.Loop() {
		for i, entityId := range entities {
			world.RemoveEntity(entityId)
			entities[i], _ = volt.CreateEntityWithComponents2(world, testTransform{}, testTag{})
		}
	}

	b.ReportAllocs()
}

// BenchmarkCreateRemoveLargeVolt: remove+recreate cycle of entities carrying 8
// components (the create-large-entities scenario of go-ecs-benchmarks).
func BenchmarkCreateRemoveLargeVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})
	volt.RegisterComponent[testC3](world, &volt.ComponentConfig[testC3]{})
	volt.RegisterComponent[testC4](world, &volt.ComponentConfig[testC4]{})
	volt.RegisterComponent[testC5](world, &volt.ComponentConfig[testC5]{})
	volt.RegisterComponent[testC6](world, &volt.ComponentConfig[testC6]{})
	volt.RegisterComponent[testC7](world, &volt.ComponentConfig[testC7]{})
	volt.RegisterComponent[testC8](world, &volt.ComponentConfig[testC8]{})

	create := func() volt.EntityId {
		entityId, _ := volt.CreateEntityWithComponents8(world, testTransform{}, testTag{}, testC3{}, testC4{}, testC5{}, testC6{}, testC7{}, testC8{})
		return entityId
	}

	entities := make([]volt.EntityId, ENTITIES_COUNT)
	for i := range entities {
		entities[i] = create()
	}

	for b.Loop() {
		for i, entityId := range entities {
			world.RemoveEntity(entityId)
			entities[i] = create()
		}
	}

	b.ReportAllocs()
}

// BenchmarkCreateRemoveFragmentedVolt: the same remove+recreate cycle as
// BenchmarkCreateRemoveVolt, in a world that already holds 256 other
// archetypes (every subset of 8 tags). It exposes how the cost of resolving an
// entity's archetype scales with the number of archetypes in the world.
func BenchmarkCreateRemoveFragmentedVolt(b *testing.B) {
	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for subset := range 256 {
		entityId := world.CreateEntity()
		for bit := range 8 {
			if subset&(1<<bit) != 0 {
				if err := world.AddTag(volt.TAGS_INDICES+volt.TagId(bit), entityId); err != nil {
					b.Fatal(err)
				}
			}
		}
	}

	entities := make([]volt.EntityId, ENTITIES_COUNT)
	for i := range entities {
		entities[i], _ = volt.CreateEntityWithComponents2(world, testTransform{}, testTag{})
	}

	for b.Loop() {
		for i, entityId := range entities {
			world.RemoveEntity(entityId)
			entities[i], _ = volt.CreateEntityWithComponents2(world, testTransform{}, testTag{})
		}
	}

	b.ReportAllocs()
}
