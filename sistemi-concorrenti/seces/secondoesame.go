package main

import (
	"fmt"
	"math/rand"
	"time"
)

// COSTANTI:
const MAXBUFF = 100
const MAXPROC = 10

const MAX = 20
const NT = 4
const NP = 4
const pesi = 0
const corsi = 1

type Richiesta struct {
	id  int
	ack chan int
}

var array [NT]chan int

var Area [2]string = [2]string{"pesi", "corsi"}

// CANALI:
var entraPT chan Richiesta
var uscitaPT chan Richiesta

var entraUtente [2]chan Richiesta
var uscitaUtente [2]chan Richiesta

// CANALI per terminazione:
var done = make(chan bool)
var termina = make(chan bool)

// funzioni:
func when(b bool, c chan Richiesta) chan Richiesta {
	if !b {
		return nil
	}
	return c
}

func sleepRandTime(timeLimit int) {
	if timeLimit > 0 {
		time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
	}
}

// GOROUTINE:
func Utente(id int) {
	fmt.Printf("[Utente %d] Inizio\n", id)
	scelta := rand.Intn(2) //scelta area
	tempolezione := rand.Intn(20)
	r := Richiesta{tempolezione, make(chan int)}

	fmt.Printf("[Utente %d] Vuole entrare in  %s\n", id, Area[scelta])
	entraUtente[scelta] <- r
	<-r.ack
	sleepRandTime(tempolezione) //TEMPO ARBITRARIO NEL QUALE STA NELLA ZONA SCELTA

	uscitaUtente[scelta] <- r
	<-r.ack
	fmt.Printf("[Utente %d] Uscito dall'area  %s\n", id, Area[scelta])
	sleepRandTime(2)

	done <- true
}

func PersonalTrainer(id int) {
	fmt.Printf("[PT %d] Inizio\n", id)
	for {

		r := Richiesta{id, make(chan int)}

		fmt.Printf("[PT %d] Vuole entrare in  %s\n", id, Area[corsi])
		entraPT <- r
		<-r.ack
		lezione := 0
		lezione = <-array[id]
		sleepRandTime(lezione) //TEMPO concordato NEL QUALE STA a lezione

		uscitaPT <- r
		<-r.ack
		fmt.Printf("[PT %d] Uscito dall'area  %s\n", id, Area[corsi])
		sleepRandTime(20)
	}
	done <- true
}

func palestra() {
	// variabili di stato:

	nUtenti := 0 //Numero utenti triennali
	nPT := 0     //Numero utenti magistrali
	nPTliberi := 0
	var PTliberi [NT]bool
	for i := 0; i < len(PTliberi); i++ {
		PTliberi[i] = false
	}
	fmt.Printf("[SERVER] Inizio\n")

	for {

		select {

		case r := <-when(nUtenti+nPT < MAX, entraPT): //PT
			nPT++
			nPTliberi++
			PTliberi[r.id] = true
			r.ack <- 1
		case r := <-when(nUtenti+nPT < MAX && nPTliberi > 0, entraUtente[corsi]): //PT
			nUtenti++
			nPTliberi--
			for i := 0; i < len(PTliberi); i++ {
				if PTliberi[i] == true {
					PTliberi[i] = false
					array[i] <- r.id
				}
			}
			r.ack <- 1
		case r := <-uscitaPT: //PT
			nPT--
			r.ack <- 1
		case r := <-uscitaUtente[corsi]: //PT
			nUtenti--
			r.ack <- 1
		case r := <-when(nStudenti_t+nStudenti_m < MAX && ((len(entraUtente_biblioteca[0]) == 0 && len(entraUtente_biblioteca[1]) == 0 && len(entraUtente_biblioteca[2]) == 0) || nStudenti_m > nStudenti_t), entraUtente_biblioteca[3]): //tr+laureando
			nStudenti_t++
			r.ack <- 1
		case r := <-uscitaUtente[pesi]:
			nStudenti_m--
			r.ack <- 1
		case r := <-uscitaUtente[corsi]:
			nStudenti_m--
			r.ack <- 1
		case <-termina:
			done <- true

		}
	}
}

func main() {
	fmt.Printf("[MAIN] Inizio\n\n")
	rand.Seed(time.Now().UnixNano())

	// Inizializzazione canali
	entraPT = make(chan Richiesta, MAXBUFF)
	uscitaPT = make(chan Richiesta, MAXBUFF)
	for i := 0; i < len(array); i++ {
		array[i] = false
	}
	for i := 0; i < 2; i++ {
		entraUtente[i] = make(chan Richiesta, MAXBUFF)
		entraUtente[i] = make(chan Richiesta, MAXBUFF)
	}

	// Esecuzione goroutine
	go palestra()

	for i := 0; i < MAXPROC; i++ {
		go Utente(i)
	}
	for i := 0; i < NT; i++ {
		go PersonalTrainer(i)
	}
	// Join goroutine
	for i := 0; i < MAXPROC; i++ {
		<-done
	}
	for i := 0; i < NT; i++ { //attesa PT
		<-done
	}
	// chiudo:
	termina <- true
	<-done
	fmt.Printf("\n\n[MAIN] Fine\n")
}
