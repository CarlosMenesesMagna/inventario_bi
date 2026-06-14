package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

// ListarUsuarios devuelve todos los técnicos y empleados registrados
func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) //validar metodo
		return
	}
	//realizar query a la base de datos para obtener usuarios
	rows, err := database.DB.Query("SELECT id, nombre_completo, email, empresa_id FROM usuarios")
	if err != nil {
		log.Println("Error al consultar usuarios: ", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	//variable que almacenará los usuarios obtenidos de la base de datos
	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.NombreCompleto, &u.Email, &u.EmpresaID); err != nil {
			continue
		}
		usuarios = append(usuarios, u)
	}
	//Devolver el JSON con la lista de usuarios al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// InsertarUsuario agrega una nueva persona al sistema
func InsertarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) //metodo permitido
		return
	}
	//variable que almacenará los datos del nuevo usuario enviados por el cliente
	var nuevoUsuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&nuevoUsuario); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Validación exigimos Nombre y Empresa (el ID 1 que es Vector, por ejemplo).
	// No exigimos correo, por si en el Excel dice "N/A".
	if nuevoUsuario.NombreCompleto == "" || nuevoUsuario.EmpresaID == 0 {
		http.Error(w, "El nombre completo y la empresa son obligatorios", http.StatusBadRequest)
		return
	}
	// Insertar el nuevo usuario en la base de datos y obtener su ID generado
	query := `INSERT INTO usuarios (nombre_completo, email, empresa_id) VALUES ($1, $2, $3) RETURNING id`
	err := database.DB.QueryRow(query, nuevoUsuario.NombreCompleto, nuevoUsuario.Email, nuevoUsuario.EmpresaID).Scan(&nuevoUsuario.ID)
	if err != nil {
		log.Println("Error al insertar usuario: ", err)
		http.Error(w, "Error al guardar el usuario", http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el nuevo usuario creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Usuario registrado con éxito",
		"usuario": nuevoUsuario,
	})
}
