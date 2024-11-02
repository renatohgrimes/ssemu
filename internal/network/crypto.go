package network

import "ssemu/internal/resources"

func Decrypt(data []byte) {
	key := resources.GetCypherKey()
	dataLength := len(data)
	for i := 0; i < dataLength; i++ {
		data[i] = ((data[i] >> 1) & 0x7F) | ((data[i] << 7) & 0x80)
		data[i] ^= key[(i % 0x28)]
	}
}

func Encrypt(data []byte) {
	key := resources.GetCypherKey()
	dataLength := len(data)
	for i := 0; i < dataLength; i++ {
		data[i] ^= key[(i % 0x28)]
		data[i] = ((data[i] & 0x7F) << 0x01) | ((data[i] & 0x80) >> 0x07)
	}
}
