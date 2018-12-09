package main

import (
	"fmt"
	"time"
)

const MAXPROC = 10
const MAXRES = 3
const MAXBUFF = 20
const TUR = 0
const EVE = 1

var richiesta [2]chan int     //Richieste
var rilascio [2]chan int      //Canale di rilascio, solitamente ne basta uno singolo
var risorsa [MAXPROC]chan int //Array delle risorse, da attribuire
var done = make(chan int)     //Canale per la sincronizzazione
var termina = make(chan int)  //Canale per la sincronizzazione
var contatori [2]int          //Contatori
var turUp = true

func client(i int, tipo int) { //Funzione che rappresenta il cliente ed effettua la richiesta

	richiesta[tipo] <- i
	r := <-risorsa[i]
	time.Sleep(time.Second * 2)
	rilascio[tipo] <- r
	done <- i
}

func when(b bool, c chan int) chan int { //Funzione che ci rappresenta la condizione logica dei comandi con guardia
	if !b {
		return nil
	}
	return c
}

func turPrio(disponibili int) bool {
	if turUp {
		return disponibili > 0
	} else {
		return disponibili > 0 && len(richiesta[EVE]) == 0
	}
}

func evePrio(disponibili int) bool {
	if !turUp {
		return disponibili > 0
	} else {
		return disponibili > 0 && len(richiesta[TUR]) == 0
	}
}

func server() {

	var disponibili int = MAXRES //Tiene traccia delle risorse disponibili
	var res, p, i int
	var libera [MAXRES]bool //Tiene traccia dello stato delle risorse
	//numer di processi sopsesi

	// inizializzazioni:
	for i := 0; i < MAXRES; i++ {
		libera[i] = true
	}

	for {
		time.Sleep(time.Second * 1)
		select {
		case res = <-rilascio[EVE]: //Rilascio della risorsa
			contatori[EVE]++
			if contatori[EVE] == 3 {
				contatori[EVE] = 0
				turUp = true
				fmt.Printf("10 EVENTI inverto!\n")
			}
			disponibili++
			libera[res] = true
			fmt.Printf("Esce un evento\n")
		case res = <-rilascio[TUR]:
			contatori[TUR]++
			if contatori[TUR] == 3 {
				contatori[TUR] = 0
				turUp = false
				fmt.Printf("10 TURISTI inverto!\n")
			}
			disponibili++
			libera[res] = true
			fmt.Printf("Esce un turista \n")

		case p = <-when(turPrio(disponibili), richiesta[TUR]): //Richiesta TUR

			for i = 0; i < MAXRES && !libera[i]; i++ {
			}
			libera[i] = false
			disponibili--
			risorsa[p] <- i
			fmt.Printf("Arrivato un turista \n")

		case p = <-when(evePrio(disponibili), richiesta[EVE]): //Richiesta EVE
			for i = 0; i < MAXRES && !libera[i]; i++ {
			}
			libera[i] = false
			disponibili--
			risorsa[p] <- i
			fmt.Printf("Arrivato un eventi \n")

		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("FINE !!!!!!")
			done <- 1
			return

		}
	}
}

func main() {
	var cli int
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)

	//inizializzazione canali

	for i := 0; i < 2; i++ {
		richiesta[i] = make(chan int, MAXBUFF)
		rilascio[i] = make(chan int, MAXBUFF)
	}
	for i := 0; i < cli; i++ {
		risorsa[i] = make(chan int, MAXBUFF)
	}
	tipo := 1
	for i := 0; i < cli; i++ {
		go client(i, tipo)
		tipo = 1 - tipo
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1
	<-done //attesa terminazione server

}
