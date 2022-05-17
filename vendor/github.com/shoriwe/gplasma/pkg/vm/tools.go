package vm

import (
	"github.com/shoriwe/gplasma/pkg/tools"
)

func (p *Plasma) Repeat(context *Context, content []*Value, times int64) ([]*Value, *Value) {
	copyFunctions := map[int64]*Value{}
	var result []*Value
	if times == 0 {
		for _, object := range content {
			copyObject, getError := object.Get(p, context, Copy)
			if getError != nil {
				copyFunctions[object.Id()] = nil
				continue
			}
			copyFunctions[object.Id()] = copyObject
		}
	}
	for i := int64(0); i < times; i++ {
		for _, object := range content {
			copyFunction := copyFunctions[object.Id()]
			if copyFunction == nil {
				result = append(result, object)
				continue
			}
			objectCopy, success := p.CallFunction(context, copyFunction)
			if !success {
				return nil, objectCopy
			}
			result = append(result, objectCopy)
		}
	}
	return result, nil
}

func (p *Plasma) Equals(context *Context, leftHandSide *Value, rightHandSide *Value) (bool, *Value) {
	equals, getError := leftHandSide.Get(p, context, Equals)
	if getError != nil {
		// Try with the rightHandSide
		var rightEquals *Value
		rightEquals, getError = rightHandSide.Get(p, context, RightEquals)
		if getError != nil {
			return false, getError
		}
		result, success := p.CallFunction(context, rightEquals, rightHandSide)
		if !success {
			return false, result
		}
		return p.QuickGetBool(context, result)
	}
	result, success := p.CallFunction(context, equals, rightHandSide)
	if !success {
		// Try with the rightHandSide
		var rightEquals *Value
		rightEquals, getError = rightHandSide.Get(p, context, RightEquals)
		if getError != nil {
			return false, getError
		}
		result, success = p.CallFunction(context, rightEquals, rightHandSide)
		if !success {
			return false, result
		}
	}
	return p.QuickGetBool(context, result)
}

func (p *Plasma) ContentEquals(context *Context, leftHandSide *Value, rightHandSide *Value) (*Value, bool) {
	leftHandSideLength := len(leftHandSide.Content)
	rightHandSideLength := len(rightHandSide.Content)
	if leftHandSideLength != rightHandSideLength {
		return p.GetFalse(), true
	}
	var rightEquals *Value
	var comparisonResult *Value
	var success bool
	var callError *Value
	var comparisonBool bool

	for i := 0; i < leftHandSideLength; i++ {
		leftEquals, getError := rightHandSide.Content[i].Get(p, context, Equals)
		if getError != nil {
			rightEquals, getError = rightHandSide.Content[i].Get(p, context, RightEquals)
			if getError != nil {
				return getError, false
			}
			comparisonResult, success = p.CallFunction(context, rightEquals, leftHandSide.Content[i])
		} else {
			comparisonResult, success = p.CallFunction(context, leftEquals, rightHandSide.Content[i])
		}
		if !success {
			return comparisonResult, false
		}
		comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
		if callError != nil {
			return callError, false
		}
		if !comparisonBool {
			return p.GetFalse(), true
		}
	}
	return p.GetTrue(), true
}

func (p *Plasma) ContentNotEquals(context *Context, leftHandSide *Value, rightHandSide *Value) (*Value, bool) {
	leftHandSideLength := len(leftHandSide.Content)
	rightHandSideLength := len(rightHandSide.Content)
	if leftHandSideLength != rightHandSideLength {
		return p.GetTrue(), true
	}
	var rightEquals *Value
	var comparisonResult *Value
	var success bool
	var callError *Value
	var comparisonBool bool

	for i := 0; i < leftHandSideLength; i++ {
		leftEquals, getError := rightHandSide.Content[i].Get(p, context, Equals)
		if getError != nil {
			rightEquals, getError = rightHandSide.Content[i].Get(p, context, RightEquals)
			if getError != nil {
				return getError, false
			}
			comparisonResult, success = p.CallFunction(context, rightEquals, leftHandSide.Content[i])
		} else {
			comparisonResult, success = p.CallFunction(context, leftEquals, rightHandSide.Content[i])
		}
		if !success {
			return comparisonResult, false
		}
		comparisonBool, callError = p.QuickGetBool(context, comparisonResult)
		if callError != nil {
			return callError, false
		}
		if !comparisonBool {
			return p.GetTrue(), true
		}
	}
	return p.GetFalse(), true
}

