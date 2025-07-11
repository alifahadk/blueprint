// Package workerpool is a plugin that wraps the server-side of a service to limit concurrent method executions
// using a worker/thread pool.
//
// # Wiring Spec Usage
//
// To use the workerpool plugin in your wiring spec, apply it to an application-level service instance:
//
//	workerpool.Create(spec, "my_service", 10, 100)
//
// This will allow up to 10 concurrent requests, queuing up to 100 additional calls.
//
// # Description
//
// When applied, the workerpool plugin wraps the server-side implementation of a service with a worker pool
// that limits the number of concurrent method calls to N (workers). Requests beyond the queue limit will block.
// This plugin provides backpressure and concurrency control.
//
// You can continue applying other application-level modifiers after applying workerpool.
//
// # Artifacts Generated
//
// During compilation, the workerpool plugin generates a server-side wrapper. It also relies on runtime
// support code provided in [runtime/plugins/workerpool].
package workerpool

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/pointer"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"golang.org/x/exp/slog"
)

// Create wraps the server-side of a service with a workerpool limiting concurrency.
// `numWorkers` is the number of concurrent workers (goroutines),
// and `queueSize` is the number of calls that can wait in the queue.
func Create(spec wiring.WiringSpec, serviceName string, numWorkers int, queueSize int) {
	poolName := serviceName + ".workerpool"

	// Get the pointer metadata
	ptr := pointer.GetPointer(spec, serviceName)
	if ptr == nil {
		slog.Error("Unable to create workerpool for " + serviceName + " as it is not a pointer")
		return
	}

	// Add the workerpool as a destination modifier on the server side
	serverNext := ptr.AddDstModifier(spec, poolName)

	// Define the worker pool
	spec.Define(poolName, &WorkerPool{}, func(namespace wiring.Namespace) (ir.IRNode, error) {
		pool := &WorkerPool{
			PoolName:   poolName,
			MaxWorkers: numWorkers,
			QueueSize:  queueSize,
		}
		poolNamespace, err := namespace.DeriveNamespace(poolName, &serverWorkerPoolNamespace{pool})
		if err != nil {
			return nil, err
		}
		return pool, poolNamespace.Get(serverNext, &pool.Service)
	})
}

// A wiring.NamespaceHandler used to build ServerWorkerPool IRNodes
type serverWorkerPoolNamespace struct {
	*WorkerPool
}

// Implements wiring.NamespaceHandler
func (pool *WorkerPool) Accepts(nodeType any) bool {
	_, ok := nodeType.(golang.Node)
	return ok
}

// Implements wiring.NamespaceHandler
func (pool *WorkerPool) AddEdge(name string, edge ir.IRNode) error {
	pool.Edges = append(pool.Edges, edge)
	return nil
}

// Implements wiring.NamespaceHandler
func (pool *WorkerPool) AddNode(name string, node ir.IRNode) error {
	pool.Nodes = append(pool.Nodes, node)
	return nil
}
