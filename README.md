# Volt
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/akmonengine/volt)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/akmonengine/volt)
[![Go Report Card](https://goreportcard.com/badge/github.com/akmonengine/volt)](https://goreportcard.com/report/github.com/akmonengine/volt)
![Tests](https://img.shields.io/github/actions/workflow/status/akmonengine/volt/code_coverage.yml?label=tests)
![Codecov](https://img.shields.io/codecov/c/github/akmonengine/volt)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/akmonengine/volt)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/akmonengine/volt)


Volt is an ECS(entity-component-system) oriented for games development with Go.
It is inspired by the documentation available here: https://github.com/SanderMertens/ecs-faq

There is many ways to write an ECS, and Volt is based on the Archetype paradigm.

## Knowledge
### Entity
An entity is the end object in a game (e.g. a character). It is only defined by
its identifier called EntityId. This identifier is randomly generated, its type uint64 avoiding to generate twice the same id.
It is also required to set a name for each entity, only used to easily retrieve them when required.

Looking at the benchmark, a scene can handle between 100.000 to 1.000.000 depending on your machine and the complexity of the project.
But of course, the lower the better, as it will allow the project to run on slower computers.

### Component
An entity is composed from 1 to N Component(s).
It is a structure of properties, and should not contain any logic by itself (meaning no functions).
The Components are manipulated by Systems.

A Component is defined by its ComponentId.

### System
A system is a specialized tool that fetches entities, filtered by their Components, and transforms the datas.
For example: the audio could be managed by a system, or the graphics managed by a render system.

Volt does not directly implements Systems, but allows you to create Queries that you can use in your own specific tools.

### Query
A Query is a search tool for the set of entities that possess (at least) the list of ComponentId provided.
It is then possible to iterate over the result of this search within a System, in order to manipulate the Components.

### Archetype
In an ECS (Entity-Component-System), an Archetype is the set of Entities that share the same ComponentId.
The Archetype itself is not publicly exposed, but is instead managed internally and represents a major structure within Volt.

### Storage
Using the Structure Of Arrays (SoA) paradigm, Components are persisted in a dedicated storage for each ComponentId. This allows for cache hits during read phases within Query iterations, resulting in significantly improved performance compared to an Array of Structures (AoS) model.

## Basic Usage

- Create a World to contain all the datas
```go 
world := volt.CreateWorld()
```
- Create your components, and implement the ComponentInterface with GetComponentId()
```go 
const (
    transformComponentId = iota
)

type transformComponent struct {
x, y, z float64
}

func (t transformComponent) GetComponentId() volt.ComponentId {
    return transformComponentId
}
type transformConfiguration struct {
	x, y, z float64
}
```
- Register the component for the world to use it. The BuilderFn is optional, it allows to initialize and customize the component at its creation.
```go 
volt.RegisterComponent[transformComponent](world, &ComponentConfig[transformComponent]{BuilderFn: func(component any, configuration any) {
    conf := configuration.(*transformConfiguration)
    transformComponent := component.(*transformComponent)
	transformComponent.x = conf.x
	transformComponent.y = conf.y
	transformComponent.z = conf.z
}})
```
- Create the entity
```go 
entityId := world.CreateEntity("entityName")
```
- Add the component to the entity
```go 
component := volt.ConfigureComponent[transformComponent](&scene.World, transformConfiguration{x: 1.0, y: 2.0, z: 3.0})
volt.AddComponent(&scene.World, entity, component)
```
- Remove the component to the entity
```go
err := RemoveComponent[testTransform](world, entityId)
if err != nil {
	fmt.Println(err)
}
```
- Delete the entity
```go
world.RemoveEntity(entityId)
```
## Searching for an entity
- Knowing an entity by its name, you can get its identifier:
```go
entityId := world.SearchEntity("entityName")
```
- The reversed search is also possible, fetching its name by its idenfier:
```go
entityName := world.GetEntityName(entityId)
```

## Queries
The most powerful feature is the possibility to query entities with a given set of Components.
For example, in the Rendering system of the game engine, a query will fetch only for the entities having a Mesh & Transform:
```go
query := volt.CreateQuery2[transformComponent, meshComponent](world, []volt.OptionalComponent{meshComponentId})
for result := range query.Foreach(nil) {
    transformData(result.A)
}
```
The Foreach function receives a function to pre-filter the results, and returns an iterator.

For faster performances, you can use concurrency with the function ForeachChannel:
```go
query := volt.CreateQuery2[transformComponent, meshComponent](world, []volt.OptionalComponent{meshComponent})
queryChannel := query.ForeachChannel(1000, nil)

runWorkers(4, func(workerId int) {
    for results := range queryChannel {
        for result := range results {
            transformData(result.A)
        }
    }
})

func runWorkers(workersNumber int, worker func(int)) {
  var wg sync.WaitGroup
  
  for i := range workersNumber {
    i := i
    wg.Add(1)
    go func(worker func(int)) {
        defer wg.Done()
        worker(i)
    }(worker)
  }
  
  wg.Wait()
}
```

Queries exist for 1 to 8 Components.

You can also get the number of entities, without looping on each:
```go
total := query.Count()
```
Or get the entities identifiers as a slice:
```go
entitiesIds := query.FetchAll()
```

## Benchmark
Few ECS tools exist for Go. Arche and unitoftime/ecs are probably the most looked at, and the most optimized.
In the benchmark folder, this module is compared to both of them.

- Go - v1.23.1
- Volt - v1.1.0
- [Arche - v0.13.2](https://github.com/mlange-42/arche)
- [UECS - v0.0.2-0.20240727195554-03fbb2d998cf](https://github.com/unitoftime/ecs)

The given results were produced by a ryzen 7 5800x, with 100.000 entities:

| Benchmark feature (entities count)                                | Time/Operation  | Bytes/Operation | Allocations/Operation |
|-------------------------------------------------------------------|-----------------|-----------------|-----------------------|
| BenchmarkCreateEntityVolt (100000)                                | 28019814 ns/op  | 27509201 B/op   | 203267 allocs/op      |
| BenchmarkCreateEntityArche (100000)                               | 13090540 ns/op  | 43679965 B/op   | 1631 allocs/op        |
| BenchmarkCreateEntityUECS (100000)                                | 33813923 ns/op  | 49120107 B/op   | 200148 allocs/op      |
| BenchmarkIterateVolt (100000)                                     | 327920 ns/op    | 128 B/op        | 5 allocs/op           |
| BenchmarkIterateConcurrentlyVolt (100000) - 16 concurrent workers | 95396 ns/op     | 3274 B/op       | 93 allocs/op          |
| BenchmarkIterateArche (100000)                                    | 302350 ns/op    | 354 B/op        | 4 allocs/op           |
| BenchmarkIterateUECS (100000)                                     | 234814 ns/op    | 152 B/op        | 4 allocs/op           |
| BenchmarkAddVolt (100000)                                         | 27868065 ns/op  | 4750176 B/op    | 300002 allocs/op      |
| BenchmarkAddArche (100000)                                        | 6204147 ns/op   | 3329050 B/op    | 200000 allocs/op      |
| BenchmarkAddUECS (100000)                                         | 71806194 ns/op  | 16643870 B/op   | 500010 allocs/op      |
| BenchmarkRemoveVolt (100000)                                      | 19352594 ns/op  | 200000 B/op     | 100000 allocs/op      |
| BenchmarkRemoveArche (100000)                                     | 6796755 ns/op   | 3300006 B/op    | 200000 allocs/op      |
| BenchmarkRemoveUECS (100000)                                      | 76231994 ns/op  | 17092955 B/op   | 600002 allocs/op      |

These results show a few things:
- Arche is the fastest tool for writes operations. In our game development though we would rather lean towards fastest read operations, because the games loops will read way more often than write.
- Unitoftime/ecs is the fastest tool for read operations on one thread only, but the writes are currently way slower than Arche and Volt (except on the Create benchmark).
- Volt is a good compromise, an in-between: fast enough add/remove operations, and almost as fast as Arche and UECS for reads on one thread.
Volt uses the new iterators from go1.23, which in their current implementation are slower than using a function call in the for-loop inside the Query (as done in UECS). See https://github.com/golang/go/issues/69411.
This means, if the Go team finds a way to improve the performances from the iterators, we can hope to acheive near performances as UECS.
- Thanks to the iterators, Volt provides a simple way to use goroutines for read operations. The data is received through a channel of iterator.
As seen in the results, though not totally comparable, this allows way faster reading operations than any other implementation, and to use all the CPU capabilities to perform hard work on the components.
- It might be doable to use goroutines in Arche and UECS, but I could not find this feature natively? Creating chunks of the resulted slices would generate a lot of memory allocations and is not desirable.

## What is to come next ?
- Hopefully the ticket https://github.com/golang/go/issues/69411 could have a huge impact to boost the performances for queries.
- Tags (zero sized types) are useful to query entities with specific features: for example, in a renderer, to get only the entities with the boolean isCulled == false.
This would hugely reduce the loops operations in some scenarios. Currently we can use the filters on the iterators, but it does not avoid the fact that every entity (with the given components) is looped by the renderer.
- For now the system is not designed to manage writes on a concurrent way: it means it is not safe to add/remove components in queries 
  using multiples threads/goroutines. I need to figure out how to implement this, though I never met the need for this feature myself.

## Sources
- https://github.com/SanderMertens/ecs-faq
- https://skypjack.github.io/2019-02-14-ecs-baf-part-1/
- https://ajmmertens.medium.com/building-an-ecs-1-where-are-my-entities-and-components-63d07c7da742
- https://github.com/unitoftime/ecs

## Contributing Guidelines

See [how to contribute](CONTRIBUTING.md).

## Licence
This project is distributed under the [Apache 2.0 licence](LICENCE.md).