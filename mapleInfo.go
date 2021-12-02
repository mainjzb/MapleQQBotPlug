package main

var TotalExpCollection = []struct {
	Level int64
	Exp   int64
}{
	{210, 50192858013},
	{220, 226834057694},
	{230, 888805728115},
	{240, 2780379685705},
	{250, 7764451421743},
	{255, 14465974056466},
	{260, 21509339594499},
	{265, 36314635528483},
	{270, 51875150349788},
	{275, 84583665273612},
	{280, 166871327183154},
	{285, 410235606102699},
	{290, 1129981080450215},
	{295, 3258615365414436},
	{300, 10103007996648098},
}

var AllClass = []struct {
	Names []string
}{
	{[]string{"等级", "", "排名"}},
	{[]string{"联盟"}},
	{[]string{"Warrior", "战士", "英雄", "黑骑", "圣骑"}},
	{[]string{"Magician", "法师", "冰雷", "火毒", "主教"}},
	{[]string{"Bowman", "弓手", "弓箭手", "神射手", "神射", "弩手"}},
	{[]string{"Thief", "飞侠", "双刀", "刀飞", "标飞"}},
	{[]string{"Pirate", "海盗", "船长", "队长", "火炮"}},
	{[]string{"Aran", "战神"}},
	{[]string{"Evan", "龙神"}},
	{[]string{"Mercedes", "双弩"}},
	{[]string{"Phantom", "幻影"}},
	{[]string{"Jett", "杰特"}},
	{[]string{"Mihile", "米哈尔", "米哈哈"}},
	{[]string{"Luminous", "夜光", "夜光法师"}},
	{[]string{"Kaiser", "凯撒", "狂龙"}},
	{[]string{"Angelic Buster", "天使", "ab", "AB"}},
	{[]string{"Hayato", "剑豪"}},
	{[]string{"Kanna", "阴阳师", "娜神"}},
	{[]string{"Xenon", "煎饼", "尖兵"}},
	{[]string{"Zero", "神之子", "神子"}},
	{[]string{"Beast Tamer", "BT", "bt", "林志玲", "林之灵", "lzl"}},
	{[]string{"Shade", "隐月"}},
	{[]string{"Kinesis", "超能", "超能力者", "超能力"}},
	{[]string{"Blaster", "爆破", "爆破者"}},
	{[]string{"Cadena", "卡姐", "卡德娜"}},
	{[]string{"Illium", "黑皮", "圣经使徒", "圣晶使徒"}},
	{[]string{"Ark", "牙科", "亚克"}},
	{[]string{"Pathfinder", "pf", "开拓者", "古迹猎人", "PF"}},
	{[]string{"Hoyoung", "虎影"}},
	{[]string{"Adele", "阿呆", "阿黛尔"}},
	{[]string{"Dawn Warrior", "dw", "DW", "魂骑"}},
	{[]string{"Blaze Wizard", "BW", "bw", "炎术士"}},
	{[]string{"Wind Archer", "WA", "wa", "风铃"}},
	{[]string{"Night Walker", "NW", "nw", "夜行", "夜行者"}},
	{[]string{"Thunder Breaker", "TB", "tb", "奇袭", "奇袭者"}},
	{[]string{"Demon Slayer", "DS", "ds", "红毛"}},
	{[]string{"Battle Mage", "BM", "bm", "战法", "幻灵", "唤灵"}},
	{[]string{"Wild Hunter", "WH", "wh", "豹弩"}},
	{[]string{"Mechanicr", "轮椅", "机械"}},
	{[]string{"Demon Avenger", "白毛", "DA"}},
	{[]string{"Kain", "卡因", "卡隐", "卡影"}},
	{[]string{"Lara", "lara", "lala", "拉拉"}},
}

