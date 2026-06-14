package main

import (
	"fmt"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/handlers" // 1. Importamos la nueva carpeta de handlers
)

func main() {
	database.ConectarDB()

	// Rutas de la API test
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API del Inventario Nacional iniciada y estructurada correctamente.")
	})

	// rutas para activos
	http.HandleFunc("/api/activos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListarActivos(w, r) // Llamamos a la función para listar activos
		case http.MethodPost:
			handlers.InsertarActivo(w, r) // Llamamos a la función para insertar un nuevo activo
		case http.MethodPut:
			handlers.ActualizarActivo(w, r) // Llamamos a la función para actualizar un activo existente
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
