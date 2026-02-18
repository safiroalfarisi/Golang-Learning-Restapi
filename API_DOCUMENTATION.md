# Employee Management API

## Endpoints

### Base URL: `/api/employees`
*Semua endpoints memerlukan JWT token di header Authorization: `Bearer <token>`*

## CRUD Operations

### 1. Get All Employees
- **GET** `/api/employees`
- **Response**: Array of employee objects

### 2. Get Employee by ID
- **GET** `/api/employees/:id`
- **Response**: Single employee object

### 3. Create Employee
- **POST** `/api/employees`
- **Body**:
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "position": "Developer",
    "role": "staff",
    "password": "password123"
}
```

### 4. Update Employee
- **PUT** `/api/employees/:id`
- **Body**: (fields yang ingin diupdate saja)
```json
{
    "name": "John Doe Updated",
    "position": "Senior Developer",
    "role": "admin"
}
```

### 5. Delete Employee
- **DELETE** `/api/employees/:id`
- **Response**: Confirmation message

## Import/Export Operations

### 6. Import Employees from CSV
- **POST** `/api/employees/import`
- **Content-Type**: `multipart/form-data`
- **Body**: Form dengan field `file` (CSV file)
- **CSV Format**:
  - Header: `Name,Email,Position,Role`
  - Default password untuk semua employee: `password123`
- **Example**:
```csv
Name,Email,Position,Role
John Doe,john@example.com,Developer,staff
Jane Smith,jane@example.com,Manager,admin
```

### 7. Export Employees to CSV
- **GET** `/api/employees/export`
- **Response**: CSV file download
- **Filename**: `employees.csv`

## Example Usage

### 1. Update Employee (curl)
```bash
curl -X PUT http://localhost:3000/api/employees/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "position": "Senior Developer"
  }'
```

### 2. Import CSV (curl)
```bash
curl -X POST http://localhost:3000/api/employees/import \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@employees.csv"
```

### 3. Export CSV (curl)
```bash
curl -X GET http://localhost:3000/api/employees/export \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o employees.csv
```

## Notes

- Semua password akan di-hash menggunakan bcrypt
- Import CSV menggunakan password default: `password123`
- Delete operation menggunakan soft delete (data tidak benar-benar terhapus)
- Template CSV tersedia di file `employees_template.csv`