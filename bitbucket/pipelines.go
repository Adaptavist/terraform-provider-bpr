package bitbucket

import (
	"context"
	"fmt"
	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func makePipelineRequest(d *schema.ResourceData, meta interface{}) (*model.PostPipelineRequest, error) {
	cli := meta.(*client.Client)
	pipeline := builders.Pipeline()

	for key, value := range d.Get("variables").(map[string]interface{}) {
		pipeline.Variable(key, fmt.Sprintf("%v", value), false)
	}

	target := builders.Target()

	if pattern := d.Get("pipeline").(string); pattern != "" {
		target.Pattern(pattern)
	}

	if tag := d.Get("tag").(string); tag != "" {
		// retrieve the tag from Bitbucket as we'll need the commit hash
		tagResponse, err := cli.GetTag(model.GetTagRequest{
			Workspace:  valueAsPointer("workspace", d),
			Repository: valueAsPointer("repository", d),
			Tag:        tag,
		})

		if err != nil {
			return nil, err
		}

		target.Tag(tag, tagResponse.Target.Hash)
	}

	if branch := d.Get("branch").(string); branch != "" {
		target.Branch(branch)
	}

	pipeline.Target(target.Build())

	return &model.PostPipelineRequest{
		Workspace:  valueAsPointer("workspace", d),
		Repository: valueAsPointer("repository", d),
		Pipeline:   pipeline.Build(),
	}, nil
}

type RunPipelineResult struct {
	error error
	pipeline model.Pipeline
}

// runPipeline will POST the pipeline but will monitor its status until completion.
func runPipeline(ctx context.Context, meta interface{}, request model.PostPipelineRequest) (*model.Pipeline, error) {
	cli := meta.(*client.Client)

	result := make(chan RunPipelineResult, 1)

	run := func(ctx context.Context) {
		response, err := cli.PostPipeline(request)

		if err != nil {
			result <- RunPipelineResult{error: err}
			return
		}

		for {
			// has the pipeline completed
			if response.CompletedOn != nil {
				result <- RunPipelineResult{pipeline: *response}
				return
			}

			time.Sleep(time.Second * 5) // TODO add configuration

			response, err = cli.GetPipeline(model.GetPipelineRequest{
				Workspace:  request.Workspace,
				Repository: request.Repository,
				Pipeline:   response,
			})

			if err != nil {
				result <- RunPipelineResult{pipeline: *response}
				return
			}

			if response.State.Result.HasError() {
				result <- RunPipelineResult{
					fmt.Errorf("pipeline finished with error %v", response.State.Result.Error),
					*response,
				}
				return
			}
		}
	}

	go run(ctx)

	select {
		case p := <-result:
			fmt.Println("[DEBUG]: returning result")
			return &p.pipeline, p.error
		case <-ctx.Done():
			fmt.Println("[DEBUG]: done")
			return nil, ctx.Err()
	}
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
