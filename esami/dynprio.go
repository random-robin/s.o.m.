package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX_PIENE = 10
const MAX_VUOTE = 10
const MAX_PROC = 10

const PA = 2
const PV = 1
const X = 2
const Y = 3
const Z = 4

const K = 10

const bancomat = 1
const bonifico = 1
const contanti = 0

var acquisti [2]chan int
var vendite [2]chan int
var ACKacquirente [MAX_PROC]chan bool
var ACKfornitore [MAX_PROC]chan bool
var done = make(chan bool)
var termina chan bool

func acquirente(id int) {
	pagamento := rand.Intn(2)
	acquisti[pagamento] <- id
	<-ACKacquirente[id]
	fmt.Printf("Ho finito\n")
	done <- true
}

func fornitore(id int) {
	pagamento := rand.Intn(2)
	vendite[pagamento] <- id
	<-ACKfornitore[id]
	fmt.Printf("Ho finito")
	done <- true
}
func ditta() {
	cassa := PV * X
	conto := PV * X
	piene := 0
	vuote := 0
	for {

		select {
		case id := <-when(acquisti[bancomat], piene >= Y && vuote+Z < MAX_VUOTE && (cassa > K || len(acquisti[contanti]) == 0)):
			piene = piene + Y
			vuote = vuote + Z
			conto = conto + PA*Y
			ACKacquirente[id] <- true
		case id := <-when(acquisti[contanti], piene >= Y && vuote+Z < MAX_VUOTE && (cassa < K || len(acquisti[bancomat]) == 0)):
			piene = piene + Y
			vuote = vuote + Z
			cassa = cassa + PA*Y
			ACKacquirente[id] <- true
		case id := <-when(vendite[bonifico], piene + X <= MAX_PIENE && conto >= PV*X && (cassa < K || len(vendite[contanti]) == 0)):
			vuote = 0
			piene = piene + X
			conto = conto - PV*X
			ACKfornitore[id]<-true
		case id := <-when(vendite[contanti], piene + X <= MAX_PIENE && conto >= PV*X && (cassa > K || len(vendite[bonifico]) == 0)):
			vuote = 0
			piene = piene + X
			cassa = cassa - PV*X
			ACKfornitore[id]<-true

		case <-termina:
			done <- true
			break
		}
	}
}

func when(canale chan int, cond bool) chan int {
	if cond {
		return canale
	}
	return nil
}

func main() {
	//Canale per la sincronizzazione
	done = make(chan bool)
	termina = make(chan bool)
	//inizializzazione del random
	rand.Seed(time.Now().UTC().UnixNano())

	//Init dei canali per ack e comunizazione
	for i := 0; i < 2; i++ {
		acquisti[i] = make(chan int, 100)
		vendite[i] = make(chan int, 100)
	}
	for i := 0; i < 10; i++ {
		ACKacquirente[i] = make(chan bool)
		ACKfornitore[i] = make(chan bool)
	}

	//Lettura dei parametri da input
	var numeroFornitori int
	var numeroAcquirenti int
	fmt.Printf("Inserisci i fornitori\n");
	fmt.Scanf("%d", &numeroFornitori)
	fmt.Printf("Inserisici gli acquirenti\n");
	fmt.Scanf("%d", &numeroAcquirenti)

	//Server start
	go ditta()

	for i:=0; i < numeroFornitori; i++ {
		go fornitore(i)
	}
	for i:=0; i < numeroAcquirenti; i++ {
		go acquirente(i)
	}
	//Client start

	//Join dei clienti
	for i:=0; i<numeroAcquirenti+numeroFornitori; i++ {
		<-done
	}
	//Terminazione del server
	termina <- true

	//Raccolta terminazione del server
	<-done

}
