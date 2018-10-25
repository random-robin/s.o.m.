#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

#define K 2

void *calcola_max (void *vector) {
	int *dataVector;
	int i, max;
	dataVector = (int*) vector;
	max = dataVector[0];
	for (i = 1; i < K; i++) {
		if (max < dataVector[i]) max = dataVector[i];
	}

	int *toReturn = (int*) malloc(sizeof(int));
	*toReturn = max;
	printf("Il massimo trovato: %d\n", *toReturn);
	pthread_exit((void*) toReturn);
}

int main (int argc, char **argv) {
	int n_values, i, n_threads;
	printf("Inserisci il numero di valori dell'array\n");
	scanf("%d%*c", &n_values);
	//Creazione del vettore dinamico
	int *vector = (int*) malloc(sizeof(int) * n_values);

	//Creazione delle variabili per i thread
	n_threads = n_values / K;
	pthread_t my_threads[n_threads];

	//Assegnamento dei valori
	for (i = 0; i < n_values; i++) {
		printf("Inserisci il %d numero\n", i + 1);
		scanf("%d%*c", &(vector[i]));
	}

	//Lancio dei thread
	for (i = 0; i < n_threads; i++) {
		pthread_create(&(my_threads[i]), NULL, calcola_max, (void*)&(vector[K*i]));
	}
	
	int *max;
	pthread_join(my_threads[0], (void*)&max);

	for (i = 1; i < n_threads; i++) {
		int *result;
		pthread_join(my_threads[i], (void*) &result);
		if ((*result) > (*max)) *max = (*result);
	}

	printf("Il massimo totale: %d\n", *max);
	free(vector);
	vector = NULL;

}
