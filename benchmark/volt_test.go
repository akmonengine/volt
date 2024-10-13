package benchmark

import (
	"github.com/akmonengine/volt"
	"math/rand/v2"
	"strconv"
	"testing"
)

func BenchmarkCreateEntityVolt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := volt.CreateWorld(ENTITIES_COUNT)
		volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
		volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

		for j := range ENTITIES_COUNT {
			volt.CreateEntityWithComponents2(world, strconv.Itoa(j),
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
	b.StopTimer()

	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.CreateEntity(strconv.Itoa(i))
		volt.AddComponent[testTransform](world, id, testTransform{})
		volt.AddComponent[testTag](world, id, testTag{})
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		query := volt.CreateQuery2[testTransform, testTag](world, []volt.OptionalComponent{})
		for result := range query.Foreach(nil) {
			transformData(result.A)
		}
	}

	b.ReportAllocs()
}

func BenchmarkIterateConcurrentlyVolt(b *testing.B) {
	b.StopTimer()

	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.CreateEntity(strconv.Itoa(i))
		volt.AddComponent[testTransform](world, id, testTransform{})
		volt.AddComponent[testTag](world, id, testTag{})
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		query := volt.CreateQuery2[testTransform, testTag](world, []volt.OptionalComponent{})
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

func BenchmarkAddVolt(b *testing.B) {
	b.StopTimer()

	world := volt.CreateWorld(ENTITIES_COUNT)
	volt.RegisterComponent[testTransform](world, &volt.ComponentConfig[testTransform]{})
	volt.RegisterComponent[testTag](world, &volt.ComponentConfig[testTag]{})

	entities := make([]volt.EntityId, 0, ENTITIES_COUNT)
	for j := range ENTITIES_COUNT {
		entityId := world.CreateEntity(strconv.Itoa(j))
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
	for j := range ENTITIES_COUNT {
		entityId := world.CreateEntity(strconv.Itoa(j))
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
