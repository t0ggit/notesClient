package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notesClient/models/dto"
	"os"
	"strconv"
	"strings"
)

var (
	notesDebug        = false
	autoClearTerminal = true
)

func main() {
	for {
		var debugTag string
		if notesDebug {
			debugTag = "(DEBUG)"
		} else {
			debugTag = "_______"
		}
		var autoClearTag string
		if !autoClearTerminal {
			autoClearTag = "(NO AUTO-CLEAR)"
		} else {
			autoClearTag = "_______________"
		}
		fmt.Printf("___%s____[Notes App - CLI Client]_________%s______\n", autoClearTag, debugTag)
		fmt.Printf("[1]Create  [2]Get  [3]Update  [4]Delete  [5]GetAll  [9]Quit  [Clear]\n")
		fmt.Print(">> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()
		command = strings.ToLower(command)
		switch command {
		case "1", "create", "add":
			createNote()
		case "2", "get", "show":
			getNote()
		case "3", "update", "edit":
			updateNote()
		case "4", "delete", "remove":
			deleteNote()
		case "5", "get-all", "list":
			getAllNotes()
		case "9", "q", "quit", "exit":
			return
		case "clear":
			clearTerminal()
		case "debug":
			notesDebug = !notesDebug
			if autoClearTerminal {
				clearTerminal()
			}
		case "autoclear":
			autoClearTerminal = !autoClearTerminal
			if autoClearTerminal {
				clearTerminal()
			}
		case "help", "?":
			if autoClearTerminal {
				clearTerminal()
			}
			fmt.Println("[Notes App - CLI Client - Help]")
			fmt.Println("Commands:")
			fmt.Println("     [1 | create | add] - create a new note")
			fmt.Println("       [2 | get | show] - get details of a specific note")
			fmt.Println("    [3 | update | edit] - update an existing note")
			fmt.Println("  [4 | delete | remove] - delete a note")
			fmt.Println("   [5 | get-all | list] - get details of all notes")
			fmt.Println("  [9 | q | quit | exit] - exit the application")
			fmt.Println("                [clear] - clear the terminal")
			fmt.Println("                [debug] - toggle debug mode")
			fmt.Println("            [autoclear] - toggle auto clear terminal")
			fmt.Println("             [help | ?] - display this help message")
			fmt.Print("(help) enter to quit: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if autoClearTerminal {
				clearTerminal()
			}
		case "":
			if doYouWannaQuit() {
				return
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		default:
			if autoClearTerminal {
				clearTerminal()
			}
			fmt.Println("unknown command")
		}
	}
}

func createNote() {
	if autoClearTerminal {
		clearTerminal()
	}

	note := dto.NewNote()

	for note.Name == "" {
		fmt.Print("(create) Name: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				if autoClearTerminal {
					clearTerminal()
				}
				return
			}
		} else {
			note.Name = enteredText
		}
	}

	for note.LastName == "" {
		fmt.Print("(create) LastName: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				if autoClearTerminal {
					clearTerminal()
				}
				return
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		} else {
			note.LastName = enteredText
		}
	}

	for note.Content == "" {
		fmt.Print("(create) Note: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				if autoClearTerminal {
					clearTerminal()
				}
				return
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		} else {
			note.Content = enteredText
		}
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if notesDebug {
		fmt.Println("jsonData:", string(jsonData))
	}

	resp, err := http.Post("http://localhost:8080/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}

	HandleResponseBody(body)
}

