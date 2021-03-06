package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/windmilleng/tilt/internal/container"
	"github.com/windmilleng/tilt/pkg/model"
)

func imageID(s string) model.TargetID {
	return model.TargetID{
		Type: model.TargetTypeImage,
		Name: model.TargetName(s),
	}
}

func TestOneAndOnlyLiveUpdatedContainerID(t *testing.T) {
	set := BuildResultSet{
		imageID("a"): NewLiveUpdateBuildResult(imageID("a"), nil, []container.ID{"cA"}),
		imageID("b"): NewLiveUpdateBuildResult(imageID("b"), nil, []container.ID{"cB"}),
	}
	assert.Equal(t, "", string(set.OneAndOnlyLiveUpdatedContainerID()))

	set = BuildResultSet{
		imageID("a"): NewLiveUpdateBuildResult(imageID("a"), nil, []container.ID{"cA"}),
		imageID("b"): NewLiveUpdateBuildResult(imageID("b"), nil, []container.ID{"cA"}),
		imageID("c"): NewLiveUpdateBuildResult(imageID("c"), nil, []container.ID{""}),
	}
	assert.Equal(t, "cA", string(set.OneAndOnlyLiveUpdatedContainerID()))
}
