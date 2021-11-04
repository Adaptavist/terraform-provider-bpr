package bitbucket

import (
	"context"
	"fmt"
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
	// adding some timeout functionality to stop of waiting forever!
	request, err := makePipelineRequest(d, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	response, err := runPipeline(ctx, meta, *request)

	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("[DEBUG]: %s\n", *response.UUID)

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

func resourceBitbucketPipelineRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Do not try to read from the API, as the pipeline doesn't always exist
	return nil
}

func resourceBitbucketPipelineDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
