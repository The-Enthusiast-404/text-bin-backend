# Text-Bin API

The Text-Bin API allows users to store and retrieve text snippets, providing functionalities such as creating new texts, fetching existing texts, and performing health checks on the service.

## Base URL

The base URL for the API is: `https://textbin.theenthusiast.dev/v1/`

## Authentication

There's no authentication setup at the moment. All endpoints are public.

## Endpoints

### Health Check

- **GET** `/healthcheck`
  - **Description**: Checks the health of the API service.
  - **Response**: `200 OK` if the service is healthy.

### Text Management

#### Create Text

- **POST** `/texts`
  - **Description**: Creates a new text snippet.
  - **Body**:
    ```json
    {
      "title": "string",
      "content": "string",
      "expires_in": "int",
      "expires_unit": "string",
      "format": "string",
    }
    ```
  - **Response**: `201 Created` with the created text snippet.

#### Fetch Text

- **GET** `/texts/:id`
  - **Description**: Fetches a text snippet by its ID.
  - **Response**: `200 OK` with the text snippet.

#### Update Text

- **PATCH** `/texts/:id`
  - **Description**: Updates an existing text snippet.
  - **Body**:
    ```json
    {
      "title": "string",
      "content": "string",
      "expires_in": "int",
      "expires_unit": "string",
      "format": "string",
    }
    ```
  - **Response**: `200 OK` with the updated text snippet.

#### Delete Text

- **DELETE** `/texts/:id`
  - **Description**: Deletes a text snippet by its ID.
  - **Response**: `204 No Content`.


## Examples



## Development

To set up a local development environment:

1. Clone the repository.
2. Install dependencies: `go mod tidy`.
3. Start the server: `go run cmd/api/`.

## Deployment

Refer to the [Makefile](Makefile) for deployment commands and the [setup](remote/setup/01.sh) script for initial server setup. (Will be added soon)

## Contributing

We welcome contributions! Please see our [contribution guidelines](CONTRIBUTING.md) for details.

## Support

If you have any questions or need support, please [open an issue](https://github.com/your/repository/issues) on GitHub.
