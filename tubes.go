package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/aquasecurity/table"
	"github.com/nsf/termbox-go"
)

type Jurusan string

type Mahasiswa struct {
	ID       string
	Nama     string
	Jurusan  string
	NilaiTes float64
	Status   string
}

type User struct {
	Name string
	Pass string
	Role string // "admin", "student"
}

func dummy() ([]Jurusan, []Mahasiswa) {
	jurusans := []Jurusan{
		"Teknik Informatika",
		"Sistem Informasi",
		"Manajemen",
		"Akuntansi",
		"Teknik Sipil",
		"Arsitektur",
		"Psikologi",
		"Biologi",
		"Matematika",
		"Fisika",
	}

	mahasiswas := []Mahasiswa{
		{ID: "1", Nama: "Alfa", Jurusan: "Teknik Informatika", NilaiTes: 80.5, Status: "✅"},
		{ID: "2", Nama: "Beta", Jurusan: "Sistem Informasi", NilaiTes: 65.4, Status: "❌"},
		{ID: "3", Nama: "Charlie", Jurusan: "Manajemen", NilaiTes: 90.2, Status: "✅"},
		{ID: "4", Nama: "Delta", Jurusan: "Akuntansi", NilaiTes: 70.3, Status: "❌"},
		{ID: "5", Nama: "Echo", Jurusan: "Teknik Sipil", NilaiTes: 85.6, Status: "✅"},
		{ID: "6", Nama: "Foxtrot", Jurusan: "Arsitektur", NilaiTes: 72.1, Status: "❌"},
		{ID: "7", Nama: "Golf", Jurusan: "Psikologi", NilaiTes: 77.8, Status: "✅"},
		{ID: "8", Nama: "Hotel", Jurusan: "Biologi", NilaiTes: 89.3, Status: "✅"},
		{ID: "9", Nama: "India", Jurusan: "Matematika", NilaiTes: 60.5, Status: "❌"},
		{ID: "10", Nama: "Juliet", Jurusan: "Fisika", NilaiTes: 95.0, Status: "✅"},
	}
	return jurusans, mahasiswas
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func highlightText(text string) string {
	return fmt.Sprintf("\033[7;34m%s\033[0m", text)
}

func tambahJurusan(jurusan []Jurusan, nama string) []Jurusan {
	jurusan = append(jurusan, Jurusan(nama))
	return jurusan
}

func hapusJurusan(jurusan []Jurusan, nama string) []Jurusan {
	for i := 0; i < len(jurusan); i++ {
		if jurusan[i] == Jurusan(nama) {
			jurusan = append(jurusan[:i], jurusan[i+1:]...)
			break
		}
	}
	return jurusan
}

func editJurusan(jurusan []Jurusan, mahasiswa []Mahasiswa, oldName, newName string) ([]Jurusan, []Mahasiswa) {
	for i := 0; i < len(jurusan); i++ {
		if jurusan[i] == Jurusan(oldName) {
			jurusan[i] = Jurusan(newName)
			break
		}
	}

	for i := 0; i < len(mahasiswa); i++ {
		if mahasiswa[i].Jurusan == oldName {
			mahasiswa[i].Jurusan = newName
		}
	}

	return jurusan, mahasiswa
}

func viewJurusan(jurusan []Jurusan) {
	t := table.New(os.Stdout)
	t.SetHeaders("No", "Nama Jurusan")
	for i, j := range jurusan {
		t.AddRow(strconv.Itoa(i+1), string(j))
	}

	t.Render()
}

func tambahMahasiswa(mhs []Mahasiswa, id string, nama, jurusan string, nilaiTes float64) []Mahasiswa {
	status := "ditolak"
	if nilaiTes >= 75 {
		status = "diterima"
	}
	mhs = append(mhs, Mahasiswa{ID: id, Nama: nama, Jurusan: jurusan, NilaiTes: nilaiTes, Status: status})
	return mhs
}

func editMahasiswa(mhs []Mahasiswa, id string, nama, jurusan string, nilaiTes float64) []Mahasiswa {
	for i := 0; i < len(mhs); i++ {
		if mhs[i].ID == id {
			mhs[i].Nama = nama
			mhs[i].Jurusan = jurusan
			mhs[i].NilaiTes = nilaiTes
			if nilaiTes >= 75 {
				mhs[i].Status = "✅"
			} else {
				mhs[i].Status = "❌"
			}
			break
		}
	}
	return mhs
}

func hapusMahasiswa(mhs []Mahasiswa, id string) []Mahasiswa {
	for i := 0; i < len(mhs); i++ {
		if mhs[i].ID == id {
			mhs = append(mhs[:i], mhs[i+1:]...)
			break
		}
	}
	return mhs
}

func viewMhs(mhs []Mahasiswa, jurusanFilter string) {
	t := table.New(os.Stdout)
	t.SetHeaders("ID", "Nama", "Jurusan", "Nilai Tes", "Status")

	for _, m := range mhs {
		if jurusanFilter == "" || m.Jurusan == jurusanFilter {
			t.AddRow(m.ID, m.Nama, m.Jurusan, fmt.Sprintf("%.2f", m.NilaiTes), m.Status)
		}
	}

	t.Render()
}

func sortNilai(mhs []Mahasiswa, ascending bool) {
	if ascending {
		// Sequential search
		for i := 0; i < len(mhs)-1; i++ {
			for j := i + 1; j < len(mhs); j++ {
				if mhs[i].NilaiTes > mhs[j].NilaiTes {
					mhs[i], mhs[j] = mhs[j], mhs[i]
				}
			}
		}
	} else {
		binarySort(mhs)
	}
}

func binarySort(mhs []Mahasiswa) {
	for i := 1; i < len(mhs); i++ {
		key := mhs[i]
		low, high := 0, i-1

		for low <= high {
			mid := (low + high) / 2
			if mhs[mid].NilaiTes < key.NilaiTes {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}

		for j := i - 1; j >= low; j-- {
			mhs[j+1] = mhs[j]
		}
		mhs[low] = key
	}
}

func drawMenu[T any](options []T, selected int) {
	clearScreen()
	for i, option := range options {
		var optionStr string
		if str, ok := any(option).(string); ok {
			optionStr = str
		} else {
			optionStr = fmt.Sprintf("%v", option)
		}
		if i == selected {
			fmt.Println(highlightText(optionStr))
		} else {
			fmt.Println(optionStr)
		}
	}
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func enableTermbox() {
	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to reinitialize termbox:", err)
		os.Exit(1)
	}
}

func MenuControl[l any](list []l, controller *int) bool {
	for {
		drawMenu(list, (*controller))
		event := termbox.PollEvent()
		if event.Key == termbox.KeyArrowUp {
			if (*controller) > 0 {
				(*controller)--
			}
		} else if event.Key == termbox.KeyArrowDown {
			if (*controller) < len(list)-1 {
				(*controller)++
			}
		} else if event.Key == termbox.KeyEnter {
			return true
		} else if event.Key == termbox.KeyEsc {
			return false
		}
	}
}

func exportJurusanToCSV(jurusan []Jurusan, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// headers
	writer.Write([]string{"No", "Nama Jurusan"})

	// data
	for i, j := range jurusan {
		writer.Write([]string{strconv.Itoa(i + 1), string(j)})
	}

	fmt.Println("Jurusan data exported to", filename)
}

func exportMahasiswaToCSV(mahasiswa []Mahasiswa, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Nama", "Jurusan", "Nilai Tes", "Status"})

	for _, m := range mahasiswa {
		writer.Write([]string{
			m.ID, m.Nama, m.Jurusan,
			fmt.Sprintf("%.2f", m.NilaiTes), m.Status,
		})
	}

	fmt.Println("Mahasiswa data exported to", filename)
}

func exportToText(jurusan []Jurusan, mahasiswa []Mahasiswa, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("Daftar Jurusan:\n")
	for i, j := range jurusan {
		writer.WriteString(fmt.Sprintf("%d. %s\n", i+1, j))
	}
	writer.WriteString("\n")

	writer.WriteString("Daftar Mahasiswa:\n")
	for _, m := range mahasiswa {
		writer.WriteString(fmt.Sprintf(
			"ID: %s, Nama: %s, Jurusan: %s, Nilai Tes: %.2f, Status: %s\n",
			m.ID, m.Nama, m.Jurusan, m.NilaiTes, m.Status))
	}

	fmt.Println("Data exported to", filename)
}

func importJurusanFromCSV(filename string) ([]Jurusan, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jurusan []Jurusan
	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 2 {
			return nil, errors.New("invalid CSV format for Jurusan")
		}

		jurusan = append(jurusan, Jurusan(record[1]))
	}

	return jurusan, nil
}

func importMahasiswaFromCSV(filename string) ([]Mahasiswa, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mahasiswa []Mahasiswa
	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 5 {
			return nil, errors.New("invalid CSV format for Mahasiswa")
		}

		nilaiTes, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		mahasiswa = append(mahasiswa, Mahasiswa{
			ID:       record[0],
			Nama:     record[1],
			Jurusan:  record[2],
			NilaiTes: nilaiTes,
			Status:   record[4],
		})
	}

	return mahasiswa, nil
}

