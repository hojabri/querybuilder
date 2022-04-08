package querybuilder_test

import (
	"fmt"
	"github.com/hojabri/querybuilder"
	"log"
)

func ExampleSelect() {
	sq := querybuilder.Select()
	
	// Sample01
	query, args, err := sq.
		Table("table1").
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample01: query:%s args:%v\n", query, args)
	
	// Sample02
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,c3").
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample02: query:%s args:%v\n", query, args)
	
	// Sample03
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=true").
		Where("c2=?", 10).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample03: query:%s args:%v\n", query, args)
	
	// Sample04
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=true").
		Where("c2=? OR c3>?", 10, 20).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample04: query:%s args:%v\n", query, args)
	
	// Sample05
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,c3").
		Where("c1=?", true).
		Where(querybuilder.In("c2", 10, 20)).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample05: query:%s args:%v\n", query, args)
	
	// Sample06
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,SUM(c3) AS total").
		Where("c1=?", 1).
		Group("c1,c2").
		Having("SUM(c3)>?", 100).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample06: query:%s args:%v\n", query, args)
	
	// Sample07
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2,SUM(c3) AS total,AVG(c4) AS average").
		Where("c1=?", 1).
		Where("c2=?", true).
		Group("c1,c2").
		Having("SUM(c3)>?", 100).
		Having("AVG(c4)<?", 0.1).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample07: query:%s args:%v\n", query, args)
	
	// Sample08
	query, args, err = sq.
		Table("table1").
		Columns("id,c1,c2,c3").
		Joins("table2", "table1.id = table2.t_id", querybuilder.JoinLeft).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample08: query:%s args:%v\n", query, args)
	
	// Sample09
	query, args, err = sq.
		Table("table1 t1").
		Columns("t1.id,t2.c3").
		Joins("table2 t2", "t1.id = t2.t_id", querybuilder.JoinInner).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample09: query:%s args:%v\n", query, args)
	
	// Sample10
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2").
		Order("c1", querybuilder.OrderDesc).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample10: query:%s args:%v\n", query, args)
	
	// Sample11
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2").
		Order("c1", querybuilder.OrderDesc).
		Order("c2", querybuilder.OrderAsc).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample11: query:%s args:%v\n", query, args)
	
	// Sample12
	query, args, err = sq.
		Table("table1").
		Columns("c1,c2").
		Limit(20).
		Offset(0).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample12: query:%s args:%v\n", query, args)
	
	// Output:
	//Sample01: query:SELECT * FROM table1 args:[]
	//Sample02: query:SELECT c1,c2,c3 FROM table1 args:[]
	//Sample03: query:SELECT c1,c2,c3 FROM table1 WHERE (c1=true) AND (c2=?) args:[10]
	//Sample04: query:SELECT c1,c2,c3 FROM table1 WHERE (c1=true) AND (c2=? OR c3>?) args:[10 20]
	//Sample05: query:SELECT c1,c2,c3 FROM table1 WHERE (c1=?) AND (c2 IN (?,?)) args:[true 10 20]
	//Sample06: query:SELECT c1,c2,SUM(c3) AS total FROM table1 WHERE (c1=?) GROUP BY c1,c2 HAVING (SUM(c3)>?) args:[1 100]
	//Sample07: query:SELECT c1,c2,SUM(c3) AS total,AVG(c4) AS average FROM table1 WHERE (c1=?) AND (c2=?) GROUP BY c1,c2 HAVING (SUM(c3)>?) AND (AVG(c4)<?) args:[1 true 100 0.1]
	//Sample08: query:SELECT id,c1,c2,c3 FROM table1 LEFT JOIN table2 ON table1.id = table2.t_id args:[]
	//Sample09: query:SELECT t1.id,t2.c3 FROM table1 t1 JOIN table2 t2 ON t1.id = t2.t_id args:[]
	//Sample10: query:SELECT c1,c2 FROM table1 ORDER BY c1 DESC args:[]
	//Sample11: query:SELECT c1,c2 FROM table1 ORDER BY c1 DESC,c2 ASC args:[]
	//Sample12: query:SELECT c1,c2 FROM table1 LIMIT 20 OFFSET 0 args:[]
}

