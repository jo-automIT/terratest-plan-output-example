package testing

import (
	"fmt"
	"testing"
	"reflect"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the simple Terraform module in examples/terraform-basic-example using Terratest.
func TestTerraformBasicExample(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/terraform-basic-example",

		// Variables to pass to our Terraform code using -var-file options
		VarFiles: []string{"varfile.tfvars"},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, terraformOptions)

	assert.Equal(t, "test", PlanOutput(plan, "example"))
	assert.Equal(t, "test2", PlanOutput(plan, "example2"))
	assert.Equal(t, []string{"test"}, PlanOutputList(plan, "example_list"))
	assert.Equal(t, map[string]interface{}{"expected": "test"}, PlanOutputMap(plan, "example_map"))
}

func PlanOutput(p *terraform.PlanStruct, key string) interface{} {
	return p.RawPlan.PlannedValues.Outputs[key].Value
}

func PlanOutputList(p *terraform.PlanStruct, key string) []string {
	planOutRaw := p.RawPlan.PlannedValues.Outputs[key].Value
	var out []string
	switch reflect.TypeOf(planOutRaw).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(planOutRaw)

		for i := 0; i < s.Len(); i++ {
			out = append(out, s.Index(i).Elem().String())
		}
	}
	return out
}

func PlanOutputMap(p *terraform.PlanStruct, key string) map[string]interface{} {
	planOutRaw := p.RawPlan.PlannedValues.Outputs[key].Value
	fmt.Println(planOutRaw)
	switch reflect.TypeOf(planOutRaw).Kind() {
	case reflect.Map:
		s := reflect.ValueOf(planOutRaw)
		i := s.Interface()
		out := i.(map[string]interface{})
		return out
	}
	return nil
}
