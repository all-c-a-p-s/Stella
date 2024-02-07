function first(nums: int[5]) -> int = {
  nums[0]
}

function main() -> IO = {
  let mut test: int[5] = [1, 2, 3, 4, 5]
  println!(first(test))
}
