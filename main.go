package main

import (
    "fmt"
    "bufio"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "syscall"
)

var PATH string = "/usr/bin"

func parseArgs(input string) []string {
    var start uint32 = 0
    var args []string
    cInput := strings.TrimSpace(input)
    for i := 0 ; i< len(cInput); i++ {
        if (cInput[i] == ' ') {
            args = append(args, cInput[start:i])
            start = uint32(i)+1
        }
    }
    args = append(args, cInput[start:len(cInput)])
    return args
}

func run_command_loop() {
    var input string
    var arguments []string

    reader := bufio.NewReader(os.Stdin)

    for true {
        ch, _ := reader.ReadByte()

        if ch == '\n' {
            if strings.TrimSpace(input) == "exit" {
                fmt.Println("Exiting shell...")
                os.Exit(0)
            }
            arguments = parseArgs(input)
            args := arguments[1:]
            cmd := exec.Command(arguments[0], args...)

            cmd.Stdin = os.Stdin
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr

            err := cmd.Run() 
            if err != nil {
                if exitError, ok := err.(*exec.ExitError); ok {
                    fmt.Fprintf(os.Stderr, "Command exited with error: %v (Exit Code: %d)\n", err, exitError.ExitCode())
                } else {
                    fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
                }
            }
            fmt.Print("> ")
            input = ""
            fmt.Printf("The input should be zero : %s ?", input)

		}
        input += string(ch)
        fmt.Printf("The input is increasing: %s ", input)
    }
}

func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT)

    fmt.Print("> ")
    run_command_loop()
    select {
    case sig := <-sigChan:
        fmt.Printf("Signal received: %v\n", sig)
        fmt.Println("Performing cleanup...")
        os.Exit(0)
    }
}
