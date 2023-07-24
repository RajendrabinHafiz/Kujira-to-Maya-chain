//go:build !testnet && !stagenet
// +build !testnet,!stagenet

package mayachain

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

func importPreRegistrationMAYANames(ctx cosmos.Context, mgr Manager) error {
	oneYear := fetchConfigInt64(ctx, mgr, constants.BlocksPerYear)
	names, err := getPreRegisterMAYANames(ctx, ctx.BlockHeight()+oneYear)
	if err != nil {
		return err
	}

	for _, name := range names {
		mgr.Keeper().SetMAYAName(ctx, name)
	}
	return nil
}

func migrateStoreV96(ctx cosmos.Context, mgr Manager) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v88", "error", err)
		}
	}()

	err := importPreRegistrationMAYANames(ctx, mgr)
	if err != nil {
		ctx.Logger().Error("fail to migrate store to v88", "error", err)
	}
}

func migrateStoreV102(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v102", "error", err)
		}
	}()

	if err := mgr.BeginBlock(ctx); err != nil {
		ctx.Logger().Error("fail to initialise block", "error", err)
		return
	}

	// Remove stuck txs
	hashes := []string{
		"D19F621FAD0AE81688E4AF40EA9D0CD132A8A4DBFF3EA56F443E2D9083915F17",
		"A03C0A41909D85B2DF2F7E9D5D13F6E0AF89F366F6B580C0CCC13F5CEC0A9872",
		"7B7CC323ED0AD04DCB26DF1DEB46DE02B85345499336D043CBB5582EB77D22DB",
		"CE29D8AD79314333E307265529256304E26FC0B538B19B2D07578BE3D6252CE4",
		"B6EBF457EB1817E852722CB9F51C26E45C35F58B2445048FE4BD38FD1A603894",
		"199838DB755A6199AB401AAB1D56D296C66B0972001CB033B9CDC4217E636270",
		"674BDD72DF068A95EFA5DD94C4691A1D492A3342DA368DF0799ADD4D344D694D",
		"6F9C3D5AD6221159191540CE55704BCBB446626B209C852DC29C5C0AC7A24A82",
		"F3A3041FD304B11B8EBB748C9BE964E1FCCE0004770B109F5F9B72114F7FB9B9",
		"4912D98B5C8D9D090CAD2732754F39FFB324DA7008A19A0235DD77A4AB8EF3E3",
		"C82029C6D3F7D8D226E9B13F09CD05CF30FEA15F6C96BB8D49E20A4E063F6E82",
		"46C783972218F50015281F28222F5DC46FA3926EABB93549A383180C43064F96",
		"081819200E3ACC82CC8D95DFA87A6F0D87704154922022F777FFA5AD82B1BEF0",
		"9C3B5774352256A37CC3B26B82287458C4D4DEDC988342E6A2088A1800ACE992",
		"8B651D92B0374FA4E97834E86D35601940F90E104B800BEA836685E28452A953",
		"EA57B9FC879E981598732F6112255D756593D354DC88712665FEDC354374AD41",
		"672D551F02D6030A77745E25C0C8768347BBDA35DD7AA61C02751C86799D7C18",
		"116237EDC4814A9F684D8FCBC58FB5ADED2A9386B5ED0F1E627BCEFA8246474C",
		"3EE6362906A180279B0B9470221017465A1AA25807EFDB5A7B9342A95E120E2B",
		"1B04FF39247F519BC01F88FB1AC6843223FF351C47DC1D96B0FEC782645463F1",
		"25BC9C71B8F4D071A684A327D6E2657DA1D01D241E419C5C705D3690B5653C2E",
		"758B2A8DF6BC62F1A922DEA5E75F585A9BDB39CFC01152E4C74FBF929C5B5777",
		"A502C210DE19555884464A27408E16C378D9327BDC155EE3076F7D3D8CC8B963",
		"95AD2B7B2EE2E2CBEF272B20AB271400CDF57EE8EC170F8B265554A9FC24542A",
		"279B076B50ADCF2DE06CF129DE6B4917754F56FFB7CF4644FF0429CFC49A0D23",
		"D427493DEF0DFE953194E2E3C633C7EF3AAEC38F77B06B9E1EFAEFFF2071D58C",
		"B6F9C4CB2ABC7FB80B336950E559DD3020CC44C8F92A6AD9D3449612A5A232CC",
		"DDB8A4FD768443BF36187EA6147469A3D4975ED0CD8B4DDB2140EF4B924C7817",
		"C6E972C90798E33317DAC162D7B419AF825540F352A7CF38A5AB1297EAA866E9",
		"313F924DC160D565573F3B9D0A47F378A099606FFE4D059947B5377AE98E9F65",
		"E8340277F7E2310DDD2435A52EC1CC7C07C6D33FCD1F4ABD513FF23B6B19990F",
		"CDFDEBF26E28789F7C272813524C7F3766A9B82AFF55CBAD9AF347121061171B",
		"9F10CA47145E9F6B6EB4297272D9DF9999A75937CD154D6C4DAEC3DBFE14C3D6",
		"1D843810B79E7ED1CDD5424B0FDBD6158ADE77479E1C006F2583E5263E26E667",
		"1A486C051F7478CE845D67E019EE30DEF58D61B8EE408FE43E6CD520DE45518F",
		"EC461F353F95D933723BDAE7945B970A7F45DCA68A671D06B0FE9AA206686EFC",
		"66F228BD65D1A82C6D78C234D1A86F1C7E118A1051D87FE6546F708E401720FA",
		"233D1D0FB660BDA2C3C13C9B6C2BD0E96E81E05EC93C43A526CE0B782CA4ADA1",
		"8C867538F1C5A564C1C82206CBF0B96277B66E630BF13473E51D27BAE8B1994E",
		"BEA5B29954A3634B37CC0D73EB30EB8427ABF58900521A413F0A66C73AD6742A",
		"F3FA329499C42BF258B4D79E43ADEEB1E9C56FB60D4A9390B12F4946A554642A",
		"B26E3D4D8458DD43DB3B6424F1310B457053BDA95DC22D3936FCC373B49C95AC",
		"972E49601D4BF9949C3B91162399249B4AC997ED1BA830DB6DBC7DF44ABBEF3D",
		"D4B8E0F61978046D1205B5DC857BBD887214BB7054113B499FAADF7105F4CFE0",
		"97C8E399272FD9C64C2E2F1E2E32804157BFBA71504B4B838850F2590F87D781",
		"50C0CCD601689011E88B54358001EA2C6B1E8C0AC6794D1A7D8C95A74256071C",
		"640702E326CB6B61CC7285B0ADCE6DEF0694E9CFB629FD32C34A5475B5391E9E",
		"C60AE6A164FE3BE8B2BE87543B25B0F36E199E1CE466CB09482D9ED7D2D78BA9",
		"A58A823D9E467B368713D65090DCCFAAD92D1C8D6F2B57E3933EB8ACB9946031",
		"8EA259B4E7D15FBB6703C0A1248A137BCBAA7255ABEC09CEFAC2FB34DF7BC2F5",
		"B20036A869329CA1CDF966F0443B8B524A2CF6AB4F4ECA7C359D61A0A167F36B",
		"05F77D6640AE44FCFCB30FBDB8E76F82C4FD75E05A6DC48271EC499A1A09C378",
		"8E70E838DCEAD4D1763B4B40E59942EEFE5492B631D9CFB303A1DB0F7075F835",
		"C268C6821C3A8C19B435B4591F216A20E8DAD9AD2C17EF59F7CF9BD4DA2B4536",
		"9E1DF502EE17709E267AEDF673CF94B188698C0CBB9A6FAAAF57EAE20D043495",
		"FA194EDF6818312E6B28AB1D228A44B8623415595ED1716E7B7A92CB3DFCDE36",
		"941D3F4B252B735C2D358A368724DC809ED9CC63D6ED4426E369E75175EDF0F0",
		"24B077C67D4F835D176B701EEC59FA1C14143A0E3ACFF64077632FB3CEBD2851",
		"C4E86318378C561AD16DF9697F09224D254A314BED36EC7AC6C0B7F35FAB5CDB",
		"BEFFC122704DB5525A9511411A942F7F06EECF6386C104BB05622EDAE94D8096",
		"5422336EF4134851F601A74AA30C5E47702CC08111775ADC3944F4F0B467CD4F",
		"308FDED05E0F39E103A6E3898A497A1F28806ED7DEAB2D88F825E95CC4942D53",
		"C41ADECBC9D85D956D3246CEFD350E54CCABDA2B315793FD2625D30BEA0763C4",
		"333C9BC7B7479D4A675307B63AB2372C89C9C21A75C379BB6FE8EA8FB83813A0",
		"A967482A359194C6B3E0045F68B2E11CD275B29FF7F3A7F6129902D90FFA7055",
		"9DD6CFA490E5ED47BAE45E1CEE141329C411D8BAF5642758CCF3749D13862076",
	}
	removeTransactions(ctx, mgr, hashes...)

	// Rebalance asgard vs real balance for RUNE
	vaultPubkey, err := common.NewPubKey("mayapub1addwnpepq0tgksv4kjn0ya5n4gt2546dnasw84nr3zdtdzfud9z984p8pvmnu5t3qsy")
	if err != nil {
		ctx.Logger().Error("fail to get vault pubkey", "error", err)
		return
	}
	vault, err := mgr.Keeper().GetVault(ctx, vaultPubkey)
	if err != nil {
		ctx.Logger().Error("fail to get vault", "error", err)
		return
	}

	vault.SubFunds(common.NewCoins(common.NewCoin(common.RUNEAsset, cosmos.NewUint(3947_32403277))))

	if err := mgr.Keeper().SetVault(ctx, vault); err != nil {
		ctx.Logger().Error("fail to set vault", "error", err)
		return
	}

	// Remove retiring vault
	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, RetiringVault)
	if err != nil {
		ctx.Logger().Error("fail to get retiring asgard vaults", "error", err)
		return
	}
	for _, v := range vaults {
		runeAsset := v.GetCoin(common.RUNEAsset)
		v.SubFunds(common.NewCoins(runeAsset))
		if err := mgr.Keeper().SetVault(ctx, v); err != nil {
			ctx.Logger().Error("fail to save vault", "error", err)
		}
	}

	// Add LPs from unobserved txs
	lps := []struct {
		MayaAddress string
		ThorAddress string
		TxID        string
		Amount      cosmos.Uint
		Tier        int64
	}{
		// users which don't have an LP position yet
		{
			MayaAddress: "maya142m4adpj57hkrymqe5n8zzcxm5cqccpn3a6w6y",
			ThorAddress: "thor1jzzaw44tr0cxgxaah7h2sen2ck03lllw882wn2",
			TxID:        "73217ACF7F4061089236E29588825603FB4025E40AC5835586ED0B7959BE4A1F",
			Amount:      cosmos.NewUint(2_00000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya15kg7dfew844rdh5esgkrdevp78yhf4fjryjcfu",
			ThorAddress: "thor1hd9p0fllkwkgj9epe3nynr253az7uclxs4g2uw",
			TxID:        "5594D2500BB36F70ADB4063B4D7A331DCE884D2C34373EDBD69022C33E31CD0F",
			Amount:      cosmos.NewUint(1_00000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya17lllslx89rrxu0dh6y9ctz0aa2j82tljnuuy9s",
			ThorAddress: "thor1vmq7vwk8t6sxg730aps5vqetm905ndtmcvdq69",
			TxID:        "10376393CBF1C9E92CCBBDF582FFE9896FC04E82C2E9C641B4CB18A23559E43E",
			Amount:      cosmos.NewUint(2_00000000),
			Tier:        1,
		},
		{
			MayaAddress: "maya1f40wek6sj6uay6nplxpe2c6pj98c5uq78xspa4",
			ThorAddress: "thor1f40wek6sj6uay6nplxpe2c6pj98c5uq783wdt9",
			TxID:        "2A7297AD1EB1F1C53C90241264E78F067DE94F1C80588C208E8B7B5D86B3B9E7",
			Amount:      cosmos.NewUint(5388_00000000),
			Tier:        1,
		},
		{
			MayaAddress: "maya1j8pswr7vpf9jjmhrn0xlwvzla2f9yfxwcwtj0p",
			ThorAddress: "thor1y9h2yk95c6uqp29xglkgyf9kqxqnu28nn6vwwz",
			TxID:        "D5BEA6C8B3170B418ACD67B8C8A44A60CD0A66696B9B691E7A7471E393F5E8B4",
			Amount:      cosmos.NewUint(1_00000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya1k83lm2nyrd7vgl8h9xcjhwu9kr2zecauslje79",
			ThorAddress: "thor1k83lm2nyrd7vgl8h9xcjhwu9kr2zecausgv4g4",
			TxID:        "DEF2BC77DFCDA774C81D921C8846886FFF804D462F0E6BFF78DBAA1ADDF72E68",
			Amount:      cosmos.NewUint(24_83513163),
		},
		{
			MayaAddress: "maya1p3hcnlfdla2647rpersykfatplvhkehd2duspa",
			ThorAddress: "thor1p3hcnlfdla2647rpersykfatplvhkehd26zuhd",
			TxID:        "C4F73CFBAC15565CCAED86B66EB405AE9E36F712F0457F8353D050FF37D636BB",
			Amount:      cosmos.NewUint(1_00000000),
		},
		{
			MayaAddress: "maya1pn03td7tzsftp6xz25r5fas43dgqynpf0lyan5",
			ThorAddress: "thor1pn03td7tzsftp6xz25r5fas43dgqynpf0g639y",
			TxID:        "A85BE46FFDD915D2074EC85C8E5B63B0407EFDD44CC6094CCC9A616A7FFB0494",
			Amount:      cosmos.NewUint(1_00000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya1s0ry4c65c7k020vgpykjfy5rkqv8d7yn60lzx6",
			ThorAddress: "thor1s0ry4c65c7k020vgpykjfy5rkqv8d7yn6cpws2",
			TxID:        "4A1CA0E1D87869C5083F6BBD2042BF5DA5545B01ECE9CD7922F11D8AB715B261",
			Amount:      cosmos.NewUint(1_00000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya1vwslytml73dclz0h4enc2xluf4z03esrt36n6r",
			ThorAddress: "thor1vwslytml73dclz0h4enc2xluf4z03esrtxylvn",
			TxID:        "237CBC3570DA3AE95D15F6E7C04A50EF3799A4106434A9A831A11BEDA8EB0FF6",
			Amount:      cosmos.NewUint(36_00000000),
			Tier:        1,
		},
		{
			MayaAddress: "maya1wlx25u0692nvxllg57tgt45h53hjsgggzlgavn",
			ThorAddress: "thor1cjlsyrzmfpldxhmz4j3yzyc0f6dp57lhv6cm2r",
			TxID:        "C8D1F65C6C6559D4A23E8BB47533E86CC25D8C41FA8382EC2C6FBF868953AB23",
			Amount:      cosmos.NewUint(1_50000000),
			Tier:        3,
		},
		{
			MayaAddress: "maya1zgtzwkd9qaagvwedgnmxeh9tsqc8wdsjwjxf6e",
			ThorAddress: "thor1zgtzwkd9qaagvwedgnmxeh9tsqc8wdsjw9c9vf",
			TxID:        "79A5288200EB347569B7E3707A822E72B2DB1CCD52BC035323DE2B1DC44273B3",
			Amount:      cosmos.NewUint(499_98000000),
			Tier:        1,
		},
		// users which already have an existing LP position
		{
			MayaAddress: "maya10nqg4w30e9dnm0qg7swa8qsyqevuemwx78dpdx",
			ThorAddress: "thor10nqg4w30e9dnm0qg7swa8qsyqevuemwx7sndmk",
			Amount:      cosmos.NewUint(5_58000000),
		},
		{
			MayaAddress: "maya14sanmhejtzxxp9qeggxaysnuztx8f5jra7vedl",
			ThorAddress: "thor14sanmhejtzxxp9qeggxaysnuztx8f5jrafj4m0",
			Amount:      cosmos.NewUint(958_08765797),
		},
		{
			MayaAddress: "maya17w5n2r7akuunq9e296y22qrljh3qqegf6usf5x",
			ThorAddress: "thor17w5n2r7akuunq9e296y22qrljh3qqegf6tw9zk",
			Amount:      cosmos.NewUint(1400_00000000),
		},
		{
			MayaAddress: "maya1a4v8ajttgx5u822k2s8zms3phvytz3at2a7mgj",
			ThorAddress: "thor1a4v8ajttgx5u822k2s8zms3phvytz3at22qh7z",
			Amount:      cosmos.NewUint(1_000000),
		},
		{
			MayaAddress: "maya1fdl7xga4sxhwlfs48fhkgwen88003g3hl006pn",
			ThorAddress: "thor1fdl7xga4sxhwlfs48fhkgwen88003g3hlc3khr",
			Amount:      cosmos.NewUint(1_00000000),
		},
		{
			MayaAddress: "maya1hh03993slyvggmvdl7q4xperg5n7l86pufhkwr",
			ThorAddress: "thor1wlzhcxs0r4yh4pswj8zfqlp7dnp95p4kxn0dcr",
			Amount:      cosmos.NewUint(4_30000000),
		},
		{
			MayaAddress: "maya1j42xpqgfdyagr57pxkxgmryzdfy2z4l65mjzf9",
			ThorAddress: "thor1j42xpqgfdyagr57pxkxgmryzdfy2z4l65vvwl4",
			Amount:      cosmos.NewUint(2_00000000),
		},
		{
			MayaAddress: "maya1j6ep9yljeswft03w2qunqx8my9e2efph5ywhls",
			ThorAddress: "thor1jj4xufkxrjd4d3uswh0ztgr0yan3mdcdxh3tgn",
			Amount:      cosmos.NewUint(2_00000000),
		},
		{
			MayaAddress: "maya1jwq4zu4v3tfktwemwh2lwwnlu3nvvrhuhs6k0h",
			ThorAddress: "thor1jwq4zu4v3tfktwemwh2lwwnlu3nvvrhuh8y6e8",
			Amount:      cosmos.NewUint(285_40743565),
		},
		{
			MayaAddress: "maya1ka2v9exy8ata00pch87wgzf9dsmyag94tq8mug",
			ThorAddress: "thor1ka2v9exy8ata00pch87wgzf9dsmyag94theh2c",
			Amount:      cosmos.NewUint(978_00000000),
		},
		{
			MayaAddress: "maya1mj8yhw3jqljfcggkjd77k9t7jlcw0uur0yfurh",
			ThorAddress: "thor1mj8yhw3jqljfcggkjd77k9t7jlcw0uur0nhs48",
			Amount:      cosmos.NewUint(341_00000000),
		},
		{
			MayaAddress: "maya1ppdzsyugtsdtd6dpvzzg2746qfdfmux7k2ydal",
			ThorAddress: "thor1z9xhmhtxn5gxd4ugfuxk7hg9hp03tw3qtqs3f3",
			Amount:      cosmos.NewUint(1_00000000),
		},
		{
			MayaAddress: "maya1qdhqqlg5kcn9hz7wf8wzw8hj8ujrjplnz669c9",
			ThorAddress: "thor1ru7upan5aj2hmzlevrztd6gn5r5z8jxrcjzmup",
			Amount:      cosmos.NewUint(1_00000000),
		},
		{
			MayaAddress: "maya1qtcst64ea585s7gtek3daj2xe59hgn8q7j0ccl",
			ThorAddress: "thor1qtcst64ea585s7gtek3daj2xe59hgn8q7935w0",
			Amount:      cosmos.NewUint(2998_00000000),
		},
	}

	pool, err := mgr.Keeper().GetPool(ctx, common.RUNEAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return
	}

	for _, sender := range lps {
		address, err := common.NewAddress(sender.MayaAddress)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			continue
		}

		lp, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, address)
		if err != nil {
			ctx.Logger().Error("fail to get liquidity provider", "error", err)
			continue
		}

		pool.PendingInboundAsset = pool.PendingInboundAsset.Add(sender.Amount)
		lp.PendingAsset = lp.PendingAsset.Add(sender.Amount)
		lp.LastAddHeight = ctx.BlockHeight()
		if sender.TxID != "" {
			txID, err := common.NewTxID(sender.TxID)
			if err != nil {
				ctx.Logger().Error("fail to parse txID", "error", err)
				continue
			}
			lp.PendingTxID = txID
		}

		if lp.AssetAddress.IsEmpty() {
			thorAdd, err := common.NewAddress(sender.ThorAddress)
			if err != nil {
				ctx.Logger().Error("fail to parse address", "address", sender.MayaAddress, "error", err)
				continue
			}
			lp.AssetAddress = thorAdd
		}

		mgr.Keeper().SetLiquidityProvider(ctx, lp)
		if sender.Tier != 0 {
			if err := mgr.Keeper().SetLiquidityAuctionTier(ctx, lp.CacaoAddress, sender.Tier); err != nil {
				ctx.Logger().Error("fail to set liquidity auction tier", "address", lp.CacaoAddress, "error", err)
				continue
			}
		}

		if err := mgr.Keeper().SetPool(ctx, pool); err != nil {
			ctx.Logger().Error("fail to set pool", "address", pool.Asset, "error", err)
			return
		}

		evt := NewEventPendingLiquidity(pool.Asset, AddPendingLiquidity, lp.CacaoAddress, cosmos.ZeroUint(), lp.AssetAddress, sender.Amount, common.TxID(""), common.TxID(sender.TxID))
		if err := mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
			continue
		}
	}

	// Remove duplicated THOR address LP position
	// https://mayanode.mayachain.info/mayachain/liquidity_auction_tier/thor.rune/maya1dy6c9tmu7qgpd6cw2unumew3sknduwx7s0myr6?height=488436
	// https://mayanode.mayachain.info/mayachain/liquidity_auction_tier/thor.rune/maya1yf0sglxse7jkq0laddtve2fskkrv6vzclu3u6e?height=488436
	add1, err := common.NewAddress("maya1dy6c9tmu7qgpd6cw2unumew3sknduwx7s0myr6")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}

	lp1, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, add1)
	if err != nil {
		ctx.Logger().Error("fail to get liquidity provider", "error", err)
		return
	}

	add2, err := common.NewAddress("maya1yf0sglxse7jkq0laddtve2fskkrv6vzclu3u6e")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}

	lp2, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, add2)
	if err != nil {
		ctx.Logger().Error("fail to get liquidity provider", "error", err)
		return
	}
	lp2.PendingAsset = lp2.PendingAsset.Add(lp1.PendingAsset)

	mgr.Keeper().SetLiquidityProvider(ctx, lp2)
	if err := mgr.Keeper().SetLiquidityAuctionTier(ctx, lp2.CacaoAddress, 0); err != nil {
		ctx.Logger().Error("fail to set liquidity auction tier", "error", err)
	}
	mgr.Keeper().RemoveLiquidityProvider(ctx, lp1)

	// Mint cacao
	toMint := common.NewCoin(common.BaseAsset(), cosmos.NewUint(9_900_000_000_00000000))
	if err := mgr.Keeper().MintToModule(ctx, ModuleName, toMint); err != nil {
		ctx.Logger().Error("fail to mint cacao", "error", err)
		return
	}

	if err = mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, ReserveName, common.NewCoins(toMint)); err != nil {
		ctx.Logger().Error("fail to send cacao to reserve", "error", err)
		return
	}

	// 150214766379119 de BTC Asgard a reserve
	// 473657580023 de ETH Asgard a reserve
	// 24192844274670 de RUNE asgard a reserve
	for _, asset := range []common.Asset{common.BTCAsset, common.ETHAsset, common.RUNEAsset} {
		pool, err := mgr.Keeper().GetPool(ctx, asset)
		if err != nil {
			ctx.Logger().Error("fail to get pool", "error", err)
			return
		}
		switch asset {
		case common.BTCAsset:
			pool.BalanceCacao = pool.BalanceCacao.Sub(cosmos.NewUint(1_501_734_01773759))
		case common.ETHAsset:
			pool.BalanceCacao = pool.BalanceCacao.Sub(cosmos.NewUint(4736_57580023))
		case common.RUNEAsset:
			pool.BalanceCacao = pool.BalanceCacao.Sub(cosmos.NewUint(211_877_34242261))
		}

		if err = mgr.Keeper().SetPool(ctx, pool); err != nil {
			ctx.Logger().Error("fail to set pool", "error", err)
			return
		}
	}

	// Sum of all the above will be sent
	asgardToReserve := common.NewCoin(common.BaseAsset(), cosmos.NewUint(1_717_347_93596043))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, AsgardName, ReserveName, common.NewCoins(asgardToReserve)); err != nil {
		ctx.Logger().Error("fail to send asgard to reserve", "error", err)
		return
	}

	// 164293529917265 de itzamna a reserve
	itzamnaToReserve := common.NewCoin(common.BaseAsset(), cosmos.NewUint(1_642_935_29917265))
	itzamnaAcc, err := cosmos.AccAddressFromBech32("maya18z343fsdlav47chtkyp0aawqt6sgxsh3vjy2vz")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}

	if err := mgr.Keeper().SendFromAccountToModule(ctx, itzamnaAcc, ReserveName, common.NewCoins(itzamnaToReserve)); err != nil {
		ctx.Logger().Error("fail to send itzamna to reserve", "error", err)
		return
	}

	// FROM RESERVE TXS
	// 8_910_000_500_00000000 from reserve to itzamna
	reserveToItzamna := common.NewCoin(common.BaseAsset(), cosmos.NewUint(8_910_001_000_00000000))
	if err := mgr.Keeper().SendFromModuleToAccount(ctx, ReserveName, itzamnaAcc, common.NewCoins(reserveToItzamna)); err != nil {
		ctx.Logger().Error("fail to send reserve to itzamna", "error", err)
		return
	}

	// Remove Slash points from genesis nodes
	for _, genesis := range GenesisNodes {
		acc, err := cosmos.AccAddressFromBech32(genesis)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			continue
		}

		mgr.Keeper().ResetNodeAccountSlashPoints(ctx, acc)
	}
}

