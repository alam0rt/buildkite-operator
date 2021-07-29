package pipelines

import (
	"errors"

	"github.com/alam0rt/go-buildkite/v2/buildkite"
)

// PipelineAL defines what methods are required to manage pipelines
type PipelineAL interface {
	Get() (buildkite.Pipeline, error)
	Create(pipelineInput *buildkite.CreatePipeline) error
	Update(pipeline *buildkite.Pipeline) error
	Exists() (bool, error)
}

type buildkitePipeline struct {
	client       buildkite.Client
	organization string
	nameSlug     string
}

func NewBuildkitePipelineAL(organization, nameSlug, accessToken string) (PipelineAL, error) {
	config, err := buildkite.NewTokenConfig(accessToken, false)
	if err != nil {
		return nil, err
	}

	client := buildkite.NewClient(config.Client())

	pipeline := &buildkitePipeline{
		client:       *client,
		organization: organization,
		nameSlug:     nameSlug,
	}
	return pipeline, nil
}

func (p *buildkitePipeline) Get() (buildkite.Pipeline, error) {
	pipeline := &buildkite.Pipeline{}
	pipeline, httpResp, err := p.client.Pipelines.Get(p.organization, p.nameSlug)
	if err != nil {
		return *pipeline, err
	}

	if httpResp.Response.StatusCode == 404 {
		err = errors.New("pipeline not found")
		return *pipeline, err
	}

	return *pipeline, nil
}

func (p *buildkitePipeline) Update(pipeline *buildkite.Pipeline) error {
	updateResp, err := p.client.Pipelines.Update(p.organization, pipeline)
	if err != nil {
		return err
	}

	if updateResp.Response.StatusCode == 200 {
		return nil
	}

	err = errors.New("there was an unknown exception updating the pipeline")
	return err
}

func (p *buildkitePipeline) Create(pipelineInput *buildkite.CreatePipeline) error {
	_, _, err := p.client.Pipelines.Create(p.organization, pipelineInput)
	if err != nil {
		return err
	}
	return nil
}

func (p *buildkitePipeline) Exists() (bool, error) {
	_, resp, err := p.client.Pipelines.Get(p.organization, p.nameSlug)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}

	if resp.Response.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}
