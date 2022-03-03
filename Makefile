VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/0xhelloweb3/go-coin-wallet

pkgCore = ${module}/core/eth
pkgUtil = ${module}/util 
pkgGasNow= ${module}/gasnow
pkgGasNow= ${module}/constants

pkgAll =  $(pkgCore)
#$(pkgCore) $(pkgUtil) $(pkgGasNow) $(pkgConstants)

buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/wallet.aar ${pkgAll}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios ${pkgAll}