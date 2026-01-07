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
const NS = 4
const NM = 4
const salita = 0
const discesa = 1
const ATTESA = 3

type Richiesta struct {
	id  int
	ack chan int
}

var Area [2]string = [2]string{"pesi", "corsi"}

// CANALI:

var entraauto [2]chan Richiesta
var entracamper [2]chan Richiesta
var entraspazzaneve [2]chan Richiesta

var uscitaauto [2]chan Richiesta
var uscitacamper [2]chan Richiesta
var uscitaspazzaneve [2]chan Richiesta

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
func auto(id int) {
	//fmt.Printf("[auto %d] Inizio\n", id)

	r := Richiesta{id, make(chan int)}

	entraauto[salita] <- r
	<-r.ack
	tempo := rand.Intn(ATTESA)
	fmt.Printf("[auto %d] sta salendo\n", id)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE sale
	risposta := 0
	uscitaauto[salita] <- r
	risposta = <-r.ack
	fmt.Printf("[auto %d] e' salita\n", id)

	tempo = rand.Intn(ATTESA)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE STA parcheggiato
	r.id = risposta
	entraauto[discesa] <- r
	<-r.ack
	fmt.Printf("[auto %d] sta scendendo\n", id)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE Scende
	uscitaauto[discesa] <- r
	<-r.ack
	fmt.Printf("[auto %d] e' sceso\n", id)

	done <- true
}
func camper(id int) {
	//	fmt.Printf("[camper %d] Inizio\n", id)

	r := Richiesta{id, make(chan int)}

	entracamper[salita] <- r
	<-r.ack
	fmt.Printf("[camper %d] sta salsalendo\n", id)

	tempo := rand.Intn(ATTESA)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE sale

	uscitacamper[salita] <- r
	<-r.ack
	fmt.Printf("[camper %d] e' salito\n", id)

	tempo = rand.Intn(ATTESA)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE STA parcheggiato

	entracamper[discesa] <- r
	<-r.ack
	fmt.Printf("[camper %d] sta scendendo\n", id)

	sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE Scende
	uscitacamper[discesa] <- r
	<-r.ack
	fmt.Printf("[camper %d] e' sceso\n", id)

	done <- true
}
func spazzaneve() {
	for {
		//	fmt.Printf("[spazzaneve] Inizio\n")

		r := Richiesta{0, make(chan int)}

		entraspazzaneve[discesa] <- r
		<-r.ack
		fmt.Printf("[spazzaneve] sta scendendo\n")

		tempo := rand.Intn(ATTESA)

		sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE sale

		uscitaspazzaneve[discesa] <- r
		<-r.ack
		fmt.Printf("[spazzaneve %d] e' sceso\n")

		tempo = rand.Intn(ATTESA)

		sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE STA parcheggiato

		entraspazzaneve[salita] <- r
		<-r.ack
		fmt.Printf("[spazzaneve] sta salendo\n")

		sleepRandTime(tempo) //TEMPO ARBITRARIO NEL QUALE Scende
		uscitaspazzaneve[salita] <- r

		<-r.ack
		fmt.Printf("[spazzaneve %d] e' salito\n")

	}

}

func strada() {
	// variabili di stato:

	autodiscesa := 0
	autosalita := 0
	campersalita := 0
	camperdiscesa := 0
	spazzaneve := false
	salitalibera := true
	discesalibera := true
	autoparked := 0
	autoparkedMAX := 0
	camperparkedMAX := 0
	fmt.Printf("[SERVER] Inizio\n")

	for {

		select {

		case r := <-when(autodiscesa == 0 && autosalita == 0 && camperdiscesa == 0 && campersalita == 0, entraspazzaneve[discesa]): //
			spazzaneve = true
			salitalibera = false
			discesalibera = false
			r.ack <- 1
		case r := <-uscitaspazzaneve[discesa]: //
			spazzaneve = false
			salitalibera = true
			discesalibera = true
			r.ack <- 1
		case r := <-when(spazzaneve == false && len(entraspazzaneve[discesa]) == 0 && autosalita == 0 && campersalita == 0, entracamper[discesa]):
			camperdiscesa++
			camperparkedMAX--
			salitalibera = false
			discesalibera = true
			r.ack <- 1
		case r := <-when(spazzaneve == false && len(entraspazzaneve[discesa]) == 0 && len(entracamper[discesa]) == 0 && campersalita == 0, entraauto[discesa]):
			autodiscesa++
			if r.id == 1 {
				autoparked--

			} else if r.id == 2 {
				autoparkedMAX--

			}
			if salitalibera && discesalibera {
				print("ciao")
			}
			r.ack <- 1
		case r := <-uscitaauto[discesa]: //
			autodiscesa--
			r.ack <- 1
		case r := <-uscitacamper[discesa]: //
			camperdiscesa--
			if camperdiscesa == 0 {
				salitalibera = true
				discesalibera = true
			}

			r.ack <- 1
			//4 len case r := <-when(len(entraauto[salita])==0&&autodiscesa == 0 && autosalita == 0 && camperdiscesa == 0 && campersalita == 0, entraspazzaneve[discesa]): //
			spazzaneve = true
			salitalibera = false
			discesalibera = false
			r.ack <- 1
		case r := <-uscitaspazzaneve[salita]: //
			spazzaneve = false
			salitalibera = true
			discesalibera = true
			r.ack <- 1

		case r := <-when(autoparkedMAX+camperparkedMAX < NM && spazzaneve == false && camperdiscesa == 0 && autodiscesa == 0 && len(entraspazzaneve[discesa]) == 0 && len(entracamper[discesa]) == 0 && len(entraauto[discesa]) == 0, entracamper[salita]):
			campersalita++
			camperparkedMAX++
			salitalibera = true
			discesalibera = false
			r.ack <- 1
			//fatto
		case r := <-when((autoparked < NS || autoparkedMAX+camperparkedMAX < NM) && spazzaneve == false && camperdiscesa == 0 && len(entraspazzaneve[discesa]) == 0 && len(entracamper[discesa]) == 0 && len(entraauto[discesa]) == 0 && len(entracamper[salita]) == 0, entraauto[salita]):
			autosalita++
			risposta := 0
			if autoparked < NS {
				autoparked++
				risposta = 1
			} else {
				autoparkedMAX++
				risposta = 2
			}
			r.ack <- risposta
			//fatto
		case r := <-uscitaauto[salita]: //
			autosalita--
			r.ack <- 1
		case r := <-uscitacamper[salita]: //fatto
			campersalita--
			if campersalita == 0 {
				salitalibera = true
				discesalibera = true
			}

			r.ack <- 1
		case <-termina:
			done <- true

		}
	}
}

func main() {
	fmt.Printf("[MAIN] Inizio\n\n")
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 2; i++ {
		entraauto[i] = make(chan Richiesta)
		entracamper[i] = make(chan Richiesta)
		entraspazzaneve[i] = make(chan Richiesta)

		uscitaauto[i] = make(chan Richiesta)
		uscitacamper[i] = make(chan Richiesta)
		uscitaspazzaneve[i] = make(chan Richiesta)
	}

	// Esecuzione goroutine
	go strada()

	for i := 0; i < MAXPROC; i++ {
		go auto(i)
	}

	for i := 0; i < MAXPROC/2; i++ {
		go camper(i)
	}

	go spazzaneve()

	// Join goroutine
	for i := 0; i < MAXPROC; i++ {
		<-done
	}
	for i := 0; i < MAXPROC/2; i++ { //attesa PT
		<-done
	}

	// chiudo:
	termina <- true
	<-done
	fmt.Printf("\n\n[MAIN] Fine\n")
}
