package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/aktsk/medit/cmd"

	prompt "github.com/c-bata/go-prompt"
)

var appPID string
var addrCache []cmd.Found

func executor(in string) {
	if in == "ps" {
		if pid, err := cmd.Plist(); err != nil {
			log.Fatal(err)
		} else if pid != "" {
			appPID = pid
		}

	} else if strings.HasPrefix(in, "attach") {
		slice := strings.Split(in, " ")
		var pid string
		if len(slice) > 1 {
			pid = slice[1]
		} else if appPID != "" {
			pid = appPID
		} else {
			fmt.Println("PID cannot be specified.")
		}
		cmd.Attach(pid)

	} else if strings.HasPrefix(in, "find") {
		inputSlice := strings.Split(in, " ")
		dataType := "all"
		targetVal := inputSlice[1]
		if len(inputSlice) < 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}
		if len(inputSlice) == 3 {
			targetVal = inputSlice[2]
			dataType = inputSlice[1]
		}
		foundAddr, err := cmd.Find(appPID, targetVal, dataType)
		if err != nil {
			fmt.Println(err)
		}
		addrCache = foundAddr

	} else if strings.HasPrefix(in, "filter") {
		if len(addrCache) == 0 {
			fmt.Println("No previous results. ")
			return
		}
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}

		foundAddr, err := cmd.Filter(appPID, slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}
		addrCache = foundAddr

	} else if strings.HasPrefix(in, "patch") {
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}
		err := cmd.Patch(appPID, slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}

	} else if in == "detach" {
		if err := cmd.Detach(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else if in == "exit" {
		os.Exit(0)

	} else if in == "" {

	} else {
		fmt.Println("Command not found.")
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "attach", Description: "Attach to the specified process."},
		{Text: "attach <pid>", Description: "Attach to the process specified on the command line."},
		{Text: "find   <int>", Description: "Search the specified integer."},
		{Text: "find   <datatype> <int>", Description: "Types can be specified are string, word, dword, qword."},
		{Text: "filter <int>", Description: "Filter previous search results that match the current search results."},
		{Text: "patch  <int>", Description: "Write the specified value on the address found by search."},
		{Text: "detach", Description: "Detach from the attached process."},
		{Text: "ps", Description: "Find the target process and if there is only one, specify it as the target."},
		{Text: "exit"},
	}
}

func main() {
	// for ptrace attach
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if pid, err := cmd.Plist(); err != nil {
		log.Fatal(err)
	} else if pid != "" {
		appPID = pid
	}
	addrCache = []cmd.Found{}
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("medit: simple MEmory eDIT tool"),
		prompt.OptionPrefix("> "),
		prompt.OptionInputTextColor(prompt.Cyan),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionPreviewSuggestionTextColor(prompt.Green),
		prompt.OptionDescriptionTextColor(prompt.DarkGray),
	)
	p.Run()
}
