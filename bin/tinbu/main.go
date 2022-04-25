package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/advbet/tinbu"
)

func main() {
	var url string

	flag.StringVar(&url, "url", "http://www.lotterynumbersxml.com/lotterydata/.../lottery.xml", "URL for main lottery feed XML document")
	flag.Parse()

	c := tinbu.Client{
		URL: url,
	}

	games, err := c.Load(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for _, game := range games {
		fmt.Println("ID:", game.ID)
		fmt.Println("Name:", game.Name)
		fmt.Println("States/provinces:")
		for _, s := range game.StateProvs {
			fmt.Println("\tID:", s.ID)
			fmt.Println("\tName:", s.Name)
			fmt.Println("\tCountry:", s.Country)
		}
		fmt.Println("Update time:", game.UpdateTime)
		fmt.Println("Last draw:", game.LastDrawNumbers)
		fmt.Println("Last draw date:", game.LastDrawDate.Format("2006-01-02"))
		fmt.Println("Next draw date:", game.NextDrawDate.Format("2006-01-02"))
		if game.Jackpot != nil {
			fmt.Println("Jackpot:", game.Jackpot)
		}
		fmt.Println()
	}
}
