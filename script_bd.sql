-- =====================================================
-- BASE DE DATOS
-- =====================================================

DROP DATABASE IF EXISTS retos_academicos;
CREATE DATABASE retos_academicos;
USE retos_academicos;

-- =====================================================
-- TABLA USUARIOS (LOGIN)
-- =====================================================

CREATE TABLE usuarios (
    id_usuario INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    puntos_totales INT NOT NULL DEFAULT 0,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- TABLA RETOS (CATÁLOGO DE RETOS)
-- =====================================================

CREATE TABLE retos (
    id_reto INT AUTO_INCREMENT PRIMARY KEY,
    id_usuario INT NOT NULL,
    materia VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    meta INT NOT NULL,
    puntos_otorgados INT DEFAULT 20,
    fecha_limite DATE,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario)
        ON DELETE CASCADE
);

-- =====================================================
-- TABLA USUARIO_RETOS (RELACIÓN M:N)
-- =====================================================

CREATE TABLE usuario_retos (
    id_usuario INT NOT NULL,
    id_reto INT NOT NULL,
    progreso INT DEFAULT 0,
    estado ENUM('activo', 'completado') DEFAULT 'activo',
    fecha_union TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id_usuario, id_reto),
    
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario)
        ON DELETE CASCADE,
    FOREIGN KEY (id_reto) REFERENCES retos(id_reto)
        ON DELETE CASCADE
);

-- =====================================================
-- TABLA LOGROS (CATÁLOGO)
-- =====================================================

CREATE TABLE logros (
    id_logro INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL,
    puntos_requeridos INT DEFAULT 0,
    retos_requeridos INT DEFAULT 0
);

-- =====================================================
-- TABLA USUARIO_LOGROS (RELACIÓN)
-- =====================================================

CREATE TABLE usuario_logros (
    id_usuario INT,
    id_logro INT,
    fecha_obtenido TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id_usuario, id_logro),
    
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id_usuario)
        ON DELETE CASCADE,
    FOREIGN KEY (id_logro) REFERENCES logros(id_logro)
        ON DELETE CASCADE
);

-- =====================================================
-- VISTA DE RANKING
-- =====================================================

CREATE VIEW ranking AS
SELECT 
    id_usuario,
    nombre,
    puntos_totales
FROM usuarios
ORDER BY puntos_totales DESC;

-- =====================================================
-- TRIGGER: COMPLETAR RETO AUTOMÁTICAMENTE
-- =====================================================

DELIMITER $$

CREATE TRIGGER completar_usuario_reto
BEFORE UPDATE ON usuario_retos
FOR EACH ROW
BEGIN
    DECLARE v_meta INT;
    SELECT meta INTO v_meta FROM retos WHERE id_reto = NEW.id_reto;
    IF NEW.progreso >= v_meta AND OLD.estado = 'activo' THEN
        SET NEW.estado = 'completado';
    END IF;
END$$

DELIMITER ;

-- =====================================================
-- TRIGGER: SUMAR PUNTOS AL COMPLETAR RETO
-- =====================================================

DELIMITER $$

CREATE TRIGGER sumar_puntos_usuario_reto
AFTER UPDATE ON usuario_retos
FOR EACH ROW
BEGIN
    DECLARE v_puntos INT;
    IF NEW.estado = 'completado' AND OLD.estado = 'activo' THEN
        SELECT puntos_otorgados INTO v_puntos FROM retos WHERE id_reto = NEW.id_reto;
        UPDATE usuarios
        SET puntos_totales = puntos_totales + v_puntos
        WHERE id_usuario = NEW.id_usuario;
    END IF;
END$$

DELIMITER ;

-- =====================================================
-- DATOS INICIALES DE LOGROS
-- =====================================================

INSERT INTO logros (nombre, descripcion, puntos_requeridos, retos_requeridos)
VALUES 
('Primer Paso', 'Completa tu primer reto', 0, 1),
('Constante', 'Alcanza 100 puntos acumulados', 100, 0),
('Productivo', 'Completa 5 retos', 0, 5);

-- =====================================================
-- USUARIO DE PRUEBA
-- (password: 123456 encriptar en backend en producción)
-- =====================================================

INSERT INTO usuarios (nombre, email, password)
VALUES ('Anna', 'anna@email.com', '123456');

-- =====================================================
-- RETO DE PRUEBA
-- =====================================================

INSERT INTO retos (id_usuario, materia, descripcion, meta, puntos_otorgados, fecha_limite)
VALUES (1, 'Matemáticas', 'Resolver 20 ejercicios', 20, 20, '2026-03-30');

-- =====================================================
-- USUARIO SE UNE AL RETO DE PRUEBA
-- =====================================================

INSERT INTO usuario_retos (id_usuario, id_reto, progreso, estado)
VALUES (1, 1, 0, 'activo');