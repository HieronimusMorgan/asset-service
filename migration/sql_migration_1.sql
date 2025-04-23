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
-- Indexes to speed up queries
CREATE INDEX idx_asset_status_id ON asset_status (asset_status_id);
CREATE INDEX idx_asset_status_name ON asset_status (status_name);

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
    user_client_id VARCHAR(50) NOT NULL,
    category_name     VARCHAR(255) NOT NULL,
    description       TEXT,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by        VARCHAR(255),
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by        VARCHAR(255),
    deleted_at        TIMESTAMP,
    deleted_by        VARCHAR(255)
);
-- Indexes to speed up queries
CREATE INDEX idx_asset_category_id ON asset_category (asset_category_id);
CREATE INDEX idx_asset_category_name ON asset_category (category_name);

-- Table for Asset
CREATE TABLE asset
(
    asset_id             SERIAL PRIMARY KEY,
    user_client_id       VARCHAR(50)  NOT NULL,
    serial_number        VARCHAR(100) DEFAULT NULL,
    name                 VARCHAR(100) NOT NULL,
    description          TEXT,
    barcode              VARCHAR(100)   DEFAULT NULL,
    category_id          INT          NOT NULL,
    status_id            INT          NOT NULL,
    purchase_date        DATE,
    expiry_date          DATE           DEFAULT NULL,
    warranty_expiry_date DATE         DEFAULT NULL,
    price                DECIMAL(40, 2) DEFAULT 0,
    stock                INT            DEFAULT 0,
    notes                TEXT         DEFAULT NULL,
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
-- Indexes to speed up queries
CREATE INDEX idx_asset_id ON asset (asset_id);
CREATE INDEX idx_asset_user ON asset (user_client_id);
CREATE INDEX idx_asset_serial ON asset (serial_number);
CREATE INDEX idx_asset_name ON asset (name);
CREATE INDEX idx_asset_barcode ON asset (barcode);
CREATE INDEX idx_asset_category ON asset (category_id);
CREATE INDEX idx_asset_status ON asset (status_id);

CREATE TABLE asset_stock
(
    stock_id         SERIAL PRIMARY KEY,
    asset_id         INT         NOT NULL,
    user_client_id   VARCHAR(50) NOT NULL,
    initial_quantity INT NOT NULL CHECK (initial_quantity >= 0),                           -- Can be 0 or more
    latest_quantity  INT NOT NULL CHECK (latest_quantity >= 0),                            -- Can be 0 or more
    change_type      VARCHAR(50) NOT NULL CHECK (change_type IN ('INCREASE', 'DECREASE')), -- Defines stock adjustments
    quantity         INT NOT NULL CHECK (quantity > 0),                                    -- The amount added or removed
    reason           TEXT DEFAULT NULL,                                                    -- Optional reason for stock adjustment
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);

-- Indexes to speed up queries
CREATE INDEX idx_asset_stock_id ON asset_stock (stock_id);
CREATE INDEX idx_asset_stock_asset ON asset_stock (asset_id);
CREATE INDEX idx_asset_stock_user ON asset_stock (user_client_id);

CREATE TABLE asset_image
(
    image_id       SERIAL PRIMARY KEY,
    user_client_id VARCHAR(50) NOT NULL,
    asset_id       INT         NOT NULL,
    image_url      TEXT        NOT NULL,
    file_type      VARCHAR(50) NOT NULL,
    file_size      BIGINT,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255),
    CONSTRAINT fk_asset FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE
);
-- Indexes to speed up queries
CREATE INDEX idx_asset_image_id ON asset_image (image_id);
CREATE INDEX idx_asset_image_user ON asset_image (user_client_id);
CREATE INDEX idx_asset_image_asset ON asset_image (asset_id);

CREATE TABLE asset_stock_history
(
    stock_history_id  SERIAL PRIMARY KEY,
    asset_id          INT         NOT NULL,
    user_client_id    VARCHAR(50) NOT NULL,
    stock_id          INT         NOT NULL,                                -- Reference to asset_stock
    change_type       VARCHAR(50) NOT NULL CHECK (change_type IN ('INCREASE', 'DECREASE', 'ADJUSTMENT')),
    previous_quantity INT         NOT NULL CHECK (previous_quantity >= 0), -- Before update
    new_quantity      INT         NOT NULL CHECK (new_quantity >= 0),      -- After update
    quantity_changed  INT         NOT NULL CHECK (quantity_changed > 0),   -- Change difference
    reason            TEXT      DEFAULT NULL,                              -- Optional adjustment reason
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by        VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id) ON DELETE CASCADE,
    FOREIGN KEY (stock_id) REFERENCES asset_stock (stock_id) ON DELETE CASCADE
);

-- Indexes to speed up queries
CREATE INDEX idx_asset_stock_history_id ON asset_stock_history (stock_history_id);
CREATE INDEX idx_asset_stock_history_stock ON asset_stock_history (stock_id);
CREATE INDEX idx_asset_stock_history_asset ON asset_stock_history (asset_id);
CREATE INDEX idx_asset_stock_history_user ON asset_stock_history (user_client_id);

