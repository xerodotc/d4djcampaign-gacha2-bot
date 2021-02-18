package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xerodotc/d4djcampaign-gacha2-bot/session"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Character: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	character := strings.ToLower(strings.TrimSpace(line))

	fmt.Println("Rolling until:", strings.Title(character))

	gotCharacter := false
	var serial string

	for !gotCharacter {
		fmt.Println("Starting new session...")

		s, err := session.NewGachaSession()
		if err != nil {
			panic(err)
		}

		for s.GetRollCount() < session.RollLimit && !gotCharacter {
			if err := s.Roll(); err != nil {
				panic(err)
			}

			if s.GetAlternateCharacter() == "" {
				fmt.Println("First roll:", strings.Title(s.GetCurrentCharacter()))
			} else {
				fmt.Printf("Next roll: %s vs %s\n", strings.Title(s.GetCurrentCharacter()), strings.Title(s.GetAlternateCharacter()))

				if s.GetAlternateCharacter() == character {
					fmt.Println("Switching...")
					if err := s.SwitchCharacter(); err != nil {
						panic(err)
					}
				}
			}

			if s.GetCurrentCharacter() == character {
				fmt.Println("Obtaining serial...")
				if err := s.ObtainSerial(); err != nil {
					panic(err)
				}

				serial = s.GetSerial()
				gotCharacter = true
			}
		}
	}

	fmt.Printf("Serial for %s: %s\n", strings.Title(character), serial)
}
