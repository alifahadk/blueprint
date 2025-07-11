package workloadgen

import (
	"context"
	"errors"
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"

	"github.com/blueprint-uservices/blueprint/runtime/core/workload"
)

// Workload specific flags
var outfile = flag.String("outfile", "stats.csv", "Outfile where individual request information will be stored")
var duration = flag.String("duration", "1m", "Duration for which the workload should be run")
var tput = flag.Int64("tput", 100, "Desired throughput")

type ComplexWorkload interface {
	ImplementsComplexWorkload(ctx context.Context) error
}

type complexWldGen struct {
	ComplexWorkload
	ValidUsers []string
	ValidItems []string

	frontend workflow.FrontendService
}

func NewComplexWorkload(ctx context.Context, frontend workflow.FrontendService) (ComplexWorkload, error) {
	w := &complexWldGen{frontend: frontend}
	return w, nil
}

type FnType func() error

func statWrapper(fn FnType) workload.Stat {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	s := workload.Stat{}
	s.Start = start.UnixNano()
	s.Duration = duration.Nanoseconds()
	s.IsError = (err != nil)
	return s
}

func (w *complexWldGen) RunSearchHandler(ctx context.Context) workload.Stat {
	color, length, breadth, height, price, startDate, endDate := GenSearchHandler()
	return statWrapper(func() error {
		_, err := w.frontend.Search(ctx, color, length, breadth, height, price, startDate, endDate)

		return err
	})
}

func (w *complexWldGen) RunUserHandler(ctx context.Context) workload.Stat {
	username, email, address := GenUserHandler()
	return statWrapper(func() error {
		err := w.frontend.CreateUserProfile(ctx, username, email, address)
		// No need to add this user to the valid user list
		return err
	})
}

func (w *complexWldGen) RunItemHandler(ctx context.Context) workload.Stat {
	color, length, breadth, height, price, username := GenItemHandler()
	return statWrapper(func() error {
		_, err := w.frontend.AddItem(ctx, color, length, breadth, height, price, username)
		// No need to add this item to the valid item list
		return err
	})
}

func (w *complexWldGen) RunReservationHandler(ctx context.Context) workload.Stat {
	luggageid, text, username, rating := GenReservationHandler(w.ValidItems, w.ValidUsers)
	return statWrapper(func() error {
		_, err := w.frontend.MakeReservation(ctx, luggageid, text, username, rating)

		return err
	})
}

func (w *complexWldGen) RunReviewHandler(ctx context.Context) workload.Stat {
	luggageid, text, username, rating := GenReviewHandler(w.ValidItems, w.ValidUsers)
	return statWrapper(func() error {
		return w.frontend.WriteReview(ctx, luggageid, text, username, rating)
	})
}

func (w *complexWldGen) InitializeSystem(ctx context.Context) error {

	// Initialize 100 users
	for i := 1; i <= 100; i++ {
		username := "user_" + strconv.Itoa(i)
		_, email, address := GenUserHandler()
		err := w.frontend.CreateUserProfile(ctx, username, email, address)
		if err != nil {
			return err
		}
		w.ValidUsers = append(w.ValidUsers, username)

		// For each user, add upto 10 items
		num_items := rand.Intn(10) + 1
		for j := 0; j < num_items; j++ {
			color, length, breadth, height, price, _ := GenItemHandler()
			item_id, err := w.frontend.AddItem(ctx, color, length, breadth, height, price, username)

			if err != nil {
				return err
			}
			w.ValidItems = append(w.ValidItems, item_id)
		}
	}

	return nil
}

func (w *complexWldGen) Run(ctx context.Context) error {

	err := w.InitializeSystem(ctx)
	if err != nil {
		return errors.New("Failed to initialize the system with users and items")
	}
	// Configure an open-loop workload for executing requests against the initialized system
	wrk := workload.NewWorkload()
	// Configure the workload with the client side generators for the various APIs and their respective proportions
	wrk.AddAPI("SearchHandler", w.RunSearchHandler, 70)
	wrk.AddAPI("ItemHandler", w.RunItemHandler, 20)
	wrk.AddAPI("UserHandler", w.RunUserHandler, 2)
	wrk.AddAPI("ReservationHandler", w.RunReservationHandler, 5)
	wrk.AddAPI("ReviewHandler", w.RunReviewHandler, 3)
	// Initialize the engine
	engine, err := workload.NewEngine(*outfile, *tput, *duration, wrk)
	if err != nil {
		return err
	}
	// Run the workload
	engine.RunOpenLoop(ctx)
	// Print statistics from the workload
	return engine.PrintStats()
}

func (w *complexWldGen) ImplementsComplexWorkload(context.Context) error {
	return nil
}
