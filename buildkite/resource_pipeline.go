package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipeline,
		Read:   ReadPipeline,
		Update: UpdatePipeline,
		Delete: DeletePipeline,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"web_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"builds_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"step": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"env": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"timeout_in_minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"agent_query_rules": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

type Pipeline struct {
	Id         string `json:"id,omitempty"`
	Slug       string `json:"slug,omitempty"`
	WebURL     string `json:"web_url,omitempty"`
	BuildsURL  string `json:"builds_url,omitempty"`
	Repository string `json:"repository,omitempty"`
	Name       string `json:"name,omitempty"`
	Steps      []Step `json:"steps"`
}

type Step struct {
	Type             string            `json:"type"`
	Name             string            `json:"name,omitempty"`
	Command          string            `json:"command,omitempty"`
	Environment      map[string]string `json:"env,omitempty"`
	TimeoutInMinutes int               `json:"timeout_in_minutes,omitempty"`
	AgentQueryRules  []string          `json:"agent_query_rules,omitempty"`
}

func CreatePipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	req := preparePipelineRequestPayload(d)
	res := &Pipeline{}

	err := client.Post([]string{"pipelines"}, req, res)
	if err != nil {
		return err
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func ReadPipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	slug := d.Id()

	res := &Pipeline{}

	err := client.Get([]string{"pipelines", slug}, res)
	if err != nil {
		if _, ok := err.(*notFound); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func UpdatePipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	slug := d.Id()

	req := preparePipelineRequestPayload(d)
	res := &Pipeline{}

	err := client.Patch([]string{"pipelines", slug}, req, res)
	if err != nil {
		return err
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func DeletePipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	slug := d.Id()

	return client.Delete([]string{"pipelines", slug})
}

func updatePipelineFromAPI(d *schema.ResourceData, p *Pipeline) {
	d.SetId(p.Slug)
	d.Set("id", p.Id)
	d.Set("name", p.Name)
	d.Set("repository", p.Repository)
	d.Set("web_url", p.WebURL)
	d.Set("slug", p.Slug)
	d.Set("builds_url", p.BuildsURL)
}

func preparePipelineRequestPayload(d *schema.ResourceData) *Pipeline {
	req := &Pipeline{}

	req.Name = d.Get("name").(string)
	req.Slug = d.Get("slug").(string)
	req.Repository = d.Get("repository").(string)
	stepsI := d.Get("step").([]interface{})
	req.Steps = make([]Step, len(stepsI))

	for i, stepI := range stepsI {
		stepM := stepI.(map[string]interface{})
		req.Steps[i] = Step{
			Type:            stepM["type"].(string),
			Name:            stepM["name"].(string),
			Command:         stepM["command"].(string),
			Environment:     map[string]string{},
			AgentQueryRules: make([]string, len(stepM["agent_query_rules"].([]interface{}))),
		}

		for k, vI := range stepM["env"].(map[string]interface{}) {
			req.Steps[i].Environment[k] = vI.(string)
		}

		for j, vI := range stepM["agent_query_rules"].([]interface{}) {
			req.Steps[i].AgentQueryRules[j] = vI.(string)
		}
	}

	return req
}
