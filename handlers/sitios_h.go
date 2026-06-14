package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

// ListarSitios devuelve todas las plantas u oficinas registradas
func ListarSitios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, nombre FROM sitios")
	if err != nil {
		log.Println("Error al consultar sitios: ", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sitios []models.Sitio
	for rows.Next() {
		var s models.Sitio
		if err := rows.Scan(&s.ID, &s.Nombre); err != nil {
			continue
		}
		sitios = append(sitios, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sitios)
}

// InsertarSitio permite agregar una nueva planta (Ej: Tocopilla, Laguna Verde)
func InsertarSitio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var nuevoSitio models.Sitio
	if err := json.NewDecoder(r.Body).Decode(&nuevoSitio); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if nuevoSitio.Nombre == "" {
		http.Error(w, "El nombre del sitio es obligatorio", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO sitios (nombre) VALUES ($1) RETURNING id`
	err := database.DB.QueryRow(query, nuevoSitio.Nombre).Scan(&nuevoSitio.ID)
	if err != nil {
		log.Println("Error al insertar sitio: ", err)
		http.Error(w, "Error al guardar el sitio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Sitio registrado con éxito",
		"sitio":   nuevoSitio,
	})
}
