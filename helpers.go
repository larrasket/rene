package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func isValidTimeFormat(time string) bool {
	timeFormat := "^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$"
	match, _ := regexp.MatchString(timeFormat, time)
	if !match {
		return false
	}
	timeComponents := strings.Split(time, ":")
	hours, _ := strconv.Atoi(timeComponents[0])
	minutes, _ := strconv.Atoi(timeComponents[1])
	if hours >= 0 && hours <= 23 && minutes >= 0 && minutes <= 59 {
		return true
	}
	return false
}

func timeFrameReader(msg string) (tf string) {
	fmt.Print(msg)
	reading := func(s bool) (res string) {
		fmt.Scanln(&res)
		for !isValidTimeFormat(res) {
			if res == "n" && s {
				return res
			}
			fmt.Print("Time is not in hh:mm fromat, please try again: ")
			fmt.Scanln(&res)
		}
		return
	}
	start := reading(true)
	if start == "n" {
		return "00:00-00:00"
	}
	fmt.Print("Enter end hour: ")
	end := reading(false)
	return start + "-" + end
}

func (f *field) FieldReader() error {
	if f.readFunc != nil {
		res := f.readFunc(f.message)
		f.value = res
	} else {
		fmt.Print(f.message)

		var res string
		_, err := fmt.Scanf("%s", &res)
		if err != nil {
			return err
		}
		f.value = res
	}
	return nil

}
