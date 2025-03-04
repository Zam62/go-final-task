package calculation

import (
	"fmt"
	"strconv"
	"strings"
)

func CalcOld(expression string) (float64, error) {
	if expression == "" {
		return 0, NewExpressionTooShortError()
	}

	bracketCount := 0
	for _, char := range expression {
		if char == '(' {
			bracketCount++
		} else if char == ')' {
			bracketCount--
		}
		if bracketCount < 0 {
			return 0, NewBracketsExpressionError(fmt.Errorf("некорректное расположение скобок"))
		}
	}
	if bracketCount != 0 {
		return 0, NewBracketsExpressionError(fmt.Errorf("некорректное расположение скобок"))
	}

	return evalExpression(expression)
}

func Calc(expression string) (float64, error) {
	arr := strings.Split(expression, "")
	var stack []string
	var queue []string //выходная строка

	//Проверка
	count_digits := 0
	count_operators := 0
	count_brackets := 0

	for i := 0; i < len(arr); i++ {
		ch := arr[i]

		if !IsDigit(ch) &&
			!IsOperation(ch) &&
			ch != "(" && ch != ")" {
			return 0, fmt.Errorf("invalid sybmol")
		}

		if IsDigit(ch) {
			count_digits++
		}
		if IsOperation(ch) {
			count_operators++
		}
		if ch == "(" {
			count_brackets++
		}
		if ch == ")" {
			count_brackets--
		}
		if count_brackets < 0 {
			return 0, fmt.Errorf("invalid input")
		}
	}

	if count_operators == count_digits {
		return 0, fmt.Errorf("invalid input")
	}

	//получим в обратной польской записи

	//цикл по всем символам
	for i := 0; i < len(arr); i++ {
		s := arr[i]

		if IsDigit(s) {
			//если число - помещаем в выходную строку
			queue = append(queue, s)
		}
		if IsOperation(s) {
			//оператор
			if len(stack) == 0 || stack[len(stack)-1] == "(" {
				//Если в стеке пусто, или в стеке открывающая скобка
				//добавляем оператор в стек
				stack = PushToStack(s, stack)
			} else if GetPriority(s) > GetPriority(stack[len(stack)-1]) {
				//Если входящий оператор имеет более высокий приоритет чем вершина stack
				//добавляем оператор в стек
				stack = PushToStack(s, stack)
			} else if GetPriority(s) <= GetPriority(stack[len(stack)-1]) {
				//Если оператор имеет более низкий или равный приоритет, чем в стеке,
				//выгружаем POP в очередь (QUEUE),
				//пока не увидите оператор с меньшим приоритетом или левую скобку на вершине (TOP),
				stack, queue = PopStackToQueue(s, stack, queue)
				// затем добавьте (PUSH) входящий оператор в стек (STACK).
				stack = PushToStack(s, stack)
			}
		} else if s == "(" {
			//Если входящий элемент является левой скобкой,
			// поместите (PUSH) его в стек (STACK).
			stack = PushToStack(s, stack)
		} else if s == ")" {
			//Если входящий элемент является правой скобкой,
			//выгружаем стек (POP) и добавляем его элементы в очередь (QUEUE),
			//пока не увидите левую круглую скобку.
			//Удалите найденную скобку из стека (STACK).
			stack, queue = PopStackToQueue(s, stack, queue)
		}

	}
	// В конце выражения выгрузите стек (POP) в очередь (QUEUE)
	stack, queue = PopStackToQueue("", stack, queue)

	//   посчитаем выражение
	//   queue_new := make([]string, cap(queue), cap(queue))
	stack_new := make([]float64, cap(stack), cap(queue))

	for i := 0; i < len(queue); i++ {
		curr := queue[i]
		curr_float64, _ := strconv.ParseFloat(strings.TrimSpace(curr), 64)
		if IsDigit(curr) {
			stack_new = append(stack_new, curr_float64)
		} else {
			a := stack_new[len(stack_new)-1]
			stack_new = RemoveItemInIntSlice(stack_new, len(stack_new)-1)

			b := stack_new[len(stack_new)-1]
			stack_new = RemoveItemInIntSlice(stack_new, len(stack_new)-1)

			stack_new = append(stack_new, Operation(b, a, curr))
		}

	}

	return stack_new[len(stack_new)-1], nil
}

func evalExpression(expression string) (float64, error) {
	lastOpen := -1
	for i, char := range expression {
		if char == '(' {
			lastOpen = i
		} else if char == ')' && lastOpen != -1 {
			innerResult, err := evalSimpleExpression(expression[lastOpen+1 : i])
			if err != nil {
				return 0, NewBracketsExpressionError(err)
			}

			newExpr := expression[:lastOpen]
			newExpr += fmt.Sprintf("%g", innerResult)
			newExpr += expression[i+1:]

			return evalExpression(newExpr)
		}
	}

	return evalSimpleExpression(expression)
}

