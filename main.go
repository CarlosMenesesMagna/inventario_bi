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

	// Rutas de la API
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API del Inventario Nacional iniciada y estructurada correctamente.")
	})

	// 2. Registramos la nueva URL para listar activos
	http.HandleFunc("/api/activos", handlers.ListarActivos)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
