# CLAUDE.md — Flutter Project

## Project Context
- **Name:** [project_name]
- **Platform targets:** [iOS | Android | macOS | Web]
- **Min SDK:** iOS [17+] / Android [API 24+]
- **Dart SDK:** 3.11+ (arm64 macOS)
- **State management:** [Riverpod | BLoC]
- **Backend:** [API base URL or "BFF pattern"]

---

## Architecture
Feature-first + Clean Architecture layers:
- `data/` → DTOs, repositories impl, data sources (HTTP, local DB)
- `domain/` → entities, repository interfaces, use cases
- `presentation/` → screens, widgets, state notifiers/blocs

---

## Key Features & Flows
<!--
Descreva os fluxos principais:
- Auth: email/senha → JWT → secure storage → auto-refresh
- Onboarding: 3 telas → permissões → home
- Offline: cache local com Hive, sync na reconexão
-->

---

## Package Inventory
<!--
Liste os packages em uso — Claude não deve sugerir alternativas sem perguntar:
-->
```yaml
# State
riverpod: ^2.x
flutter_riverpod: ^2.x
riverpod_annotation: ^2.x

# Navigation
go_router: ^x.x

# Network
dio: ^x.x
retrofit: ^x.x

# Local storage
hive_flutter: ^x.x
flutter_secure_storage: ^x.x

# UI
cached_network_image: ^x.x
flutter_svg: ^x.x

# Utils
intl: ^x.x
logger: ^x.x
```

---

## Flavor / Environment Setup
```dart
// lib/core/constants/env.dart
enum Env { dev, staging, prod }
// Accessed via: AppEnv.current
```

---

## Navigation Structure
<!--
Documente a estrutura de rotas go_router:
/ → SplashScreen
/auth → LoginScreen
/auth/register → RegisterScreen
/home → HomeScreen (ShellRoute)
  /home/dashboard → DashboardScreen
  /home/profile → ProfileScreen
-->

---

## MCP / Debug Notes
- Para usar flutter MCP: `flutter run --debug --host-vmservice-port=9100 --enable-vm-service --disable-service-auth-codes`
- VM Service porta padrão do projeto: `9100`
- MCP configurado: `claude mcp add dart-flutter -- dart mcp-server --force-roots-fallback`

---

## Known Gotchas
<!--
- Xcode: se der linking error, rodar `cd ios && pod install --repo-update`
- Homebrew warnings de coreutils: adicionar `/opt/homebrew/opt/coreutils/libexec/gnubin` ao PATH
- DDS proxy: a porta que o Flutter mostra no output pode diferir da --host-vmservice-port
- Build macOS: requer `NSMicrophoneUsageDescription` no Info.plist mesmo sem usar mic
-->
