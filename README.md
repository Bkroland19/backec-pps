# Point Prevalence Survey API

A Go Gin-based REST API for managing Point Prevalence Survey data with PostgreSQL integration and Swagger documentation.

## Features

-    **CSV Data Import**: Upload and import data from CSV files for patients, antibiotics, indications, optional variables, and specimens
-    **RESTful API**: Complete CRUD operations with proper relationships and joins
-    **PostgreSQL Integration**: Robust database operations with GORM
-    **Swagger Documentation**: Interactive API documentation
-    **Docker Support**: Easy deployment with Docker and Docker Compose
-    **Filtering & Pagination**: Advanced querying capabilities
-    **Statistics**: Comprehensive analytics and reporting

## Project Structure

```
Point Prevalence Survey/
├── config/                 # Configuration management
├── database/              # Database connection and migrations
├── docs/                  # Swagger documentation
├── handlers/              # HTTP request handlers
├── models/                # Database models
├── routes/                # API route definitions
├── services/              # Business logic services
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose setup
└── README.md             # This file
```

## Data Models

The API handles five main data entities:

1. **Patients**: Main patient records with demographic and clinical information
2. **Antibiotics**: Antibiotic usage data linked to patients
3. **Indications**: Treatment indications and diagnoses
4. **Optional Variables**: Additional treatment variables
5. **Specimens**: Microbiology specimen data

## API Endpoints

### Patients

-    `GET /api/v1/patients` - List all patients with filtering
-    `GET /api/v1/patients/{id}` - Get specific patient with all related data
-    `GET /api/v1/patients/{id}/antibiotics` - Get patient's antibiotics
-    `GET /api/v1/patients/{id}/indications` - Get patient's indications
-    `GET /api/v1/patients/{id}/optional-vars` - Get patient's optional variables
-    `GET /api/v1/patients/{id}/specimens` - Get patient's specimens
-    `GET /api/v1/patients/stats` - Get patient statistics

### Antibiotics

-    `GET /api/v1/antibiotics` - List all antibiotics with filtering
-    `GET /api/v1/antibiotics/{id}` - Get specific antibiotic
-    `GET /api/v1/antibiotics/stats` - Get antibiotic usage statistics
-    `GET /api/v1/antibiotics/patient/{patient_id}` - Get antibiotics by patient

### Specimens

-    `GET /api/v1/specimens` - List all specimens with filtering
-    `GET /api/v1/specimens/{id}` - Get specific specimen
-    `GET /api/v1/specimens/stats` - Get specimen statistics
-    `GET /api/v1/specimens/patient/{patient_id}` - Get specimens by patient

### Upload

-    `POST /api/v1/upload/patients` - Upload patients CSV
-    `POST /api/v1/upload/antibiotics` - Upload antibiotics CSV
-    `POST /api/v1/upload/indications` - Upload indications CSV
-    `POST /api/v1/upload/optional-vars` - Upload optional variables CSV
-    `POST /api/v1/upload/specimens` - Upload specimens CSV

### Health Check

-    `GET /health` - API health status

## Quick Start

### Prerequisites

-    Go 1.21 or higher
-    PostgreSQL 12 or higher
-    Docker and Docker Compose (optional)

### Using Docker Compose (Recommended)

1. **Clone and navigate to the project directory:**

     ```bash
     cd "Point Prevalence Survey"
     ```

2. **Start the services:**

     ```bash
     docker-compose up -d
     ```

3. **Access the API:**
     - API: http://localhost:8080
     - Swagger Documentation: http://localhost:8080/swagger/index.html
     - Health Check: http://localhost:8080/health

### Manual Setup

1. **Install dependencies:**

     ```bash
     go mod download
     ```

2. **Set up PostgreSQL database:**

     ```sql
     CREATE DATABASE point_prevalence_survey;
     ```

3. **Configure environment variables:**

     ```bash
     cp config.env.example config.env
     # Edit config.env with your database credentials
     ```

4. **Run the application:**
     ```bash
     go run main.go
     ```

## Configuration

The application uses environment variables for configuration. Create a `config.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=point_prevalence_survey
DB_SSLMODE=disable
SERVER_PORT=8080
```

## CSV Upload

The API supports uploading CSV files for each data type. The CSV files should match the expected format:

-    **Patients**: Main survey data with patient information
-    **Antibiotics**: Antibiotic usage data with ATC codes and classifications
-    **Indications**: Treatment indications and diagnoses
-    **Optional Variables**: Additional treatment variables
-    **Specimens**: Microbiology specimen data

### Example CSV Upload

```bash
curl -X POST http://localhost:8080/api/v1/upload/patients \
  -F "file=@patients.csv"
```

## API Usage Examples

### Get all patients with filtering

```bash
curl "http://localhost:8080/api/v1/patients?region=Karamoja&page=1&limit=10"
```

### Get specific patient with all related data

```bash
curl "http://localhost:8080/api/v1/patients/uuid:86489018-001c-4a53-8c13-8635211e7d4e"
```

### Get patient statistics

```bash
curl "http://localhost:8080/api/v1/patients/stats"
```

### Get antibiotics by class

```bash
curl "http://localhost:8080/api/v1/antibiotics?class=Penicillins"
```

## Database Schema

The application uses GORM for database operations and automatically creates the following tables:

-    `patients` - Main patient records
-    `antibiotics` - Antibiotic usage data
-    `indications` - Treatment indications
-    `optional_vars` - Optional treatment variables
-    `specimens` - Microbiology specimens

All tables are properly linked with foreign key relationships.

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o main .
```

### Generating Swagger Documentation

```bash
swag init
```

## Docker Commands

### Build and run with Docker Compose

```bash
docker-compose up --build
```

### Stop services

```bash
docker-compose down
```

### View logs

```bash
docker-compose logs -f
```

### Access database

```bash
docker-compose exec postgres psql -U postgres -d point_prevalence_survey
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue in the repository.
