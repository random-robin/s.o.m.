#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <semaphore.h>

#define K 2

typedef struct {
	char film[K][40];
	int voti[K];
	int pareri;
	pthread_mutex_t m;
}SONDAGGIO;

SONDAGGIO sondaggio;

void init_sondaggio(){
	strcpy(sondaggio.film[0], "Guerre Stellari");
	strcpy(sondaggio.film[1], "Bello figo swagger");
	int i;
	for (i = 0; i < K; i++) {
		sondaggio.voti[i] = 0;
	}
	sondaggio.pareri = 0;
}

void stampa_media(){
	//printf("Chiamato %d\n", sondaggio.pareri);
	if(sondaggio.pareri > 0) {
		for (int i = 0; i < K; i++) {
			printf("%s\tvoto medio:\t%f\n", sondaggio.film[i], (sondaggio.voti[i]/((float)sondaggio.pareri)));
		}
	}

}
void *inizia_sondaggio(void *param) {
	int i, voto;
	for (i = 0; i < K; i++) {
		//Creazione del numero random [1..10]
		voto = (rand() %10) + 1;
		printf("Inserisco il voto %d al film %s\n", voto, sondaggio.film[i]);
		//Faccio il lock sul mutex e inizio a modificare la struttura condivisa
		pthread_mutex_lock(&(sondaggio.m));
		sondaggio.voti[i] += voto;
		sondaggio.pareri = sondaggio.pareri + 1;
		pthread_mutex_unlock(&(sondaggio.m));
		stampa_media();
	}

}
void stampa_max_film() {
	char *max_film = sondaggio.film[0];
	int max = sondaggio.voti[0];
	for (int i = 1; i < K; i++) {
		if (max < sondaggio.voti[i]) {
			max = sondaggio.voti[i];
			max_film = sondaggio.film[i];
		}
	}
	printf("%s\t Voti: %d\n", max_film, max);
}

int main(int argc, char **argv) {
	int i, n_persone;
	init_sondaggio();
	printf("Quante persone vuoi che ci siano\n");
	scanf("%d%*c", &n_persone);
	pthread_t persone[n_persone];
	for (i = 0; i < n_persone; i++) {
		pthread_create(&(persone[i]), NULL, inizia_sondaggio, NULL);
	}

	for (i = 0; i < n_persone; i++) {
		int *result;
		pthread_join(persone[i], (void*) &(result));
	}
	printf("--------MEDIA FINALE--------\n");
	stampa_media();
	printf("-------FILM CON PIU VOTI----\n");
	stampa_max_film();

}

