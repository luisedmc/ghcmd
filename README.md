<div align="center">

# GHCMD

**Forget your browser, use GitHub from the terminal.**

<!-- [About](#about) •
[Usage](#usage) •
[Screenshots](#screenshots) -->

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

## :zap: Usage

First of all, you will need a GitHub API Key. You can know more about and get one [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic) and also check the [GitHub API Documentation](https://developer.github.com/v3/).

About the application, it's actually really simple. The status bar at the bottom of the screen shows your API Key status and how to navigate through the application. Check the Screenshots section to see how it looks like.

## :dart: Features

In the current version, you can:

| Service             | Description                                   |
| ------------------- | --------------------------------------------- |
| `Search Repository` | Search for a specific repository from an user |
| `Create Repository` | Create a repository in your GitHub account    |

For now, it is a work in progress and only supports a few commands. I will be always trying to add more features and improve the existing ones.

## :camera_flash: Screenshots

<div align="center">
    Main view</br>
    <img src="/docs/main.png" alt="Main view" width=380 height=380>
</div>

## :page_facing_up: License

- [MIT](https://raw.githubusercontent.com/luisedmc/ghcmd/master/LICENSE)
