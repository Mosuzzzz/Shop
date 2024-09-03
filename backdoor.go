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
)

type Product struct {
	Product_name string    `json:"product_name"`
	Price		float64    `json:"price"`
}

var Productlist map[string]*Product
func main() {
	Productlist = ReadFile()
	for{
		clearScreen()
		Welcome()
		fmt.Println("WHAT YOU WANT TO DO WITH OUR SHOP?")
		fmt.Println("1.Add Product.")
		fmt.Println("2.Delete Product.")
		fmt.Println("3.List Product.")
		fmt.Println("4.Exit")
		choice := 0
		fmt.Scan(&choice)
		switch choice {
		case 1:
			clearScreen()
			AddProduct()
		case 2:
			clearScreen()
			DeleteProduct()
		case 3:
			clearScreen()
			ListProduct()
		case 4:
			clearScreen()
			fmt.Println("Exit...")
			return
		}
		fmt.Println("Please Enter to continue")
		fmt.Scanln()
		fmt.Println("Continuing with the program...")
		
		
	}
}


func Welcome()  {
	fmt.Print(`
		Welcome
`)
}

func AddProduct()  {
	var id string
	var InputProduct string
	var price float64
	fmt.Print("Enter Product ID : ")
	fmt.Scanln(&id)
	fmt.Print("Enter Product name : ")
	fmt.Scanln(&InputProduct)
	fmt.Print("Enter Price of Product : ")
	fmt.Scanln(&price)
	if id == ""||InputProduct == ""|| price == 0.0 {
		fmt.Println("Input anything plaese")
		return
	}
	_, exists := Productlist[id]
	if exists {
		fmt.Println("Already have this Product ID")
		return
	}

	Productlist[id] = &Product{Product_name: InputProduct,Price: price,}
	savedata()

}

func savedata() {
	data, err := json.Marshal(Productlist)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("save.json", data, fs.FileMode(0644))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK HURBLE")
}


func ListProduct()  {
	if len(Productlist) == 0 {
		fmt.Println("Emtry Peoduct")
	}
	KEY := make([]string, 0 ,len(Productlist))
	for k := range Productlist{
		KEY = append(KEY,k)
	}
	sort.Strings(KEY)
	for count,v := range KEY {
		Productlists := Productlist[v]
		fmt.Printf("%v. ID : %v Product Name : %v Price : %v$\n",count+1,v,Productlists.Product_name,Productlists.Price)
	}
}

func DeleteProduct()  {
	ListProduct()
	if len(Productlist) == 0 {
		return
	}
	id := ""
	fmt.Println("WHAT PRODUCT YOU WANT TO DELETE?")
	fmt.Print("Enter ID : ")
	fmt.Scanln(&id)
	_, exists := Productlist[id]
	if !exists {
		fmt.Println("Ain't got this Product.")
		return
	}
	for v := range Productlist {
		if v == id {
			delete(Productlist,v)
		}
	}

	savedata()
	
}

func ReadFile() map[string]*Product {
	Products,err := ioutil.ReadFile("save.json")
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