package main

import (
	"project-its/controllers"
	"project-its/initializers"
	"project-its/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://testing-production-9641.up.railway.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://localhost:8000"
		},
		MaxAge: 12 * time.Hour,
	}))

	//Websocket route
	r.GET("/ws", controllers.WebSocketHandler)

	// Route yang tidak memerlukan autentikasi
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Terapkan middleware autentikasi ke semua route selanjutnya
	r.Use(middleware.TokenAuthMiddleware())

	// Routes for User
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// r.GET("/updateAll", controllers.UpdateAllSheets)
	r.GET("/exportAll", controllers.ExportAllSheets)

	// Setup session store
	store = cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// logout must be after middleware
	r.POST("/logout", controllers.Logout)

	r.GET("/user", controllers.UserIndex)
	r.PUT("/user/:id", controllers.UserUpdate)
	r.DELETE("/user/:id", controllers.UserDelete)

	// Routes for MeetingList
	r.GET("/meetingSchedule", controllers.MeetingListIndex)
	r.POST("/meetingSchedule", controllers.MeetingListCreate)
	r.GET("/meetingSchedule/:id", controllers.MeetingListShow)
	r.PUT("/meetingSchedule/:id", controllers.MeetingListUpdate)
	r.DELETE("/meetingSchedule/:id", controllers.MeetingListDelete)
	r.GET("/exportMeetingList", controllers.CreateExcelMeetingList)
	r.POST("/uploadMeetingList", controllers.ImportExcelMeetingList)

	r.POST("/uploadFileMeetingList", controllers.UploadHandlerMeetingList)
	r.GET("/downloadMeetingList/:id/:filename", controllers.DownloadFileHandlerMeetingList)
	r.DELETE("/deleteMeetingList/:id/:filename", controllers.DeleteFileHandlerMeetingList)
	r.GET("/filesMeetingList/:id", controllers.GetFilesByIDMeetingList)

	// Routes for Meeting
	r.GET("/meetings", controllers.MeetingIndex)
	r.POST("/meetings", controllers.MeetingCreate)
	r.GET("/meetings/:id", controllers.MeetingShow)
	r.PUT("/meetings/:id", controllers.MeetingUpdate)
	r.DELETE("/meetings/:id", controllers.MeetingDelete)
	r.GET("/exportMeeting", controllers.CreateExcelMeeting)
	r.POST("/uploadMeeting", controllers.ImportExcelMeeting)

	r.POST("/uploadFileMeeting", controllers.UploadHandlerMeeting)
	r.GET("/downloadMeeting/:id/:filename", controllers.DownloadFileHandlerMeeting)
	r.DELETE("/deleteMeeting/:id/:filename", controllers.DeleteFileHandlerMeeting)
	r.GET("/filesMeeting/:id", controllers.GetFilesByIDMeeting)

	// Routes for Memo
	r.GET("/memo", controllers.MemoIndex)
	r.POST("/memo", controllers.MemoCreate)
	r.GET("/memo/:id", controllers.MemoShow)
	r.PUT("/memo/:id", controllers.MemoUpdate)
	r.DELETE("/memo/:id", controllers.MemoDelete)
	r.GET("/exportMemo", controllers.ExportMemoHandler)
	r.POST("/uploadMemo", controllers.ImportExcelMemo)

	r.GET("/beritaAcara", controllers.BeritaAcaraIndex)
	r.POST("/beritaAcara", controllers.BeritaAcaraCreate)
	r.GET("/beritaAcara/:id", controllers.BeritaAcaraShow)
	r.PUT("/beritaAcara/:id", controllers.BeritaAcaraUpdate)
	r.DELETE("/beritaAcara/:id", controllers.BeritaAcaraDelete)
	r.GET("/exportBeritaAcara", controllers.ExportBeritaAcaraHandler)
	r.POST("/uploadBeritaAcara", controllers.ImportExcelBeritaAcara)

	r.POST("/uploadFileBeritaAcara", controllers.UploadHandlerBeritaAcara)
	r.GET("/downloadBeritaAcara/:id/:filename", controllers.DownloadFileHandlerBeritaAcara)
	r.DELETE("/deleteBeritaAcara/:id/:filename", controllers.DeleteFileHandlerBeritaAcara)
	r.GET("/filesBeritaAcara/:id", controllers.GetFilesByIDBeritaAcara)

	r.GET("/surat", controllers.SuratIndex)
	r.POST("/surat", controllers.SuratCreate)
	r.GET("/surat/:id", controllers.SuratShow)
	r.PUT("/surat/:id", controllers.SuratUpdate)
	r.DELETE("/surat/:id", controllers.SuratDelete)
	r.GET("/exportSurat", controllers.ExportSuratHandler)
	r.POST("/uploadSurat", controllers.ImportExcelSurat)

	r.POST("/uploadFileSurat", controllers.UploadHandlerSurat)
	r.GET("/downloadSurat/:id/:filename", controllers.DownloadFileHandlerSurat)
	r.DELETE("/deleteSurat/:id/:filename", controllers.DeleteFileHandlerSurat)
	r.GET("/filesSurat/:id", controllers.GetFilesByIDSurat)

	r.GET("/sk", controllers.SkIndex)
	r.POST("/sk", controllers.SkCreate)
	r.GET("/sk/:id", controllers.SkShow)
	r.PUT("/sk/:id", controllers.SkUpdate)
	r.DELETE("/sk/:id", controllers.SkDelete)
	r.GET("/exportSk", controllers.ExportSkHandler)
	r.POST("/uploadSk", controllers.ImportExcelSk)

	r.POST("/uploadFileSk", controllers.UploadHandlerSk)
	r.GET("/downloadSk/:id/:filename", controllers.DownloadFileHandlerSk)
	r.DELETE("/deleteSk/:id/:filename", controllers.DeleteFileHandlerSk)
	r.GET("/filesSk/:id", controllers.GetFilesByIDSk)

	r.POST("/uploadFileMemo", controllers.UploadHandlerMemo)
	r.GET("/downloadMemo/:id/:filename", controllers.DownloadFileHandlerMemo)
	r.DELETE("/deleteMemo/:id/:filename", controllers.DeleteFileHandlerMemo)
	r.GET("/filesMemo/:id", controllers.GetFilesByIDMemo)

	// Project routes
	r.POST("/Project", controllers.ProjectCreate)
	r.PUT("/Project/:id", controllers.ProjectUpdate)
	r.GET("/Project", controllers.ProjectIndex)
	r.GET("/Project/:id", controllers.ProjectShow)
	r.DELETE("/Project/:id", controllers.ProjectDelete)
	r.GET("/exportProject", controllers.ExportProjectHandler)
	r.POST("/uploadProject", controllers.ImportExcelProject)

	r.POST("/uploadFileProject", controllers.UploadHandlerProject)
	r.GET("/downloadProject/:id/:filename", controllers.DownloadFileHandlerProject)
	r.DELETE("/deleteProject/:id/:filename", controllers.DeleteFileHandlerProject)
	r.GET("/filesProject/:id", controllers.GetFilesByIDProject)

	// Notif Calendar
	r.GET("/notifications", controllers.GetNotifications)
	r.DELETE("/notifications/:id", controllers.DeleteNotification)

	// Timeline Project routes
	r.GET("/timelineProject", controllers.GetEventsProject)
	r.POST("/timelineProject", controllers.CreateEventProject)
	r.DELETE("/timelineProject/:id", controllers.DeleteEventProject)
	r.GET("/resourceProject", controllers.GetResourcesProject)
	r.POST("/resourceProject", controllers.CreateResourceProject)
	r.DELETE("/resourceProject/:id", controllers.DeleteResourceProject)
	r.GET("/exportTimelineProject", controllers.ExportTimelineProjectToExcel)

	// Timeline Desktop routes
	r.GET("/timelineDesktop", controllers.GetEventsDesktop)
	r.POST("/timelineDesktop", controllers.CreateEventDesktop)
	r.DELETE("/timelineDesktop/:id", controllers.DeleteEventDesktop)
	r.GET("/resourceDesktop", controllers.GetResourcesDesktop)
	r.POST("/resourceDesktop", controllers.CreateResourceDesktop)
	r.DELETE("/resourceDesktop/:id", controllers.DeleteResourceDesktop)
	r.GET("/exportTimelineDesktop", controllers.ExportTimelineDesktopToExcel)

	// Booking Rapat routes
	r.GET("/booking-rapat", controllers.GetEventsBookingRapat)
	r.POST("/booking-rapat", controllers.CreateEventBookingRapat)
	r.DELETE("/booking-rapat/:id", controllers.DeleteEventBookingRapat)
	r.GET("/exportBookingRapat", controllers.ExportBookingRapatToExcel)

	// jadwal Rapat routes
	r.GET("/jadwal-rapat", controllers.GetEventsRapat)
	r.POST("/jadwal-rapat", controllers.CreateEventRapat)
	r.DELETE("/jadwal-rapat/:id", controllers.DeleteEventRapat)
	r.GET("/exportRapat", controllers.ExportJadwalRapatToExcel)

	// Jadwal Cuti routes
	r.GET("/jadwal-cuti", controllers.GetEventsCuti)
	r.POST("/jadwal-cuti", controllers.CreateEventCuti)
	r.DELETE("/jadwal-cuti/:id", controllers.DeleteEventCuti)
	r.GET("/exportCuti", controllers.ExportJadwalCutiToExcel)

	// Perdin routes
	r.POST("/Perdin", controllers.PerdinCreate)
	r.PUT("/Perdin/:id", controllers.PerdinUpdate)
	r.GET("/Perdin", controllers.PerdinIndex)
	r.DELETE("/Perdin/:id", controllers.PerdinDelete)
	r.GET("/Perdin/:id", controllers.PerdinShow)
	r.GET("/exportPerdin", controllers.CreateExcelPerdin)
	r.POST("/uploadPerdin", controllers.ImportExcelPerdin)

	r.POST("/uploadFilePerdin", controllers.UploadHandlerPerdin)
	r.GET("/downloadPerdin/:id/:filename", controllers.DownloadFileHandlerPerdin)
	r.DELETE("/deletePerdin/:id/:filename", controllers.DeleteFileHandlerPerdin)
	r.GET("/filesPerdin/:id", controllers.GetFilesByIDPerdin)

	// Surat  Masuk routes
	r.POST("/SuratMasuk", controllers.SuratMasukCreate)
	r.PUT("/SuratMasuk/:id", controllers.SuratMasukUpdate)
	r.GET("/SuratMasuk", controllers.SuratMasukIndex)
	r.DELETE("/SuratMasuk/:id", controllers.SuratMasukDelete)
	r.GET("/SuratMasuk/:id", controllers.SuratMasukShow)
	r.GET("/exportSuratMasuk", controllers.CreateExcelSuratMasuk)
	r.POST("/importSuratMasuk", controllers.ImportExcelSuratMasuk)

	r.POST("/uploadFileSuratMasuk", controllers.UploadHandlerSuratMasuk)
	r.GET("/downloadSuratMasuk/:id/:filename", controllers.DownloadFileHandlerSuratMasuk)
	r.DELETE("/deleteSuratMasuk/:id/:filename", controllers.DeleteFileHandlerSuratMasuk)
	r.GET("/filesSuratMasuk/:id", controllers.GetFilesByIDSuratMasuk)

	// Surat  Keluar routes
	r.POST("/SuratKeluar", controllers.SuratKeluarCreate)
	r.PUT("/SuratKeluar/:id", controllers.SuratKeluarUpdate)
	r.GET("/SuratKeluar", controllers.SuratKeluarIndex)
	r.DELETE("/SuratKeluar/:id", controllers.SuratKeluarDelete)
	r.GET("/SuratKeluar/:id", controllers.SuratKeluarShow)
	r.GET("/exportSuratKeluar", controllers.CreateExcelSuratKeluar)
	r.POST("/importSuratKeluar", controllers.ImportExcelSuratKeluar)

	r.POST("/uploadFileSuratKeluar", controllers.UploadHandlerSuratKeluar)
	r.GET("/downloadSuratKeluar/:id/:filename", controllers.DownloadFileHandlerSuratKeluar)
	r.DELETE("/deleteSuratKeluar/:id/:filename", controllers.DeleteFileHandlerSuratKeluar)
	r.GET("/filesSuratKeluar/:id", controllers.GetFilesByIDSuratKeluar)

	// Routes for Arsip
	r.GET("/Arsip", controllers.ArsipIndex)
	r.POST("/Arsip", controllers.ArsipCreate)
	r.PUT("/Arsip/:id", controllers.ArsipUpdate)
	r.GET("/Arsip/:id", controllers.ArsipShow)
	r.DELETE("/Arsip/:id", controllers.ArsipDelete)

	r.GET("/exportArsip", controllers.CreateExcelArsip)
	r.POST("/uploadArsip", controllers.ImportExcelArsip)

	r.POST("/upload", controllers.UploadHandlerArsip)
	r.GET("/files/:id", controllers.GetFilesByIDArsip)
	r.GET("/download/:id/:filename", controllers.DownloadFileHandlerArsip)
	r.DELETE("/delete/:id/:filename", controllers.DeleteFileHandlerArsip)

	// Request routes
	r.GET("/request", controllers.RequestIndex)
	r.GET("/AccRequest/:id", controllers.AccRequest)
	r.GET("/CancelRequest/:id", controllers.CancelRequest)
	// r.GET("/request/:id", controllers.RequestShow)
	// r.PUT("/request/:id", controllers.RequestUpdate)
	// r.DELETE("/request/:id", controllers.RequestDelete)

	r.Run()
}
