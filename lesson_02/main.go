package main

import (
	"strconv"
)

func main() {
	println(FibonacciIterative(10))
	println(FibonacciRecursive(10))
	println(IsPrime(29))
	println(IsPrime(16))
	println(IsBinaryPalindrome(7))
	println(IsBinaryPalindrome(6))
	println(ValidParentheses("[{}]"))
	println(ValidParentheses("[{]}"))
	println(Increment("101"))
}

func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	prev, current := 0, 1
	for i := 2; i <= n; i++ {
		next := prev + current
		prev = current
		current = next
	}
	return current
}

func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsBinaryPalindrome(n int) bool {
	binStr := strconv.FormatInt(int64(n), 2)
	i, j := 0, len(binStr)-1
	for i < j {
		if binStr[i] != binStr[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func ValidParentheses(s string) bool {
	stack := []rune{}
	for _, ch := range s {
		switch ch {
		case '(', '[', '{':
			stack = append(stack, ch)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func Increment(num string) int {
	value, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		println("Invalid binary number")
	}
	return int(value + 1)
}
