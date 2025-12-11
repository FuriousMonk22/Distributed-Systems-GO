package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

type ConfigData struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type KeyValue struct {
	Key   string
	Value int
}

func main() {
	//Config
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

	//Server
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server pornit pe", address)

	//Accept conn
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go connect(conn)
	}
}

func connect(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	//Handshake
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println(name, "s-a conectat.")
	conn.Write([]byte("Conectat.\n"))

	//Procesare cereri
	for {
		var lines []string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(name, "s-a deconectat.")
				return
			}
			line = strings.TrimSpace(line)

			if line == "" {
				break
			}
			lines = append(lines, line)
		}

		if len(lines) == 0 {
			continue
		}

		fmt.Println("Request de la", name, ":", lines[0])

		// Trimitem raspunsul inapoi
		rsp := rezolva(lines)
		conn.Write([]byte(rsp + "\n"))
	}
}

func rezolva(args []string) string {
	if len(args) == 0 {
		return "Eroare: date lipsa"
	}
	cmd := args[0]
	data := args[1:]

	if len(data) == 0 {
		return "0.00"
	}

	switch cmd {
	case "ex3":
		return ex3(data)
	case "ex9":
		return ex9(data)
	case "ex11":
		return ex11(data)
	}
	return "Comanda necunoscuta"
}

func Map(document string, valuefunc func(string) int) []KeyValue {
	words := strings.Fields(document)
	acc := make(map[string]int)
	for _, word := range words {
		acc[word] += valuefunc(word)
	}
	keyValues := make([]KeyValue, 0, len(acc))
	for k, v := range acc {
		keyValues = append(keyValues, KeyValue{Key: k, Value: v})
	}
	return keyValues
}

func Reduce(keyValues []KeyValue) int {
	sum := 0
	for _, kv := range keyValues {
		sum += kv.Value
	}
	return sum
}

// ex3
func ex3Func(s string) int {
	runes := []rune(s)
	for i, r := range runes {
		//daca e vocala
		if strings.ContainsRune("aeiouAEIOU", r) {
			//daca exista
			if i+1 < len(runes) {
				//daca nu e p, nu e pasareasca
				if strings.ToLower(string(runes[i+1])) != "p" {
					return 0
				}
			}
		}
	}
	return 1
}

func ex3(lines []string) string {
	results := make(chan int, len(lines))
	for _, line := range lines {
		go func(l string) {
			results <- Reduce(Map(l, ex3Func))
		}(line)
	}
	return calcAvg(results, len(lines))
}

// ex11
func ex11Func(s string) int {
	u, l, d, sym := false, false, false, false
	for _, r := range s {
		if unicode.IsUpper(r) {
			u = true
		} else if unicode.IsLower(r) {
			l = true
		} else if unicode.IsDigit(r) {
			d = true
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			sym = true
		}
	}
	if u && l && d && sym {
		return 1
	}
	return 0
}

func ex11(lines []string) string {
	results := make(chan int, len(lines))
	for _, line := range lines {
		go func(l string) {
			results <- Reduce(Map(l, ex11Func))
		}(line)
	}
	return calcAvg(results, len(lines))
}

// ex9
func ex9Map(document string) []KeyValue {
	words := strings.Fields(document)
	buckets := make(map[string]int)

	// suf pus in buckets
	for _, w := range words {
		if len(w) >= 2 {
			suf := w[len(w)-2:]
			buckets[suf]++
		}
	}

	// nr perechile din bucket
	totalPairs := 0
	for _, count := range buckets {
		if count%2 == 0 {
			totalPairs += count / 2
			continue
		}
	}

	//trimitem total pairs
	return []KeyValue{{Key: "perechi", Value: totalPairs}}
}

func ex9(lines []string) string {
	results := make(chan int, len(lines))
	for _, line := range lines {
		go func(l string) {
			results <- Reduce(ex9Map(l))
		}(line)
	}
	return calcAvg(results, len(lines))
}

func calcAvg(results chan int, n int) string {
	sum := 0
	for i := 0; i < n; i++ {
		sum += <-results
	}
	avg := float64(sum) / float64(n)
	return fmt.Sprintf("%.2f", avg)
}

//go run tema2/server/server.go
//go run tema2/client/client.go

//apap paprc apnap mipnipm copil
//cepr program lepu zepcep golang tema
//par impar papap gepr

//stele mele borcan vajnic straşnic
//crocodil garnisit muşețel făurit arhanghel noapte
//lampă sine cine torişte

//sadsa1@A cevaA!4 saar aaastrfb
//aaabbbccc !Caporal1 ddanube jahfjksgfjhs ajsdas urs
//scoica Coral!@12 arac karnak
