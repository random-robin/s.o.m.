package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX = 3
const MAXBUFF = 100
const MAXPROC = 10

const GRASSO = 1
const MAGRO = 0

const NORD = 0
const SUD = 1

var entrata [2][2]chan int //TIPO DIREZIONE
var uscita [2][2]chan int  //TIPO DIREZIONE
var done chan bool
var ACK [MAXPROC][2]chan int
var termina chan bool

// Usata per implementare la condizione logica nei comandi con guardia
func when(cond bool, channel chan int) chan int {
	if cond {
		return channel
	}
	return nil
}

func utente(id int, tipologia int, direzione int) {
	fmt.Printf("Sono l'utente %d e sono partito\n", id)
	var tt int
	tt = rand.Intn(5) + 1
	time.Sleep(time.Duration(tt) * time.Second)
	entrata[direzione][tipologia] <- id
	<-ACK[id][direzione]
	time.Sleep(time.Duration(tt) * time.Second)
	uscita[direzione][tipologia] <- id
	done <- true
}

func server() {
	var cont = [][]int{[]int{0, 0}, []int{0, 0}} //DIREZIONE TIPO
	for {
		select {
		//Caso MAGRO-NORD
		case x := <-when(cont[SUD][GRASSO] == 0 &&
			cont[NORD][MAGRO]+cont[NORD][GRASSO]+cont[SUD][MAGRO] < MAX, entrata[NORD][MAGRO]):
			cont[NORD][MAGRO]++
			ACK[x][NORD] <- 1
			fmt.Println("Sono MAGRO da NORD e sono ENTRATO")

		//Caso magri da sud
		case x := <-when((cont[NORD][GRASSO] == 0) &&
			(cont[SUD][GRASSO]+cont[SUD][MAGRO]+cont[NORD][MAGRO] < MAX), entrata[SUD][MAGRO]):
			cont[SUD][MAGRO]++
			ACK[x][SUD] <- 1
			fmt.Println("Sono MAGRO da SUD e sono ENTRATO")

		//Caso grasso da nord
		case x := <-when(cont[SUD][GRASSO] == 0 &&
			cont[SUD][MAGRO] == 0 &&
			(cont[NORD][MAGRO]+cont[NORD][GRASSO]) < MAX &&
			len(entrata[NORD][MAGRO]) == 0 &&
			len(entrata[SUD][MAGRO]) == 0, entrata[NORD][GRASSO]):

			cont[NORD][GRASSO]++
			ACK[x][NORD] <- 1
			fmt.Println("Sono GRASSO da NORD e sono ENTRATO")

		//Caso grasso da sud
		case x := <-when(cont[NORD][GRASSO] == 0 &&
			cont[NORD][MAGRO] == 0 &&
			(cont[SUD][MAGRO]+cont[SUD][GRASSO]) < MAX &&
			len(entrata[SUD][MAGRO]) == 0 &&
			len(entrata[NORD][MAGRO]) == 0, entrata[SUD][GRASSO]):

			cont[SUD][GRASSO]++
			ACK[x][SUD] <- 1
			fmt.Println("Sono GRASSO da SUD e sono ENTRATO")

		//Uscita magro da nord
		case <-uscita[NORD][MAGRO]:
			cont[NORD][MAGRO]--
			fmt.Println("Sono MAGRO da NORD e sono USCITO")

		//Uscita grasso da nord
		case <-uscita[NORD][GRASSO]:
			cont[NORD][GRASSO]--
			fmt.Println("Sono GRASSO da NORD e sono USCITO")

		//Uscita magro da sud
		case <-uscita[SUD][MAGRO]:
			cont[SUD][MAGRO]--
			fmt.Println("Sono MAGRO da SUD e sono USCITO")

		//Uscita grasso da sud
		case <-uscita[SUD][GRASSO]:
			cont[SUD][GRASSO]--
			fmt.Println("Sono GRASSO da SUD e sono USCITO")

		//Terminazione
		case <-termina:
			done <- true
			return
		}
	}
}

func main() {
	fmt.Println("Ciao")
	done = make(chan bool)
	termina = make(chan bool)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			entrata[i][j] = make(chan int, MAXBUFF)
			uscita[i][j] = make(chan int, MAXBUFF)
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 2; j++ {
			ACK[i][j] = make(chan int, MAXBUFF)
		}
	}
	go server()

	/*Bad but goot test*/
	for i := 0; i < 2; i++ {
		go utente(i, MAGRO, NORD)
		go utente(i+2, GRASSO, SUD)
		go utente(i+4, MAGRO, SUD)
		go utente(i+7, GRASSO, NORD)
	}
	for i := 0; i < 6; i++ {
		<-done
	}
	termina <- true
	<-done
	fmt.Println("Ho finito")
}
