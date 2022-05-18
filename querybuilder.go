package querybuilder

// Select creates new SelectQuery
func Select(name string) *SelectQuery {
	sq := SelectQuery{}
	sq.table = name
	return &sq
}

// SelectByDriver creates new SelectQuery by specifying driver name
func SelectByDriver(name string, driver DriverName) *SelectQuery {
	sq := SelectQuery{driver: driver}
	sq.table = name
	return &sq
}

// Insert creates new InsertQuery
func Insert(name string) *InsertQuery {
	iq := InsertQuery{}
	iq.table = name
	return &iq
}

// InsertByDriver creates new InsertQuery by specifying driver name
func InsertByDriver(name string, driver DriverName) *InsertQuery {
	iq := InsertQuery{driver: driver}
	iq.table = name
	return &iq
}

// Update creates new UpdateQuery
func Update(name string) *UpdateQuery {
	uq := UpdateQuery{}
	uq.table = name
	return &uq
}

// UpdateByDriver creates new UpdateQuery by specifying driver name
func UpdateByDriver(name string, driver DriverName) *UpdateQuery {
	uq := UpdateQuery{driver: driver}
	uq.table = name
	return &uq
}

// Delete creates new DeleteQuery
func Delete(name string) *DeleteQuery {
	dq := DeleteQuery{}
	dq.table = name
	return &dq
}

// DeleteByDriver creates new DeleteQuery by specifying driver name
func DeleteByDriver(name string, driver DriverName) *DeleteQuery {
	dq := DeleteQuery{driver: driver}
	dq.table = name
	return &dq
}
