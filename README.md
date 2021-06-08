# Packer SSH key plugin

Packer plugin used to generate SSH keys.

For the full list of available features for this plugin see [documentation](https://www.packer.io/docs/datasources/sshkey).

Packer 1.7.3 or later is required.

## Usage example
```hcl
packer {
  required_plugins {
    sshkey = {
      version = ">= 0.1.0"
      source = "github.com/ivoronin/sshkey"
    }
  }
}

data "sshkey" "install" {
}

source "qemu" "install" {
  ssh_username              = "root"
  ssh_private_key_file      = data.sshkey.install.private_key_path
  ssh_clear_authorized_keys = true
  http_content = {
    "/preseed.cfg" = templatefile("preseed.cfg.pkrtpl", {
        "ssh_public_key" : data.sshkey.install.public_key
    })
  }
  <...>
}

build {
  sources = ["source.qemu.install"]
}
```

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.
