package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Cambia los datos según tu configuración de base de datos.
	dsn := "proyectoFinal_hunterhave:root@tcp(hax.h.filess.io:3307)/proyectoFinal_hunterhave"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verificar que la conexión funcione.
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	fmt.Println("Conexión exitosa a la base de datos.")

	// Crear una nueva instancia de Gin
	router := gin.Default()

	// Rutas para Insertar, Actualizar, Eliminar, Obtener por ID y Obtener todos los usuarios
	router.POST("/usuario", insertarUsuario)
	router.PUT("/usuario/:id", actualizarUsuario)
	router.DELETE("/usuario/:id", eliminarUsuario)
	router.GET("/usuario/:id", obtenerUsuarioPorID) // Nueva ruta para obtener un usuario por ID
	router.GET("/usuarios", obtenerUsuarios)        // Nueva ruta para obtener todos los usuarios

	// Correr el servidor en el puerto 8080
	router.Run(":8080")
}

type NombreUsuario struct {
	Nombre string `json:"nombre"`
}

// Insertar un nuevo usuario
func insertarUsuario(c *gin.Context) {
	var entrada NombreUsuario
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := db.Exec("INSERT INTO usuario (nombre) VALUES (?)", entrada.Nombre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"message": "Usuario insertado correctamente", "id": id})
}

// Actualizar un usuario por ID

func actualizarUsuario(c *gin.Context) {
	id := c.Param("id")
	var entrada NombreUsuario
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if entrada.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre no puede estar vacío"})
		return
	}

	result, err := db.Exec("UPDATE usuario SET nombre = ? WHERE id = ?", entrada.Nombre, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado correctamente"})
}

// Eliminar un usuario por ID
func eliminarUsuario(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM usuario WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}

// Obtener un usuario por ID
func obtenerUsuarioPorID(c *gin.Context) {
	id := c.Param("id")
	var nombre string

	err := db.QueryRow("SELECT nombre FROM usuario WHERE id = ?", id).Scan(&nombre)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"nombre": nombre,
	})
}

// Obtener todos los usuarios
func obtenerUsuarios(c *gin.Context) {
	rows, err := db.Query("SELECT id, nombre FROM usuario")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []map[string]interface{}

	for rows.Next() {
		var id int
		var nombre string
		err := rows.Scan(&id, &nombre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usuario := map[string]interface{}{
			"id":     id,
			"nombre": nombre,
		}
		usuarios = append(usuarios, usuario)
	}

	c.JSON(http.StatusOK, usuarios)
}
