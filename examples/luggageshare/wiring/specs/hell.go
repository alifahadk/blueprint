package specs

/*import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/gotests"
	"github.com/blueprint-uservices/blueprint/plugins/hello"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/mysql"
	wf "github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
)

var hell = cmdbuilder.SpecOption{
	Name:        "hell",
	Description: "adds log instrumentation to each service",
	Build:       makeHelloSpec,
}

func makeHelloSpec(spec wiring.WiringSpec) ([]string, error) {
	applyHelloDefaults := func(spec wiring.WiringSpec, serviceName string) {
		hello.Instrument(spec, serviceName)
	}

	luggage_db := mysql.Container(spec, "luggage_db")
	user_db := mongodb.Container(spec, "user_db")
	review_db := mongodb.Container(spec, "review_db")
	reserv_db := mongodb.Container(spec, "reserv_db")

	luggage_service := wf.Service[workflow.LuggageService](spec, "luggage_service", luggage_db)
	applyHelloDefaults(spec, luggage_service)

	user_service := wf.Service[workflow.UserProfileService](spec, "user_service", user_db)
	applyHelloDefaults(spec, user_service)

	review_service := wf.Service[workflow.ReviewService](spec, "review_service", review_db)
	applyHelloDefaults(spec, review_service)

	reserv_service := wf.Service[workflow.ReservationService](spec, "reserv_service", reserv_db)
	applyHelloDefaults(spec, reserv_service)

	search_service := wf.Service[workflow.SearchService](spec, "search_service", reserv_service, luggage_service)
	applyHelloDefaults(spec, search_service)

	frontend_service := wf.Service[workflow.FrontendService](spec, "frontend_service", search_service, reserv_service, review_service, user_service, luggage_service)
	applyHelloDefaults(spec, frontend_service)

	tests := gotests.Test(spec, allServices...)
	cntrs = append(cntrs, tests)

	return cntrs, nil
}*/
