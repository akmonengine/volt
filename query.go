package volt

import (
	"iter"
	"math"
	"slices"
)

type OptionalComponent ComponentId

type Query1[A ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult1[A ComponentInterface] struct {
	EntityId EntityId
	A        *A
}
type QueryResultChunk1[A ComponentInterface] struct {
	EntityId []EntityId
	A        []A
}

func CreateQuery1[A ComponentInterface](world *World, optionalComponents []OptionalComponent) Query1[A] {
	var a A
	return Query1[A]{
		World:              world,
		componentsIds:      world.getComponentsIds(a),
		optionalComponents: optionalComponents,
	}
}

func (query *Query1[A]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query1[A]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query1[A]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query1[A]) GetEntitiesIds() []EntityId {
	var entities []EntityId
	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query1[A]) Foreach(filterFn func(QueryResult1[A]) bool) iter.Seq[QueryResult1[A]] {
	return func(yield func(QueryResult1[A]) bool) {
		storageA := getStorage[A](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}

				result := QueryResult1[A]{
					EntityId: entityId,
					A:        dataA,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query1[A]) ForeachChannel(chunkSize int, filterFn func(QueryResult1[A]) bool) <-chan iter.Seq[QueryResult1[A]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult1[A]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk1[A]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult1[A]) bool) {
					queryResult := QueryResult1[A]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query2[A, B ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult2[A, B ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
}
type QueryResultChunk2[A, B ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
}

func CreateQuery2[A, B ComponentInterface](world *World, optionalComponents []OptionalComponent) Query2[A, B] {
	var a A
	var b B
	return Query2[A, B]{
		World:              world,
		componentsIds:      world.getComponentsIds(a, b),
		optionalComponents: optionalComponents,
	}
}

func (query *Query2[A, B]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query2[A, B]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query2[A, B]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query2[A, B]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query2[A, B]) Foreach(filterFn func(QueryResult2[A, B]) bool) iter.Seq[QueryResult2[A, B]] {
	return func(yield func(QueryResult2[A, B]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]

			var result QueryResult2[A, B]
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					result.A = &sliceA[i]
				}
				if sliceB != nil {
					result.B = &sliceB[i]
				}
				result.EntityId = entityId

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query2[A, B]) ForeachChannel(chunkSize int, filterFn func(QueryResult2[A, B]) bool) <-chan iter.Seq[QueryResult2[A, B]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult2[A, B]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk2[A, B]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult2[A, B]) bool) {
					queryResult := QueryResult2[A, B]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query3[A, B, C ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult3[A, B, C ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
}
type QueryResultChunk3[A, B, C ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
}

func CreateQuery3[A, B, C ComponentInterface](world *World, optionalComponents []OptionalComponent) Query3[A, B, C] {
	var a A
	var b B
	var c C
	return Query3[A, B, C]{
		World:              world,
		componentsIds:      world.getComponentsIds(a, b, c),
		optionalComponents: optionalComponents,
	}
}

func (query *Query3[A, B, C]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query3[A, B, C]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query3[A, B, C]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query3[A, B, C]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query3[A, B, C]) Foreach(filterFn func(QueryResult3[A, B, C]) bool) iter.Seq[QueryResult3[A, B, C]] {
	return func(yield func(QueryResult3[A, B, C]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}

				result := QueryResult3[A, B, C]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query3[A, B, C]) ForeachChannel(chunkSize int, filterFn func(QueryResult3[A, B, C]) bool) <-chan iter.Seq[QueryResult3[A, B, C]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult3[A, B, C]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk3[A, B, C]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult3[A, B, C]) bool) {
					queryResult := QueryResult3[A, B, C]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query4[A, B, C, D ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}

type QueryResult4[A, B, C, D ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
	D        *D
}

type QueryResultChunk4[A, B, C, D ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
	D        []D
}

func CreateQuery4[A, B, C, D ComponentInterface](world *World, optionalComponents []OptionalComponent) Query4[A, B, C, D] {
	var a A
	var b B
	var c C
	var d D
	return Query4[A, B, C, D]{
		World:              world,
		componentsIds:      world.getComponentsIds(a, b, c, d),
		optionalComponents: optionalComponents,
	}
}

func (query *Query4[A, B, C, D]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query4[A, B, C, D]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query4[A, B, C, D]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query4[A, B, C, D]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query4[A, B, C, D]) Foreach(filterFn func(QueryResult4[A, B, C, D]) bool) iter.Seq[QueryResult4[A, B, C, D]] {
	return func(yield func(QueryResult4[A, B, C, D]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			var dataD *D
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}
				if sliceD != nil {
					dataD = &sliceD[i]
				}

				result := QueryResult4[A, B, C, D]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
					D:        dataD,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query4[A, B, C, D]) ForeachChannel(chunkSize int, filterFn func(QueryResult4[A, B, C, D]) bool) <-chan iter.Seq[QueryResult4[A, B, C, D]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult4[A, B, C, D]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk4[A, B, C, D]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}
				if sliceD != nil {
					result.D = sliceD[i : i+end : i+end]
				}
				channel <- func(yield func(QueryResult4[A, B, C, D]) bool) {
					queryResult := QueryResult4[A, B, C, D]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}
						if result.D != nil {
							queryResult.D = &result.D[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query5[A, B, C, D, E ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult5[A, B, C, D, E ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
	D        *D
	E        *E
}
type QueryResultChunk5[A, B, C, D, E ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
	D        []D
	E        []E
}

func CreateQuery5[A, B, C, D, E ComponentInterface](world *World, optionalComponents []OptionalComponent) Query5[A, B, C, D, E] {
	var a A
	var b B
	var c C
	var d D
	var e E
	return Query5[A, B, C, D, E]{
		World:              world,
		componentsIds:      world.getComponentsIds(a, b, c, d, e),
		optionalComponents: optionalComponents,
	}
}

func (query *Query5[A, B, C, D, E]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query5[A, B, C, D, E]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query5[A, B, C, D, E]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query5[A, B, C, D, E]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query5[A, B, C, D, E]) Foreach(filterFn func(QueryResult5[A, B, C, D, E]) bool) iter.Seq[QueryResult5[A, B, C, D, E]] {
	return func(yield func(QueryResult5[A, B, C, D, E]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			var dataD *D
			var dataE *E
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}
				if sliceD != nil {
					dataD = &sliceD[i]
				}
				if sliceE != nil {
					dataE = &sliceE[i]
				}

				result := QueryResult5[A, B, C, D, E]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
					D:        dataD,
					E:        dataE,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query5[A, B, C, D, E]) ForeachChannel(chunkSize int, filterFn func(QueryResult5[A, B, C, D, E]) bool) <-chan iter.Seq[QueryResult5[A, B, C, D, E]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult5[A, B, C, D, E]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk5[A, B, C, D, E]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}
				if sliceD != nil {
					result.D = sliceD[i : i+end : i+end]
				}
				if sliceE != nil {
					result.E = sliceE[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult5[A, B, C, D, E]) bool) {
					queryResult := QueryResult5[A, B, C, D, E]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}
						if result.D != nil {
							queryResult.D = &result.D[k]
						}
						if result.E != nil {
							queryResult.E = &result.E[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query6[A, B, C, D, E, F ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult6[A, B, C, D, E, F ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
	D        *D
	E        *E
	F        *F
}
type QueryResultChunk6[A, B, C, D, E, F ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
	D        []D
	E        []E
	F        []F
}

func CreateQuery6[A, B, C, D, E, F ComponentInterface](world *World, optionalComponents []OptionalComponent) Query6[A, B, C, D, E, F] {
	var a A
	var b B
	var c C
	var d D
	var e E
	var f F
	return Query6[A, B, C, D, E, F]{World: world,
		componentsIds:      world.getComponentsIds(a, b, c, d, e, f),
		optionalComponents: optionalComponents,
	}
}

func (query *Query6[A, B, C, D, E, F]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query6[A, B, C, D, E, F]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query6[A, B, C, D, E, F]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query6[A, B, C, D, E, F]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query6[A, B, C, D, E, F]) Foreach(filterFn func(QueryResult6[A, B, C, D, E, F]) bool) iter.Seq[QueryResult6[A, B, C, D, E, F]] {
	return func(yield func(QueryResult6[A, B, C, D, E, F]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			var dataD *D
			var dataE *E
			var dataF *F
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}
				if sliceD != nil {
					dataD = &sliceD[i]
				}
				if sliceE != nil {
					dataE = &sliceE[i]
				}
				if sliceF != nil {
					dataF = &sliceF[i]
				}

				result := QueryResult6[A, B, C, D, E, F]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
					D:        dataD,
					E:        dataE,
					F:        dataF,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query6[A, B, C, D, E, F]) ForeachChannel(chunkSize int, filterFn func(QueryResult6[A, B, C, D, E, F]) bool) <-chan iter.Seq[QueryResult6[A, B, C, D, E, F]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult6[A, B, C, D, E, F]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk6[A, B, C, D, E, F]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}
				if sliceD != nil {
					result.D = sliceD[i : i+end : i+end]
				}
				if sliceE != nil {
					result.E = sliceE[i : i+end : i+end]
				}
				if sliceF != nil {
					result.F = sliceF[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult6[A, B, C, D, E, F]) bool) {
					queryResult := QueryResult6[A, B, C, D, E, F]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}
						if result.D != nil {
							queryResult.D = &result.D[k]
						}
						if result.E != nil {
							queryResult.E = &result.E[k]
						}
						if result.F != nil {
							queryResult.F = &result.F[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query7[A, B, C, D, E, F, G ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult7[A, B, C, D, E, F, G ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
	D        *D
	E        *E
	F        *F
	G        *G
}
type QueryResultChunk7[A, B, C, D, E, F, G ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
	D        []D
	E        []E
	F        []F
	G        []G
}

func CreateQuery7[A, B, C, D, E, F, G ComponentInterface](world *World, optionalComponents []OptionalComponent) Query7[A, B, C, D, E, F, G] {
	var a A
	var b B
	var c C
	var d D
	var e E
	var f F
	var g G
	return Query7[A, B, C, D, E, F, G]{World: world,
		componentsIds:      world.getComponentsIds(a, b, c, d, e, f, g),
		optionalComponents: optionalComponents,
	}
}

func (query *Query7[A, B, C, D, E, F, G]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query7[A, B, C, D, E, F, G]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query7[A, B, C, D, E, F, G]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query7[A, B, C, D, E, F, G]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query7[A, B, C, D, E, F, G]) Foreach(filterFn func(QueryResult7[A, B, C, D, E, F, G]) bool) iter.Seq[QueryResult7[A, B, C, D, E, F, G]] {
	return func(yield func(QueryResult7[A, B, C, D, E, F, G]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)
		storageG := getStorage[G](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]
			sliceG := storageG.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			var dataD *D
			var dataE *E
			var dataF *F
			var dataG *G
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}
				if sliceD != nil {
					dataD = &sliceD[i]
				}
				if sliceE != nil {
					dataE = &sliceE[i]
				}
				if sliceF != nil {
					dataF = &sliceF[i]
				}
				if sliceG != nil {
					dataG = &sliceG[i]
				}

				result := QueryResult7[A, B, C, D, E, F, G]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
					D:        dataD,
					E:        dataE,
					F:        dataF,
					G:        dataG,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query7[A, B, C, D, E, F, G]) ForeachChannel(chunkSize int, filterFn func(QueryResult7[A, B, C, D, E, F, G]) bool) <-chan iter.Seq[QueryResult7[A, B, C, D, E, F, G]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult7[A, B, C, D, E, F, G]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)
		storageG := getStorage[G](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]
			sliceG := storageG.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk7[A, B, C, D, E, F, G]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}
				if sliceD != nil {
					result.D = sliceD[i : i+end : i+end]
				}
				if sliceE != nil {
					result.E = sliceE[i : i+end : i+end]
				}
				if sliceF != nil {
					result.F = sliceF[i : i+end : i+end]
				}
				if sliceG != nil {
					result.G = sliceG[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult7[A, B, C, D, E, F, G]) bool) {
					queryResult := QueryResult7[A, B, C, D, E, F, G]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}
						if result.D != nil {
							queryResult.D = &result.D[k]
						}
						if result.E != nil {
							queryResult.E = &result.E[k]
						}
						if result.F != nil {
							queryResult.F = &result.F[k]
						}
						if result.G != nil {
							queryResult.G = &result.G[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}

type Query8[A, B, C, D, E, F, G, H ComponentInterface] struct {
	World              *World
	componentsIds      []ComponentId
	optionalComponents []OptionalComponent
}
type QueryResult8[A, B, C, D, E, F, G, H ComponentInterface] struct {
	EntityId EntityId
	A        *A
	B        *B
	C        *C
	D        *D
	E        *E
	F        *F
	G        *G
	H        *H
}
type QueryResultChunk8[A, B, C, D, E, F, G, H ComponentInterface] struct {
	EntityId []EntityId
	A        []A
	B        []B
	C        []C
	D        []D
	E        []E
	F        []F
	G        []G
	H        []H
}

func CreateQuery8[A, B, C, D, E, F, G, H ComponentInterface](world *World, optionalComponents []OptionalComponent) Query8[A, B, C, D, E, F, G, H] {
	var a A
	var b B
	var c C
	var d D
	var e E
	var f F
	var g G
	var h H
	return Query8[A, B, C, D, E, F, G, H]{World: world,
		componentsIds:      world.getComponentsIds(a, b, c, d, e, f, g, h),
		optionalComponents: optionalComponents,
	}
}

func (query *Query8[A, B, C, D, E, F, G, H]) GetComponentsIds() []ComponentId {
	return query.componentsIds
}

func (query *Query8[A, B, C, D, E, F, G, H]) filter() []Archetype {
	var componentsIds []ComponentId

	for _, componentId := range query.componentsIds {
		if !slices.Contains(query.optionalComponents, OptionalComponent(componentId)) {
			componentsIds = append(componentsIds, componentId)
		}
	}
	archetypes := query.World.getArchetypesForComponentsIds(componentsIds...)

	return archetypes
}

func (query *Query8[A, B, C, D, E, F, G, H]) Count() int {
	count := 0
	archetypes := query.filter()

	for _, archetype := range archetypes {
		count += len(archetype.entities)
	}

	return count
}

func (query *Query8[A, B, C, D, E, F, G, H]) GetEntitiesIds() []EntityId {
	var entities []EntityId

	archetypes := query.filter()

	for _, archetype := range archetypes {
		entities = append(entities, archetype.entities...)
	}

	return entities
}

func (query *Query8[A, B, C, D, E, F, G, H]) Foreach(filterFn func(QueryResult8[A, B, C, D, E, F, G, H]) bool) iter.Seq[QueryResult8[A, B, C, D, E, F, G, H]] {
	return func(yield func(QueryResult8[A, B, C, D, E, F, G, H]) bool) {
		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)
		storageG := getStorage[G](query.World)
		storageH := getStorage[H](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]
			sliceG := storageG.archetypesComponentsEntities[archetype.Id]
			sliceH := storageH.archetypesComponentsEntities[archetype.Id]
			var dataA *A
			var dataB *B
			var dataC *C
			var dataD *D
			var dataE *E
			var dataF *F
			var dataG *G
			var dataH *H
			for i, entityId := range archetype.entities {
				if sliceA != nil {
					dataA = &sliceA[i]
				}
				if sliceB != nil {
					dataB = &sliceB[i]
				}
				if sliceC != nil {
					dataC = &sliceC[i]
				}
				if sliceD != nil {
					dataD = &sliceD[i]
				}
				if sliceE != nil {
					dataE = &sliceE[i]
				}
				if sliceF != nil {
					dataF = &sliceF[i]
				}
				if sliceG != nil {
					dataG = &sliceG[i]
				}
				if sliceH != nil {
					dataH = &sliceH[i]
				}
				result := QueryResult8[A, B, C, D, E, F, G, H]{
					EntityId: entityId,
					A:        dataA,
					B:        dataB,
					C:        dataC,
					D:        dataD,
					E:        dataE,
					F:        dataF,
					G:        dataG,
					H:        dataH,
				}

				if filterFn != nil && !filterFn(result) {
					continue
				}

				if !yield(result) {
					return
				}
			}
		}
	}
}

func (query *Query8[A, B, C, D, E, F, G, H]) ForeachChannel(chunkSize int, filterFn func(QueryResult8[A, B, C, D, E, F, G, H]) bool) <-chan iter.Seq[QueryResult8[A, B, C, D, E, F, G, H]] {
	if chunkSize == 0 {
		panic("chunk size must be greater than zero")
	}

	channelsCount := math.Ceil(float64(query.Count()) / float64(chunkSize))
	channel := make(chan iter.Seq[QueryResult8[A, B, C, D, E, F, G, H]], int(channelsCount))

	go func() {
		defer close(channel)

		storageA := getStorage[A](query.World)
		storageB := getStorage[B](query.World)
		storageC := getStorage[C](query.World)
		storageD := getStorage[D](query.World)
		storageE := getStorage[E](query.World)
		storageF := getStorage[F](query.World)
		storageG := getStorage[G](query.World)
		storageH := getStorage[H](query.World)

		archetypes := query.filter()
		for _, archetype := range archetypes {
			sliceA := storageA.archetypesComponentsEntities[archetype.Id]
			sliceB := storageB.archetypesComponentsEntities[archetype.Id]
			sliceC := storageC.archetypesComponentsEntities[archetype.Id]
			sliceD := storageD.archetypesComponentsEntities[archetype.Id]
			sliceE := storageE.archetypesComponentsEntities[archetype.Id]
			sliceF := storageF.archetypesComponentsEntities[archetype.Id]
			sliceG := storageG.archetypesComponentsEntities[archetype.Id]
			sliceH := storageH.archetypesComponentsEntities[archetype.Id]

			for i := 0; i < len(archetype.entities); i += chunkSize {
				result := QueryResultChunk8[A, B, C, D, E, F, G, H]{}
				end := min(chunkSize, len(archetype.entities[i:]))

				// Set the capacity of each chunk so that appending to a chunk does
				// not modify the original slice.
				result.EntityId = archetype.entities[i : i+end : i+end]
				if sliceA != nil {
					result.A = sliceA[i : i+end : i+end]
				}
				if sliceB != nil {
					result.B = sliceB[i : i+end : i+end]
				}
				if sliceC != nil {
					result.C = sliceC[i : i+end : i+end]
				}
				if sliceD != nil {
					result.D = sliceD[i : i+end : i+end]
				}
				if sliceE != nil {
					result.E = sliceE[i : i+end : i+end]
				}
				if sliceF != nil {
					result.F = sliceF[i : i+end : i+end]
				}
				if sliceG != nil {
					result.G = sliceG[i : i+end : i+end]
				}
				if sliceH != nil {
					result.H = sliceH[i : i+end : i+end]
				}

				channel <- func(yield func(QueryResult8[A, B, C, D, E, F, G, H]) bool) {
					queryResult := QueryResult8[A, B, C, D, E, F, G, H]{}
					for k := range result.EntityId {
						queryResult.EntityId = result.EntityId[k]

						if result.A != nil {
							queryResult.A = &result.A[k]
						}
						if result.B != nil {
							queryResult.B = &result.B[k]
						}
						if result.C != nil {
							queryResult.C = &result.C[k]
						}
						if result.D != nil {
							queryResult.D = &result.D[k]
						}
						if result.E != nil {
							queryResult.E = &result.E[k]
						}
						if result.F != nil {
							queryResult.F = &result.F[k]
						}
						if result.G != nil {
							queryResult.G = &result.G[k]
						}
						if result.H != nil {
							queryResult.H = &result.H[k]
						}

						if filterFn != nil && !filterFn(queryResult) {
							continue
						}

						if !yield(queryResult) {
							return
						}
					}
				}
			}
		}
	}()

	return channel
}
