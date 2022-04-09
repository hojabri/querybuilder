![GitHub CI](https://github.com/hojabri/querybuilder/actions/workflows/go.yml/badge.svg)
# querybuilder
querybuilder simply builds SQL queries for you.

If you need to create complex and dynamic queries (runtime queries based on some conditions), you can use this library to create them at runtime.

## install

    go get github.com/hojabri/querybuilder

## usage
There are four types of SQL queries which can be built: **SELECT**, **INSERT**, **UPDATE** and **DELETE**

Import library

    import "github.com/hojabri/querybuilder"

below are examples to create these queries:

### SELECT
To build SELECT queries, you need to first call `querybuilder.Select()` and then use a combination of below functions:

- `Table(name string)` which gets the name of the table

- `Columns(query string, args ...interface{})` gets the name of columns in the `query` parameter and optional arguments in the `args` parameter

- `Joins(tableName string, on string, joinType JoinType, args ...interface{})` to specify join tables. It gets the name of join table in the `tableName` parameter, join condition in the `on` parameter, join type in the `joinType` parameter and optional args in the `args` parameter.

join types can be one of:
`JoinInner`, `JoinLeft` or `JoinRight`

- `Where(query string, args ...interface{})` specifies the condition for the SELECT query. you can define the condition in the `query` parameter and it's arguments in the optional `args` parameter.

_Note:_ you can have many `Where` functions in any order

- `Having(query string, args ...interface{})` to use a Having conditions for SELECT queries with Groups. the parameter usage is the same as `Where` function.

_Note:_ you can have many `Having` functions in any order

- `Group(query string)` to specify GROUP BY queries. (Samples in the examples section)
- `Order(column string, direction OrderDirection)` to specify ORDER BY part of the SELECT queries. It gets column name in the `column` parameter and order direction in the `direction` parameter.
direction can be one of `OrderAsc` or `OrderDesc`
- `Limit(limit int64)` specifies LIMIT part of the SELECT query to have pagination. It accepts an `int64` value.
- `Offset(offset int64)` specifies OFFSET part of the SELECT query to have pagination. It accepts an `int64` value.
- `Build()` after specifying all SELECT functions, you need to call this method to create your final query string and also final arguments.
- `Rebind(query string)` after your final query string is ready, you can call this method to rebind your query string based on the database driver.


####Sample 1
```go
	query, args, err := querybuilder.Select().
		Table("table1").
		Build()

```

Output:

    query:  SELECT * FROM table1
    args:   []

####Sample 2
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,c3").
		Build()
```

Output:

    query:  SELECT c1,c2,c3 FROM table1
    args:   []
####Sample 3
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=true").
		Where("c2=?", 10).
		Build()
```

Output:

    query:  SELECT c1,c2,c3 FROM table1 WHERE (c1=true) AND (c2=?)
    args:   [10]

####Sample 4
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=true").
		Where("c2=? OR c3>?", 10, 20).
		Build()
```
Output:

    query:  SELECT c1,c2,c3 FROM table1 WHERE (c1=true) AND (c2=? OR c3>?)
    args:   [10 20]
####Sample 5
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=?", true).
		Where(querybuilder.In("c2", 10, 20)).
		Build()
```
Output:

    query:  SELECT c1,c2,c3 FROM table1 WHERE (c1=?) AND (c2 IN (?,?))
    args:   [true 10 20]
####Sample 6
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,SUM(c3) AS total").
		Where("c1=?", 1).
		Group("c1,c2").
		Having("SUM(c3)>?", 100).
		Build()
```
Output:

    query:  SELECT c1,c2,SUM(c3) AS total FROM table1 WHERE (c1=?) GROUP BY c1,c2 HAVING (SUM(c3)>?)
    args:   [1 100]
####Sample 7
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2,SUM(c3) AS total,AVG(c4) AS average").
		Where("c1=?", 1).
		Where("c2=?", true).
		Group("c1,c2").
		Having("SUM(c3)>?", 100).
		Having("AVG(c4)<?", 0.1).
		Build()
```
Output:

    query:  SELECT c1,c2,SUM(c3) AS total,AVG(c4) AS average FROM table1 WHERE (c1=?) AND (c2=?) GROUP BY c1,c2 HAVING (SUM(c3)>?) AND (AVG(c4)<?)
    args:   [1 true 100 0.1]
####Sample 8
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("id,c1,c2,c3").
		Joins("table2", "table1.id = table2.t_id", querybuilder.JoinLeft).
		Build()
```
Output:

    query:  SELECT id,c1,c2,c3 FROM table1 LEFT JOIN table2 ON table1.id = table2.t_id
    args:   []
####Sample 9
```go
	query, args, err = querybuilder.Select().
		Table("table1 t1").
		Columns("t1.id,t2.c3").
		Joins("table2 t2", "t1.id = t2.t_id", querybuilder.JoinInner).
		Build()
```
Output:

    query:  SELECT t1.id,t2.c3 FROM table1 t1 JOIN table2 t2 ON t1.id = t2.t_id
    args:   []
####Sample 10
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2").
		Order("c1", querybuilder.OrderDesc).
		Build()
```
Output:

    query:  SELECT c1,c2 FROM table1 ORDER BY c1 DESC
    args:   []
####Sample 11
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2").
		Order("c1", querybuilder.OrderDesc).
		Order("c2", querybuilder.OrderAsc).
		Build()
```
Output:

    query:  SELECT c1,c2 FROM table1 ORDER BY c1 DESC,c2 ASC
    args:   []
####Sample 12
```go
	query, args, err = querybuilder.Select().
		Table("table1").
		Columns("c1,c2").
		Limit(20).
		Offset(0).
		Build()
```
Output:

    query:  SELECT c1,c2 FROM table1 LIMIT 20 OFFSET 0
    args:   []

### INSERT
To build INSERT queries, you need to first call `querybuilder.Insert()` and then use a combination of below functions:
- `Table(name string)` which gets the name of the table
- `MapValues(columnValues map[string]interface{})` you can specify columns and values to be inserted to table as a `map` object. (column name in string as the `key` of the map and the value in the `value` of the map)
- `StructValues(structure interface{})` another and in some case better choice is to use any existing struct as an input for this function. It automatically extracts all columns and values from the `struct` type.
the column name will the same as Struct field name, except you specify them in the struct `db` tags. For example:
```go
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
}
```
_Note:_ if you want to skip a column to be used for insert query, you can use `"-"` for the `db` tag.

- `Build()` after specifying all INSERT functions, you need to call this method to create your final query string and also final arguments.
- `Rebind(query string)` after your final query string is ready, you can call this method to rebind your query string based on the database driver.


Sample struct type for insert examples
```go
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
	}
	sampleImage := []byte("img")
