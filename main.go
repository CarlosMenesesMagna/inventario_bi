package main

import (
	"fmt"
	"log"
	"net/http"

	"inventario_bi/database"
	"inventario_bi/handlers"
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
			handlers.ListarActivos(w, r) // función para listar activos
		case http.MethodPost:
			handlers.InsertarActivo(w, r) // función para insertar un nuevo activo
		case http.MethodPut:
			handlers.ActualizarActivo(w, r) // función para actualizar un activo existente
		case http.MethodDelete:
			handlers.EliminarActivo(w, r) // función para eliminar un activo
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// rutas para asignaciones
	http.HandleFunc("/api/asignaciones", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.InsertarAsignacion(w, r) // función para insertar una nueva asignación
		case http.MethodGet:
			handlers.ListarAsignaciones(w, r) // función para listar asignaciones
		case http.MethodPut:
			handlers.ActualizarAsignacion(w, r) // función para actualizar una asignación existente
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	//rutas api sitios
	http.HandleFunc("/api/sitios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListarSitios(w, r) // función para listar sitios
		case http.MethodPost:
			handlers.InsertarSitio(w, r) // función para insertar un nuevo sitio
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	//ruta api/usuarios
	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListarUsuarios(w, r) // función para listar usuarios
		case http.MethodPost:
			handlers.InsertarUsuario(w, r) // función para insertar un nuevo usuario
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	//ruta api/modelos
	http.HandleFunc("/api/modelos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListarModelos(w, r) // función para listar modelos
		case http.MethodPost:
			handlers.InsertarModelo(w, r) // función para insertar un nuevo modelo
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	//ruta api/empresas
	http.HandleFunc("/api/empresas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListarEmpresas(w, r) // función para listar empresas
		case http.MethodPost:
			handlers.InsertarEmpresa(w, r) // función para insertar una nueva empresa
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Println("Servidor corriendo en http://localhost:8080 con CORS habilitado")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

// enableCORS es un Middleware que agrega las cabeceras necesarias para permitir
// que un frontend (como React) se conecte sin ser bloqueado por el navegador.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Le decimos al navegador: "Sí, acepto peticiones desde cualquier origen (*)"
		// En producción estricta, en vez de "*" pondrías "http://localhost:5173"
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Los navegadores envían una petición fantasma llamada "OPTIONS" antes del POST/PUT/DELETE
		// para preguntar si tienen permiso. Si es OPTIONS, le decimos OK y cortamos ahí.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Si todo está bien, dejamos que la petición siga su camino hacia tus handlers
		next.ServeHTTP(w, r)
	})
}
