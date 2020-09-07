package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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
		if tableWorkers[id].active == true {
			fmt.Println("Can't work on busy channel!")
			return
		}
		s := strconv.Itoa(id)
		tableWorkers[id] = Worker{"Worker#" + s, j, true}
		fmt.Println("Worker:", id, ",random number:", j)
		workers <- tableWorkers[id]
		time.Sleep(2 * time.Second)
		tableWorkers[id].active = false
	}
}
func insertIntodb(worker <-chan Worker, db *sql.DB) {
	for v := range worker {
		fmt.Println("Read value", v, "from worker channel.")
		if err := numbersql.InsertRow(db, v.name, v.number); err != nil {
			log.Println("Failed to insert row to db: ", err)
		}
		numbersql.SelectAllData(db)
	}
}

var numbers = make(chan int)
var workers = make(chan Worker)

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
	go insertIntodb(workers, db)
	randomNumber()

}
