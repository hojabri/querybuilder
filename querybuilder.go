package querybuilder

var Driver DriverName

// Select creates new SelectQuery
func Select(name string) *SelectQuery {
	sq := SelectQuery{}
	sq.table = name
	return &sq
}

// Insert creates new InsertQuery
func Insert(name string) *InsertQuery {
	iq := InsertQuery{}
	iq.table = name
	return &iq
}

// Update creates new UpdateQuery
func Update(name string) *UpdateQuery {
	uq := UpdateQuery{}
	uq.table = name
	return &uq
}

// Delete creates new DeleteQuery
func Delete(name string) *DeleteQuery {
	dq := DeleteQuery{}
	dq.table = name
	return &dq
}

// Rebind transforms a query table QUESTION to the DB driver's bindvar type.
func Rebind(query string) string {
	return rebind(BindType(Driver), query)
}