-- Table for Maintenance Type
CREATE TABLE asset_maintenance_type
(
    maintenance_type_id   SERIAL PRIMARY KEY,           -- Unique ID for each type
    user_client_id        VARCHAR(50)         NOT NULL,
    maintenance_type_name VARCHAR(100) UNIQUE NOT NULL, -- Name of maintenance type
    description           TEXT,                         -- Description of what this maintenance involves
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by            VARCHAR(255),
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by            VARCHAR(255),
    deleted_at            TIMESTAMP,
    deleted_by            VARCHAR(255)
);

CREATE INDEX idx_maintenance_type_id ON asset_maintenance_type (maintenance_type_id);
CREATE INDEX idx_maintenance_type_name ON asset_maintenance_type (maintenance_type_name);
CREATE INDEX idx_maintenance_type_user ON asset_maintenance_type (user_client_id);

-- Table for Asset Maintenance
CREATE TABLE asset_maintenance
(
    id                  SERIAL PRIMARY KEY,
    user_client_id      VARCHAR(50) NOT NULL,
    asset_id            INT         NOT NULL, -- Reference to asset
    maintenance_type_id INT         NOT NULL,
    maintenance_date    DATE        NOT NULL,
    maintenance_details TEXT,
    maintenance_cost    DECIMAL(15, 2),       -- Cost of maintenance
    performed_by        VARCHAR(255),         -- Who performed the maintenance
    interval_days       INT,                  -- Maintenance interval in days
    next_due_date       DATE DEFAULT NULL,    -- Scheduled next maintenance
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by          VARCHAR(255),
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by          VARCHAR(255),
    deleted_at          TIMESTAMP,
    deleted_by          VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
    FOREIGN KEY (maintenance_type_id) REFERENCES asset_maintenance_type (maintenance_type_id) ON DELETE SET NULL
);

CREATE INDEX idx_maintenance_id ON asset_maintenance (id);
CREATE INDEX idx_maintenance_asset ON asset_maintenance (asset_id);
CREATE INDEX idx_maintenance_type ON asset_maintenance (maintenance_type_id);

CREATE TABLE asset_maintenance_record
(
    maintenance_record_id SERIAL PRIMARY KEY,
    user_client_id      VARCHAR(50) NOT NULL,
    asset_id            INT         NOT NULL,   -- Reference to asset
    maintenance_id      INT         NOT NULL,   -- Reference to asset_maintenance
    maintenance_type_id INT         NOT NULL,
    maintenance_date    DATE        NOT NULL,
    maintenance_details TEXT,
    maintenance_cost    DECIMAL(15, 2),         -- Cost of maintenance
    performed_by        VARCHAR(255),           -- Who performed the maintenance
    interval_days       INT,                    -- Maintenance interval in days
    next_due_date       DATE      DEFAULT NULL, -- Scheduled next maintenance
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by            VARCHAR(255),
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by            VARCHAR(255),
    deleted_at            TIMESTAMP,
    deleted_by            VARCHAR(255),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
    FOREIGN KEY (maintenance_id) REFERENCES asset_maintenance (id),
    FOREIGN KEY (maintenance_type_id) REFERENCES asset_maintenance_type (maintenance_type_id)
);

-- Indexes to speed up queries
CREATE INDEX idx_maintenance_record_id ON asset_maintenance_record (maintenance_record_id);
CREATE INDEX idx_maintenance_record_user ON asset_maintenance_record (user_client_id);
CREATE INDEX idx_maintenance_record_asset ON asset_maintenance_record (asset_id);
CREATE INDEX idx_maintenance_record_maintenance ON asset_maintenance_record (maintenance_id);
CREATE INDEX idx_maintenance_record_type ON asset_maintenance_record (maintenance_type_id);

CREATE TABLE asset_group
(
    asset_group_id   SERIAL PRIMARY KEY,
    asset_group_name VARCHAR(255) NOT NULL,
    description      TEXT,
    owner_user_id    INT          NOT NULL,
    invitation_token VARCHAR(100) UNIQUE DEFAULT NULL,
    max_uses         INT                 DEFAULT NULL,
    current_uses     INT                 DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255)
);
-- Indexes to speed up queries
CREATE INDEX idx_asset_group_id ON asset_group (asset_group_id);
CREATE INDEX idx_asset_group_owner ON asset_group (owner_user_id);


ALTER TABLE asset_group
    ADD CONSTRAINT chk_usage_valid
        CHECK (current_uses >= 0 AND (max_uses IS NULL OR max_uses >= current_uses));

CREATE TABLE asset_group_member
(
    asset_group_id INT NOT NULL,
    user_id        INT NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255),
    PRIMARY KEY (user_id, asset_group_id),
    FOREIGN KEY (asset_group_id) REFERENCES asset_group (asset_group_id)
);
CREATE INDEX idx_asset_group_member_user ON asset_group_member (user_id);
CREATE INDEX idx_asset_group_member_group ON asset_group_member (asset_group_id);

