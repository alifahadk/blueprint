package tests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/simplenosqldb"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/require"
)

var reservationServiceRegistry = registry.NewServiceRegistry[workflow.ReservationService]("reservation_service")

func init() {
	reservationServiceRegistry.Register("local", func(ctx context.Context) (workflow.ReservationService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewReservationServiceImpl(ctx, db)
	})
}

func TestReservationServiceMakeReservation(t *testing.T) {
	ctx := context.Background()
	service, err := reservationServiceRegistry.Get(ctx)
	require.NoError(t, err)

	user := "test_user"
	luggageid := "test_luggage"
	luggageid_alt := "another_test_luggage"

	startdate := "2001-01-01"
	enddate := "2001-01-03"

	startdate_reserved := "2001-01-02"
	enddate_reserved := "2001-01-10"

	stardate_free := "2001-03-01"
	enddate_free := "2001-04-01"

	// Make reservation
	{
		status, err := service.MakeReservation(ctx, luggageid, user, startdate, enddate)
		require.NoError(t, err)
		require.Equal(t, true, status)
	}

	// Make reservation for period already reserved
	{
		status, err := service.MakeReservation(ctx, luggageid, user, startdate_reserved, enddate_reserved)
		require.Error(t, err)
		require.Equal(t, false, status)
	}

	// Make reservation after period that is reserved
	{
		status, err := service.MakeReservation(ctx, luggageid, user, stardate_free, enddate_free)
		require.NoError(t, err)
		require.Equal(t, true, status)
	}

	// Make a second reservation for another luggage
	{
		status, err := service.MakeReservation(ctx, luggageid_alt, user, startdate, enddate)
		require.NoError(t, err)
		require.Equal(t, true, status)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestReservationServiceCheckAvailability(t *testing.T) {
	ctx := context.Background()
	service, err := reservationServiceRegistry.Get(ctx)
	require.NoError(t, err)

	user := "test_user"
	luggageid := "test_luggage"
	luggageid_alt := "another_test_luggage"
	ids := []string{luggageid, luggageid_alt}

	startdate := "2001-01-01"
	enddate := "2001-02-01"

	startdate_alt := "2001-01-27"
	enddate_alt := "2001-02-05"

	start_both_avail := "2001-03-01"
	end_both_avail := "2001-04-01"
	expected_both_avail := []string{luggageid, luggageid_alt}

	start_one_avail := "2001-01-03"
	end_one_avail := "2001-01-05"
	expected_one_avail := []string{luggageid_alt}

	start_none_avail := "2001-01-27"
	end_none_avail := "2001-01-28"
	expected_none_avail := []string{}

	// Make reservations
	{
		_, err := service.MakeReservation(ctx, luggageid, user, startdate, enddate)
		require.NoError(t, err)
		_, err = service.MakeReservation(ctx, luggageid_alt, user, startdate_alt, enddate_alt)
		require.NoError(t, err)
	}

	// Check availability both are available
	{
		avail, err := service.CheckAvailability(ctx, ids, start_both_avail, end_both_avail)
		require.NoError(t, err)
		require.ElementsMatch(t, expected_both_avail, avail)
	}

	// Check availability one is available
	{
		avail, err := service.CheckAvailability(ctx, ids, start_one_avail, end_one_avail)
		require.NoError(t, err)
		require.ElementsMatch(t, expected_one_avail, avail)
	}

	// Check availability none are available
	{
		avail, err := service.CheckAvailability(ctx, ids, start_none_avail, end_none_avail)
		require.Error(t, err)
		require.ElementsMatch(t, expected_none_avail, avail)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}
