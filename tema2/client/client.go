package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type ConfigData struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	file, err := os.Open("./tema2/server/config.json")
	if err != nil {
		fmt.Println("Lipsa config.json")
		return
	}
	var cfg ConfigData
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		fmt.Println("Eroare JSON:", err)
		return
	}
	file.Close()

	//conn
	address := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Eroare conectare:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	//Handshake
	fmt.Print("Numele tau: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	conn.Write([]byte(name + "\n"))

	reply, _ := serverReader.ReadString('\n')
	fmt.Print("Server: " + reply)

	for {
		fmt.Println("\n--- MENIU ---")
		fmt.Println("ex3, ex9, ex11, exit")
		fmt.Print(">> ")

		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" {
			break
		}
		if cmd != "ex3" && cmd != "ex9" && cmd != "ex11" {
			fmt.Println("Comanda invalida.")
			continue
		}

		// Citire date de la tastatura
		var bufferToSend []string
		bufferToSend = append(bufferToSend, cmd) //Prima linie e comanda

		fmt.Print("Nr linii: ")
		var n int
		fmt.Scanf("%d\n", &n)

		for i := 0; i < n; i++ {
			fmt.Printf("Linia %d: ", i+1)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			bufferToSend = append(bufferToSend, text)
		}

		for _, line := range bufferToSend {
			conn.Write([]byte(line + "\n"))
		}

		conn.Write([]byte("\n"))

		// Asteptam raspunsul
		response, _ := serverReader.ReadString('\n')
		fmt.Print("Rezultat: " + response)
	}
}
