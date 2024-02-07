function divisible(n: int, divisor: int, previous: int) -> bool = {
  // lowest integer factor must always necessarily be a prime
  let mut to_return: bool = false
  if (previous + divisor) > n {
    to_return = false
  } else if (previous + divisor) == n {
    to_return = true
  } else {
    to_return = divisible(n, divisor, previous + divisor)
    // recursive call with next integer
  }
  to_return
}

function next_divisor(n: int, divisor: int) -> int = {
  let mut to_return: int = 0
  if divisible(n, divisor, 0) {
    to_return = divisor
  } else {
    to_return = next_divisor(n, divisor+1)
  }
  to_return
}

function divide(n: int) -> int = {
  let mut to_return: int = 0
  if n == 1 {
    println!(n)
  } else {
    let d: int = next_divisor(n, 2)
    println!(d)
    let divided: int = n / d
    let res: int = divide(divided)
    to_return = res
  }
  to_return
}

function main() -> IO = {
  let num_to_factorise: int = 228944
  let res: int = divide(num_to_factorise)
  if res == -1 {
    print!(res)
  }
}