func (p *Plasma) ContentContains(context *Context, source *Value, value *Value) (*Value, bool) {
	valueRightEquals, getError := value.Get(p, context, RightEquals)
	if getError != nil {
		return getError, false
	}
	for _, tupleValue := range source.Content {
		callResult, success := p.CallFunction(context, valueRightEquals, tupleValue)
		if !success {
			return callResult, false
		}
		var boolValue *Value
		if callResult.IsTypeById(BoolId) {
			boolValue = callResult
		} else {
			var boolValueToBool *Value
			boolValueToBool, getError = callResult.Get(p, context, ToBool)
			if getError != nil {
				return getError, false
			}
			callResult, success = p.CallFunction(context, boolValueToBool)
			if !success {
				return callResult, false
			}
			if !callResult.IsTypeById(BoolId) {
				return p.NewInvalidTypeError(context, callResult.TypeName(), BoolName), false
			}
			boolValue = callResult
		}
		if boolValue.Bool {
			return p.GetTrue(), true
		}
	}
	return p.GetFalse(), true
}

func (p *Plasma) ContentCopy(context *Context, source *Value) (*Value, bool) {
	var copiedObjects []*Value
	for _, contentObject := range source.Content {
		objectCopy, getError := contentObject.Get(p, context, Copy)
		if getError != nil {
			return getError, false
		}
		copiedObject, success := p.CallFunction(context, objectCopy)
		if !success {
			return copiedObject, false
		}
		copiedObjects = append(copiedObjects, copiedObject)
	}
	if source.BuiltInTypeId == TupleId {
		return p.NewTuple(context, false, copiedObjects), true
	}
	return p.NewArray(context, false, copiedObjects), true
}

func (p *Plasma) ContentIndex(context *Context, source *Value, indexObject *Value) (*Value, bool) {
	if indexObject.IsTypeById(IntegerId) {
		index, calcError := tools.CalcIndex(indexObject.Integer, len(source.Content))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Content), indexObject.Integer), false
		}
		return source.Content[index], true
	} else if indexObject.IsTypeById(TupleId) {
		if len(indexObject.Content) != 2 {
			return p.NewInvalidNumberOfArgumentsError(context, len(indexObject.Content), 2), false
		}
		startIndex, calcError := tools.CalcIndex(indexObject.Content[0].Integer, len(source.Content))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Content), indexObject.Content[0].Integer), false
		}
		var targetIndex int
		targetIndex, calcError = tools.CalcIndex(indexObject.Content[1].Integer, len(source.Content))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Content), indexObject.Content[1].Integer), false
		}
		return p.NewArray(context, false, source.Content[startIndex:targetIndex]), true
	}
	return p.NewInvalidTypeError(context, indexObject.TypeName(), IntegerName, TupleName), false
}

func (p *Plasma) ContentToString(context *Context, source *Value) (*Value, bool) {
	var opening string
	var closing string
	switch source.BuiltInTypeId {
	case ArrayId:
		opening = "["
		closing = "]"
	case TupleId:
		opening = "("
		closing = ")"
	}
	var objectString *Value
	var success bool
	result := ""
	for index, contentObject := range source.Content {
		if index != 0 {
			result += ", "
		}
		objectToString, getError := contentObject.Get(p, context, ToString)
		if getError != nil {
			return getError, false
		}
		objectString, success = p.CallFunction(context, objectToString)
		if !success {
			return objectString, false
		}
		result += objectString.String
	}
	return p.NewString(context, false, opening+result+closing), true
}

func (p *Plasma) ContentAssign(context *Context, source, indexObject *Value, value *Value) (*Value, bool) {
	if !indexObject.IsTypeById(IntegerId) {
		return p.NewInvalidTypeError(context, indexObject.GetClass(p).Name, IntegerName), false
	}
	index, calcError := tools.CalcIndex(indexObject.Integer, len(source.Content))
	if calcError != nil {
		return p.NewIndexOutOfRange(context, len(source.Content), indexObject.Integer), false
	}
	source.Content[index] = value
	return p.GetNone(), true
}

