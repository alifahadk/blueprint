package workloadgen

func GenReviewHandler(valid_items []string, valid_users []string) (luggageid string, text string, username string, rating int64) {
	// TODO: Generate request parameters
	// luggageid must be of a valid luggage item
	// username must be of a valid user

	return
}

func GenUserHandler() (username string, email string, address string) {
	// TODO: Generate request parameters

	return
}

func GenSearchHandler() (color string, length int64, breadth int64, height int64, price float64, startDate string, endDate string) {
	// TODO: Generate request parameters

	// Note: As there are a limited number of items, it is not necessary that all search executions will have non-0 items returned

	return
}

func GenItemHandler() (color string, length int64, breadth int64, height int64, price float64, username string) {
	// TODO: Generate request parameters

	return
}

func GenReservationHandler(valid_items []string, valid_users []string) (luggageid string, username string, startdate string, enddate string) {
	// TODO: Generate request parameters
	// luggageid must be of a valid luggage item
	// username must be of a valid user

	// Note: As there are a limited number of items, it is normal that all reservation requests will not be successful.

	return
}
