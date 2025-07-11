// Blueprint: Auto-generated by HealthChecker Plugin
package healthcheck

import (
	"context"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
)

type LuggageService_HealthChecker interface {
	AddItem(ctx context.Context, item workflow.LuggageItem) (error)
	Cleanup(ctx context.Context) (error)
	FindItems(ctx context.Context, color string, length int64, breadth int64, height int64, price float64) ([]workflow.LuggageItem, error)
	GetItemById(ctx context.Context, id string) (workflow.LuggageItem, error)
	Health(ctx context.Context) (string, error)
	
}

type LuggageService_HealthCheckHandler struct {
	Service workflow.LuggageService
}

func New_LuggageService_HealthCheckHandler(ctx context.Context, service workflow.LuggageService) (*LuggageService_HealthCheckHandler, error) {
	handler := &LuggageService_HealthCheckHandler{}
	handler.Service = service
	return handler, nil
}


func (handler *LuggageService_HealthCheckHandler) AddItem(ctx context.Context, item workflow.LuggageItem) (error) {
	return handler.Service.AddItem(ctx, item)
}

func (handler *LuggageService_HealthCheckHandler) Cleanup(ctx context.Context) (error) {
	return handler.Service.Cleanup(ctx)
}

func (handler *LuggageService_HealthCheckHandler) FindItems(ctx context.Context, color string, length int64, breadth int64, height int64, price float64) ([]workflow.LuggageItem, error) {
	return handler.Service.FindItems(ctx, color, length, breadth, height, price)
}

func (handler *LuggageService_HealthCheckHandler) GetItemById(ctx context.Context, id string) (workflow.LuggageItem, error) {
	return handler.Service.GetItemById(ctx, id)
}

func (handler *LuggageService_HealthCheckHandler) Health(ctx context.Context) (string, error) {
	return "Healthy", nil
}
