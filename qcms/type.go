package qcms

/*
#cgo CFLAGS: -Iqcms
#include "qcms.h"
*/
import "C"

type DataType C.qcms_data_type

const (
	DATA_RGB_8  DataType = C.QCMS_DATA_RGB_8
	DATA_RGBA_8 DataType = C.QCMS_DATA_RGBA_8
)
