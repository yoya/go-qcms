package qcms

/*
#cgo CFLAGS: -Iqcms
#include <stdlib.h>
#include "qcms.h"
*/
import "C"

import (
	"unsafe"
)

type Profile struct {
	prof *C.qcms_profile
}

func OpenProfileFromFile(filename string) *Profile {
	csfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(csfilename))
	csmode := C.CString("r")
	defer C.free(unsafe.Pointer(csmode))
	return &Profile{prof: C.qcms_profile_from_path(csfilename)}
}

func OpenProfileFromMem(profdata []byte) *Profile {
	data := unsafe.Pointer(&profdata[0])
	dataLen := C.size_t(len(profdata))
	return &Profile{prof: C.qcms_profile_from_memory(data, dataLen)}
}

func Create_sRGBProfile() *Profile {
	return &Profile{prof: C.qcms_profile_sRGB()}
}

func (prof *Profile) CloseProfile() {
	if prof.prof != nil {
		C.qcms_profile_release(prof.prof)
	}
}
