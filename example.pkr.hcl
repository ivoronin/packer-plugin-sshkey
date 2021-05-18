data "sshkey" "example" {
}

source "null" "example" {
  communicator = "none"
}

build {
  sources = ["sources.null.example"]
  provisioner "shell-local" {
    inline = [
      "echo PUBLIC KEY:",
      "echo ${data.sshkey.example.public_key}",
      "echo PRIVATE KEY \\(${data.sshkey.example.private_key_path}\\):",
      "cat \"${data.sshkey.example.private_key_path}\"",
    ]
  }
}
