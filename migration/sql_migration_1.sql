-- Function and Trigger to update `update_at`
CREATE
OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
= NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Table for Asset Status
CREATE TABLE asset_status
(
    asset_status_id SERIAL PRIMARY KEY,
    status_name     VARCHAR(255) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by      VARCHAR(255),
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMP,
    deleted_by      VARCHAR(255)
);

-- Insert default assets statuses
INSERT INTO asset_status (status_name, description, created_by, updated_by)
VALUES ('Available', 'The assets is available for use', 'system', 'system'),
       ('In Use', 'The assets is currently being used', 'system', 'system'),
       ('Under Maintenance', 'The assets is undergoing maintenance', 'system', 'system'),
       ('Retired', 'The assets is no longer in use', 'system', 'system');

-- Table for Asset Category
CREATE TABLE asset_category
(
    asset_category_id SERIAL PRIMARY KEY,
    category_name     VARCHAR(255) NOT NULL,
    description       TEXT,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by        VARCHAR(255),
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by        VARCHAR(255),
    deleted_at        TIMESTAMP,
    deleted_by        VARCHAR(255)
);

INSERT INTO asset_category (category_name, description, created_by, updated_by)
VALUES ('Electronics', 'Electronic devices and gadgets', 'system', 'system'),
       ('Furniture', 'Office and home furniture', 'system', 'system'),
       ('Vehicles', 'Company vehicles', 'system', 'system'),
       ('Machinery', 'Industrial machinery and tools', 'system', 'system'),
       ('Software', 'Licensed software assets', 'system', 'system'),
       ('Hardware', 'Computer hardware components', 'system', 'system'),
       ('Stationery', 'Office stationery and supplies', 'system', 'system'),
       ('Appliances', 'Home and office appliances', 'system', 'system'),
       ('Real Estate', 'Buildings and land owned by the company', 'system', 'system'),
       ('Tools', 'Hand and power tools', 'system', 'system'),
       ('Medical Equipment', 'Healthcare devices and machines', 'system', 'system'),
       ('Security Equipment', 'Cameras, alarms, and other security devices', 'system', 'system'),
       ('Laboratory Equipment', 'Scientific and testing tools', 'system', 'system'),
       ('Books', 'Books and reference materials', 'system', 'system'),
       ('IT Equipment', 'Computers, servers, and network devices', 'system', 'system'),
       ('Art', 'Paintings and decorative items', 'system', 'system'),
       ('Audio Equipment', 'Sound systems and related items', 'system', 'system'),
       ('Video Equipment', 'Cameras and video recording devices', 'system', 'system'),
       ('Clothing', 'Uniforms and protective clothing', 'system', 'system'),
       ('Miscellaneous', 'Other uncategorized assets', 'system', 'system');

-- Table for Asset
CREATE TABLE asset
(
    asset_id             SERIAL PRIMARY KEY,
    user_client_id       VARCHAR(50)  NOT NULL,
    asset_code           VARCHAR(100)   DEFAULT NULL,
    name                 VARCHAR(100) NOT NULL,
    description          TEXT,
    barcode              VARCHAR(100)   DEFAULT NULL,
    category_id          INT          NOT NULL,
    status_id            INT          NOT NULL,
    purchase_date        DATE,
    expiry_date          DATE           DEFAULT NULL,
    warranty_expiry_date DATE           DEFAULT NULL, -- Warranty expiration date
    insurance_policy     JSONB          DEFAULT '{}'::JSONB,
    price                DECIMAL(40, 2) DEFAULT 0,
    stock                INT            DEFAULT 0,
    general              JSONB          DEFAULT '{}'::JSONB,
    is_wishlist          BOOLEAN        DEFAULT FALSE,
    created_at           TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    created_by           VARCHAR(255),
    updated_at           TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    updated_by           VARCHAR(255),
    deleted_at           TIMESTAMP,
    deleted_by           VARCHAR(255),
    FOREIGN KEY (category_id) REFERENCES asset_category (asset_category_id),
    FOREIGN KEY (status_id) REFERENCES asset_status (asset_status_id)
);

-- Table for Maintenance Type
CREATE TABLE asset_maintenance_type
(
    type_id     SERIAL PRIMARY KEY,           -- Unique ID for each type
    type_name   VARCHAR(100) UNIQUE NOT NULL, -- Name of maintenance type
    description TEXT,                         -- Description of what this maintenance involves
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP,
    deleted_by  VARCHAR(255)
);

