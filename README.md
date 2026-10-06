# Devtool CLI 🚀

`devtool-cli` is an elite, high-performance, modular command-line interface built in Go. It acts as a multi-tool for developers, automating the most common (and annoying) daily tasks like fetching APIs, formatting JSON, sorting messy folders, decoding JWTs, converting Unix timestamps, and much more.

It features **interactive UI components** (spinners and progress bars), stream-based processing for massive files, and extreme concurrency for blazing-fast execution.

## 📦 Installation

*(Assuming you have GoReleaser set up, binary downloads will appear in the Releases tab. For now, you can install via `go install`)*

```bash
go install github.com/GhanshyamJha05/devtool-cli@latest
```

## 🛠️ Features & Commands

### 1. Network & API Tools
*   **`devtool get <url>`**: Fetch data from any API and pretty-print the JSON response, status code, and latency.
    *   `--save <filename>`: Stream massive files directly to disk without crashing memory, complete with a real-time visual progress bar.
*   **`devtool serve`**: Instantly spin up a local static HTTP server in your current directory.
    *   `--port <number>`: Specify a custom port (default is 8080).

### 2. File Organization
*   **`devtool clean <folder>`**: Scan a chaotic directory (like your Downloads folder) and automatically sort files into categorized subfolders (`Images`, `Code`, `Documents`, etc.) concurrently.
    *   `--dry-run`: Preview what will happen without actually moving any files.
*   **`devtool undo <folder>`**: Instantly revert a previous `clean` operation. Files are put back exactly where they were using an auto-generated `.devtool-undo.json` mapping.

### 3. Developer Utilities
*   **`devtool jwt <token>`**: Decode and inspect a JSON Web Token securely on your local machine. It pretty-prints the Header and Payload, and clearly tells you if the token is expired or how much time is left.
*   **`devtool time <timestamp|now>`**: Convert Unix epoch timestamps (in seconds or milliseconds) to human-readable Local and UTC time, or get the current epoch stamp by passing `now`.
*   **`devtool format <file.json>`**: Read a minified or ugly JSON file and print a beautifully indented, human-readable version.
    *   `--in-place` (`-i`): Batch format an entire directory of JSON files instantly.
*   **`devtool uuid`**: Quickly generate and print a random UUID (Version 4) to your terminal.
*   **`devtool hash <file>`**: Compute the cryptographic hash of a file (defaults to `SHA-256`).
    *   `--md5`: Use MD5 instead.
    *   `--string "text"`: Hash raw text instead of a file.
*   **`devtool base64 <string>`**: Encode raw text into a Base64 string.
    *   `--decode` (`-d`): Decode a Base64 string back into readable text.

### 4. Self-Updater
*   **`devtool update`**: Automatically pings the GitHub repository, checks for the latest release, and safely replaces your current executable with the newest version so you never fall behind.

## ⚙️ Configuration

`devtool-cli` supports a configuration file! By default, it looks for `~/.devtool-cli.yaml`. 
You can use this to customize the CLI's behavior. For example, you can define custom file categories for the `clean` command:

```yaml
categories:
  Music:
    - .mp3
    - .wav
    - .flac
  Design:
    - .fig
    - .psd
```

## 🐛 Debugging
Every single command supports a global `--verbose` (`-v`) flag. If something isn't working right, tack on `-v` to get highly detailed, step-by-step debug logs.

## 🤝 Contributing
Pull requests are welcome! The project uses Cobra for CLI structuring, Viper for configuration, and Pterm for interactive UI components. 

1. Fork the repo.
2. Create your feature branch (`git checkout -b feature/cool-idea`).
3. Commit your changes (`git commit -m 'feat: Add some cool idea'`).
4. Push to the branch (`git push origin feature/cool-idea`).
5. Open a Pull Request.
