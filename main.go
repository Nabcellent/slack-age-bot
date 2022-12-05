package main

import (
	"context"
	"fmt"
	"github.com/shomali11/slacker"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("...Command Events...")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println("--------------------------------------------------------------------------------------------")
	}
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error: ", err)
		}
	}

	bot := slacker.NewClient(viper.GetString("SLACK_BOT_TOKEN"), viper.GetString("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("My year of birth is - <year>", &slacker.CommandDefinition{
		Description: "YOB calculator",
		Examples:    []string{"My year of birth is - <year>"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println(err)
			}

			age := time.Now().Year() - yob

			r := fmt.Sprintf("Your age is: %d", age)

			if err := response.Reply(r); err != nil {
				fmt.Println("Failed to reply!")
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
