-- db/mssql/init.sql
-- SQL Server initialization script: create users table if not exists

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
BEGIN
    CREATE TABLE users (
        id NVARCHAR(255) PRIMARY KEY,
        name NVARCHAR(255) NOT NULL
    );
END
