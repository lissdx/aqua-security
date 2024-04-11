CREATE TYPE "resource_type" AS ENUM (
    'repository',
    'image'
    );

CREATE TYPE "repository_source_type" AS ENUM (
    'github',
    'bitbucket',
    'gitlab'
    );

CREATE TYPE "image_source_type" AS ENUM (
    'dockerhub',
    'ecr',
    'jfrog'
    );

CREATE TYPE "architecture_type" AS ENUM (
    'arm',
    'amd'
    );

CREATE TYPE "highest_severity" AS ENUM (
    'high',
    'medium',
    'low'
    );

CREATE TABLE "input_repository" (
       "indx" bigserial PRIMARY KEY NOT NULL,
       "id" uuid UNIQUE NOT NULL,
       "name" varchar(250) NOT NULL,
       "url" varchar(1024) NOT NULL,
       "type" resource_type NOT NULL DEFAULT 'repository',
       "source" repository_source_type NOT NULL,
       "size" bigint NOT NULL,
       "last_push" timestamptz NOT NULL,
       "created_date_timestamp" timestamptz NOT NULL
);

ALTER TABLE input_repository ADD CONSTRAINT input_repository_check_length_name CHECK (length(name) >= 3 AND length(name) <= 250);
ALTER TABLE input_repository ADD CONSTRAINT input_repository_check_length_url CHECK (length(url) >= 8 AND length(url) <= 1024);
ALTER TABLE input_repository ADD CONSTRAINT input_repository_check_size CHECK (size>0);

CREATE TABLE "input_image" (
        "indx" bigserial PRIMARY KEY NOT NULL,
        "id" uuid UNIQUE NOT NULL,
        "repository_id" uuid,
        "name" varchar(250) NOT NULL,
        "url" varchar(1024) NOT NULL,
        "type" resource_type NOT NULL DEFAULT 'image',
        "source" image_source_type NOT NULL,
        "number_of_layers" integer NOT NULL,
        "architecture" architecture_type NOT NULL,
        "created_date_timestamp" timestamptz NOT NULL
);

ALTER TABLE input_image ADD CONSTRAINT input_image_check_length_name CHECK (length(name) >= 3 AND length(name) <= 250);
ALTER TABLE input_image ADD CONSTRAINT input_image_check_length_url CHECK (length(url) >= 8 AND length(url) <= 1024);
ALTER TABLE input_image ADD CONSTRAINT input_image_check_number_of_layers CHECK (number_of_layers > 0);


CREATE TABLE "scan" (
       "indx" bigserial PRIMARY KEY NOT NULL,
       "scan_id" integer UNIQUE NOT NULL,
       "resource_id" uuid NOT NULL,
       "resource_type" resource_type NOT NULL,
       "highest_severity" highest_severity NOT NULL,
       "total_findings" integer NOT NULL,
       "is_reported" boolean NOT NULL default false,
       "scan_date_timestamp" timestamptz NOT NULL
);

CREATE UNIQUE INDEX "scan_resource" ON "scan" ("scan_id", "resource_id");
