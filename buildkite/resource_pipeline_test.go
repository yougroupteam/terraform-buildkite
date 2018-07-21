package buildkite

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPipeline_basic_unknown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicUnknown,
				Check:  testAccCheckBuildkitePipelineBasicAttributesFactory("unknown"),
			},
		},
	})
}

func TestAccPipeline_basic_beanstalk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicBeanstalk,
				Check:  testAccCheckBuildkitePipelineBasicAttributesFactory("beanstalk"),
			},
		},
	})
}

func TestAccPipeline_basic_github(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicGithub,
				Check:  testAccCheckBuildkitePipelineBasicAttributesFactory("github"),
			},
		},
	})
}

func TestAccPipeline_basic_bitbucket(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicBitbucket,
				Check:  testAccCheckBuildkitePipelineBasicAttributesFactory("bitbucket"),
			},
		},
	})
}

func TestAccPipeline_basic_gitlab(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicGitlab,
				Check:  testAccCheckBuildkitePipelineBasicAttributesFactory("gitlab"),
			},
		},
	})
}

func testAccCheckBuildkitePipelineExists(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		res := new(Pipeline)

		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline ID is set")
		}

		err := client.Get([]string{"pipelines", rs.Primary.ID}, res)

		if err != nil {
			return err
		}

		if res.Slug != rs.Primary.ID {
			return fmt.Errorf("Pipeline not found")
		}

		return nil
	}
}

func testAccCheckBuildkitePipelineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_pipeline" {
			continue
		}

		res := new(Pipeline)

		err := client.Get([]string{"pipelines", rs.Primary.ID}, res)
		if err == nil {
			if res.Slug == rs.Primary.ID {
				return fmt.Errorf("Pipeline still exists")
			}
		}

		// Verify the error
		if _, ok := err.(*notFound); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckBuildkitePipelineBasicAttributesFactory(repoProvider string) resource.TestCheckFunc {
	pipelineStateId := fmt.Sprintf("buildkite_pipeline.test_%v", repoProvider)
	pipelineName := fmt.Sprintf("tf-acc-basic-%v", repoProvider)

	return resource.ComposeTestCheckFunc(
		testAccCheckBuildkitePipelineExists(pipelineStateId),
		resource.TestCheckResourceAttr(pipelineStateId, "id", pipelineName),
		resource.TestCheckResourceAttr(pipelineStateId, "slug", pipelineName),
		resource.TestCheckResourceAttr(pipelineStateId, "name", pipelineName),
		resource.TestCheckResourceAttrSet(pipelineStateId, "repository"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.#", "1"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.agent_query_rules.#", "0"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.artifact_paths", ""),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.branch_configuration", ""),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.command", "echo 'Hello World'"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.concurrency", "0"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.env.%", "0"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.name", "test"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.parallelism", "0"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.timeout_in_minutes", "0"),
		resource.TestCheckResourceAttr(pipelineStateId, "step.0.type", "script"),
		resource.TestCheckResourceAttr(pipelineStateId, "default_branch", "master"),
		resource.TestCheckResourceAttr(pipelineStateId, "branch_configuration", ""),
		resource.TestCheckResourceAttr(pipelineStateId, "description", ""),
		resource.TestCheckResourceAttr(pipelineStateId, "env.%", "0"),
		resource.TestCheckResourceAttrSet(pipelineStateId, "builds_url"),
		resource.TestCheckResourceAttrSet(pipelineStateId, "web_url"),
	)
}

const testAccPipeline_basicUnknown = `
resource "buildkite_pipeline" "test_unknown" {
  name = "tf-acc-basic-unknown"
  repository = "git@example.com:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicBeanstalk = `
resource "buildkite_pipeline" "test_beanstalk" {
  name = "tf-acc-basic-beanstalk"
  repository = "git@terraform-provider-buildkite.git.beanstalkapp.com:/terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`
const testAccPipeline_basicGithub = `
resource "buildkite_pipeline" "test_github" {
  name = "tf-acc-basic-github"
  repository = "git@github.com:saymedia/terraform-provider-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicBitbucket = `
resource "buildkite_pipeline" "test_bitbucket" {
  name = "tf-acc-basic-bitbucket"
  repository = "git@bitbucket.org:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicGitlab = `
resource "buildkite_pipeline" "test_gitlab" {
  name = "tf-acc-basic-gitlab"
  repository = "git@gitlab.com:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`
