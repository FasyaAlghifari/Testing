package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"project-its/initializers"
	"project-its/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type SuratKeluarRequest struct {
	ID       uint    `gorm:"primaryKey"`
	NoSurat  *string `json:"no_surat"`
	Title    *string `json:"title"`
	From     *string `json:"from"`
	Pic      *string `json:"pic"`
	Tanggal  *string `json:"tanggal"`
	CreateBy string  `json:"create_by"`
	Version  uint    `gorm:"default:1"`
}

func UploadHandlerSuratKeluar(c *gin.Context) {
	id := c.PostForm("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File diperlukan"})
		return
	}

	// Konversi id dari string ke uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	baseDir := "C:/UploadedFile/suratkeluar"
	dir := filepath.Join(baseDir, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	filePath := filepath.Join(dir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// Menyimpan metadata file ke database
	newFile := models.File{
		UserID:      uint(userID), // Gunakan userID yang sudah dikonversi
		FilePath:    filePath,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
	}
	result := initializers.DB.Create(&newFile)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan metadata file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDSuratKeluar(c *gin.Context) {
	id := c.Param("id")

	var files []models.File
	result := initializers.DB.Where("user_id = ?", id).Find(&files)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data file"})
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.FileName)
	}

	c.JSON(http.StatusOK, gin.H{"files": fileNames})
}

func DeleteFileHandlerSuratKeluar(c *gin.Context) {
	encodedFilename := c.Param("filename")
	filename, err := url.QueryUnescape(encodedFilename)
	if err != nil {
		log.Printf("Error decoding filename: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	id := c.Param("id")
	log.Printf("Received ID: %s and Filename: %s", id, filename) // Tambahkan log ini

	baseDir := "C:/UploadedFile/suratkeluar"
	fullPath := filepath.Join(baseDir, id, filename)

	log.Printf("Attempting to delete file at path: %s", fullPath)

	// Hapus file dari sistem file
	err = os.Remove(fullPath)
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Hapus metadata file dari database
	result := initializers.DB.Where("file_path = ?", fullPath).Delete(&models.File{})
	if result.Error != nil {
		log.Printf("Error deleting file metadata: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func DownloadFileHandlerSuratKeluar(c *gin.Context) {
	id := c.Param("id")
	filename := c.Param("filename")
	baseDir := "C:/UploadedFile/suratkeluar"
	fullPath := filepath.Join(baseDir, id, filename)

	log.Printf("Full path for download: %s", fullPath)

	// Periksa keberadaan file di database
	var file models.File
	result := initializers.DB.Where("file_path = ?", fullPath).First(&file)
	if result.Error != nil {
		log.Printf("File not found in database: %v", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "File tidak ditemukan"})
		return
	}

	// Periksa keberadaan file di sistem file
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("File not found in system: %s", fullPath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File tidak ditemukan di sistem file"})
		return
	}

	log.Printf("File downloaded successfully: %s", fullPath)
	c.File(fullPath)
}

func SuratKeluarCreate(c *gin.Context) {
	// Get data off req body
	var requestBody SuratKeluarRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	requestBody.CreateBy = c.MustGet("username").(string)

	var tanggal *time.Time // Deklarasi variabel tanggal sebagai pointer ke time.Time
	if requestBody.Tanggal != nil && *requestBody.Tanggal != "" {
		// Parse the date string only if it's not nil and not empty
		parsedTanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		tanggal = &parsedTanggal
	}

	surat_keluar := models.SuratKeluar{
		NoSurat:  requestBody.NoSurat,
		Title:    requestBody.Title,
		From:     requestBody.From,
		Pic:      requestBody.Pic,
		Tanggal:  tanggal, // Gunakan tanggal yang telah diparsing, bisa jadi nil jika input kosong
		CreateBy: requestBody.CreateBy,
	}

	result := initializers.DB.Create(&surat_keluar)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarIndex(c *gin.Context) {

	// Get models from DB
	var surat_keluar []models.SuratKeluar
	initializers.DB.Find(&surat_keluar)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var surat_keluar models.SuratKeluar

	initializers.DB.First(&surat_keluar, id)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarUpdate(c *gin.Context) {

	var requestBody SuratKeluarRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

	var surat_keluar models.SuratKeluar
	initializers.DB.First(&surat_keluar, id)

	if err := initializers.DB.First(&surat_keluar, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "surat_keluar tidak ditemukan"})
		return
	}

	// Cek apakah versi yang diberikan cocok dengan versi di database
	if requestBody.Version != surat_keluar.Version {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Conflict: Data has been modified by another user",
		})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)
	surat_keluar.CreateBy = requestBody.CreateBy

	if requestBody.Tanggal != nil {
		tanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		surat_keluar.Tanggal = &tanggal
	}

	if requestBody.NoSurat != nil {
		surat_keluar.NoSurat = requestBody.NoSurat
	} else {
		surat_keluar.NoSurat = surat_keluar.NoSurat // gunakan nilai yang ada dari database
	}

	if requestBody.Title != nil {
		surat_keluar.Title = requestBody.Title
	} else {
		surat_keluar.Title = surat_keluar.Title // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != nil {
		surat_keluar.Pic = requestBody.Pic
	} else {
		surat_keluar.Pic = surat_keluar.Pic // gunakan nilai yang ada dari database
	}

	if requestBody.CreateBy != "" {
		surat_keluar.CreateBy = requestBody.CreateBy
	} else {
		surat_keluar.CreateBy = surat_keluar.CreateBy // gunakan nilai yang ada dari database
	}

	surat_keluar.Version ++

	initializers.DB.Model(&surat_keluar).Updates(surat_keluar)

	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the Surat Keluar
	var surat_keluar models.SuratKeluar

	if err := initializers.DB.First(&surat_keluar, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "surat_keluar not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&surat_keluar).Error; err != nil {
		c.JSON(404, gin.H{"error": "Surat Keluar Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"SuratKeluar": "Deleted",
	})
}

func CreateExcelSuratKeluar(c *gin.Context) {
	dir := "C:\\excel"
	baseFileName := "its_report_suratkeluar"
	filePath := filepath.Join(dir, baseFileName+".xlsx")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, append "_new" to the file name
		baseFileName += "_new"
	}

	fileName := baseFileName + ".xlsx"

	// File does not exist, create a new file
	f := excelize.NewFile()

	// Define sheet names
	sheetName := "SURAT KELUAR"

	// Create sheets and set headers for "SAG" only
	f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "No Surat")
	f.SetCellValue(sheetName, "B1", "Title Of Letter")
	f.SetCellValue(sheetName, "C1", "From")
	f.SetCellValue(sheetName, "D1", "Pic")
	f.SetCellValue(sheetName, "E1", "Date Issue")

	styleHeader, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4F81BD"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FFFFFF",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	f.SetCellStyle(sheetName, "A1", "E1", styleHeader)

	f.SetColWidth(sheetName, "A", "A", 27)
	f.SetColWidth(sheetName, "B", "B", 40)
	f.SetColWidth(sheetName, "C", "C", 20)
	f.SetColWidth(sheetName, "D", "D", 20)
	f.SetColWidth(sheetName, "E", "E", 20)
	f.SetRowHeight(sheetName, 1, 20)

	// Fetch initial data from the database
	var surat_keluars []models.SuratKeluar
	initializers.DB.Find(&surat_keluars)

	// Write initial data to the "SAG" sheet
	surat_keluarSheetName := "SURAT KELUAR"
	for i, surat_keluar := range surat_keluars {
		tanggalString := surat_keluar.Tanggal.Format("2 January 2006")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("A%d", rowNum), *surat_keluar.NoSurat)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("B%d", rowNum), *surat_keluar.Title)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("C%d", rowNum), *surat_keluar.From)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("D%d", rowNum), *surat_keluar.Pic)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("E%d", rowNum), tanggalString)

		styleData, err := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		f.SetCellStyle(surat_keluarSheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("E%d", rowNum), styleData)
	}

	// Delete the default "Sheet1" sheet
	if err := f.DeleteSheet("Sheet1"); err != nil {
		panic(err) // Handle error jika bukan error "sheet tidak ditemukan"
	}

	// Save the newly created file
	buf, err := f.WriteToBuffer()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}

	// Serve the file to the client
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Write(buf.Bytes())
}

