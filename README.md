# SWIFT Code Manager

An application for managing bank SWIFT codes. It allows searching, adding, deleting, and browsing SWIFT codes for different countries.

## Features
- **Search for a SWIFT code**: Enter a SWIFT code to get information about the bank.
- **Search for SWIFT codes by country**: Enter a country code (ISO2) to get a list of SWIFT codes for that country.
- **Add a new SWIFT code**: Add a new SWIFT code to the database.
- **Delete a SWIFT code**: Remove an existing SWIFT code from the database.

## Technologies
- **Go**: The programming language used to build the application.
- **Fyne**: A library for creating the graphical user interface (GUI).
- **SQLite**: The database used to store SWIFT codes.

## Requirements
- **Go** (version 1.20 or higher)
- **SQLite** (optional, if you want to manage the database manually)

## Installation and Running

### Using Docker

Build the Docker image:
   ```bash
   docker build -t swift-code-app .

```

Run the Docker container (e.g., on port 8080):

```
docker run -it --rm -p 8080:8080 swift-code-app
```
Clone the repository:
```bash
git clone https://github.com/twoj-uzytkownik/swift-code-app.git
cd swift-code-app

```
Initialize the Go module:
```bash
go mod init swift-code-app

```

Download dependencies:

```bash
go mod tidy

```
<img width="599" alt="image" src="https://github.com/user-attachments/assets/0c6a0a22-4961-484c-865f-4d2321e1cabb" />
<img width="798" alt="image" src="https://github.com/user-attachments/assets/2f39c837-0c55-4809-bfd5-fd55a723588f" />


