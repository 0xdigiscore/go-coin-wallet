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
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/eth-wallet.aar ${pkgAll}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/eth-wallet.xcframework ${pkgAll}

packageAll:
	rm -rf ${outdir}/*
	@make buildAllAndroid && make buildAllIOS
	@cd ${outdir} && mkdir android && mv eth-wallet* android
	@cd ${outdir} && tar czvf android.tar.gz android/*
	@cd ${outdir} && tar czvf eth-wallet.xcframework eth-wallet.xcframework/*