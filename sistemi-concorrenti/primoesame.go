package main

import (
	"fmt"
	"math/rand"
	"time"
)

// COSTANTI:
const MAXBUFF = 100
const MAXPROC = 10
const NT = 4
const MAX = 8
const FUN = 0
const FISIO = 1

type Richiesta struct {
	id        int
	documento int
	ack       chan int
}
type Richiesta struct {
	badge int
	ack   chan int
}

var array [MAXPROC]int

var Area [2]string = [2]string{"portineria", "biblioteca"}

// CANALI:
var entraUtente_portineria chan Richiesta
var uscitaUtente_portineria chan Richiesta

var entraUtente_biblioteca [4]chan Richiesta
var uscitaUtente_biblioteca [4]chan Richiesta

// CANALI per terminazione:
var done = make(chan bool)
var termina = make(chan bool)
var TerminaCentro = make(chan bool)

// funzioni:
func when_portineria(b bool, c chan Richiesta) chan Richiesta {
	if !b {
		return nil
	}
	return c
}

func when_biblioteca(b bool, c chan Richiesta) chan Richiesta {
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
	badge := rand.Intn(4)          //priorita' badge
	documento := rand.Intn(400000) //priorita' badge

	r := Richiesta{id, documento, make(chan int)}
	r1 := Richiesta{badge, make(chan int)}
	portineria := 0
	biblioteca := 1
	fmt.Printf("[Utente %d] Vuole entrare in  %s\n", id, Area[portineria])
	entraUtente_portineria <- r
	<-r.ack
	r.documento = -1

	fmt.Printf("[Utente %d] Vuole entrare in  %s\n", id, Area[biblioteca])

	entraUtente_biblioteca[r1.badge] <- r1
	<-r1.ack
	sleepRandTime(7) //TEMPO ARBITRARIO NEL QUALE STA NELLA ZONA SCELTA

	uscitaUtente_biblioteca[r1.badge] <- r1
	<-r1.ack

	uscitaUtente_portineria <- r
	r.documento = <-r.ack //////////senza uguale il valore a sinistra deve essere un canale!!!!
	fmt.Printf("[Utente %d] Uscito dall'area  %s\n", id, Area[portineria])
	sleepRandTime(2)

	done <- true
}

// SERVER:
func portineria() {
	for {
		select {

		case r := <-when_portineria(true, entraUtente_portineria):
			array[r.id] = r.documento
			r.ack <- 1
		case r := <-when_portineria(true, uscitaUtente_portineria):

			r.ack <- array[r.id]
		}
	}
}
func biblioteca() {
	// variabili di stato:

	nStudenti_t := 0 //Numero utenti triennali
	nStudenti_m := 0 //Numero utenti magistrali
	fmt.Printf("[SERVER] Inizio\n")

	for {

		select {

		case r := <-when_biblioteca(nStudenti_t+nStudenti_m < MAX && ((len(entraUtente_biblioteca[2]) == 0 && len(entraUtente_biblioteca[3]) == 0) || nStudenti_m <= nStudenti_t), entraUtente_biblioteca[0]): //mag+laureando
			nStudenti_m++
			r.ack <- 1
		case r := <-when_biblioteca(nStudenti_t+nStudenti_m < MAX && ((len(entraUtente_biblioteca[2]) == 0 && len(entraUtente_biblioteca[3]) == 0 && len(entraUtente_biblioteca[0]) == 0) || nStudenti_m <= nStudenti_t), entraUtente_biblioteca[1]): //mag+laureando
			nStudenti_m++
			r.ack <- 1
		case r := <-when_biblioteca(nStudenti_t+nStudenti_m < MAX && ((len(entraUtente_biblioteca[0]) == 0 && len(entraUtente_biblioteca[1]) == 0) || nStudenti_m > nStudenti_t), entraUtente_biblioteca[2]): //tr+laureando
			nStudenti_t++
			r.ack <- 1
		case r := <-when_biblioteca(nStudenti_t+nStudenti_m < MAX && ((len(entraUtente_biblioteca[0]) == 0 && len(entraUtente_biblioteca[1]) == 0 && len(entraUtente_biblioteca[2]) == 0) || nStudenti_m > nStudenti_t), entraUtente_biblioteca[3]): //tr+laureando
			nStudenti_t++
			r.ack <- 1
		case r := <-uscitaUtente_biblioteca[0]:
			nStudenti_m--
			r.ack <- 1
		case r := <-uscitaUtente_biblioteca[1]:
			nStudenti_m--
			r.ack <- 1
		case r := <-uscitaUtente_biblioteca[2]:
			nStudenti_t--
			r.ack <- 1
		case r := <-uscitaUtente_biblioteca[3]:
			nStudenti_t--
			r.ack <- 1

		}
	}
}

func main() {
	fmt.Printf("[MAIN] Inizio\n\n")
	rand.Seed(time.Now().UnixNano())

	// Inizializzazione canali
	entraUtente_portineria = make(chan Richiesta, MAXBUFF)
	uscitaUtente_portineria = make(chan Richiesta, MAXBUFF)

	for i := 0; i < 4; i++ {
		entraUtente_biblioteca[i] = make(chan Richiesta, MAXBUFF)
		entraUtente_biblioteca[i] = make(chan Richiesta, MAXBUFF)
	}

	// Esecuzione goroutine
	go portineria()
	go biblioteca()

	for i := 0; i < MAXPROC; i++ {
		go Utente(i)
	}
	// Join goroutine
	for i := 0; i < MAXPROC; i++ {
		<-done
	}
	fmt.Printf("\nTutti gli utenti sono terminati!\n \n")
	TerminaCentro <- true
	for i := 0; i < MAXPROC/2; i++ { //attesa bagnini
		<-done
	}
	fmt.Printf("\nTutti i bagnini sono terminati!\n \n")
	// chiudo:
	termina <- true
	<-done
	fmt.Printf("\n\n[MAIN] Fine\n")
}
