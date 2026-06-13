package database

import (
	"database/sql"
	"fmt"
	"log"

	// Driver de PostgreSQL
	_ "github.com/lib/pq"
)

// DB es una variable global que mantiene la conexión viva
// Empieza con mayúscula para poder usarla en otras partes del programa
var DB *sql.DB

// ConectarDB abre la conexión con PostgreSQL en Docker
func ConectarDB() {
	// Cadena de conexión a tu Docker local
	connStr := "postgres://admin:superpassword123@localhost:5432/inventario_db?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al configurar la conexión a la base de datos: ", err)
	}

	// Verificar que realmente hay comunicación con un Ping
	err = DB.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos. ¿Docker está corriendo?: ", err)
	}

	fmt.Println("¡Conexión exitosa a PostgreSQL desde el módulo database! 🚀")
}
