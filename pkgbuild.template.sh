pkgname=gogh
pkgver=$VERSION
pkgrel=1
pkgdesc='GO GitHub project manager'
arch=('x86_64')
url="https://github.com/kyoh86/$pkgname"
license=('MIT')
makedepends=('go')
depends=('glibc')
source=("$url/archive/refs/tags/v$pkgver.tar.gz")
options=('zipman')
sha256sums=(.)
prepare(){
  cd "$pkgname-$pkgver"
  mkdir -p build/
}
build() {
  cd "$pkgname-$pkgver"
  export LDFLAGS="-linkmode=external -s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=$(date)"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  go build -ldflags="${LDFLAGS}" -buildmode=pie -trimpath -mod=readonly -modcacherw -o build ./cmd/...
  go run -ldflags="${LDFLAGS}" -tags man ./cmd/gogh man
}
check() {
  cd "$pkgname-$pkgver"
  make test
}
package() {
  cd "$pkgname-$pkgver"
  install -Dm755 build/$pkgname "$pkgdir/usr/bin/$pkgname"
  if [ -f LICENSE ]; then
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  fi
  if [ -f "$pkgname.1" ]; then
    install -Dm644 "$pkgname.1" "$pkgdir/usr/share/man/man1/$pkgname.1"
  fi
}
