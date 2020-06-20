# Tav

Tav is a statically typed, compiled language. It is heavily inspired by C, Go and Jai. The goal is to develop a language that is powerful, yet easily readable.  Keywords are short, any variable's type can be infered at compile time and the grammar is natural.

	main : fn i32{
		x := 123;
		y i32 = x * 2;
		
		s : string = "hello world";
		puts(s);
		
		ret y;
	}

### Compile time function execution:
As Tav use LLVM ir, functions can be JIT compiled, allowing for any function to be compiled at compile time.
In C, a comppile time square function would look like this:
	
	#define SQUARE (x) (x*x)
	int main(void)
	{
	       int x = SQUARE(2);
	       return 0;
	}
In Tav, the code would look like this:
	
	square : fn i32 (val : i32) {
		ret val*val;
	}
	main : fn i32 {
		x := #run square(2);
		ret 0;
	} 
This allows us to run functions at compile time, and at runtime as the compiler doesn't see any difference.
