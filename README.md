# BSEStarMF Go Project

This project provides a simple Go HTTP server that interacts with the BSE Star MF demo API for authentication.

## How to Run the Project

1.  **Prerequisites:**
    *   Go (version 1.24.5 or higher)

2.  **Clone the Repository (if not already cloned):**
    ```bash
git clone <repository_url>
cd BSEStarMF
    ```

3.  **Install Dependencies:**
    Navigate to the project root directory and run:
    ```bash
go mod tidy
    ```

4.  **Run the Application:**
    You can run the application directly from the `cmd/bsemf` directory:
    ```bash
go run cmd/bsemf/main.go
    ```
    Alternatively, you can build an executable and then run it:
    ```bash
go build -o bsemf_app cmd/bsemf/main.go
./bsemf_app
    ```

5.  **Access the Application:**
    The server will start on port `8080`. You can access the API using a tool like `curl` or Postman.

## API Endpoints

### `POST /api/auth`

*   **Description:** Authenticates a user against the BSE Star MF demo API. It sends a SOAP request with the provided credentials and returns the response from the BSE API.
*   **URL:** `http://localhost:8080/api/auth`
*   **Method:** `POST`
*   **Request Body (JSON):**
    ```json
{
  "user_id": "<YOUR_USER_ID>",
  "password": "<YOUR_PASSWORD>",
  "pass_key": "<YOUR_PASS_KEY>"
}
    ```
    *Replace `<YOUR_USER_ID>`, `<YOUR_PASSWORD>`, and `<YOUR_PASS_KEY>` with your actual credentials.*

*   **Response Body (JSON):**
    ```json
{
  "code": "<RESPONSE_CODE>",
  "encrypted_password": "<ENCRYPTED_PASSWORD>"
}
    ```
    *   `code`: A status code from the BSE API (e.g., "101" for success, or other codes for errors).
    *   `encrypted_password`: The encrypted password returned by the BSE API.