package controllers

import (
	"net/http"

	"api-gateway/src/services"
	"api-gateway/src/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	DockerService *services.DockerService
}

// NewAdminController crea una nueva instancia de AdminController.
func NewAdminController(dockerService *services.DockerService) *AdminController {
	return &AdminController{DockerService: dockerService}
}

// ListInstances devuelve la lista de contenedores.
func (ac *AdminController) ListInstances(c *gin.Context) {
	containers, err := ac.DockerService.ListContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, containers)
}

// CreatetInstance crea e inicia un nuevo contenedor.
// Se espera recibir un JSON con { "image": "nombre-imagen", "name": "nombre-contenedor", "port": "puerto" }
func (ac *AdminController) CreatetInstance(c *gin.Context) {
	var req struct {
		Image string `json:"image"`
		Name  string `json:"name"`
		Port  string `json:"port"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrada inválida"})
		return
	}

	// Crear el contenedor.
	id, err := ac.DockerService.CreateContainer(req.Image, req.Name, req.Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Iniciar el contenedor.
	if err := ac.DockerService.StartContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.ReloadNginxConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"container_id": id})
}

// StartInstance inicia un contenedor dado su ID (pasado como parámetro en la URL).
func (ac *AdminController) StartInstance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta el ID del contenedor"})
		return
	}

	if err := ac.DockerService.StartContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizamos la configuración de Nginx.
	if err := utils.ReloadNginxConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"started": id})
}

// StopInstance detiene un contenedor dado su ID (pasado como parámetro en la URL).
func (ac *AdminController) StopInstance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta el ID del contenedor"})
		return
	}

	if err := ac.DockerService.StopContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.ReloadNginxConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stopped": id})
}

// RemoveInstance elimina (remueve) un contenedor dado su ID (pasado como parámetro en la URL).
func (ac *AdminController) RemoveInstance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta el ID del contenedor"})
		return
	}

	// Remover el contenedor.
	if err := ac.DockerService.RemoveContainer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizamos la configuración de Nginx.
	if err := utils.ReloadNginxConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"removed": id})
}
