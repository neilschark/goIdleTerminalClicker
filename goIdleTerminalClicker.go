package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//cleanup function here
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func CostCalculator(increase_value int, base int) int {
	returnCost := int(round(math.Pow(float64(increase_value), 2))) * base

	return returnCost
}

func ClearOutput() {
	cmd := exec.Command("clear") //Linux
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func AddIncomeToMoney(money int, factories []Factory) int {
	for _, factory := range factories {
		money = money + CostCalculator(factory.Level, factory.BaseProductionPerSecond)
	}
	return money
}

func GetCommandFromUserRoutine(userCommand *string) {

	for {
		var input string
		fmt.Scanln(&input)
		*userCommand = input
	}
}

func UpgradeFactory(factories []Factory, numberOfFactory int, money *int) error {

	cost := CostCalculator(factories[numberOfFactory].Level, factories[numberOfFactory].BaseUpgradeCost)

	if cost < *money {
		factories[numberOfFactory].UpgradeLevelOfFactory()
		*money -= cost
	} else {
		return errors.New("not enough money")
	}

	return nil
}

type Factory struct {
	Name                    string
	Level                   int
	BaseUpgradeCost         int
	BaseProductionPerSecond int
}

func (factory *Factory) UpgradeLevelOfFactory() {
	factory.Level += 1
}

func main() {
	SetupCloseHandler()
	factories := []Factory{}

	factory_names := [5]string{"screw", "tool", "car", "bus", "airplane"}

	for _, name := range factory_names {
		factory := Factory{name, 1, 100, 1}
		factories = append(factories, factory)
	}
	for _, factory := range factories {
		fmt.Println(factory)
	}

	var money int = 100
	var messageToUser string = ""
	var userCommand string = ""

	go GetCommandFromUserRoutine(&userCommand)

	//main loop
	for {
		start := time.Now()
		ClearOutput()

		// calculations here
		money = AddIncomeToMoney(money, factories)

		fmt.Println("Your money:", money)
		fmt.Print("Your factories:\n\n")

		for index, factory := range factories {
			fmt.Println("")
			fmt.Println("[", index, "]")
			fmt.Println("Name: ", factory.Name)
			fmt.Println("Upgrade cost: ", CostCalculator(factory.Level, factory.BaseUpgradeCost))
			fmt.Println("Level: ", factory.Level)
			fmt.Println("Money per second: ", CostCalculator(factory.Level, factory.BaseProductionPerSecond))
		}

		fmt.Println("\nInput number of factory to upgrade and press enter:")

		fmt.Println("Last command pressed: ", userCommand)

		if userCommand != "" {
			if userCommandAsInt, err := strconv.Atoi(userCommand); err == nil {
				err := UpgradeFactory(factories, userCommandAsInt, &money)
				if err != nil {
					messageToUser = err.Error()
				} else {
					messageToUser = "Upgraded"
				}
			}
			userCommand = ""
		}

		fmt.Println(messageToUser)

		elapsed := time.Since(start)

		if !(elapsed > 1*time.Second) {
			time.Sleep((1 * time.Second) - elapsed) // wait until second is finished
		}
	}
}
