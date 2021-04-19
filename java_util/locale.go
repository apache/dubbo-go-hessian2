/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_util

import "fmt"

// LocaleEnum is Locale enumeration value
type LocaleEnum int

// Locale struct enum
const (
	ENGLISH LocaleEnum = iota
	FRENCH
	GERMAN
	ITALIAN
	JAPANESE
	KOREAN
	CHINESE
	SIMPLIFIED_CHINESE
	TRADITIONAL_CHINESE
	FRANCE
	GERMANY
	ITALY
	JAPAN
	KOREA
	CHINA
	PRC
	TAIWAN
	UK
	US
	CANADA
	CANADA_FRENCH
	ROOT
)

// Locale => java.util.Locale
type Locale struct {
	// ID is used to implement enumeration
	id     LocaleEnum
	Lang   string
	County string
}

func (locale *Locale) String() string {
	if len(locale.County) != 0 {
		return fmt.Sprintf("%s_%s", locale.Lang, locale.County)
	}
	return locale.Lang
}

// LocaleHandle => com.alibaba.com.caucho.hessian.io.LocaleHandle object
type LocaleHandle struct {
	Value string `hessian:"value"`
}

func (LocaleHandle) JavaClassName() string {
	return "com.alibaba.com.caucho.hessian.io.LocaleHandle"
}

// locales is all const Locale struct slice
// localeMap is key = locale.String() value = locale struct
var (
	locales   []Locale            = make([]Locale, 22, 22)
	localeMap map[string](Locale) = make(map[string](Locale), 22)
)

// init java.util.Locale static object
func init() {
	locales[ENGLISH] = Locale{
		id:     ENGLISH,
		Lang:   "en",
		County: "",
	}
	locales[FRENCH] = Locale{
		id:     FRENCH,
		Lang:   "fr",
		County: "",
	}
	locales[GERMAN] = Locale{
		id:     GERMAN,
		Lang:   "de",
		County: "",
	}
	locales[ITALIAN] = Locale{
		id:     ITALIAN,
		Lang:   "it",
		County: "",
	}
	locales[JAPANESE] = Locale{
		id:     JAPANESE,
		Lang:   "ja",
		County: "",
	}
	locales[KOREAN] = Locale{
		id:     KOREAN,
		Lang:   "ko",
		County: "",
	}
	locales[CHINESE] = Locale{
		id:     CHINESE,
		Lang:   "zh",
		County: "",
	}
	locales[SIMPLIFIED_CHINESE] = Locale{
		id:     SIMPLIFIED_CHINESE,
		Lang:   "zh",
		County: "CN",
	}
	locales[TRADITIONAL_CHINESE] = Locale{
		id:     TRADITIONAL_CHINESE,
		Lang:   "zh",
		County: "TW",
	}
	locales[FRANCE] = Locale{
		id:     FRANCE,
		Lang:   "fr",
		County: "FR",
	}
	locales[GERMANY] = Locale{
		id:     GERMANY,
		Lang:   "de",
		County: "DE",
	}
	locales[ITALY] = Locale{
		id:     ITALY,
		Lang:   "it",
		County: "it",
	}
	locales[JAPAN] = Locale{
		id:     JAPAN,
		Lang:   "ja",
		County: "JP",
	}
	locales[KOREA] = Locale{
		id:     KOREA,
		Lang:   "ko",
		County: "KR",
	}
	locales[CHINA] = locales[SIMPLIFIED_CHINESE]
	locales[PRC] = locales[SIMPLIFIED_CHINESE]
	locales[TAIWAN] = locales[TRADITIONAL_CHINESE]
	locales[UK] = Locale{
		id:     UK,
		Lang:   "en",
		County: "GB",
	}
	locales[US] = Locale{
		id:     US,
		Lang:   "en",
		County: "US",
	}
	locales[CANADA] = Locale{
		id:     CANADA,
		Lang:   "en",
		County: "CA",
	}
	locales[CANADA_FRENCH] = Locale{
		id:     CANADA_FRENCH,
		Lang:   "fr",
		County: "CA",
	}
	locales[ROOT] = Locale{
		id:     ROOT,
		Lang:   "",
		County: "",
	}
	for _, locale := range locales {
		localeMap[locale.String()] = locale
	}
}

// ToLocale get locale from enum
func ToLocale(e LocaleEnum) *Locale {
	return &locales[e]
}

// GetLocaleFromHandler is use LocaleHandle get Locale
func GetLocaleFromHandler(localeHandler *LocaleHandle) *Locale {
	locale := localeMap[localeHandler.Value]
	return &locale
}
