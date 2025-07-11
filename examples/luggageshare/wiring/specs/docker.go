package specs

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/clientpool"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
	"github.com/blueprint-uservices/blueprint/plugins/gotests"
	"github.com/blueprint-uservices/blueprint/plugins/healthchecker"
	"github.com/blueprint-uservices/blueprint/plugins/hello"
	"github.com/blueprint-uservices/blueprint/plugins/http"
	"github.com/blueprint-uservices/blueprint/plugins/linuxcontainer"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/mysql"
	wf "github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
)

var Docker = cmdbuilder.SpecOption{
	Name:        "docker",
	Description: "Deploys each service in a separate container with http, uses mongodb as NoSQL database backends, and applies a number of modifiers",
	Build:       makeDockerSpec,
}

func makeDockerSpec(spec wiring.WiringSpec) ([]string, error) {
	applyDockerDefaults := func(spec wiring.WiringSpec, serviceName string) string {
		hello.Instrument(spec, serviceName)
		healthchecker.AddHealthCheckAPI(spec, serviceName)
		http.Deploy(spec, serviceName)
		goproc.Deploy(spec, serviceName)
		return linuxcontainer.Deploy(spec, serviceName)
	}
	var cntrs []string
	var allServices []string

	luggage_db := mysql.Container(spec, "luggage_db")
	user_db := mongodb.Container(spec, "user_db")
	review_db := mongodb.Container(spec, "review_db")
	reserv_db := mongodb.Container(spec, "reserv_db")

	luggage_service := wf.Service[workflow.LuggageService](spec, "luggage_service", luggage_db)
	clientpool.Create(spec, luggage_service, 1)
	luggage_cntr := applyDockerDefaults(spec, luggage_service)
	cntrs = append(cntrs, luggage_cntr)
	allServices = append(allServices, luggage_service)

	user_service := wf.Service[workflow.UserProfileService](spec, "user_service", user_db)
	user_cntr := applyDockerDefaults(spec, user_service)
	cntrs = append(cntrs, user_cntr)
	allServices = append(allServices, user_service)

	review_service := wf.Service[workflow.ReviewService](spec, "review_service", review_db)
	review_cntr := applyDockerDefaults(spec, review_service)
	cntrs = append(cntrs, review_cntr)
	allServices = append(allServices, review_service)

	reserv_service := wf.Service[workflow.ReservationService](spec, "reserv_service", reserv_db)
	reserv_cntr := applyDockerDefaults(spec, reserv_service)
	cntrs = append(cntrs, reserv_cntr)
	allServices = append(allServices, reserv_service)

	search_service := wf.Service[workflow.SearchService](spec, "search_service", reserv_service, luggage_service)
	search_cntr := applyDockerDefaults(spec, search_service)
	cntrs = append(cntrs, search_cntr)
	allServices = append(allServices, search_service)

	frontend_service := wf.Service[workflow.FrontendService](spec, "frontend_service", search_service, reserv_service, review_service, user_service, luggage_service)
	frontend_cntr := applyDockerDefaults(spec, frontend_service)
	cntrs = append(cntrs, frontend_cntr)
	allServices = append(allServices, frontend_service)

	tests := gotests.Test(spec, allServices...)
	cntrs = append(cntrs, tests)

	return cntrs, nil
}