func ExampleInsert() {
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
	}
	sampleImage := []byte("img")
	
	sq := querybuilder.Insert()
	
	// Sample01 with map[string]any or map[string]interface{} as input
	query, args, err := sq.
		Table("table1").
		MapValues(map[string]any{"field1": "value1", "field2": 10}).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample01: query:%s args:%v\n", query, args)
	
	// Sample02 with Structure as input
	query, args, err = sq.
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
	fmt.Printf("Sample02: query:%s args:%v\n", query, args)
	
	// Sample03 with Structure as input - skipping null value for pointers
	query, args, err = sq.
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Grade: 10,
		}).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample03: query:%s args:%v\n", query, args)
	
	// Sample04 nil column/value
	_, _, err = sq.
		Table("table1").
		Build()
	if err != nil {
		fmt.Printf("Sample04: err: %s", err)
	}
	
	// Output:
	//Sample01: query:INSERT INTO table1(field1,field2) VALUES(?,?) args:[value1 10]
	//Sample02: query:INSERT INTO table1(name,email,image,grade) VALUES(?,?,?,?) args:[Omid o.hojabri@gmail.com [105 109 103] 10]
	//Sample03: query:INSERT INTO table1(name,email,grade) VALUES(?,?,?) args:[Omid o.hojabri@gmail.com 10]
	//Sample04: err: column/value map is empty
}

func ExampleUpdate() {
	type sampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
	}
	sampleImage := []byte("img")
	
	sq := querybuilder.Update()
	
	// Sample01 with map[string]any or map[string]interface{} as input
	query, args, err := sq.
		Table("table1").
		MapValues(map[string]any{"field1": "value1", "field2": 10}).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample01: query:%s args:%v\n", query, args)
	
	// Sample02 with Structure as input
	query, args, err = sq.
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
	fmt.Printf("Sample02: query:%s args:%v\n", query, args)
	
	// Sample03 with Structure as input - skipping null value for pointers
	query, args, err = sq.
		Table("table1").
		StructValues(sampleStructType{
			Name:  "Omid",
			Email: "o.hojabri@gmail.com",
			ID:    nil,
			Order: 1,
			Grade: 10,
		}).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample03: query:%s args:%v\n", query, args)
	
	// Sample04 nil column/value
	_, _, err = sq.
		Table("table1").
		Build()
	if err != nil {
		fmt.Printf("Sample04: err: %s", err)
	}
	
	// Output:
	//Sample01: query:UPDATE table1 SET field1=?,field2=? args:[value1 10]
	//Sample02: query:UPDATE table1 SET name=?,email=?,image=?,grade=? args:[Omid o.hojabri@gmail.com [105 109 103] 10]
	//Sample03: query:UPDATE table1 SET name=?,email=?,grade=? args:[Omid o.hojabri@gmail.com 10]
	//Sample04: err: column/value map is empty
	
}

func ExampleDelete() {
	sq := querybuilder.Delete()
	
	// Sample01
	query, args, err := sq.
		Table("table1").
		Where("id=?", 10).
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample01: query:%s args:%v\n", query, args)
	
	// Sample02
	query, args, err = sq.
		Table("table1").
		Where("id=?", 10).
		Where("email=? OR name=?", "o.hojabri@gmail.com", "Omid").
		Build()
	if err != nil {
		log.Printf("err: %s", err)
	}
	fmt.Printf("Sample02: query:%s args:%v\n", query, args)
	
	// Sample03- wrong number of arguments
	_, _, err = sq.
		Table("table1").
		Where("id=?", 10).
		Where("email=? OR name=?", "o.hojabri@gmail.com").
		Build()
	if err != nil {
		fmt.Printf("Sample03: err: %s", err)
	}
	
	// Output:
	//Sample01: query:DELETE FROM table1 WHERE (id=?) args:[10]
	//Sample02: query:DELETE FROM table1 WHERE (id=?) AND (email=? OR name=?) args:[10 o.hojabri@gmail.com Omid]
	//Sample03: err: wrong number of arguments
	
}
