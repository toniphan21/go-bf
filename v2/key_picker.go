package bf

import "math"

type Key uint32

const uintSize = 32 << (^uint(0) >> 63) // 32 or 64

type KeyPicker struct {
	Source []uint
}

//func (p *KeyPicker) Pick(configs []ConfigBlock) [][]Key {
//	return nil
//}

func (p *KeyPicker) pickKey(size, index int) Key {
	end := index + size
	indexStart := index / uintSize
	indexEnd := end / uintSize

	remainderStart := index % uintSize
	remainderEnd := end % uintSize

	var maskStart uint = math.MaxUint << remainderStart
	var maskEnd uint = math.MaxUint >> (uintSize - remainderEnd)
	if indexStart == indexEnd {
		return Key(p.Source[indexStart] & maskStart & maskEnd >> remainderStart)
	}

	result := Key(p.Source[indexStart] & maskStart >> remainderStart)
	if remainderEnd == 0 {
		return result
	}
	return result | Key(p.Source[indexEnd]&maskEnd)<<(uintSize-remainderStart)
}