func migrateStoreV104(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v104", "error", err)
		}
	}()

	// Select the least secure ActiveVault Asgard for all outbounds.
	// Even if it fails (as in if the version changed upon the keygens-complete block of a churn),
	// updating the voter's FinalisedHeight allows another MaxOutboundAttempts for LackSigning vault selection.
	activeAsgards, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil || len(activeAsgards) == 0 {
		ctx.Logger().Error("fail to get active asgard vaults", "error", err)
		return
	}
	if len(activeAsgards) > 1 {
		signingTransactionPeriod := mgr.GetConstants().GetInt64Value(constants.SigningTransactionPeriod)
		activeAsgards = mgr.Keeper().SortBySecurity(ctx, activeAsgards, signingTransactionPeriod)
	}
	vaultPubKey := activeAsgards[0].PubKey

	// Refund failed synth swaps back to users
	// These swaps were refunded because the target amount set by user was higher than the swap output
	// but because there were a bug in calculating the fee of synth swaps they were treated as zombie coins,
	// and thus we failed to generate the out tx of refund. (keep in mind that the refund event is emitted)
	// Since they are all inbound transactions, we can refund them back to users without deducting fee (see refundTransactions implementation)
	failedSwaps := []adhocRefundTx{
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      300000000,
			inboundHash: "86AC0A216FA3138E3B1EE15D66DEBCBE46D8A62B45EA6D33E07DE044D4BD638E",
		}, {
			toAddr:      "maya1x5979k5wqgq58f4864glr7w2rtgyuqqm6l2zhx",
			asset:       "THOR/RUNE",
			amount:      26142918750,
			inboundHash: "9FC3C8886CD432338B4E4A388DF718B3EE03B257CA2D87792A9D3AFE4AC76DA6",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5047059000,
			inboundHash: "86964E9623839AEBD7D4E74CC777F917AC5DACA850B322F07E7CD6F9A8ACEC1F",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5273550000,
			inboundHash: "1A9D4E7000FE5EF4E292378F1EA075D69DE4DAF2FD5258AC5C2C6E495F28B843",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5366524000,
			inboundHash: "7BADEBA845A889750BF9477B8A01870F109EAAE46E55EF032EF868540F6DB4C1",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      6347294000,
			inboundHash: "C71D1260FD4FE208CFA70440847250716F8C852674592B55EDA390FF840E1C8E",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      6347446000,
			inboundHash: "703E35046FC628B48CA06DD8FE9A95151ECC447C9A55FA7172CC6ED0F97540C4",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5389962000,
			inboundHash: "B1F5B7C9B8AA46A96A72D1E10BD083172810669B912E46BFF2713B4D6237C42C",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5390043000,
			inboundHash: "D55E1265E68D4605118EC02ADE7FB2FD2A91AD35878E701093D0C82B8D624A04",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      1271000,
			inboundHash: "A20F86A30BED39CFC4734EEA1C50680CC32002974E2FE5CE82BB22B26643D618",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      1271000,
			inboundHash: "F3FE1EFE4181E1F81048CCCD366A0E624A98C8C9ED9DC304E3DC32BF2FD3050D",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      1272000,
			inboundHash: "11BFC34721FEBE40CD080432B379F4D9C43DCA147653AFF82849B82838C1B4FD",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      621000,
			inboundHash: "31722DE22A5243DAB294529F3323B6708E8B3040C0205D5602F4F3F5D4218712",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      847000,
			inboundHash: "FFBBAE0420A1F7D1371F837BCA89D697EAC6E7D90835767C70D5B05584F95CD1",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      237000,
			inboundHash: "AD516C4C23A984336DC2BEE188CD7B607F31E342192BD0D05371A0B1AC127234",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      237000,
			inboundHash: "42C110B63651F47B10066C47334296E9E28E006A6481B675A8DAA27946843B81",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      429000,
			inboundHash: "EC192B4327CE11A03611FB5EFEEF3E133C8937040B36B4312A1743BACB4FFA88",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      384000,
			inboundHash: "A43B1F63D2B3092B80F0321DFFE81179BD3EE7209B1EA035D573A83F68EB7177",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      384000,
			inboundHash: "811035D9F7A177199F2BD84B90F82477AC68B2898D0F99B9EEB524766AC914DD",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      429000,
			inboundHash: "F5A01D066C001EB138E8DC4FA21B36917FCA8DCB07289B0E8575FD9B500C4C59",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      429000,
			inboundHash: "EB9F7D4AB93920447EC3423A8D4F1E92C43AFC42B9D18C1362EC5322DFB5ADDB",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      265000,
			inboundHash: "0C1359A466FDB89F450D02AB5C36F1073179D9333D05616FC33D8946058498BC",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      310000,
			inboundHash: "6AB46CD16A3BD9015570C2CD086CEF7BF75ADEBD98C7732149024C45F8458602",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      237000,
			inboundHash: "A45D65ADC584C687DBB696DD10E4E7E9FDC1E81FBA5525BEBF978A03EA2B93BD",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      310000,
			inboundHash: "4554C137395882EB69F01017D64993CCDCA7263AAFB755ADBFE5FDA6A00AB8A4",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "86567487D4E2B1E05B1A86EFE7A7A548849B8F82E17A8325484F855B92633D9B",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "83767F9F779657C0D770975ADFF3A92BCD1EE7A1C999985EA9DB066FBB44610F",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "34DE850E4F1AEFB5EE32C9FA2446D85B76F531E3006463B5F390570439FD96DE",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "18308912F70B0C58FF53BFC7617514C1279EDE4537FA29C822CA3801BBF7C82B",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      265000,
			inboundHash: "7C0EAF6ADE8B9A6DEF3CCF2F462826DDEEFD71C60B29023CE107115616B614BC",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      265000,
			inboundHash: "45B2F0E6338DCD7799026821D3F86A1C794F80A19ED5EE8CA3E5EC649A194B4C",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      265000,
			inboundHash: "2E22F59FC3B69CF871A411B2A057FDCA2DF00469819EFB2A946BE05E22373362",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      310000,
			inboundHash: "3DB087973B69B3A32EA4FA5B16579B3C2EEE6A0E070C03EF8DDE578A12B399FD",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      265000,
			inboundHash: "841A7B59C3E20A58A20C939BCD45800E69D2316F74B84F990B9DCC2E5D43D632",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "E394B2D44226421FC4FCCAD3C0F58D80EE6FE3F70B93E5FA1B699923EBF73588",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "79122693813DB28FC79D9454FA4327523ADAF71F4675BD36BF677301A568090F",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "34068637E11234A6AFC0C85A318CB58666623A8E36626D4E265B10551E1C7166",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "05AEE74CB1B9A2AD4CA3BAC63DB4D4ADA0ADF1EB345D6BC94CCCD0672669E1F8",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "BDBA06245E439EA80C1DEB8295449DF6EF3FC22D0F4D64FAFF3C0095D64413CB",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "A99BAE4BB38CA491D8457D16F2577144879285D97A44E3624848D70F1FD5963B",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      502000,
			inboundHash: "4541912B35D1D6A27A6263B6A7E608AFEB1687984765B2669DAE2612188AD4B9",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      846000,
			inboundHash: "59E5738F0DDB3B1D7F3DF8EFBD691C61480BDB909B57E971C5FFF03054A0EC3B",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      846000,
			inboundHash: "C48B393A050D7983F83C429BD3D17E2C300FAC2EFCB5B18F45E68E42952BC126",
		}, {
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      20021610000,
			inboundHash: "48D8788931772A5566C922C098C579FBBEE2B2793057B487FE3AE2AC2F3C8ED9",
		}, {
			toAddr:      "maya1x5979k5wqgq58f4864glr7w2rtgyuqqm6l2zhx",
			asset:       "ETH/USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
			amount:      82605080070,
			inboundHash: "E0320F7459B83A9F86695C9D0DB78B916F69FF5E408F94E521889F3F0C3CE086",
		},
	}
	refundTransactions(ctx, mgr, vaultPubKey.String(), failedSwaps...)

	// 1st user tier fix
	// User with address "maya1dy6c9tmu7qgpd6cw2unumew3sknduwx7s0myr6" and "maya1yf0sglxse7jkq0laddtve2fskkrv6vzclu3u6e" which had
	// should have been allocated an amount during the cacao donation in the last store migration but seems that
	// there was a problem with the migration and the amount was not allocated. So we change his/her tier to 1
	// and allocate the attribution amount manually from reserve.
	// The changes are as the following:
	// 1. Change Tier from 0 -> 1
	// 2. Overwrite LP Units from 0 -> 38089_5898484080 LP Units
	// 3. Pending Asset from 3210_34000000 -> 0
	// 4. Asset Deposit Value from 0 -> 3273_80071698
	// 5. Cacao Deposit Value from 0 -> 38089_5898484080 (Same as LP Units)
	// 6. Move 38827_9343263458 CACAO from Reserve to Asgard module (CACAO deposit value + Change difference between asset deposit value and pending asset with CACAO denom)
	// 7. Increase by 38827_9343263458 the CACAO on Asgard Pool for RUNE (CACAO deposit value + Change difference between asset deposit value and pending asset with CACAO denom)
	// 8. Increase by 38089_5898484080 the LP UNITS on Asgard Pool for RUNE
	// 9. Move 3210_34000000 Asset from Pending_Asset in Asgard Pool for RUNE to Balance_Asset in Asgard Pool for RUNE
	// 10. Emit Add Liquidity Event
	addr1, err := common.NewAddress("maya1yf0sglxse7jkq0laddtve2fskkrv6vzclu3u6e")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}
	lp1, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, addr1)
	if err != nil {
		ctx.Logger().Error("fail to get liquidity provider", "error", err)
		return
	}
	lp1.Units = cosmos.NewUint(38089_5898484080)
	lp1.PendingAsset = cosmos.ZeroUint()
	lp1.AssetDepositValue = cosmos.NewUint(3273_80071698)
	lp1.CacaoDepositValue = cosmos.NewUint(38089_5898484080)
	mgr.Keeper().SetLiquidityProvider(ctx, lp1)
	if err := mgr.Keeper().SetLiquidityAuctionTier(ctx, lp1.CacaoAddress, 1); err != nil {
		ctx.Logger().Error("fail to set liquidity auction tier", "error", err)
	}

	reserve2Asgard1 := common.NewCoin(common.BaseAsset(), cosmos.NewUint(38827_9343263458))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, AsgardName, common.NewCoins(reserve2Asgard1)); err != nil {
		ctx.Logger().Error("fail to send reserve to asgard", "error", err)
		return
	}
	pool, err := mgr.Keeper().GetPool(ctx, common.RUNEAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return
	}
	addedCacao1 := cosmos.NewUint(38827_9343263458)
	pool.BalanceCacao = pool.BalanceCacao.Add(addedCacao1)
	addedLPUnits := cosmos.NewUint(38089_5898484080)
	pool.LPUnits = pool.LPUnits.Add(addedLPUnits)
	pendingAsset2Balance := cosmos.NewUint(3210_34000000)
	pool.PendingInboundAsset = pool.PendingInboundAsset.Sub(pendingAsset2Balance)
	pool.BalanceAsset = pool.BalanceAsset.Add(pendingAsset2Balance)
	evt1 := NewEventAddLiquidity(
		pool.Asset,
		addedLPUnits,
		lp1.CacaoAddress,
		addedCacao1,
		pendingAsset2Balance,
		common.TxID(""),
		common.TxID(""),
		lp1.AssetAddress,
	)

	// 2nd user tier fix
	// 1. Change Tier from 3 -> 1
	// 2. Increase Asset Deposit Value from 333_46986565 to 507_5548724869
	// 3. Increase Cacao Deposit Value from 3879_8117483183 to 5905_2332716199
	// 4. Increase LP UNITS from 3879_8117483183 to 5905_2332716199
	// 5. Move 4050_84304660322 CACAO from Reserve module to Asgard module (twice as much on purpose, to account for asset side, will be armed away)
	// 6. Increase in 4050_84304660322 CACAO the balance_cacao of Asgard Pool for RUNE (twice as much on purpose, to account for asset side, will be arbed away)
	// 7. Increase by 2025_4215233016 the LP UNITS on Asgard Pool for RUNE
	// 8. Emit Add Liquidity Event
	addr2, err := common.NewAddress("maya1jwq4zu4v3tfktwemwh2lwwnlu3nvvrhuhs6k0h")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}
	lp2, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, addr2)
	if err != nil {
		ctx.Logger().Error("fail to get liquidity provider", "error", err)
		return
	}
	lp2.AssetDepositValue = cosmos.NewUint(507_5548724869)
	lp2.CacaoDepositValue = cosmos.NewUint(5905_2332716199)
	lp2.Units = cosmos.NewUint(5905_2332716199)
	mgr.Keeper().SetLiquidityProvider(ctx, lp2)
	if err := mgr.Keeper().SetLiquidityAuctionTier(ctx, lp2.CacaoAddress, 1); err != nil {
		ctx.Logger().Error("fail to set liquidity auction tier", "error", err)
	}

	reserve2Asgard2 := common.NewCoin(common.BaseAsset(), cosmos.NewUint(4050_84304660322))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, AsgardName, common.NewCoins(reserve2Asgard2)); err != nil {
		ctx.Logger().Error("fail to send reserve to asgard", "error", err)
		return
	}
	addedCacao2 := cosmos.NewUint(4050_84304660322)
	pool.BalanceCacao = pool.BalanceCacao.Add(addedCacao2)
	addedLPUnits2 := cosmos.NewUint(2025_4215233016)
	pool.LPUnits = pool.LPUnits.Add(addedLPUnits2)
	evt2 := NewEventAddLiquidity(
		pool.Asset,
		addedLPUnits2,
		lp2.CacaoAddress,
		addedCacao2,
		cosmos.ZeroUint(),
		common.TxID(""),
		common.TxID(""),
		common.Address(""),
	)

	err = mgr.Keeper().SetPool(ctx, pool)
	if err != nil {
		ctx.Logger().Error("fail to set pool", "error", err)
		return
	}
	if err := mgr.EventMgr().EmitEvent(ctx, evt1); err != nil {
		ctx.Logger().Error("fail to emit event", "error", err)
		return
	}
	if err := mgr.EventMgr().EmitEvent(ctx, evt2); err != nil {
		ctx.Logger().Error("fail to emit event", "error", err)
		return
	}
}

