# lanyard

_lanyard_ is a tool to manage `ssh-agent` across multiple terminals; meaning you only have to unlock your keys once at startup, and you won't have to type in the password every time you go to do something.

lanyard was heavily inspired by [keychain](https://github.com/danielrobbins/keychain), although i wrote this because keychain is written in python and is thus quite slow.

## installation

### with Nix (recommended)

add the input to your `flake.nix`:

```nix
{

  inputs = {
    # ...
    lanyard = {
      url = "github:indium114/lanyard";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

}
```

and configure it with home-manager:

```nix
{
  inputs, # remember to take in inputs
  ...
}:

{

  imports = [
    inputs.lanyard.homeModules.default
  ];

  programs.lanyard = {
    enable = true;
    enableNushellIntegration = true;
    keys = ''
      id_ed25519
    ''; # list of the keys you want to manage; stored in ~/.ssh/*
  };

}
```

### manually

1. install the program as you prefer (build from source or download the binary from the releases tab on the right)
2. download the `shell/lanyard.nu` file, and follow the instructions inside.
3. edit `~/.config/lanyard/keys.conf` and add all  of the keys you want to manage; separated by newlines

## usage

the first time you open a shell after you boot your system (or after you configure lanyard for the first time), you will be asked to unlock your ssh keys. if you say yes, you will have to put in the passwords for each configured key. if you answer no, you will be prompted again when you open another terminal.

you can use the `lanyard wipe` command to remove all keys from the agent, in the event you need to do that for whatever reason.
