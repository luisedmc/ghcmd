<div align="center">

# ghcmd

**A Terminal User Interface for Github written in Golang**

</div>

## About

GHCMD is a simple terminal user interface for GitHub. The main goal of this application is to provide an intuitive and fast way to interact with GitHub without leaving the terminal. Keep in mind that it is not a full client and does not provide all the features that GitHub has.

Tools and Libraries used:

- [Go](https://go.dev/)
- [go-github](https://github.com/google/go-github)
- [GitHub API](https://developer.github.com/v3/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)

## Installation

```
#  Clone the repository
git clone https://github.com/luisedmc/ghcmd.git

#  Go to the project directory
cd ghcmd

#  Run the application
go run .
```

You can also build the application and run it as a binary file in any directory. To do so, you build using `go build` and you can check [here](https://zwbetz.com/how-to-add-a-binary-to-your-path-on-macos-linux-windows/) how to add the binary to your path.

## Usage

When you first start the application you'll need a GitHub personal access token. GitHub provides a simple step-by-step guide to create this token, which you can find [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic).

Once you have your token, you can enter it in the terminal and start using the application.

## Token

Your API key is stored locally. There's no option to change it currently, so if you want to change your key you'll need to delete the token file, which is located at:

| OS      | Path                                                          |
| ------- | ------------------------------------------------------------- |
| macOS   | `~/Library/Application Support/ghcmd/token`                   |
| Linux   | `~/.config/ghcmd/token` (or `$XDG_CONFIG_HOME/ghcmd/token`)  |
| Windows | `%AppData%\ghcmd\token` (e.g. `C:\Users\<user>\AppData\Roaming\ghcmd\token`) |

## Features

In the current version, you can:

| Service             | Description                                   |
| ------------------- | --------------------------------------------- |
| `Search Repository` | Search for a specific repository from an user |
| `Create Repository` | Create a repository in your GitHub account    |

For now, it is a work in progress and only supports a few commands. I will be always trying to add more features and improve the existing ones.

## Screenshots

<div align=center>
    This is how a successful search looks like<br>
    <img src="/docs/example.gif" width=700 height=400>
</div>

## License

- [MIT](https://raw.githubusercontent.com/luisedmc/ghcmd/master/LICENSE)
