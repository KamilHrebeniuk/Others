package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type BossCreated struct {
	first     int
	second    int
	operation int
	result    int
	response  chan bool
}

type WorkerGet struct {
	first     int
	second    int
	operation int
	result    int
	response  chan BossCreated
}

type WorkerCreated struct {
	value    int
	response chan bool
}

type ClientGet struct {
	response chan int
}

type MachineInform struct {
	request  int
	state    bool
	response chan bool
}

type MachineCheck struct {
	request  int
	response chan bool
}

type WorkerSend struct {
	first     int
	second    int
	operation int
	result    int
	target    int
	response  chan bool
}

type MachineRequest struct {
	first     int
	second    int
	operation int
	result    int
	target    int
	response  chan WorkerSend
}

type ServiceSender struct {
	target   int
	response chan bool
}

type ServiceReceiver struct {
	response chan int
}
type ServiceRespond struct {
	worker int
	target   int
	response chan bool
}

type ServiceFixer struct {
	target   int
	response chan bool
}

type ServiceWorkerTask struct {
	target   int
	response chan int
}

type ServiceWorkerDone struct {
	worker int
	target   int
	response chan bool
}

type ServiceCheck struct{
	response chan int
}

var mode = 1

var TODOCapacity = 0
var storageCapacity = 0
var bossFrequency = 0
var workersAmount = 0
var workersFrequency = 0
var clientsAmount = 0
var clientsFrequency = 0
var additionMachinesAmount = 0
var multiplicationMachinesAmount = 0
var machinesFrequency = 0
var patientPenalty = 0
var simulationTime = 0
var failureProbability = 0
var serviceWorkersAmount = 0
var serviceWorkersFrequency = 0
var serviceFrequency = 0

