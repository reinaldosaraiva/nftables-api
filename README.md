## Go API for managing nftables
This is a Go API that provides a RESTful interface for managing nftables, a powerful and flexible firewall framework in the Linux kernel.

### Features
- CRUD operations for nftables rules: Create, Read, Update, and Delete rules easily through the API.
- Easy to use: Simple and intuitive API design ensures a smooth and efficient development experience.

### Example .env
It is necessary to create a .env file in the root directory of the project.
```
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=nftablesdb
WEB_SERVER_PORT=8000
JWT_SECRET=secret
JWT_EXPIRE_IN=30000
```