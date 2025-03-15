package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

type SwiftCode struct {
    Address        string `json:"address"`
    BankName       string `json:"bankName"`
    CountryISO2    string `json:"countryISO2"`
    CountryName    string `json:"countryName"`
    IsHeadquarter  bool   `json:"isHeadquarter"`
    SwiftCode      string `json:"swiftCode" gorm:"primary_key"`
}

const API_URL = "http://localhost:8080/v1/swift-codes"

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("SWIFT Code Manager")
    myWindow.Resize(fyne.NewSize(600, 400))

    swiftCodeEntry := widget.NewEntry()
    resultLabel := widget.NewLabel("")
    countryCodeEntry := widget.NewEntry()

    var swiftCodes []SwiftCode
    resultList := widget.NewList(
        func() int {
            return len(swiftCodes)
        },
        func() fyne.CanvasObject {
            return container.NewVBox(
                widget.NewLabel("Bank:"),
                widget.NewLabel("Adres:"),
            )
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            box := o.(*fyne.Container)
            labels := box.Objects
            labels[0].(*widget.Label).SetText("Bank: " + swiftCodes[i].BankName)
            labels[1].(*widget.Label).SetText("Adres: " + swiftCodes[i].Address)
        },
    )

    searchButton := widget.NewButton("Wyszukaj kod SWIFT", func() {
        swiftCode := swiftCodeEntry.Text
        if swiftCode == "" {
            resultLabel.SetText("Wprowadź kod SWIFT")
            return
        }

        resp, err := http.Get(fmt.Sprintf("%s/%s", API_URL, swiftCode))
        if err != nil {
            resultLabel.SetText(fmt.Sprintf("Błąd: %v", err))
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            var code SwiftCode
            if err := json.NewDecoder(resp.Body).Decode(&code); err != nil {
                resultLabel.SetText(fmt.Sprintf("Błąd dekodowania: %v", err))
                return
            }

            resultLabel.SetText(fmt.Sprintf(
                "Kod SWIFT: %s\nBank: %s\nAdres: %s\nKraj: %s (%s)\nGłówna siedziba: %v",
                code.SwiftCode, code.BankName, code.Address, code.CountryName, code.CountryISO2, code.IsHeadquarter,
            ))
        } else {
            resultLabel.SetText("Nie znaleziono kodu SWIFT")
        }
    })

    searchCountryButton := widget.NewButton("Wyszukaj kody SWIFT dla kraju", func() {
        countryCode := countryCodeEntry.Text
        if countryCode == "" {
            resultLabel.SetText("Wprowadź kod kraju (ISO2)")
            return
        }

        resp, err := http.Get(fmt.Sprintf("%s/country/%s", API_URL, countryCode))
        if err != nil {
            resultLabel.SetText(fmt.Sprintf("Błąd: %v", err))
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            var codes []SwiftCode
            if err := json.NewDecoder(resp.Body).Decode(&codes); err != nil {
                resultLabel.SetText(fmt.Sprintf("Błąd dekodowania: %v", err))
                return
            }

            resultWindow := myApp.NewWindow("Wyniki wyszukiwania")
            resultWindow.Resize(fyne.NewSize(800, 600))

            resultList := widget.NewList(
                func() int {
                    return len(codes)
                },
                func() fyne.CanvasObject {
                    return container.NewVBox(
                        widget.NewLabel("Bank:"),
                        widget.NewLabel("Adres:"),
                        widget.NewLabel("Kraj:"),
                        widget.NewLabel("Główna siedziba:"),
                    )
                },
                func(i widget.ListItemID, o fyne.CanvasObject) {
                    box := o.(*fyne.Container)
                    labels := box.Objects
                    labels[0].(*widget.Label).SetText("Bank: " + codes[i].BankName)
                    labels[1].(*widget.Label).SetText("Adres: " + codes[i].Address)
                    labels[2].(*widget.Label).SetText("Kraj: " + codes[i].CountryName)
                    labels[3].(*widget.Label).SetText("Główna siedziba: " + fmt.Sprintf("%v", codes[i].IsHeadquarter))
                },
            )

            resultWindow.SetContent(container.NewBorder(
                widget.NewLabel(fmt.Sprintf("Wyniki wyszukiwania dla kraju: %s", countryCode)),
                nil, nil, nil,
                container.NewVScroll(resultList),
            ))

            resultWindow.Show()
        } else {
            resultLabel.SetText("Nie znaleziono kodów SWIFT dla tego kraju")
        }
    })

    addSwiftCodeButton := widget.NewButton("Dodaj nowy kod SWIFT", func() {
        swiftCodeEntry := widget.NewEntry()
        bankNameEntry := widget.NewEntry()
        addressEntry := widget.NewEntry()
        countryISO2Entry := widget.NewEntry()
        countryNameEntry := widget.NewEntry()
        isHeadquarterCheck := widget.NewCheck("Główna siedziba", func(bool) {})

        form := &widget.Form{
            Items: []*widget.FormItem{
                {Text: "Kod SWIFT", Widget: swiftCodeEntry},
                {Text: "Nazwa banku", Widget: bankNameEntry},
                {Text: "Adres", Widget: addressEntry},
                {Text: "Kod kraju (ISO2)", Widget: countryISO2Entry},
                {Text: "Nazwa kraju", Widget: countryNameEntry},
                {Text: "Główna siedziba", Widget: isHeadquarterCheck},
            },
            OnSubmit: func() {
                code := SwiftCode{
                    SwiftCode:     swiftCodeEntry.Text,
                    BankName:      bankNameEntry.Text,
                    Address:       addressEntry.Text,
                    CountryISO2:   countryISO2Entry.Text,
                    CountryName:   countryNameEntry.Text,
                    IsHeadquarter: isHeadquarterCheck.Checked,
                }

                jsonData, err := json.Marshal(code)
                if err != nil {
                    dialog.ShowError(fmt.Errorf("Błąd podczas tworzenia danych JSON"), myWindow)
                    return
                }

                resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(jsonData))
                if err != nil {
                    dialog.ShowError(fmt.Errorf("Błąd podczas wysyłania żądania: %v", err), myWindow)
                    return
                }
                defer resp.Body.Close()

                if resp.StatusCode == http.StatusOK {
                    dialog.ShowInformation("Sukces", "Kod SWIFT został dodany", myWindow)
                } else {
                    dialog.ShowError(fmt.Errorf("Nie udało się dodać kodu SWIFT"), myWindow)
                }
            },
        }

        addWindow := myApp.NewWindow("Dodaj nowy kod SWIFT")
        addWindow.SetContent(form)
        addWindow.Resize(fyne.NewSize(400, 300))
        addWindow.Show()
    })

    deleteSwiftCodeButton := widget.NewButton("Usuń kod SWIFT", func() {
        swiftCodeEntry := widget.NewEntry()

        form := &widget.Form{
            Items: []*widget.FormItem{
                {Text: "Kod SWIFT", Widget: swiftCodeEntry},
            },
            OnSubmit: func() {
                swiftCode := swiftCodeEntry.Text
                if swiftCode == "" {
                    dialog.ShowError(fmt.Errorf("Wprowadź kod SWIFT"), myWindow)
                    return
                }

                req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", API_URL, swiftCode), nil)
                if err != nil {
                    dialog.ShowError(fmt.Errorf("Błąd podczas tworzenia żądania: %v", err), myWindow)
                    return
                }

                client := &http.Client{}
                resp, err := client.Do(req)
                if err != nil {
                    dialog.ShowError(fmt.Errorf("Błąd podczas wysyłania żądania: %v", err), myWindow)
                    return
                }
                defer resp.Body.Close()

                if resp.StatusCode == http.StatusOK {
                    dialog.ShowInformation("Sukces", "Kod SWIFT został usunięty", myWindow)
                } else if resp.StatusCode == http.StatusNotFound {
                    dialog.ShowError(fmt.Errorf("Nie znaleziono kodu SWIFT"), myWindow)
                } else {
                    dialog.ShowError(fmt.Errorf("Nie udało się usunąć kodu SWIFT"), myWindow)
                }
            },
        }

        deleteWindow := myApp.NewWindow("Usuń kod SWIFT")
        deleteWindow.SetContent(form)
        deleteWindow.Resize(fyne.NewSize(300, 150))
        deleteWindow.Show()
    })

    content := container.NewVBox(
        widget.NewLabel("Wyszukaj kod SWIFT:"),
        swiftCodeEntry,
        searchButton,
        resultLabel,
        widget.NewLabel("Wyszukaj kody SWIFT dla kraju:"),
        countryCodeEntry,
        searchCountryButton,
        container.NewVScroll(resultList),
        addSwiftCodeButton,
        deleteSwiftCodeButton,
    )

    myWindow.SetContent(content)
    myWindow.ShowAndRun()
}