// migrateStoreV105 is complementory migration to migration v104
// it will refund another 17 failed synth swaps txs back to users
func migrateStoreV105(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v105", "error", err)
		}
	}()

	// Select the least secure ActiveVault Asgard for all outbounds.
	// Even if it fails (as in if the version changed upon the keygens-complete block of a churn),
	// updating the voter's FinalisedHeight allows another MaxOutboundAttempts for LackSigning vault selection.
	activeAsgards, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil || len(activeAsgards) == 0 {
		ctx.Logger().Error("fail to get active asgard vaults", "error", err)
		return
	}
	if len(activeAsgards) > 1 {
		signingTransactionPeriod := mgr.GetConstants().GetInt64Value(constants.SigningTransactionPeriod)
		activeAsgards = mgr.Keeper().SortBySecurity(ctx, activeAsgards, signingTransactionPeriod)
	}
	vaultPubKey := activeAsgards[0].PubKey

	// Refund failed synth swaps back to users
	// These swaps were refunded because the target amount set by user was higher than the swap output
	// but because there were a bug in calculating the fee of synth swaps they were treated as zombie coins,
	// and thus we failed to generate the out tx of refund. (keep in mind that the refund event is emitted)
	// Since they are all inbound transactions, we can refund them back to users without deducting fee (see refundTransactions implementation)
	failedSwaps := []adhocRefundTx{
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      6175167000,
			inboundHash: "8EECEE5C27795B96E8465D3234DEC050219AC591D899D038D2F11A1EFCE00E72",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      1271000,
			inboundHash: "E31EBA09AA7E64DE5F1209656956286C4883196B0E85A075764600ABC57ACDB6",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      384000,
			inboundHash: "6777B04215485FC495A88FA5D76C1873E250756FFF5E23577CA3CEEB4E042B0C",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      384000,
			inboundHash: "E34798700D6034A3D8C82F80E7FCC4AC0F68574FCB7FD018EFA7E90A2594A44F",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      384000,
			inboundHash: "2D129E0E58A762263272DB2548B432912E995F2A09CFF4A6C06A4DF8534290C7",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      846000,
			inboundHash: "0172C67339320D14E477DCEB64F9FC4FABEE67DF233F08A81EB4D061F1820AC1",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "BTC/BTC",
			amount:      846000,
			inboundHash: "86E8363E44B4EF0B32A894FD3011AC6AB8EC7AAE3EA2F65ACD8D0D15DB1299C7",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5314334000,
			inboundHash: "DBBDA76A5315F25787041BF95A65FC19BD2464B637BA9ED322CD8A52C1CE447E",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5927862000,
			inboundHash: "8D40B6E45B676638764FB38A998FAD782514AF2DDDB840A809A6CB65C854DF70",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      8527629000,
			inboundHash: "A860164DFDC3B0E76B871FF93A509B80736486B622AE59D0EE77ECE5F0E39D6A",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      8527450000,
			inboundHash: "3B99680D1927C6A3D909B964378443E8D5C71F9DA2A3E7FF4AFF16C7B6E08FA5",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      20123240000,
			inboundHash: "53BA8317F50DFB97FF30235BA479F3E3F78E29FEFA90BB2C113891F121D79C04",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      7565885000,
			inboundHash: "DD862F71E427F5DE280F2CAA49E007B77A1B64E30896AA93B1EC782374CDAB04",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      20107080000,
			inboundHash: "62D7648EC776A7B68FFAB23844EB5AA2C967F7E7CA97379E03607682B312B33E",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "THOR/RUNE",
			amount:      20103040000,
			inboundHash: "BD90238148001ADF4B485D98D549AB472F5AC1881E8F67DBA70CC5C80E979803",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5252251000,
			inboundHash: "D028821B72FD5A37C092771FF9F5039C7A7E04FFDDEF8793A0FFE0BD73156733",
		},
		{
			toAddr:      "maya1gyap83aenguyhce3a0y3gprap32ypuc99vtzlc",
			asset:       "ETH/USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
			amount:      5060695000,
			inboundHash: "89266DF89E689C79DE4ACCB7312FEC85CE57CF92852A222228C3388F6FBDDA57",
		},
		{
			toAddr:      "maya1x5979k5wqgq58f4864glr7w2rtgyuqqm6l2zhx",
			asset:       "THOR/RUNE",
			amount:      72494713125,
			inboundHash: "8671D17BFD6040531470C89D0412116EE2909396BB6C54E037535DFD529E67D2",
		},
	}
	refundTransactions(ctx, mgr, vaultPubKey.String(), failedSwaps...)

	// Refunding USDT coins that mistakenly got sent to the vault (mayapub1addwnpepqwuwsax7p3raecsn2k9uvqyykanlvhw47asz836se2h0nyg6knug6n9hklq) by "transfer" txs back to user
	// transaction hashes are: 0xda4306037c838dcaed92775ecd515441e4a932b1bcbeef1199bf37a29274575d and 0xa6d765192856e982deae51bfc817f612c30344402ca72fbe526e8c534b91d048 on eth mainnet
	maxGas, err := mgr.GasMgr().GetMaxGas(ctx, common.ETHChain)
	if err != nil {
		ctx.Logger().Error("fail to get max gas", "error", err)
		return
	}
	toi := TxOutItem{
		Chain:       common.ETHChain,
		InHash:      common.BlankTxID,
		ToAddress:   common.Address("0x2510d455bF4a9b829C0CfD579543918D793F7762"),
		Coin:        common.NewCoin(common.USDTAsset, cosmos.NewUint(191_970_000+96_448_216)),
		MaxGas:      common.Gas{maxGas},
		GasRate:     int64(mgr.GasMgr().GetGasRate(ctx, common.ETHChain).Uint64()),
		VaultPubKey: common.PubKey("mayapub1addwnpepqwuwsax7p3raecsn2k9uvqyykanlvhw47asz836se2h0nyg6knug6n9hklq"),
	}
	if err := mgr.TxOutStore().UnSafeAddTxOutItem(ctx, mgr, toi); err != nil {
		ctx.Logger().Error("fail to save tx out item for refund transfers", "error", err)
		return
	}
}
