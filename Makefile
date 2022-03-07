VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/coming-chat/wallet-SDK

pkgCore = ${module}/core/eth

pkgEth =  $(pkgCore)
#$(pkgCore) $(pkgUtil) $(pkgGasNow) $(pkgConstants)

pkgPolka = ${module}/wallet

buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/ethwalletcore.aar ${pkgEth}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/ethwalletcore.xcframework ${pkgEth}

buildAllSDKIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/Wallet.xcframework ${pkgEth} ${pkgPolka}

packageAll:
	rm -rf ${outdir}/*
	@make buildAllAndroid && make buildAllIOS
	@cd ${outdir} && mkdir android && mv eth-wallet* android
	@cd ${outdir} && tar czvf android.tar.gz android/*
	@cd ${outdir} && tar czvf eth-wallet.xcframework eth-wallet.xcframework/*