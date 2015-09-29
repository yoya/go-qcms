package qcms

/*
#cgo CFLAGS: -Iqcms
#include <stdlib.h>
#include "qcms.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Transform struct {
	trans *C.qcms_transform
}

func CreateTransform(src_prof *Profile, src_type DataType, dst_prof *Profile, dst_type DataType) *Transform {
	transform := C.qcms_transform_create(
		src_prof.prof, C.qcms_data_type(src_type),
		dst_prof.prof, C.qcms_data_type(dst_type),
		C.QCMS_INTENT_PERCEPTUAL)
	return &Transform{trans: transform}
}

func (trans *Transform) DeleteTransform() {
	if trans.trans != nil {
		C.qcms_transform_release(trans.trans)
	}
}

func (trans *Transform) DoTransform(inputBuffer []uint8, outputBuffer []uint8, length int) error {
	inputLen := len(inputBuffer)
	outputLen := len(outputBuffer)
	if inputLen < length {
		return fmt.Errorf("DoTransform: inputLen(%d) < length(%d)", inputLen, length)
	}
	if outputLen < length {
		return fmt.Errorf("DoTransform: outputLen(%d) < length(%d)", outputLen, length)
	}
	inputPtr := unsafe.Pointer(&inputBuffer[0])
	outputPtr := unsafe.Pointer(&outputBuffer[0])
	length /= 4 // XXX?
	C.qcms_transform_data(trans.trans, inputPtr, outputPtr, C.size_t(length))
	return nil
}
