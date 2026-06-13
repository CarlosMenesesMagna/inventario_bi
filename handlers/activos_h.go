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

	query := `SELECT id, serial, host_name, modelo_id, anio_compra, anio_renovacion, ciclo_de_vida, disposition, status, notas_tecnicas FROM activos`
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