func (p *Plasma) ContentIterator(context *Context, source *Value) (*Value, bool) {
	result := p.NewIterator(context, false)
	information := struct {
		index  int
		length int
	}{
		0,
		len(source.Content),
	}
	result.Set(
		p, context,
		Next,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					if information.index >= information.length {
						return p.NewIndexOutOfRange(context, information.length, int64(information.index)), false
					}
					r := source.Content[information.index]
					information.index++
					return r, true
				},
			),
		),
	)
	result.Set(
		p, context,
		HasNext,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					return p.InterpretAsBool(information.index < information.length), true
				},
			),
		),
	)
	return result, true
}

func (p *Plasma) BytesIterator(context *Context, source *Value) (*Value, bool) {
	result := p.NewIterator(context, false)
	information := struct {
		index  int
		length int
	}{
		0,
		len(source.Bytes),
	}
	result.Set(
		p, context,
		Next,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					if information.index >= information.length {
						return p.NewIndexOutOfRange(context, information.length, int64(information.index)), false
					}
					r := self.Bytes[information.index]
					information.index++
					return p.NewInteger(context, false, int64(r)), true
				},
			),
		),
	)
	result.Set(
		p, context,
		HasNext,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					return p.InterpretAsBool(information.index < information.length), true
				},
			),
		),
	)
	return result, true
}

func (p *Plasma) StringIterator(context *Context, source *Value) (*Value, bool) {
	result := p.NewIterator(context, false)
	asRune := []rune(source.String)
	information := struct {
		index  int
		length int
	}{
		0,
		len(asRune),
	}
	result.Set(
		p, context,
		Next,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					if information.index >= information.length {
						return p.NewIndexOutOfRange(context, information.length, int64(information.index)), false
					}
					r := asRune[information.index]
					information.index++
					return p.NewString(context, false, string(r)), true
				},
			),
		),
	)
	result.Set(
		p, context,
		HasNext,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					return p.InterpretAsBool(information.index < information.length), true
				},
			),
		),
	)
	return result, true
}

func (p *Plasma) InterpretAsBool(expression bool) *Value {
	if expression {
		return p.GetTrue()
	}
	return p.GetFalse()
}

func (p *Plasma) StringIndex(context *Context, source *Value, indexObject *Value) (*Value, bool) {
	if indexObject.IsTypeById(IntegerId) {
		index, getIndexError := tools.CalcIndex(indexObject.Integer, len(source.String))
		if getIndexError != nil {
			return p.NewIndexOutOfRange(context, len(source.String), indexObject.Integer), false
		}
		return p.NewString(context, false, string(source.String[index])), true
	} else if indexObject.IsTypeById(TupleId) {
		if len(indexObject.Content) != 2 {
			return p.NewInvalidNumberOfArgumentsError(context, len(indexObject.Content), 2), false
		}
		startIndex, calcError := tools.CalcIndex(indexObject.Content[0].Integer, len(source.String))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.String), indexObject.Content[0].Integer), false
		}
		var targetIndex int
		targetIndex, calcError = tools.CalcIndex(indexObject.Content[1].Integer, len(source.String))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.String), indexObject.Content[1].Integer), false
		}
		return p.NewString(context, false, source.String[startIndex:targetIndex]), true
	}
	return p.NewInvalidTypeError(context, indexObject.TypeName(), IntegerName, TupleName), false
}

func (p *Plasma) BytesIndex(context *Context, source *Value, indexObject *Value) (*Value, bool) {
	if indexObject.IsTypeById(IntegerId) {
		index, calcError := tools.CalcIndex(indexObject.Integer, len(source.Bytes))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Bytes), indexObject.Integer), false
		}
		return p.NewInteger(context, false, int64(source.Bytes[index])), true
	} else if indexObject.IsTypeById(TupleId) {
		if len(indexObject.Content) != 2 {
			return p.NewInvalidNumberOfArgumentsError(context, len(indexObject.Content), 2), false
		}
		startIndex, calcError := tools.CalcIndex(indexObject.Content[0].Integer, len(source.Bytes))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Bytes), indexObject.Content[0].Integer), false
		}
		var targetIndex int
		targetIndex, calcError = tools.CalcIndex(indexObject.Content[1].Integer, len(source.Bytes))
		if calcError != nil {
			return p.NewIndexOutOfRange(context, len(source.Bytes), indexObject.Content[1].Integer), false
		}
		return p.NewBytes(context, false, source.Bytes[startIndex:targetIndex]), true
	}
	return p.NewInvalidTypeError(context, indexObject.TypeName(), IntegerName, TupleName), false
}

