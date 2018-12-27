package main

import (
	"fmt"
	"math/rand"
	"time"
)

const CAPIENZA = 2
const TRIENNALE = 0
const MAGISTRALE = 1
const NON_LAUR = 0
const LAUR = 1
const MAXPROC = 10

var done chan bool
var terminaBiblioteca chan bool
var terminaPortineria chan bool
var ACKricevuto [MAXPROC]chan bool

var richiestaBiblioteca [2][2]chan int
var consegnaTesserino chan int
var rilascioTesserino chan int
var uscitaBiblioteca [2]chan int

func when(canale chan int, cond bool) chan int {
	if cond {
		return canale
	}
	return nil
}


func studente(id int) {
	consegnaTesserino <- id
	<-ACKricevuto[id]
	tipo := rand.Intn(2)
	laur := rand.Intn(2)
	fmt.Printf("%d: tipo:%d laur:%d\n",id,tipo,laur)
	richiestaBiblioteca[tipo][laur] <- id
	<-ACKricevuto[id]
	fmt.Printf("%d: TIPO:%d LAUR:%d sono entrato\n", id,tipo, laur)
	time.Sleep(5 *time.Second)
	uscitaBiblioteca[tipo]<- id
	<-ACKricevuto[id]
	rilascioTesserino <- id
	<-ACKricevuto[id]
	fmt.Printf("%d: TIPO:%d LAUR:%d sono uscito\n",id,tipo, laur)
	done<- true
}
//rand.Intn(numero) per generazione casuale
func portineria() {
	tesserini := 0
	for {
		select {
		case id:=<-consegnaTesserino:
			tesserini++
			ACKricevuto[id]<-true
			fmt.Printf("Tesserino consegnato\n")

		case id:= <-rilascioTesserino:
			tesserini--
			ACKricevuto[id]<-true
			fmt.Printf("Tesserino ritirato\n")
		case <-terminaPortineria:
			fmt.Printf("Portineria chiusa\n")
			done <- true
			break
		}
	}

}

func getPrioTriennaleNon(triennali, magistrali int) bool {
	if magistrali <= triennali {
		return len(richiestaBiblioteca[MAGISTRALE][NON_LAUR]) == 0 && len(richiestaBiblioteca[MAGISTRALE][LAUR]) == 0
	} else {
		return len(richiestaBiblioteca[TRIENNALE][LAUR]) == 0
	}
}

func getPrioTriennale(triennali, magistrali int) bool {
	if magistrali <= triennali {
		return len(richiestaBiblioteca[MAGISTRALE][NON_LAUR]) == 0 && len(richiestaBiblioteca[MAGISTRALE][LAUR]) == 0
	} else {
		return true
	}
}

func getPrioMagistraleNon(triennali, magistrali int) bool {
	if triennali < magistrali {
		return len(richiestaBiblioteca[TRIENNALE][NON_LAUR]) == 0 && len(richiestaBiblioteca[TRIENNALE][LAUR]) == 0
	} else {
		return len(richiestaBiblioteca[MAGISTRALE][LAUR]) == 0
	}
}

func getPrioMagistrale(triennali, magistrali int) bool {
	if triennali < magistrali {
		return len(richiestaBiblioteca[TRIENNALE][NON_LAUR]) == 0 && len(richiestaBiblioteca[TRIENNALE][LAUR]) == 0
	} else {
		return true
	}
}

func biblioteca() {
	numeroTriennali := 0
	numeroMagistrali := 0
	for {
		select {

		case id := <-when(richiestaBiblioteca[TRIENNALE][NON_LAUR], numeroMagistrali+numeroTriennali < CAPIENZA && getPrioTriennaleNon(numeroTriennali, numeroMagistrali)):
			numeroTriennali++
			ACKricevuto[id]<-true

		case id := <-when(richiestaBiblioteca[TRIENNALE][LAUR], numeroMagistrali+numeroTriennali < CAPIENZA && getPrioTriennale(numeroTriennali, numeroMagistrali)):
			numeroTriennali++
			ACKricevuto[id]<-true

		case id := <-when(richiestaBiblioteca[MAGISTRALE][NON_LAUR], numeroMagistrali+numeroTriennali < CAPIENZA && getPrioMagistraleNon(numeroTriennali, numeroMagistrali)):
			numeroMagistrali++
			ACKricevuto[id]<-true

		case id := <-when(richiestaBiblioteca[MAGISTRALE][LAUR], numeroMagistrali+numeroTriennali < CAPIENZA && getPrioMagistrale(numeroTriennali, numeroMagistrali)):
			numeroMagistrali++
			ACKricevuto[id]<-true

		case id:=  <-uscitaBiblioteca[MAGISTRALE]:
			numeroMagistrali = numeroMagistrali - 1
			ACKricevuto[id]<-true
		case id:= <-uscitaBiblioteca[TRIENNALE]:
			numeroTriennali = numeroTriennali - 1
			ACKricevuto[id]<-true
		case <-terminaBiblioteca:
			done <- true
			break
		}
	}
}

func main() {
	//Canale per la sincronizzazione
	done = make(chan bool)
	terminaBiblioteca = make(chan bool)
	terminaPortineria = make(chan bool)
	consegnaTesserino = make(chan int, 100)
	rilascioTesserino = make(chan int, 100)
	//inizializzazione del random
	rand.Seed(time.Now().UTC().UnixNano())

	//Init dei canali per ack e comunizazione
	for i := 0; i < 2; i++ {
		uscitaBiblioteca[i] = make(chan int, 100)
		for j:=0; j<2; j++ {
			richiestaBiblioteca[i][j] = make(chan int,100)
		}
	}
	var numeroStudenti int
	fmt.Printf("Inserisici gli studenti\n")
	fmt.Scanf("%d", &numeroStudenti)

	go portineria()
	//Server start
	go biblioteca()

	for i:=0; i<MAXPROC; i++{
		ACKricevuto[i] = make(chan bool)
	}
	//Client start
	for i := 0; i < numeroStudenti; i++ {
		go studente(i)
	}

	//Join dei clienti
	for i := 0; i < numeroStudenti; i++ {
		<-done
	}
	//Terminazione del biblioteca
	terminaBiblioteca <- true

	//Raccolta terminazione del biblioteca
	<-done
	terminaPortineria <- true
	<-done

}
