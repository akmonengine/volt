package volt

import (
	"testing"
)

func Test_getStorage(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})

	s1 := getStorage[testTransform](world)
	if s1 == nil {
		t.Errorf("RegisterComponent did not generate storage for ComponentId %d", testTransformId)
	}
	s2 := getStorage[testTag](world)
	if s2 != nil {
		t.Errorf("storage found for not registered ComponentId %d", testTagId)
	}
}

func TestWorld_getStorageForComponentId(t *testing.T) {
	world := CreateWorld(1024)
	RegisterComponent[testTransform](world, &ComponentConfig[testTransform]{ID: testTransformId})

	s, err := world.getStorageForComponentId(testTransformId)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s == nil {
		t.Errorf("getStorageForComponentId() did not find storage for %d", testTransformId)
	}
	s, err = world.getStorageForComponentId(testTagId)
	if err == nil || s != nil {
		t.Errorf("getStorageForComponentId() returned a storage for the not registered ComponentId %d", testTransformId)
	}
}
