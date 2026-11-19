# Flutter DevTool 🚀 (Wireless Hotspot Development)

A lightweight CLI developer tool written in Go that enables you to build, run, debug, and hot-reload **Flutter applications on real Android devices without a USB cable, without an external Wi-Fi router, and without installing any third-party apps on your phone**.

Designed for developers on the go (coffee shops, trains, classrooms, coworking spaces, travel) where you only have your laptop and your smartphone's personal hotspot.

---

## ⚡ The Problem It Solves

Normally, running Flutter on a physical Android device requires either:
1. **A physical USB cable** tethering your phone to your computer.
2. **An external Wi-Fi router** connected to both devices (which frequently fails on public or office networks due to *AP Client Isolation* blocking peer-to-peer traffic).
3. **Third-party apps** (such as Shizuku, LADB, or Termux) installed on the phone to enable Wireless ADB.

When your phone broadcasts a Personal Hotspot to your laptop without an external router:
- Android OS detects that its active internet connection is **Cellular (Mobile Data)** rather than a Wi-Fi client connection, and by default grays out the **Wireless debugging** toggle.
- When pairing is opened, Android often displays its internal 464XLAT cellular carrier address (`192.0.0.2`), which cannot be reached over Wi-Fi.

**Flutter DevTool** automates the entire hotspot discovery, IP resolution, native ADB pairing, dynamic port scanning, and Flutter execution pipeline cleanly from your laptop terminal.

---

## 🌟 Key Features

* **Zero 3rd-Party Software on Phone**: Uses native Android 11+ Wireless Debugging built directly into Android Developer Options.
* **Automatic Hotspot Discovery**: Automatically identifies the active wireless network adapter and extracts the phone's gateway IP in under 50ms.
* **Automated Native Pairing & Port Scanning**: Performs native `adb pair` and automatically scans for the active dynamic wireless debugging port on the phone.
* **Persistent Connection**: Automatically activates persistent port `5555` on the phone so you don't have to re-pair on future sessions.
* **Flutter Integration & Clean Logs**: Launches your Flutter app directly on the wireless device and tags application logs (`[APP]`) while stripping boilerplate Gradle noise.
* **Built-in TCP Relay Engine**: Includes a multi-session bidirectional TCP relay server with generation tracking and reconnect message queuing.

---

## 📋 Prerequisites

