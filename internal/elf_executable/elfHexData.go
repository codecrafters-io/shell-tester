package elf_executable

func getELFHexData() (string, string) {
	elfDataPart1 :=
		`# ELF header -- always 64 bytes

	# Magic number
	7f 45 4c 46 

	# 02 = 64 bit
	02

	# 01 = little-endian
	01

	# 01 = ELF version
	01

	# 00 = Target ABI -- usually left at zero for static executables
	00

	# 00 = Target ABI version -- usually left at zero for static executables
	00

	# 7 bytes undefined
	00 00 00 00 00 00 00 

	# 02 = executable binary
	02 00 

	# 3E = amd64 architecture
	3e 00

	# 1 = ELF version
	01 00 00 00 

	# 0x400078 = start address: right after this header and the program
	#  header, which take 0x78 bytes, if the binary is mapped into 
	#  memory at address 0x400000)
	78 00 40 00 00 00 00 00

	# 0x40 = offset to program header, right after this header which is 0x40 bytes long 
	40 00 00 00 00 00 00 00

	# 0xC0 = offset to section header, which is after the program text and the 
	#  string table
	c0 00 00 00 00 00 00 00 

	# 0x00000000 = architecture specific flags
	00 00 00 00

	# 0x40 = size of this header, always 64 bytes
	40 00

	# 0x38 = size of a program header, always 56 bytes
	38 00

	# 1 = number of program header
	01 00

	# 0x40 = size of a section header, always 64 bytes
	40 00

	# 3 = number of sections headers 
	03 00

	# 2 = index of the section header that references the stringtable
	02 00 

	######################################################################
	# Program header -- always 56 bytes

	# 1 = type of entry: loadable segment
	01 00 00 00 

	# 0x05 = segment-dependent flags: executable | readable
	05 00 00 00 

	# 0x0 = Offset within file
	00 00 00 00 00 00 00 00 

	# 0x400000 = load position in virtual memory
	00 00 40 00 00 00 00 00 

	# 0x400000 = load position in physical memory
	00 00 40 00 00 00 00 00

	# 0xB0 = size of the loaded section the file
	B0 00 00 00 00 00 00 00

	# 0xB0 = size of the loaded section in memory 
	B0 00 00 00 00 00 00 00

	# 0x200000 = alignment boundary for sections
	00 00 20 00 00 00 00 00

	######################################################################
	# Program code -- 42 bytes

	# mov 0x01,%rax -- sys_write
	48 c7 c0 01 00 00 00 

	# mov 0x01,%rdi -- file descriptor, stdout
	48 c7 c7 01 00 00 00

	# mov 0x4000a2,%rsi -- location of string
	48 c7 c6 a2 00 40 00

	# mov 0x0d,%rdx -- size of string, 13 bytes
	48 c7 c2 0d 00 00 00

	# syscall
	0f 05

	# mov 0x3c,$rax -- exit program
	48 c7 c0 3c 00 00 00 

	# xor %rdi,%rdi -- exit code, 0
	48 31 ff 

	# syscall
	0f 05 

	# Text "Hello, world\n" -- total 13 bytes including the null
	# 48 65 6c 6c 6f 2c 20 77 6f 72 6c 64 0a
	`

	elfDataPart2 :=
		`# String table. ".shstrab\0.text\0" -- 16 bytes

	2e 73 68 73 74 72 74 61 62 00 2e 74 65 78 74 00 

	# Section header table, section 0 -- unused -- 64 bytes

	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 
	00 00 00 00 00 00 00 00 

	# Section header table, section 1 -- program text -- 64 bytes

	# 0x0A = offset to the name of the section in the stringtable
	0a 00 00 00 

	# 1 = type: program data
	01 00 00 00

	# 0x06 flags = executable | occupies memory
	06 00 00 00 00 00 00 00 

	# 0x400078 address in virtual memory of this section
	78 00 40 00 00 00 00 00

	# 0x78 = offset in the file of this section (start of program code)
	78 00 00 00 00 00 00 00

	# 0x38 = size of this section in the file: 56 bytes
	38 00 00 00 00 00 00 00 

	# sh_link -- not used for this section
	00 00 00 00 00 00 00 00

	# 0x01 = alignment code: default??
	01 00 00 00 00 00 00 00

	# sh_entsize: not used
	00 00 00 00 00 00 00 00

	# Section header table, section 2 -- stringtable

	# 0x0 = offset to the name of the section in the stringtable
	00 00 00 00

	# 3 = type: string table 
	03 00 00 00 

	# 0x0 = flags
	00 00 00 00 00 00 00 00 

	# 0x0 = address in virtual memory (not used)
	00 00 00 00 00 00 00 00 

	# 0xB0 = offset in the file of this section (start of string table)
	b0 00 00 00 00 00 00 00

	# 0x10 = size of this section in the file: 16 bytes
	10 00 00 00 00 00 00 00

	# sh_link -- not used for this section
	00 00 00 00 00 00 00 00

	# 0x01 = alignment code: default??
	01 00 00 00 00 00 00 00

	# sh_entsize: not used
	00 00 00 00 00 00 00 00`

	return elfDataPart1, elfDataPart2
}
