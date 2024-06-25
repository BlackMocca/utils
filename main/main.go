package main

import (
	"fmt"
	"time"

	"github.com/Blackmocca/utils"
)

func main() {
	fmt.Println(utils.NewTimestampFromTime(time.Now().UTC()))
	fmt.Println(utils.NewTimestampFromNow())
	fmt.Println(utils.NewTimestampFromString("2024-12-12 15:14:00"))
	fmt.Println(utils.NewDateFromNow())
	fmt.Println(utils.NewDateFromTime(time.Now().UTC()))
	fmt.Println(utils.NewDateFromString("2024-12-12"))
}
