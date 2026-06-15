package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

// ListarModelos devuelve el catálogo de equipos soportados
func ListarModelos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, tipo, modelo FROM modelos_hardware")
	if err != nil {
		log.Println("Error al consultar modelos: ", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var modelos []models.ModeloHardware
	for rows.Next() {
		var m models.ModeloHardware
		if err := rows.Scan(&m.ID, &m.Tipo, &m.Modelo); err != nil {
			continue
		}
		modelos = append(modelos, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modelos)
}

// InsertarModelo agrega un nuevo tipo de equipo al catálogo
func InsertarModelo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var nuevoModelo models.ModeloHardware
	if err := json.NewDecoder(r.Body).Decode(&nuevoModelo); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if nuevoModelo.Tipo == "" || nuevoModelo.Modelo == "" {
		http.Error(w, "El tipo (ej: Laptop) y el modelo (ej: Latitude 5440) son obligatorios", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO modelos_hardware (tipo, modelo) VALUES ($1, $2) RETURNING id`
	err := database.DB.QueryRow(query, nuevoModelo.Tipo, nuevoModelo.Modelo).Scan(&nuevoModelo.ID)
	if err != nil {
		log.Println("Error al insertar modelo: ", err)
		http.Error(w, "Error al guardar el modelo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Modelo registrado con éxito",
		"modelo":  nuevoModelo,
	})
}
