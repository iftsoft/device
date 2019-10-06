package driver

func CheckCardPan(card string) bool {
	size := len(card)
	if size < 13 { return false }
	if size > 24 { return false }

	buff := make([]byte, 24)
	for i:=0; i< size; i++ {
		buff[i] = card[size-i-1] - '0'
	}
	var sum byte = 0
	for i:=0; i<24; i++ {
		if i%2 != 0 {
			buff[i] *= 2
		}
		if buff[i] > 9 {
			buff[i] -= 9
		}
		sum += buff[i]
	}
	if sum%10 == 0 {
		return true
	}
	return false
}

func CheckBarCode(card string) bool {
	size := len(card)
	if size < 13 { return false }
	if size > 24 { return false }

	buff := make([]byte, 24)
	for i:=0; i< size; i++ {
		buff[i] = card[size-i-1] - '0'
	}
	var sum byte = 0
	for i:=0; i<24; i++ {
		if i%2 != 0 {
			buff[i] *= 3
		}
		buff[i] = buff[i]%10
		sum += buff[i]
	}
	if sum%10 == 0 {
		return true
	}
	return false
}

func CheckBankAccount(acnt, bank string) bool {
	len1 := len(bank);
	len2 := len(acnt);
	if len1<5 || len2<5 { return false }

	mask := []byte { 1, 3, 7, 1, 3, 3, 7, 1, 3, 0, 1, 3, 7, 1, 3, 7, 1, 3, 7, 0 }
	buff := []byte { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }
	sum  := 0

	for i:=0; i<5; i++	{
		buff[i] = bank[i] - '0'
	}
	for i:=0; i<len2; i++	{
		buff[i+5] = acnt[i] - '0'
	}
	for i:=0; i<20; i++	{
		sum += int(mask[i]*buff[i]) % 10
	}
	key := byte((((sum+len2)%10)*7)%10)
	if key == buff[9] {
		return true
	}
	return false
}

func CheckFirmRegistryCode(code string) bool {
	if len(code) != 8 { return false }
	getFirmRegKey := func(buff, keys []byte) byte {
		var sum byte
		for i:=0; i<7; i++ {
			sum += buff[i] * keys[i];
		}
		return sum%11
	}
	keys1 := []byte { 1, 2, 3, 4, 5, 6, 7, 0 }
	keys2 := []byte { 7, 1, 2, 3, 4, 5, 6, 0 }
	keys3 := []byte { 3, 4, 5, 6, 7, 8, 9, 0 }
	keys4 := []byte { 9, 3, 4, 5, 6, 7, 8, 0 }
	buff  := []byte { 0, 0, 0, 0, 0, 0, 0, 0 }
	for i:=0; i<8; i++	{
		buff[i] = code[i] - '0'
	}
	var key byte = 0
	if buff[0] < 3  || buff[0] >=6 	{
		key = getFirmRegKey(buff, keys1)
		if key == 10 {
			key = getFirmRegKey(buff, keys3)
		}
		if key == 10 {
			key = 0
		}
	} else {
		key = getFirmRegKey(buff, keys2);
		if key == 10 {
			key = getFirmRegKey(buff, keys4)
		}
		if key == 10 {
			key = 0
		}
	}
	if key == buff[7] {
		return true
	}
	return false
}

// ABCDEFGHIZ – tax code
// check_sum = A2*(-1)+B2*5+C2*7+D2*9+E2*4+F2*6+G2*10+H2*5+I2*7
// last_byte = check_sum % 11
func CheckPersonalTaxNumber(inn string) bool {
	if len(inn) != 8 { return false }
	buff  := []int { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }
	for i:=0; i<10; i++	{
		buff[i] = int(inn[i] - '0')
	}
	sum := 0
	sum += buff[0] * -1;	// A2*(-1)
	sum += buff[1] * 5;		// B2*5
	sum += buff[2] * 7;		// C2*7
	sum += buff[3] * 9;		// D2*9
	sum += buff[4] * 4;		// E2*4
	sum += buff[5] * 6;		// F2*6
	sum += buff[6] * 10;	// G2*10
	sum += buff[7] * 5;		// H2*5
	sum += buff[8] * 7;		// I2*7
	if (sum%11)%10 == buff[9] {
		return true
	}
	return false
}




/*


//////////////////////////////////////////////////////////////////////
// При вводе ИНН осуществляется проверка корректности номера по следующему правилу:
// ABCDEFGHIZ – ИНН
// Вычисляем контрольную сумму = A2*(-1)+B2*5+C2*7+D2*9+E2*4+F2*6+G2*10+H2*5+I2*7
// Вычисляем остаток от деления контрольной суммы на 11.
// Если остаток = Z, ИНН – верный. В другом случае – не верный

bool IftCheckPersonalTaxNumber(LPCTSTR inn)
{
	unsigned char buff[10]  = { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 };

	if( _tcslen(inn) != 10 )	return false;
	for(int i=0; i<10; i++)		buff[i] = inn[i] - '0';

	int sum  = 0;
	sum += buff[0] * -1;	// A2*(-1)
	sum += buff[1] * 5;		// B2*5
	sum += buff[2] * 7;		// C2*7
	sum += buff[3] * 9;		// D2*9
	sum += buff[4] * 4;		// E2*4
	sum += buff[5] * 6;		// F2*6
	sum += buff[6] * 10;	// G2*10
	sum += buff[7] * 5;		// H2*5
	sum += buff[8] * 7;		// I2*7
	if( (sum%11)%10 == buff[9] )	return true;
	return false;
}

*/
