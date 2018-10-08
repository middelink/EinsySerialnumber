# EinsySerialNumber

## TL;DR

This tool allows one to show or set the serial number in the Prusa Edition
of the 32u2 bootloader. It requires access to a ISP (In System Programming)
tool like the AVRISP, USBISP, USBASP or many others. Also, having access to
the [Go compiler](https://golang.org/dl/) helps :)

## Why?

I have multiple Prusa and Prusa clone machines, each having their own Einsy
logic board. All expect the official Prusa machine have no serial number in
their USB descriptors, making it imposible to map the machines consistently
to a usb tty device. E.g. 5 machines having
`/dev/serial/by-id/usb-Prusa_Research__prusa3d.com__Original_Prusa_i3_MK3____________________-if00`
does not make finding the right one any easier.

## Command Line Flags

Usage: EinsySerialNumber [-create] [-filename filename] [SerialNumber]

Usage of EinsySerialnumber:
  -create
    	Create the file if it does not exist (default true)
  -filename string
    	Name fo the eeprom image we insert the serial number in (default "32u2.eep")

## Building the binary

* Clone, download, copy/paste the source files onto your local disk.
* Execute `go build .` to create the EinsySerialNumber binary.
* (optional) Copy the binary to /usr/local/sbin.

## How to use it

Use [avrdude](https://www.nongnu.org/avrdude/) to extract the current
content of the 32u2 EEPROM:

```
sudo avrdude -C /usr/local/etc/avrdude.conf -v -p atmega32u2 -c usbasp -P usb -U eeprom:r:32u2.eep:r
```

If you get an error message, stop here and figure out what is wrong. If you
cannot read the 32u2's EEPROM, there is no use continuing.

Run the EinsySerialnumber tool without arguments, it will list the current
serial number. If you are using a China clone or a Einsy board directly from
Ultimachine, there is likely no serial number shown.

```
$ ./EinsySerialnumber 
Serial number: CZPX****X004X******
```

Now we're ready to setup our own serial number, I prefer to use the machine
name, the date I build it and a serial number all seperated by X'es. But it
is totally up to you. It is advisable to use 19 characters though.

```
./EinsySerialnumber NOTOX20180604X00001
Serial number: NOTOX20180604X00001
```

Now write the modified EEPROM image back to the 32u2:

```
sudo avrdude -C /usr/local/etc/avrdude.conf -v -p atmega32u2 -c usbasp -P usb -U eeprom:w:32u2.eep:r
```

## Technical stuff

For those interested, the serial number seems to be placed in the EEPROM,
address 0x100 till 0x100+19, so only 20 valid charcters.

For record keeping, I will also add the binary image (32u2.bin) of the flash
to the repository so people can reflash the 32u2 if their version is not
identical as to the Prusa one. Use at your own risk as it is a binary image
I do not have the source of.

## Small rant

Sadly I cannot find the source of the 32u2 USB\>serial converter, as such I can
only speculate it has some special commands to read/write the serial number
without having to jump to all of the above hoops. Hopefully one day it will be
released as opensource as well so we can verify this assumption.