var classURL = map[string]string{
	"Warrior":         "https://maplestory.nexon.net/api/ranking?id=job&id2=1&rebootIndex=1&page_index=",
	"Magician":        "https://maplestory.nexon.net/api/ranking?id=job&id2=2&rebootIndex=1&page_index=",
	"Bowman":          "https://maplestory.nexon.net/api/ranking?id=job&id2=3&rebootIndex=1&page_index=",
	"Thief":           "https://maplestory.nexon.net/api/ranking?id=job&id2=4&rebootIndex=1&page_index=",
	"Pirate":          "https://maplestory.nexon.net/api/ranking?id=job&id2=5&rebootIndex=1&page_index=",
	"Aran":            "https://maplestory.nexon.net/api/ranking?id=job&id2=21&rebootIndex=1&page_index=",
	"Evan":            "https://maplestory.nexon.net/api/ranking?id=job&id2=22&rebootIndex=1&page_index=",
	"Mercedes":        "https://maplestory.nexon.net/api/ranking?id=job&id2=23&rebootIndex=1&page_index=",
	"Phantom":         "https://maplestory.nexon.net/api/ranking?id=job&id2=24&rebootIndex=1&page_index=",
	"Jett":            "https://maplestory.nexon.net/api/ranking?id=job&id2=201&rebootIndex=1&page_index=",
	"Mihile":          "https://maplestory.nexon.net/api/ranking?id=job&id2=202&rebootIndex=1&page_index=",
	"Luminous":        "https://maplestory.nexon.net/api/ranking?id=job&id2=203&rebootIndex=1&page_index=",
	"Kaiser":          "https://maplestory.nexon.net/api/ranking?id=job&id2=204&rebootIndex=1&page_index=",
	"Angelic Buster":  "https://maplestory.nexon.net/api/ranking?id=job&id2=205&rebootIndex=1&page_index=",
	"Hayato":          "https://maplestory.nexon.net/api/ranking?id=job&id2=206&rebootIndex=1&page_index=",
	"Kanna":           "https://maplestory.nexon.net/api/ranking?id=job&id2=207&rebootIndex=1&page_index=",
	"Xenon":           "https://maplestory.nexon.net/api/ranking?id=job&id2=208&rebootIndex=1&page_index=",
	"Zero":            "https://maplestory.nexon.net/api/ranking?id=job&id2=210&rebootIndex=1&page_index=",
	"Beast Tamer":     "https://maplestory.nexon.net/api/ranking?id=job&id2=211&rebootIndex=1&page_index=",
	"Shade":           "https://maplestory.nexon.net/api/ranking?id=job&id2=212&rebootIndex=1&page_index=",
	"Kinesis":         "https://maplestory.nexon.net/api/ranking?id=job&id2=214&rebootIndex=1&page_index=",
	"Blaster":         "https://maplestory.nexon.net/api/ranking?id=job&id2=215&rebootIndex=1&page_index=",
	"Cadena":          "https://maplestory.nexon.net/api/ranking?id=job&id2=216&rebootIndex=1&page_index=",
	"Illium":          "https://maplestory.nexon.net/api/ranking?id=job&id2=217&rebootIndex=1&page_index=",
	"Ark":             "https://maplestory.nexon.net/api/ranking?id=job&id2=218&rebootIndex=1&page_index=",
	"Pathfinder":      "https://maplestory.nexon.net/api/ranking?id=job&id2=219&rebootIndex=1&page_index=",
	"Hoyoung":         "https://maplestory.nexon.net/api/ranking?id=job&id2=220&rebootIndex=1&page_index=",
	"Adele":           "https://maplestory.nexon.net/api/ranking?id=job&id2=221&rebootIndex=1&page_index=",
	"Dawn Warrior":    "https://maplestory.nexon.net/api/ranking?id=job&id2=11&rebootIndex=1&page_index=",
	"Blaze Wizard":    "https://maplestory.nexon.net/api/ranking?id=job&id2=12&rebootIndex=1&page_index=",
	"Wind Archer":     "https://maplestory.nexon.net/api/ranking?id=job&id2=13&rebootIndex=1&page_index=",
	"Night Walker":    "https://maplestory.nexon.net/api/ranking?id=job&id2=14&rebootIndex=1&page_index=",
	"Thunder Breaker": "https://maplestory.nexon.net/api/ranking?id=job&id2=15&rebootIndex=1&page_index=",
	"Demon Slayer":    "https://maplestory.nexon.net/api/ranking?id=job&id2=31&rebootIndex=1&page_index=",
	"Battle Mage":     "https://maplestory.nexon.net/api/ranking?id=job&id2=32&rebootIndex=1&page_index=",
	"Wild Hunter":     "https://maplestory.nexon.net/api/ranking?id=job&id2=33&rebootIndex=1&page_index=",
	"Mechanicr":       "https://maplestory.nexon.net/api/ranking?id=job&id2=34&rebootIndex=1&page_index=",
	"Demon Avenger":   "https://maplestory.nexon.net/api/ranking?id=job&id2=209&rebootIndex=1&page_index=",
	"Kain":            "https://maplestory.nexon.net/api/ranking?id=job&id2=222&rebootIndex=1&page_index=",
	"Lara":            "https://maplestory.nexon.net/api/ranking?id=job&id2=223&rebootIndex=1&page_index",
	"联盟":              "https://maplestory.nexon.net/api/ranking?id=legion&id2=45&page_index=",
	"等级":              "https://maplestory.nexon.net/api/ranking?id=world&id2=45&rebootIndex=0&page_index=",
}

