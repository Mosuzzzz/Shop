package main



import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
)

type Product struct {
	Product_name string  `json:"product_name"`
	Price        float64 `json:"price"`
	Count        int     `json:"count"`
}


var Productlist map[string]*Product
var Chartlist map[string]*Product
var accountBalance float64

func main() {
	accountBalance = Bookbank()
	Productlist = ReadFile()
	Chartlist = ReadFileChart()
	cleanUpCart()
	for {
		clearScreen()
		Welcome()
		fmt.Println("Do you want to buy something?")
		fmt.Println("1.List Product in chart.")
		fmt.Println("2.GO to buy something.")
		fmt.Println("3.Delete Product in chart.")
		fmt.Println("4.Deposit")
		fmt.Println("5.Withdraw")
		fmt.Println("6.CheckMONEY")
		fmt.Println("7.Pay the bill.")
		fmt.Println("8.Exit")
		choice := 0
		fmt.Scan(&choice)
		switch choice {
		case 1:
			clearScreen()
			ListinChart()
		case 2:
			clearScreen()
			Shopping()
		case 3:
			clearScreen()
			DeleteProduct()
		case 4 :
			clearScreen()
			deposit(&accountBalance)
		case 5 :
			clearScreen()
			withdraw(&accountBalance)
		case 6:
			clearScreen()
			Displaybaland(accountBalance)
		case 7:
			clearScreen()
			pay()
		case 8:
			clearScreen()
			fmt.Println("Exit...")
			return
		}
		fmt.Println("Please Enter to continue")
		fmt.Scanln()
		fmt.Println("Continuing with the program...")

	}
}

func pay() {
	sum := 0
	ListinChart()
	for _,v := range Chartlist {
		sum += int(v.Price)
	}
	fmt.Println("Totol is ",sum)
	if accountBalance < float64(sum) {
		fmt.Println("You have enough Money.")
		return
	}
	accountBalance -= float64(sum)
	for v := range Chartlist {
			delete(Chartlist, v)
		savedata()
	}
	save(&accountBalance)
}

func Shopping() {
	id := ""
	count := 0
	ListProduct()
	fmt.Println("Enter Product ID you want to buy : ")
	fmt.Scanln(&id)
	product, exists := Productlist[id]
	if !exists {
		fmt.Println("Product not found!")
		return
	}
	fmt.Println("How much you want to buy?")
	fmt.Scanln(&count)
	if count == 0 {
		return
	}
	Chartlist[id] = &Product{
		Product_name: product.Product_name,
		Price:        product.Price*float64(count),
		Count:        count,
	}
	savedata()

}

func Welcome() {
	fmt.Print(`
		Welcome
`)
}


func Bookbank()(accountBalance float64){
	Bytes, err := os.ReadFile("Money.data")
	if err != nil {
		return 500
	}
	BalanceText := string(Bytes)
	accountBalance, _ = strconv.ParseFloat(BalanceText, 64)
	return 
}

func Displaybaland(accountBalance float64) {
	fmt.Println(" > Your balance is ", accountBalance, " $")
}

func withdraw(accountBalance *float64) {
	amount := getInput(" Withdraw amount", 99.99999999999999999)
	*accountBalance -= amount
	if *accountBalance < amount {
		*accountBalance = 0
	}
	Displaybaland(*accountBalance)
	save(accountBalance)
}


func deposit(accountBalance *float64) {
	depositAmount := getInput(" Deposit amount", 0.00000001)
	*accountBalance += depositAmount
	Displaybaland(*accountBalance)
	save(accountBalance)
}

func getInput(prompt string, limit float64) float64 {
	fmt.Printf(" %v: ", prompt)
	amount := 0.0
	fmt.Scan(&amount)
	if amount < limit {
		fmt.Println(" Invalid value")
		return 0
	}
	return amount
}

func save(accountBalance *float64) {
	data := fmt.Sprint(*accountBalance)
	os.WriteFile("Money.data",[]byte(data),0664)
}

func savedata() {
	data, err := json.Marshal(Chartlist)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("chart.json", data, fs.FileMode(0644))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sucessful!!")
}

func returnEmtry() bool  {
	if len(Chartlist) == 0 {
		fmt.Println("Emtry Peoduct")
		return false
	}
	return true
}

func ListinChart() {
	if !returnEmtry(){
		return
	}
	KEY := make([]string, 0, len(Chartlist))
	for k := range Chartlist {
		KEY = append(KEY, k)
	}
	sort.Strings(KEY)
	for count, v := range KEY {
		Chartlists := Chartlist[v]
		fmt.Printf("%v. ID : %v Product Name : %v Price : %v$ Count : %v\n", count+1, v, Chartlists.Product_name, Chartlists.Price,Chartlists.Count)
	}
}


func ListProduct() {
	if len(Productlist) == 0 {
		fmt.Println("Emtry Peoduct")
	}
	KEY := make([]string, 0, len(Productlist))
	for k := range Productlist {
		KEY = append(KEY, k)
	}
	sort.Strings(KEY)
	for count, v := range KEY {
		Productlists := Productlist[v]
		fmt.Printf("%v. ID : %v Product Name : %v Price : %v$\n", count+1, v, Productlists.Product_name, Productlists.Price,)
	}
}

func DeleteProduct() {
	ListinChart()
	id := ""
	fmt.Println("WHAT PRODUCT YOU WANT TO DELETE IN YOUR CHART?")
	fmt.Print("Enter ID : ")
	fmt.Scanln(&id)
	
	count, exists := Chartlist[id]
	if !exists {
		fmt.Println("Ain't got this Product.")
		return
	}

	product := Productlist[id]
	for {
		fmt.Println("Select")
		fmt.Println("1.Delete some")
		fmt.Println("2.Delete ALL")
		choice := 0
		fmt.Scan(&choice)
		switch choice {
		case 1:
			sum := 0
			fmt.Println("How much you want to delete : ")
			fmt.Scanln(&sum)
			count.Count -= sum
			if count.Count == 0 {
				for v := range Chartlist {
					if v == id {
						delete(Chartlist, v)
					}
					savedata()
					return
				}
			}
			if sum == 0  {
				return
			}
		sum = int(product.Price) * count.Count
		count.Price = float64(sum)
		savedata()
		return
		case 2:
			for v := range Chartlist {
				if v == id {
					delete(Chartlist, v)
				}
				savedata()
				return
			}
		}
		
	}
}

func ReadFile() map[string]*Product {
	Products, err := ioutil.ReadFile("../save.json")
	if err != nil {
		fmt.Println(err)
		return make(map[string]*Product)
	}
	ProductData := map[string]*Product{}
	err = json.Unmarshal(Products, &ProductData)
	if err != nil {
		fmt.Println(err)
		return make(map[string]*Product)
	}

	return ProductData
}

func ReadFileChart() map[string]*Product {
	Products, err := ioutil.ReadFile("chart.json")
	if err != nil {
		fmt.Println(err)
		return make(map[string]*Product)
	}
	ProductData := map[string]*Product{}
	err = json.Unmarshal(Products, &ProductData)
	if err != nil {
		fmt.Println(err)
		return make(map[string]*Product)
	}

	return ProductData
}


func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Print("\033[2J")
		fmt.Print("\033[H")
	}
}

func cleanUpCart() {
	
	for id := range Chartlist {
		if _, exists := Productlist[id]; !exists {
			delete(Chartlist, id)
		}
	}
savedata()
}