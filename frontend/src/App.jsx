import { useState, useEffect } from 'react'
import aesLogo from './assets/AesLogo.png'

function App() {
  // Estados Globales
  const [vistaActiva, setVistaActiva] = useState('dashboard') // 'dashboard' o 'ingresar'
  const [inventario, setInventario] = useState([])
  const [sitios, setSitios] = useState([])
  const [busqueda, setBusqueda] = useState("")
  const [sitioSeleccionado, setSitioSeleccionado] = useState("Todos")
  const [cargando, setCargando] = useState(true)

  // Estados del Formulario de Ingreso
  const [nuevoEquipo, setNuevoEquipo] = useState({
    serial: '', host_name: '', modelo_id: 1, anio_compra: '', anio_renovacion: '', ciclo_de_vida: 'DENTRO', disposicion: 'BODEGA', status: 'Activo', notas_tecnicas: ''
  })
  const [mensajeForm, setMensajeForm] = useState(null)

  // Cargar datos (se ejecuta al abrir y cuando volvemos al dashboard)
  const cargarTodo = async () => {
    setCargando(true)
    try {
      const [resActivos, resAsignaciones, resSitios, resUsuarios] = await Promise.all([
        fetch('http://localhost:8080/api/activos'),
        fetch('http://localhost:8080/api/asignaciones'),
        fetch('http://localhost:8080/api/sitios'),
        fetch('http://localhost:8080/api/usuarios')
      ])

      const activos = await resActivos.json() || []
      const asignaciones = await resAsignaciones.json() || []
      const sitiosBD = await resSitios.json() || []
      const usuarios = await resUsuarios.json() || []

      setSitios(sitiosBD)

      const datosCruzados = activos.map(activo => {
        const asignacionActiva = asignaciones.find(a => a.activo_id === activo.id && a.fecha_devolucion === null)
        let sitioNombre = "Bodega Central"
        let usuarioNombre = "N/A"
        let ubicacion = "Bodega"

        if (asignacionActiva) {
          const sitio = sitiosBD.find(s => s.id === asignacionActiva.sitio_id)
          const usuario = usuarios.find(u => u.id === asignacionActiva.usuario_id)
          sitioNombre = sitio ? sitio.nombre : "Desconocido"
          usuarioNombre = usuario ? usuario.nombre_completo : "N/A"
          ubicacion = asignacionActiva.ubicacion_fisica || "N/A"
        }

        return { ...activo, sitioNombre, usuarioNombre, ubicacion }
      })

      setInventario(datosCruzados)
      setCargando(false)
    } catch (error) {
      console.error("Error conectando al backend:", error)
      setCargando(false)
    }
  }

  useEffect(() => {
    if (vistaActiva === 'dashboard') {
      cargarTodo()
    }
  }, [vistaActiva])

  // Lógica para enviar el formulario a Golang
  const manejarSubmitIngreso = async (e) => {
    e.preventDefault() // Evita que la página se recargue
    try {
      // Convertimos los IDs y años a números para que Golang no reclame
      const payload = {
        ...nuevoEquipo,
        modelo_id: parseInt(nuevoEquipo.modelo_id),
        anio_compra: parseInt(nuevoEquipo.anio_compra),
        anio_renovacion: parseInt(nuevoEquipo.anio_renovacion)
      }

      const respuesta = await fetch('http://localhost:8080/api/activos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })

      if (respuesta.ok) {
        setMensajeForm({ tipo: 'exito', texto: '¡Equipo ingresado correctamente a la base de datos!' })
        // Limpiamos el formulario
        setNuevoEquipo({ serial: '', host_name: '', modelo_id: 1, anio_compra: '', anio_renovacion: '', ciclo_de_vida: 'DENTRO', disposicion: 'BODEGA', status: 'Activo', notas_tecnicas: '' })
      } else {
        setMensajeForm({ tipo: 'error', texto: 'Error al ingresar el equipo. Verifica los datos.' })
      }
    } catch (error) {
      setMensajeForm({ tipo: 'error', texto: 'Fallo la conexión con el servidor.' })
    }
  }

  // Filtro de búsqueda
  const inventarioFiltrado = inventario.filter(equipo => {
    const coincideSitio = sitioSeleccionado === "Todos" || equipo.sitioNombre === sitioSeleccionado
    const termino = busqueda.toLowerCase()
    const coincideBusqueda = 
      equipo.serial.toLowerCase().includes(termino) ||
      equipo.host_name.toLowerCase().includes(termino) ||
      equipo.usuarioNombre.toLowerCase().includes(termino) ||
      equipo.ubicacion.toLowerCase().includes(termino)

    return coincideSitio && coincideBusqueda
  })

  return (
    <div style={{ display: 'flex', height: '100vh', width: '100vw', backgroundColor: '#f4f6f9' }}>
      
      {/* MENU LATERAL */}
      <div style={{ width: '260px', backgroundColor: '#1e293b', color: 'white', display: 'flex', flexDirection: 'column', boxShadow: '4px 0 10px rgba(0,0,0,0.1)', zIndex: 10 }}>
        <div style={{ padding: '20px', backgroundColor: '#ffffff', textAlign: 'center', borderBottom: '4px solid #27b54d' }}>
          <img src={aesLogo} alt="AES Chile Logo" style={{ width: '160px', height: 'auto', marginBottom: '5px' }} />
          <p style={{ margin: 0, fontSize: '13px', color: '#64748b', fontWeight: 'bold', letterSpacing: '1px' }}>INVENTARIO BI</p>
        </div>
        
        <nav style={{ padding: '25px 15px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <button 
            onClick={() => setVistaActiva('dashboard')}
            style={{ backgroundColor: vistaActiva === 'dashboard' ? '#3154f4' : 'transparent', color: vistaActiva === 'dashboard' ? 'white' : '#cbd5e1', border: 'none', padding: '14px', borderRadius: '8px', fontWeight: 'bold', cursor: 'pointer', textAlign: 'left', fontSize: '15px', transition: '0.2s' }}>
            📊 Vista General
          </button>
          <button 
            onClick={() => setVistaActiva('ingresar')}
            style={{ backgroundColor: vistaActiva === 'ingresar' ? '#3154f4' : 'transparent', color: vistaActiva === 'ingresar' ? 'white' : '#cbd5e1', border: 'none', padding: '14px', borderRadius: '8px', fontWeight: 'bold', cursor: 'pointer', textAlign: 'left', transition: '0.2s', fontSize: '15px' }}>
            💻 Ingresar Equipo
          </button>
          <button style={{ backgroundColor: 'transparent', color: '#cbd5e1', border: 'none', padding: '14px', borderRadius: '8px', fontWeight: 'bold', cursor: 'pointer', textAlign: 'left', transition: '0.2s', fontSize: '15px' }}>
            🔄 Asignar / Devolver
          </button>
        </nav>
      </div>

      {/* CONTENIDO DINÁMICO */}
      <div style={{ flex: 1, padding: '40px', overflowY: 'auto' }}>
        
        {/* === VISTA: DASHBOARD === */}
        {vistaActiva === 'dashboard' && (
          <>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
              <div>
                <h1 style={{ color: '#0f172a', margin: '0 0 8px 0', fontSize: '32px' }}>Panel de Activos</h1>
                <p style={{ color: '#64748b', margin: 0, fontSize: '16px' }}>Búsqueda y gestión en tiempo real</p>
              </div>
              <input 
                type="text" 
                placeholder="🔍 Buscar persona, serial o host..." 
                value={busqueda}
                onChange={(e) => setBusqueda(e.target.value)}
                style={{ padding: '14px 20px', width: '380px', borderRadius: '10px', border: '1px solid #cbd5e1', fontSize: '15px', outline: 'none', boxShadow: '0 2px 5px rgba(0,0,0,0.05)' }}
              />
            </div>

            <div style={{ display: 'flex', gap: '12px', marginBottom: '25px', overflowX: 'auto', paddingBottom: '5px' }}>
              <button onClick={() => setSitioSeleccionado("Todos")} style={{ padding: '10px 20px', borderRadius: '8px', border: 'none', cursor: 'pointer', backgroundColor: sitioSeleccionado === "Todos" ? '#3154f4' : '#e2e8f0', color: sitioSeleccionado === "Todos" ? 'white' : '#475569', fontWeight: 'bold', fontSize: '14px' }}>Todas las Plantas</button>
              {sitios.map(sitio => (
                <button key={sitio.id} onClick={() => setSitioSeleccionado(sitio.nombre)} style={{ padding: '10px 20px', borderRadius: '8px', border: 'none', cursor: 'pointer', backgroundColor: sitioSeleccionado === sitio.nombre ? '#3154f4' : '#e2e8f0', color: sitioSeleccionado === sitio.nombre ? 'white' : '#475569', fontWeight: 'bold', fontSize: '14px' }}>{sitio.nombre}</button>
              ))}
            </div>

            <div style={{ overflowX: 'auto', backgroundColor: '#fff', borderRadius: '12px', border: '1px solid #e2e8f0', boxShadow: '0 4px 15px rgba(0,0,0,0.03)' }}>
              {cargando ? (
                <p style={{ padding: '30px', textAlign: 'center', color: '#64748b' }}>Sincronizando con base de datos...</p>
              ) : (
                <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', fontSize: '15px' }}>
                  <thead>
                    <tr style={{ backgroundColor: '#f8fafc', borderBottom: '2px solid #e2e8f0' }}>
                      <th style={{ padding: '18px', color: '#475569', fontWeight: '600' }}>Equipo</th>
                      <th style={{ padding: '18px', color: '#475569', fontWeight: '600' }}>Ubicación</th>
                      <th style={{ padding: '18px', color: '#475569', fontWeight: '600' }}>Usuario</th>
                      <th style={{ padding: '18px', color: '#475569', fontWeight: '600' }}>Estado</th>
                    </tr>
                  </thead>
                  <tbody>
                    {inventarioFiltrado.length > 0 ? (
                      inventarioFiltrado.map((equipo) => (
                        <tr key={equipo.id} style={{ borderBottom: '1px solid #e2e8f0' }}>
                          <td style={{ padding: '18px' }}><div style={{ fontWeight: 'bold', color: '#0f172a' }}>{equipo.host_name}</div><div style={{ color: '#64748b', fontSize: '13px', marginTop: '4px' }}>{equipo.serial}</div></td>
                          <td style={{ padding: '18px' }}><div style={{ fontWeight: 'bold', color: '#334155' }}>{equipo.sitioNombre}</div><div style={{ color: '#64748b', fontSize: '13px', marginTop: '4px' }}>{equipo.ubicacion}</div></td>
                          <td style={{ padding: '18px', fontWeight: '500', color: equipo.usuarioNombre === 'N/A' ? '#94a3b8' : '#0f172a' }}>{equipo.usuarioNombre}</td>
                          <td style={{ padding: '18px' }}><span style={{ backgroundColor: equipo.disposicion === 'EN USO' ? '#e6f7eb' : '#fee2e2', color: equipo.disposicion === 'EN USO' ? '#27b54d' : '#991b1b', padding: '6px 12px', borderRadius: '20px', fontSize: '12px', fontWeight: '700' }}>{equipo.disposicion}</span></td>
                        </tr>
                      ))
                    ) : (
                      <tr><td colSpan="4" style={{ padding: '50px', textAlign: 'center', color: '#64748b' }}>No se encontraron resultados en el inventario.</td></tr>
                    )}
                  </tbody>
                </table>
              )}
            </div>
          </>
        )}

        {/* === VISTA: INGRESAR EQUIPO === */}
        {vistaActiva === 'ingresar' && (
          <div style={{ maxWidth: '800px', backgroundColor: 'white', padding: '30px', borderRadius: '12px', boxShadow: '0 4px 15px rgba(0,0,0,0.05)', border: '1px solid #e2e8f0' }}>
            <h2 style={{ color: '#0f172a', marginBottom: '20px', fontSize: '24px' }}>Alta de Nuevo Equipo</h2>
            <p style={{ color: '#64748b', marginBottom: '30px' }}>Ingresa los datos del hardware recién adquirido o descubierto.</p>

            {mensajeForm && (
              <div style={{ padding: '15px', marginBottom: '20px', borderRadius: '8px', backgroundColor: mensajeForm.tipo === 'exito' ? '#e6f7eb' : '#fee2e2', color: mensajeForm.tipo === 'exito' ? '#27b54d' : '#991b1b', fontWeight: 'bold' }}>
                {mensajeForm.texto}
              </div>
            )}

            <form onSubmit={manejarSubmitIngreso} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
              
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Host Name (Ej: CLANGLP123)</label>
                <input required type="text" value={nuevoEquipo.host_name} onChange={e => setNuevoEquipo({...nuevoEquipo, host_name: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1' }} />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Número de Serie</label>
                <input required type="text" value={nuevoEquipo.serial} onChange={e => setNuevoEquipo({...nuevoEquipo, serial: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1' }} />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>ID del Modelo (1=Laptop, etc.)</label>
                <input required type="number" min="1" value={nuevoEquipo.modelo_id} onChange={e => setNuevoEquipo({...nuevoEquipo, modelo_id: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1' }} />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Año de Compra</label>
                <input required type="number" min="2010" max="2030" value={nuevoEquipo.anio_compra} onChange={e => setNuevoEquipo({...nuevoEquipo, anio_compra: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1' }} />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Año de Renovación</label>
                <input required type="number" min="2010" max="2035" value={nuevoEquipo.anio_renovacion} onChange={e => setNuevoEquipo({...nuevoEquipo, anio_renovacion: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1' }} />
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Disposición (Estado actual)</label>
                <select value={nuevoEquipo.disposicion} onChange={e => setNuevoEquipo({...nuevoEquipo, disposicion: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1', backgroundColor: 'white' }}>
                  <option value="BODEGA">En Bodega</option>
                  <option value="EN USO">En Uso</option>
                  <option value="DE BAJA">De Baja / Desecho</option>
                </select>
              </div>

              <div style={{ gridColumn: '1 / -1', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ fontWeight: 'bold', color: '#334155', fontSize: '14px' }}>Notas Técnicas (Opcional)</label>
                <textarea rows="3" value={nuevoEquipo.notas_tecnicas} onChange={e => setNuevoEquipo({...nuevoEquipo, notas_tecnicas: e.target.value})} style={{ padding: '10px', borderRadius: '6px', border: '1px solid #cbd5e1', resize: 'vertical' }} placeholder="Ej: Rayón en la tapa, cargador alternativo..."></textarea>
              </div>

              <div style={{ gridColumn: '1 / -1', marginTop: '10px' }}>
                <button type="submit" style={{ width: '100%', padding: '14px', backgroundColor: '#3154f4', color: 'white', border: 'none', borderRadius: '8px', fontWeight: 'bold', fontSize: '16px', cursor: 'pointer', transition: '0.2s', boxShadow: '0 4px 6px rgba(49, 84, 244, 0.3)' }}>
                  💾 Guardar Equipo en Base de Datos
                </button>
              </div>

            </form>
          </div>
        )}

      </div>
    </div>
  )
}

export default App