func (p *Plasma) Hash(context *Context, value *Value) (*Value, bool) {
	objectHashFunc, getError := value.Get(p, context, Hash)
	if getError != nil {
		return getError, false
	}
	objectHash, success := p.CallFunction(context, objectHashFunc)
	if !success {
		return objectHash, false
	}
	if !objectHash.IsTypeById(IntegerId) {
		return p.NewInvalidTypeError(context, objectHash.TypeName(), IntegerName), false
	}
	return objectHash, true
}

func (p *Plasma) HashIndexAssign(context *Context, hash *Value, key *Value, value *Value) (*Value, bool) {
	indexHash, success := p.Hash(context, key)
	if !success {
		return indexHash, false
	}
	if !indexHash.IsTypeById(IntegerId) {
		return p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName), false
	}
	keyValues, found := hash.KeyValues[indexHash.Integer]
	if found {
		hash.AddKeyValue(indexHash.Integer, &KeyValue{
			Key:   key,
			Value: value,
		})
		return p.GetNone(), true
	}
	indexObjectEquals, getError := key.Get(p, context, Equals)
	if getError != nil {
		return getError, false
	}
	var equals *Value
	for index, keyValue := range keyValues {
		equals, success = p.CallFunction(context, indexObjectEquals, keyValue.Key)
		if !success {
			return equals, false
		}
		equalsBool, callError := p.QuickGetBool(context, equals)
		if callError != nil {
			return callError, false
		}
		if equalsBool {
			hash.KeyValues[indexHash.Integer][index].Value = value
			return p.GetNone(), true
		}
	}
	hash.KeyValues[indexHash.Integer] = append(
		hash.KeyValues[indexHash.Integer],
		&KeyValue{
			Key:   key,
			Value: value,
		},
	)
	return p.GetNone(), true
}

func (p *Plasma) HashEquals(context *Context, leftHandSide *Value, rightHandSide *Value) (*Value, bool) {
	if !leftHandSide.IsTypeById(HashTableId) && !rightHandSide.IsTypeById(HashTableId) {
		return p.NewInvalidTypeError(context, leftHandSide.Name, HashName), false
	} else if !leftHandSide.IsTypeById(HashTableId) {
		return p.GetFalse(), true
	} else if !rightHandSide.IsTypeById(HashTableId) {
		return p.GetFalse(), true
	}
	rightIndex, getError := rightHandSide.Get(p, context, Index)
	if getError != nil {
		return getError, false
	}
	for key, leftValue := range leftHandSide.KeyValues {
		// Check if other has the key
		rightValue, ok := rightHandSide.KeyValues[key]
		if !ok {
			return p.GetFalse(), true
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return p.GetFalse(), true
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, success := p.CallFunction(context, rightIndex, entry.Key)
			if !success {
				return p.GetFalse(), true
			}
		}
	}
	return p.GetTrue(), true
}

func (p *Plasma) HashNotEquals(context *Context, leftHandSide *Value, rightHandSide *Value) (*Value, bool) {
	if !leftHandSide.IsTypeById(HashTableId) && !rightHandSide.IsTypeById(HashTableId) {
		return p.NewInvalidTypeError(context, leftHandSide.Name, HashName), false
	} else if !leftHandSide.IsTypeById(HashTableId) {
		return p.GetTrue(), true
	} else if !rightHandSide.IsTypeById(HashTableId) {
		return p.GetTrue(), true
	}

	rightIndex, getError := rightHandSide.Get(p, context, Index)
	if getError != nil {
		return getError, false
	}
	for key, leftValue := range leftHandSide.KeyValues {
		// Check if other has the key
		rightValue, ok := rightHandSide.KeyValues[key]
		if !ok {
			return p.GetTrue(), true
		}
		// Check if the each entry one has the same length
		if len(leftValue) != len(rightValue) {
			return p.GetTrue(), true
		}
		// Start comparing the entries
		for _, entry := range leftValue {
			_, success := p.CallFunction(context, rightIndex, entry.Key)
			if !success {
				return p.GetTrue(), true
			}
		}
	}
	return p.GetFalse(), true
}

func (p *Plasma) HashContent(context *Context, value *Value) (*Value, bool) {
	tupleHash := XXPrime5 ^ p.Seed()
	for _, contentObject := range value.Content {
		objectHash, success := p.Hash(context, contentObject)
		if !success {
			return objectHash, false
		}
		tupleHash += uint64(objectHash.Integer) * XXPrime2
		tupleHash = (tupleHash << 31) | (tupleHash >> 33)
		tupleHash *= XXPrime1
		tupleHash &= (1 << 64) - 1
	}
	return p.NewInteger(context, false, int64(tupleHash)), true
}

