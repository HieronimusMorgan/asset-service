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
VALUES ('Electronics', 'Electronic devices and gadgets', 'admin', 'admin'),
       ('Furniture', 'Office and home furniture', 'admin', 'admin'),
       ('Vehicles', 'Company vehicles', 'admin', 'admin'),
       ('Machinery', 'Industrial machinery and tools', 'admin', 'admin'),
       ('Software', 'Licensed software assets', 'admin', 'admin'),
       ('Hardware', 'Computer hardware components', 'admin', 'admin'),
       ('Stationery', 'Office stationery and supplies', 'admin', 'admin'),
       ('Appliances', 'Home and office appliances', 'admin', 'admin'),
       ('Real Estate', 'Buildings and land owned by the company', 'admin', 'admin'),
       ('Tools', 'Hand and power tools', 'admin', 'admin'),
       ('Medical Equipment', 'Healthcare devices and machines', 'admin', 'admin'),
       ('Security Equipment', 'Cameras, alarms, and other security devices', 'admin', 'admin'),
       ('Laboratory Equipment', 'Scientific and testing tools', 'admin', 'admin'),
       ('Books', 'Books and reference materials', 'admin', 'admin'),
       ('IT Equipment', 'Computers, servers, and network devices', 'admin', 'admin'),
       ('Art', 'Paintings and decorative items', 'admin', 'admin'),
       ('Audio Equipment', 'Sound systems and related items', 'admin', 'admin'),
       ('Video Equipment', 'Cameras and video recording devices', 'admin', 'admin'),
       ('Clothing', 'Uniforms and protective clothing', 'admin', 'admin'),
       ('Miscellaneous', 'Other uncategorized assets', 'admin', 'admin');

-- Table for Asset
CREATE TABLE asset
(
    asset_id       SERIAL PRIMARY KEY,
    user_client_id VARCHAR(255) NOT NULL,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    category_id    INT          NOT NULL,
    status_id      INT          NOT NULL,
    purchase_date  DATE,
    expiry_date    DATE,
    value          DECIMAL(40, 2),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255),
    FOREIGN KEY (category_id) REFERENCES asset_category (asset_category_id),
    FOREIGN KEY (status_id) REFERENCES asset_status (asset_status_id)
);

-- Table for Asset Maintenance
CREATE TABLE asset_maintenance
(
    id                  SERIAL PRIMARY KEY,
    asset_id            INT  NOT NULL,
    maintenance_date    DATE NOT NULL,
    maintenance_details TEXT,
    maintenance_cost    DECIMAL(15, 2),
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
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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