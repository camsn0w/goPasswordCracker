package main

import (
	"bufio"
	"cdec"
	_ "cdec"
	"crypto/md5"
	"encoding/hex"
	"os"
	"sort"
	"strings"
	"unicode"
)

func processWord(wordList []string) bool {
	hash := os.Args[1]
	for _, word := range wordList {
		for _, val := range recPerm(word) {
			if checkHash(val, hash) {
				return true
			}
		}
	}
	return false
}
func replaceStuff(wordIn string, place int) string {
	s := strings.ToLower(string(wordIn[place]))

	word := strings.Split(wordIn, "")

	switch s {
	case "i":
		word[place] = "1"
	case "s":
		word[place] = "5"
	case "o":
		word[place] = "0"
	case "a":
		word[place] = "4"
	case "e":
		word[place] = "3"
	case "t":
		word[place] = "7"
	default:
		return ""
	}

	return strings.Join(word, "")
}
func checkHash(word string, hash string) bool {
	hsh := md5.New()
	hsh.Write([]byte(word))
	hashed := hex.EncodeToString(hsh.Sum(nil))
	if hashed == hash {
		print("Password = " + word)
		return true
	}
	return false
}
func fileRead(fileName string) *bufio.Reader {
	print("Filename: " + fileName + "\n")
	inFile, err := os.Open(fileName)
	rdr := bufio.NewReader(inFile)
	if err != nil {
		print("\nERROR during fileRead\n ")
		print(err.Error())

		return nil
	}
	return rdr
}
func sortDict(mappy map[string]float64) ([]int, map[int]string) {
	flipped := map[int]string{}
	vals := []int{}
	for k, v := range mappy {
		notLetters := false
		for _, letter := range k {
			if !unicode.IsLetter(letter) {
				delete(mappy, k)
				notLetters = true
				break
			}
		}
		if notLetters {
			continue
		}
		for flipped[int(v)] != "" {
			v++
		}
		mappy[k] = v
		flipped[int(v)] = k
		vals = append(vals, int(v))
	}
	sort.Ints(vals)

	sort.Sort(sort.Reverse(sort.IntSlice(vals)))
	return vals, flipped
}
func addSymbols(word string) []string {
	symbols := []string{"^", "%", "$", "@", "!", "#"}
	result := []string{}
	for _, symbol := range symbols {
		result = append(result, symbol+word)
		result = append(result, symbol+word+symbol)
		result = append(result, word+symbol)
	}
	return result
}
func recPermHelper(listOWords []string, pos int) []string {
	if !strings.ContainsAny(listOWords[pos], "isoaet") {
		symboled := []string{}
		for _, word := range listOWords {
			symboled = append(symboled, addSymbols(word)...)
		}
		listOWords = append(listOWords, symboled...)
		return listOWords
	}
	for _, word := range listOWords[pos:] {
		for index := range word {
			temp := replaceStuff(word, index)
			if temp != "" {
				listOWords = append(listOWords, temp)
			}

		}
		pos++
	}

	return recPermHelper(listOWords, pos)
}
func recPerm(word string) []string {
	return recPermHelper([]string{word}, 0)
}
func main() {
	reader := fileRead(os.Args[2])
	commonWords := []string{}
	freqs, err := cdec.ScanFreqsFromReader(reader)
	if err != nil {
		print(err.Error())
	}
	nums, sortedDict := sortDict(freqs.Freqs)
	for i := range nums {
		commonWords = append(commonWords, sortedDict[i])
	}
	processWord(commonWords)

}
