package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/models"
)

// registra a quien y donde se le entrego el equipo
func InsertarAsignacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "metodo no permitido, debe ser POST", http.StatusMethodNotAllowed)
		return
	}

	var nuevaAsignacion models.Asignacion                   //variable con el modelo de datos de asignacion
	err := json.NewDecoder(r.Body).Decode(&nuevaAsignacion) //se guarda en err para luego validar si hubo un error al decodificar el json
	if err != nil {
		log.Println("Error al decodificar la solicitud de asignación", err)
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	//Validacion equipo y sitio
	if nuevaAsignacion.ActivoID == 0 || nuevaAsignacion.SitioID == 0 {
		http.Error(w, "ActivoID y SitioID son obligatorios", http.StatusBadRequest)
		return
	}

	//Insertar en la base de datos
	query := `
		INSERT INTO asignaciones (activo_id, usuario_id, sitio_id, ubicacion_fisica, fecha_entrega)
		VALUES ($1, $2, $3, $4, CURRENT_DATE)
		RETURNING id
	`
	err = database.DB.QueryRow(
		query,
		nuevaAsignacion.ActivoID,
		nuevaAsignacion.UsuarioID,
		nuevaAsignacion.SitioID,
		nuevaAsignacion.UbicacionFisica,
	).Scan(&nuevaAsignacion.ID) //obtener el ID generado por la base de datos

	if err != nil {
		log.Println("Error al insertar la asignación en la base de datos", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	//Responder al cliente con la asignación creada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":    "Equipo Asignado con exito",
		"asignacion": nuevaAsignacion,
	})
}

// funcion para listar asignaciones
func ListarAsignaciones(w http.ResponseWriter, r *http.Request) {
	//Validar método http
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido, debe ser GET", http.StatusMethodNotAllowed)
		return
	}

	//Consultar en la base de datos todas las asignaciones
	rows, err := database.DB.Query(`
		SELECT id, activo_id, usuario_id, sitio_id, ubicacion_fisica, fecha_entrega, fecha_devolucion 
		FROM asignaciones
	`)
	//validar si hubo un error al consultar la base de datos
	if err != nil {
		log.Println("Error al consultar asignaciones: ", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	defer rows.Close() //cerrar la consulta al finalizar la función para liberar recursos

	//variable que almacenará las asignaciones obtenidas de la base de datos
	var asignaciones []models.Asignacion

	//recorre la variable rows obtenida de la consulta y escanea cada fila en una variable de tipo Asignacion, luego la agrega al slice de asignaciones
	for rows.Next() {
		var a models.Asignacion
		err := rows.Scan(
			&a.ID,
			&a.ActivoID,
			&a.UsuarioID,
			&a.SitioID,
			&a.UbicacionFisica,
			&a.FechaEntrega,
			&a.FechaDevolucion,
		)
		if err != nil {
			log.Println("Error al escanear fila de asignaciones: ", err)
			continue
		}
		asignaciones = append(asignaciones, a)
	}

	//Devolver el JSON con la lista de asignaciones al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asignaciones)
}

func ActualizarAsignacion(w http.ResponseWriter, r *http.Request) {
	//obtener id de asignacion
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de asignación es requerido", http.StatusBadRequest)
		return
	}

	//ejecutar el update en bd para marcar la fecha de devolución
	query := `
		UPDATE asignaciones
		SET fecha_devolucion = CURRENT_DATE
		WHERE id = $1 AND fecha_devolucion IS NULL
	`
	resultado, err := database.DB.Exec(query, idStr)
	if err != nil {
		log.Println("Error al registrar la devolución en la BD: ", err)
		http.Error(w, "Error interno al procesar la devolución", http.StatusInternalServerError)
		return
	}

	//validar si realmente se modifico la fila, puede que el id no existiera o que ya tuviera fecha de devolución
	filasAfectadas, err := resultado.RowsAffected()
	if filasAfectadas == 0 {
		http.Error(w, "La asignación no existe o el equipo ya había sido devuelto antes", http.StatusNotFound)
		return
	}
	//respuesta de exito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Devolución registrada con éxito",
	})
}
