package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX_PIENE = 10
const K = 10

var done chan bool
var termina chan bool

type Richiesta struct {
	id  int
	ack chan int
}

func printTipo(typ int) string {
        switch typ {
        case TRIENNALE:
                return "triennale"
        case MAGISTRALE:
                return "magistrale"
        default:
                return ""
        }
}

func when(canale chan int, cond bool) chan int {
	if cond {
		return canale
	}
	return nil
}

func when(canale chan richiesta, cond bool) chan richiesta {
	if cond {
		return canale
	}
	return nil
}

//rand.Intn(numero) per generazione casuale

func server() {
	for {
		select {

		case <-termina:
			done <- true
			break
		}
	}
}
func cliente(id int) {
	richiesta := Richiesta{id: id, ack: make(chan int)}
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
		vendita[i] = make(chan int, 100)
	}
	for i := 0; i < 10; i++ {
		ACKacquirente[i] = make(chan bool)
		ACKfornitore[i] = make(chan bool)
	}

	//Lettura dei parametri da input
	var numeroFornitori int
	var numeroAcquirenti int
	fmt.Printf("Inserisci i fornitori\n")
	fmt.Scanf("%d", &numeroFornitori)
	fmt.Printf("Inserisici gli acquirenti\n")
	fmt.Scanf("%d", &numeroAcquirenti)

	//Server start
	go server()

	//Client start
	for i := 0; i < numeroFornitori; i++ {
		go fornitore(i)
	}

	//Join dei clienti
	for i := 0; i < numeroFornitori; i++ {
		<-done
	}
	//Terminazione del server
	termina <- true

	//Raccolta terminazione del server
	<-done

}
