package buildkite

import (
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
							Type: schema.TypeList,
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

func CreatePipeline(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func ReadPipeline(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func UpdatePipeline(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func DeletePipeline(d *schema.ResourceData, meta interface{}) error {
	return nil
}
