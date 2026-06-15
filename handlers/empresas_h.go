package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

// ListarEmpresas devuelve la lista de empresas contratistas
func ListarEmpresas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, nombre FROM empresas_contratistas")
	if err != nil {
		log.Println("Error al consultar empresas: ", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var empresas []models.EmpresaContratista
	for rows.Next() {
		var e models.EmpresaContratista
		if err := rows.Scan(&e.ID, &e.Nombre); err != nil {
			continue
		}
		empresas = append(empresas, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(empresas)
}

// InsertarEmpresa agrega una nueva empresa al catálogo
func InsertarEmpresa(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var nuevaEmpresa models.EmpresaContratista
	if err := json.NewDecoder(r.Body).Decode(&nuevaEmpresa); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if nuevaEmpresa.Nombre == "" {
		http.Error(w, "El nombre de la empresa es obligatorio", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO empresas_contratistas (nombre) VALUES ($1) RETURNING id`
	err := database.DB.QueryRow(query, nuevaEmpresa.Nombre).Scan(&nuevaEmpresa.ID)
	if err != nil {
		log.Println("Error al insertar empresa: ", err)
		http.Error(w, "Error al guardar la empresa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Empresa registrada con éxito",
		"empresa": nuevaEmpresa,
	})
}
