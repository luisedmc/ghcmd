<div align="center">

# GHCMD

**Forget your browser, use GitHub from the terminal.**

</div>

## :pushpin: About

GHCMD is a simple command line tool for GitHub with a terminal user interface. The main goal of this application is to provide an intuitive and fast way to interact with GitHub without leaving the terminal. Keep in mind that it is not a full client and does not provide all the features that GitHub has.

Tools and Libraries used:

- [Go](https://go.dev/)
- [go-github](https://github.com/google/go-github)
- [GitHub API](https://developer.github.com/v3/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [teacup](https://github.com/mistakenelf/teacup)
- [goleveldb](https://github.com/syndtr/goleveldb)

## :rocket: Installation

```
#  Clone the repository
git clone https://github.com/luisedmc/ghcmd.git

#  Go to the project directory
cd ghcmd

#  Run the application
go run .
```

You can also build the application and run it as a binary file in any directory. To do so, you build using `go build` and you can check [here](https://zwbetz.com/how-to-add-a-binary-to-your-path-on-macos-linux-windows/) how to add the binary to your path.

## :zap: Usage

The only thing you need is a GitHub API Key. You can know more about and get one [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) and also check the [GitHub API Documentation](https://developer.github.com/v3/). The first time you run the application, you will be asked to insert your token. After that, you will be able to use the application.

## :dart: Features

All the main functionalities are implemented by now. <br>
In the current version, you can:

| Service             | Description                                   |
| ------------------- | --------------------------------------------- |
| `Search Repository` | Search for a specific repository from an user |
| `Create Repository` | Create a repository in your GitHub account    |

Also, your API Key is stored locally in a database. There's no option to change it yet, so if you want to insert a new one you will need to delete the database file. It is located at `./db/data`

For now, it is a work in progress and only supports a few commands. I will be always trying to add more features and improve the existing ones.

## :camera_flash: Screenshots

<div align=center>
    This is how a successful search looks like<br>
    <img src="/docs/example.gif" width=700 height=400>
</div>

## :page_facing_up: License

- [MIT](https://raw.githubusercontent.com/luisedmc/ghcmd/master/LICENSE)
