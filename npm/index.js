const path = require('path')
const { execFileSync } = require('child_process')

const PROJECT_NAME = 'mcp-time'

function getBinaryPath() {
  // Lookup table for all platforms and binary distribution packages
  const BINARY_DISTRIBUTION_PACKAGES = {
    'linux-x64': `${PROJECT_NAME}-linux-x64`,
    'linux-arm64': `${PROJECT_NAME}-linux-arm64`,
    'darwin-x64': `${PROJECT_NAME}-darwin-x64`,
    'darwin-arm64': `${PROJECT_NAME}-darwin-arm64`,
    'win32-x64': `${PROJECT_NAME}-win32-x64`,
  }

  // Windows binaries end with .exe so we need to special case them.
  const binaryName = process.platform === 'win32' ? `${PROJECT_NAME}.exe` : PROJECT_NAME

  // Determine package name for this platform
  const platformSpecificPackageName =
    BINARY_DISTRIBUTION_PACKAGES[`${process.platform}-${process.arch}`]

  try {
    // Resolving will fail if the optionalDependency was not installed
    return require.resolve(`${platformSpecificPackageName}/bin/${binaryName}`)
  } catch (e) {
    // Fallback to the binary downloaded by the postinstall script
    return path.join(__dirname, binaryName)
  }
}

module.exports.runBinary = function (...args) {
  execFileSync(getBinaryPath(), args, {
    stdio: 'inherit',
  })
}
