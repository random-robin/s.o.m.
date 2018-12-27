package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXRES = 3
const LOCALI = 0
const OSPITI = 1
const MAXPROC = 10

var done chan bool
var termina chan bool

var richiestaBiglietto chan int
var richiestaControllo [2]chan int
var ACKrichiesta [MAXPROC]chan int
var terminaBiglietteria chan bool
var terminaStadio chan bool
var rilascio [2]chan int
var risorsa [MAXPROC]chan int //Array delle risorse, da attribuire

func when(canale chan int, cond bool) chan int {
	if cond {
		return canale
	}
	return nil
}

//rand.Intn(numero) per generazione casuale

func biglietteria() {
	for {
		select {
		case id := <-richiestaBiglietto:
			tipo := rand.Intn(2)
			ACKrichiesta[id] <- tipo

		case <-terminaBiglietteria:
			done <- true
			fmt.Printf("Biglietteria terminata!!!\n")
			break
		}
	}
}

func getPrioLocali(locali int, ospiti int) bool {
	if ospiti > locali {
		return len(richiestaControllo[OSPITI]) == 0
	} else {
		return true
	}
}

func getPrioOspiti(locali int, ospiti int) bool {
	if locali > ospiti {
		return len(richiestaControllo[LOCALI]) == 0
	} else {
		return true
	}
}

func tifoso(id int) {
	//Richiedo biglietto
	richiestaBiglietto <- id
	tipoBiglietto := <-ACKrichiesta[id]
	fmt.Printf("%d: biglietto di tipo %d\n", id, tipoBiglietto)

	//Prendo una guardia dal pool
	richiestaControllo[tipoBiglietto] <- id
	indice := <-risorsa[id]
	fmt.Printf("%d: inizio controllo\n", id)
	time.Sleep(5 * time.Second)
	fmt.Printf("%d: fine controllo\n", id)
	rilascio[tipoBiglietto] <- indice
	fmt.Printf("%d: entro nello stadio e ho il biglietto %d \n", id, tipoBiglietto)
	done <- true

}
func stadio() {

	var libera [MAXRES]bool //Tiene traccia dello stato delle risorse
	numeroLocali := 0
	numeroOspiti := 0
	var i int
	disponibili := MAXRES //Tiene traccia delle risorse disponibili

	for i = 0; i < MAXRES; i++ {
		libera[i] = true
	}

	for {
		select {
		case id := <-when(richiestaControllo[LOCALI], getPrioLocali(numeroLocali, numeroOspiti) && disponibili > 0):
			for i = 0; i < MAXRES && !libera[i]; i++ {
			}
			libera[i] = false
			disponibili--
			risorsa[id] <- i
			fmt.Printf("[server]  sicurezza allocata al cliente \n")
			//Caso richiesta di locali
		case id := <-when(richiestaControllo[OSPITI], getPrioOspiti(numeroLocali, numeroOspiti) && disponibili > 0):
			for i = 0; i < MAXRES && !libera[i]; i++ {
			}
			libera[i] = false
			disponibili--
			risorsa[id] <- i
			fmt.Printf("[server]  sicurezza allocata al cliente \n")
			//caso richiesta ospiti
		case id := <-rilascio[OSPITI]:
			disponibili++
			numeroOspiti++
			libera[id] = true
			fmt.Printf("sicurezza resistuita\n")
		case id := <-rilascio[LOCALI]:
			disponibili++
			numeroLocali++
			libera[id] = true
			fmt.Printf("sicurezza resistuita\n")

		case <-terminaStadio:
			done <- true
			break
		}
	}
}

func main() {
	//Canale per la sincronizzazione
	done = make(chan bool)
	richiestaBiglietto = make(chan int, 100)
	terminaStadio = make(chan bool)
	terminaBiglietteria = make(chan bool)
	//inizializzazione del random
	rand.Seed(time.Now().UTC().UnixNano())

	//Init dei canali per ack e comunizazione
	for i := 0; i < 2; i++ {
		rilascio[i] = make(chan int, 100)
		richiestaControllo[i] = make(chan int, 100)
	}
	for i := 0; i < MAXPROC; i++ {
		ACKrichiesta[i] = make(chan int)
		risorsa[i] = make(chan int)
	}

	//Lettura dei parametri da input
	var numeroTifosi int
	fmt.Printf("Inserisci i tifosi\n")
	fmt.Scanf("%d", &numeroTifosi)

	//Server start
	go stadio()
	go biglietteria()

	//Client start
	for i := 0; i < numeroTifosi; i++ {
		go tifoso(i)
	}

	//Join dei clienti
	for i := 0; i < numeroTifosi; i++ {
		<-done
	}
	//Terminazione del server
	terminaStadio <- true
	terminaBiglietteria <- true

	//Raccolta terminazione del server
	<-done
	<-done

}
