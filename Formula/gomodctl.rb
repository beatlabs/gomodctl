# This file was generated by GoReleaser. DO NOT EDIT.
class Gomodctl < Formula
  desc "search,check and update go modules"
  homepage "https://github.com/beatlabs/gomodctl"
  version "0.1.0"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/beatlabs/gomodctl/releases/download/v0.1.0/gomodctl_Darwin_x86_64.tar.gz"
    sha256 "3a4e2d4044830dbcd3ac95ca3b28ebf95184b9b11f4f794343be2c92d274ac94"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/beatlabs/gomodctl/releases/download/v0.1.0/gomodctl_Linux_x86_64.tar.gz"
      sha256 "9b6be2cd9fe1a4407653e31a8d2e779abd3313fe858e62598600b3dcd70ca549"
    end
  end

  def install
    bin.install "gomodctl"
  end

  test do
    system "#{bin/gomodctl}"
  end
end
