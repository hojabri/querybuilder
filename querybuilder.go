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

// Update creates new UpdateQuery
func Update() *UpdateQuery {
	uq := UpdateQuery{}
	return &uq
}

// UpdateByDriver creates new UpdateQuery by specifying driver name
func UpdateByDriver(driver DriverName) *UpdateQuery {
	uq := UpdateQuery{driver: driver}
	return &uq
}

// Delete creates new DeleteQuery
func Delete() *DeleteQuery {
	dq := DeleteQuery{}
	return &dq
}

// DeleteByDriver creates new DeleteQuery by specifying driver name
func DeleteByDriver(driver DriverName) *DeleteQuery {
	dq := DeleteQuery{driver: driver}
	return &dq
}
