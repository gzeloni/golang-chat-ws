package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8080/ws" // Endereço do servidor WebSocket

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Digite seu nome de usuário: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Erro ao ler mensagem do servidor:", err)
				return
			}

			if !strings.HasPrefix(string(message), username+": ") {
				fmt.Printf("%s\n", message)
				fmt.Print("> ")
			}
		}
	}()

	fmt.Printf("Conectado ao servidor como '%s'. Digite sua mensagem ou 'sair' para sair.\n", username)

	for {
		fmt.Print("> ")
		scanner.Scan()
		text := scanner.Text()
		if text == "sair" {
			break
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(username+": "+text))
		if err != nil {
			fmt.Println("Erro ao enviar mensagem para o servidor:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler a entrada do usuário:", err)
	}

	fmt.Println("Desconectando do servidor.")
}
