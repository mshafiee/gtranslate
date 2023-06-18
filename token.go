package gtranslate

import (
	"strconv"
	"strings"
)

const googleTranslateTKK = "448487.932609646"

func shiftLeftOrRightThenSumOrXor(num int32, optString string) int32 {
	// Iterate through the characters in the option string, 3 characters at a time
	for i := 0; i < len(optString)-2; i += 3 {
		// Initialize accumulator with the ASCII value of the third character in each triplet
		acc := int32(optString[i+2])

		// If the ASCII value is greater than or equal to 'a' (97 in decimal), it's a lowercase letter from 'a' to 'f'
		// Subtract 87 to convert it into a number from 10 to 15, for handling hexadecimal numbers
		if 'a' <= acc {
			acc = acc - 87
		} else {
			// If it's not a lowercase letter, it's a digit. Subtract the ASCII value of '0' to convert it into a number
			acc -= '0' // Convert from ASCII
		}

		// If the second character of the triplet is '+', perform a right shift on num by acc places
		// Otherwise, perform a left shift
		if optString[i+1] == '+' {
			acc = int32(uint32(num) >> uint32(acc))
		} else {
			acc = num << acc
		}

		// Initialize a mask with all bits set to 1
		bitmask32 := 4294967295

		// If the first character of the triplet is '+', add num and the bitwise AND of acc and bitmask32
		// Otherwise, perform a bitwise XOR between num and acc
		if optString[i] == '+' {
			num += acc & int32(bitmask32)
		} else {
			num ^= acc
		}
	}
	// Return the manipulated num value
	return num
}

func transformQuery(query string) []int32 {
	// Initialize a slice to store the UTF-8 encoded bytes
	var bytesArray []int32
	// Initialize an index variable
	idx := 0

	// Iterate over each character in the input string
	for i := 0; i < len(query); i++ {
		// Get the Unicode code point of the current character
		charCode := int32(rune(query[i]))

		// If charCode is less than 128, it's a 1-byte UTF-8 character
		// Simply append it to the bytesArray
		if charCode < 128 {
			bytesArray = append(bytesArray, charCode)
			idx++
		} else {
			// If charCode is less than 2048, it's a 2-byte UTF-8 character
			// The first byte is created by shifting charCode right by 6 bits, then bitwise OR with 192 (11000000 in binary)
			if charCode < 2048 {
				bytesArray = append(bytesArray, (charCode>>6)|192)
				idx++
			} else {
				// If charCode is in the surrogate pair range, it's a 4-byte UTF-8 character
				// The first byte is created by shifting charCode right by 18 bits, then bitwise OR with 240 (11110000 in binary)
				// The second byte is created by shifting charCode right by 12 bits, bitwise AND with 63 (00111111 in binary), then bitwise OR with 128 (10000000 in binary)
				if (charCode&64512) == 55296 && i+1 < len(query) && (int(rune(query[i+1]))&64512) == 56320 {
					charCode = 65536 + ((charCode & 1023) << 10) + (int32(rune(query[i+1])) & 1023)
					i++
					bytesArray = append(bytesArray, (charCode>>18)|240)
					bytesArray = append(bytesArray, ((charCode>>12)&63)|128)
					idx++
				} else {
					// If charCode is not in the surrogate pair range, it's a 3-byte UTF-8 character
					// The first byte is created by shifting charCode right by 12 bits, then bitwise OR with 224 (11100000 in binary)
					bytesArray = append(bytesArray, (charCode>>12)|224)
				}
				// The second byte (for 3-byte UTF-8) or the third byte (for 4-byte UTF-8) is created by shifting charCode right by 6 bits, bitwise AND with 63 (00111111 in binary), then bitwise OR with 128 (10000000 in binary)
				bytesArray = append(bytesArray, ((charCode>>6)&63)|128)
				idx++
			}
			// The last byte is created by bitwise AND of charCode with 63 (00111111 in binary), then bitwise OR with 128 (10000000 in binary)
			bytesArray = append(bytesArray, (charCode&63)|128)
			idx++
		}
	}
	// Return the UTF-8 encoded bytes
	return bytesArray
}

func calcHash(query string, windowTkk string) string {
	// Split the tkk string on the '.' character
	tkkSplited := strings.Split(windowTkk, ".")
	// Convert the first part of the split tkk string to an integer
	tkkIndex, _ := strconv.Atoi(tkkSplited[0])
	// Convert the second part of the split tkk string to an integer
	tkkKey, _ := strconv.Atoi(tkkSplited[1])

	// Transform the query string into a sequence of integers
	bytesArray := transformQuery(query)

	// Start the encoding round with the tkk index value
	encondingRound := tkkIndex
	// Loop over each byte in the array
	for i := 0; i < len(bytesArray); i++ {
		// Add the current byte to the encoding round
		encondingRound += int(bytesArray[i])
		// Perform a sequence of shifts and sums/XORs on the encoding round
		encondingRound = int(shiftLeftOrRightThenSumOrXor(int32(encondingRound), "+-a^+6"))
	}
	// Perform another sequence of shifts and sums/XORs on the encoding round
	encondingRound = int(shiftLeftOrRightThenSumOrXor(int32(encondingRound), "+-3^+b+-f"))

	// XOR the encoding round with the tkk key
	encondingRound ^= tkkKey
	// If the encoding round is less than or equal to 0, perform a bitwise AND operation with 2147483647 and add 2147483648
	if encondingRound <= 0 {
		encondingRound = (encondingRound & 2147483647) + 2147483648
	}

	// Normalize the result by taking the modulus of the encoding round with 1000000
	normalizedResult := encondingRound % 1000000
	// Return the normalized result as a string, concatenated with '.', concatenated with the string representation of the normalized result XORed with the tkk index
	return strconv.Itoa(normalizedResult) + "." + strconv.Itoa(normalizedResult^tkkIndex)
}

func getToken(query string) string {
	// Call the calcHash function with the input query and the constant googleTranslateTKK
	// The googleTranslateTKK is likely a key or seed used in the hash calculation
	// Return the result of the hash calculation
	return calcHash(query, googleTranslateTKK)
}
