package main

import (
	"context"
	"fmt"
	"github.com/shomali11/slacker"
	"log"
	"os"
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
	//os.Setenv("SLACK_BOT_TOKEN", "xoxb-4467573704419-4453199093895-u5iWgHWWyhqXBus2KPoFpMog")
	//os.Setenv("SLACK_APP_TOKEN", "xapp-1-A04DB49FCF9-4453159406951-458af352f8945280190d82dc0eed2be75b623ba978fa95579d0ad52615a11395")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

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

			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
