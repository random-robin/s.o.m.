#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <semaphore.h>
#include <unistd.h>

#define K 2
#define DOWNLOADS_MAX 3

typedef struct {
	char film[K][40];
	int voti[K];
	int pareri;
	pthread_mutex_t m;
}SONDAGGIO;

SONDAGGIO sondaggio;

//Variabili per la sincronizzazione a barriera

sem_t mutex;
sem_t barrier;
int completati = 0;
int n_persone;

//Per realizzare il concetto di priorita utilizzo il concetto di semaforo privato

//Un semaforo per ogni coda, quindi 11

sem_t download_sem[11]; //Vanno inizializzati a 0
int counter[11]; //Contatore dei processi in coda va inizializzato a 0 anche lui
sem_t download_mutex; //Va inizializzato a 1 come ogni mutex
int disponibili = DOWNLOADS_MAX; //Massimo numero di download contemporanei


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

//Implementazione del gestore con i metodi richiedi e rilascia

void richiedi_download(int media, int id) {
	printf("%d: richiedo il download\n", id);
	sem_wait(&download_mutex);
	if (!(disponibili > 0)) {
		counter[media]++;
		printf("Sono %d e mi accodo in %d\n", id, media);
		sem_post(&download_mutex);
		sem_wait(&(download_sem[media]));
		counter[media]--;
		printf("Ho media %d e sono stato risvegliato\n", media);	
	}
	disponibili--;
	printf("Risorsa allocata al thread %d\n", id);
	sem_post(&download_mutex);
}

void rilascia_download(int id) {
	sem_wait(&download_mutex);
	int priority = -1;
	disponibili++;
	printf("%d: Rilascio della risorsa\n", id);
	for (int i = 11; i >= 0 && priority == -1; i--) {
		if (counter[i] > 0) {
			priority = i;
		}
	}
	if (priority != -1) {
		printf("Trovato thread in attesa con media %d\n", priority);
		sem_post(&(download_sem[priority]));
	} else {
		printf("%d: Non ho trovato nessuno rilascio il mutex\n", id);
		sem_post(&download_mutex);
	}
}

void *inizia_sondaggio(void *param) {
	int i, voto, id, media;
	media = 0;
	for (i = 0; i < K; i++) {
		//Creazione del numero random [1..10]
		voto = (rand() %10) + 1;
		printf("Inserisco il voto %d al film %s\n", voto, sondaggio.film[i]);
		//Faccio il lock sul mutex e inizio a modificare la struttura condivisa
		pthread_mutex_lock(&(sondaggio.m));
		sondaggio.voti[i] += voto;
		sondaggio.pareri = sondaggio.pareri + 1;
		pthread_mutex_unlock(&(sondaggio.m));
		media += voto;
		stampa_media();
	}

	//Implementazione della sincronizzazione a barriera
	sem_wait(&mutex);
	completati++;
	id = completati;
	if (completati == n_persone) {
		sem_post(&barrier);
		printf("Abbiamo finito tutti\n");
	}
	sem_post(&mutex);
	//Meccanismo a tornello
	sem_wait(&barrier);
	sem_post(&barrier);
	media = media / K;
	//Implementezione della richiesta al download
	richiedi_download(media, id);
	printf("%d: Sto scaricando\n", id);
	sleep(3);
	printf("%d: Ho finito\n", id);
	rilascia_download(id);
	stampa_max_film();

}

int main(int argc, char **argv) {
	int i;
	init_sondaggio();
	//Inizializzazione dei semafori per la barriera
	
	sem_init(&mutex, 0, 1);
	sem_init(&barrier, 0, 0);
	
	//Inizializzazione semafori privati e contatori e mutex
	sem_init(&download_mutex, 0, 1);
	for (int k = 0; k < 12; k++) {
		sem_init(&(download_sem[k]), 0, 0);
		counter[k] = 0;
	}

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

}

