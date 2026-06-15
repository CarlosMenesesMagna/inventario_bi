package models

// --- ESTRUCTURAS DE DATOS (Modelos) ---

type Sitio struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type EmpresaContratista struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Usuario struct {
	ID             int    `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
	Email          string `json:"email"`
	EmpresaID      int    `json:"empresa_id"`
}

type ModeloHardware struct {
	ID     int    `json:"id"`
	Tipo   string `json:"tipo"`
	Modelo string `json:"modelo"`
}

type Activo struct {
	ID             int    `json:"id"`
	Serial         string `json:"serial"`
	HostName       string `json:"host_name"`
	ModeloID       int    `json:"modelo_id"`
	AnioCompra     string `json:"anio_compra"`
	AnioRenovacion string `json:"anio_renovacion"`
	CicloDeVida    string `json:"ciclo_de_vida"`
	Disposicion    string `json:"disposicion"`
	Status         string `json:"status"`
	NotasTecnicas  string `json:"notas_tecnicas"`
}
