package main

/* checks if a given error has occurred */
func isError(e error){
  if e != nil {
    panic(e)
  }
}

// func Max(x int, y int) int {
//   if x > y {
//     return x
//   }
//   return y
// }

// func Min(x int, y int) int {
//   if y < x {
//     return y
//   }
//   return x
// }