func generateAkunMhs(mhs []Mahasiswa) []User {
	var usr []User

	for _, m := range mhs {
		user := User{
			Name: m.Nama,
			Pass: m.Nama + "123",
		}
		usr = append(usr, user)
	}
	return usr
}

func main() {
	jurus, _ := importJurusanFromCSV("jurusan.csv")
	mhs, _ := importMahasiswaFromCSV("mahasiswa.csv")
	// user := generateAkunMhs()

	options := []string{
		"Jurusan : Tampilkan",
		"Jurusan : Tambah",
		"Jurusan : Edit",
		"Jurusan : Hapus",
		"Mahasiswa : Tampilkan",
		"Mahasiswa : Tambah",
		"Mahasiswa : Edit",
		"Mahasiswa : Hapus",
		"Mahasiswa : Tampilkan per Jurusan",
		"Mahasiswa : Urutkan dari Nilai",
		"Export : Jurusan ke CSV",
		"Export : Mahasiswa ke CSV",
		"Export : Semua ke Text",
		"Extra : ganti data dengan dummy",
	}

	selected := 0
	var loginInfo [2]string
	signedIn := false

	enableTermbox()

	for !signedIn {
		loginField := []string{"Username\t" + loginInfo[0], "Password\t" + loginInfo[1], "\nLogin"}
		if MenuControl(loginField, &selected) {
			if selected == 2 {
				if loginInfo[0] == "admin" && loginInfo[1] == "1234" {
					signedIn = true
					break
				} else if loginInfo[0] == "mahasiswa" && loginInfo[1] == "1234" {
					options = []string{
						"Jurusan : Tampilkan",
						"Mahasiswa : Tampilkan",
						"Mahasiswa : Edit",
					}
					signedIn = true
					break
				}
			}
			clearScreen()
			termbox.Close()
			loginInfo[selected] = readInput(highlightText(loginField[selected]))
		} else {

		}
		enableTermbox()
	}

	for {
		if MenuControl(options, &selected) {
			clearScreen()
			termbox.Close()
			switch options[selected] {
			case "Jurusan : Tampilkan":
				viewJurusan(jurus)
			case "Jurusan : Tambah":
				nama := readInput("Masukkan Nama Jurusan: ")
				jurus = tambahJurusan(jurus, nama)
				fmt.Println("Jurusan berhasil ditambahkan!")
			case "Jurusan : Edit":
				enableTermbox()
				selected := 0
				if MenuControl(jurus, &selected) {
					clearScreen()
					termbox.Close()
					newName := readInput("Masukkan Nama Baru: ")
					jurus, mhs = editJurusan(jurus, mhs, string(jurus[selected]), newName)
				} else {
					termbox.Close()
				}
			case "Jurusan : Hapus":
				enableTermbox()
				i := 0
				if MenuControl(jurus, &i) {
					clearScreen()
					termbox.Close()
					jurus = hapusJurusan(jurus, string(jurus[i]))
				} else {
					termbox.Close()
				}
			case "Mahasiswa : Tampilkan":
				viewMhs(mhs, "")
			case "Mahasiswa : Tambah":
				normalField := [4]string{"ID\t", "Nama\t", "Jurusan\t", "Nilai\t"}
				var option [4]string
				selected := 0
				finished := false
				for !finished {
					fields := []string{"ID\t" + option[0], "Nama\t" + option[1], "Jurusan\t" + option[2], "Nilai\t" + option[3], "\nTambah"}
					enableTermbox()
					if MenuControl(fields, &selected) {
						if selected == 4 {
							finished = true
							option3, _ := strconv.ParseFloat(option[3], 64)
							mhs = tambahMahasiswa(mhs, option[0], option[1], option[2], option3)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
						option[selected] = readInput(highlightText(normalField[selected]))
					} else {
						termbox.Close()
						finished = true
					}
					fmt.Println("Mahasiswa berhasil ditambah!")
				}
			case "Mahasiswa : Edit":
				ID := readInput("Masukkan ID Mahasiswa yang akan diedit: ")
				normalField := [4]string{"ID\t", "Nama\t", "Jurusan\t", "Nilai\t"}
				var option [4]string
				for _, m := range mhs {
					if m.ID == ID {
						option[0] = m.ID
						option[1] = m.Nama
						option[2] = m.Jurusan
						option[3] = strconv.FormatFloat(m.NilaiTes, 'g', -1, 64)
					}
				}
				if option[0] == "" {
					fmt.Println("\nTidak ada calon Mahasiswa dengan ID tersebut!")
					break
				}
				selected := 0
				finished := false
				for !finished {
					fields := []string{"ID\t" + option[0], "Nama\t" + option[1], "Jurusan\t" + option[2], "Nilai\t" + option[3], "\nEdit"}
					enableTermbox()
					if MenuControl(fields, &selected) {
						if selected == 4 {
							finished = true
							option3, _ := strconv.ParseFloat(option[3], 64)
							mhs = editMahasiswa(mhs, option[0], option[1], option[2], option3)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
						option[selected] = readInput(highlightText(normalField[selected]))
					} else {
						termbox.Close()
						finished = true
					}
					fmt.Println("Mahasiswa berhasil diedit!")
				}
			case "Mahasiswa : Hapus":
				ID := readInput("Masukkan ID Mahasiswa yang akan dihapus: ")
				var option [4]string
				for _, m := range mhs {
					if m.ID == ID {
						option[0] = m.ID
						option[1] = m.Nama
						option[2] = m.Jurusan
						option[3] = strconv.FormatFloat(m.NilaiTes, 'g', -1, 64)
					}
				}
				if option[0] == "" {
					fmt.Println("\nTidak ada calon Mahasiswa dengan ID tersebut!")
					break
				}
				selected := 0
				finished := false
				for !finished {
					enableTermbox()
					if MenuControl([]string{"ID\t" + option[0], "Nama\t" + option[1], "Jurusan\t" + option[2], "Nilai\t" + option[3], "\nHapus"}, &selected) {
						if selected == 4 {
							finished = true
							mhs = hapusMahasiswa(mhs, ID)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
					} else {
						termbox.Close()
						finished = true
					}
				}
				fmt.Println("Mahasiswa berhasil dihapus!")
			case "Mahasiswa : Tampilkan per Jurusan":
				enableTermbox()
				selected := 0
				if MenuControl(jurus, &selected) {
					termbox.Close()
					clearScreen()
					viewMhs(mhs, string(jurus[selected]))
				} else {
					termbox.Close()
				}
			case "Mahasiswa : Urutkan dari Nilai":
				ascending := readInput("Urutkan Ascending? (y/n): ") == "y"
				sortNilai(mhs, ascending)
				viewMhs(mhs, "")
			case "Export : Jurusan ke CSV":
				filename := readInput("Masukkan nama file (contoh: jurusan.csv): ")
				exportJurusanToCSV(jurus, filename)
			case "Export : Mahasiswa ke CSV":
				filename := readInput("Masukkan nama file (contoh: mahasiswa.csv): ")
				exportMahasiswaToCSV(mhs, filename)
			case "Export : Semua ke Text":
				filename := readInput("Masukkan nama file (contoh: data.txt): ")
				exportToText(jurus, mhs, filename)
			case "Extra : ganti data dengan dummy":
				jurus, mhs = dummy()
			}
			fmt.Print("\nTekan Enter untuk kembali ke Menu")
			readInput("")
			enableTermbox()
		} else {
			clearScreen()
			fmt.Println("\nProgram Berhenti")
			return
		}
	}
}
