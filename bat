# Auth Service
.\migration.bat auth create create_credentials_table
.\migration.bat auth up
.\migration.bat auth down
.\migration.bat auth status

# User Service
.\migration.bat user create create_users_table
.\migration.bat user up

# Project Service
.\migration.bat project create create_projects_table
.\migration.bat project up

# Deployment Service
.\migration.bat deployment up

# Build Service
.\migration.bat build up

# Runtime Service
.\migration.bat runtime up