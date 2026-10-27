# Movie Database Microservice API

A high-performance RESTful microservice built with **Go**, **Gin**, and **MongoDB**, utilizing the **Controller-Service-DAO** architectural pattern. It features relational aggregation pipelines (`$lookup`, `$unwind`), nested JSON responses (DTOs), name-based searching, custom pagination (`start`/`skip`, `limit`), and field projection (`fields`, `exclude`).

---

## Project Directory Structure

```plaintext
.
├── cmd
│   └── server
│       ├── main.go
│       └── router.go
├── controllers
│   ├── actor_controller.go
│   ├── director_controller.go
│   └── film_controller.go
├── dao
│   ├── actor_dao.go
│   ├── director_dao.go
│   └── film_dao.go
├── go.mod
├── go.sum
├── models
│   ├── actors.go
│   ├── directors.go
│   ├── films.go
│   ├── films_pop.go
│   ├── query_options.go
│   └── setup.go
└── service
    ├── actor_service.go
    ├── director_service.go
    └── film_service.go
```

---

## Prerequisites

- **Go** (v1.18 or higher recommended)
- **MongoDB** (running locally or via a cloud instance like MongoDB Atlas)
- **HTTPie** or **cURL** (for testing endpoints)

---

## Getting Started & Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd task3-moviesAPI
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Configure Environment Variables:**  
   Create a `.env` file in the root directory (or configure your environment variables directly):
   ```env
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=movies_db
   ```

4. **Run the Application:**
   ```bash
   go run cmd/server/main.go
   ```
   The server will start on port `8080` (or your configured `PORT`).

---

## API Endpoints & Usage

### 🎬 Films (Populated & Search Endpoints)

#### 1. Get All Populated Films
Returns all films with nested director objects and full actor arrays instead of raw IDs.

- **Endpoint:** `GET /api/v1/films/populated`
- **Example (HTTPie):**
  ```bash
  http GET localhost:8080/api/v1/films/populated
  ```

---

#### 2. Search Films by Director Name
Performs a case-insensitive search by director name with support for pagination and field projections.

- **Endpoint:** `GET /api/v1/films/search/director`
- **Query Parameters:**
  - `name` (Required): Director's name (e.g., `Christopher Nolan`)
  - `start` or `skip` (Optional): Number of elements to skip (default: `0`)
  - `limit` (Optional): Maximum number of elements to return
  - `fields` (Optional): Comma-separated list of fields to include (e.g., `title,year`)
  - `exclude` (Optional): Comma-separated list of fields to exclude (e.g., `actors`)
- **Example (HTTPie):**
  ```bash
  http GET "localhost:8080/api/v1/films/search/director?name=Christopher%20Nolan&start=0&limit=5&fields=title,year"
  ```

---

#### 3. Search Films by Actor Name
Performs a case-insensitive search across associated actors with support for pagination and field projections.

- **Endpoint:** `GET /api/v1/films/search/actor`
- **Query Parameters:**
  - `name` (Required): Actor's name (e.g., `Matthew McConaughey`)
  - `start` or `skip` (Optional): Number of elements to skip
  - `limit` (Optional): Maximum number of elements to return
  - `fields` (Optional): Fields to include (e.g., `title,year,director`)
  - `exclude` (Optional): Fields to exclude (e.g., `director`)
- **Example (HTTPie):**
  ```bash
  http GET "localhost:8080/api/v1/films/search/actor?name=Matthew&exclude=actors"
  ```

---

## Inserting Sample Data via `mongosh`

To quickly populate mock data to test your search and pagination filters, open your MongoDB shell (`mongosh`) and run:

```javascript
use movies_db;

const directorId = new ObjectId();
const actor1Id = new ObjectId();
const actor2Id = new ObjectId();

db.directors.insertOne({
  _id: directorId,
  name: "Christopher Nolan",
  bio: "British-American filmmaker."
});

db.actors.insertMany([
  { _id: actor1Id, name: "Matthew McConaughey", biography: "Actor" },
  { _id: actor2Id, name: "Anne Hathaway", biography: "Actress" }
]);

db.films.insertOne({
  title: "Interstellar",
  year: 2014,
  director_id: directorId,
  actor_ids: [actor1Id, actor2Id]
});

print("Mock data seeded successfully!");
```