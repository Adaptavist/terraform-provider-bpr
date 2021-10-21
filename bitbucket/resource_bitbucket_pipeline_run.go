package bitbucket

import (
	"context"
	"fmt"
	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jeremywohl/flatten"
)

func resourceBitbucketPipelineRun() *schema.Resource {
	return &schema.Resource{
		Description:   "A Bitbucket pipeline to be invoked",
		CreateContext: resourceBitbucketPipelineInvoke,
		ReadContext:   resourceBitbucketPipelineRead,
		UpdateContext: resourceBitbucketPipelineInvoke,
		DeleteContext: resourceBitbucketPipelineDelete,
		Schema:        resourceBitbucketPipelineSchema(),
	}
}

func resourceBitbucketPipelineSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"workspace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The owner of the target pipelines repository",
		},
		"repository": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The repository containing the target pipeline",
		},
		"tag": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"branch"},
			Description:   "The tag containing the target pipeline",
		},
		"branch": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"tag"},
			Description:   "The branch containing the target pipeline",
		},
		"pipeline": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "default",
			Description: "Name of custom pipeline to run",
		},
		"variables": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Map of variables for the pipeline",
		},
		"output": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The full pipeline output",
		},
		"outputs": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Map of outputs, if we can grab them from the pipeline text output",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The UUID of the pipeline, use for looking up its status",
		},
	}
}

func resourceBitbucketPipelineInvoke(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.FromErr(err)
		}

		target.Tag(tag, tagResponse.Target.Hash)
	}

	if branch := d.Get("branch").(string); branch != "" {
		target.Branch(branch)
	}

	pipeline.Target(target.Build())

	request := model.PostPipelineRequest{
		Workspace:  valueAsPointer("workspace", d),
		Repository: valueAsPointer("repository", d),
		Pipeline:   pipeline.Build(),
	}

	response, err := cli.RunPipeline(request)

	if err != nil {
		return diag.FromErr(err)
	}

	steps, err := getPipelineSteps(d, meta, response)

	if err != nil {
		return diag.FromErr(err)
	}

	// extract the pipeline output and its containing outputs

	output := ""
	outputs := map[string]interface{}{}

	for _, step := range steps {
		log, err := getPipelineStepLog(d, meta, response, &step)

		if err != nil {
			return diag.FromErr(err)
		}

		output = output + "\n" + log
		logData, err := extractOutputs(log)

		if err != nil {
			return diag.FromErr(err)
		}

		for k, v := range logData {
			outputs[k] = v
		}
	}

	// flatten outputs

	flat, err := flatten.Flatten(outputs, "", flatten.DotStyle)

	if err != nil {
		return diag.FromErr(err)
	}

	// set the pipeline outputs (extracted from the pipeline)

	err = d.Set("outputs", flat)

	if err != nil {
		return diag.FromErr(err)
	}

	// set the pipeline output

	err = d.Set("output", output)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*response.UUID)

	return resourceBitbucketPipelineRead(ctx, d, meta)
}

func resourceBitbucketPipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Do not try to read from the API, as the pipeline doesn't always exist
	return nil
}

func resourceBitbucketPipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
