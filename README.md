# When My Meeting

A cross-platform desktop app that shows your upcoming meetings and helps you keep track of your schedule.

## Features

- 📅 Calendar integration
- 🔔 Meeting notifications
- 🖥️ System tray / menu bar support
- ⚙️ Configurable notification settings
- 🐧 Linux
- 🪟 Windows
- 🍎 macOS

## Downloads

Get the latest version from the [releases](https://github.com/makaksel/when-my-meeting/releases) page.

Available packages:

- Linux — `.deb`
- Windows — `.exe`
- macOS — `.dmg`

## Development

### Requirements

- Go 1.25+
- Fyne
- nfpm (for Linux packages)

### Run locally

```bash
make run
```

### Build

Build packages for the current platform:

```bash
make build-deb
make build-windows
make build-macos
```

Specify a version:

```bash
make build-macos VERSION=0.1.0
```

Build artifacts are placed in dist/.

---

# When My Meeting

Небольшое кроссплатформенное десктопное приложение, которое отображает предстоящие встречи и помогает следить за расписанием.

## Возможности

- 📅 Интеграция с календарем
- 🔔 Уведомления о встречах
- 🖥️ Поддержка системного трея / строки меню
- ⚙️ Настраиваемые параметры уведомлений
- 🐧 Linux
- 🪟 Windows
- 🍎 macOS

## Загрузка

Последнюю версию можно найти на странице [releases](https://github.com/makaksel/when-my-meeting/releases).

Доступные пакеты:

- Linux — `.deb`
- Windows — `.exe`
- macOS — `.dmg`

## Разработка

### Требования

- Go 1.25+
- Fyne
- nfpm (для пакетов Linux)

### Запуск локально

```bash
make run
```

### Сборка

Сборка пакетов для текущей платформы:

```bash
make build-deb
make build-windows
make build-macos
```

Указание версии:

```bash
make build-macos VERSION=0.1.0
```

Артефакты сборки размещаются в папке `dist/`.
