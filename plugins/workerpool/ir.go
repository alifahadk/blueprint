package workerpool

import (
	"fmt"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint/stringutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/service"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gocode"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gogen"
	"log/slog"
	"path/filepath"
	"strings"
)

type WorkerPool struct {
	golang.Service
	golang.GeneratesFuncs

	PoolName   string
	QueueSize  int
	MaxWorkers int
	Server     golang.Service
	Edges      []ir.IRNode
	Nodes      []ir.IRNode
}

func (pool *WorkerPool) Name() string {
	return pool.PoolName
}

func (pool *WorkerPool) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%v = ServerWorkerPool(%v) {\n", pool.PoolName, pool.Server.Name()))
	var children []string
	for _, child := range pool.Nodes {
		children = append(children, child.String())
	}
	b.WriteString(stringutil.Indent(strings.Join(children, "\n"), 2))
	b.WriteString("\n}")
	return b.String()
}

func (pool *WorkerPool) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	// WorkerPool doesn't modify the server interface, just wraps calls
	return pool.Server.GetInterface(ctx)
}

// Implements golang.Service golang.ProvidesInterface
func (pool *WorkerPool) AddInterfaces(module golang.ModuleBuilder) error {
	/* WorkerPool doesn't modify the client's interface and doesn't introduce new interfaces */
	for _, node := range pool.Nodes {
		if n, valid := node.(golang.ProvidesInterface); valid {
			if err := n.AddInterfaces(module); err != nil {
				return err
			}
		}
	}
	return nil
}

// Implements golang.GeneratesFuncs
func (pool *WorkerPool) GenerateFuncs(module golang.ModuleBuilder) error {
	// 1. Avoid duplicate generation for this service
	iface, err := golang.GetGoInterface(module, pool.Service)
	if err != nil {
		return err
	}
	if module.Visited(iface.Name + "_ServerWorkerPool") {
		return nil
	}

	// 2. Recursively generate code for all inner nodes
	for _, node := range pool.Nodes {
		if n, valid := node.(golang.GeneratesFuncs); valid {
			if err := n.GenerateFuncs(module); err != nil {
				return err
			}
		}
	}

	// 3. Prepare template arguments
	args, err := pool.getTemplateArgs(module)
	if err != nil {
		return err
	}

	// 4. Build namespace file for instantiation logic
	nsBuilder, err := gogen.NewNamespaceBuilder(
		module,
		args.Service.BaseName+"_ServerWorkerPool",
		args.PackageShortName,
		args.PackageShortName,
		args.ConstructorName,
	)
	if err != nil {
		return err
	}

	// Add instantiation for each child node
	for _, node := range pool.Nodes {
		if inst, ok := node.(golang.Instantiable); ok {
			if err := inst.AddInstantiation(nsBuilder); err != nil {
				return err
			}
		}
	}

	// Build the namespace file
	if err := nsBuilder.Build(); err != nil {
		return err
	}

	// 5. Generate the actual server wrapper Go file using the embedded template
	fileName := filepath.Join(module.Info().Path, args.PackageShortName, args.FileName)
	return gogen.ExecuteTemplateToFile("serverworkerpool", workerPoolTemplate, args, fileName)
}

// Implements golang.Service and golang.Instantiable
func (pool *WorkerPool) AddInstantiation(builder golang.NamespaceBuilder) error {
	if builder.Visited(pool.PoolName) {
		return nil
	}

	args, err := pool.getTemplateArgs(builder.Module())
	if err != nil {
		return err
	}

	builder.Import(args.PackageName)

	slog.Info(fmt.Sprintf("Instantiating ServerWorkerPool %v in %v/%v", pool.PoolName, builder.Info().Package.PackageName, builder.Info().FileName))

	code, err := gogen.ExecuteTemplate("serverworkerpool_instantiation", workerPoolBuildTemplate, args)
	if err != nil {
		return err
	}

	return builder.Declare(pool.PoolName, code)
}

func (pool *WorkerPool) getTemplateArgs(module golang.ModuleBuilder) (*templateArgs, error) {
	var err error
	args := &templateArgs{}

	args.Service, err = golang.GetGoInterface(module, pool.Service)
	if err != nil {
		return nil, err
	}

	args.WrappedService = pool.Service.Name()
	args.InstanceName = pool.PoolName
	args.MaxWorkers = pool.MaxWorkers
	args.QueueSize = pool.QueueSize

	args.PoolName = args.Service.Name + "_ServerWorkerPool"
	args.PackageShortName = "workerpool"
	args.PackageName = module.Info().Name + "/" + args.PackageShortName

	args.FileName = args.Service.BaseName + "_workerpool.go"
	args.ConstructorName = fmt.Sprintf("New_%v_ServerWorkerPool", args.Service.BaseName)
	args.Imports = gogen.NewImports(args.PackageName)

	args.Imports.AddPackages(
		"context", "sync", "fmt", "log",
		"github.com/blueprint-uservices/blueprint/runtime/plugins/workerpool",
		"github.com/blueprint-uservices/blueprint/runtime/plugins/golang",
	)

	return args, nil
}

type templateArgs struct {
	InstanceName     string
	WrappedService   string
	MaxWorkers       int
	QueueSize        int
	PoolName         string
	PackageShortName string
	PackageName      string
	FileName         string
	ConstructorName  string
	Service          *gocode.ServiceInterface
	Imports          *gogen.Imports
}

const workerPoolTemplate = `
package {{ .PackageShortName }}

{{ .Imports }}

type {{ .PoolName }} struct {
	target {{ .Service.Name }}
	queue  chan workRequest
	wg     sync.WaitGroup
}

type workRequest struct {
	fn   func()
	done chan struct{}
}

func {{ .ConstructorName }}(svc {{ .Service.Name }}) {{ .Service.Name }} {
	pool := &{{ .PoolName }}{
		target: svc,
		queue:  make(chan workRequest, {{ .QueueSize }}),
	}

	for i := 0; i < {{ .MaxWorkers }}; i++ {
		go pool.worker()
	}

	return pool
}

func (p *{{ .PoolName }}) worker() {
	for req := range p.queue {
		req.fn()
		close(req.done)
	}
}

{{ range .Service.Methods }}
func (p *{{ $.PoolName }}) {{ .Name }}({{ .ArgsStr }}) {{ .ReturnStr }} {
	done := make(chan struct{})
	var ret struct {
		{{- range .ReturnVars }}
		{{ .Name }} {{ .Type }}
		{{- end }}
	}

	p.queue <- workRequest{
		fn: func() {
			{{- if .HasReturn }}
			{{ if .HasNamedReturns }}
			{{- range .ReturnVars }}
			ret.{{ .Name }} =
			{{- end }}
			{{ else }}
			res := 
			{{ end -}}
			p.target.{{ .Name }}({{ .CallArgsStr }})
			{{ end }}
		},
		done: done,
	}

	<-done

	{{- if .HasReturn }}
	{{ if .HasNamedReturns }}
	return {{ range .ReturnVars }}ret.{{ .Name }}, {{ end }}
	{{ else }}
	return res
	{{ end }}
	{{- end }}
}
{{ end }}
`

const workerPoolBuildTemplate = `
{{ .InstanceName }} := {{ .ConstructorName }}({{ .WrappedService }})
`