func evalSimpleExpression(expression string) (float64, error) {
	var numbers []float64
	var operators []rune
	currentNumber := ""
	lastWasOperator := true

	for i, char := range expression {
		if char == ' ' {
			continue
		}

		if (char >= '0' && char <= '9') || char == '.' {
			currentNumber += string(char)
			lastWasOperator = false
			continue
		}

		if len(currentNumber) > 0 {
			num := stringToFloat64(currentNumber)
			numbers = append(numbers, num)
			currentNumber = ""
		}

		if char == '-' && (lastWasOperator || i == 0) {
			currentNumber = "-"
			lastWasOperator = true
			continue
		}

		// if IsOperation(char) {
		// 	if lastWasOperator {
		// 		return 0, NewConsecutiveOperatorsError()
		// 	}
		// 	operators = append(operators, char)
		// 	lastWasOperator = true
		// 	continue
		// }

		return 0, NewInvalidCharError(char)
	}

	if len(currentNumber) > 0 {
		num := stringToFloat64(currentNumber)
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return 0, NewExpressionTooShortError()
	}

	if len(numbers) != len(operators)+1 {
		return 0, NewInvalidOperatorPositionError()
	}

	for i := 0; i < len(operators); {
		if operators[i] == '*' || operators[i] == '/' {
			if operators[i] == '*' {
				numbers[i+1] = numbers[i] * numbers[i+1]
			} else {
				if numbers[i+1] == 0 {
					return 0, NewDivisionByZeroError()
				}
				numbers[i+1] = numbers[i] / numbers[i+1]
			}
			numbers = append(numbers[:i], numbers[i+1:]...)
			operators = append(operators[:i], operators[i+1:]...)
		} else {
			i++
		}
	}

	result := numbers[0]
	for i := 0; i < len(operators); i++ {
		if operators[i] == '+' {
			result += numbers[i+1]
		} else if operators[i] == '-' {
			result -= numbers[i+1]
		}
	}

	return result, nil
}

func Operation(a, b float64, ch string) float64 {
	if ch == "*" {
		return a * b
	}
	if ch == "/" {
		return a / b
	}
	if ch == "+" {
		return a + b
	}
	if ch == "-" {
		return a - b
	}
	return 0
}

func IsOperation(ch string) bool {
	// проверяет является ли символ операцией
	return ch == "*" || ch == "/" || ch == "+" || ch == "-"
}

func IsDigit(ch string) bool {
	// проверяет является ли символ числом
	return ch == "0" || ch == "1" || ch == "2" || ch == "3" || ch == "4" || ch == "5" || ch == "6" || ch == "7" || ch == "8" || ch == "9"
}

func IsBracket(ch string) bool {
	// проверяет является ли символ скобкой
	return ch == "(" || ch == ")"
}

func GetPriority(ch string) int {
	if ch == "*" ||
		ch == "/" {
		return 1
	}
	return 0
}

func stringToFloat64(ch string) float64 {
	var result float64
	isNegative := false

	if ch == "" {
		return 0
	}

	if ch[0] == '-' {
		isNegative = true
		ch = ch[1:]
	}

	for _, c := range ch {
		if c >= '0' && c <= '9' {
			result = result*10 + float64(c-'0')
		}
	}

	if isNegative {
		result = -result
	}

	return result
}

func RemoveItemInSlice(slice []string, ind int) []string {
	// Удаляем значение из slice
	return append(slice[:ind], slice[ind+1:]...)
}

func RemoveItemInIntSlice(slice []float64, ind int) []float64 {
	return append(slice[:ind], slice[ind+1:]...)
}

func PushToStack(s string, stack []string) []string {
	stack = append(stack, s)
	return stack
}

func PopStackToQueue(s string, stack, queue []string) ([]string, []string) {
	//сделаем копию стека
	stack_temp := stack
	for i := 0; i < len(stack); i++ {
		stack_temp[i] = stack[i]
	}
	if s != "" {
		//если s != "" идем в обратном порядке до "(" или с меньшим приоритетом
		for i := len(stack) - 1; i >= 0; i-- {
			curr := stack[i]
			if GetPriority(curr) < GetPriority(s) || curr == "(" {
				//stop
				if curr == "(" {
					stack_temp = RemoveItemInSlice(stack_temp, i)
				}
				break
			} else {
				queue = append(queue, curr)
				stack_temp = RemoveItemInSlice(stack_temp, i)
			}
		}
	} else {
		for i := len(stack) - 1; i >= 0; i-- {
			curr := stack[i]
			queue = append(queue, curr)
			stack_temp = RemoveItemInSlice(stack_temp, i)
		}
	}

	return stack_temp, queue
}