var LevelExp = map[int64]int64{
	1:   15,
	2:   34,
	3:   57,
	4:   92,
	5:   135,
	6:   372,
	7:   560,
	8:   840,
	9:   1242,
	10:  1242,
	11:  1242,
	12:  1242,
	13:  1242,
	14:  1242,
	15:  1490,
	16:  1788,
	17:  2145,
	18:  2574,
	19:  3088,
	20:  3705,
	21:  4446,
	22:  5335,
	23:  6402,
	24:  7682,
	25:  9218,
	26:  11061,
	27:  13273,
	28:  15927,
	29:  19112,
	30:  19112,
	31:  19112,
	32:  19112,
	33:  19112,
	34:  19112,
	35:  22934,
	36:  27520,
	37:  33024,
	38:  39628,
	39:  47553,
	40:  51357,
	41:  55465,
	42:  59902,
	43:  64694,
	44:  69869,
	45:  75458,
	46:  81494,
	47:  88013,
	48:  95054,
	49:  102658,
	50:  110870,
	51:  119739,
	52:  129318,
	53:  139663,
	54:  150836,
	55:  162902,
	56:  175934,
	57:  190008,
	58:  205208,
	59:  221624,
	60:  221624,
	61:  221624,
	62:  221624,
	63:  221624,
	64:  221624,
	65:  238245,
	66:  256113,
	67:  275321,
	68:  295970,
	69:  318167,
	70:  342029,
	71:  367681,
	72:  395257,
	73:  424901,
	74:  456768,
	75:  488741,
	76:  522952,
	77:  559558,
	78:  598727,
	79:  640637,
	80:  685481,
	81:  733464,
	82:  784806,
	83:  839742,
	84:  898523,
	85:  961419,
	86:  1028718,
	87:  1100728,
	88:  1177778,
	89:  1260222,
	90:  1342136,
	91:  1429374,
	92:  1522283,
	93:  1621231,
	94:  1726611,
	95:  1838840,
	96:  1958364,
	97:  2085657,
	98:  2221224,
	99:  2365603,
	100: 2365603,
	101: 2365603,
	102: 2365603,
	103: 2365603,
	104: 2365603,
	105: 2519367,
	106: 2683125,
	107: 2857528,
	108: 3043267,
	109: 3241079,
	110: 3451749,
	111: 3676112,
	112: 3915059,
	113: 4169537,
	114: 4440556,
	115: 4729192,
	116: 5036589,
	117: 5363967,
	118: 5712624,
	119: 6083944,
	120: 6479400,
	121: 6900561,
	122: 7349097,
	123: 7826788,
	124: 8335529,
	125: 8877338,
	126: 9454364,
	127: 10068897,
	128: 10723375,
	129: 11420394,
	130: 12162719,
	131: 12953295,
	132: 13795259,
	133: 14691950,
	134: 15646926,
	135: 16663976,
	136: 17747134,
	137: 18900697,
	138: 20129242,
	139: 21437642,
	140: 22777494,
	141: 24201087,
	142: 25713654,
	143: 27320757,
	144: 29028304,
	145: 30842573,
	146: 32770233,
	147: 34818372,
	148: 36994520,
	149: 39306677,
	150: 41763344,
	151: 44373553,
	152: 47146900,
	153: 50093581,
	154: 53224429,
	155: 56550955,
	156: 60085389,
	157: 63840725,
	158: 67830770,
	159: 72070193,
	160: 76574580,
	161: 81360491,
	162: 86445521,
	163: 91848366,
	164: 97588888,
	165: 103688193,
	166: 110168705,
	167: 117054249,
	168: 124370139,
	169: 132143272,
	170: 138750435,
	171: 145687956,
	172: 152972353,
	173: 160620970,
	174: 168652018,
	175: 177084618,
	176: 185938848,
	177: 195235790,
	178: 204997579,
	179: 215247457,
	180: 226009829,
	181: 237310320,
	182: 249175836,
	183: 261634627,
	184: 274716358,
	185: 288452175,
	186: 302874783,
	187: 318018522,
	188: 333919448,
	189: 350615420,
	190: 368146191,
	191: 386553500,
	192: 405881175,
	193: 426175233,
	194: 447483994,
	195: 469858193,
	196: 493351102,
	197: 518018657,
	198: 543919589,
	199: 571115568,
	200: 2207026470,
	201: 2471869646,
	202: 2768494003,
	203: 3100713283,
	204: 3472798876,
	205: 3889534741,
	206: 4356278909,
	207: 4879032378,
	208: 5464516263,
	209: 6120258214,
	210: 9792413142,
	211: 10869578587,
	212: 12065232231,
	213: 13392407776,
	214: 14865572631,
	215: 19325244420,
	216: 21064516417,
	217: 22960322894,
	218: 25026751954,
	219: 27279159629,
	220: 43646655406,
	221: 46701921284,
	222: 49971055773,
	223: 53469029677,
	224: 57211861754,
	225: 74375420280,
	226: 78094191294,
	227: 81998900858,
	228: 86098845900,
	229: 90403788195,
	230: 144646061112,
	231: 148985442945,
	232: 153455006233,
	233: 158058656419,
	234: 162800416111,
	235: 211640540944,
	236: 217989757172,
	237: 224529449887,
	238: 231265333383,
	239: 238203293384,
	240: 381125269414,
	241: 392559027496,
	242: 404335798320,
	243: 416465872269,
	244: 428959848437,
	245: 557647802968,
	246: 574377237057,
	247: 591608554168,
	248: 609356810793,
	249: 627637515116,
	250: 1313764762354,
	251: 1326902409977,
	252: 1340171434076,
	253: 1353573148416,
	254: 1367108879900,
	255: 1380779968699,
	256: 1394587768385,
	257: 1408533646068,
	258: 1422618982528,
	259: 1436845172353,
	260: 2902427248153,
	261: 2931451520634,
	262: 2960766035840,
	263: 2990373696198,
	264: 3020277433159,
	265: 3050480207490,
	266: 3080985009564,
	267: 3111794859659,
	268: 3142912808255,
	269: 3174341936337,
	270: 6412170711400,
	271: 6476292418514,
	272: 6541055342699,
	273: 6606465896125,
	274: 6672530555086,
	275: 13478511721273,
	276: 14826362893400,
	277: 16308999182740,
	278: 17939899101014,
	279: 19733889011115,
	280: 39862455802452,
	281: 43848701382697,
	282: 48233571520966,
	283: 53056928673062,
	284: 58362621540368,
	285: 117892495511543,
	286: 129681745062697,
	287: 142649919568966,
	288: 156914911525862,
	289: 172606402678448,
	290: 348664933410464,
	291: 383531426751510,
	292: 421884569426661,
	293: 464073026369327,
	294: 510480329006259,
	295: 1031170264592640,
	296: 1134287291051900,
	297: 1247716020157090,
	298: 1372487622172800,
	299: 2058731433259200,
	300: 1,
}

