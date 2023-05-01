package bot

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadQuestions() []string {
	var questions []string
	fileQuestions, errq := os.Open("bot/questions.txt")
	if errq != nil {
		log.Fatal(errq)
	}
	defer fileQuestions.Close()

	scannerQ := bufio.NewScanner(fileQuestions)

	for scannerQ.Scan() {
		questions = append(questions, scannerQ.Text())
	}

	if err := scannerQ.Err(); err != nil {
		log.Fatal(err)
	}

	return questions

}

func LoadAnswers() map[string][]string {

	var answers = make(map[string][]string)

	fileAnswers, err := os.Open("bot/answers.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fileAnswers.Close()

	scannerAnswers := bufio.NewScanner(fileAnswers)
	var q string
	for scannerAnswers.Scan() {
		line := scannerAnswers.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "¿") {
			q = line
		} else {
			answers[q] = append(answers[q], line)
		}
	}
	if err := scannerAnswers.Err(); err != nil {
		log.Fatal(err)
	}

	return answers

}
