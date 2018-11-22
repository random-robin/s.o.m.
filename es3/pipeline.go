package main

import "fmt"

const MAXLEVELS = 10

var pipeline [MAXLEVELS]chan int

/*Canali per la sincronizzazione*/
var done = make(chan int)
var termina = make(chan int)

func pass(receiver <-chan int, sender chan<- int) {
	fmt.Println("Sono in attesa")
	x := <-receiver
	y := x + 1
	fmt.Printf("Ho ricevuto %d  e invio %d\n", x, y)
	sender <- y
}
func consumatore(receive <-chan int) {
	x := <-receive
	fmt.Printf("Ho ricevuto %d\n", x)
	done <- 1
}

func main() {
	/*Inizalizzazione dei canali*/
	for i := 0; i < MAXLEVELS; i++ {
		pipeline[i] = make(chan int)
	}
	fmt.Println("Inserisci il numero di stadi della pipeline")
	var nstadi int
	fmt.Scanf("%d", &nstadi)

	go consumatore(pipeline[nstadi])
	for i := 0; i < MAXLEVELS-1 && i < nstadi; i++ {
		go pass(pipeline[i], pipeline[i+1])
	}

	pipeline[0] <- 5
	<-done
}
