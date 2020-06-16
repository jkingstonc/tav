define i32 @main() {
0:
	%1 = alloca i32
	store i32 1, i32* %1
	%2 = alloca i32
	store i32 2, i32* %2
	ret i32 99
}
