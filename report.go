package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var yearRe = regexp.MustCompile(`\b(?:19|20)\d{2}\b`)

/* Helpers */

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func yearOf(date string) (int, bool) {
	m := yearRe.FindString(date)
	if m == "" {
		return 0, false
	}
	y, _ := strconv.Atoi(m)
	return y, true
}

/* Report */

func (b *Books) Report() {

	clearScreen()
	fmt.Println(Cyan + Bold + "------------------------" + Reset)
	fmt.Println(Cyan + Bold + "VK-BOOKS" + Reset)
	fmt.Println(Cyan + Bold + "------------------------" + Reset)

	if len(b.Books) == 0 {
		fmt.Println(Gray + "No books yet. Add one with: books add" + Reset)
		return
	}

	b.printStats()
	b.printLatest(5)
	b.printPerYear()
	b.printTop("Top genres", func(bk Book) string { return bk.Genre }, 5)
	b.printTop("Top authors", func(bk Book) string { return bk.Author }, 5)
	b.printMostReread()
}

func (b *Books) printStats() {
	pages, rereads := 0, 0
	withPages := 0 // books that have a usable page count
	for _, bk := range b.Books {
		if p := atoi(bk.Pages); p > 0 {
			pages += p
			withPages++
		}
		rereads += max(atoi(bk.ReadCount)-1, 0) // reads beyond the first
	}

	fmt.Printf("\n%sBooks:%s %d   %sPages:%s %d   %sRe-reads:%s %d\n",
		Bold, Reset, len(b.Books),
		Bold, Reset, pages,
		Bold, Reset, rereads)

	if withPages > 0 {
		avg := int(math.Round(float64(pages) / float64(withPages)))
		fmt.Printf("%sAverage length:%s %d pages\n", Bold, Reset, avg)
	}
}

func (b *Books) printLatest(n int) {
	fmt.Println(Bold + "\nLatest" + Reset)

	// Last n books, kept in the order they were added
	start := max(len(b.Books)-n, 0)
	recent := b.Books[start:]

	// Width of the id column, so titles line up
	idW := 0
	for _, bk := range recent {
		idW = max(idW, len(strconv.Itoa(bk.Id)))
	}

	for i, bk := range recent {
		if i == len(recent)-1 {
			// Newest book: stands out
			fmt.Printf("\t%s%s#%-*d  %s %s- %s%s\n",
				Bold, Green, idW, bk.Id, bk.Name, Reset+Green, bk.Author, Reset)
			continue
		}
		fmt.Printf("\t%s#%-*d%s  %s%s%s %s- %s%s\n",
			Gray, idW, bk.Id, Reset,
			Yellow, bk.Name, Reset,
			Gray, bk.Author, Reset)
	}
}

func (b *Books) printPerYear() {
	counts := make(map[int]int)
	unknown := 0
	for _, bk := range b.Books {
		if y, ok := yearOf(bk.Date); ok {
			counts[y]++
		} else {
			unknown++
		}
	}
	if len(counts) == 0 {
		return
	}

	years := make([]int, 0, len(counts))
	for y := range counts {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years))) // newest first

	// Scale bars so a big year can't overflow the line
	maxCount := 0
	for _, c := range counts {
		maxCount = max(maxCount, c)
	}
	const maxBar = 20

	thisYear := time.Now().Year()

	fmt.Println(Bold + "\nBooks per year" + Reset)
	for _, y := range years {
		c := counts[y]
		barLen := max(int(math.Round(float64(c)/float64(maxCount)*maxBar)), 1)

		color := Cyan
		if y == thisYear {
			color = Yellow
		}
		fmt.Printf("\t%d  %3d  %s%s%s\n", y, c, color, strings.Repeat("█", barLen), Reset)
	}

	if unknown > 0 {
		fmt.Printf("\t%s%d without a valid date%s\n", Gray, unknown, Reset)
	}
}

func (b *Books) printTop(title string, field func(Book) string, n int) {
	counts := make(map[string]int)
	for _, bk := range b.Books {
		if key := strings.TrimSpace(field(bk)); key != "" {
			counts[key]++
		}
	}
	if len(counts) == 0 {
		return
	}

	type kv struct {
		key   string
		count int
	}
	rows := make([]kv, 0, len(counts))
	for k, c := range counts {
		rows = append(rows, kv{k, c})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].key < rows[j].key
	})
	if len(rows) > n {
		rows = rows[:n]
	}

	nameW := 0
	for _, r := range rows {
		nameW = max(nameW, utf8.RuneCountInString(r.key))
	}

	fmt.Println(Bold + "\n" + title + Reset)
	for _, r := range rows {
		fmt.Printf("\t%-*s  %2d  %s%s%s\n",
			nameW, r.key, r.count, Cyan, strings.Repeat("█", min(r.count, 20)), Reset)
	}
}