func (p *Plasma) HashToContent(context *Context, source *Value, target uint8) (*Value, bool) {
	var keys []*Value
	for _, keyValues := range source.KeyValues {
		for _, keyValue := range keyValues {
			keys = append(keys, keyValue.Key)
		}
	}
	switch target {
	case ArrayId:
		return p.NewArray(context, false, keys), true
	case TupleId:
		return p.NewTuple(context, false, keys), true
	}
	return p.NewInvalidTypeError(context, "UNKNOWN", HashName), false
}

func (p *Plasma) HashIterator(context *Context, source *Value) (*Value, bool) {
	tuple, success := p.HashToContent(context, source, TupleId)
	if !success {
		return tuple, false
	}
	information := struct {
		index  int
		length int
	}{
		0,
		len(tuple.Content),
	}
	result := p.NewIterator(context, false)
	result.Set(
		p, context,
		Next,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					if information.index >= information.length {
						return p.NewIndexOutOfRange(context, information.length, int64(information.index)), false
					}
					r := tuple.Content[information.index]
					information.index++
					return r, true
				},
			),
		),
	)
	result.Set(
		p, context,
		HasNext,
		p.NewFunction(
			context,
			false,
			context.PeekSymbolTable(),
			NewBuiltInClassFunction(
				result,
				0,
				func(self *Value, _ ...*Value) (*Value, bool) {
					return p.InterpretAsBool(information.index < information.length), true
				},
			),
		),
	)
	return result, true
}

func (p *Plasma) HashIndex(context *Context, source *Value, indexObject *Value) (*Value, bool) {
	indexObjectHash, getError := indexObject.Get(p, context, Hash)
	if getError != nil {
		return getError, false
	}
	indexHash, success := p.CallFunction(context, indexObjectHash)
	if !success {
		return indexHash, false
	}
	if !indexHash.IsTypeById(IntegerId) {
		return p.NewInvalidTypeError(context, indexHash.TypeName(), IntegerName), false
	}
	keyValues, found := source.KeyValues[indexHash.Integer]
	if !found {
		return p.NewKeyNotFoundError(context, indexObject), false
	}
	var indexObjectEquals *Value
	indexObjectEquals, getError = indexObject.Get(p, context, Equals)
	if getError != nil {
		return getError, false
	}
	var equals *Value
	for _, keyValue := range keyValues {
		equals, success = p.CallFunction(context, indexObjectEquals, keyValue.Key)
		if !success {
			return equals, false
		}
		equalsBool, callError := p.QuickGetBool(context, equals)
		if callError != nil {
			return callError, false
		}
		if equalsBool {
			return keyValue.Value, true
		}
	}
	return p.NewKeyNotFoundError(context, indexObject), false
}

func (p *Plasma) HashContains(context *Context, source *Value, value *Value) (*Value, bool) {
	valueHashFunc, getError := value.Get(p, context, Hash)
	if getError != nil {
		return getError, false
	}
	valueHashObject, success := p.CallFunction(context, valueHashFunc)
	if !success {
		return valueHashObject, false
	}
	if !valueHashObject.IsTypeById(IntegerId) {
		return p.NewInvalidTypeError(context, valueHashObject.TypeName(), IntegerName), false
	}
	valueHash := valueHashObject.Integer
	entries, found := source.KeyValues[valueHash]
	if !found {
		return p.GetFalse(), true
	}
	var valueEquals *Value
	valueEquals, getError = value.Get(p, context, RightEquals)
	if getError != nil {
		return getError, false
	}
	var comparisonResult *Value
	for _, entry := range entries {
		comparisonResult, success = p.CallFunction(context, valueEquals, entry.Key)
		if !success {
			return comparisonResult, false
		}
		comparisonResultBool, callError := p.QuickGetBool(context, comparisonResult)
		if callError != nil {
			return callError, false
		}
		if comparisonResultBool {
			return p.GetTrue(), true
		}
	}
	return p.GetFalse(), true
}

func (p *Plasma) IndexCall(context *Context, source *Value, index *Value) (*Value, bool) {
	indexOperation, getError := source.Get(p, context, Index)
	if getError != nil {
		return getError, false
	}
	return p.CallFunction(context, indexOperation, index)
}

