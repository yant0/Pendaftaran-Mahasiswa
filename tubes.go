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
	Role string // "admin"
}

// /$$   /$$                                     /$$$$$$             /$$                          /$$$$$$
// | $$  | $$                                    |_  $$_/            | $$                         /$$__  $$
// | $$  | $$  /$$$$$$$  /$$$$$$   /$$$$$$         | $$   /$$$$$$$  /$$$$$$    /$$$$$$   /$$$$$$ | $$  \__//$$$$$$   /$$$$$$$  /$$$$$$
// | $$  | $$ /$$_____/ /$$__  $$ /$$__  $$        | $$  | $$__  $$|_  $$_/   /$$__  $$ /$$__  $$| $$$$   |____  $$ /$$_____/ /$$__  $$
// | $$  | $$|  $$$$$$ | $$$$$$$$| $$  \__/        | $$  | $$  \ $$  | $$    | $$$$$$$$| $$  \__/| $$_/    /$$$$$$$| $$      | $$$$$$$$
// | $$  | $$ \____  $$| $$_____/| $$              | $$  | $$  | $$  | $$ /$$| $$_____/| $$      | $$     /$$__  $$| $$      | $$_____/
// |  $$$$$$/ /$$$$$$$/|  $$$$$$$| $$             /$$$$$$| $$  | $$  |  $$$$/|  $$$$$$$| $$      | $$    |  $$$$$$$|  $$$$$$$|  $$$$$$$
//  \______/ |_______/  \_______/|__/            |______/|__/  |__/   \___/   \_______/|__/      |__/     \_______/ \_______/ \_______/

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

func generateAkunMhs(mhs []Mahasiswa) []User {
	var usr []User

	for _, m := range mhs {
		user := User{
			Name: m.Nama,
			Pass: m.Nama + "123",
			Role: "mhs",
		}
		usr = append(usr, user)
	}
	return usr
}

func autentikasi(userList []User, username string, password string) (bool, string) {
	for _, u := range userList {
		if username == u.Name && password == u.Pass {
			return true, u.Role
		}
	}
	return false, ""
}

// /$$$$$$$$                              /$$     /$$                               /$$ /$$   /$$
// | $$_____/                             | $$    |__/                              | $$|__/  | $$
// | $$    /$$   /$$ /$$$$$$$   /$$$$$$$ /$$$$$$   /$$  /$$$$$$  /$$$$$$$   /$$$$$$ | $$ /$$ /$$$$$$   /$$   /$$
// | $$$$$| $$  | $$| $$__  $$ /$$_____/|_  $$_/  | $$ /$$__  $$| $$__  $$ |____  $$| $$| $$|_  $$_/  | $$  | $$
// | $$__/| $$  | $$| $$  \ $$| $$        | $$    | $$| $$  \ $$| $$  \ $$  /$$$$$$$| $$| $$  | $$    | $$  | $$
// | $$   | $$  | $$| $$  | $$| $$        | $$ /$$| $$| $$  | $$| $$  | $$ /$$__  $$| $$| $$  | $$ /$$| $$  | $$
// | $$   |  $$$$$$/| $$  | $$|  $$$$$$$  |  $$$$/| $$|  $$$$$$/| $$  | $$|  $$$$$$$| $$| $$  |  $$$$/|  $$$$$$$
// |__/    \______/ |__/  |__/ \_______/   \___/  |__/ \______/ |__/  |__/ \_______/|__/|__/   \___/   \____  $$
// .																							       /$$  | $$
// .																							      |  $$$$$$/
// .																							       \______/
func tambahJurusan(jurusan []Jurusan, nama string) []Jurusan {
	jurusan = append(jurusan, Jurusan(nama))
	return jurusan
}

func tambahMahasiswa(mhs []Mahasiswa, id string, nama, jurusan string, nilaiTes float64) []Mahasiswa {
	status := "ditolak"
	if nilaiTes >= 75 {
		status = "✅"
	} else {
		status = "❌"
	}
	mhs = append(mhs, Mahasiswa{ID: id, Nama: nama, Jurusan: jurusan, NilaiTes: nilaiTes, Status: status})
	return mhs
}

