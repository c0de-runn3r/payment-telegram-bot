package main

import (
	"log"
	"os"

	tgClient "github.com/c0de_runn3r/payments-telegram-bot/clients/telegram"
	event_consumer "github.com/c0de_runn3r/payments-telegram-bot/consumer/event-consumer"
	"github.com/c0de_runn3r/payments-telegram-bot/events/telegram"
	storage "github.com/c0de_runn3r/payments-telegram-bot/files_storage"

	"github.com/joho/godotenv"
)

// to glue everything here

const batchSize = 100

func main() {
	token, host := processENV()

	eventsProcessor := telegram.New(
		tgClient.New(host, token),
	)
	log.Print("service started")
	storage.CreateAndMigrateDB()
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	go eventsProcessor.EveryHourCheck()
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func processENV() (token string, host string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token = os.Getenv("TOKEN")
	host = os.Getenv("HOST")
	return token, host
}
