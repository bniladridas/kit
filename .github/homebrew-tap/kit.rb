class Kit < Formula
  desc "Lightweight developer CLI tool for GitHub workflows"
  homepage "https://github.com/bniladridas/kit"
  version "1.0.0"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/bniladridas/kit/releases/download/v#{version}/kit_#{version}_darwin_arm64.tar.gz"
    sha256 "<darwin_arm64_sha256>"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/bniladridas/kit/releases/download/v#{version}/kit_#{version}_darwin_amd64.tar.gz"
    sha256 "<darwin_amd64_sha256>"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/bniladridas/kit/releases/download/v#{version}/kit_#{version}_linux_amd64.tar.gz"
    sha256 "<linux_amd64_sha256>"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/bniladridas/kit/releases/download/v#{version}/kit_#{version}_linux_arm64.tar.gz"
    sha256 "<linux_arm64_sha256>"
  end

  def install
    bin.install "kit"
  end

  test do
    system "#{bin}/kit", "version"
  end
end
