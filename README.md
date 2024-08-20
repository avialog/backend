# Avialog backend

## Golang-powered backend application for mobile app that allows airplane pilots to record their flights and track career.

## API Reference 🚧🚧🚧 **IN PROGRESS** 🚧🚧🚧

EC2 hostname: `ec2-13-48-148-72.eu-north-1.compute.amazonaws.com`

#### Get server health info

```http
  GET /api/info
```
##

#### Get user profile [**JWT REQUIRED**]

```http
  GET /api/profile
```
##

#### Update user profile [**JWT REQUIRED**]

```http
  PUT /api/profile
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `userRequest`| `json` | **Required**. JSON body of user profile information to update.`


##

#### Get user contacts [**JWT REQUIRED**]

```http
  GET /api/contacts
```

This endpoint retrieves a list of contacts for a user. The user ID is obtained from the context.
##

#### Insert a new contact [**JWT REQUIRED**]

```http
  POST /api/contacts
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `contactRequest`| `json` | **Required**. JSON body of contact information to insert.

##

#### Update an existing contact [**JWT REQUIRED**]

```http
  PUT /api/contacts/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id` | `int` | **Required**. Id of contact to update.
| `contactRequest`| `json` | **Required**. JSON body of contact information to update.

##

#### Delete an existing contact [**JWT REQUIRED**]

```http
   DELETE /api/contacts/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`| `int` | **Required**.  Id of contact to delete`




## Technologies and libraries used:

- 🐹 **Go (Golang)**: Efficient and performant language for backend development.
- 🍸 **GIN Framework**: Lightweight HTTP framework for building APIs.
- 🗃️ **PostgreSQL (Database)**: Utilized for persistent data storage.
- 🔐 **Firebase ApiKeyAuth**: Token-based authentication for secure communication.
- 📬 **Postman**: API development and testing tool.
- 💻 **IntelliJ GoLand**: Integrated development environment for Go.
-  ☘ **Ginkgo**: Ginkgo is a testing framework for Go designed to help you write expressive tests.
- Ω   **Gomega**: Gomega is a matcher/assertion library. It is best paired with the Ginkgo BDD test framework, but can be adapted for use in other contexts too.
- ♻  **Swaggo package**: Swag converts Go annotations to Swagger Documentation 2.0.  
- 🔑 **Validator package**: Validator is a validation library for Go that ensures data integrity and adherence to specified rules and constraints. It's used to validate user inputs, API requests, and other data structures within the application, ensuring they meet defined criteria before processing. This helps maintain data consistency, security, and overall application reliability.
-  🖨**Gomock**: gomock is a mocking framework for the Go programming language. It integrates well with Go's built-in testing package, but can be used in other contexts too.
- 🎫 **golangci-lint**:  golangci-lint is a powerful static analysis tool for Go that helps identify and fix various code issues, ensuring code quality and adherence to best practices.

## Firebase authorization simplified architecture:
![auth](https://github.com/avialog/backend/assets/7630626/e522040a-cecb-4f65-8c2e-346214e2a561)


## Database schema:

![db](https://github.com/avialog/backend/assets/7630626/9919821f-b24c-4215-b7f4-719b0ca96ace)


