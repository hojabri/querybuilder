package querybuilder

// Select creates new SelectQuery
func Select() *SelectQuery {
	sq := SelectQuery{}
	return &sq
}

// SelectByDriver creates new SelectQuery by specifying driver name
func SelectByDriver(driver DriverName) *SelectQuery {
	sq := SelectQuery{driver: driver}
	return &sq
}

// Insert creates new InsertQuery
func Insert() *InsertQuery {
	iq := InsertQuery{}
	return &iq
}

// InsertByDriver creates new InsertQuery by specifying driver name
func InsertByDriver(driver DriverName) *InsertQuery {
	iq := InsertQuery{driver: driver}
	return &iq
}
