package main

type phoneBook interface {
	nextNumber() string
	hasNumbersLeft() bool
}

type pbook struct {
	numbers []string
	index   int
}

func (p *pbook) nextNumber() string {
	p.index++
	return p.numbers[p.index]
}

func (p *pbook) hasNumbersLeft() bool {
	return len(p.numbers) > p.index
}

func newPhoneBook(numbers []string) phoneBook {
	return &pbook{
		numbers: numbers,
		index:   -1,
	}
}
