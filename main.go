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

type Worker interface {
	Work(int) error
	GetName() string
}

type FileWorker struct {
	name string
}

func (f *FileWorker) Work(number int) error {
	fmt.Println("Read value", number, "from worker file.")
	fp, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	newLine := f.GetName() + " " + strconv.Itoa(number)
	if _, err = fmt.Fprintln(fp, newLine); err != nil {
		return err
	}
	fmt.Println("File appended successfully.")
	return nil
}
func (f *FileWorker) GetName() string {
	return f.name
}

type DatabaseWorker struct {
	name string
	db   *sql.DB
}

func (d *DatabaseWorker) Work(number int) error {
	if err := numbersql.InsertRow(d.db, d.name, number); err != nil {
		log.Println("Failed to insert row to db: ", err)
		return err
	}
	if err := numbersql.SelectAllData(d.db); err != nil {
		return err
	}
	return nil
}
func (d *DatabaseWorker) GetName() string {
	return d.name
}
func randomNumber() {
	for {
		numbers <- rand.Intn(100000)
		time.Sleep(time.Second)
	}
}
func worker(w Worker, numbers <-chan int) {
	fmt.Println()
	for j := range numbers {
		if err := w.Work(j); err != nil {
			log.Println("Failed to do work:", err)
		}
	}
	time.Sleep(2 * time.Second)
}

var numbers = make(chan int)

const CountOfWorkers = 3

func main() {
	db, err := numbersql.CreateDb("numbers")
	if err != nil {
		log.Println("Failed to open connection:", err)
		return
	}
	workers := []Worker{
		&DatabaseWorker{name: "DatabaseWorker#1", db: db},
		&DatabaseWorker{name: "DatabaseWorker#2", db: db},
		&FileWorker{name: "FileWorker#1"},
	}
	if err = numbersql.InitDb(db); err != nil {
		log.Println("Failed to init db: ", err)
		return
	}
	defer db.Close()
	for w := 0; w < CountOfWorkers; w++ {
		go worker(workers[w], numbers)
	}
	randomNumber()
}
