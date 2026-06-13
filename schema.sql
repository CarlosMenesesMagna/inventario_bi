-- 1. ENUMS (Filtros cerrados para evitar errores de tipeo)
CREATE TYPE tipo_equipo AS ENUM ('Laptop', 'Desktop');
CREATE TYPE status_activo AS ENUM ('Activo', 'Bodega', 'Baja', 'Mantenimiento', 'Malo', 'Robado');
CREATE TYPE ciclo_vida_activo AS ENUM ('DENTRO', 'FUERA');
CREATE TYPE disposicion_activo AS ENUM ('EN USO', 'VENTA', 'DESECHO', 'DEVOLUCION', 'ELIMINAR AF', 'OTRO', 'POOL');

-- 2. TABLAS MAESTRAS (Para no repetir marcas, modelos ni empresas de contratistas)
CREATE TABLE sitios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE -- Ej: 'Angamos', 'Cochrane', 'Bolero'
);

CREATE TABLE empresas_contratistas (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE -- Ej: 'AES Andes', 'Vector', 'Axinntus', 'Symmetric'
);

CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nombre_completo VARCHAR(150) NOT NULL,
    email VARCHAR(150),
    empresa_id INT REFERENCES empresas_contratistas(id)
);

CREATE TABLE modelos_hardware (
    id SERIAL PRIMARY KEY,
    tipo tipo_equipo NOT NULL,       -- Solo 'Laptop' o 'Desktop'
    modelo VARCHAR(100) NOT NULL UNIQUE -- Ej: 'Latitude 5440', 'Optiplex 7050'
);

-- 3. LA TABLA DURO: LOS ASSETS (La máquina física)
CREATE TABLE activos (
    id SERIAL PRIMARY KEY,
    serial VARCHAR(100) NOT NULL,
    host_name VARCHAR(100) NOT NULL, 
    modelo_id INT NOT NULL REFERENCES modelos_hardware(id),
    anio_compra VARCHAR(20) NOT NULL,       -- Soporta números o 'Externo'
    anio_renovacion VARCHAR(20) NOT NULL,   -- Soporta números o 'NO APLICA'
    ciclo_de_vida ciclo_vida_activo NOT NULL,
    disposicion disposicion_activo NOT NULL,
    status status_activo NOT NULL,
    notas_tecnicas TEXT                     -- Para 'Falla placa', 'Pantalla rota', etc.
);

-- 4. EL HISTORIAL (Quién lo tiene hoy y dónde está físicamente)
CREATE TABLE asignaciones (
    id SERIAL PRIMARY KEY,
    activo_id INT NOT NULL REFERENCES activos(id),
    usuario_id INT REFERENCES usuarios(id), -- Puede ser NULL si quedó en bodega
    sitio_id INT NOT NULL REFERENCES sitios(id),
    ubicacion_fisica VARCHAR(150),          -- Ej: 'Taller mecánico', 'Garita CCTV'
    fecha_entrega DATE DEFAULT CURRENT_DATE,
    fecha_devolucion DATE
);