func (b *Books) printMostReread() {
	best := -1
	bestReads := 1
	for i, bk := range b.Books {
		if r := atoi(bk.ReadCount); r > bestReads {
			best, bestReads = i, r
		}
	}
	if best < 0 {
		return
	}

	bk := b.Books[best]
	fmt.Println(Bold + "\nMost re-read" + Reset)
	fmt.Printf("\t%s#%d%s  %s%s%s %s- %s  (%d times)%s\n",
		Gray, bk.Id, Reset,
		Yellow, bk.Name, Reset,
		Gray, bk.Author, bestReads, Reset)
}

/* History */

func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	return string(r[:w-1]) + "…"
}

func padRight(s string, w int) string {
	if n := utf8.RuneCountInString(s); n < w {
		return s + strings.Repeat(" ", w-n)
	}
	return s
}

func padLeft(s string, w int) string {
	if n := utf8.RuneCountInString(s); n < w {
		return strings.Repeat(" ", w-n) + s
	}
	return s
}

func (b *Books) PrintHistory() {
	fmt.Println(Cyan + Bold + "------------------------" + Reset)
	fmt.Println(Cyan + Bold + "HISTORY" + Reset)
	fmt.Println(Cyan + Bold + "------------------------" + Reset)

	if len(b.Books) == 0 {
		fmt.Println(Gray + "No books yet. Add one with: books add" + Reset)
		return
	}

	const (
		maxTitle   = 32
		maxAuthor  = 22
		maxGenre   = 14
		maxOpinion = 80
	)

	// Group label for a book: its year, or "Unknown"
	label := func(bk Book) string {
		if y, ok := yearOf(bk.Date); ok {
			return strconv.Itoa(y)
		}
		return "Unknown"
	}

	// Column widths and books per year
	idW, titleW, authorW, genreW, pagesW := 0, 0, 0, 0, 0
	perYear := make(map[string]int)
	pagesText := func(bk Book) string {
		if p := atoi(bk.Pages); p > 0 {
			return strconv.Itoa(p) + "p"
		}
		return "-"
	}
	for _, bk := range b.Books {
		perYear[label(bk)]++
		idW = max(idW, len(strconv.Itoa(bk.Id)))
		titleW = max(titleW, utf8.RuneCountInString(truncate(bk.Name, maxTitle)))
		authorW = max(authorW, utf8.RuneCountInString(truncate(bk.Author, maxAuthor)))
		genreW = max(genreW, utf8.RuneCountInString(truncate(bk.Genre, maxGenre)))
		pagesW = max(pagesW, utf8.RuneCountInString(pagesText(bk)))
	}

	current := ""
	for i, bk := range b.Books {
		// New year heading when the year changes
		if l := label(bk); l != current {
			current = l
			fmt.Printf("\n%s%s%s %s(%d)%s\n", Cyan+Bold, l, Reset, Gray, perYear[l], Reset)
		}

		isLast := i == len(b.Books)-1
		titleColor := Yellow
		if isLast {
			titleColor = Bold + Green // newest book stands out
		}

		line := fmt.Sprintf("\t%s#%-*d%s  %s%s%s  %s%s%s  %s%s%s  %s%s%s  %s%s%s",
			Gray, idW, bk.Id, Reset,
			titleColor, padRight(truncate(bk.Name, maxTitle), titleW), Reset,
			Gray, padRight(truncate(bk.Author, maxAuthor), authorW), Reset,
			Cyan, padRight(truncate(bk.Genre, maxGenre), genreW), Reset,
			Gray, padLeft(pagesText(bk), pagesW), Reset,
			Gray, bk.Date, Reset)

		if r := atoi(bk.ReadCount); r > 1 {
			line += fmt.Sprintf("  %s×%d%s", Purple, r, Reset)
		}
		fmt.Println(line)

		// Opinion on its own line, if there is one
		if op := strings.TrimSpace(bk.Opinion); op != "" {
			fmt.Printf("\t%s%s\"%s\"%s\n", Gray, Italic, truncate(op, maxOpinion), Reset)
		}
	}

	fmt.Printf("\n%s%d books%s\n", Bold, len(b.Books), Reset)
}
