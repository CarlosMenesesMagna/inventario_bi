package models

type Asignacion struct {
	ID              int     `json:"id"`
	ActivoID        int     `json:"activo_id"`
	UsuarioID       *int    `json:"usuario_id"` // Puntero * porque puede ser NULL (Bodega)
	SitioID         int     `json:"sitio_id"`
	UbicacionFisica string  `json:"ubicacion_fisica"` // Ej: 'Taller mecánico', 'Garita 3'
	FechaEntrega    string  `json:"fecha_entrega"`
	FechaDevolucion *string `json:"fecha_devolucion"` // Puntero * porque arranca siendo NULL
}
