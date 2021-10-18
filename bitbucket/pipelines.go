package bitbucket

import (
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPipeline(d *schema.ResourceData, meta interface{}) (*model.Pipeline, error) {
	cli := meta.(*client.Client)
	request := model.GetPipelineRequest{
		Workspace:  valueAsPointer("workspace", d),
		Repository: valueAsPointer("repository", d),
		Pipeline: &model.Pipeline{
			UUID: valueAsPointer("id", d),
		},
	}
	// We only really care about getting the outputs from the pipeline
	return cli.GetPipeline(request)
}

func getPipelineSteps(d *schema.ResourceData, meta interface{}, pipeline *model.Pipeline) (model.PipelineSteps, error) {
	cli := meta.(*client.Client)
	return cli.GetPipelineSteps(model.GetPipelineRequest{
		Workspace:  valueAsPointer("workspace", d),
		Repository: valueAsPointer("repository", d),
		Pipeline:   pipeline,
	})
}

func getPipelineStepLog(d *schema.ResourceData, meta interface{}, pipeline *model.Pipeline, step *model.PipelineStep) (string, error) {
	cli := meta.(*client.Client)
	log, err := cli.GetPipelineStepLog(model.GetPipelineStepRequest{
		Workspace:    valueAsPointer("workspace", d),
		Repository:   valueAsPointer("repository", d),
		Pipeline:     pipeline,
		PipelineStep: step,
	})
	return string(log), err
}
