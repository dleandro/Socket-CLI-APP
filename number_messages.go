package main

import (
	"log"
	"os"
)

type NumberMessages struct {
	messageList []string
	lastSummary Summary
	currentSummary Summary
} 

type Summary struct {
	numberOfUniques int 
	totalNumberOfUniques int
	numberOfDuplicates int
}

func (m *NumberMessages) handleNewNumber(inputNumber string) error {
	m.logNumber(inputNumber)
	m.messageList = append(m.messageList, inputNumber)

	return nil
}

func (m *NumberMessages) logNumber(inputNumber string) {
	file, _ := openLogFile("./" + LOG_FILE_NAME)

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

func (m *NumberMessages) getSummary() {
	differenceBetweenLastSummaryOfDuplicates := m.currentSummary.numberOfDuplicates - m.lastSummary.numberOfDuplicates
	differenceBetweenLastSummaryOfUniques :=  m.currentSummary.numberOfUniques - m.lastSummary.numberOfUniques

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