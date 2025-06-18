# My First Go URL Shortener

This is my very first project built with Go! It's a simple URL shortening service that converts long URLs into short, shareable links, featuring both a RESTful API and an integrated web interface.

## Technologies Used

* **Backend:**
    * Go (using `net/http` for server)
    * [Bun](https://bun.uptrace.dev/) (ORM for database interaction)

* **Database:** MySQL
* **Frontend:** HTML, CSS, JavaScript (served directly by the Go application)

## How to Use

Follow these steps to get the URL Shortener running on your local machine:

### Prerequisites

* Go (1.21+) installed
* MySQL server running locally

### Setup & Run

1.  **Database Setup:**
    * Connect to your MySQL server (e.g., `mysql -u root -p`).
    * Create a database (e.g., `url_parserer_db`).
    * Create a user (e.g., `bunuser` with password `bunpassword`) and grant all privileges on `url_parserer_db` to this user.
        ```sql
        CREATE DATABASE IF NOT EXISTS url_parserer_db;
        CREATE USER 'bunuser'@'localhost' IDENTIFIED BY 'bunpassword';
        GRANT ALL PRIVILEGES ON url_parserer_db.* TO 'bunuser'@'localhost';
        FLUSH PRIVILEGES;
        ```

2.  **Clone Repository:**
    ```bash
    git clone [https://github.com/your-username/your-repo-name.git](https://github.com/your-username/your-repo-name.git) # Replace with your actual repo
    cd your-repo-name
    ```

3.  **Create `.env` file:**
    In the root directory of the project, create a file named `.env` and populate it with your database credentials and application settings:
    ```dotenv
    # Application Configuration
    PORT=8080
    BASE_URL="http://localhost:8080/"

    # Database Configuration (MySQL)
    DSN="bunuser:bunpassword@tcp(127.0.0.1:3306)/url_parserer_db?parseTime=true&loc=Local"
    ```
    **Remember to replace `bunuser` and `bunpassword` with your actual MySQL credentials.**

4.  **Install Dependencies & Run:**
    ```bash
    go mod tidy
    go run main.go
    ```
    The application will start, connect to the database, create tables (if they don't exist), and listen on the configured port (default: `8080`).

## API Endpoints & Web Interface

Once the application is running:

1.  **Access the Web Interface:**
    Open your web browser and go to `http://localhost:8080/`. You'll see a simple interface to shorten URLs and look up existing ones.

2.  **Shorten URL (POST Request):**
    * **URL:** `http://localhost:8080/api/set`
    * **Method:** `POST`
    * **Headers:** `Content-Type: application/json`
    * **Body:** `{"url_string": "your_long_url_here"}`
    * **Example (`curl`):**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"url_string": "[https://www.google.com](https://www.google.com)"}' http://localhost:8080/api/set
        ```

3.  **Retrieve Original URL Details (GET Request):**
    * **URL:** `http://localhost:8080/api/get?url={shortCode}`
    * **Method:** `GET`
    * **Example (`curl`):**
        ```bash
        curl "http://localhost:8080/api/get?url=yWnYLAY"
        ```

4.  **Redirect Short URL (GET Request):**
    * **URL:** `http://localhost:8080/{shortCode}`
    * **Method:** `GET`
    * Access this directly in your browser or with `curl -L` (to follow redirects).
    * **Example (`curl`):**
        ```bash
        curl -v "http://localhost:8080/yWnYLAY"
        ```

---
**Note:** This is my initial foray into Go programming, and I'm learning along the way! Feedback and suggestions are welcome.
