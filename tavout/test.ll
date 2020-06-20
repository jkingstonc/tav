declare i32 @printf(i8* %formatter, i32 %value)

declare i32 @puts(i8* %string)

define void @main() {
main_body:
	%0 = alloca [13 x i8]
	store [13 x i8] c"hello world\00\0A", [13 x i8]* %0
	%1 = bitcast [13 x i8]* %0 to i8*
	%2 = call i32 @puts(i8* %1)
	ret void
}
