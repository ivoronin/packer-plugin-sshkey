# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "SSH Key"
  description = "The SSHkey plugin can be used for generating SSH keys for configuring private key authentication."
  identifier = "packer/ivoronin/sshkey"
  component {
    type = "data-source"
    name = "SSH Key"
    slug = "sshkey"
  }
}
