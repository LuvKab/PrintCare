//go:build windows

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/image/bmp"
)

var (
	winspool              = syscall.NewLazyDLL("winspool.drv")
	gdi32                 = syscall.NewLazyDLL("gdi32.dll")
	procEnumPrintersW     = winspool.NewProc("EnumPrintersW")
	procCreateDCW         = gdi32.NewProc("CreateDCW")
	procDeleteDC          = gdi32.NewProc("DeleteDC")
	procStartDocW         = gdi32.NewProc("StartDocW")
	procEndDoc            = gdi32.NewProc("EndDoc")
	procStartPage         = gdi32.NewProc("StartPage")
	procEndPage           = gdi32.NewProc("EndPage")
	procSetStretchBltMode = gdi32.NewProc("SetStretchBltMode")
	procStretchDIBits     = gdi32.NewProc("StretchDIBits")
	procGetDeviceCaps     = gdi32.NewProc("GetDeviceCaps")
)

const (
	printerEnumLocal      = 0x00000002
	printerEnumConnection = 0x00000004
	horzres               = 8
	vertres               = 10
	coloroncolor          = 3
	srcCopy               = 0x00CC0020
	dibRGBColors          = 0
)

type printerInfo2 struct {
	ServerName         *uint16
	PrinterName        *uint16
	ShareName          *uint16
	PortName           *uint16
	DriverName         *uint16
	Comment            *uint16
	Location           *uint16
	DevMode            uintptr
	SepFile            *uint16
	PrintProcessor     *uint16
	Datatype           *uint16
	Parameters         *uint16
	SecurityDescriptor uintptr
	Attributes         uint32
	Priority           uint32
	DefaultPriority    uint32
	StartTime          uint32
	UntilTime          uint32
	Status             uint32
	Jobs               uint32
	AveragePPM         uint32
}

type docInfoW struct {
	Size       int32
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
	Type       uint32
}

type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

func EnumPrinters() ([]string, error) {
	var needed, returned uint32
	flags := uint32(printerEnumLocal | printerEnumConnection)

	procEnumPrintersW.Call(uintptr(flags), 0, 2, 0, 0,
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)))

	if needed == 0 {
		return nil, nil
	}

	buf := make([]byte, needed)
	ret, _, err := procEnumPrintersW.Call(uintptr(flags), 0, 2,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)))
	if ret == 0 {
		return nil, fmt.Errorf("EnumPrinters: %v", err)
	}

	printers := make([]string, 0, returned)
	sz := unsafe.Sizeof(printerInfo2{})
	for i := uint32(0); i < returned; i++ {
		info := (*printerInfo2)(unsafe.Pointer(&buf[uintptr(i)*sz]))
		if info.PrinterName != nil {
			name := syscall.UTF16ToString(unsafe.Slice(info.PrinterName, 1024))
			printers = append(printers, name)
		}
	}
	return printers, nil
}

func PrintImage(printerName, imagePath string, paperSource int) error {
	img, err := loadImage(imagePath)
	if err != nil {
		return fmt.Errorf("load image: %v", err)
	}

	printerW, _ := syscall.UTF16PtrFromString(printerName)
	driverW, _ := syscall.UTF16PtrFromString("WINSPOOL")

	hdc, _, err := procCreateDCW.Call(
		uintptr(unsafe.Pointer(driverW)),
		uintptr(unsafe.Pointer(printerW)), 0, 0)
	if hdc == 0 {
		return fmt.Errorf("CreateDC: %v", err)
	}
	defer procDeleteDC.Call(hdc)

	docNameW, _ := syscall.UTF16PtrFromString("Printhead Maintainer")
	di := docInfoW{
		Size:    int32(unsafe.Sizeof(docInfoW{})),
		DocName: docNameW,
	}

	ret, _, err := procStartDocW.Call(hdc, uintptr(unsafe.Pointer(&di)))
	if int32(ret) <= 0 {
		return fmt.Errorf("StartDoc: %v", err)
	}

	ret, _, err = procStartPage.Call(hdc)
	if int32(ret) <= 0 {
		procEndDoc.Call(hdc)
		return fmt.Errorf("StartPage: %v", err)
	}

	pageW, _, _ := procGetDeviceCaps.Call(hdc, uintptr(horzres))
	pageH, _, _ := procGetDeviceCaps.Call(hdc, uintptr(vertres))
	procSetStretchBltMode.Call(hdc, uintptr(coloroncolor))

	bounds := img.Bounds()
	bmi, bits := imageToDIB(img)

	procStretchDIBits.Call(hdc,
		0, 0, pageW, pageH,
		0, 0, uintptr(bounds.Dx()), uintptr(bounds.Dy()),
		uintptr(unsafe.Pointer(&bits[0])),
		uintptr(unsafe.Pointer(bmi)),
		uintptr(dibRGBColors), uintptr(srcCopy))

	procEndPage.Call(hdc)
	procEndDoc.Call(hdc)
	return nil
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.HasSuffix(strings.ToLower(path), ".bmp") {
		return bmp.Decode(f)
	}
	img, _, err := image.Decode(f)
	return img, err
}

func imageToDIB(img image.Image) (*bitmapInfoHeader, []byte) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	stride := ((w*3 + 3) / 4) * 4
	bits := make([]byte, stride*h)

	for y := 0; y < h; y++ {
		row := (h - 1 - y) * stride
		for x := 0; x < w; x++ {
			r, g, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			off := row + x*3
			bits[off] = byte(bl >> 8)
			bits[off+1] = byte(g >> 8)
			bits[off+2] = byte(r >> 8)
		}
	}

	bmi := &bitmapInfoHeader{
		Size:     uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		Width:    int32(w),
		Height:   int32(h),
		Planes:   1,
		BitCount: 24,
	}
	return bmi, bits
}
