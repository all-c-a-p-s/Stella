package main

import "fmt"

func evaluate(board [9]string, side string) string {
  var evaluation string  = "Unknown"
  var i int  = 0
  var draw bool  = true
  for i < 9 {
    if board[i] == "_" {
      draw = false

    }
    i = i + 1

  }
  if draw {
    evaluation = "Draw"

  }
  if ( ( board[0] == board[1] ) && ( board[1] == board[2] ) ) && ( board[0] != "_" ) {
    if board[0] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  } else if ( ( board[3] == board[4] ) && ( board[4] == board[5] ) ) && ( board[3] != "_" ) {
    if board[3] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  } else if ( ( board[6] == board[7] ) && ( board[7] == board[8] ) ) && ( board[6] != "_" ) {
    if board[6] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  }
  if ( ( board[0] == board[3] ) && ( board[3] == board[6] ) ) && ( board[0] != "_" ) {
    if board[0] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  } else if ( ( board[1] == board[4] ) && ( board[4] == board[7] ) ) && ( board[1] != "_" ) {
    if board[1] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  } else if ( ( board[2] == board[5] ) && ( board[5] == board[8] ) ) && ( board[2] != "_" ) {
    if board[2] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  }
  if ( ( board[0] == board[4] ) && ( board[4] == board[8] ) ) && ( board[0] != "_" ) {
    if board[0] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  } else if ( ( board[2] == board[4] ) && ( board[4] == board[6] ) ) && ( board[2] != "_" ) {
    if board[2] == side {
      evaluation = "Win"

    } else  {
      evaluation = "Loss"

    }

  }
  return evaluation

}
func opposite_evaluation(evaluation string) string {
  var res string  = "Unknown"
  if evaluation == "Win" {
    res = "Loss"

  } else if evaluation == "Loss" {
    res = "Win"

  } else if evaluation == "Draw" {
    res = "Draw"

  }
  return res

}
func opposite_side(side string) string {
  var res string  = "O"
  if side == "O" {
    res = "X"

  }
  return res

}
func minimax(board [9]string, side string) string {
  var res string  = "Loss"
  var evaluation string  = evaluate(board, side)
  if evaluation == "Unknown" {
    var square int  = 0
    for square < 9 {
      if board[square] == "_" {
        var copy [9]string = board
        copy[square] = side
        var opponent string  = opposite_side(side)
        var opponent_perspective string  = minimax(copy, opponent)
        var conditional_evaluation string  = opposite_evaluation(opponent_perspective)
        if conditional_evaluation == "Win" {
          res = "Win"
          break

        } else if conditional_evaluation == "Draw" {
          res = "Draw"

        }

      }
      square = square + 1

    }

  } else  {
    res = evaluation

  }
  return res

}
func main() {
  var board [9]string = [9]string{"X", "X", "_", "_", "_", "_", "_", "_", "_"}
  var future_board [9]string = [9]string{"X", "X", "_", "_", "_", "_", "_", "_", "_"}
  var side string  = "X"
  var square int  = 0
  var found_a_move bool  = false
  for square < 9 {
    if board[square] == "_" {
      var copy [9]string = board
      copy[square] = side
      var opponent string  = opposite_side(side)
      var opponent_perspective string  = minimax(copy, opponent)
      var conditional_evaluation string  = opposite_evaluation(opponent_perspective)
      if conditional_evaluation == "Win" {
        var i int  = 0
        for i < 9 {
          future_board[i] = copy[i]
          i = i + 1

        }
        break

      } else if conditional_evaluation == "Draw" {
        var i int  = 0
        for i < 9 {
          future_board[i] = copy[i]
          i = i + 1

        }

      } else if ! found_a_move {
        var i int  = 0
        for i < 9 {
          future_board[i] = copy[i]
          i = i + 1

        }
        found_a_move = true

      }

    }
    square = square + 1

  }
  var row int  = 0
  for row < 3 {
    var col int  = 0
    for col < 3 {
      var i int  = row * 3 + col
      fmt.Print( future_board[i] )
      col = col + 1

    }
    fmt.Print( "\n" )
    row = row + 1

  }

}


