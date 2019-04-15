package main

/*
Stworzone przez
Kamil Hrebeniuk
236735
 */

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var TODOCapacity = 0
var storageCapacity = 0
var bossFrequency = 0
var workersAmount = 0
var workersFrequency = 0
var clientsAmount = 0
var clientsFrequency = 0

var mode = 0


type task struct{
	first     int
	second    int
	operation int
}

type product struct{
	result int
}

func initialize(){
	file, err := os.Open("parameters.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var i = 0
	for scanner.Scan() {
		switch i {
		case 0: TODOCapacity, err = strconv.Atoi(scanner.Text())
		case 1: storageCapacity, err = strconv.Atoi(scanner.Text())
		case 2: bossFrequency, err = strconv.Atoi(scanner.Text())
		case 3: workersAmount, err = strconv.Atoi(scanner.Text())
		case 4: workersFrequency, err = strconv.Atoi(scanner.Text())
		case 5: clientsAmount, err = strconv.Atoi(scanner.Text())
		case 6: clientsFrequency, err = strconv.Atoi(scanner.Text())
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func changeMode(jobs <-chan task, results chan <- product){
	fmt.Println("Wpisz 0 by wyciszyc")
	fmt.Println("Wpisz 1 by obserwowac dzialanie")
	fmt.Println("Wpisz 2 by sprawdzic stan magazynu")
	for true {
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		mode = int(char - '0')
		fmt.Println()
		fmt.Println("Zmiana trybu na: ", mode)
		fmt.Println()

		if mode == 2{
			fmt.Println()
			fmt.Println("Wielkosc TO TO: ", cap(jobs))
			fmt.Println("Zapelnienie TO TO: ", len(jobs))
			fmt.Println("Wielkosc magazynu: ", cap(results))
			fmt.Println("Zapelnienie magazynu: ", len(results))
			fmt.Println()
		}
	}
}

func boss(myRand *rand.Rand, jobs chan<- task){

	for true {
		if cap(jobs) > len(jobs) {
			var newTask = task{myRand.Intn(100), myRand.Intn(100), myRand.Intn(3)}
			jobs <- newTask
			time.Sleep(time.Millisecond * time.Duration(bossFrequency))
			if mode == 1{
				fmt.Println("Tworze: ", newTask.first, " ", newTask.second, " ", newTask.operation)
			}
		} else {
			if mode == 1{
				fmt.Println("Nie ma miejsca")
			}
		}
	}
}

func worker(id int, jobs <-chan task, results chan <- product) {
	for myTask := range jobs {
		var myProduct product
		switch myTask.operation{
		case 0: myProduct.result = myTask.first + myTask.second
		case 1: myProduct.result = myTask.first - myTask.second
		case 2: myProduct.result = myTask.first * myTask.second
		case 3: myProduct.result = myTask.first / myTask.second
		}

		time.Sleep(time.Millisecond * time.Duration(workersFrequency))

		if mode == 1{
			fmt.Println("Opracowano: ", myTask.first, " ", myTask.second, " ", myProduct.result)
		}

		results <- myProduct
	}
}

func client(id int, results <- chan product){
	for myResult := range results {
		time.Sleep(time.Millisecond * time.Duration(clientsFrequency))
		if mode == 1{
			fmt.Println("Kupiono: ", myResult.result)
		}
	}
}

func main() {

	initialize()

	seed := rand.NewSource(time.Now().UnixNano())
	myRand := rand.New(seed)

	jobs := make(chan task, TODOCapacity)
	results := make(chan product, storageCapacity)

	go changeMode(jobs, results)

	for w := 1; w <= workersAmount; w++ {
		go worker(w, jobs, results)
	}

	go boss(myRand, jobs)


	for w := 1; w <= clientsAmount; w++ {
		go client(w, results)
	}

	for range results {
		time.Sleep(time.Hour * 1000000)
	}





}