package main1

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "os"
    "strings"

    _ "github.com/mattn/go-sqlite3"
)


func main1() {
    db, err := sql.Open("sqlite3", "./swift_codes1.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS swift_codes (
        swift_code TEXT PRIMARY KEY,
        bank_name TEXT,
        address TEXT,
        country_iso2 TEXT,
        country_name TEXT,
        is_headquarter BOOLEAN
    )`)
    if err != nil {
        panic("failed to create table")
    }

    file, err := os.Open("data.csv")
    if err != nil {
        panic("failed to open CSV file")
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        panic("failed to read CSV file")
    }

    for i, record := range records {
        if i == 0 {
            continue
        }

        swiftCode := record[1]
        isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

        swiftCodeData := SwiftCode{
            SwiftCode:     swiftCode,
            BankName:      record[3], 
            Address:       record[4],  
            CountryISO2:   strings.ToUpper(record[0]),  
            CountryName:   strings.ToUpper(record[6]),  
            IsHeadquarter: isHeadquarter,
        }

        _, err := db.Exec("INSERT INTO swift_codes VALUES (?, ?, ?, ?, ?, ?)",
            swiftCodeData.SwiftCode, swiftCodeData.BankName, swiftCodeData.Address, swiftCodeData.CountryISO2, swiftCodeData.CountryName, swiftCodeData.IsHeadquarter)
        if err != nil {
            fmt.Printf("Error inserting record %d: %v\n", i, err)
        }
    }

    fmt.Println("Data import completed!")
}