package querybuilder

// Select creates new SelectQuery
func Select() *SelectQuery {
	sq := SelectQuery{}
	return &sq
}

// SelectByDriver creates new SelectQuery vy specifying driver name
func SelectByDriver(driver DriverName) *SelectQuery {
	sq := SelectQuery{driver: driver}
	return &sq
}
