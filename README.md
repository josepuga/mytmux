# MyTmux

## What is MyTmux?

MyTmux allows you to create one or more tmux sessions with their corresponding windows and, optionally, execute commands in those windows at creation time.

Everything is configured through an ini file with a very simple syntax.

## How to use

You only need to run `mytmux` and specify the configuration file. For example:

```bash
mytmux ~/.config/tmux-workspaces.ini
```

## FConfiguration file

The configuration consists of an ini file defining the sessions and windows to be created:

```ini
[mi_session]
window1=/ruta/por/defecto
window2=/otra/ruta;df -h
logs=/tmp;ls *.log

[mi_other_session]
Home=~;fastfetch
#(...)
```

You can add as many sections as you want to the ini file (each section represents a tmux session).

Each `key=value` entry inside a section corresponds to a window.

The command after the `;` is optional. If specified, it will be executed automatically when the window is created.

Below is an example of my personal configuration file to show its potential. I use the `ai` session for various AI-related projects:

```ini
[ai]
Ollama=/mnt/ai/ollama;ollama list && podman ps
ComfyUI server=/mnt/ai/ComfyUI;source venv/bin/activate && python main.py --listen 0.0.0.0 2>&1 | tee  /tmp/comfyui.log
ComfyUI shell=/mnt/ai/ComfyUI;source venv/bin/activate
Fooocus=/mnt/ai/Fooocus;source venv/bin/activate && echo run with: python entry_with_update.py
```

### A real example

#### This is my configuration.

I have some terminals running AI programs, Ollama, ComfyUI Server, ComfyUI shell (with venv) and Fooocus.

![Ini Content](doc/ini-content.png)

To Create a `Tmux` session I run `mytmux` with that file

![Create Session](doc/session-created.png)

And After attach it, all tabs has been created and the commands runned.

![Tab Fooocus](doc/tmux-tab-fooocus.png)
![Tab ComfyUI server](doc/tmux-tab-comfyui-server.png)

### Explanation

This creates a session called `ai` with 4 windows:

`Ollama | ComfyUi server | ComfyUI shell | Fooocus`

I run Ollama using two containers: the GUI and the server. When starting the session, I want to see which models are available and the current state of the containers, hence the command.

For ComfyUI and Fooocus, the virtual environment is activated automatically, along with other commands that are not necessary to explain here.

## TODO

- Implement a default active window. Currently, the active window is always the last one created.
- Add pane support. Panes are not implemented because I never use them; I prefer having a full window and switching tabs using the keyboard.

```

```
