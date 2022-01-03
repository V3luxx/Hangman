package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// va lire le fichier en argument
var mot = []string{}

func Load(filename string) bool {
	if strings.TrimSpace(filename) == "" {
		return false
	}
	f, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	word := ""

	for _, char := range string(f) {
		if char == '\n' {
			mot = append(mot, word)
			word = ""
		} else {
			word += string(char)
		}
	}
	return true
}

func Randmot() string {

	// renvoi un mot completement au hasard du fichier words

	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(mot))
	return mot[i]
}

type Hangman struct {
	level    string   // la ou en est la partie
	Letters  []string // Lettres a trouver
	Found    []string // lettres trouvées
	tried    []string // lettres tentées
	Tleft    int      // Tours restant
	position string
}

func New(turns int, word string) (*Hangman, error) {
	if len(word) < 3 {
		return nil, fmt.Errorf("Le mot '%s' doit avoir plus de 3 lettres, la y'en a : %v", word, len(word))
	}

	letters := strings.Split(strings.ToUpper(word), "")
	found := make([]string, len(letters))
	for i := 0; i < len(letters); i++ {
		found[i] = "_"
	}

	// affiche les premieres lettres en utilisant un compteur

	count := 0
	for {
		index := rand.Intn(len(found))
		if found[index] == "_" {
			found[index] = strings.ToUpper(string(word[index]))
			count++
		}
		if count >= len(word)/2 {
			break
		}
	}
	// initialisation des valeurs de la structure hangman
	g := &Hangman{
		level:   "",
		Letters: letters,
		Found:   found,
		tried:   []string{},
		Tleft:   turns,
	}

	return g, nil
}

// fonction qui gere l'etat de la partie en fonction de la valeur de "level"
func (g *Hangman) EssaiL(Essai string) {
	Essai = strings.ToUpper(Essai)

	//cas ou tu donne le mot
	if len(Essai) > 1 {
		correct := true
		if len(Essai) == len(g.Letters) {
			for i, c := range g.Letters {
				if c != string(Essai[i]) {
					correct = false
					g.level = "gagné"
					break
				}
			}
			// si faux ou si plus de tours
		} else {
			correct = false
			g.tried = append(g.tried, Essai)
			g.Tleft -= 2
			if g.Tleft < 0 {
				g.Tleft = 0
				g.level = "perdu"
			}
			return
		}
		// si la lettre est bien dans le mot
		if correct {
			for _, c := range g.Letters {
				g.ShowL(string(c))
			}
			g.level = "gagné"
		}
	}

	switch g.level {
	case "gagné", "perdu":
		return
	}

	if Lmot(Essai, g.tried) {
		g.level = "dejavu"
	} else if Lmot(Essai, g.Letters) {
		g.level = "juste"
		g.ShowL(Essai)

		if Gagné(g.Letters, g.Found) {
			g.level = "gagné"
		}
	} else {
		g.level = "raté"
		g.Tperdu(Essai)

		if g.Tleft <= 0 {
			g.level = "perdu"
		}
	}
}

func Gagné(letters []string, foundLetters []string) bool {
	for i := range letters {
		if letters[i] != foundLetters[i] {
			return false
		}
	}
	return true
}

func (g *Hangman) ShowL(Essai string) {
	g.tried = append(g.tried, Essai)
	for i, l := range g.Letters {
		if l == Essai {
			g.Found[i] = Essai
		}
	}
}

func (g *Hangman) Tperdu(Essai string) {
	g.Tleft--
	g.tried = append(g.tried, Essai)
}

func Lmot(Essai string, letters []string) bool {
	for _, l := range letters {
		if l == Essai {
			return true
		}
	}
	return false
}

//////// affichage

func Draw(g *Hangman, Essai string) {
	Showlevel(g, Essai)
	g.Hangmanpos()
	fmt.Println()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

///ca nous sert a afficher le petit pendu
func (g *Hangman) Hangmanpos() {
	data, err := ioutil.ReadFile("hangman.txt")
	content := string(data)
	positions := strings.Split(content, "\n\n")
	var n int = 10 - g.Tleft
	if err != nil {
		fmt.Println(err)

	} else {
		for i := 0; i < n; i++ {
			g.position = positions[i]
		}
		if n > 0 {
			fmt.Println(positions[n-1])
		}
	}

}

func drawLetters(g []string) {
	for _, c := range g {
		fmt.Printf("%v ", c)
	}
	fmt.Println()
}

func errfunc(err error) {
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
}

//cette fonction nous permet de suivre l'état de la partie
func Showlevel(g *Hangman, Essai string) {
	fmt.Println("mot à trouver: ", g.Found)

	fmt.Println("Lettres utilisées: ", g.tried)

	fmt.Println("Il reste :", g.Tleft, "tours")

	switch g.level {
	case "perdu":
		fmt.Print("Et Non!! c'étais")
		drawLetters(g.Found)
	case "gagné":
		fmt.Print("Bien joué la veinasse, c'etais bien :")
		drawLetters(g.Found)
	case "Juste":
		fmt.Print("Bienvu bogoss!")
	case "dejavu":
		fmt.Printf("T'a deja essayé '%s'...", Essai)
	case "raté":
		fmt.Printf("Dommage, je vois pas '%s' dans le mot...", Essai)
	}
	fmt.Println()
}

// Partie ou le joueur interagit

var reader = bufio.NewReader(os.Stdin)

// lis les demande de l'utilisateur
func ReadGuess() (guess string, err error) {
	valid := false
	for !valid {
		fmt.Print("c'est quoi ta lettre le sang ? ")
		guess, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		guess = strings.TrimSpace(guess)

		valid = true
	}
	return guess, nil
}

// executuion des fonctions, et donc mise en place du jeu en general

func main() {
	if len(os.Args) < 2 {
		wordsfile := "./words.txt"
		Load(wordsfile)
	} else {
		Load(os.Args[1])
	}

	g, err := New(10, Randmot())

	if err != nil {
		fmt.Printf("partie impossible %v\n", err)
		os.Exit(1)
	}
	guess := ""
	for {
		Draw(g, guess)

		switch g.level {
		case "gagné", "perdu":
			os.Exit(0)
		}

		l, err := ReadGuess()
		if err != nil {
			fmt.Printf("Erreur %v", err)
			os.Exit(1)
		}
		guess = l

		g.EssaiL(guess)
	}
}