func getNote() {
	note := dto.NewNote()

	fmt.Print("(get) ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	parsedID, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		fmt.Println("Error: cannot parse ID:", err)
		return
	}
	note.ID = parsedID

	if autoClearTerminal {
		clearTerminal()
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error: cannot marshal note:", err)
		return
	}

	if notesDebug {
		fmt.Println("jsonData:", string(jsonData))
	}

	resp, err := http.Post("http://localhost:8080/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}

	HandleResponseBody(body)
}

func updateNote() {
	note := dto.NewNote()

	fmt.Print("(update) ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	parsedID, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		fmt.Println("Error: cannot parse ID:", err)
		return
	}
	note.ID = parsedID

	// Проверяем существование записи с таким ID с помощью запроса на "/get"
	getJsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error: cannot marshal note:", err)
		return
	}

	if notesDebug {
		fmt.Println("getJsonData:", string(getJsonData))
	}

	getResp, err := http.Post("http://localhost:8080/get", "application/json", bytes.NewBuffer(getJsonData))
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	getRespBody, err := io.ReadAll(getResp.Body)
	getResp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}
	getRespDTO := dto.Response{}
	err = json.Unmarshal(getRespBody, &getRespDTO)
	if err != nil {
		fmt.Println("Error: cannot unmarshal response:", err)
		return
	}
	if getRespDTO.Result == "ERROR" {
		fmt.Println("Error:", getRespDTO.Error)
		return
	}
	oldNote := dto.Note{}
	err = json.Unmarshal(getRespDTO.Data, &oldNote)
	if err != nil {
		fmt.Println("Error: cannot unmarshal data:", err)
		return
	}

	if autoClearTerminal {
		clearTerminal()
	}

	for note.Name == "" {
		fmt.Printf("[upd#%d] Name (%s): ", note.ID, oldNote.Name)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				return
			}
			if doYouWannaNoUpdate(oldNote.Name) {
				note.Name = oldNote.Name
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		} else {
			note.Name = enteredText
		}
	}

	for note.LastName == "" {
		fmt.Printf("[upd#%d] LastName (%s): ", note.ID, oldNote.LastName)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				return
			}
			if doYouWannaNoUpdate(oldNote.LastName) {
				note.LastName = oldNote.LastName
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		} else {
			note.LastName = enteredText
		}
	}

	for note.Content == "" {
		fmt.Printf("[upd#%d] Note (%s): ", note.ID, oldNote.Content)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enteredText := scanner.Text()
		if enteredText == "" {
			if doYouWannaQuit() {
				return
			}
			if doYouWannaNoUpdate(oldNote.Content) {
				note.Content = oldNote.Content
			} else {
				if autoClearTerminal {
					clearTerminal()
				}
			}
		} else {
			note.Content = enteredText
		}
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error: cannot marshal note:", err)
		return
	}

	if notesDebug {
		fmt.Println("jsonData:", string(jsonData))
	}

	resp, err := http.Post("http://localhost:8080/update", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}

	HandleResponseBody(body)
}

func deleteNote() {
	note := dto.NewNote()

	fmt.Print("(delete) ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	parsedID, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		fmt.Println("Error: cannot parse ID:", err)
		return
	}
	note.ID = parsedID

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error: cannot marshal note:", err)
		return
	}

	if notesDebug {
		fmt.Println("jsonData:", string(jsonData))
	}

	resp, err := http.Post("http://localhost:8080/delete", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}

	HandleResponseBody(body)
}

func getAllNotes() {
	if autoClearTerminal {
		clearTerminal()
	}

	resp, err := http.Get("http://localhost:8080/get-all")
	if err != nil {
		fmt.Println("Error: request failed:", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error: cannot read body:", err)
		return
	}

	HandleResponseBody(body)
}

func HandleResponseBody(body []byte) {
	resp := dto.Response{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error: cannot unmarshal response:", err)
		return
	}

	if resp.Result == "ERROR" {
		fmt.Println("Error:", resp.Error)
		return
	}

	if resp.Result == "OK" && resp.Data != nil {
		data := []dto.Note{}
		err = json.Unmarshal(resp.Data, &data)
		if err != nil {
			data := dto.Note{}
			err = json.Unmarshal(resp.Data, &data)
			if err != nil {
				fmt.Println("Error: cannot unmarshal data:", err)
				return
			}
			PrintNote(data)
			fmt.Println()
			return
		}

		for _, note := range data {
			PrintNote(note)
		}
	}
	fmt.Println()
}

func PrintNote(note dto.Note) {
	// Вывод заметки с "шапкой" (номер, фамилия, имя) и контентом
	if note.ID > 0 && note.Name != "" && note.LastName != "" && note.Content != "" {
		fmt.Printf("\n┌┌──── Note #%d (by %s, %s)", note.ID, note.LastName, note.Name)
		fmt.Printf("\n└└ %s", note.Content)
	}

	// Вывод только номера заметки (после создания)
	if note.ID > 0 && note.Name == "" && note.LastName == "" && note.Content == "" {
		fmt.Printf("\n├├─── New Note is #%d\n", note.ID)
	}
}

func clearTerminal() {
	fmt.Print("\033[2J")
}

func doYouWannaQuit() bool {
	fmt.Print("└───── Do you wanna quit? (y): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	enteredAnswer := strings.ToLower(scanner.Text())
	if enteredAnswer == "y" || enteredAnswer == "yes" {
		return true
	}
	return false
}

func doYouWannaNoUpdate(oldValue string) bool {
	fmt.Printf("      ┌───(old)──→ \"%s\" ", oldValue)
	fmt.Print("\n      └─ Do you wanna leave it without update? (y): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	enteredAnswer := strings.ToLower(scanner.Text())
	if enteredAnswer == "y" || enteredAnswer == "yes" {
		return true
	}
	return false
}