```
####Sample 1
```go
	query, args, err := querybuilder.Insert().
		Table("table1").
		MapValues(map[string]interface{}{"field1": "value1", "field2": 10}).
		Build()
```
Output:

    query:  INSERT INTO table1(field1,field2) VALUES(?,?)
    args:   [value1 10]
####Sample 2
```go
	query, args, err = querybuilder.Insert().
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Image: &sampleImage,
			Grade: 10,
		}).
		Build()
```
Output:

    query:  INSERT INTO table1(name,email,image,grade) VALUES(?,?,?,?)
    args:   [Omid o.hojabri@gmail.com [105 109 103] 10]
####Sample 3
```go
	query, args, err = querybuilder.Insert().
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Grade: 10,
		}).
		Build()
```
Output:

    query:  INSERT INTO table1(name,email,grade) VALUES(?,?,?)
    args:   [Omid o.hojabri@gmail.com 10]
### UPDATE
To build UPDATE queries, you need to first call `querybuilder.UPDATE()` and then use a combination of below functions:
- `Table(name string)` which gets the name of the table
- `MapValues(columnValues map[string]interface{})` you can specify columns and values to be updated in the table as a `map` object. (column name in string as the `key` of the map and the value in the `value` of the map)
- `StructValues(structure interface{})` another and in some case better choice is to use any existing struct as an input for this function. It automatically extracts all columns and values from the `struct` type.
  the column name will the same as Struct field name, except you specify them in the struct `db` tags. For example:
```go
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
}
```
_Note:_ if you want to skip a column to be used for update query, you can use `"-"` for the `db` tag.

- `Where(query string, args ...interface{})` specifies the condition for the UPDATE query. you can define the condition in the `query` parameter and it's arguments in the optional `args` parameter.

_Note:_ you can have many `Where` functions in any order

- `Build()` after specifying all UPDATE functions, you need to call this method to create your final query string and also final arguments.
- `Rebind(query string)` after your final query string is ready, you can call this method to rebind your query string based on the database driver.



Sample struct type for update examples
```go
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
	}
	sampleImage := []byte("img")
