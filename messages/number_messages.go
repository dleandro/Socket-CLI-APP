package messages

import (
	"log"
	"os"
	"sync"
	"time"
	"errors"
	"unicode"
)

const numberOfCharsPerMessage = 9

var numberMessagesReceived = make(chan *NumberMessages)

type NumberMessages struct {
	messageList    []string // TODO: could've used a map to avoid duplicates? it would be a different variable just for the non duplicates
	lastSummary    Summary
	currentSummary Summary
}

type Summary struct {
	numberOfUniques      int
	totalNumberOfUniques int
	numberOfDuplicates   int
}

// Init function is the entry point for this file
func Init(quit chan interface{}, wg sync.WaitGroup) *NumberMessages {
	nm := &NumberMessages{}

	wg.Add(1)
	go nm.scheduleSummary(quit, wg)  // need to add this go routine to the wait group

	return nm
}

func (nm *NumberMessages) scheduleSummary(quit chan interface{}, wg sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(10 * time.Second) {
		select {
		case <- quit:
			return
		default:
			nm.printSummary()
			nm.transferSummary()
			nm.resetCurrentSummary()
		}
	}
}

// HandleMessage exported for commands calling
func (nm *NumberMessages) HandleMessage(inputMessage string) error {
	err := checkIfMessageIsValid(inputMessage)

	if err != nil {
		return err
	}

	return nm.handleNewNumber(inputMessage)
}

func checkIfMessageIsValid(input string) error {
	// check for numberOfDigits
	if len(input) != numberOfCharsPerMessage {
		return errors.New("input message has incorrect number of digits")
	}

	// check if every char is a digit
	for _, char := range input {
		if !unicode.IsDigit(char) {
			return errors.New("input message has invalid chars")
		}
	}

	return nil
}

func (m *NumberMessages) handleNewNumber(inputNumber string) error {
	m.logNumber(inputNumber) // TODO: could've put this in a go routine
	m.messageList = append(m.messageList, inputNumber)

	return nil
}

func (m *NumberMessages) logNumber(inputNumber string) {
	file, _ := openLogFile("./" + LOG_FILE_NAME) // TODO: not well handled

	logger := log.New(file, "", log.LstdFlags)
	logger.SetFlags(0)

	if m.shouldNumberBeLogged(inputNumber) {
		logger.Println(inputNumber)
	}
}

func openLogFile(path string) (*os.File, error) {
	file, err := os.OpenFile(LOG_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return file, nil
}

func (m *NumberMessages) shouldNumberBeLogged(input string) bool {
	// could use a sort and then a search algorithm if this is not good enough
	for _, n := range m.messageList {
		if n == input {
			m.currentSummary.numberOfDuplicates++
			return false
		}
	}
	m.currentSummary.numberOfUniques++
	m.currentSummary.totalNumberOfUniques++
	return true
}

func (m *NumberMessages) printSummary() {
	differenceBetweenLastSummaryOfDuplicates := m.currentSummary.numberOfDuplicates - m.lastSummary.numberOfDuplicates
	differenceBetweenLastSummaryOfUniques := m.currentSummary.numberOfUniques - m.lastSummary.numberOfUniques

	log.Printf("Received %d unique numbers, %d duplicates. Unique total: %d", differenceBetweenLastSummaryOfUniques, differenceBetweenLastSummaryOfDuplicates, m.currentSummary.totalNumberOfUniques)
}

func (m *NumberMessages) transferSummary() {
	m.lastSummary = m.currentSummary
}

func (m *NumberMessages) resetCurrentSummary() {
	m.currentSummary.numberOfDuplicates = 0
	m.currentSummary.numberOfUniques = 0
	m.currentSummary.totalNumberOfUniques = m.lastSummary.totalNumberOfUniques
}
