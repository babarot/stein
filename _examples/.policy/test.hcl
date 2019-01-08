rule "test" {
  description = "Test"

  expressions = [true]

  report {
    level   = "ERROR"
    message = "${color("red", "red")}"
  }

  debug = []
}
