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
-- Insert default asset statuses with better categorization
INSERT INTO asset_status (status_name, description, created_by, updated_by)
VALUES
    -- ✅ General Asset Lifecycle Statuses
    ('Available', 'The asset is available for use or allocation.', 'system', 'system'),
    ('In Use', 'The asset is currently being used by an individual or department.', 'system', 'system'),
    ('Under Maintenance', 'The asset is undergoing scheduled maintenance or repair.', 'system', 'system'),
    ('Retired', 'The asset is no longer in operational use but not yet disposed.', 'system', 'system'),
    ('Lost', 'The asset is lost or missing and needs further tracking.', 'system', 'system'),
    ('Damaged', 'The asset is damaged and requires repair.', 'system', 'system'),
    ('Stolen', 'The asset has been reported stolen and needs recovery.', 'system', 'system'),
    ('Disposed', 'The asset has been officially disposed or recycled.', 'system', 'system'),
    ('Sold', 'The asset has been sold or transferred to another party.', 'system', 'system'),

    -- ✅ Wishlist Statuses
    ('Wishlist - Pending', 'Asset is added to the wishlist but not yet purchased.', 'system', 'system'),
    ('Wishlist - Purchased', 'Asset from the wishlist has been successfully acquired.', 'system', 'system'),
    ('Wishlist - Removed', 'Asset has been removed from the wishlist.', 'system', 'system'),

    -- ✅ Asset Assignment & Operational Statuses
    ('Reserved', 'The asset is reserved for a specific user or project.', 'system', 'system'),
    ('Checked Out', 'The asset is checked out to a user or department.', 'system', 'system'),
    ('Checked In', 'The asset has been returned and is available again.', 'system', 'system'),

    -- ✅ Asset Financial & Depreciation Statuses
    ('Under Warranty', 'The asset is covered under warranty for repairs or replacements.', 'system', 'system'),
    ('Warranty Expired', 'The asset is no longer covered under warranty.', 'system', 'system'),
    ('Depreciated', 'The asset has undergone financial depreciation.', 'system', 'system'),

    -- ✅ Miscellaneous Statuses
    ('Awaiting Disposal', 'The asset is pending disposal or recycling.', 'system', 'system'),
    ('Inactive', 'The asset is temporarily inactive but not retired.', 'system', 'system');

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
VALUES
    -- ✅ Electronics & IT Equipment
    ('Electronics', 'Electronic devices, gadgets, and computing equipment', 'system', 'system'),
    ('Computers & Laptops', 'Desktops, laptops, and tablets', 'system', 'system'),
    ('Mobile Phones', 'Smartphones, feature phones, and accessories', 'system', 'system'),
    ('Audio Equipment', 'Speakers, headphones, and microphones', 'system', 'system'),
    ('Video Equipment', 'Cameras, camcorders, and video production gear', 'system', 'system'),
    ('Networking Devices', 'Routers, switches, and access points', 'system', 'system'),
    ('Storage Devices', 'Hard drives, SSDs, and NAS storage', 'system', 'system'),

    -- ✅ Furniture & Office Equipment
    ('Furniture', 'Office and home furniture', 'system', 'system'),
    ('Office Desks', 'Workstations and office desks', 'system', 'system'),
    ('Office Chairs', 'Ergonomic and standard chairs for workspace', 'system', 'system'),
    ('Conference Tables', 'Tables used in meeting and conference rooms', 'system', 'system'),
    ('Cabinets & Shelves', 'Storage units including file cabinets and shelves', 'system', 'system'),

    -- ✅ Vehicles & Transport
    ('Vehicles', 'Company-owned vehicles for transport and logistics', 'system', 'system'),
    ('Company Cars', 'Sedans and SUVs used by employees', 'system', 'system'),
    ('Trucks', 'Heavy-duty vehicles for transportation', 'system', 'system'),
    ('Motorcycles', 'Two-wheeled vehicles for company use', 'system', 'system'),
    ('Electric Vehicles', 'Battery-powered electric vehicles', 'system', 'system'),

    -- ✅ Machinery & Tools
    ('Machinery', 'Industrial machinery and large equipment', 'system', 'system'),
    ('Hand Tools', 'Basic hand tools like hammers and screwdrivers', 'system', 'system'),
    ('Power Tools', 'Electric drills, saws, grinders, and similar tools', 'system', 'system'),

    -- ✅ Software & Hardware
    ('Software', 'Licensed software, applications, and digital assets', 'system', 'system'),
    ('Hardware', 'Computer peripherals and hardware components', 'system', 'system'),

    -- ✅ Appliances & Home Equipment
    ('Appliances', 'Electrical and non-electrical appliances for office or home use', 'system', 'system'),
    ('Kitchen Appliances', 'Refrigerators, microwaves, and coffee machines', 'system', 'system'),
    ('Cleaning Equipment', 'Vacuum cleaners, air purifiers, and sanitation devices', 'system', 'system'),

    -- ✅ Medical & Laboratory Equipment
    ('Medical Equipment', 'Healthcare and hospital equipment', 'system', 'system'),
    ('Laboratory Equipment', 'Scientific research and testing instruments', 'system', 'system'),

    -- ✅ Security & Surveillance
    ('Security Equipment', 'Cameras, alarms, and surveillance devices', 'system', 'system'),
    ('CCTV Cameras', 'Security cameras for surveillance', 'system', 'system'),
    ('Alarm Systems', 'Intruder detection and emergency alarms', 'system', 'system'),

    -- ✅ Office Supplies & Stationery
    ('Stationery', 'Office supplies like pens, paper, and notebooks', 'system', 'system'),
    ('Printers & Scanners', 'Devices for printing and scanning documents', 'system', 'system'),

    -- ✅ Real Estate & Properties
    ('Real Estate', 'Buildings, land, and company-owned properties', 'system', 'system'),

    -- ✅ Other Miscellaneous Categories
    ('Art', 'Paintings, sculptures, and decorative items', 'system', 'system'),
    ('Books & Reference', 'Books, manuals, and reference materials', 'system', 'system'),
    ('Clothing & Uniforms', 'Work uniforms and protective clothing', 'system', 'system'),
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
    warranty_expiry_date DATE DEFAULT NULL,
    price                DECIMAL(40, 2) DEFAULT 0,
    stock                INT            DEFAULT 0,
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
VALUES
    -- ✅ General Maintenance Types
    ('Battery Replacement', 'Replace worn-out battery with a new one to maintain performance.', 'system', 'system'),
    ('Software Update', 'Update firmware and OS for better performance and security.', 'system', 'system'),
    ('Cleaning Service', 'Perform deep cleaning to remove dust and debris from components.', 'system', 'system'),
    ('Hardware Repair', 'Fix or replace damaged hardware components of the asset.', 'system', 'system'),
    ('Annual Inspection', 'General check-up and inspection to ensure the asset is in optimal condition.', 'system',
     'system'),

    -- ✅ IT & Network Maintenance
    ('Firmware Upgrade', 'Upgrade device firmware to the latest version.', 'system', 'system'),
    ('Networking Maintenance', 'Check and repair network connections, cables, and routers.', 'system', 'system'),
    ('Security Patch Update', 'Apply the latest security updates and patches.', 'system', 'system'),
    ('Cloud Backup Check', 'Ensure that cloud backup systems are properly configured and up to date.', 'system',
     'system'),
    ('General Diagnostics', 'Perform full system diagnostics to detect potential issues.', 'system', 'system'),
    ('Performance Optimization', 'Optimize asset performance through software and hardware adjustments.', 'system',
     'system'),

    -- ✅ Electrical & Mechanical Maintenance
    ('Cooling System Check', 'Inspect and clean cooling fans and heat sinks to prevent overheating.', 'system',
     'system'),
    ('Electrical Testing', 'Ensure safe electrical operation and check for voltage stability.', 'system', 'system'),
    ('Oil & Lubrication', 'Apply lubrication to mechanical parts to prevent wear and tear.', 'system', 'system'),
    ('Parts Replacement', 'Replace broken or worn-out components with new parts.', 'system', 'system'),
    ('Sensor Calibration', 'Calibrate sensors to maintain accuracy and efficiency.', 'system', 'system'),
    ('Motor Servicing', 'Check and maintain motors in mechanical and industrial assets.', 'system', 'system'),

    -- ✅ Security & Surveillance Equipment Maintenance
    ('CCTV Camera Maintenance', 'Inspect and clean surveillance cameras for optimal performance.', 'system', 'system'),
    ('Alarm System Check', 'Test and verify the functionality of alarm systems.', 'system', 'system'),
    ('Fire Safety Inspection', 'Ensure fire safety equipment is in working condition.', 'system', 'system'),

    -- ✅ Vehicle & Transportation Maintenance
    ('Engine Tuning', 'Fine-tune vehicle engines for better efficiency and performance.', 'system', 'system'),
    ('Brake Inspection', 'Check and replace brake pads, fluids, and related components.', 'system', 'system'),
    ('Tire Replacement', 'Inspect and replace tires for better safety and performance.', 'system', 'system'),
    ('Fuel System Maintenance', 'Clean fuel injectors and ensure proper fuel flow.', 'system', 'system'),
    ('Transmission Check', 'Inspect transmission systems to prevent failure.', 'system', 'system'),

    -- ✅ Heavy Equipment & Industrial Machinery Maintenance
    ('Hydraulic System Inspection', 'Check and maintain hydraulic systems in heavy machinery.', 'system', 'system'),
    ('Welding & Structural Repair', 'Inspect and reinforce metal structures.', 'system', 'system'),
    ('Conveyor Belt Maintenance', 'Check for misalignment and damage in conveyor belts.', 'system', 'system'),

    -- ✅ Medical & Laboratory Equipment Maintenance
    ('Medical Device Calibration', 'Ensure medical equipment provides accurate readings.', 'system', 'system'),
    ('Sterilization Service', 'Sterilize medical and laboratory equipment to maintain hygiene.', 'system', 'system'),
    ('Oxygen System Check', 'Inspect and maintain oxygen supply systems.', 'system', 'system'),

    -- ✅ Office Equipment & Appliances Maintenance
    ('Printer & Scanner Service', 'Clean and repair printers, scanners, and copiers.', 'system', 'system'),
    ('Air Conditioner Service', 'Check and refill refrigerant, clean filters in AC systems.', 'system', 'system'),
    ('Refrigerator Maintenance', 'Ensure proper cooling and clean condenser coils in refrigerators.', 'system',
     'system'),

    -- ✅ Miscellaneous & General Maintenance
    ('Furniture Repair', 'Fix loose hinges, screws, and broken parts in furniture.', 'system', 'system'),
    ('Painting & Coating', 'Repaint assets to maintain aesthetic appeal and prevent rust.', 'system', 'system'),

    -- ✅ The "Other" Type (Must Always be the Last Entry)
    ('Other', 'Any other maintenance not covered in predefined types.', 'system', 'system');

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