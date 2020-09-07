package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Daniorocket/Routine_workers/numbersql"
)

type Worker struct {
	name   string
	number int
	active bool
}

var tableWorkers []Worker

func randomNumber() {
	for {
		numbers <- rand.Intn(100000)
		time.Sleep(time.Second)
	}
}
func worker(id int, numbers <-chan int) {
	fmt.Println()
	for j := range numbers {
		s := strconv.Itoa(id)
		tableWorkers[id] = Worker{"Worker#" + s, j, true}
		fmt.Println("Worker:", id, ",random number:", j)
		if id == 0 {
			workersFile <- tableWorkers[id]
		}
		if id == 1 || id == 2 {
			workersDb <- tableWorkers[id]
		}
	}
	time.Sleep(2 * time.Second)
}
func chooseAction(db *sql.DB) {
	for {
		select {
		case <-workersDb:
			{
				insertIntoDb(workersDb, db)
			}
		case <-workersFile:
			{
				insertIntoFile(workersFile)
			}
		default:
			fmt.Println("Can't receive reply from worker channel ")
		}
	}
}
func insertIntoDb(worker <-chan Worker, db *sql.DB) {
	for v := range worker {
		fmt.Println("Read value", v, "from worker db.")
		if err := numbersql.InsertRow(db, v.name, v.number); err != nil {
			log.Println("Failed to insert row to db: ", err)
			return
		}
		//numbersql.SelectAllData(db)
	}
}
func insertIntoFile(worker <-chan Worker) {
	for v := range worker {
		fmt.Println("Read value", v, "from worker file.")
		f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		newLine := v.name + " " + strconv.Itoa(v.number)
		_, err = fmt.Fprintln(f, newLine)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("File appended successfully")
	}
}

var numbers = make(chan int)
var workersDb = make(chan Worker)
var workersFile = make(chan Worker)

const CountOfWorkers = 3

func main() {
	db, err := numbersql.CreateDb("numbers")
	if err != nil {
		log.Println("Failed to open connection:", err)
		return
	}
	if err = numbersql.InitDb(db); err != nil {
		log.Println("Failed to init db: ", err)
		return
	}
	defer db.Close()
	for i := 0; i < CountOfWorkers; i++ {
		tableWorkers = append(tableWorkers, Worker{"", 0, false})
	}
	for w := 0; w < CountOfWorkers; w++ {
		go worker(w, numbers)
	}
	go chooseAction(db)
	randomNumber()

}
