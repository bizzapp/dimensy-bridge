# Database Migration Scripts

## Create Client IP Whitelist Table

```sql
CREATE TABLE IF NOT EXISTS client_ip_whitelists (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_client_ip_whitelists_client_id 
        FOREIGN KEY (client_id) REFERENCES clients(id) 
        ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_client_ip_whitelists_client_id 
    ON client_ip_whitelists(client_id);

CREATE INDEX IF NOT EXISTS idx_client_ip_whitelists_ip_address 
    ON client_ip_whitelists(ip_address);

CREATE INDEX IF NOT EXISTS idx_client_ip_whitelists_deleted_at 
    ON client_ip_whitelists(deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_client_ip_whitelists_client_ip 
    ON client_ip_whitelists(client_id, ip_address) 
    WHERE deleted_at IS NULL;
```

## Using GORM Auto Migration

Add to your migration setup:

```go
// internal/config/db.go or migration file

import "dimensy-bridge/internal/model"

func MigrateDatabase(db *gorm.DB) error {
    return db.AutoMigrate(
        // ... existing models
        &model.ClientIPWhitelist{},
    )
}
```

Then call in main:

```go
func main() {
    db := config.InitDB()
    config.MigrateDatabase(db)
    // ... rest of code
}
```

## Insert Sample Data

```sql
-- Sample whitelist entries
INSERT INTO client_ip_whitelists (client_id, ip_address, description, is_active)
VALUES 
    (1, '192.168.1.100', 'Office Server', true),
    (1, '10.0.0.5', 'Development Machine', true),
    (2, '203.0.113.42', 'Production Server', true),
    (2, '198.51.100.0/24', 'VPN Network', false);
```

