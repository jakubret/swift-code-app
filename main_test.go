package main

import (
    "database/sql"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB() *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:") // Baza danych w pamięci
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`CREATE TABLE swift_codes (
        swift_code TEXT PRIMARY KEY,
        bank_name TEXT,
        address TEXT,
        country_iso2 TEXT,
        country_name TEXT,
        is_headquarter BOOLEAN
    )`)
    if err != nil {
        panic(err)
    }

    // Dodanie przykładowych danych do testów
    _, err = db.Exec(`INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter) 
                      VALUES ('ABCDEFGHXXX', 'Bank of Earth', '123 Main St', 'US', 'United States', true)`)
    if err != nil {
        panic(err)
    }

    return db
}

func TestGetSwiftCodeNotFound(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    req, err := http.NewRequest("GET", "/v1/swift-codes/NONEXISTENT", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
    }
}

func TestAddSwiftCodeInvalidJSON(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    jsonStr := `{"address":"123 Main St","bankName":"Bank of Earth","countryISO2":"US","countryName":"United States","isHeadquarter":true,"swiftCode":"ABCDEFGHXXX"`

    req, err := http.NewRequest("POST", "/v1/swift-codes", strings.NewReader(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(addSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
}

func TestAddSwiftCodeMissingFields(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    jsonStr := `{"address":"123 Main St","bankName":"Bank of Earth","countryISO2":"US","isHeadquarter":true}`

    req, err := http.NewRequest("POST", "/v1/swift-codes", strings.NewReader(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(addSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
}

func TestDeleteSwiftCode(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    req, err := http.NewRequest("DELETE", "/v1/swift-codes/ABCDEFGHXXX", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(deleteSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestDeleteSwiftCodeNotFound(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    req, err := http.NewRequest("DELETE", "/v1/swift-codes/NONEXISTENT", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(deleteSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
    }
}

func TestGetSwiftCodesByCountry(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    req, err := http.NewRequest("GET", "/v1/swift-codes/country/US", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getSwiftCodesByCountry)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestGetSwiftCodesByCountryNotFound(t *testing.T) {
    db = setupTestDB()
    defer db.Close()

    req, err := http.NewRequest("GET", "/v1/swift-codes/country/XX", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getSwiftCodesByCountry)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK { // Zwraca pustą listę, więc status powinien być 200
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}