INSERT INTO asset_maintenance_type (type_name, description, created_by, updated_by)
VALUES ('Battery Replacement', 'Replace worn-out battery with a new one', 'system', 'system'),
       ('Software Update', 'Update firmware and OS for better performance', 'system', 'system'),
       ('Cleaning Service', 'Perform deep cleaning to remove dust and debris', 'system', 'system'),
       ('Hardware Repair', 'Fix or replace damaged hardware components', 'system', 'system'),
       ('Annual Inspection', 'General check-up and inspection of the asset', 'system', 'system'),
       ('Firmware Upgrade', 'Upgrade device firmware to latest version', 'system', 'system'),
       ('Cooling System Check', 'Inspect and clean cooling fans and heat sinks', 'system', 'system'),
       ('Electrical Testing', 'Ensure safe electrical operation', 'system', 'system'),
       ('Oil & Lubrication', 'Apply lubrication to mechanical parts', 'system', 'system'),
       ('Parts Replacement', 'Replace any broken or worn-out components', 'system', 'system'),
       ('Sensor Calibration', 'Calibrate sensors for accurate readings', 'system', 'system'),
       ('Security Patch Update', 'Apply latest security updates and fixes', 'system', 'system'),
       ('Networking Maintenance', 'Check and repair network connections', 'system', 'system'),
       ('Cloud Backup Check', 'Ensure cloud backups are up to date', 'system', 'system'),
       ('General Diagnostics', 'Perform full diagnostics to detect issues', 'system', 'system'),
       ('Performance Optimization', 'Improve asset performance and efficiency', 'system', 'system'),
       ('Other', 'Any other maintenance not covered in predefined types', 'system', 'system');

-- Table for Asset Maintenance
CREATE TABLE asset_maintenance
(
    id                  SERIAL PRIMARY KEY,
    asset_id         INT NOT NULL,      -- Reference to asset
    type_id          INT NOT NULL,
    maintenance_date    DATE NOT NULL,
    maintenance_details TEXT,
    maintenance_cost DECIMAL(15, 2),    -- Cost of maintenance
    performed_by     VARCHAR(255),      -- Who performed the maintenance
    next_due_date    DATE DEFAULT NULL, -- Scheduled next maintenance
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by          VARCHAR(255),
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by          VARCHAR(255),
    deleted_at          TIMESTAMP,
    deleted_by          VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id)
);

CREATE TABLE asset_maintenance_record
(
    maintenance_record_id SERIAL PRIMARY KEY,
    asset_id              INT          NOT NULL, -- Reference to the asset
    maintenance_date      DATE         NOT NULL, -- Date of maintenance
    maintenance_type      VARCHAR(255) NOT NULL, -- Type of maintenance (e.g., Repair, Inspection)
    maintenance_details   TEXT,                  -- Details of the maintenance work
    maintenance_cost      DECIMAL(15, 2),        -- Cost of maintenance
    performed_by          VARCHAR(255),          -- Who performed the maintenance
    next_due_date         DATE,                  -- Date of the next maintenance (if applicable)
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by            VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by            VARCHAR(255),
    deleted_at            TIMESTAMP,
    deleted_by            VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id)
);

CREATE TABLE asset_tags
(
    tag_id      SERIAL PRIMARY KEY,
    tag_name    VARCHAR(255) NOT NULL UNIQUE, -- Name of the tag
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255)
);

CREATE TABLE asset_tag_map
(
    asset_id INT NOT NULL,
    tag_id   INT NOT NULL,
    PRIMARY KEY (asset_id, tag_id),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
    FOREIGN KEY (tag_id) REFERENCES asset_tags (tag_id)
);

CREATE TABLE asset_audit_log
(
    log_id       SERIAL PRIMARY KEY,
    table_name   VARCHAR(255) NOT NULL,
    action       VARCHAR(255) NOT NULL,
    old_data     TEXT,
    new_data     TEXT,
    performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    performed_by VARCHAR(255)
);

CREATE TABLE cron_jobs
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL UNIQUE,
    schedule         VARCHAR(255) NOT NULL, -- Cron expression
    is_active        BOOLEAN   DEFAULT TRUE,
    description      TEXT,
    last_executed_at TIMESTAMP,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
);

INSERT INTO cron_jobs (name, schedule, is_active, description, created_by)
VALUES ('asset_maintenance', '* * * * *', true, 'Check Maintenance Asset', 'system');


-- Triggers to update `update_at`
CREATE TRIGGER trigger_update_cron_jobs
    BEFORE UPDATE
    ON cron_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Triggers to update `update_at`
CREATE TRIGGER trigger_update_asset_status
    BEFORE UPDATE
    ON asset_status
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_asset_category
    BEFORE UPDATE
    ON asset_category
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_asset
    BEFORE UPDATE
    ON asset
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_asset_maintenance
    BEFORE UPDATE
    ON asset_maintenance
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();