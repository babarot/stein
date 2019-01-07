rule "test" {
  description = "Test"

  # ignore_cases = ["hoge"]

  expressions = [true]
  report {
    level   = "ERROR"
    message = "${color("red", "red")}"
  }
  debug = []

  # "${join(", ", maplist(var.namespace_name_map, jsonpath(".metadata.namespace")))}",

  # "${join(", ", list(1, 2, ["a"]))}",
}

rule "test2" {
  description = "Test"

  # ignore_cases = ["hoge"]

  expressions = [true]
  report {
    level   = "ERROR"
    message = "${color("red", "red")}"
  }

  # debug = ["${maphoge("hoge", "hogehogehoge")}"]
  # debug = ["${color("white", "Test")}"]

  # "${join(", ", maplist(var.namespace_name_map, jsonpath(".metadata.namespace")))}",

  # "${join(", ", list(1, 2, ["a"]))}",
}
