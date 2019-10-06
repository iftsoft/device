package linker

const (
	STX byte = 0x02
	ETX = 0x03
	EOT = 0x04
	ENQ = 0x05
	ACK = 0x06
	EOL = 0x0D
	DLE = 0x10
	NAK = 0x15
	FS  = 0x1C
	SEP = 0x2E
	CMD = 0x43
	NEG = 0x4E
	POS = 0x50
)

const POLYNOMIAL = 0x8408

func CalcLRC(data []byte) byte {
	var out byte = 0
	if data == nil || len(data) == 0 {
		return 0
	}
	for i:=0; i<len(data); i++ {
		out ^= data[i]
	}
	return out
}

func CalcCRC16(data []byte) uint16 {
	var CRC uint16 = 0
	for i:=0; i<len(data); i++ {
		CRC ^= uint16(data[i])
		for j:=0; j<8; j++ {
			if (CRC & 0x0001) != 0 {
				CRC >>= 1
				CRC ^= POLYNOMIAL;
			} else {
				CRC >>= 1
			}
		}
	}
	return CRC
}

