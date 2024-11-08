package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Sesuaikan dengan aturan CORS
}

var clients = make(map[*websocket.Conn]bool) // Mengelola koneksi WebSocket
var broadcast = make(chan models.Notification) // Channel untuk broadcast notifikasi

func init() {
	// Goroutine untuk menerima notifikasi dari channel dan mengirim ke semua klien
	go func() {
		for notification := range broadcast {
			for client := range clients {
				err := client.WriteJSON(notification)
				if err != nil {
					log.Printf("WebSocket error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()
}

func WebSocketHandler(c *gin.Context) {
	// Upgrade koneksi HTTP ke WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
}

func SetNotification(title string, startTime time.Time, category string) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return
	}

	startTime, err = time.ParseInLocation(time.RFC3339, startTime.Format(time.RFC3339), loc)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return
	}
	log.Printf("Parsed start time in WIB: %v", startTime)

	notificationTime24 := startTime.Add(-24 * time.Hour)
	notificationTime1 := startTime.Add(-1 * time.Hour)

	notification := models.Notification{
		Title:    title,
		Start:    startTime,
		Category: category,
	}
	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
		return
	}

	go func() {
		time.Sleep(time.Until(notificationTime24))
		broadcast <- models.Notification{
			Title:    title,
			Start:    notificationTime24,
			Category: category + " - 24 jam sebelum event",
		}
		log.Printf("Notifikasi 24 jam dikirim untuk event %s pada %s", title, notificationTime24)
	}()

	go func() {
		time.Sleep(time.Until(notificationTime1))
		broadcast <- models.Notification{
			Title:    title,
			Start:    notificationTime1,
			Category: category + " - 1 jam sebelum event",
		}
		log.Printf("Notifikasi 1 jam dikirim untuk event %s pada %s", title, notificationTime1)
	}()
}

func GetNotifications(c *gin.Context) {
	var notifications []models.Notification
	if err := initializers.DB.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":"Notifikasi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi ditemukan"})
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID yang diterima untuk dihapus: %s", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID harus disertakan"})
		return
	}

	if err := initializers.DB.Where("id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil dihapus"})
}
