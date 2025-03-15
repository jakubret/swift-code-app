package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestGetSwiftCodeNotFound(t *testing.T) {
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

func TestGetSwiftCodeInvalidFormat(t *testing.T) {
    req, err := http.NewRequest("GET", "/v1/swift-codes/INVALID", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getSwiftCode)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }
}

func TestAddSwiftCodeInvalidJSON(t *testing.T) {
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
    req, err := http.NewRequest("GET", "/v1/swift-codes/country/XX", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getSwiftCodesByCountry)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
    }
}