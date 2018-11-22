package main

import (
	"fmt"
	"time"
)

const MAXPROC = 10
const MAXRES = 3
const MAXBUFF = 20

var richiesta = make(chan int, MAXBUFF)
var rilascio = make(chan int, MAXBUFF)
var risorsa [MAXPROC]chan int
var done = make(chan int)
var termina = make(chan int)

func client(i int) {

	richiesta <- i
	r := <-risorsa[i]
	fmt.Printf("\n [client %d] uso della risorsa %d\n", i, r)
	time.Sleep(time.Second * 2)
	rilascio <- r
	done <- i
}

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func server() {

	var disponibili int = MAXRES
	var res, p, i int
	var libera [MAXRES]bool
	var sospesi = 0

	// inizializzazioni:
	for i := 0; i < MAXRES; i++ {
		libera[i] = true
	}

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("nuovo ciclo server")
		select {
		case res = <-rilascio:
			if sospesi == 0 {
				disponibili++
				libera[res] = true
				fmt.Printf("[server]  restituita risorsa: %d  \n", res)
			} else {
				sospesi--
				risorsa[i] <- res
				fmt.Printf("[server]  risvegliato processo %d - allocazione risorsa %d  \n", i, res)
			}

		case p = <-when(disponibili > 0, richiesta):

			//if disponibili > 0 {
			for i = 0; i < MAXRES && !libera[i]; i++ {
			}
			libera[i] = false
			disponibili--
			risorsa[p] <- i
			fmt.Printf("[server]  allocata risorsa %d a cliente %d \n", i, p)
			/*} else {
				sospesi++
				fmt.Printf("[server]  cliente %d in attesa .. \n", p)
				//bloccato = true
			}*/
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
	for i := 0; i < cli; i++ {
		risorsa[i] = make(chan int, MAXBUFF)
	}

	for i := 0; i < cli; i++ {
		go client(i)
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1
	<-done //attesa terminazione server

}