CREATE TABLE asset_group_asset
(
    asset_group_id INT NOT NULL,
    asset_id       INT NOT NULL,
    user_id        INT NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255),
    PRIMARY KEY (asset_id, asset_group_id, user_id),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
    FOREIGN KEY (asset_group_id) REFERENCES asset_group (asset_group_id)
);
CREATE INDEX idx_asset_group_asset_user ON asset_group_asset (user_id);
CREATE INDEX idx_asset_group_asset_group ON asset_group_asset (asset_group_id);
CREATE INDEX idx_asset_group_asset_asset ON asset_group_asset (asset_id);


CREATE TABLE asset_group_permission
(
    permission_id   SERIAL PRIMARY KEY,
    permission_name VARCHAR(100) UNIQUE NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by      VARCHAR(255),
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMP,
    deleted_by      VARCHAR(255)
);
-- Indexes to speed up queries
CREATE INDEX idx_asset_group_permission_id ON asset_group_permission (permission_id);
CREATE INDEX idx_asset_group_permission_name ON asset_group_permission (permission_name);
-- Insert default family_permission
INSERT INTO asset_group_permission (permission_name, description)
VALUES ('Admin', 'Full control over family members and assets'),
       ('Manage', 'Manage family members and permissions'),
       ('Read-Write', 'Read and Write access to assets'),
       ('Read', 'Read/View assets');

CREATE TABLE asset_group_member_permission
(
    asset_group_id INT NOT NULL,
    user_id        INT NOT NULL,
    permission_id  INT NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255),
    PRIMARY KEY (user_id, asset_group_id, permission_id)
);
CREATE INDEX idx_asset_group_member_permission_user ON asset_group_member_permission (user_id);
CREATE INDEX idx_asset_group_member_permission_group ON asset_group_member_permission (asset_group_id);
CREATE INDEX idx_asset_group_member_permission_permission ON asset_group_member_permission (permission_id);

CREATE TABLE asset_group_invitation
(
    invitation_id      SERIAL PRIMARY KEY,
    asset_group_id     INT         NOT NULL,
    invited_user_id    INT         NOT NULL,                   -- user being invited
    invited_user_token VARCHAR(100) UNIQUE,
    invited_by_user_id INT         NOT NULL,                   -- user who sent the invitation
    status             VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'expired'
    message            TEXT,                                   -- optional message included in the invitation
    invited_at         TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    responded_at       TIMESTAMP,                              -- when accepted/rejected
    expired_at TIMESTAMP,                                      -- when the invitation expires
    created_at         TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    created_by         VARCHAR(255),
    updated_at         TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    updated_by         VARCHAR(255),
    deleted_at         TIMESTAMP,
    deleted_by         VARCHAR(255),

    FOREIGN KEY (asset_group_id) REFERENCES asset_group (asset_group_id)
);
CREATE INDEX idx_asset_group_invitation_user ON asset_group_invitation (invited_user_id);
CREATE INDEX idx_asset_group_invitation_group ON asset_group_invitation (asset_group_id);
CREATE INDEX idx_asset_group_invitation_status ON asset_group_invitation (status);
CREATE INDEX idx_asset_group_invitation_token ON asset_group_invitation (invited_user_token);
CREATE INDEX idx_asset_group_invitation_invited_by ON asset_group_invitation (invited_by_user_id);

ALTER TABLE asset_group_invitation
    ADD CONSTRAINT chk_invite_status
        CHECK (status IN ('pending', 'accepted', 'rejected', 'expired'));

CREATE TABLE asset_tags
(
    tag_id      SERIAL PRIMARY KEY,
    tag_name    VARCHAR(255) NOT NULL UNIQUE, -- Name of the tag
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255)
);
CREATE INDEX idx_asset_tags_id ON asset_tags (tag_id);
CREATE INDEX idx_asset_tags_name ON asset_tags (tag_name);

CREATE TABLE asset_tag_map
(
    asset_id INT NOT NULL,
    tag_id   INT NOT NULL,
    PRIMARY KEY (asset_id, tag_id),
    FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
    FOREIGN KEY (tag_id) REFERENCES asset_tags (tag_id)
);
CREATE INDEX idx_asset_tag_map_asset ON asset_tag_map (asset_id);
CREATE INDEX idx_asset_tag_map_tag ON asset_tag_map (tag_id);

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
CREATE INDEX idx_asset_audit_log_table ON asset_audit_log (table_name);

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
-- Indexes to speed up queries
CREATE INDEX idx_cron_jobs_id ON cron_jobs (id);

INSERT INTO cron_jobs (name, schedule, is_active, description, created_by)
VALUES ('asset_maintenance', '0 5 * * *', true, 'Check and schedule maintenance for assets', 'system'),
       ('asset_image_cleanup', '*/1 * * * *', true, 'Cleanup old and unused asset images', 'system'),
       ('image_cleanup_unused', '*/1 * * * *', true, 'Remove unused images from asset storage', 'system');

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