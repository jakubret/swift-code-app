package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB


func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./swift_codes1.db")
    if err != nil {
        log.Fatal(err)
    }
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS swift_codes (
        swift_code TEXT PRIMARY KEY,
        bank_name TEXT,
        address TEXT,
        country_iso2 TEXT,
        country_name TEXT,
        is_headquarter BOOLEAN
    )`)
    if err != nil {
        log.Fatal(err)
    }
}

func getSwiftCode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    swiftCode := vars["swiftCode"]

    var code SwiftCode
    err := db.QueryRow("SELECT * FROM swift_codes WHERE swift_code = ?", swiftCode).
        Scan(&code.SwiftCode, &code.BankName, &code.Address, &code.CountryISO2, &code.CountryName, &code.IsHeadquarter)
    if err != nil {
        http.Error(w, "SWIFT code not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(code)
}

func getSwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    countryISO2 := vars["countryISO2code"]

    rows, err := db.Query("SELECT * FROM swift_codes WHERE country_iso2 = ?", countryISO2)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var codes []SwiftCode
    for rows.Next() {
        var code SwiftCode
        if err := rows.Scan(&code.SwiftCode, &code.BankName, &code.Address, &code.CountryISO2, &code.CountryName, &code.IsHeadquarter); err != nil {
            http.Error(w, "Error reading data", http.StatusInternalServerError)
            return
        }
        codes = append(codes, code)
    }
    json.NewEncoder(w).Encode(codes)
}

func addSwiftCode(w http.ResponseWriter, r *http.Request) {
    var code SwiftCode
    if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    _, err := db.Exec("INSERT INTO swift_codes VALUES (?, ?, ?, ?, ?, ?)",
        code.SwiftCode, code.BankName, code.Address, code.CountryISO2, code.CountryName, code.IsHeadquarter)
    if err != nil {
        http.Error(w, "Failed to insert data", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "SWIFT code added successfully"})
}

func deleteSwiftCode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    swiftCode := vars["swiftCode"]

    _, err := db.Exec("DELETE FROM swift_codes WHERE swift_code = ?", swiftCode)
    if err != nil {
        http.Error(w, "Failed to delete", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "SWIFT code deleted successfully"})
}

func main() {
    initDB()
    defer db.Close()

    r := mux.NewRouter()
    r.HandleFunc("/v1/swift-codes/{swiftCode}", getSwiftCode).Methods("GET")
    r.HandleFunc("/v1/swift-codes/country/{countryISO2code}", getSwiftCodesByCountry).Methods("GET")
    r.HandleFunc("/v1/swift-codes", addSwiftCode).Methods("POST")
    r.HandleFunc("/v1/swift-codes/{swiftCode}", deleteSwiftCode).Methods("DELETE")

    log.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}