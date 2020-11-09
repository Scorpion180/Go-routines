package main

import (
	"fmt"
	"sync"
	"time"
)

type Proceso struct {
	id    int
	print chan bool
	stop  chan int
}

type Procesos struct {
	proceso []Proceso
}

func ProcesoPrincipal(id int, flg chan bool,
	stop chan int, response chan bool) {
	i := uint64(0)
	for {
		select {
		case <-flg:
			fmt.Printf("id %d: %d\n", id, i)
			i = i + 1
		case cmp := <-stop:
			if cmp == id {
				response <- true
				return
			}
		default:
			i = i + 1
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func Printer(flg chan bool, close chan int, response chan bool) {
	for {
		select {
		case cmp := <-close:
			if cmp == -2 {
				response <- true
				return
			}
		default:
			flg <- true
		}
	}
}

func Close(id int, close chan int, response chan bool,
	flg chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-response:
			//response <- true
			return
		default:
			close <- id
		}
	}
}

func remove(slice []Proceso, id int) []Proceso {
	for s, element := range slice {
		if element.id == id {
			return append(slice[:s], slice[s+1:]...)
		}
	}
	return slice
}

func main() {
	var opc int = 0
	var cont int = 0
	var idTemp int = 0
	var wg sync.WaitGroup
	procesos := Procesos{}
	printFlag := make(chan bool)
	closeRoutine := make(chan int)
	response := make(chan bool)
	//printFlag <- true
	for opc != 4 {
		fmt.Println("Administrador de procesos \n Opciones")
		fmt.Println("1.- Agregar proceso")
		fmt.Println("2.- Mostrar procesos")
		fmt.Println("3.- Terminar proceso")
		fmt.Println("4.- Salir")
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			p := Proceso{id: cont, print: printFlag, stop: closeRoutine}
			go ProcesoPrincipal(cont, printFlag, closeRoutine, response)
			cont++
			procesos.proceso = append(procesos.proceso, p)
		case 2:
			go Printer(printFlag, closeRoutine, response)
			wg.Add(1)
			fmt.Scanln(&idTemp)
			go Close(-2, closeRoutine, response, printFlag, &wg)
			wg.Wait()
		case 3:
			fmt.Println("Ingrese el ID")
			fmt.Scanln(&idTemp)
			wg.Add(1)
			go Close(idTemp, closeRoutine, response, printFlag, &wg)
			wg.Wait()
			procesos.proceso = remove(procesos.proceso, idTemp)
		}

	}
}
