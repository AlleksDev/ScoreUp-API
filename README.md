# ScoreUp API

API REST + WebSocket en tiempo real para una plataforma de retos académicos gamificada. Construida con **Go**, **Gin**, **MySQL** y **Google Wire** siguiendo arquitectura hexagonal.

---

## Tabla de contenidos

- [Configuración](#configuración)
- [Autenticación](#autenticación)
- [Endpoints](#endpoints)
  - [Usuarios (público)](#-usuarios-público)
  - [Retos (protegido)](#-retos-protegido)
  - [Logros (protegido)](#-logros-protegido)
  - [Usuario-Retos (protegido)](#-usuario-retos-protegido)
  - [Usuario-Logros (protegido)](#-usuario-logros-protegido)
- [WebSocket](#-websocket)
- [Resumen rápido](#resumen-rápido-de-endpoints)

---

## Configuración

```bash
# Clonar el repositorio
git clone https://github.com/AlleksDev/ScoreUp-API.git
cd ScoreUp-API

# Crear archivo .env
echo 'DB_DSN=usuario:contraseña@tcp(host:3306)/retos_academicos?parseTime=true' > .env
echo 'JWT_SECRET=tu_secreto_aqui' >> .env
echo 'PORT=8080' >> .env

# Ejecutar
go run .
```

**Variables de entorno:**

| Variable     | Descripción                            | Ejemplo                                                       |
|-------------|----------------------------------------|---------------------------------------------------------------|
| `DB_DSN`    | DSN de conexión a MySQL                | `root:pass@tcp(localhost:3306)/retos_academicos?parseTime=true` |
| `JWT_SECRET`| Secreto para firmar tokens JWT (HS256) | `mi_clave_secreta_2026`                                        |
| `PORT`      | Puerto del servidor (default: 8080)    | `8080`                                                         |

---

## Autenticación

Los endpoints protegidos requieren un token JWT en el header:

```
Authorization: Bearer <TOKEN>
```

El token se obtiene mediante el endpoint de login y contiene los siguientes claims:

| Claim     | Tipo   | Descripción                  |
|-----------|--------|------------------------------|
| `user_id` | int    | ID del usuario               |
| `email`   | string | Email del usuario            |
| `name`    | string | Nombre del usuario           |
| `exp`     | int    | Expiración (Unix timestamp). Duración: 24 horas |

**Errores de autenticación comunes:**

| Código | Respuesta                                                        |
|--------|------------------------------------------------------------------|
| `401`  | `{"error": "Authorization header required"}`                     |
| `401`  | `{"error": "Formato de token inválido. Use: Bearer <token>"}`    |
| `401`  | `{"error": "Token inválido o expirado"}`                         |

---

## Endpoints

### 👤 Usuarios (público)

Estas rutas **no** requieren autenticación.

---

#### `POST /api/users/register`

Registra un nuevo usuario.

**Request body:**

```json
{
  "nombre": "Juan Pérez",
  "email": "juan@correo.com",
  "password": "miPassword123",
  "phone": "3312345678"
}
```

| Campo      | Tipo   | Requerido | Validación       |
|-----------|--------|-----------|------------------|
| `nombre`  | string | Sí        | mínimo 3 caracteres |
| `email`   | string | Sí        | formato email válido |
| `password`| string | Sí        | mínimo 6 caracteres |
| `phone`   | string | Sí        | —                |

**Respuestas:**

| Código | Respuesta                                                              |
|--------|------------------------------------------------------------------------|
| `201`  | `{"message": "Usuario creado exitosamente", "id": 1, "email": "juan@correo.com"}` |
| `400`  | `{"error": "Datos inválidos: ..."}`                                     |
| `409`  | `{"error": "...ya existe..."}`                                          |
| `500`  | `{"error": "..."}`                                                      |

---

#### `POST /api/users/login`

Inicia sesión y devuelve un token JWT.

**Request body:**

```json
{
  "email": "juan@correo.com",
  "password": "miPassword123"
}
```

| Campo      | Tipo   | Requerido |
|-----------|--------|-----------|
| `email`   | string | Sí        |
| `password`| string | Sí        |

**Respuestas:**

| Código | Respuesta                                                  |
|--------|------------------------------------------------------------|
| `200`  | `{"message": "Login exitoso", "token": "eyJhbGciOi..."}`  |
| `400`  | `{"error": "Datos inválidos: ..."}`                         |
| `401`  | `{"error": "Email o contraseña incorrectos"}`               |
| `500`  | `{"error": "..."}`                                          |

---

#### `GET /api/users/rank`

Obtiene el ranking de usuarios ordenado por puntuación.

**Request body:** ninguno

**Respuesta exitosa (`200`):**

```json
{
  "ranking": [
    {
      "ID": 3,
      "Name": "María López",
      "TotalScore": 150
    },
    {
      "ID": 1,
      "Name": "Juan Pérez",
      "TotalScore": 80
    }
  ]
}
```

| Código | Respuesta                                |
|--------|------------------------------------------|
| `200`  | `{"ranking": [...]}`                     |
| `500`  | `{"error": "Error obteniendo ranking"}`  |

---

### 🎯 Retos (protegido)

Requieren header `Authorization: Bearer <TOKEN>`.

---

#### `POST /api/retos`

Crea un nuevo reto. El `user_id` del creador se extrae del token JWT.

**Request body:**

```json
{
  "subject": "Matemáticas Discretas",
  "description": "Resolver 10 ejercicios de grafos",
  "goal": 10,
  "points_awarded": 25,
  "deadline": "2026-06-15"
}
```

| Campo            | Tipo    | Requerido | Validación / Default             |
|-----------------|---------|-----------|----------------------------------|
| `subject`       | string  | Sí        | —                                |
| `description`   | string  | Sí        | —                                |
| `goal`          | int     | Sí        | mínimo 1                         |
| `points_awarded`| int     | No        | default: `20`                    |
| `deadline`      | string  | No        | formato `YYYY-MM-DD`             |

**Respuestas:**

| Código | Respuesta                                                       |
|--------|----------------------------------------------------------------|
| `201`  | `{"message": "Reto creado exitosamente", "id": 5}`             |
| `400`  | `{"error": "Datos inválidos: ..."}` o `{"error": "Formato de fecha inválido, use YYYY-MM-DD"}` |
| `401`  | `{"error": "Usuario no autenticado"}`                           |
| `500`  | `{"error": "..."}`                                              |

> 📡 **WebSocket:** Al crearse el reto, se envía la lista completa de retos actualizada al canal `retos`.

---

#### `GET /api/retos`

Obtiene todos los retos.

**Respuesta exitosa (`200`):**

```json
{
  "retos": [
    {
      "ID": 1,
      "UserID": 3,
      "Subject": "Matemáticas Discretas",
      "Description": "Resolver 10 ejercicios de grafos",
      "Goal": 10,
      "PointsAwarded": 25,
      "Deadline": "2026-06-15T00:00:00Z",
      "CreatedAt": "2026-02-27T10:30:00Z"
    }
  ]
}
```

---

#### `GET /api/retos/:id`

Obtiene un reto por su ID.

**Respuestas:**

| Código | Respuesta                        |
|--------|----------------------------------|
| `200`  | `{"reto": {...}}`                |
| `400`  | `{"error": "ID inválido"}`       |
| `404`  | `{"error": "..."}`               |

---

#### `GET /api/retos/mine`

Obtiene los retos creados por el usuario autenticado.

**Respuestas:**

| Código | Respuesta                                |
|--------|------------------------------------------|
| `200`  | `{"retos": [...]}`                       |
| `401`  | `{"error": "Usuario no autenticado"}`    |
| `500`  | `{"error": "..."}`                       |

---

#### `PUT /api/retos/:id`

Actualiza un reto existente.

**Request body:** (misma estructura que POST)

```json
{
  "subject": "Matemáticas Discretas (actualizado)",
  "description": "Resolver 15 ejercicios de grafos",
  "goal": 15,
  "points_awarded": 30,
  "deadline": "2026-07-01"
}
```

**Respuestas:**

| Código | Respuesta                                          |
|--------|---------------------------------------------------|
| `200`  | `{"message": "Reto actualizado exitosamente"}`     |
| `400`  | `{"error": "ID inválido"}` / `{"error": "Datos inválidos: ..."}` / `{"error": "Formato de fecha inválido, use YYYY-MM-DD"}` |
| `500`  | `{"error": "..."}`                                  |

> 📡 **WebSocket:** Se envía la lista actualizada de retos al canal `retos`.

---

#### `DELETE /api/retos/:id`

Elimina un reto.

**Respuestas:**

| Código | Respuesta                                          |
|--------|---------------------------------------------------|
| `200`  | `{"message": "Reto eliminado exitosamente"}`       |
| `400`  | `{"error": "ID inválido"}`                          |
| `500`  | `{"error": "..."}`                                  |

> 📡 **WebSocket:** Se envía la lista actualizada de retos al canal `retos`.

---

### 🏆 Logros (protegido)

Requieren header `Authorization: Bearer <TOKEN>`.

---

#### `POST /api/logros`

Crea un nuevo logro.

**Request body:**

```json
{
  "name": "Primer reto completado",
  "description": "Completar tu primer reto académico",
  "required_points": 0,
  "required_retos": 1
}
```

| Campo             | Tipo   | Requerido | Default |
|------------------|--------|-----------|---------|
| `name`           | string | Sí        | —       |
| `description`    | string | Sí        | —       |
| `required_points`| int    | No        | `0`     |
| `required_retos` | int    | No        | `0`     |

**Respuestas:**

| Código | Respuesta                                                  |
|--------|-----------------------------------------------------------|
| `201`  | `{"message": "Logro creado exitosamente", "id": 2}`       |
| `400`  | `{"error": "Datos inválidos: ..."}`                        |
| `500`  | `{"error": "..."}`                                         |

---

#### `GET /api/logros`

Obtiene todos los logros.

**Respuesta exitosa (`200`):**

```json
{
  "logros": [
    {
      "ID": 1,
      "Name": "Primer reto completado",
      "Description": "Completar tu primer reto académico",
      "RequiredPoints": 0,
      "RequiredRetos": 1
    }
  ]
}
```

---

#### `GET /api/logros/:id`

Obtiene un logro por su ID.

**Respuestas:**

| Código | Respuesta                        |
|--------|----------------------------------|
| `200`  | `{"logro": {...}}`               |
| `400`  | `{"error": "ID inválido"}`       |
| `404`  | `{"error": "..."}`               |

---

#### `PUT /api/logros/:id`

Actualiza un logro existente.

**Request body:** (misma estructura que POST)

**Respuestas:**

| Código | Respuesta                                           |
|--------|------------------------------------------------------|
| `200`  | `{"message": "Logro actualizado exitosamente"}`      |
| `400`  | `{"error": "ID inválido"}` / `{"error": "Datos inválidos: ..."}` |
| `500`  | `{"error": "..."}`                                    |

---

#### `DELETE /api/logros/:id`

Elimina un logro.

**Respuestas:**

| Código | Respuesta                                          |
|--------|---------------------------------------------------|
| `200`  | `{"message": "Logro eliminado exitosamente"}`      |
| `400`  | `{"error": "ID inválido"}`                          |
| `500`  | `{"error": "..."}`                                  |

---

### 🔗 Usuario-Retos (protegido)

Gestiona la relación M:N entre usuarios y retos (participación, progreso).
Requieren header `Authorization: Bearer <TOKEN>`.

---

#### `POST /api/usuario-retos`

Unirse a un reto enviando el ID en el body.

**Request body:**

```json
{
  "reto_id": 5
}
```

**Respuestas:**

| Código | Respuesta                                       |
|--------|------------------------------------------------|
| `201`  | `{"message": "Unido al reto exitosamente"}`     |
| `400`  | `{"error": "Datos inválidos: ..."}`              |
| `401`  | `{"error": "Usuario no autenticado"}`            |
| `500`  | `{"error": "..."}`                               |

---

#### `POST /api/usuario-retos/:retoId/join`

Unirse a un reto usando el ID como parámetro de ruta.

**Respuestas:** (mismas que el anterior)

---

#### `PUT /api/usuario-retos/:retoId/progress`

Actualiza el progreso del usuario en un reto. Si el progreso alcanza la meta, se marca como completado, se suman puntos al score del usuario y se evalúan logros.

**Request body:**

```json
{
  "progress": 7
}
```

| Campo      | Tipo | Requerido | Validación |
|-----------|------|-----------|------------|
| `progress`| int  | Sí        | mínimo 0   |

**Respuesta exitosa (`200`):**

```json
{
  "message": "Progreso actualizado exitosamente",
  "completed": true,
  "logros_awarded": [
    {
      "ID": 1,
      "Name": "Primer reto completado"
    }
  ]
}
```

> `completed` indica si el reto se completó con esta actualización.  
> `logros_awarded` solo aparece si se desbloquearon logros nuevos.

| Código | Respuesta                                        |
|--------|--------------------------------------------------|
| `200`  | `{"message": "...", "completed": bool, ...}`     |
| `400`  | `{"error": "ID de reto inválido"}` / `{"error": "Datos inválidos: ..."}` |
| `401`  | `{"error": "Usuario no autenticado"}`            |
| `500`  | `{"error": "..."}`                               |

> 📡 **WebSocket:** Al actualizar el progreso se envía el ranking actualizado al canal `rank`.

---

#### `GET /api/usuario-retos`

Obtiene los retos en los que participa el usuario autenticado.

**Respuesta exitosa (`200`):**

```json
{
  "usuario_retos": [
    {
      "UserID": 1,
      "RetoID": 5,
      "Progress": 7,
      "Status": "completado",
      "JoinedAt": "2026-02-20T08:00:00Z"
    }
  ]
}
```

---

#### `GET /api/usuario-retos/:retoId/participants`

Obtiene todos los participantes de un reto.

**Respuestas:**

| Código | Respuesta                            |
|--------|--------------------------------------|
| `200`  | `{"usuario_retos": [...]}`           |
| `400`  | `{"error": "ID de reto inválido"}`   |
| `500`  | `{"error": "..."}`                   |

---

#### `DELETE /api/usuario-retos/:retoId`

Abandonar un reto.

**Respuestas:**

| Código | Respuesta                                             |
|--------|------------------------------------------------------|
| `200`  | `{"message": "Abandonaste el reto exitosamente"}`     |
| `400`  | `{"error": "ID de reto inválido"}`                     |
| `401`  | `{"error": "Usuario no autenticado"}`                  |
| `500`  | `{"error": "..."}`                                     |

---

### 🎖 Usuario-Logros (protegido)

Gestiona los logros obtenidos por cada usuario.
Requieren header `Authorization: Bearer <TOKEN>`.

---

#### `POST /api/usuario-logros/evaluate`

Evalúa y asigna automáticamente los logros que el usuario haya desbloqueado según sus puntos y retos completados.

**Request body:** ninguno

**Respuesta exitosa (`200`):**

```json
{
  "message": "Evaluación de logros completada",
  "logros_awarded": [
    {
      "ID": 2,
      "Name": "Maestro de retos"
    }
  ]
}
```

---

#### `GET /api/usuario-logros`

Obtiene los logros del usuario autenticado.

**Respuesta exitosa (`200`):**

```json
{
  "usuario_logros": [
    {
      "UserID": 1,
      "LogroID": 2,
      "ObtainedAt": "2026-02-27T12:00:00Z"
    }
  ]
}
```

---

#### `DELETE /api/usuario-logros/:logroId`

Remueve un logro del usuario.

**Respuestas:**

| Código | Respuesta                                              |
|--------|---------------------------------------------------------|
| `200`  | `{"message": "Logro removido del usuario exitosamente"}`|
| `400`  | `{"error": "ID de logro inválido"}`                      |
| `401`  | `{"error": "Usuario no autenticado"}`                    |
| `500`  | `{"error": "..."}`                                       |

---

## 📡 WebSocket

Conexión WebSocket para recibir actualizaciones en tiempo real.

### Conectar

```
ws://<HOST>:<PORT>/ws?role=alumno&user_id=1&channel=retos
```

| Parámetro  | Tipo   | Requerido | Descripción                                   |
|-----------|--------|-----------|-----------------------------------------------|
| `role`    | string | Sí        | Rol del cliente (ej: `alumno`, `profesor`)     |
| `user_id` | string | Sí        | ID del usuario                                 |
| `channel` | string | No        | Canal de suscripción: `retos` o `rank`         |

**Errores:**

| Código | Respuesta                                              |
|--------|---------------------------------------------------------|
| `400`  | `{"error": "role y user_id son requeridos"}`            |
| `503`  | `{"error": "Límite de conexiones WebSocket alcanzado"}` |

### Canales disponibles

| Canal    | Se dispara cuando                             | Datos recibidos                         |
|---------|-----------------------------------------------|-----------------------------------------|
| `retos` | Se crea, actualiza o elimina un reto          | `{"retos": [<lista completa de retos>]}`|
| `rank`  | Se actualiza el progreso de un usuario-reto   | `{"ranking": [<ranking actualizado>]}`  |

### Ejemplo de conexión (JavaScript)

```javascript
const ws = new WebSocket('ws://184.72.233.162:8080/ws?role=alumno&user_id=1&channel=retos');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Actualización recibida:', data);
};
```

### Configuración técnica

| Parámetro           | Valor  |
|--------------------|--------|
| Max conexiones     | 2048   |
| Ping interval      | 54s    |
| Pong timeout       | 60s    |
| Write timeout      | 10s    |
| Max tamaño mensaje | 4 KB   |

---

## Resumen rápido de endpoints

| Método   | Ruta                                         | Auth | Descripción                      | WS  |
|----------|----------------------------------------------|------|----------------------------------|-----|
| `POST`   | `/api/users/register`                        | No   | Registrar usuario                | —   |
| `POST`   | `/api/users/login`                           | No   | Login → JWT                      | —   |
| `GET`    | `/api/users/rank`                            | No   | Ranking de usuarios              | —   |
| `POST`   | `/api/retos`                                 | JWT  | Crear reto                       | 📡 `retos` |
| `GET`    | `/api/retos`                                 | JWT  | Listar todos los retos           | —   |
| `GET`    | `/api/retos/mine`                            | JWT  | Mis retos creados                | —   |
| `GET`    | `/api/retos/:id`                             | JWT  | Obtener reto por ID              | —   |
| `PUT`    | `/api/retos/:id`                             | JWT  | Actualizar reto                  | 📡 `retos` |
| `DELETE` | `/api/retos/:id`                             | JWT  | Eliminar reto                    | 📡 `retos` |
| `POST`   | `/api/logros`                                | JWT  | Crear logro                      | —   |
| `GET`    | `/api/logros`                                | JWT  | Listar logros                    | —   |
| `GET`    | `/api/logros/:id`                            | JWT  | Obtener logro por ID             | —   |
| `PUT`    | `/api/logros/:id`                            | JWT  | Actualizar logro                 | —   |
| `DELETE` | `/api/logros/:id`                            | JWT  | Eliminar logro                   | —   |
| `POST`   | `/api/usuario-retos`                         | JWT  | Unirse a reto (body)             | —   |
| `POST`   | `/api/usuario-retos/:retoId/join`            | JWT  | Unirse a reto (param)            | —   |
| `PUT`    | `/api/usuario-retos/:retoId/progress`        | JWT  | Actualizar progreso              | 📡 `rank` |
| `GET`    | `/api/usuario-retos`                         | JWT  | Mis retos (participación)        | —   |
| `GET`    | `/api/usuario-retos/:retoId/participants`    | JWT  | Participantes de un reto         | —   |
| `DELETE` | `/api/usuario-retos/:retoId`                 | JWT  | Abandonar reto                   | —   |
| `POST`   | `/api/usuario-logros/evaluate`               | JWT  | Evaluar logros del usuario       | —   |
| `GET`    | `/api/usuario-logros`                        | JWT  | Mis logros obtenidos             | —   |
| `DELETE` | `/api/usuario-logros/:logroId`               | JWT  | Remover logro del usuario        | —   |
| `GET`    | `/ws`                                        | No   | Conexión WebSocket               | —   |