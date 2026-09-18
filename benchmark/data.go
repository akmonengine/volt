package benchmark

import (
	"github.com/akmonengine/volt"
	"sync"
)

const ENTITIES_COUNT = 100000
const WORKERS = 16

const (
	testTransformId = iota
	testTagId
)

type testTransform struct {
	x, y, z float64
}

func (t testTransform) GetComponentId() volt.ComponentId {
	return testTransformId
}

type testTag struct {
	x, y, z int
}

func (t testTag) GetComponentId() volt.ComponentId {
	return testTagId
}

// Six more small components, to build "large" entities (8 components) like the
// create-large-entities scenario of go-ecs-benchmarks.
const (
	testC3Id = iota + 2
	testC4Id
	testC5Id
	testC6Id
	testC7Id
	testC8Id
)

type testC3 struct{ v float64 }
type testC4 struct{ v float64 }
type testC5 struct{ v float64 }
type testC6 struct{ v float64 }
type testC7 struct{ v float64 }
type testC8 struct{ v float64 }

func (testC3) GetComponentId() volt.ComponentId { return testC3Id }
func (testC4) GetComponentId() volt.ComponentId { return testC4Id }
func (testC5) GetComponentId() volt.ComponentId { return testC5Id }
func (testC6) GetComponentId() volt.ComponentId { return testC6Id }
func (testC7) GetComponentId() volt.ComponentId { return testC7Id }
func (testC8) GetComponentId() volt.ComponentId { return testC8Id }

func transformData(tr *testTransform) {
	tr.x += 1.0
	tr.y += 2.0
	tr.z += 3.0
}

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
