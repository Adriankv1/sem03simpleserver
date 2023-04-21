package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "strconv"
    "strings"
    "sync"

    "github.com/Adriankv1/is105sem03/mycrypt"
    //"github.com/Adriankv1/misc/conv"
)

func CelsiusToFahrenheit(value float64) float64 { // Celsius to Fahrenheit
	return (value * 9 / 5) + 32
}

func main() {
    var wg sync.WaitGroup

    server, err := net.Listen("tcp", "172.17.0.3:")
    if err != nil {
        log.Fatal(err)
    }
    defer server.Close()

    log.Printf("bundet til %s", server.Addr().String())

    wg.Add(1)
    go func() {
        defer wg.Done()

        for {
            log.Println("før server.Accept() kallet")

            conn, err := server.Accept()
            if err != nil {
                log.Println(err)
                continue
            }

            wg.Add(1)
            go func(c net.Conn) {
                defer wg.Done()
                defer conn.Close()

                for {
                    buf := make([]byte, 1024)
                    n, err := c.Read(buf)
                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return
                    }
                    dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                    log.Println("Dekrypter melding: ", string(dekryptertMelding))

                    msgString := string(dekryptertMelding)

                    switch msgString {
                    case "ping":
                        kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
                        log.Println("Kryptert melding: ", string(kryptertMelding))
                        _, err = c.Write([]byte(string(kryptertMelding)))
                    case "Kjevik":
                        parts := strings.Split(msgString, " ")
                        celsius, _ := strconv.ParseFloat(parts[1], 64)
                        fahrenheit := CelsiusToFarenheit(celsius)
                        kryptertMelding := mycrypt.Krypter([]rune(fmt.Sprintf("Temperaturen i Fahrenheit er %.2f", fahrenheit)), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                        log.Println("Kryptert melding: ", string(kryptertMelding))
                        _, err = c.Write([]byte(string(kryptertMelding)))
                    default:
                        _, err = c.Write(buf[:n])
                    }

                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return
                    }
                }
            }(conn)
        }
    }()

    wg.Wait()
}
