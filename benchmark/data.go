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