func initialize() {
	file, err := os.Open("parameters.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var i = 0
	for scanner.Scan() {
		switch i {
		case 0:
			TODOCapacity, err = strconv.Atoi(scanner.Text())
		case 1:
			storageCapacity, err = strconv.Atoi(scanner.Text())
		case 2:
			bossFrequency, err = strconv.Atoi(scanner.Text())
		case 3:
			workersAmount, err = strconv.Atoi(scanner.Text())
		case 4:
			workersFrequency, err = strconv.Atoi(scanner.Text())
		case 5:
			clientsAmount, err = strconv.Atoi(scanner.Text())
		case 6:
			clientsFrequency, err = strconv.Atoi(scanner.Text())
		case 7:
			additionMachinesAmount, err = strconv.Atoi(scanner.Text())
		case 8:
			multiplicationMachinesAmount, err = strconv.Atoi(scanner.Text())
		case 9:
			machinesFrequency, err = strconv.Atoi(scanner.Text())
		case 10:
			patientPenalty, err = strconv.Atoi(scanner.Text())
		case 11:
			simulationTime, err = strconv.Atoi(scanner.Text())
		case 12:
			failureProbability, err = strconv.Atoi(scanner.Text())
		case 13:
			serviceWorkersAmount, err = strconv.Atoi(scanner.Text())
		case 14:
			serviceWorkersFrequency, err = strconv.Atoi(scanner.Text())
		case 15:
			serviceFrequency, err = strconv.Atoi(scanner.Text())
		}
		i++
	}
	fmt.Println("A: ", workersAmount)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func changeMode() {
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

func boss(BossWrites chan *BossCreated) {
	for {
		write := &BossCreated{
			first:     rand.Intn(100),
			second:    rand.Intn(100),
			operation: rand.Intn(2) + 1,
			result:    0,
			response:  make(chan bool)}
		if mode == 1 {
			fmt.Println("Tworze: ", write.first, " & ", write.second)
		}
		BossWrites <- write
		<-write.response
		time.Sleep(time.Millisecond * time.Duration(bossFrequency))
	}
}

func worker(w int, WorkerReads chan *WorkerGet, MachineChecks chan *MachineCheck, WorkerSends chan *WorkerSend, WorkerReceives chan *MachineRequest, WorkerWrites chan *WorkerCreated, ServiceSenders chan *ServiceSender) {
	var patient = rand.Intn(2)
	var stats = 0
	var failure = false
	var myTask = BossCreated{}
	if mode == 1 {
		if patient == 0 {
			fmt.Println("Jestem ", w, ". Urodzilem sie niecierpliwy.")
		} else {
			fmt.Println("Jestem ", w, ". Urodzilem sie cierpliwy.")
		}
	}
	for {
		if !failure {
			read := &WorkerGet{
				first:     0,
				second:    0,
				operation: 0,
				result:    0,
				response:  make(chan BossCreated)}
			WorkerReads <- read
			myTask = <-read.response
		}

		if myTask.operation > 0 {
			if mode == 1 {
				if myTask.operation == 1 {
					fmt.Println(w, " Otrzymalem: ", myTask.first, " + ", myTask.second)
				} else {
					fmt.Println(w, " Otrzymalem: ", myTask.first, " * ", myTask.second)
				}

				var myMachine = -1
				var requestTarget = 0
				if patient == 0 {
					var lookingForMachine = true
					for lookingForMachine {
						if myTask.operation == 1 {
							requestTarget = rand.Intn(additionMachinesAmount)
						} else {
							requestTarget = rand.Intn(multiplicationMachinesAmount) + additionMachinesAmount
						}

						read := &MachineCheck{
							request:  requestTarget,
							response: make(chan bool)}

						MachineChecks <- read
						lookingForMachine = <-read.response
						myMachine = read.request
						time.Sleep(time.Millisecond * time.Duration(patientPenalty))
					}
				}

				if myTask.operation == 1 {
					myMachine = rand.Intn(additionMachinesAmount)
				} else {
					myMachine = rand.Intn(multiplicationMachinesAmount) + additionMachinesAmount
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
					target:    myMachine,
					response:  make(chan WorkerSend)}
				WorkerReceives <- read2
				myTask2 := <-read2.response

				if myTask2.result == -1 {
					failure = true
					fmt.Println("FAILURE ", myTask2.target)
					write := &ServiceSender{
						target:   myTask2.target,
						response: make(chan bool)}
					ServiceSenders <- write
					<-write.response
				} else {
					failure = false
				}

				if !failure {
					//Wielu workerów może mieć ten sam target. Dane z miejsca targetu nie znikają mimo pobrania
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
							fmt.Println("Jestem niecierpliwy ", w, ". Wykonalem ", stats)
						} else {
							fmt.Println("Jestem cierpliwy ", w, ". Wykonalem ", stats)
						}
					}
					time.Sleep(time.Millisecond * time.Duration(workersFrequency))
				}
			}
		}
	}
}

func machine(m int, t int, MachineInforms chan *MachineInform, MachineReads chan *MachineRequest, MachineWrites chan *WorkerSend, ServiceFixers chan *ServiceFixer) {
	if mode == 1 {
		fmt.Println("Maszyna nr: ", m, " typu ", t, " startuje")
	}

	var broken = false

	for {
		if !broken {
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

			switch t {
			case 0:
				myTask.result = myTask.first + myTask.second
				if mode == 1 {
					fmt.Println(m, " Wykonałem: ", myTask.first, " + ", myTask.second, " = ", myTask.result)
				}
				if rand.Intn(100) < failureProbability {
					myTask.result = -1
					broken = true
				}
			case 1:
				myTask.result = myTask.first * myTask.second
				if mode == 1 {
					fmt.Println(m, " Wykonałem: ", myTask.first, " * ", myTask.second, " = ", myTask.result)
				}
				if rand.Intn(100) < failureProbability {
					myTask.result = -1
				}
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

			fmt.Println("Target: ", m)

			time.Sleep(time.Millisecond * time.Duration(machinesFrequency))

		} else {
			read := &ServiceFixer{
				target:   m,
				response: make(chan bool)}
			ServiceFixers <- read
			broken = <-read.response
		}
	}
}

func client(c int, ClientReads chan *ClientGet) {
	for {
		read := &ClientGet{
			response: make(chan int)}
		ClientReads <- read
		receiver := <-read.response

		if receiver > 0 && mode == 1 {
			fmt.Println(c, " Kupiono: ", receiver)
		}
		time.Sleep(time.Millisecond * time.Duration(clientsFrequency))
	}
}

func serviceWorker(s int, ServiceWorkerTasks chan *ServiceWorkerTask, ServiceWorkerDones chan *ServiceWorkerDone) {
	for {
		read := &ServiceWorkerTask{
			target:   s,
			response: make(chan int)}
		ServiceWorkerTasks <- read
		machineToFix := <-read.response
		if machineToFix >= 0 {
			fmt.Println("Fixxxxx: ", machineToFix)
		}
		time.Sleep(time.Millisecond * time.Duration(serviceWorkersFrequency))

		write := &ServiceWorkerDone{
			worker:   s,
			target:   machineToFix,
			response: make(chan bool)}
		ServiceWorkerDones <- write
		<-write.response
	}
}

func service(ServiceReceivers chan *ServiceReceiver, ServiceResponds chan *ServiceRespond, ServiceCheckers chan * ServiceCheck) {
	var ServiceList = make(map[int]bool)

	for i := 0; i < additionMachinesAmount+multiplicationMachinesAmount; i++ {
		ServiceList[i] = false
	}

	for {
		read := &ServiceReceiver{
			response: make(chan int)}
		ServiceReceivers <- read
		receiver := <-read.response
		if receiver >= 0 {
			fmt.Println("Awaria: ", receiver)
			if ServiceList[receiver] == false {
				ServiceList[receiver] = true

				fmt.Println("Zgłaszam awarie: ", receiver)

				write := &ServiceRespond{
					worker:    rand.Intn(serviceWorkersAmount),
					target:    receiver,
					response:  make(chan bool)}
				ServiceResponds <- write
				<-write.response
			}
		}

	//	time.Sleep(time.Millisecond * time.Duration(serviceFrequency))
			read2 := &ServiceCheck{
				response: make(chan int)}
			ServiceCheckers <- read2
			tmp := <- read2.response
			if tmp >= 0 {
				fmt.Println("Zrobione: ", tmp)
			}
			if tmp >= 0 {
				ServiceList[tmp] = false
			}
		time.Sleep(time.Millisecond * time.Duration(serviceFrequency))
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

	ServiceSenders := make(chan *ServiceSender)

	ServiceReceivers := make(chan *ServiceReceiver)
	ServiceResponds := make(chan *ServiceRespond)

	ServiceWorkerTasks := make(chan *ServiceWorkerTask)
	ServiceWorkerDones := make(chan *ServiceWorkerDone)

	ServiceFixers := make(chan *ServiceFixer)

	ServiceCheckers := make(chan *ServiceCheck)

	go func() {
		var TODOList = make(map[int]BossCreated)
		var MachineList = make(map[int]bool)
		var MachineReceiveStorage = make(map[int]WorkerSend)
		var MachineOutputStorage = make(map[int]WorkerSend)
		var Storage = make(map[int]int)
		var ServiceWorkerTODOList = make(map[int]int)
		var MachineFixList = make(map[int]bool)
		for i:=0; i < additionMachinesAmount + multiplicationMachinesAmount; i++{
			MachineFixList[i] = false
		}
		ServiceTODOList := list.New()
		ServiceDoneList := list.New()
		var BossTODOCounter = 0
		var WorkerTODOCounter = 0
		var WorkerStorageCounter = 0
		var ClientStorageCounter = 0

		for {
			select {
			//Boss is creating
			case write := <-BossWrites:
				TODOList[BossTODOCounter] = *write
				write.response <- true
				BossTODOCounter++


			//Worker is receiving
			case read := <-WorkerReads:
				if TODOList[WorkerTODOCounter].operation != 1 && TODOList[WorkerTODOCounter].operation != 2 {
					read.response <- BossCreated{
						operation: 0,
					}
				} else {
					read.response <- TODOList[WorkerTODOCounter]
					TODOList[WorkerTODOCounter] = BossCreated{operation: 0}
					WorkerTODOCounter++
				}

			//Machine is informing about state
			case write := <-MachineInforms:
				MachineList[write.request] = write.state
				write.response <- true

			//Machine availability is being tested
			case read := <-MachineChecks:
				read.response <- MachineList[read.request]

			//Worker sends task to machine
			case write := <-WorkerSends:
				MachineReceiveStorage[write.target] = *write
				write.response <- true

			//Machine is getting task
			case read := <-MachineReads:
				read.response <- MachineReceiveStorage[read.target]

			//Machine is sending answer
			case write := <-MachineWrites:
				MachineOutputStorage[write.target] = *write
				write.response <- true

			//Worker is getting answer
			case read := <-WorkerReceives:
				read.response <- MachineOutputStorage[read.target]

			//Worker is creating
			case write := <-WorkerWrites:
				if write.value > 0 {
					Storage[WorkerStorageCounter] = write.value
					write.response <- true
					WorkerStorageCounter++
				} else {
					write.response <- false
				}

			//Client is receiving
			case read := <-ClientReads:
				if Storage[ClientStorageCounter] == 0 {
					read.response <- 0
				} else {
					read.response <- Storage[ClientStorageCounter]
					ClientStorageCounter++
				}

			//Service is being informed
			case write := <-ServiceSenders:
				if write.target >= 0 {
					ServiceTODOList.PushBack(write.target)
					MachineFixList[write.target] = true
					write.response <- true
				} else {
					write.response <- false
				}

			//Service gets information
			case read := <-ServiceReceivers:
				if ServiceTODOList.Len() > 0 {
					var target= ServiceTODOList.Front()
					fmt.Println("Naprawiam: ", target.Value)
					ServiceTODOList.Remove(target)
					var tmp= target.Value
					read.response <- tmp.(int)
				} else {
					read.response <- -1
				}

			//Service put task
			case write := <-ServiceResponds:
				ServiceWorkerTODOList[write.worker] = write.target
				fmt.Println("Ustawiam workerowi: ", write.worker, " zadanie: ", write.target)
				write.response <- true

			//Service worker get task
			case read := <-ServiceWorkerTasks:
				read.response <- ServiceWorkerTODOList[read.target]

			//Service worker sets repaired
			case write := <-ServiceWorkerDones:
				ServiceWorkerTODOList[write.worker] = -1
				MachineFixList[write.target] = false
				ServiceDoneList.PushBack(write.target)
				write.response <- true

			//Machine receives
			case read := <-ServiceFixers:
				read.response <- MachineFixList[read.target]

			//Service checks done tasks
			case read := <-ServiceCheckers:
				if ServiceDoneList.Len() > 0 {
					var target= ServiceDoneList.Front()
					ServiceDoneList.Remove(target)
					read.response <- target.Value.(int)
				} else{
					read.response <- -1
				}
			}
		}
	}()

	time.Sleep(time.Second * 3)

	//BOSS
	go boss(BossWrites)

	//Worker
	for w := 0; w < workersAmount; w++ {
		go worker(w, WorkerReads, MachineChecks, WorkerSends, WorkerReceives, WorkerWrites, ServiceSenders)
	}

	//Machine
	for m := 0; m < additionMachinesAmount; m++ {
		go machine(m, 0, MachineInforms, MachineReads, MachineWrites, ServiceFixers)
	}
	for m := additionMachinesAmount; m < additionMachinesAmount+multiplicationMachinesAmount; m++ {
		go machine(m, 1, MachineInforms, MachineReads, MachineWrites, ServiceFixers)
	}

	//Client
	for c := 0; c < clientsAmount; c++ {
		go client(c, ClientReads)
	}

	//Service
	go service(ServiceReceivers, ServiceResponds, ServiceCheckers)

	//Service workers
	for s := 0; s < serviceWorkersAmount; s++ {
		go serviceWorker(s, ServiceWorkerTasks, ServiceWorkerDones)
	}

	time.Sleep(time.Second * time.Duration(simulationTime))
}
