package pipelines

import (
	"fmt"
	"testing"

	"github.com/alam0rt/buildkite-operator/pkg/pipelines/mocks"
)

func TestCreate(t *testing.T) {
	mockPipeline := mocks.PipelineAL{}

	mockPipeline.On("Exists", nil).Return(nil, false)
	exists, _ := mockPipeline.Exists()
	fmt.Print(exists)

}
