class GoNotes < Formula
  desc "Terminal note-taking app with Notion-like styling"
  homepage "https://github.com/HugoBritez/go-notes"
  url "https://github.com/HugoBritez/go-notes/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "ecb3651605f5198d0736c86add5668c46d1a001a32b7e3f6ac274a458386bb61"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"note"), "main.go"
  end

  test do
    system "#{bin}/note", "--help"
  end
end