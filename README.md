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

A Component is defined by its ComponentId, ranging between [0;2048].

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
- Create your components, and implement the ComponentInterface with GetComponentId(). Your ComponentId should range between [0;2048].
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
entityId := world.CreateEntity()
```
**Important**: the entity will receive a unique identifier. When the entity is removed, this id can be used again and assigned to a new entity.

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

## Queries
The most powerful feature is the possibility to query entities with a given set of Components.
For example, in the Rendering system of the game engine, a query will fetch only for the entities having a Mesh & Transform:
```go
query := volt.CreateQuery2[transformComponent, meshComponent](world, volt.QueryConfiguration{OptionalComponents: []volt.OptionalComponent{meshComponentId}})
for result := range query.Foreach(nil) {
    transformData(result.A)
}
```
The Foreach function receives a function to pre-filter the results, and returns an iterator.

For faster performances, you can use concurrency with the function ForeachChannel:
```go
query := volt.CreateQuery2[transformComponent, meshComponent](world, volt.QueryConfiguration{OptionalComponents: []volt.OptionalComponent{meshComponentId}})
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

## Tags
Tags are considered like any other Component internally, except they have no structure/value attached.
They cannot be fetched using functions like _GetComponent_. Due to their simpler form, they do not need to be registered.

Tags are useful to categorize your entities.

e.g. "NPC", "STATIC", "DISABLED". For example, if you want to fetch only static content, you can query through the tag "STATIC".
The Query will return only the entities tagged, in a faster way than applying the filter function in _Query.Foreach_ to check on each entities if they are static.

e.g. to fetch only static entities:
```go
const TAG_STATIC_ID = iota + volt.TAGS_INDICES
query := volt.CreateQuery2[transformComponent, meshComponent](world, volt.QueryConfiguration{Tags: []volt.TagId{TAG_STATIC_ID}})
for result := range query.Foreach(nil) {
    transformData(result.A)
}
```
Important: the TagIds should start from volt.TAGS_INDICES, allowing a range from [2048; 65535] for TagIds.

You can Add a Tag, check if an entity Has a Tag, or Remove it:
```go
world.AddTag(TAG_STATIC_ID, entityId)
world.HasTag(TAG_STATIC_ID, entityId)
world.RemoveTag(TAG_STATIC_ID, entityId)
```

## Events
The lifecycle (creation/deletion) of entities and components can trigger events.
You can configure a callback function for each of these events, to execute your custom code:
```go
world := volt.CreateWorld(100)
world.SetEntityAddedFn(func(entityId volt.EntityId) {
    fmt.Println("A new entity has been created", entityId)
})
world.SetEntityRemovedFn(func(entityId volt.EntityId) {
    fmt.Println("An entity has been deleted", entityId)
})
world.SetComponentAddedFn(func(entityId volt.EntityId, componentId volt.ComponentId) {
    fmt.Println("The component", componentId, "is attached to the entity", entityId)
})
world.SetComponentRemovedFn(func(entityId volt.EntityId, componentId volt.ComponentId) {
fmt.Println("The component", componentId, "is removed from the entity", entityId)
})
```

## Naming entities
Volt managed the naming of entities up to the version 1.6.0. For performances reasons, this feature is removed from the v1.7.0+.
You now have to keep track of the names by yourself in your application:
- Having a simple map[name string]volt.EntityId, you can react to the events and register these. Keep in mind that if your scene has a lot
of entities, it will probably have a huge impact on the garbage collector.
- Add a MetadataComponent. To fetch an entity by its name can be very slow, so you probably do not want to name all your entities. For example:
```go
const MetadataComponentId = 0

type MetadataComponent struct {
	Name string
}

func (MetadataComponent MetadataComponent) GetComponentId() volt.ComponentId {
	return MetadataComponentId
}
volt.RegisterComponent[MetadataComponent](&world, &volt.ComponentConfig[MetadataComponent]{BuilderFn: func(component any, configuration any) {}})

func GetEntityName(world *volt.World, entityId volt.EntityId) string {
    if world.HasComponents(entityId, MetadataComponentId) {
        metadata := volt.GetComponent[MetadataComponent](world, entityId)
    
        return metadata.Name
    }
    
    return ""
}

func (scene *Scene) SearchEntity(name string) volt.EntityId {
    q := volt.CreateQuery1[MetadataComponent](&scene.World, volt.QueryConfiguration{})
    for result := range q.Foreach(nil) {
        if result.A.Name == name {
            return result.EntityId
        }
    }
    
    return 0
}
```

