function "remove_ext" {
  params = [file]
  result = replace(basename(file), ext(file), "")
}

function "white" {
  params = [file]
  result = color("white", file)
}
