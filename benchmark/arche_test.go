package benchmark

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"math/rand/v2"
	"testing"
)

func BenchmarkCreateEntityArche(b *testing.B) {
	for b.Loop() {
		world := ecs.NewWorld(ENTITIES_COUNT)
		mapper := generic.NewMap2[testTransform, testTag](&world)

		for range ENTITIES_COUNT {
			// Create a new Entity with components.
			entity := mapper.New()
			// Get the components
			tr, _ := mapper.Get(entity)
			// Initialize component fields.
			tr.x = rand.Float64() * 100
			tr.y = rand.Float64() * 100
			tr.z = rand.Float64() * 100
		}
	}

	b.ReportAllocs()
}

func BenchmarkIterateArche(b *testing.B) {
	world := ecs.NewWorld()
	mapper := generic.NewMap2[
		testTransform,
		testTag,
	](&world)

	mapper.NewBatch(ENTITIES_COUNT)

	for b.Loop() {
		filter := generic.NewFilter2[testTransform, testTag]()
		query := filter.Query(&world)
		for query.Next() {
			tr, _ := query.Get()
			transformData(tr)
		}
	}

	b.ReportAllocs()
}

func BenchmarkAddArche(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld(ENTITIES_COUNT)
	mapper := generic.NewMap1[testTag](&world)

	entities := make([]ecs.Entity, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		entities = append(entities, mapper.New())
	}
	componentId := ecs.ComponentID[testTransform](&world)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, e := range entities {
			world.Add(e, componentId)
		}

		b.StopTimer()
		for _, e := range entities {
			world.Remove(e, componentId)
		}
	}

	b.ReportAllocs()
}

func BenchmarkRemoveArche(b *testing.B) {
	b.StopTimer()

	world := ecs.NewWorld(ENTITIES_COUNT)
	mapper := generic.NewMap1[testTag](&world)

	entities := make([]ecs.Entity, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		entities = append(entities, mapper.New())
	}
	componentId := ecs.ComponentID[testTransform](&world)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, e := range entities {
			world.Add(e, componentId)
		}

		b.StartTimer()
		for _, e := range entities {
			world.Remove(e, componentId)
		}
	}

	b.ReportAllocs()
}
