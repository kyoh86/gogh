# Maintainer: kyoh86 <me@kyoh86.dev>
pkgname=gogh
pkgver=$VERSION
pkgrel=1
pkgdesc='GO GitHub project manager'
url="https://github.com/kyoh86/$pkgname"
arch=('x86_64')
license=('MIT')
makedepends=('go')
depends=('glibc')
source=("$url/archive/refs/tags/v$pkgver.tar.gz")
options=('zipman')
sha256sums=(.)
echo "${pkgname}" "${pkgver}" "${pkgrel}" "${pkgdesc}" "${url}"
echo "${arch[@]}"
echo "${license[@]}"
echo "${makedepends[@]}"
echo "${depends[@]}"
echo "${source[@]}"
echo "${options[@]}"
echo "${sha256sums[@]}"
prepare(){
  cd "$pkgname-$pkgver" || exit 1
  mkdir -p build/
}
build() {
  cd "$pkgname-$pkgver" || exit 1
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  LDF="-linkmode=external -s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=$(date --iso-8601=seconds)"
  go build \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -modcacherw \
    -ldflags="${LDF}" \
    -o build ./cmd/...
  go run -ldflags="${LDF}" -tags man ./cmd/gogh man
}
check() {
  cd "$pkgname-$pkgver" || exit 1
  go test ./...
}
package() {
  cd "$pkgname-$pkgver" || exit 1
  declare pkgdir
  install -Dm755 build/$pkgname "$pkgdir/usr/bin/$pkgname"
  if [ -f LICENSE ]; then
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  fi
  if [ -f "$pkgname.1" ]; then
    install -Dm644 "$pkgname.1" "$pkgdir/usr/share/man/man1/$pkgname.1"
  fi
}
