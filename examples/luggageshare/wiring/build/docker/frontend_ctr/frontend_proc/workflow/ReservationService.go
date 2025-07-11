package workflow

import (
	"context"
	"errors"
	"time"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"go.mongodb.org/mongo-driver/bson"
)

type ReservationService interface {
	MakeReservation(ctx context.Context, luggageid string, user string, startdate string, enddate string) (bool, error)
	CheckAvailability(ctx context.Context, item_ids []string, startdate string, enddate string) ([]string, error)

	// Cleanup method only used by tests
	Cleanup(ctx context.Context) error
}

type ReservationServiceImpl struct {
	reservationdb backend.NoSQLDatabase
}

func NewReservationServiceImpl(ctx context.Context, db backend.NoSQLDatabase) (ReservationService, error) {
	return &ReservationServiceImpl{reservationdb: db}, nil
}

func (r *ReservationServiceImpl) MakeReservation(ctx context.Context, luggageid string, user string, startdate string, enddate string) (bool, error) {
	// TODO: Implement

	// Step 1: Makes a reservation for given luggage id using date format YYYY-MM-DD
	targetStart, err := time.Parse(time.RFC3339, startdate+"T12:00:00+00:00")
	if err != nil {
		return false, err
	}
	targetEnd, err := time.Parse(time.RFC3339, enddate+"T12:00:00+00:00")
	if err != nil {
		return false, err
	}
	query := bson.D{{"luggageid", luggageid}}
	coll, err := r.reservationdb.GetCollection(ctx, "reservation", "reservation")
	if err != nil {
		return false, err
	}
	res, err := coll.FindMany(ctx, query)
	if err != nil {
		return false, err
	}
	var reservations []Reservation
	err = res.All(ctx, &reservations)
	if err != nil {
		return false, err
	}

	is_booked := false
	for _, res := range reservations {
		resStart, _ := time.Parse(time.RFC3339, res.StartDate+"T12:00:00+00:00")
		resEnd, _ := time.Parse(time.RFC3339, res.EndDate+"T12:00:00+00:00")
		if !((targetStart.Before(resStart) && targetEnd.Before(resStart)) || (targetStart.After(resEnd) && targetEnd.After(resEnd))) {
			is_booked = true
			break
		}
	}
	// Step 2: Return true if reservation was successful
	if !is_booked {
		reservation := Reservation{LuggageID: luggageid, User: user, StartDate: startdate, EndDate: enddate}
		err := coll.InsertOne(ctx, reservation)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Step 3: Return false and an error if reservation is not possible

	return false, errors.New("Item is already reserved for selected dates")
}

func (r *ReservationServiceImpl) CheckAvailability(ctx context.Context, item_ids []string, startdate string, enddate string) ([]string, error) {
	// TODO: Implement

	// Step 1: Checks if the given items are not already reserved using data format YYYY-MM-DD

	var filtered_ids []string
	coll, err := r.reservationdb.GetCollection(ctx, "reservation", "reservation")
	if err != nil {
		return filtered_ids, err
	}
	res, err := coll.FindMany(ctx, bson.D{})
	if err != nil {
		return filtered_ids, err
	}
	var reservations []Reservation
	err = res.All(ctx, &reservations)
	if err != nil {
		return filtered_ids, err
	}

	targetStart, err := time.Parse(time.RFC3339, startdate+"T12:00:00+00:00")
	if err != nil {
		return filtered_ids, err
	}
	targetEnd, err := time.Parse(time.RFC3339, enddate+"T12:00:00+00:00")
	if err != nil {
		return filtered_ids, err
	}

	for _, id := range item_ids {
		is_booked := false
		for _, res := range reservations {
			if res.LuggageID != id {
				continue
			}
			resStart, _ := time.Parse(time.RFC3339, res.StartDate+"T12:00:00+00:00")
			resEnd, _ := time.Parse(time.RFC3339, res.EndDate+"T12:00:00+00:00")
			if !((targetStart.Before(resStart) && targetEnd.Before(resStart)) || (targetStart.After(resEnd) && targetEnd.After(resEnd))) {
				is_booked = true
				break
			}
		}
		if !is_booked {
			filtered_ids = append(filtered_ids, id)
		}
	}

	// Step 2: Return a filtered list with only the items that are available in the given period
	if len(filtered_ids) != 0 {
		return filtered_ids, nil
	}

	// Step 3: Return an error if none of the items can be reservied during the given data interval

	return []string{}, errors.New("No available luggage")
}

func (r *ReservationServiceImpl) Cleanup(ctx context.Context) error {
	coll, err := r.reservationdb.GetCollection(ctx, "reservation", "reservation")
	if err != nil {
		return err
	}

	err = coll.DeleteMany(ctx, bson.D{})
	return err
}
