package main

/** To Do List Code
 *  [구현해야 하는 기능]
 *  1. 메뉴 출력
 *
 *  [기본 기능]
 *  1. 스케줄 입력
 *  2. 스케줄 출력
 *  3. 스케줄 알림
 *  4. 스케줄 내보내기
 */

// Path: /workspace/chaney11
// cd /workspace/chaney11
// go run main.go

// source ~/.profile

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	// "time"
)

type (
	Menu struct {
		Title string
	}

	MenuList struct {
		MenuFormatter

		List []Menu
	}

	MenuFormatter interface {
		ToMenu(...string) *MenuList
	}

	Schedule struct {
		Name string
		// StartTime, EndTime time.Duration
	}

	ScheduleList struct {
		sList []Schedule
	}
)

func checkErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func (List *MenuList) ToFormat() string {
	var response []byte

	for i, v := range List.List {
		var data []byte = []byte(strconv.Itoa(i+1) + ". " + v.Title + "\n")
		for _, j := range data {
			response = append(response, j)
		}
	}

	return string(response)
}

func ToMenu(data ...string) *MenuList {
	response := &MenuList{
		List: make([]Menu, len(data)),
	}

	for i, v := range data {
		response.List[i] = Menu{
			Title: v,
		}
	}

	return response
}

func (List *ScheduleList) ScheString() string {
	var response []byte

	for i, v := range List.sList {
		var data []byte = []byte(strconv.Itoa(i+1) + ". " + v.Name + "\n")
		for _, j := range data {
			response = append(response, j)
		}
	}

	return string(response)
}

func (s *ScheduleList) addSche() {
	var line string

	fmt.Print("Type here: ")

	in := bufio.NewReader(os.Stdin)

	line, _ = in.ReadString('\n')

	in.Read([]byte{}) // 버퍼에 담긴 모든 문자를 데이터로 반환, but 저장되는 곳이 없어 데이터는 소멸됨.

	// fmt.Scanf(" %v", line)

	s.sList = append(s.sList, Schedule{
		Name: string(strings.TrimSpace(line)),
	})

	/* Debug Code */
	//fmt.Printf("Added Schedule Name: %s\n", line)

}

func (s *ScheduleList) deleteSche() {
	if s.sList == nil {
		fmt.Println("List is empty!")
	} else {
		fmt.Println("Which one you want to delete. (Type the number below here)")
		fmt.Println(s.ScheString())
		fmt.Print("Type here:")

		var num int
		fmt.Scan(&num)

		s.sList = append(s.sList[:num-1], s.sList[num:]...)
	}
}
func (s *ScheduleList) makeCSV() {
	if s.sList == nil {
		fmt.Println("List is empty!")
	} else {
		fmt.Println("Make CSV file!")
		// 파일 생성
		file, err := os.Create("./ToDoList.csv")
		if err != nil {
			panic(err)
		}

		// csv writer
		w := csv.NewWriter(file)

		header := []string{"num", "ToDoList"}
		writeE := w.Write(header)
		checkErr(writeE)
		for i, record := range s.sList {
			list := []string{strconv.Itoa(i + 1), record.Name}
			listE := w.Write(list)
			checkErr(listE)
		}
		w.Flush()
	}
}
func (s *ScheduleList) SendEmail() {
	if s.sList == nil {
		fmt.Println("List is empty!")
	} else {
		auth := smtp.PlainAuth("", "chaney2225@gmail.com", "lrybdiqxljdotquj", "smtp.gmail.com")
		from := "chaney2225@gmail.com"
		to := []string{"chaney11@naver.com"}
		headerSubject := "Subject: ToDoList\r\n"
		headerBlank := "\r\n"
		body := s.ScheString()
		msg := []byte(headerSubject + headerBlank + body)

		err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Your ToDoList is successfully sended")
		}

	}
}

func main() {
	var (
		MenuForm = ToMenu(
			"Add new schedule.",
			"Remove existing schedule.",
			"Make ToDo.csv file",
			"Send ToDoList to your email",
			"Exit.",
		)
		sche ScheduleList
	)

	for {
		var num int

		if len(sche.sList) > 0 {
			fmt.Println(sche.ScheString())
		}

		fmt.Println("[To Do List]")
		fmt.Println(MenuForm.ToFormat())

		fmt.Scan(&num)

		switch num {
		case 1:
			sche.addSche()
		case 2:
			sche.deleteSche()
		case 3:
			sche.makeCSV()
		case 4:
			sche.SendEmail()
		case 5:
			os.Exit(0)
		}
	}

}
