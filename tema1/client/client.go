package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func loadConfig(filename string) (*config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <input_file>")
		return
	}
	filename := os.Args[1]

	cfg, err := loadConfig("tema1/server/config.json")
	if err != nil {
		fmt.Println("Eroare la citirea config.json în client:", err)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var name string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name = line
		break
	}

	if name == "" {
		fmt.Println("Fișierul nu conține un nume de client pe prima linie.")
		return
	}
	fmt.Println("Nume client citit din fișier:", name)

	//addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", name)

	reader := bufio.NewReader(conn)
	enc := json.NewEncoder(conn)

	reply, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
		return
	}
	reply = strings.TrimSpace(reply)
	fmt.Println("Server response:", reply)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		args := strings.Fields(line)
		fmt.Println("> Comanda din fișier:", args)
		req := Request{
			Name: name,
			Args: args,
		}

		if err := enc.Encode(req); err != nil {
			fmt.Printf("Error sending request to server: %v\n", err)
			return
		}

		reply1, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
			return
		}
		reply1 = strings.TrimSpace(reply1)
		fmt.Println("Server:", reply1)

		reply2, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Pierdut conexiune (la citirea rezultatului).")
			return
		}
		reply2 = strings.TrimSpace(reply2)
		fmt.Println("Server:", reply2)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Eroare la citirea din fișier:", err)
	}

}
