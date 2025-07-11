// Blueprint: Auto-generated by HealthChecker plugin
package healthcheck

import (
	"context"
)

type ReservationService_HealthChecker interface {
	CheckAvailability(ctx context.Context, item_ids []string, startdate string, enddate string) ([]string, error)
	Cleanup(ctx context.Context) (error)
	Health(ctx context.Context) (string, error)
	MakeReservation(ctx context.Context, luggageid string, user string, startdate string, enddate string) (bool, error)
	
}
