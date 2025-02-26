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

// GetLogs devuelve los logs de los contenedores.
// Si se envía el query parameter "service", se filtrarán los logs de ese servicio; de lo contrario, se devuelven todos los logs del proyecto "backend".
// Se puede enviar "since" (en segundos) para filtrar los logs de los últimos N segundos, y opcionalmente "until".
func (ac *AdminController) GetLogs(c *gin.Context) {
	// Ejemplo: /admin/logs?service=nginx&since=3600&until=
	service := c.Query("service")
	since := c.Query("since")
	until := c.Query("until")

	logs, err := ac.DockerService.GetLogs(service, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (ac *AdminController) GetStats(c *gin.Context) {

	id := c.Query("id")
	if id != "" {
		stats, err := ac.DockerService.GetContainerStats(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{id: stats})
		return
	}

	service := c.Query("service")
	containers, err := ac.DockerService.ListContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statsMap := make(map[string]interface{})
	for _, container := range containers {
		if service != "" {
			if container.Labels["com.docker.compose.service"] != service {
				continue
			}
		} else {
			if container.Labels["com.docker.compose.project"] != "backend" {
				continue
			}
		}

		stats, err := ac.DockerService.GetContainerStats(container.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		statsMap[container.ID] = stats
	}

	c.JSON(http.StatusOK, statsMap)
}
