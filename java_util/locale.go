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

// Locales is all const Locale struct slice
var Locales []Locale = make([]Locale, 22, 22)

// init java.util.Locale static object
func init() {
	// Useful constant for language.
	Locales[ENGLISH] = Locale{
		id:     ENGLISH,
		Lang:   "en",
		County: "",
	}
	// Useful constant for language.
	Locales[FRENCH] = Locale{
		id:     FRENCH,
		Lang:   "fr",
		County: "",
	}
	// Useful constant for language.
	Locales[GERMAN] = Locale{
		id:     GERMAN,
		Lang:   "de",
		County: "",
	}
	// Useful constant for language.
	Locales[ITALIAN] = Locale{
		id:     ITALIAN,
		Lang:   "it",
		County: "",
	}
	// Useful constant for language.
	Locales[JAPANESE] = Locale{
		id:     JAPANESE,
		Lang:   "ja",
		County: "",
	}
	// Useful constant for language.
	Locales[KOREAN] = Locale{
		id:     KOREAN,
		Lang:   "ko",
		County: "",
	}
	// Useful constant for language.
	Locales[CHINESE] = Locale{
		id:     CHINESE,
		Lang:   "zh",
		County: "",
	}
	// Useful constant for language.
	Locales[SIMPLIFIED_CHINESE] = Locale{
		id:     SIMPLIFIED_CHINESE,
		Lang:   "zh",
		County: "CN",
	}
	// Useful constant for language.
	Locales[TRADITIONAL_CHINESE] = Locale{
		id:     TRADITIONAL_CHINESE,
		Lang:   "zh",
		County: "TW",
	}
	// Useful constant for language.
	Locales[FRANCE] = Locale{
		id:     FRANCE,
		Lang:   "fr",
		County: "FR",
	}
	// Useful constant for language.
	Locales[GERMANY] = Locale{
		id:     GERMANY,
		Lang:   "de",
		County: "DE",
	}
	// Useful constant for language.
	Locales[ITALY] = Locale{
		id:     ITALY,
		Lang:   "it",
		County: "it",
	}
	// Useful constant for language.
	Locales[JAPAN] = Locale{
		id:     JAPAN,
		Lang:   "ja",
		County: "JP",
	}
	// Useful constant for language.
	Locales[KOREA] = Locale{
		id:     KOREA,
		Lang:   "ko",
		County: "KR",
	}
	// Useful constant for language.
	Locales[CHINA] = Locales[SIMPLIFIED_CHINESE]
	// Useful constant for language.
	Locales[PRC] = Locales[SIMPLIFIED_CHINESE]
	// Useful constant for language.
	Locales[TAIWAN] = Locales[TRADITIONAL_CHINESE]
	// Useful constant for language.
	Locales[UK] = Locale{
		id:     UK,
		Lang:   "en",
		County: "GB",
	}
	// Useful constant for language.
	Locales[US] = Locale{
		id:     US,
		Lang:   "en",
		County: "US",
	}
	// Useful constant for language.
	Locales[CANADA] = Locale{
		id:     CANADA,
		Lang:   "en",
		County: "CA",
	}
	// Useful constant for language.
	Locales[CANADA_FRENCH] = Locale{
		id:     CANADA_FRENCH,
		Lang:   "fr",
		County: "CA",
	}
	// Useful constant for language.
	Locales[ROOT] = Locale{
		id:     ROOT,
		Lang:   "",
		County: "",
	}
}

// GetLocaleFromHandler is use LocaleHandle get Locale
func GetLocaleFromHandler(localeHandler *LocaleHandle) *Locale {
	// Useful constant for language.
	for _, locale := range Locales {
		if localeHandler.Value == locale.String() {
			return &locale
		}
	}
	return &Locales[ROOT]
}