var class2Chinese = map[string]string{
	"Magician12":   "火毒",
	"Magician22":   "冰雷",
	"Magician32":   "主教",
	"Warrior12":    "英雄",
	"Warrior22":    "圣骑士",
	"Warrior32":    "黑骑士",
	"Thief12":      "标飞",
	"Thief22":      "刀飞",
	"Dual Blade34": "双刀",
	"Bowman12":     "弓射手",
	"Bowman22":     "弩射手",
	"Pathfinder32": "开拓者",

	"Aran12":     "战神",
	"Phantom12":  "幻影",
	"Mercedes12": "双弩",
	"Luminous12": "夜光",
	"Evan12":     "龙神",

	"Zero12":           "神之子",
	"Illium12":         "圣晶",
	"Xenon12":          "尖兵",
	"Adele12":          "阿黛尔",
	"Kanna12":          "阴阳师",
	"Angelic Buster12": "天使",
	"Mihile12":         "米哈尔",
	"Hoyoung12":        "虎影",
	"Kain12":           "卡因",
	"Demon Avenger22":  "白毛",
}

var serverName2Chinese = map[string]string{
	"Reboot (NA)": "R区",
	"Aurora":      "A区",
	"Bera":        "B区",
	"Scania":      "S区",
	"Elysium":     "E区",
}
