package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

type config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type Request struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type Response struct {
	Result []string `json:"result"`
	Error  string   `json:"error"`
}

func loadConfig(filename string) (*config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func handleConnection(conn net.Conn, cfg *config) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Printf("Client %s conectat.\n", name)
	conn.Write([]byte("Conectat\n"))

	remoteAddr := conn.RemoteAddr().String()
	dec := json.NewDecoder(reader)

	for {
		var req Request

		if err := dec.Decode(&req); err != nil {
			log.Printf("Conexiune %s închisă sau eroare la citire: %v\n", remoteAddr, err)
			return
		}

		if len(req.Args) == 0 {
			log.Printf("Client %s a trimis args gol.\n", name)
			conn.Write([]byte("nicio comanda primita\n"))
			continue
		}

		fmt.Printf("Client %s a facut request cu datele: %v\n", req.Name, req.Args)
		fmt.Println("Server a primit requestul.")
		fmt.Println("Server proceseaza datele.")

		conn.Write([]byte("primit request\n"))

		result := rezolva(req.Args)

		fmt.Println("Client " + name + " a primit răspunsul: " + result)

		conn.Write([]byte(result + "\n"))
	}
}

func rezolva(args []string) string {
	switch args[0] {
	case "ex1":
		return ex1(args)
	case "ex2":
		return ex2(args)
	case "ex3":
		return ex3(args)
	case "ex5":
		return ex5(args)
	case "ex11":
		return ex11(args)
	default:
		return "ex necunoscut"
	}
}

func ex1(args []string) string {
	if len(args) < 2 {
		return "eroare: numar insuficient de argumente pentru ex1"
	}

	words := args[1:]

	firstLen := len(words[0])
	for _, word := range words {
		if len(word) != firstLen {
			return "eroare: toate cuvintele trebuie sa aiba aceeasi lungime"
		}
	}

	out := make([]string, firstLen)
	for i := 0; i < firstLen; i++ {
		var sb strings.Builder
		for _, word := range words {
			sb.WriteByte(word[i])
		}
		out[i] = sb.String()
	}

	//casa, masa, trei, tanc, 4321 => cmtt4, aara3, ssen2, aaic1
	return strings.Join(out, " ")
}

func ex2(args []string) string {
	if len(args) < 2 {
		return "eroare: numar insuficient de argumente pentru ex2"
	}

	count := 0

	for _, word := range args[1:] {
		var digitsBuilder strings.Builder
		for _, ch := range word {
			if ch >= '0' && ch <= '9' {
				digitsBuilder.WriteRune(ch)
			}
		}

		digitsStr := digitsBuilder.String()
		if digitsStr == "" {
			continue
		}

		n, err := strconv.Atoi(digitsStr)
		if err != nil {
			continue
		}

		root := int(math.Sqrt(float64(n)))
		if root*root == n {
			count++
		}
	}
	return fmt.Sprintf("Numar de cuvinte cu numere patrate: %d", count)
}

func ex3(args []string) string {
	if len(args) < 2 {
		return "eroare: numar insuficient de argumente pentru ex3"
	}

	sum := 0

	for _, s := range args[1:] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return "eroare: toate argumentele trebuie sa fie numere intregi"
		}
		sum += reverseInt(n)
	}
	return fmt.Sprintf("Suma numerelor inversate este: %d", sum)
}

func reverseInt(n int) int {

	if n == 0 {
		return 0
	}

	if n < 0 {
		return -reverseInt(-n)
	}

	rev := 0
	for n > 0 {
		rev = rev*10 + n%10
		n /= 10
	}
	return rev
}

func ex5(args []string) string {
	return "rezultat ex5 (de implementat)"
}

func ex11(args []string) string {
	return "rezultat ex11 (de implementat)"
}

func main() {
	configData, err := loadConfig("tema1/server/config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", configData.Host, configData.Port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer ln.Close()
	fmt.Printf("Server listening on %s:%d\n", configData.Host, configData.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, configData)
	}
}
