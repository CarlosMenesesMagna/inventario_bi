package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

func ListarActivos(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT id, serial, host_name, modelo_id, anio_compra, anio_renovacion, ciclo_de_vida, disposicion, status, notas_tecnicas FROM activos`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Error al consultar activos", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return

	}
	defer rows.Close() //cerrar filas

	//crear lista dinamica donde se guarden los activos que encontramos
	var activos []models.Activo = []models.Activo{}

	//recorrer resultados fila por fila
	for rows.Next() {
		var a models.Activo
		//scan copia las variables de la fila actual hacia el struct
		err := rows.Scan(&a.ID, &a.Serial, &a.HostName, &a.ModeloID, &a.AnioCompra, &a.AnioRenovacion, &a.CicloDeVida, &a.Disposicion, &a.Status, &a.NotasTecnicas)
		if err != nil {
			log.Println("Error al escanear fila de activos", err)
			continue
		}
		activos = append(activos, a)
	}

	//config cabecera indicando que el dato que se enviará es un json

	w.Header().Set("content-type", "application/json")

	//transformar slice go en json y enviar al cliente / navegador

	json.NewEncoder(w).Encode(activos)

}

// insertar un nuevo activo a la base de datos
func InsertarActivo(w http.ResponseWriter, r *http.Request) {
	//validar que el metodo sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "metodo no permitido, debe ser POST", http.StatusMethodNotAllowed)
		return
	}

	//Crear molde vacio de activo para recibir los datos del cliente
	var nuevoActivo models.Activo

	//leer el json que viene de la peticion y copiarlo al molde
	err := json.NewDecoder(r.Body).Decode(&nuevoActivo)
	if err != nil {
		log.Println("Error al decodificar el cuerpo de la solicitud", err)
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	//preparar consulta sql se usan $1 $2... para evitar inyeccion sql
	//nota; no pasamos ID por que Postgres lo genera automaticamente
	query := `INSERT INTO activos (serial, host_name, modelo_id, anio_compra, anio_renovacion, ciclo_de_vida, disposicion, status, notas_tecnicas) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
				RETURNING id`

	//ejecutar consulta pasando los valores del struct nuevoActivo
	err = database.DB.QueryRow(
		query,
		nuevoActivo.Serial,
		nuevoActivo.HostName,
		nuevoActivo.ModeloID,
		nuevoActivo.AnioCompra,
		nuevoActivo.AnioRenovacion,
		nuevoActivo.CicloDeVida,
		nuevoActivo.Disposicion,
		nuevoActivo.Status,
		nuevoActivo.NotasTecnicas,
	).Scan(&nuevoActivo.ID) //obtener el ID generado por la base de datos

	if err != nil {
		log.Println("Error al insertar nuevo activo", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if nuevoActivo.Serial == "" || nuevoActivo.ModeloID == 0 {
		http.Error(w, "Los campos serial y modelo son obligatorios", http.StatusBadRequest)
		return
	}

	//responderle al usuario con un json indicando el exito
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Equipo Registrado con exito",
		"activo":  nuevoActivo,
	})

}
