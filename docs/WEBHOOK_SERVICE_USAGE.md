# Webhook Service Usage Documentation

## Deskripsi
WebhookService berfungsi untuk mengirimkan notifikasi dokumen dari PSrE ke client callback URL dengan method POST.

## Fitur Utama
1. **HTTP POST Request**: Mengirim payload JSON ke client callback URL
2. **Request Logging**: Mencatat semua request dan response untuk audit
3. **Document Status Update**: Memperbarui status dokumen berdasarkan webhook status
4. **Error Handling**: Penanganan error yang komprehensif
5. **Timeout Management**: HTTP request dengan timeout 30 detik

## Struktur DTO

### WebhookDocumentNotificationRequest
```go
type WebhookDocumentNotificationRequest struct {
    DocumentID string `json:"documentId" binding:"required"`
    Status     string `json:"status" binding:"required"`
    SignedAT   string `json:"signedAt,omitempty"`
}
```

## Implementasi Service

### Method: SendDocumentNotification
```go
func (s *webhookService) SendDocumentNotification(req dto.WebhookDocumentNotificationRequest) error
```

#### Alur Kerja:
1. **Cari Document**: Mencari dokumen berdasarkan `DocumentID` dari database
2. **Validasi Callback URL**: Memastikan client memiliki callback URL yang valid
3. **Persiapan Payload**: Marshal request DTO ke format JSON
4. **HTTP Request**: Membuat dan mengirim POST request ke client callback URL
5. **Logging**: Mencatat request dan response untuk audit trail
6. **Update Status**: Memperbarui status dokumen jika webhook berhasil

#### Headers yang Dikirim:
- `Content-Type: application/json`
- `User-Agent: Dimensy-Bridge-Webhook/1.0`

#### Status Mapping:
- `SIGNED` → `DOCUMENT_STATUS_SIGNED`
- `ON_PROCESS` → `DOCUMENT_STATUS_ON_PROCESS` 
- `WAITING` → `DOCUMENT_STATUS_WAITING`

## Logging System

### ClientRequestLog Fields:
- `URL`: Client callback URL yang dipanggil
- `Type`: "WEBHOOK_NOTIFICATION"
- `ClientID`: ID client pemilik dokumen
- `Body`: JSON payload yang dikirim
- `Header`: "Content-Type: application/json"
- `Response`: Status code dan response body dari client

## Error Handling

### Kemungkinan Error:
1. **Document Not Found**: Dokumen dengan DocumentID tidak ditemukan
2. **Missing Callback URL**: Client tidak memiliki callback URL yang terset
3. **Marshal Error**: Gagal mengkonversi DTO ke JSON
4. **HTTP Request Error**: Gagal membuat atau mengirim HTTP request
5. **Client Response Error**: Client merespon dengan status code >= 400

### Error Response Format:
Semua error dikembalikan dengan pesan yang informatif dan di-log untuk debugging.

## Contoh Penggunaan

### 1. Di Handler (webhook_handler.go):
```go
func (h *WebhookHandler) HandlePSRENotification(c *gin.Context) {
    var req dto.WebhookDocumentNotificationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid request payload",
        })
        return
    }
    
    err := h.webhookSvc.SendDocumentNotification(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error", 
            "message": "failed to process webhook",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "webhook received",
    })
}
```

### 2. Request Payload Example:
```json
{
    "documentId": "doc-12345",
    "status": "SIGNED",
    "signedAt": "2024-11-10T10:30:00Z"
}
```

### 3. Client Callback Response Expected:
Client callback endpoint harus merespon dengan HTTP status 200-299 untuk dianggap berhasil.

## Dependencies

### Required Repositories:
- `ClientDocumentRepository`: Untuk mencari dokumen berdasarkan external ID
- `ClientRequestLogRepository`: Untuk logging request/response

### Database Transaction:
Service menggunakan GORM DB instance untuk update status dokumen.

## Configuration

### Environment Variables:
- Tidak ada environment variable khusus, menggunakan client callback URL dari database

### Timeout:
- HTTP Client timeout: 30 detik

## Best Practices

1. **Idempotency**: Pastikan client callback endpoint bersifat idempoten
2. **Response Time**: Client callback harus merespon dalam waktu < 30 detik
3. **Error Handling**: Client sebaiknya memberikan response yang informatif
4. **Logging**: Monitor log untuk debugging dan audit
5. **Retry Logic**: Pertimbangkan implementasi retry logic untuk kasus timeout

## Monitoring

### Metrics yang Bisa Dimonitor:
- Success rate webhook calls
- Response time dari client callbacks
- Error rate per client
- Volume webhook notifications per hari

### Log Analysis:
Gunakan `ClientRequestLog` table untuk analisis:
```sql
-- Webhook success rate per client
SELECT client_id, 
       COUNT(*) as total_calls,
       SUM(CASE WHEN response LIKE 'Status: 2%' THEN 1 ELSE 0 END) as success_calls
FROM client_request_logs 
WHERE type = 'WEBHOOK_NOTIFICATION'
GROUP BY client_id;
```