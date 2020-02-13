# This file was generated by GoReleaser. DO NOT EDIT.
class Gomodctl < Formula
  desc "search,check and update go modules"
  homepage "https://github.com/beatlabs/gomodctl"
  version "0.1.2"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/beatlabs/gomodctl/releases/download/v0.1.2/gomodctl_Darwin_x86_64.tar.gz"
    sha256 "212c189f89c295c9cad95a0c30c284d2a05b4584ff9d3c84de41aecc0cc12ce7"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/beatlabs/gomodctl/releases/download/v0.1.2/gomodctl_Linux_x86_64.tar.gz"
      sha256 "27fadeb9a4d91f6036f633d93ef9a1edd73a0472fb83bc87995ebfc092085963"
    end
  end

  def install
    bin.install "gomodctl"
  end

  test do
    system "#{bin/gomodctl}"
  end
end