func viewMhs(mhs []Mahasiswa, jurusanFilter string) {
	t := table.New(os.Stdout)
	t.SetHeaders("ID", "Nama", "Jurusan", "Nilai Tes", "Status")

	for _, mahasiswa := range mhs {
		if jurusanFilter == "" || mahasiswa.Jurusan == jurusanFilter {
			t.AddRow(mahasiswa.ID, mahasiswa.Nama, mahasiswa.Jurusan, fmt.Sprintf("%.2f", mahasiswa.NilaiTes), mahasiswa.Status)
		}
	}

	t.Render()
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

func hapusJurusan(jurusan []Jurusan, nama string) []Jurusan {
	for i := 0; i < len(jurusan); i++ {
		if jurusan[i] == Jurusan(nama) {
			jurusan = append(jurusan[:i], jurusan[i+1:]...)
			break
		}
	}
	return jurusan
}

func viewJurusan(jurusan []Jurusan) {
	t := table.New(os.Stdout)
	t.SetHeaders("No", "Nama Jurusan")
	for index, nama := range jurusan {
		t.AddRow(strconv.Itoa(index+1), string(nama))
	}

	t.Render()
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

// /$$$$$$  /$$                               /$$   /$$
// /$$__  $$| $$                              |__/  | $$
// | $$  \ $$| $$  /$$$$$$   /$$$$$$   /$$$$$$  /$$ /$$$$$$   /$$$$$$/$$$$   /$$$$$$
// | $$$$$$$$| $$ /$$__  $$ /$$__  $$ /$$__  $$| $$|_  $$_/  | $$_  $$_  $$ |____  $$
// | $$__  $$| $$| $$  \ $$| $$  \ $$| $$  \__/| $$  | $$    | $$ \ $$ \ $$  /$$$$$$$
// | $$  | $$| $$| $$  | $$| $$  | $$| $$      | $$  | $$ /$$| $$ | $$ | $$ /$$__  $$
// | $$  | $$| $$|  $$$$$$$|  $$$$$$/| $$      | $$  |  $$$$/| $$ | $$ | $$|  $$$$$$$
// |__/  |__/|__/ \____  $$ \______/ |__/      |__/   \___/  |__/ |__/ |__/ \_______/
// 		     	  /$$  \ $$
// 		     	 |  $$$$$$/
// 		     	  \______/

func sortNilai(mhs []Mahasiswa, ascending bool) {
	if ascending {
		// Selection sort
		for i := 0; i < len(mhs)-1; i++ {
			minIdx := i
			for j := i + 1; j < len(mhs); j++ {
				if mhs[j].NilaiTes < mhs[minIdx].NilaiTes {
					minIdx = j
				}
			}
			if minIdx != i {
				mhs[i], mhs[minIdx] = mhs[minIdx], mhs[i]
			}
		}
	} else {
		// Insertion sort
		for i := 1; i < len(mhs); i++ {
			key := mhs[i]
			j := i - 1
			for j >= 0 && mhs[j].NilaiTes < key.NilaiTes {
				mhs[j+1] = mhs[j]
				j--
			}
			mhs[j+1] = key
		}
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

func main() {
	jurus, _ := importJurusanFromCSV("jurusan.csv")
	mhs, _ := importMahasiswaFromCSV("mahasiswa.csv")
	user := generateAkunMhs(mhs)
	user = append(user, User{"admin", "admin123", "admin"})

	pilihan := []string{
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

	ygDipilih := 0
	var loginInfo [2]string
	signedIn := false
	role := ""

	enableTermbox()

	for !signedIn {
		loginAja := []string{"Username : " + loginInfo[0], "Password : " + loginInfo[1], "\nLogin"}
		loginField := []string{"Username\t", "Password\t", "\nLogin"}
		if MenuControl(loginAja, &ygDipilih) {
			if ygDipilih == 2 {
				signedIn, role = autentikasi(user, loginInfo[0], loginInfo[1])
				if role == "mhs" {
					pilihan = []string{
						"Jurusan : Tampilkan",
						"Mahasiswa : Tampilkan",
					}
				}
			} else {
				clearScreen()
				termbox.Close()
				loginInfo[ygDipilih] = readInput(highlightText(loginField[ygDipilih]))
				enableTermbox()
			}
		} else {
			clearScreen()
			fmt.Println("\nProgram Berhenti")
			return
		}
	}

	ygDipilih = 0

	for {
		if MenuControl(pilihan, &ygDipilih) {
			clearScreen()
			termbox.Close()
			switch pilihan[ygDipilih] {
			case "Jurusan : Tampilkan":
				viewJurusan(jurus)
			case "Jurusan : Tambah":
				nama := readInput("Masukkan Nama Jurusan: ")
				jurus = tambahJurusan(jurus, nama)
				fmt.Println("Jurusan berhasil ditambahkan!")
			case "Jurusan : Edit":
				enableTermbox()
				ygDipilih := 0
				if MenuControl(jurus, &ygDipilih) {
					clearScreen()
					termbox.Close()
					newName := readInput("Masukkan Nama Baru: ")
					jurus, mhs = editJurusan(jurus, mhs, string(jurus[ygDipilih]), newName)
				} else {
					termbox.Close()
				}
			case "Jurusan : Hapus":
				enableTermbox()
				ygDipilih := 0
				if MenuControl(jurus, &ygDipilih) {
					clearScreen()
					termbox.Close()
					jurus = hapusJurusan(jurus, string(jurus[ygDipilih]))
				} else {
					termbox.Close()
				}
			case "Mahasiswa : Tampilkan":
				viewMhs(mhs, "")
			case "Mahasiswa : Tambah":
				masukanAja := [4]string{"ID\t", "Nama\t", "Jurusan\t", "Nilai\t"}
				var input [4]string
				ygDipilih := 0
				selesai := false
				for !selesai {
					Masukan := []string{"ID\t" + input[0], "Nama\t" + input[1], "Jurusan\t" + input[2], "Nilai\t" + input[3], "\nTambah"}
					enableTermbox()
					if MenuControl(Masukan, &ygDipilih) {
						if ygDipilih == 4 {
							selesai = true
							option3, _ := strconv.ParseFloat(input[3], 64)
							mhs = tambahMahasiswa(mhs, input[0], input[1], input[2], option3)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
						input[ygDipilih] = readInput(highlightText(masukanAja[ygDipilih]))
					} else {
						termbox.Close()
						selesai = true
					}
					fmt.Println("Mahasiswa berhasil ditambah!")
				}
			case "Mahasiswa : Edit":
				ID := readInput("Masukkan ID Mahasiswa yang akan diedit: ")
				masukanAja := [4]string{"ID\t", "Nama\t", "Jurusan\t", "Nilai\t"}
				var input [4]string
				for _, m := range mhs {
					if m.ID == ID {
						input[0] = m.ID
						input[1] = m.Nama
						input[2] = m.Jurusan
						input[3] = strconv.FormatFloat(m.NilaiTes, 'g', -1, 64)
					}
				}
				if input[0] == "" {
					fmt.Println("\nTidak ada calon Mahasiswa dengan ID tersebut!")
					break
				}
				ygDipilih := 0
				selesai := false
				for !selesai {
					fields := []string{"ID\t" + input[0], "Nama\t" + input[1], "Jurusan\t" + input[2], "Nilai\t" + input[3], "\nEdit"}
					enableTermbox()
					if MenuControl(fields, &ygDipilih) {
						if ygDipilih == 4 {
							selesai = true
							option3, _ := strconv.ParseFloat(input[3], 64)
							mhs = editMahasiswa(mhs, input[0], input[1], input[2], option3)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
						input[ygDipilih] = readInput(highlightText(masukanAja[ygDipilih]))
					} else {
						termbox.Close()
						selesai = true
					}
					fmt.Println("Mahasiswa berhasil diedit!")
				}
			case "Mahasiswa : Hapus":
				ID := readInput("Masukkan ID Mahasiswa yang akan dihapus: ")
				var masukan [4]string
				for _, m := range mhs {
					if m.ID == ID {
						masukan[0] = m.ID
						masukan[1] = m.Nama
						masukan[2] = m.Jurusan
						masukan[3] = strconv.FormatFloat(m.NilaiTes, 'g', -1, 64)
					}
				}
				if masukan[0] == "" {
					fmt.Println("\nTidak ada calon Mahasiswa dengan ID tersebut!")
					break
				}
				ygDipilih := 0
				selesai := false
				for !selesai {
					enableTermbox()
					if MenuControl([]string{"ID\t" + masukan[0], "Nama\t" + masukan[1], "Jurusan\t" + masukan[2], "Nilai\t" + masukan[3], "\nHapus"}, &ygDipilih) {
						if ygDipilih == 4 {
							selesai = true
							mhs = hapusMahasiswa(mhs, ID)
							clearScreen()
							termbox.Close()
							break
						}
						clearScreen()
						termbox.Close()
					} else {
						termbox.Close()
						selesai = true
					}
				}
				fmt.Println("Mahasiswa berhasil dihapus!")
			case "Mahasiswa : Tampilkan per Jurusan":
				enableTermbox()
				ygDipilih := 0
				if MenuControl(jurus, &ygDipilih) {
					termbox.Close()
					clearScreen()
					viewMhs(mhs, string(jurus[ygDipilih]))
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
			termbox.Close()
			fmt.Println("\nProgram Berhenti")
			return
		}
	}
}