1. **Laptop**:
   - Windows 10 or 11
   - [Go](https://go.dev/) (1.21+ recommended to build from source)
   - [Flutter SDK](https://docs.flutter.dev/)
   - [Android SDK Platform-Tools](https://developer.android.com/tools/releases/platform-tools) (`adb.exe`)
2. **Phone**:
   - Android 11 or higher
   - **Developer Options** enabled (*Settings $\to$ About Phone $\to$ tap "Build Number" 7 times*)

---

## 🛠️ Installation & Build

Clone the repository and build the executable:

```powershell
# Clone the repository
git clone https://github.com/your-username/devtool.git
cd devtool

# Build the binary
go build -o flutter-devtool.exe ./cmd/flutter-devtool
```

*(Optional)* Add the directory containing `flutter-devtool.exe` to your Windows `PATH` environment variable so you can run `flutter-devtool` from any folder.

---

## 🚀 Quick Start Guide (Zero 3rd-Party Apps)

### Step 1: Connect Devices

You can connect your laptop and phone in one of two ways:

#### Option A: Phone Hotspot (Most Common)
1. Turn ON **Personal Hotspot** on your phone.
2. Connect your laptop's Wi-Fi to your phone's hotspot.
3. ⚠️ **Important Tip**: Temporarily turn **OFF Mobile Data** on your phone for 30 seconds during pairing. *(This forces Android to bind Wireless Debugging to the Wi-Fi interface instead of the cellular 464XLAT `192.0.0.2` address).*

#### Option B: Laptop Hotspot (If phone restricts hotspot debugging)
1. Turn ON **Mobile hotspot** in Windows (*Press `Win + A` $\to$ click Mobile hotspot*).
2. On your phone, connect to your laptop's Wi-Fi network.

---

### Step 2: Pair Wirelessly (One-Time Setup)

1. On your phone, navigate to:
   **Settings** $\to$ **Developer options** $\to$ **Wireless debugging**
2. Toggle **Wireless debugging** to **ON**.
   *(If greyed out when using Option A, see Troubleshooting below to enable the Quick Settings developer tile).*
3. Tap **"Pair device with pairing code"**.
   A popup will show:
   * **Wi-Fi pairing code** (6 digits, e.g. `482915`)
   * **IP address & Port** (e.g. `10.235.10.94:38491` — the port after the colon is the pairing port)
4. On your laptop, open PowerShell and run:
   ```powershell
   .\flutter-devtool.exe pair
   ```
5. Follow the interactive prompts:
   * **Phone IP**: Press `Enter` to use the auto-detected IP (or type the IP shown on screen).
   * **Pairing port**: Enter the port shown in the popup (e.g. `38491`).
   * **Pairing code**: Enter the 6-digit code (e.g. `482915`).

The tool will pair, scan for the active connect port, connect ADB, and configure persistent port `5555`. You can now turn Mobile Data back on.

---

### Step 3: Run Your Flutter App

Run your Flutter project wirelessly:

```powershell
# If inside your Flutter project directory:
.\flutter-devtool.exe run

# Or point to your project path:
.\flutter-devtool.exe run --project "D:\Flutter\my_app"
```

Your app will compile, install onto your Android phone over Wi-Fi, and launch with live log streaming. Keyboard hot-reload keys (`r` for hot reload, `R` for restart, `q` to quit) work normally!

---

### Step 4: Reconnecting on Future Sessions

On future sessions while connected to the same hotspot, you don't need to re-pair! Simply run:

```powershell
# Automatically connects to phone over hotspot
.\flutter-devtool.exe connect

# Run your app
.\flutter-devtool.exe run
```

---

## 📖 CLI Command Reference

| Command | Description |
|---|---|
| `flutter-devtool hotspot` | Automatically detects the active Wi-Fi adapter, phone gateway IP, laptop IP, and reachability. |
| `flutter-devtool pair` | Interactively pairs with Android 11+ native Wireless Debugging without installing any 3rd-party apps. |
| `flutter-devtool connect` | Auto-detects the phone gateway and connects ADB over Wi-Fi (checking port 5555 and dynamic ports). |
| `flutter-devtool devices` | Lists all connected ADB devices. Prompts to connect if no devices are found. |
| `flutter-devtool run` | Builds and launches your Flutter project on the connected device with app-level log highlighting. |
| `flutter-devtool relay` | Starts the built-in TCP relay server (port 9000) for custom companion agent connections. |

---

## ❓ Troubleshooting & FAQs

### Q1: "Pairing failed: protocol fault (couldn't read status message)"
* **Cause**: Your phone displayed `192.0.0.2` in the pairing popup. `192.0.0.2` is a cellular 464XLAT modem address assigned by your carrier (e.g. Jio/Airtel/T-Mobile), NOT a Wi-Fi IP address. Your laptop cannot route packets to cellular modem interfaces.
* **Fix**: Turn **OFF Mobile Data** on your phone. Go back to Wireless debugging, tap "Pair device with pairing code" again, and use your real Wi-Fi IP (e.g. `10.x.x.x` or `192.168.x.x`). You can verify your real Wi-Fi IP by running `.\flutter-devtool.exe hotspot`.

### Q2: "Wireless debugging is greyed out when Hotspot is on"
* **Cause**: Android checks if the device is a Wi-Fi *client*. When broadcasting a hotspot on mobile data, Android thinks it has no Wi-Fi.
* **Fix A (Quick Settings Tile)**: Go to **Settings $\to$ Developer options $\to$ Quick settings developer tiles** $\to$ enable **Wireless debugging**. Now pull down your notification shade and tap the "Wireless debugging" quick tile directly to turn it ON.
* **Fix B (Laptop Hotspot)**: Invert the connection. Turn on **Mobile hotspot** in Windows (`Win + A`), and connect your phone to the laptop's Wi-Fi. The phone is now a Wi-Fi client, so the toggle immediately unlocks.

### Q3: "ADB executable not found"
* The tool automatically scans `PATH`, `ANDROID_HOME`, `ANDROID_SDK_ROOT`, `LOCALAPPDATA\Android\Sdk`, and root drive SDK paths (e.g., `D:\Android\Sdk`, `C:\Android\Sdk`). If your SDK is installed in a custom location, add its `platform-tools` directory to your system `PATH` or set the `ANDROID_HOME` environment variable.

---

## 🏗️ Project Architecture

```
devtool/
├── cmd/
│   ├── flutter-devtool/       # Main CLI entrypoint
│   │   ├── main.go            # Command router
│   │   ├── commands.go        # Command definitions & help
│   │   ├── connect.go         # Pair & connect implementations
│   │   ├── devices.go         # Device listing & status
│   │   ├── run.go             # Flutter run executor
│   │   └── selection.go       # Interactive device selection
│   └── test-client/           # Test client for TCP relay verification
├── internal/
│   ├── adb/                   # ADB discovery, device parsing & execution
│   ├── flutter/               # Flutter CLI runner & log stream filters
│   ├── hotspot/               # Routing table parsing & gateway auto-discovery
│   └── relay/                 # Bidirectional TCP relay, message queues & sessions
└── readme.md                  # Documentation
```

---

## 📄 License

MIT License. Contributions and PRs are welcome!
