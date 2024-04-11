### Aqua Security Home Assignment

#### Main Third Part Tools was used
1. github.com/fsnotify/fsnotify  (filesystem watchdog)
2. github.com/xeipuuv/gojsonschema (JSON schema validator)
3. github.com/sqlc-dev/sqlc  (sql code generator)
4. github.com/golang-migrate (DB migrate)

### General Instructions
1. Start postgresql DB  
   run `make local_db_up` or `docker-compose -f docker-compose-pg.yml up -d`
2. Start the monitoring app
   1. Easiest way
   run `make run_monitoring` or `docker-compose -f as-monitoring.yml up -d`
   NOTE: make sure that your mounted dir existed (check [as-monitoring.yml](as-monitoring.yml))
   for example run `mkdir -p test_dir/input && mkdir -p test_dir/output`
   2. Another way
      1. For the first time run you should init the DB
      so run the _aqua-security/cmd/migration/migration.go_
      2. Run the monitoring app _aqua-security/cmd/monitoring/monitoring.go_
      NOTE: before you start the monitoring app please export env (or set in your IDE)
      ```
      export GO_ENV=development_local
      export PROCESS_CONF_SUB_FOLDER=monitoring-process
      export PROCESS_NAME=DATA_UPDATER
      ```

### How It Works
**Important: The validators are configured to the highest mode**   

1. We are available to upload files into our WATCHING_DIR_PATH dir in recommended order
   1. resources[*].json (creates repository and image data)
   2. connection[*].json (paring image to repository)
   3. scans[*].json (create the scan data)  
   NOTE: any scans[*].json uploading will be triggered the output generation process

#### Key Files/Folders description
1. DB Init [./db/migration](db/migration)
2. SQL queries [./db/query/as_sql.sql](db/query/as_sql.sql)
3. main process file [./internal/processes/data_updater_process/data_updater_process.go](internal/processes/data_updater_process/data_updater_process.go)
4. config file [./configs/monitoring-process/development_local.env](configs/monitoring-process/development_local.env)
5. JSON schema validators [./configs/jsonschema](configs/jsonschema)

#### DB Insert/Upsert Policy Important Points
see [as_sql.sql](db/query/as_sql.sql)
1. table **input_repository**
   - the row will be updated if **last_push** is highest 
     (input_repository.last_push & input_repository.size fields will be updated only)
2. **"input_images"** the row will NOT be updated at all except the field _repository_id_  
   this one can be updated with connections[*].json
3. **scan** 
   - will NOT be updated at all exclude the field _is_reported_  
   the field is updated by app if the report for the scan was generated
