package benchmark

import (
	uecs "github.com/unitoftime/ecs"
	"golang.org/x/exp/rand"
	"testing"
)

func BenchmarkCreateEntityUECS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := uecs.NewWorld()

		for range ENTITIES_COUNT {
			id := world.NewId()
			world.Write(id,
				uecs.C(testTransform{
					x: rand.Float64() * 100,
					y: rand.Float64() * 100,
					z: rand.Float64() * 100,
				}),
				uecs.C(testTag{}),
			)
		}
	}

	b.ReportAllocs()
}

func BenchmarkIterateUECS(b *testing.B) {
	b.StopTimer()
	world := uecs.NewWorld()

	for i := 0; i < ENTITIES_COUNT; i++ {
		id := world.NewId()
		world.Write(id,
			uecs.C(testTransform{}),
			uecs.C(testTag{}),
		)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		query := uecs.Query2[testTransform, testTag](world)
		query.MapId(func(id uecs.Id, tr *testTransform, tag *testTag) {
			transformData(tr)
		})
	}

	b.ReportAllocs()
}

func BenchmarkAddUECS(b *testing.B) {
	b.StopTimer()

	world := uecs.NewWorld()

	entities := make([]uecs.Id, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		id := world.NewId()
		world.Write(id, uecs.C(testTag{}))

		entities = append(entities, id)
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for _, entityId := range entities {
			world.Write(entityId, uecs.C(testTransform{}))
		}

		b.StopTimer()
		for _, entityId := range entities {
			uecs.DeleteComponent(world, entityId, uecs.C(testTransform{}))
		}
	}

	b.ReportAllocs()
}

func BenchmarkRemoveUECS(b *testing.B) {
	b.StopTimer()

	world := uecs.NewWorld()

	entities := make([]uecs.Id, 0, ENTITIES_COUNT)
	for range ENTITIES_COUNT {
		id := world.NewId()
		world.Write(id, uecs.C(testTag{}))

		entities = append(entities, id)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, entityId := range entities {
			world.Write(entityId, uecs.C(testTransform{}))
		}

		b.StartTimer()
		for _, entityId := range entities {
			uecs.DeleteComponent(world, entityId, uecs.C(testTransform{}))
		}
	}

	b.ReportAllocs()
}