func (p *Plasma) StringToContent(context *Context, s *Value, target uint8) (*Value, bool) {
	var content []*Value
	for _, char := range []rune(s.String) {
		content = append(content, p.NewString(context, false, string(char)))
	}
	if target == ArrayId {
		return p.NewArray(context, false, content), true
	} else if target == TupleId {
		return p.NewTuple(context, false, content), true
	}
	panic("String to content only support ArrayId and TupleId")
}

func (p *Plasma) InterpretAsIterator(context *Context, value *Value) (*Value, bool) {
	_, foundNext := value.Get(p, context, Next)
	_, foundHasNext := value.Get(p, context, HasNext)
	if foundNext == nil && foundHasNext == nil {
		return value, true
	}
	iter, getError := value.Get(p, context, Iter)
	if getError != nil {
		return getError, false
	}
	asIter, success := p.CallFunction(context, iter)
	if !success {
		return asIter, false
	}
	return asIter, true
}

func (p *Plasma) BytesToContent(context *Context, s *Value, target uint8) (*Value, bool) {
	var newContent []*Value
	for _, byte_ := range s.Bytes {
		newContent = append(newContent,
			p.NewInteger(context, false,
				int64(byte_),
			),
		)
	}
	if target == ArrayId {
		return p.NewArray(context, false, newContent), true
	} else if target == TupleId {
		return p.NewTuple(context, false, newContent), true
	}
	panic("Bytes to content only support ArrayId and TupleId")
}

func (p *Plasma) IterToContent(context *Context, s *Value, target uint8) (*Value, bool) {
	next, nextGetError := s.Get(p, context, Next)
	if nextGetError != nil {
		return nextGetError, false
	}
	hasNext, hasNextGetError := s.Get(p, context, HasNext)
	if hasNextGetError != nil {
		return hasNextGetError, false
	}
	var content []*Value
	for {
		doesHasNext, success := p.CallFunction(context, hasNext)
		if !success {
			return doesHasNext, false
		}
		doesHasNextAsBool, interpretationError := p.QuickGetBool(context, doesHasNext)
		if interpretationError != nil {
			return interpretationError, false
		}
		if !doesHasNextAsBool {
			break
		}
		var nextValue *Value
		nextValue, success = p.CallFunction(context, next)
		if !success {
			return nextValue, false
		}
		content = append(content, nextValue)
	}
	if target == ArrayId {
		return p.NewArray(context, false, content), true
	} else if target == TupleId {
		return p.NewTuple(context, false, content), true
	}
	panic("Iter to content only support ArrayId and TupleId")
}

func (p *Plasma) UnpackValues(context *Context, source *Value, numberOfReceivers int) ([]*Value, *Value) {
	if numberOfReceivers <= 1 {
		return []*Value{source}, nil
	}
	switch source.BuiltInTypeId {
	case TupleId, ArrayId:
		return source.Content, nil
	case HashTableId:
		hashAsTuple, success := p.HashToContent(context, source, TupleId)
		if !success {
			return nil, hashAsTuple
		}
		return hashAsTuple.Content, nil
	case StringId:
		stringAsTuple, success := p.StringToContent(context, source, TupleId)
		if !success {
			return nil, stringAsTuple
		}
		return stringAsTuple.Content, nil
	case BytesId:
		bytesAsTuple, success := p.BytesToContent(context, source, TupleId)
		if !success {
			return nil, bytesAsTuple
		}
		return bytesAsTuple.Content, nil
	case IteratorId:
		iterAsTuple, success := p.IterToContent(context, source, TupleId)
		if !success {
			return nil, iterAsTuple
		}
		return iterAsTuple.Content, nil
	}
	// Transform the type to iter
	asIterInterpretation, success := p.InterpretAsIterator(context, source)
	if !success {
		return nil, asIterInterpretation
	}
	var sourceAsIter *Value
	sourceAsIter, success = p.CallFunction(context, asIterInterpretation)
	if !success {
		return nil, sourceAsIter
	}
	// The to Tuple
	var sourceIterAsTuple *Value
	sourceIterAsTuple, success = p.IterToContent(context, sourceAsIter, TupleId)
	if !success {
		return nil, sourceIterAsTuple
	}
	// Return its content
	return sourceIterAsTuple.Content, nil
}
