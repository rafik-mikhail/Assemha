package main
import (
	"fmt"
	"os"
	"database/sql"
	"net/http"
	"gopkg.in/gomail.v2"
	"strconv"
	_ "github.com/lib/pq"
)
var(
	db    *sql.DB
)
func sendMail(mail string, trip string, amount int) {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", mail)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", trip+" "+strconv.Itoa(amount))
	d := gomail.NewDialer(os.Getenv("ESERVER"), 587, os.Getenv("EGMAIL"), os.Getenv("EPASS"))
	if er := d.DialAndSend(m); er != nil {
		   fmt.Println("email error")
		   fmt.Println(er)
	}
}
func myHandler(w http.ResponseWriter, r *http.Request) {
	var html string=`<!DOCTYPE html>
	<html lang="en">
	  <body>
		<div class="container">
		  <div class="starter-template">
		  <form method="POST">
	  email:<br>
	  <input type="text" name="email" required><br>
	  trip:<br>
	  <input type="text" name="trip" required><br>
	   amount:<br>
	  <input type="number" name="amount" required><br>
	  <input type="submit" value="Submit">
	</form>
		  </div>
		</div><!-- /.container -->
	  </body>
	</html>`
	if r.Method=="POST"{
		email := r.FormValue("email")		
		trip := r.FormValue("trip")
		amount0 := r.FormValue("amount")
		amount,_ := strconv.Atoi(amount0)
		insertMepls(email, trip, amount)
	}
	fmt.Fprintf(w, html)
}
func myHandlerrr(w http.ResponseWriter, r *http.Request) {
	var html string=`<!DOCTYPE html>
	<html lang="en">
	  <body>
		<div class="container">
		  <div class="starter-template">
		  <form method="POST">
	  trip:<br>
	  <input type="text" name="trip" required><br>
	  <input type="submit" value="Submit">
	</form>
		  </div>
		</div><!-- /.container -->
	  </body>
	</html>`
	fmt.Fprintf(w, html)
	if r.Method=="POST"{
		trip := r.FormValue("trip")
		assem(trip)
	}
}
func assem(trip string){
	// search the db for the trip and get the total amount.
	emails := make([]string, 0)
	i := 0
	rows, err := db.Query("SELECT * FROM userInfo WHERE trip=$1",trip)
	if err != nil {
		fmt.Println("trip error")
		return
	}
	for rows.Next() {
		var (
			id   int
			email string
			trip string
			amount int
		)

		rows.Scan(&id, &email, &trip, &amount)
		i+=amount
	}
	rows.Close()
	j:=0
	rows1, errr := db.Query("SELECT DISTINCT email FROM userInfo WHERE trip=$1",trip)
	if errr != nil {
		fmt.Println("trip error")
		return
	}
	for rows1.Next() {
		var (
			email string
		)
		rows1.Scan(&email)
		emails = append(emails, email)
		j++
	}
	rows1.Close()
	k := i/j
	for z := 0; z<len(emails); z++{
		sendMail(emails[z], trip, k)
	}
}
func insertMepls(email string, trip string, amount int){
	stmt, _ := db.Prepare("INSERT INTO userInfo (email, trip, amount) VALUES ($1, $2, $3)")
	_, err0 := stmt.Exec(email, trip, amount)
		if err0 != nil {
			fmt.Println("insertion error")
		}
}
func main(){
	dbInfo :=fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("MYAPP_DATABASE_HOST"), 
	os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("error connecting to DB")
	} else{
		fmt.Println("DB connection successful")
	}
	//create our table
	_,errr := db.Exec(`
	CREATE TABLE userInfo (
		id SERIAL PRIMARY KEY,
		email varchar(100) NOT NULL,
		trip varchar(100) NOT NULL,
		amount int NOT NULL
	);`)
	if errr != nil{
		fmt.Println("creation error")
	}
	//Server
	http.HandleFunc("/", myHandler)
	http.HandleFunc("/done", myHandlerrr)
	if errrr := http.ListenAndServe(":"+os.Getenv("MYAPP_WEB_PORT"), nil); errrr != nil {
		fmt.Println("Server error")
	}	
}