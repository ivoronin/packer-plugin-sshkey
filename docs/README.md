The SSHkey plugin can be used for generating SSH keys for configuring private key authentication

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    sshkey = {
      source  = "github.com/ivoronin/sshkey"
      version = "~> 1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/ivoronin/sshkey
```

### Components

#### Builders

- [sshkey](/packer/integrations/ivoronin/sshkey/latest/components/data-source/sshkey) - Data source used to generate SSH keys

### Example Usage

```hcl
packer {
  required_plugins {
    sshkey = {
      version = "~> 1"
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
