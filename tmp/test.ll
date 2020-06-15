define i32 @main() {
0:
	%1 = alloca i32
	store i32 123, i32* %1
	%2 = load i32, i32* %1
	%3 = alloca i32
	store i32 %2, i32* %3
	%4 = load i32, i32* %3
	ret i32 %4
}