func excelDateToTimeSuratKeluar(excelDate int) (time.Time, error) {
	baseDate := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	days := time.Duration(excelDate) * 24 * time.Hour
	return baseDate.Add(days), nil
}

func ImportExcelSuratKeluar(c *gin.Context) {
	// Mengambil file dari form upload
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	// Simpan file sementara jika perlu
	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name()) // Hapus file sementara setelah selesai

	// Salin file dari request ke file sementara
	if _, err := file.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "Error seeking file: %v", err)
		return
	}
	if _, err := io.Copy(tempFile, file); err != nil {
		c.String(http.StatusInternalServerError, "Error copying file: %v", err)
		return
	}

	// Buka file Excel dari file sementara
	tempFile.Seek(0, 0) // Reset pointer ke awal file
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	// Pilih sheet
	sheetName := "SURAT KELUAR"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}

	log.Println("Processing rows...") // Log untuk memulai proses baris

	// Definisikan semua format tanggal yang mungkin
	dateFormats := []string{
		"2 January 2006",
		"02-Jan-06",
		"2006-01-02",
		"02-01-2006",
		"01/02/2006",
		"2006.01.02",
		"02/01/2006",
		"Jan 2, 06",
		"Jan 2, 2006",
		"01/02/06",
		"02/01/06",
		"06/02/01",
		"06/01/02",
		"06-Jan-02",
		"02-Jan-06",
		"1-Jan-06",
		"06-Jan-02",
	}

	// Loop melalui baris dan simpan ke database
	for i, row := range rows {
		if i == 0 { // Lewati baris pertama yang merupakan header
			continue
		}
		if len(row) < 5 { // Pastikan ada cukup kolom
			log.Printf("Row %d skipped: less than 5 columns filled", i+1)
			continue
		}
		noSurat := row[0]
		title := row[1]
		from := row[2]
		pic := row[3]
		tanggalString := row[4]

		var tanggal time.Time
		var parseErr error

		// Coba konversi dari serial Excel jika tanggalString adalah angka
		if serial, err := strconv.Atoi(tanggalString); err == nil {
			tanggal, parseErr = excelDateToTimeSuratKeluar(serial)
		} else {
			// Coba parse menggunakan format tanggal yang sudah ada
			for _, format := range dateFormats {
				tanggal, parseErr = time.Parse(format, tanggalString)
				if parseErr == nil {
					break // Keluar dari loop jika parsing berhasil
				}
			}
		}

		if parseErr != nil {
			log.Printf("Format tanggal tidak valid di baris %d: %v", i+1, parseErr)
			continue // Lewati baris ini jika format tanggal tidak valid
		}

		// Buat instance baru dari models.SuratKeluar untuk setiap iterasi loop
		surat_keluar := models.SuratKeluar{
			NoSurat:  &noSurat,
			Title:    &title,
			From:     &from,
			Pic:      &pic,
			Tanggal:  &tanggal,
			CreateBy: c.MustGet("username").(string),
		}

		// Simpan ke database
		if err := initializers.DB.Create(&surat_keluar).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diimpor."})
}
