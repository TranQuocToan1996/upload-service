# upload-service
Upload service with jwt auth

- Simple authentication using JWT token (with HS256)
    - Register
    - Login
    - Revoke token by time
    - Upload file image

- To do:
    - Create your own MySQL instance
    - Fix .env in api and migrations folders (MySQL, port, etc...)
    - cd to migrations folder, run command to update MySQL tables
    ```
    go run *.go
    ```
    - Import Postman collection in the docs folder
    - cd to api folder run
    ```
    go run *.go
    ```