## Benchmark
Few ECS tools exist for Go. Arche and unitoftime/ecs are probably the most looked at, and the most optimized.
In the benchmark folder, this module is compared to both of them.

- Go - v1.24.0
- Volt - v1.5.0
- [Arche - v0.15.3](https://github.com/mlange-42/arche)
- [UECS - v0.0.3](https://github.com/unitoftime/ecs)

The given results were produced by a ryzen 7 5800x, with 100.000 entities:

goos: linux
goarch: amd64
pkg: benchmark
cpu: AMD Ryzen 7 5800X 8-Core Processor             

| Benchmark                       | Iterations | ns/op     | B/op       | Allocs/op |
|---------------------------------|------------|-----------|------------|-----------|
| BenchmarkCreateEntityArche-16   | 171        | 6948273   | 11096966   | 61        |
| BenchmarkIterateArche-16        | 2704       | 426795    | 354        | 4         |
| BenchmarkAddArche-16            | 279        | 4250519   | 120089     | 100000    |
| BenchmarkRemoveArche-16         | 249        | 4821120   | 100000     | 100000    |
| BenchmarkCreateEntityUECS-16    | 34         | 37943381  | 49119549   | 200146    |
| BenchmarkIterateUECS-16         | 3885       | 287027    | 128        | 3         |
| BenchmarkAddUECS-16             | 30         | 38097927  | 4620476    | 100004    |
| BenchmarkRemoveUECS-16          | 40         | 31008811  | 3302536    | 100000    |
| BenchmarkCreateEntityVolt-16    | 49         | 27246822  | 41214216   | 200259    |
| BenchmarkIterateVolt-16         | 3651       | 329858    | 264        | 9         |
| BenchmarkIterateConcurrentlyVolt-16 | 10000      | 102732    | 3330       | 93        |
| BenchmarkAddVolt-16             | 54         | 22508281  | 4597363    | 300001    |
| BenchmarkRemoveVolt-16          | 72         | 17219355  | 400001     | 100000    |

These results show a few things:
- Arche is the fastest tool for writes operations. In our game development though we would rather lean towards fastest read operations, because the games loops will read way more often than write.
- Unitoftime/ecs is the fastest tool for read operations on one thread only, but the writes are currently way slower than Arche and Volt (except on the Create benchmark).
- Volt is a good compromise, an in-between: fast enough add/remove operations, and almost as fast as Arche and UECS for reads on one thread.
Volt uses the new iterators from go1.23, which in their current implementation are slower than using a function call in the for-loop inside the Query (as done in UECS).
This means, if the Go team finds a way to improve the performances from the iterators, we can hope to acheive near performances as UECS.
- Thanks to the iterators, Volt provides a simple way to use goroutines for read operations. The data is received through a channel of iterator.
As seen in the results, though not totally comparable, this allows way faster reading operations than any other implementation, and to use all the CPU capabilities to perform hard work on the components.
- It might be doable to use goroutines in Arche and UECS, but I could not find this feature natively? Creating chunks of the resulted slices would generate a lot of memory allocations and is not desirable.

### Other benchmarks
The creator and maintainer of Arche has published more complex benchmarks available here:
https://github.com/mlange-42/go-ecs-benchmarks

## What is to come next ?
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