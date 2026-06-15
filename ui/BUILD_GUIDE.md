DAEMON UI - BUILD GUIDE
=======================

This guide explains how to compile and run the Daemon UI in Delphi/ObjectPascal.

## REQUIREMENTS

- Delphi 10.4 Community or Enterprise Edition (or later)
- FireMonkey framework (FMX)
- Windows 7 SP1 or later (for deployment)
- At least 4GB RAM recommended

## INSTALLATION STEPS

### Step 1: Install Delphi

1. Download Delphi from: https://www.embarcadero.com/products/delphi
2. Run the installer
3. During installation, ensure you select:
   - C++ Build Tools
   - FireMonkey (FMX) framework
   - Windows 64-bit target
4. Complete the installation

### Step 2: Clone/Setup Project

```bash
cd daemon/ui/
```

### Step 3: Open Project in Delphi IDE

1. Open Delphi IDE
2. File → Open Project
3. Navigate to: `daemon/ui/DaemonUI.dproj`
4. Click Open

### Step 4: Configure Project

1. Project → Options
2. Set Target Platform to Windows 64-bit
3. Set Output Path to: `$(PROJECTDIR)\bin\`
4. Click OK

### Step 5: Build Project

#### Option A: Using Delphi IDE
- Press F9 to compile and run
- Or: Build → Build All

#### Option B: Using Command Line (RAD Studio Command Prompt)
```bash
# Navigate to UI directory
cd daemon\ui

# Compile (debug)
dcc32 DaemonUI.dpr

# Or compile (release)
dcc32 DaemonUI.dpr -DRELEASE
```

#### Option C: Using MSBuild
```bash
cd daemon\ui
msbuild DaemonUI.dproj /p:Config=Release /p:Platform=Win64
```

## DIRECTORY STRUCTURE AFTER BUILD

```
daemon/
├── ui/
│   ├── bin/
│   │   └── DaemonUI.exe          ← Compiled executable
│   ├── dcu/
│   │   └── [compiled units]
│   ├── Forms/
│   ├── Components/
│   ├── Utils/
│   ├── DaemonUI.dpr
│   ├── DaemonUI.dproj
│   └── README.md
```

## RUNNING THE APPLICATION

### Method 1: From Delphi IDE
- Press F9 after compilation

### Method 2: Direct Execution
```bash
daemon\ui\bin\DaemonUI.exe
```

### Method 3: From Explorer
- Navigate to `daemon/ui/bin/`
- Double-click `DaemonUI.exe`

## TROUBLESHOOTING

### Compilation Errors

**Error: "Unit 'FMX.xxx' not found"**
- Ensure FireMonkey is properly installed
- Project → Options → Search Path includes FireMonkey paths
- Rebuild package cache: Tools → Options → Environment → Delphi Direct

**Error: "Cannot find 'DCC.exe'"**
- Ensure RAD Studio Command Prompt is used
- Or set PATH to include: `C:\Program Files (x86)\Embarcadero\Studio\22.0\bin`

**Error: "Unit not found: 'MainForm'"**
- Verify Forms/ subdirectory exists
- Check Project → Options → Search Path includes current directory

### Runtime Errors

**Logo not displaying**
- Verify `images/daemonlogo02.png` exists relative to executable
- Check file permissions

**Dark theme not applying**
- Ensure `ApplyFLStudioTheme()` is called after form creation
- Check color values in `StyleManager.pas`

**Controls not responding**
- Verify event handlers are properly assigned
- Check component Parent properties

## DEVELOPMENT WORKFLOW

### Making Changes

1. Edit .pas file in Delphi Editor
2. Press Ctrl+S to save
3. Press F9 to recompile and test
4. Changes take effect immediately

### Adding New Components

1. Add component declaration in MainForm.pas
2. Create and configure in FormCreate()
3. Recompile (F9)

### Debugging

1. Set breakpoints by clicking margin in editor
2. Press F9 to run with debugger
3. Use Debug → Evaluate/Modify to inspect variables
4. Use Watch window to monitor values

## DEPLOYMENT

### Creating Installer

1. Tools → Deployment
2. Select target platform: Windows 64-bit
3. Right-click profile → Compile All
4. Right-click profile → Deploy All

### Standalone Executable

1. Compile in Release mode:
   - Project → Options → Compiler → Release settings
   - Press F9

2. Copy executable and dependencies:
   - `DaemonUI.exe` (main program)
   - Required DLLs from Delphi runtime libraries

3. Create installer using:
   - Advanced Installer
   - NSIS (Nullsoft Scriptable Install System)
   - Inno Setup

## COMMAND LINE OPTIONS

```bash
# Compile in verbose mode
dcc32 DaemonUI.dpr -V

# Compile with specific output directory
dcc32 DaemonUI.dpr -E..\bin\

# Compile with debug info
dcc32 DaemonUI.dpr -D

# Compile with optimization
dcc32 DaemonUI.dpr -O
```

## PERFORMANCE OPTIMIZATION

### For Release Build

1. Project → Options → Compiler
2. Code generation:
   - ☑ Optimizations
   - ☑ Use register for variables
3. Linker: ☑ Include TD32 debug info

### Runtime

- Minimize OnPaint events
- Cache results of expensive calculations
- Use TTimer for smooth animations instead of loops

## SUPPORTED PLATFORMS

| Platform | Status | Notes |
|----------|--------|-------|
| Windows 64-bit | ✅ Full | Primary target |
| Windows 32-bit | ✅ Supported | Change target in Project Options |
| macOS | ⏳ Future | Requires FMX macOS support |
| Linux | ⏳ Future | Requires FMX Linux support |

## NEXT STEPS

After successful compilation:

1. ✅ Verify UI displays correctly
2. ✅ Test transport controls (Play, Stop, Record)
3. ✅ Verify playlist functionality
4. ✅ Test step grid interaction
5. ⏳ Connect to audio engine
6. ⏳ Implement MIDI support

## ADDITIONAL RESOURCES

- Delphi Docs: https://docwiki.embarcadero.com/CodeExamples/Sydney/en/Main_Page
- FireMonkey Tutorials: https://www.embarcadero.com/blog/category/firemonkey
- RAD Studio Help: Press F1 in IDE

## SUPPORT

For compiler errors or issues:
1. Check Delphi version compatibility
2. Verify all required packages are installed
3. Clear DCU cache: Delete `/dcu/` folder and rebuild
4. Reinstall Delphi if problems persist
