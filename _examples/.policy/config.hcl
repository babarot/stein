config {
  report {
    format = "${format("[{{.Level}}]  {{.Rule}}  %s", color("{{.Message}}", "white"))}"
    style  = "console"
    color  = true
  }
}
