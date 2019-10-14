package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type BossCreated struct{
	first int
	second int
	operation int
	result int
	response chan bool
}

type WorkerGet struct{
	first int
	second int
	operation int
	result int
	response chan BossCreated
}

type WorkerCreated struct{
	value int
	response chan bool
}

type ClientGet struct{
	response chan int
}

type MachineInform struct{
	request int
	state bool
	response chan bool
}

type MachineCheck struct{
	request int
	response chan bool
}

type WorkerSend struct{
	first int
	second int
	operation int
	result int
	target int
	response chan bool
}

type MachineRequest struct{
	first int
	second int
	operation int
	result int
	target int
	response chan WorkerSend
}

var mode = 0

var TODOCapacity = 0
var storageCapacity = 0
var bossFrequency = 0
var workersAmount = 0
var workersFrequency = 0
var clientsAmount = 0
var clientsFrequency = 0
var machinesAmount = 0
var machinesFrequency = 0
var patientPenalty = 0


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
		case 7: machinesAmount, err = strconv.Atoi(scanner.Text())
		case 8: machinesFrequency, err = strconv.Atoi(scanner.Text())
		case 9: patientPenalty, err = strconv.Atoi(scanner.Text())
		}
		i++
	}
	fmt.Printf("A: ", workersAmount)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}


func changeMode(){
	fmt.Println("Wpisz 0 by wyciszyc")
	fmt.Println("Wpisz 1 by obserwowac dzialanie")
	fmt.Println("Wpisz 2 by sprawdzic stan pracownikow")
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
	}
}


func boss(BossWrites chan *BossCreated){
	for {
		write := &BossCreated{
			first:  rand.Intn(100),
			second:  rand.Intn(100),
			operation:  rand.Intn(2) + 1,
			result: 0,
			response: make(chan bool)}
		if mode == 1 {
			fmt.Println("Tworze: ", write.first, " & ", write.second)}
		BossWrites <- write
		<-write.response
		time.Sleep(time.Millisecond * time.Duration(bossFrequency))
	}
}


func worker(w int, WorkerReads chan *WorkerGet, MachineChecks chan *MachineCheck, WorkerSends chan *WorkerSend, WorkerReceives chan *MachineRequest, WorkerWrites chan *WorkerCreated){
	var patient = rand.Intn(2)
	var stats = 0
	if mode == 1{
		if patient == 0{
			fmt.Println("Jestem ", w,". Urodzilem sie niecierpliwy.")
		} else {
			fmt.Println("Jestem ", w,". Urodzilem sie cierpliwy.")
		}
	}
	for {
		read := &WorkerGet{
			first:     0,
			second:    0,
			operation: 0,
			result:    0,
			response:  make(chan BossCreated)}
		WorkerReads <- read
		myTask := <-read.response

		if myTask.operation > 0 {

			if mode == 1 {
				if myTask.operation == 1 {
					fmt.Println(w, " Otrzymalem: ", myTask.first, " + ", myTask.second)
				} else {
					fmt.Println(w, " Otrzymalem: ", myTask.first, " * ", myTask.second)
				}
			}

			var myMachine= -1
			if patient == 0 {
				var lookingForMachine= true
				for lookingForMachine {
					read := &MachineCheck{
						request:  rand.Intn(machinesAmount),
						response: make(chan bool)}

					MachineChecks <- read
					lookingForMachine = <-read.response
					myMachine = read.request
					time.Sleep(time.Millisecond * time.Duration(patientPenalty))
				}
			}

			if myMachine == -1 {
				myMachine = rand.Intn(machinesAmount)
			}

			//Przesłać zadanie do konkretnej maszyny i otrzymac odpowiedz
			write := &WorkerSend{
				first:     myTask.first,
				second:    myTask.second,
				operation: myTask.operation,
				result:    0,
				target:    myMachine,
				response:  make(chan bool)}
			WorkerSends <- write
			<-write.response

			read2 := &MachineRequest{
				first:     0,
				second:    0,
				operation: 0,
				result:    0,
				response:  make(chan WorkerSend)}
			WorkerReceives <- read2
			myTask2 := <-read2.response

			if mode == 1 {
				if myTask2.operation == 1 {
					fmt.Println(w, " Przetworzono: ", myTask2.first, " + ", myTask2.second, " = ", myTask2.result)
				} else {
					fmt.Println(w, " Przetworzono: ", myTask2.first, " * ", myTask2.second, " = ", myTask2.result)
				}
			}

			write2 := &WorkerCreated{
				value:    myTask2.result,
				response: make(chan bool)}
			WorkerWrites <- write2
			<-write2.response

			stats++
			if mode == 2 {
				if patient == 0 {
					fmt.Println("Jestem niecierpliwy ", w,". Wykonalem ", stats)
				} else{
					fmt.Println("Jestem cierpliwy ", w,". Wykonalem ", stats)
				}
			}
			time.Sleep(time.Millisecond * time.Duration(workersFrequency))
		}
	}
}


