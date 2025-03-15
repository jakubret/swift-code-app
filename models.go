package main2

type SwiftCode struct {
    Address        string `json:"address"`
    BankName       string `json:"bankName"`
    CountryISO2    string `json:"countryISO2"`
    CountryName    string `json:"countryName"`
    IsHeadquarter  bool   `json:"isHeadquarter"`
    SwiftCode      string `json:"swiftCode" gorm:"primary_key"`
}

type Branch struct {
    Address        string `json:"address"`
    BankName       string `json:"bankName"`
    CountryISO2    string `json:"countryISO2"`
    IsHeadquarter  bool   `json:"isHeadquarter"`
    SwiftCode      string `json:"swiftCode"`
}