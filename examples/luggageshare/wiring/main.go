package main

import (
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/workflow/workflowspec"
	_ "github.com/blueprint-uservices/cldrel_course/luggageshare/e2etests"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/wiring/specs"
)

func main() {

	workflowspec.AddModule("github.com/blueprint-uservices/cldrel_course/luggageshare/e2etests")

	name := "luggageapp"
	cmdbuilder.MakeAndExecute(
		name,
		specs.Docker,
	)
}
