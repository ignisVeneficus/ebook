package ascii

import (
	"strings"
)

func removeNotHungarian(r rune) []rune {
	switch r {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return []rune{r}
	case '_':
		return []rune{'_'}
	case ',':
		return []rune{','}
	case ' ':
		return []rune{' '}
	case '-':
		return []rune{'-'}
	// letters
	case 'A', 'Ⓐ', 'Ａ', 'À', 'Â', 'Ầ', 'Ấ', 'Ẫ', 'Ẩ', 'Ã', 'Ā', 'Ă', 'Ằ', 'Ắ', 'Ẵ', 'Ẳ', 'Ȧ', 'Ǡ', 'Ä', 'Ǟ', 'Ả', 'Å', 'Ǻ', 'Ǎ', 'Ȁ', 'Ȃ', 'Ạ', 'Ậ', 'Ặ', 'Ḁ', 'Ą', 'Ⱥ', 'Ɐ':
		return []rune{'A'}
	case 'Ꜳ':
		return []rune{'A', 'A'}
	case 'Æ', 'Ǽ', 'Ǣ':
		return []rune{'A', 'E'}
	case 'Ꜵ':
		return []rune{'A', 'E'}
	case 'Ꜷ':
		return []rune{'A', 'U'}
	case 'Ꜹ', 'Ꜻ':
		return []rune{'A', 'V'}
	case 'Ꜽ':
		return []rune{'A', 'Y'}
		//hungarian
	case 'Á':
		return []rune{'Á'}
	case 'B', 'Ⓑ', 'Ｂ', 'Ḃ', 'Ḅ', 'Ḇ', 'Ƀ', 'Ƃ', 'Ɓ':
		return []rune{'B'}
	case 'C', 'Ⓒ', 'Ｃ', 'Ć', 'Ĉ', 'Ċ', 'Č', 'Ç', 'Ḉ', 'Ƈ', 'Ȼ', 'Ꜿ':
		return []rune{'C'}
	case 'D', 'Ⓓ', 'Ｄ', 'Ḋ', 'Ď', 'Ḍ', 'Ḑ', 'Ḓ', 'Ḏ', 'Đ', 'Ƌ', 'Ɗ', 'Ɖ', 'Ꝺ':
		return []rune{'D'}
	case 'Ǳ', 'Ǆ':
		return []rune{'D', 'Z'}
	case 'ǲ', 'ǅ':
		return []rune{'D', 'z'}
	case 'E', 'Ⓔ', 'Ｅ', 'È', 'Ê', 'Ề', 'Ế', 'Ễ', 'Ể', 'Ẽ', 'Ē', 'Ḕ', 'Ḗ', 'Ĕ', 'Ė', 'Ë', 'Ẻ', 'Ě', 'Ȅ', 'Ȇ', 'Ẹ', 'Ệ', 'Ȩ', 'Ḝ', 'Ę', 'Ḙ', 'Ḛ', 'Ɛ', 'Ǝ':
		return []rune{'E'}
		//hungarian
	case 'É':
		return []rune{'É'}
	case 'F', 'Ⓕ', 'Ｆ', 'Ḟ', 'Ƒ', 'Ꝼ':
		return []rune{'F'}
	case 'G', 'Ⓖ', 'Ｇ', 'Ǵ', 'Ĝ', 'Ḡ', 'Ğ', 'Ġ', 'Ǧ', 'Ģ', 'Ǥ', 'Ɠ', 'Ꞡ', 'Ᵹ', 'Ꝿ':
		return []rune{'G'}
	case 'H', 'Ⓗ', 'Ｈ', 'Ĥ', 'Ḣ', 'Ḧ', 'Ȟ', 'Ḥ', 'Ḩ', 'Ḫ', 'Ħ', 'Ⱨ', 'Ⱶ', 'Ɥ':
		return []rune{'H'}
	case 'I', 'Ⓘ', 'Ｉ', 'Ì', 'Î', 'Ĩ', 'Ī', 'Ĭ', 'İ', 'Ï', 'Ḯ', 'Ỉ', 'Ǐ', 'Ȉ', 'Ȋ', 'Ị', 'Į', 'Ḭ', 'Ɨ':
		return []rune{'I'}
		//hungarian
	case 'Í':
		return []rune{'Í'}
	case 'J', 'Ⓙ', 'Ｊ', 'Ĵ', 'Ɉ':
		return []rune{'J'}
	case 'K', 'Ⓚ', 'Ｋ', 'Ḱ', 'Ǩ', 'Ḳ', 'Ķ', 'Ḵ', 'Ƙ', 'Ⱪ', 'Ꝁ', 'Ꝃ', 'Ꝅ', 'Ꞣ':
		return []rune{'K'}
	case 'L', 'Ⓛ', 'Ｌ', 'Ŀ', 'Ĺ', 'Ľ', 'Ḷ', 'Ḹ', 'Ļ', 'Ḽ', 'Ḻ', 'Ł', 'Ƚ', 'Ɫ', 'Ⱡ', 'Ꝉ', 'Ꝇ', 'Ꞁ':
		return []rune{'L'}
	case 'Ǉ':
		return []rune{'L', 'J'}
	case 'ǈ':
		return []rune{'L', 'j'}
	case 'M', 'Ⓜ', 'Ｍ', 'Ḿ', 'Ṁ', 'Ṃ', 'Ɱ', 'Ɯ':
		return []rune{'M'}
	case 'N', 'Ⓝ', 'Ｎ', 'Ǹ', 'Ń', 'Ñ', 'Ṅ', 'Ň', 'Ṇ', 'Ņ', 'Ṋ', 'Ṉ', 'Ƞ', 'Ɲ', 'Ꞑ', 'Ꞥ':
		return []rune{'N'}
	case 'Ǌ':
		return []rune{'N', 'J'}
	case 'ǋ':
		return []rune{'N', 'j'}
	case 'O', 'Ⓞ', 'Ｏ', 'Ò', 'Ô', 'Ồ', 'Ố', 'Ỗ', 'Ổ', 'Õ', 'Ṍ', 'Ȭ', 'Ṏ', 'Ō', 'Ṑ', 'Ṓ', 'Ŏ', 'Ȯ', 'Ȱ', 'Ȫ', 'Ỏ', 'Ǒ', 'Ȍ', 'Ȏ', 'Ơ', 'Ờ', 'Ớ', 'Ỡ', 'Ở', 'Ợ', 'Ọ', 'Ộ', 'Ǫ', 'Ǭ', 'Ø', 'Ǿ', 'Ɔ', 'Ɵ', 'Ꝋ', 'Ꝍ':
		return []rune{'O'}
		//hungarian
	case 'Ó':
		return []rune{'Ó'}
	case 'Ö':
		return []rune{'Ö'}
	case 'Ő':
		return []rune{'Ő'}
	case 'Ƣ':
		return []rune{'O', 'I'}
	case 'Ꝏ':
		return []rune{'O', 'O'}
	case 'Ȣ':
		return []rune{'O', 'U'}
	case '', 'Œ':
		return []rune{'O', 'E'}
	case 'P', 'Ⓟ', 'Ｐ', 'Ṕ', 'Ṗ', 'Ƥ', 'Ᵽ', 'Ꝑ', 'Ꝓ', 'Ꝕ':
		return []rune{'P'}
	case 'Q', 'Ⓠ', 'Ｑ', 'Ꝗ', 'Ꝙ', 'Ɋ':
		return []rune{'Q'}
	case 'R', 'Ⓡ', 'Ｒ', 'Ŕ', 'Ṙ', 'Ř', 'Ȑ', 'Ȓ', 'Ṛ', 'Ṝ', 'Ŗ', 'Ṟ', 'Ɍ', 'Ɽ', 'Ꝛ', 'Ꞧ', 'Ꞃ':
		return []rune{'R'}
	case 'S', 'Ⓢ', 'Ｓ', 'ẞ', 'Ś', 'Ṥ', 'Ŝ', 'Ṡ', 'Š', 'Ṧ', 'Ṣ', 'Ṩ', 'Ș', 'Ş', 'Ȿ', 'Ꞩ', 'Ꞅ':
		return []rune{'S'}
	case 'T', 'Ⓣ', 'Ｔ', 'Ṫ', 'Ť', 'Ṭ', 'Ț', 'Ţ', 'Ṱ', 'Ṯ', 'Ŧ', 'Ƭ', 'Ʈ', 'Ⱦ', 'Ꞇ':
		return []rune{'T'}
	case 'Ꜩ':
		return []rune{'T', 'Z'}
	case 'U', 'Ⓤ', 'Ｕ', 'Ù', 'Û', 'Ũ', 'Ṹ', 'Ū', 'Ṻ', 'Ŭ', 'Ǜ', 'Ǘ', 'Ǖ', 'Ǚ', 'Ủ', 'Ů', 'Ǔ', 'Ȕ', 'Ȗ', 'Ư', 'Ừ', 'Ứ', 'Ữ', 'Ử', 'Ự', 'Ụ', 'Ṳ', 'Ų', 'Ṷ', 'Ṵ', 'Ʉ':
		return []rune{'U'}
		//hungarian
	case 'Ú':
		return []rune{'Ú'}
	case 'Ü':
		return []rune{'Ü'}
	case 'Ű':
		return []rune{'Ű'}
	case 'V', 'Ⓥ', 'Ｖ', 'Ṽ', 'Ṿ', 'Ʋ', 'Ꝟ', 'Ʌ':
		return []rune{'V'}
	case 'Ꝡ':
		return []rune{'V', 'Y'}
	case 'W', 'Ⓦ', 'Ｗ', 'Ẁ', 'Ẃ', 'Ŵ', 'Ẇ', 'Ẅ', 'Ẉ', 'Ⱳ':
		return []rune{'w'}
	case 'X', 'Ⓧ', 'Ｘ', 'Ẋ', 'Ẍ':
		return []rune{'X'}
	case 'Y', 'Ⓨ', 'Ｙ', 'Ỳ', 'Ý', 'Ŷ', 'Ỹ', 'Ȳ', 'Ẏ', 'Ÿ', 'Ỷ', 'Ỵ', 'Ƴ', 'Ɏ', 'Ỿ':
		return []rune{'Y'}
	case 'Z', 'Ⓩ', 'Ｚ', 'Ź', 'Ẑ', 'Ż', 'Ž', 'Ẓ', 'Ẕ', 'Ƶ', 'Ȥ', 'Ɀ', 'Ⱬ', 'Ꝣ':
		return []rune{'Z'}

		// lower cases

	case 'a', 'ⓐ', 'ａ', 'ẚ', 'à', 'â', 'ầ', 'ấ', 'ẫ', 'ẩ', 'ã', 'ā', 'ă', 'ằ', 'ắ', 'ẵ', 'ẳ', 'ȧ', 'ǡ', 'ä', 'ǟ', 'ả', 'å', 'ǻ', 'ǎ', 'ȁ', 'ȃ', 'ạ', 'ậ', 'ặ', 'ḁ', 'ą', 'ⱥ', 'ɐ':
		return []rune{'a'}
		//hungarian
	case 'á':
		return []rune{'á'}
	case 'ꜳ':
		return []rune{'a', 'a'}
	case 'æ', 'ǽ', 'ǣ':
		return []rune{'a', 'e'}
	case 'ꜵ':
		return []rune{'a', 'o'}
	case 'ꜷ':
		return []rune{'a', 'u'}
	case 'ꜹ', 'ꜻ':
		return []rune{'a', 'v'}
	case 'ꜽ':
		return []rune{'a', 'y'}
	case 'b', 'ⓑ', 'ｂ', 'ḃ', 'ḅ', 'ḇ', 'ƀ', 'ƃ', 'ɓ':
		return []rune{'b'}
	case 'c', 'ⓒ', 'ｃ', 'ć', 'ĉ', 'ċ', 'č', 'ç', 'ḉ', 'ƈ', 'ȼ', 'ꜿ', 'ↄ':
		return []rune{'c'}
	case 'd', 'ⓓ', 'ｄ', 'ḋ', 'ď', 'ḍ', 'ḑ', 'ḓ', 'ḏ', 'đ', 'ƌ', 'ɖ', 'ɗ', 'ꝺ':
		return []rune{'d'}
	case 'ǳ', 'ǆ':
		return []rune{'d', 'z'}
	case 'e', 'ⓔ', 'ｅ', 'è', 'ê', 'ề', 'ế', 'ễ', 'ể', 'ẽ', 'ē', 'ḕ', 'ḗ', 'ĕ', 'ė', 'ë', 'ẻ', 'ě', 'ȅ', 'ȇ', 'ẹ', 'ệ', 'ȩ', 'ḝ', 'ę', 'ḙ', 'ḛ', 'ɇ', 'ɛ', 'ǝ':
		return []rune{'e'}
		//hungarian
	case 'é':
		return []rune{'é'}
	case 'f', 'ⓕ', 'ｆ', 'ḟ', 'ƒ', 'ꝼ':
		return []rune{'f'}
	case 'g', 'ⓖ', 'ｇ', 'ǵ', 'ĝ', 'ḡ', 'ğ', 'ġ', 'ǧ', 'ģ', 'ǥ', 'ɠ', 'ꞡ', 'ᵹ', 'ꝿ':
		return []rune{'g'}
	case 'h', 'ⓗ', 'ｈ', 'ĥ', 'ḣ', 'ḧ', 'ȟ', 'ḥ', 'ḩ', 'ḫ', 'ẖ', 'ħ', 'ⱨ':
		return []rune{'h'}
	case 'ƕ':
		return []rune{'h', 'v'}
	case 'i', 'ⓘ', 'ｉ', 'ì', 'î', 'ĩ', 'ī', 'ĭ', 'ï', 'ḯ', 'ỉ', 'ǐ', 'ȉ', 'ȋ', 'ị', 'į', 'ḭ', 'ɨ', 'ı':
		return []rune{'i'}
		//hungarian
	case 'í':
		return []rune{'í'}
	case 'j', 'ⓙ', 'ｊ', 'ĵ', 'ǰ', 'ɉ':
		return []rune{'j'}
	case 'k', 'ⓚ', 'ｋ', 'ḱ', 'ǩ', 'ḳ', 'ķ', 'ḵ', 'ƙ', 'ⱪ', 'ꝁ', 'ꝃ', 'ꝅ', 'ꞣ':
		return []rune{'k'}
	case 'l', 'ⓛ', 'ｌ', 'ŀ', 'ĺ', 'ľ', 'ḷ', 'ḹ', 'ļ', 'ḽ', 'ḻ', 'ſ', 'ł', 'ƚ', 'ɫ', 'ⱡ', 'ꝉ', 'ꞁ', 'ꝇ':
		return []rune{'l'}
	case 'ǉ':
		return []rune{'l', 'j'}
	case 'm', 'ⓜ', 'ｍ', 'ḿ', 'ṁ', 'ṃ', 'ɱ', 'ɯ':
		return []rune{'m'}
	case 'n', 'ⓝ', 'ｎ', 'ǹ', 'ń', 'ñ', 'ṅ', 'ň', 'ṇ', 'ņ', 'ṋ', 'ṉ', 'ƞ', 'ɲ', 'ŉ', 'ꞑ', 'ꞥ':
		return []rune{'n'}
	case 'ǌ':
		return []rune{'n', 'j'}
	case 'o', 'ⓞ', 'ｏ', 'ò', 'ô', 'ồ', 'ố', 'ỗ', 'ổ', 'õ', 'ṍ', 'ȭ', 'ṏ', 'ō', 'ṑ', 'ṓ', 'ŏ', 'ȯ', 'ȱ', 'ȫ', 'ỏ', 'ǒ', 'ȍ', 'ȏ', 'ơ', 'ờ', 'ớ', 'ỡ', 'ở', 'ợ', 'ọ', 'ộ', 'ǫ', 'ǭ', 'ø', 'ǿ', 'ɔ', 'ꝋ', 'ꝍ', 'ɵ':
		return []rune{'o'}
	//hungarian
	case 'ó':
		return []rune{'ó'}
	case 'ö':
		return []rune{'ö'}
	case 'ő':
		return []rune{'ő'}
	case 'ƣ':
		return []rune{'o', 'i'}
	case 'ȣ':
		return []rune{'o', 'u'}
	case 'ꝏ':
		return []rune{'o', 'o'}
	case '', 'œ':
		return []rune{'o', 'e'}
	case 'p', 'ⓟ', 'ｐ', 'ṕ', 'ṗ', 'ƥ', 'ᵽ', 'ꝑ', 'ꝓ', 'ꝕ':
		return []rune{'p'}
	case 'q', 'ⓠ', 'ｑ', 'ɋ', 'ꝗ', 'ꝙ':
		return []rune{'q'}
	case 'r', 'ⓡ', 'ｒ', 'ŕ', 'ṙ', 'ř', 'ȑ', 'ȓ', 'ṛ', 'ṝ', 'ŗ', 'ṟ', 'ɍ', 'ɽ', 'ꝛ', 'ꞧ', 'ꞃ':
		return []rune{'r'}
	case 's', 'ⓢ', 'ｓ', 'ß', 'ś', 'ṥ', 'ŝ', 'ṡ', 'š', 'ṧ', 'ṣ', 'ṩ', 'ș', 'ş', 'ȿ', 'ꞩ', 'ꞅ', 'ẛ':
		return []rune{'s'}
	case 't', 'ⓣ', 'ｔ', 'ṫ', 'ẗ', 'ť', 'ṭ', 'ț', 'ţ', 'ṱ', 'ṯ', 'ŧ', 'ƭ', 'ʈ', 'ⱦ', 'ꞇ':
		return []rune{'t'}
	case 'ꜩ':
		return []rune{'t', 'z'}
	case 'u', 'ⓤ', 'ｕ', 'ù', 'û', 'ũ', 'ṹ', 'ū', 'ṻ', 'ŭ', 'ǜ', 'ǘ', 'ǖ', 'ǚ', 'ủ', 'ů', 'ǔ', 'ȕ', 'ȗ', 'ư', 'ừ', 'ứ', 'ữ', 'ử', 'ự', 'ụ', 'ṳ', 'ų', 'ṷ', 'ṵ', 'ʉ':
		return []rune{'u'}
	//hungarian
	case 'ú':
		return []rune{'ú'}
	case 'ü':
		return []rune{'ü'}
	case 'ű':
		return []rune{'ű'}
	case 'v', 'ⓥ', 'ｖ', 'ṽ', 'ṿ', 'ʋ', 'ꝟ', 'ʌ':
		return []rune{'v'}
	case 'ꝡ':
		return []rune{'v', 'y'}
	case 'w', 'ⓦ', 'ｗ', 'ẁ', 'ẃ', 'ŵ', 'ẇ', 'ẅ', 'ẘ', 'ẉ', 'ⱳ':
		return []rune{'w'}
	case 'x', 'ⓧ', 'ｘ', 'ẋ', 'ẍ':
		return []rune{'x'}
	case 'y', 'ⓨ', 'ｙ', 'ỳ', 'ý', 'ŷ', 'ỹ', 'ȳ', 'ẏ', 'ÿ', 'ỷ', 'ẙ', 'ỵ', 'ƴ', 'ɏ', 'ỿ':
		return []rune{'y'}
	case 'z', 'ⓩ', 'ｚ', 'ź', 'ẑ', 'ż', 'ž', 'ẓ', 'ẕ', 'ƶ', 'ȥ', 'ɀ', 'ⱬ', 'ꝣ':
		return []rune{'z'}
	}
	return make([]rune, 0)
}
func removeHungarian(r rune) rune {
	switch r {
	case 'Á':
		return 'A'
	case 'É':
		return 'E'
	case 'Í':
		return 'I'
	case 'Ó', 'Ö', 'Ő':
		return 'O'
	case 'Ú', 'Ü', 'Ű':
		return 'U'
	case 'á':
		return 'a'
	case 'é':
		return 'e'
	case 'í':
		return 'i'
	case 'ó', 'ö', 'ő':
		return 'o'
	case 'ú', 'ü', 'ű':
		return 'u'
	default:
		return r
	}
}

func Convert(text string, removeHungarianletters bool, smallcaps bool) string {
	runes := make([]rune, 0)
	for _, r := range text {
		runes = append(runes, removeNotHungarian(r)...)
	}
	if removeHungarianletters {
		in := runes
		runes = make([]rune, 0)
		for _, r := range in {
			runes = append(runes, removeHungarian(r))
		}
	}
	ret := string(runes)
	if smallcaps {
		ret = strings.ToLower(ret)
	}
	return ret
}
