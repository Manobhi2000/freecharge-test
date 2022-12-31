package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "gorm.io/gorm"
)

type Recharge struct {
	Phoneno        string
	Sim            string
	Availableoffer string
	Amount         int
}

/*

 INSERT INTO recharges (phoneno, sim, availableoffer,amount) VALUES ('7036495079', 'jio', '10 gb extra','3950');

 INSERT INTO recharges (phoneno, sim, availableoffer,amount) VALUES ('9701585799', 'airtel', 'free amazon prime member ship','4000');

 INSERT INTO recharges (phoneno, sim, availableoffer,amount) VALUES ('7036495003', 'vi', 'free amazon zee5 subscription','3500');

*/

type Cylinderbooking struct {
	Company string
	Amount  int
}

/*

    INSERT INTO Cylinderbookings (company,amount) VALUES ('HP', '1005');

   INSERT INTO Cylinderbookings (company,amount) VALUES ('Bharath gas', '1000');


*/

type Electricitybill struct {
	Search        string
	Servicenumber string
	Amount        int
}

/*

    INSERT INTO electricitybills (search,servicenumber,amount) VALUES ('Southern Power Distribution Corporation','3100765294','400');

   INSERT INTO electricitybills (search,servicenumber,amount) VALUES ('Southern Power Distribution Corporation','3100765294','400');


*/

var db *gorm.DB

type Input struct {
	Title string `form:"title"`
}

func Inputvalidation(req *gin.Context) *Input {

	var inputs Input

	req.Bind(&inputs)

	if len(inputs.Title) < 0 {
		req.JSON(201, "invalid name")
		return nil
	}

	return &inputs

}

func APIvalidation(inputs *Input, req *gin.Context) bool {

	var responses Recharge

	db.Where("phoneno = ?", inputs.Title).First(&responses)

	if responses.Phoneno == "" {
		return false

	}

	return true

}

func rechargesdbquery(inputs *Input) *[]Recharge {

	var response []Recharge

	db.Where("phoneno = ?", inputs.Title).Find(&response)

	return &response

}

func Cylinderbookingsdbquery(inputs *Input) *[]Cylinderbooking {

	var response []Cylinderbooking

	db.Find(&response)

	return &response

}

func Sepcsdbquery(inputs *Input) *[]Electricitybill {

	var response []Electricitybill

	db.Find(&response)

	return &response

}

func API(req *gin.Context) {

	// 5. Inside api
	// above 3 steps

	inputs := Inputvalidation(req)

	if inputs == nil {
		req.JSON(201, "invalid input")
		return
	}

	fmt.Println(inputs)

	resp := APIvalidation(inputs, req)

	if resp == false {
		req.JSON(201, "wrong details ")
		return
	}

	var finalrecharge map[string]interface{} = make(map[string]interface{})

	response := rechargesdbquery(inputs)

	if response == nil {
		req.JSON(201, "db operation failed ")
		return
	}

	finalrecharge["recharge"] = response

	response1 := Cylinderbookingsdbquery(inputs)

	if response1 == nil {
		req.JSON(201, "db operation failed ")
		return
	}

	finalrecharge["Cylinderbookingsdbquery"] = response1

	response2 := Sepcsdbquery(inputs)

	if response2 == nil {
		req.JSON(201, "db operation failed ")
		return
	}

	finalrecharge["elecbill"] = response2

	req.JSON(200, finalrecharge)
}

// 3. connect to db
// 4. create table

func init() {

	var err error

	db, err = gorm.Open("mysql", "root:manobhi8686@tcp(127.0.0.1:3306)/sys?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Print("jvt", err)
		panic("db not connected")

	}
	db.AutoMigrate(&Recharge{}, &Cylinderbooking{}, &Electricitybill{})

}

func main() {

	//  1. server
	//  2. register api to server

	r := gin.Default()

	v1 := r.Group("/free")
	{

		v1.GET("/charge", API)

	}

	r.Run(":8686")

}
