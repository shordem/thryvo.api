-- Table for storing users
CREATE TABLE "users" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "is_email_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "password" VARCHAR NOT NULL,
    "role" VARCHAR NOT NULL
);

-- Table for storing verification codes
CREATE TABLE "verification_codes" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "user_id" UUID NOT NULL,
    "code" VARCHAR NOT NULL,
    "purpose" VARCHAR NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

-- Table for storing keys
CREATE TABLE "keys" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "user_id" UUID NOT NULL,
    "key" VARCHAR NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

-- Table for storing folders
CREATE TABLE "folders" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "user_id" UUID NOT NULL,
    "name" VARCHAR NOT NULL,
    "parent_id" UUID NULL,
    FOREIGN KEY ("parent_id") REFERENCES "folders" ("id") ON DELETE SET NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

-- Table for storing files
CREATE TABLE "files" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    "user_id" UUID NOT NULL,
    "folder_id" UUID NULL,
    "key" VARCHAR NOT NULL,
    "original_name" VARCHAR NOT NULL,
    "mime_type" VARCHAR NOT NULL,
    "size" BIGINT NOT NULL,
    "visibility" VARCHAR NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("folder_id") REFERENCES "folders" ("id") ON DELETE CASCADE
);

-- Create indexes to optimize queries
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_verification_codes_user_id ON verification_codes(user_id);
CREATE INDEX idx_keys_user_id ON keys(user_id);
CREATE INDEX idx_folders_user_id ON folders(user_id);
CREATE INDEX idx_folders_parent_id ON folders(parent_id);
CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_folder_id ON files(folder_id);

-- Create partial indexes for querying non-deleted records
CREATE INDEX idx_verification_codes_not_deleted ON verification_codes(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_keys_not_deleted ON keys(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_folders_not_deleted ON folders(user_id, parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_files_not_deleted ON files(user_id, folder_id) WHERE deleted_at IS NULL;

-- Create a partial unique index for non-deleted users
CREATE UNIQUE INDEX idx_users_not_deleted ON users(email) WHERE deleted_at IS NULL;