func machine(m int, MachineInforms chan *MachineInform, MachineReads chan *MachineRequest, MachineWrites chan *WorkerSend, ) {
	if mode == 1{
		fmt.Println("Maszyna nr: ", m, " startuje")
	}

	for {
		write0 := &MachineInform{
			request:  m,
			state:    true,
			response: make(chan bool)}
		MachineInforms <- write0
		<-write0.response

		read := &MachineRequest{
			first:     0,
			second:    0,
			operation: 0,
			result:    0,
			target:    m,
			response:  make(chan WorkerSend)}
		MachineReads <- read
		myTask := <-read.response

		write1 := &MachineInform{
			request:  m,
			state:    false,
			response: make(chan bool)}
		MachineInforms <- write1
		<-write1.response

		switch myTask.operation {
		case 1:
			myTask.result = myTask.first + myTask.second
			if mode == 1{
				fmt.Println("Wykonałem: ", myTask.first, " + ", myTask.second, " = ", myTask.result)}
		case 2:
			myTask.result = myTask.first * myTask.second
			if mode == 1{
			fmt.Println("Wykonałem: ", myTask.first, " * ", myTask.second, " = ", myTask.result)}
		}

		write := &WorkerSend{
			first:     myTask.first,
			second:    myTask.second,
			operation: myTask.operation,
			result:    myTask.result,
			target:    m,
			response:  make(chan bool)}
		MachineWrites <- write
		<-write.response

		time.Sleep(time.Millisecond * time.Duration(machinesFrequency))
	}
}


func client(c int, ClientReads chan *ClientGet){
	for {
		read := &ClientGet{
			response: make(chan int)}
		ClientReads <- read
		receiver := <-read.response
		if receiver > 0 {
			fmt.Println("Kupiono: ", receiver)
			time.Sleep(time.Millisecond * time.Duration(clientsFrequency))
		}
	}
}


func main() {
	go initialize()
	go changeMode()

	BossWrites := make(chan *BossCreated)
	WorkerReads := make(chan *WorkerGet)

	WorkerWrites := make(chan *WorkerCreated)
	ClientReads := make(chan *ClientGet)

	MachineInforms := make(chan *MachineInform)
	MachineChecks := make(chan *MachineCheck)

	WorkerSends := make(chan *WorkerSend)
	MachineReads := make(chan *MachineRequest)

	MachineWrites := make(chan *WorkerSend)
	WorkerReceives := make(chan *MachineRequest)



	go func() {
		var TODOList = make(map[int]BossCreated)
		var MachineList = make(map[int]bool)
		var MachineReceiveStorage = make(map[int]WorkerSend)
		var MachineOutputStorage = make(map[int]WorkerSend)
		var Storage = make(map[int]int)
		var BossTODOCounter = 0
		var WorkerTODOCounter = 0
		var WorkerStorageCounter = 0
		var ClientStorageCounter = 0


		for {
			select {
			//Boss is creating
			case write := <- BossWrites:
				TODOList[BossTODOCounter] = *write
				write.response <- true
				BossTODOCounter++

			//Worker is receiving
			case read := <- WorkerReads:
				if TODOList[WorkerTODOCounter].operation == 0 {
					read.response <- BossCreated{
						operation: 0,
					}
				} else {
					read.response <- TODOList[WorkerTODOCounter]
					WorkerTODOCounter++}

			//Worker is creating
			case write := <- WorkerWrites:
				Storage[WorkerStorageCounter] = write.value
				write.response <- true
				WorkerStorageCounter++

			//Machine is informing about state
			case write := <- MachineInforms:
				MachineList[write.request] = write.state
				write.response <- true

			//Machine availability is being tested
			case read := <- MachineChecks:
				read.response <- MachineList[read.request]

			//Worker sends task to machine
			case write := <- WorkerSends:
				MachineReceiveStorage[write.target] = *write
				write.response <- true

			//Machine is getting task
			case read := <- MachineReads:
				read.response <- MachineReceiveStorage[read.target]

			//Machine is sending answer
			case write := <- MachineWrites:
				MachineOutputStorage[write.target] = *write
				write.response <- true

			//Worker is getting answer
			case read := <- WorkerReceives:
				read.response <- MachineOutputStorage[read.target]

			//Client is receiving
			case read := <- ClientReads:
				if Storage[ClientStorageCounter] == 0 {
					read.response <- 0
				} else {
					read.response <- Storage[ClientStorageCounter]
					ClientStorageCounter++
				}
			}
		}
	}()

	time.Sleep(time.Second * 3)

	//BOSS
	go boss(BossWrites)

	//Worker
	for w := 0; w < workersAmount; w++{
		go worker(w, WorkerReads, MachineChecks, WorkerSends, WorkerReceives, WorkerWrites)
	}

	//Machine
	for m := 0; m < machinesAmount; m++{
		go machine(m, MachineInforms, MachineReads, MachineWrites)
	}

	//Client
	for c := 0; c < clientsAmount; c++ {
		go client(c, ClientReads)
	}

	time.Sleep(time.Second * 10)
}