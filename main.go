package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//d d
const (
	LocalhostPort = ":8080"
)

//PageData s
type PageData struct {
	TypeSome string
	Kolor    string
	// Justifi  string
	// Baner string
	// CheckBoxe string
}

//HomePage fzeiuh
func HomePage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("Internal Server Error", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Internal Server Error", err)
	}
}

//ResultPage zfeiub
func ResultPage(w http.ResponseWriter, r *http.Request) {

	// Initilazing each values from HTML file to Go
	StoASCII := r.FormValue("StringToAscii")
	Banner := r.FormValue("PoliceStyle")
	Color := r.FormValue("Colorize")
	Justify := r.FormValue("Justify")
	CheckColor := r.FormValue("CheckColor")

	data := PageData{TypeSome: StoASCII, Kolor: Color} // Send back values to HTML to display it on the page
	tmpl, err := template.ParseFiles("index.html")     // Parsing HTML file

	if err != nil {
		fmt.Println("Internal Server Error", err)
	} else {
		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("Internal Server Error", err)
		}
		Ascii_Art(w, StoASCII, Banner, Color, Justify, CheckColor) // Calling the function to print with the different choices
	}
}

func main() {
	fmt.Println("Please connect to", "\u001b[31m localhost", LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.Handle("/site/", http.StripPrefix("/site/", http.FileServer(http.Dir("site"))))       // Join Assets Directory to the server
	http.HandleFunc("/", HomePage)                                                             // set router
	http.HandleFunc("/ascii-art", ResultPage)
	err := http.ListenAndServe(LocalhostPort, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Ascii_Art zfeuib
func Ascii_Art(w http.ResponseWriter, args string, Banner string, colorrr string, Justify string, CheckColor string) {

	index := 1
	var file *os.File
	var arr []byte
	var word []string

	// replacing newlines from Enter key in HTTP to \n to pass the error section
	args = strings.Replace(args, string('\r'), "", -1)
	args = strings.Replace(args, string('\n'), "\\n", -1)

	// Write an error in case of invalid argument
	for _, char := range args {
		if char < ' ' && args != "" || char > '~' && args != "" {
			fmt.Println("Bad Request")
			return
		}
	}
	if args == "" || Banner == "" || colorrr == "" || Justify == "" { // if an argument is missing
		if args == "" && (Banner != "" || colorrr != "" || Justify != "") { // if only "args" is missing (if several are missing it might be returning to main page)
			fmt.Println("Not Found")
		}
		return
	}

	// Creating a HTML div to display the result
	fmt.Fprint(w, "<div class=\"Result\" style=\"color:"+colorrr+";")

	// Changing the background color when the box is checked
	if CheckColor == "CheckColor" {
		fmt.Fprint(w, "background-color: #CBCED2;")
	}
	fmt.Fprint(w, "\">")

	// Applying CSS to justify the result
	fmt.Fprint(w, "<div style=\"")
	if Justify == "left" {
		fmt.Fprint(w, " text-align: left;\">")
	} else if Justify == "right" {
		fmt.Fprint(w, " text-align: right;\">")
	} else if Justify == "center" {
		fmt.Fprint(w, " text-align: center;\">")
	}

	//Printing the result to the website
	PoliceStyle := Language(Banner)
	word = strings.Split(args, "\\n")
	for id := 0; id < len(word); id++ {
		for index = 1; index <= 8; index++ {
			fmt.Fprintln(w, "<br>")
			arr = append(arr, byte('\\'))
			arr = append(arr, byte('n'))
			for _, char := range word[id] {
				pos := ((int(rune(char)-32))*9 + index)
				PoliceStyle[pos] = strings.Replace(PoliceStyle[pos], " ", "&nbsp;", -1) // replacing GO spaces to HTML spaces
				for _, val := range PoliceStyle[pos] {
					arr = append(arr, byte(val))
				}
				fmt.Fprint(w, PoliceStyle[pos])
			}
		}
		index = 1
	}

	// Create a file where the result will be display in case of downloading
	file, _ = os.Create("assets/Ascii.txt")

	// Replacing every \n to new lines
	for index := range arr {
		if arr[index] == 92 && arr[index+1] == 110 {
			arr[index] = 10
			arr[index+1] = 0
		}
	}

	// Changing spaces from HTML to GO and Printing the result into a file
	arr = bytes.ReplaceAll(arr, []byte("&nbsp;"), []byte(" "))
	file.Write(arr)

	// Seperating the Result section from the Download section
	fmt.Fprint(w, "</div>")
	fmt.Fprint(w, "<div class=\"BottomResult\">")
	fmt.Fprint(w, "<br><br><a class=\"download\" href=\"assets/Ascii.txt\" download><button>Download to .txt</button></a>")                       // Creating download button
	fmt.Fprint(w, "<br><br><p class=\"interro\" id=\"interro2\">?<span>Justify and Colors are not taken into account when downloaded</span></p>") // Creating Help question mark to Download button
	fmt.Fprint(w, "<a href=\"ascii-art\"><button>Restart</button></a>")                                                                           // Creating Help question mark to Download button
	fmt.Fprint(w, "</div>")
	fmt.Fprint(w, "</div>")

	// If everything worked
	fmt.Println("Ok(200)")
}

//function who take all the lines from a file
func line(file *os.File) []string {
	var string []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := scanner.Text()
		string = append(string, words)
	}
	return string
}

// Language is used to open the file that is in entrance
func Language(s string) []string {
	var banner []string
	if s == "standard" || s == "shadow" || s == "thinkertoy" {
		file, _ := os.Open(s + ".txt") //we open the file we want
		banner = line(file)            //turn a file into a slice of string
	} else {
		fmt.Println("Error")
		os.Exit(1)
	}
	return banner
}
