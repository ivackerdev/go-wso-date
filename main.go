package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir conexiones de cualquier origen
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Estructura para almacenar la diferencia de tiempo
type TimeDifference struct {
	Days    float64 `json:"days"`
	Hours   float64 `json:"hours"`
	Minutes float64 `json:"minutes"`
	Seconds float64 `json:"seconds"`
}

// Función para calcular la diferencia de tiempo entre dos fechas
func getTimeDifference(now, future time.Time) TimeDifference {
	diff := future.Sub(now)

	// Convertir la diferencia a días, horas, minutos y segundos
	totalSeconds := int(diff.Seconds())
	days := totalSeconds / (24 * 3600)
	totalSeconds -= days * 24 * 3600
	hours := totalSeconds / 3600
	totalSeconds -= hours * 3600
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return TimeDifference{
		Days:    float64(days),
		Hours:   float64(hours),
		Minutes: float64(minutes),
		Seconds: float64(seconds),
	}
}

// Función para manejar las conexiones WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Mejora la solicitud HTTP a una conexión WebSocket.
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close() // Asegura que la conexión se cierre al finalizar la función.

	// Crear una fecha fija para el 28-07-2061
	future := time.Date(2061, 8, 1, 0, 0, 0, 0, time.UTC)

	ticker := time.NewTicker(1 * time.Second) // Crea un ticker que emite eventos cada segundo.
	defer ticker.Stop()                       // Asegura que el ticker se detenga para liberar recursos.

	for range ticker.C {
		now := time.Now() // Actualiza la fecha y hora actual en cada iteración.

		// Calcular la diferencia de tiempo
		diff := getTimeDifference(now, future)

		// Convertir la diferencia de tiempo a JSON
		jsonDiff, err := json.Marshal(diff)
		if err != nil {
			log.Println("JSON Marshal error:", err)
			break // Si hay un error al convertir a JSON, sale del bucle.
		}

		// Enviar el JSON al cliente WebSocket
		if err := conn.WriteMessage(websocket.TextMessage, jsonDiff); err != nil {
			log.Println("WriteMessage error:", err)
			break // Si hay un error al enviar el mensaje, sale del bucle.
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServeTLS(":443", "C:\\Certbot\\live\\lab.ivacker.dev\\cert.pem", "C:\\Certbot\\live\\lab.ivacker.dev\\privkey.pem", nil))
}