```
####Sample 1
```go
	query, args, err := querybuilder.Update().
		Table("table1").
		MapValues(map[string]interface{}{"field1": "value1", "field2": 10}).
		Build()
```
Output:

    query:  UPDATE table1 SET field1=?,field2=?
    args:   [value1 10]
####Sample 2
```go
	query, args, err = querybuilder.Update().
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Image: &sampleImage,
			Grade: 10,
		}).
		Build()
```
Output:

    query:  UPDATE table1 SET name=?,email=?,image=?,grade=?
    args:   [Omid o.hojabri@gmail.com [105 109 103] 10]
####Sample 3
```go
	query, args, err = querybuilder.Update().
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Grade: 10,
		}).
		Build()
```
Output:

    query:  UPDATE table1 SET name=?,email=?,grade=?
    args:   [Omid o.hojabri@gmail.com 10]
### DELETE
To build DELETE queries, you need to first call `querybuilder.DELETE()` and then use a combination of below functions:
- `Table(name string)` which gets the name of the table
- `Where(query string, args ...interface{})` specifies the condition for the DELETE query. you can define the condition in the `query` parameter and it's arguments in the optional `args` parameter.
- `Build()` after specifying all DELETE functions, you need to call this method to create your final query string and also final arguments.
- `Rebind(query string)` after your final query string is ready, you can call this method to rebind your query string based on the database driver.

_Note:_ you can have many `Where` functions in any order
####Sample 1
```go
	query, args, err := querybuilder.Delete().
		Table("table1").
		Where("id=?", 10).
		Build()
```
Output:

    query:  DELETE FROM table1 WHERE (id=?)
    args:   [10]
####Sample 2
```go
	query, args, err = querybuilder.Delete().
		Table("table1").
		Where("id=?", 10).
		Where("email=? OR name=?", "o.hojabri@gmail.com", "Omid").
		Build()
```
Output:

    query:  DELETE FROM table1 WHERE (id=?) AND (email=? OR name=?)
    args:   [10 o.hojabri@gmail.com Omid]

####Specifying database driver
If you want to use the `Rebind(query string)` function to rebinding the argument place-holders in your query, you need first specify the database driver.

So instead of using:

`querybuilder.Select()`

`querybuilder.Insert()`

`querybuilder.Update()`

`querybuilder.Delete()`

you need to use:

`SelectByDriver(driver DriverName)`

`InsertByDriver(driver DriverName)`

`UpdateByDriver(driver DriverName)`

`DeleteByDriver(driver DriverName)`

driver name can be one of:
```go
DriverPostgres         = "postgres"
DriverPGX              = "pgx"
DriverPqTimeout        = "pq-timeouts"
DriverCloudSqlPostgres = "cloudsqlpostgres"
DriverMySQL            = "mysql"
DriverSqlite3          = "sqlite3"
DriverOCI8             = "oci8"
DriverORA              = "ora"
DriverGORACLE          = "goracle"
DriverSqlServer        = "sqlserver"`
```

For example:
```go
	insertBuilder :=querybuilder.InsertByDriver(querybuilder.DriverPostgres)
	query, args, err = qb.
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Image: &sampleImage,
			Grade: 10,
		}).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	// query: INSERT INTO table1(name,email,grade) VALUES(?,?,?)
	query = insertBuilder.Rebind(query)
	// query: INSERT INTO table1(name,email,grade) VALUES($1,$2,$3)
```