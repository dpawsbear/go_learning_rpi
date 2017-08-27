package main

import (
	"os"
	"fmt"
	"bytes"
	"encoding/binary"
	"syscall"
	"unsafe"
	"time"
)

const(

	// define the device tree range
	BCM2837_PRI3B_DT_FILENAME                 = "/proc/device-tree/soc/ranges"
	BCM2837_PRI3_DT_PERI_BASE_ADDRESS_OFFSET  =  0x4
	BCM2837_PRI3_DT_PERI_SIZE_OFFSET          =  0X8
	BCM2837_GPIO_BASE                         =  0x00200000

	/*! GPIO register offsets from BCM2835_GPIO_BASE.
	Offsets into the GPIO Peripheral block in bytes per 6.1 Register View
	*/
	BCM2837_GPFSEL0                 =    0x0000 /*!< GPIO Function Select 0 */
	BCM2837_GPFSEL1                 =    0x0004 /*!< GPIO Function Select 1 */
	BCM2837_GPFSEL2                 =    0x0008 /*!< GPIO Function Select 2 */
	BCM2837_GPFSEL3                 =    0x000c /*!< GPIO Function Select 3 */
	BCM2837_GPFSEL4                 =    0x0010 /*!< GPIO Function Select 4 */
	BCM2837_GPFSEL5                 =    0x0014 /*!< GPIO Function Select 5 */
	BCM2837_GPSET0                  =    0x001c /*!< GPIO Pin Output Set 0 */
	BCM2837_GPSET1                  =    0x0020 /*!< GPIO Pin Output Set 1 */
	BCM2837_GPCLR0                  =    0x0028 /*!< GPIO Pin Output Clear 0 */
	BCM2837_GPCLR1                  =    0x002c /*!< GPIO Pin Output Clear 1 */
	BCM2837_GPLEV0                  =    0x0034 /*!< GPIO Pin Level 0 */
	BCM2837_GPLEV1                  =    0x0038 /*!< GPIO Pin Level 1 */
	BCM2837_GPEDS0                  =    0x0040 /*!< GPIO Pin Event Detect Status 0 */
	BCM2837_GPEDS1                  =    0x0044 /*!< GPIO Pin Event Detect Status 1 */
	BCM2837_GPREN0                  =    0x004c /*!< GPIO Pin Rising Edge Detect Enable 0 */
	BCM2837_GPREN1                  =    0x0050 /*!< GPIO Pin Rising Edge Detect Enable 1 */
	BCM2837_GPFEN0                  =    0x0058 /*!< GPIO Pin Falling Edge Detect Enable 0 */
	BCM2837_GPFEN1                  =    0x005c /*!< GPIO Pin Falling Edge Detect Enable 1 */
	BCM2837_GPHEN0                  =    0x0064 /*!< GPIO Pin High Detect Enable 0 */
	BCM2837_GPHEN1                  =    0x0068 /*!< GPIO Pin High Detect Enable 1 */
	BCM2837_GPLEN0                  =    0x0070 /*!< GPIO Pin Low Detect Enable 0 */
	BCM2837_GPLEN1                  =    0x0074 /*!< GPIO Pin Low Detect Enable 1 */
	BCM2837_GPAREN0                 =    0x007c /*!< GPIO Pin Async. Rising Edge Detect 0 */
	BCM2837_GPAREN1                 =    0x0080 /*!< GPIO Pin Async. Rising Edge Detect 1 */
	BCM2837_GPAFEN0                 =    0x0088 /*!< GPIO Pin Async. Falling Edge Detect 0 */
	BCM2837_GPAFEN1                 =    0x008c /*!< GPIO Pin Async. Falling Edge Detect 1 */
	BCM2837_GPPUD                   =    0x0094 /*!< GPIO Pin Pull-up/down Enable */
	BCM2837_GPPUDCLK0               =    0x0098 /*!< GPIO Pin Pull-up/down Enable Clock 0 */
	BCM2837_GPPUDCLK1               =    0x009c /*!< GPIO Pin Pull-up/down Enable Clock 1 */
)

var  Bcm2837_peripherals_base   uint32
var  Bcm2837_peripherals_size   uint32
var  Bcm2837_gpio               uint32

func main(){
	// find the io peripheral base and range
	f,err:= os.OpenFile(BCM2837_PRI3B_DT_FILENAME,os.O_RDONLY,0)

	if err != nil {
		fmt.Println("open range file err")
	}
	defer f.Close()

	var buf []byte = make([]byte,4)
	f.ReadAt(buf , BCM2837_PRI3_DT_PERI_BASE_ADDRESS_OFFSET )
	bytesBuffer := bytes.NewBuffer(buf)
	binary.Read(bytesBuffer , binary.BigEndian , &Bcm2837_peripherals_base )

	f.ReadAt(buf , BCM2837_PRI3_DT_PERI_SIZE_OFFSET )
	bytesBuffer  = bytes.NewBuffer(buf)
	binary.Read(bytesBuffer , binary.BigEndian , &Bcm2837_peripherals_size )

	fmt.Printf("get peripherals base:%x size:%x\n" , Bcm2837_peripherals_base , Bcm2837_peripherals_size )

	//need su execute
	if os.Geteuid() == 0 {

		/* open the master /dev/mem device */
		f, err := os.OpenFile("/dev/mem",os.O_RDWR,0)
		if err != nil {
			fmt.Println("Open mem error")
		}

		p,err := syscall.Mmap(int(f.Fd()),int64(Bcm2837_peripherals_base),int(Bcm2837_peripherals_size),syscall.PROT_READ|syscall.PROT_WRITE,syscall.MAP_SHARED)

		if err != nil {
			fmt.Println("mmap error")
		}

		//strat find the gpio register
		Bcm2837_gpio = *(*uint32)(unsafe.Pointer(&p)) + uint32( BCM2837_GPIO_BASE )

		var test uintptr = uintptr( Bcm2837_gpio + BCM2837_GPFSEL1 )
		*(*uint32)(unsafe.Pointer(test)) = ( 0x1 << 15 )

		test = uintptr( Bcm2837_gpio + BCM2837_GPSET0 )
		*(*uint32)(unsafe.Pointer(test)) = ( 0x1 << 15 )

		time.Sleep( time.Second  * 2 )

		test = uintptr( Bcm2837_gpio + BCM2837_GPCLR0 )
		*(*uint32)(unsafe.Pointer(test)) = ( 0x1 << 15 )

	}else {
		fmt.Println("please use root execute")
		panic(err)
	}
}
