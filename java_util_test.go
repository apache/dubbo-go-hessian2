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

package hessian

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_util"
)

func TestJavaUtil(t *testing.T) {
	res, err := decodeJavaResponse(`customReplyUUID`, ``, false)
	if err != nil {
		t.Error(err)
		return
	}
	m := res.(map[interface{}]interface{})

	uuid1 := &java_util.UUID{LeastSigBits: int64(-7160773830801198154), MostSigBits: int64(459021424248441700)}

	resUuid1 := m["uuid1"]
	resUuid1String := m["uuid1_string"]
	resUuid2 := m["uuid2"]
	resUuid2String := m["uuid2_string"]

	assert.NotNil(t, resUuid1)
	assert.NotNil(t, resUuid1String)
	assert.NotNil(t, resUuid2)
	assert.NotNil(t, resUuid2String)

	assert.Equal(t, uuid1, resUuid1)
	assert.Equal(t, uuid1.String(), resUuid1String)
	assert.Equal(t, (resUuid2.(*java_util.UUID)).String(), resUuid2String)
}

func TestJavaUtilLocale(t *testing.T) {
	res, err := decodeJavaResponse(`customReplyLocale`, ``, false)
	if err != nil {
		t.Error(err)
		return
	}
	m := res.(map[interface{}]interface{})

	english := m["english"]
	french := m["french"]
	german := m["german"]
	italian := m["italian"]
	japanese := m["japanese"]
	korean := m["korean"]
	chinese := m["chinese"]
	simplified_chinese := m["simplified_chinese"]
	traditional_chinese := m["traditional_chinese"]
	france := m["france"]
	germany := m["germany"]
	japan := m["japan"]
	korea := m["korea"]
	china := m["china"]
	prc := m["prc"]
	taiwan := m["taiwan"]
	uk := m["uk"]
	us := m["us"]
	canada := m["canada"]
	root := m["root"]

	assert.NotNil(t, english)
	assert.NotNil(t, french)
	assert.NotNil(t, german)
	assert.NotNil(t, italian)
	assert.NotNil(t, japanese)
	assert.NotNil(t, korean)
	assert.NotNil(t, chinese)
	assert.NotNil(t, simplified_chinese)
	assert.NotNil(t, traditional_chinese)
	assert.NotNil(t, france)
	assert.NotNil(t, germany)
	assert.NotNil(t, japan)
	assert.NotNil(t, korea)
	assert.NotNil(t, china)
	assert.NotNil(t, prc)
	assert.NotNil(t, taiwan)
	assert.NotNil(t, uk)
	assert.NotNil(t, us)
	assert.NotNil(t, canada)
	//assert.NotNil(t, canada_french)
	assert.NotNil(t, root)

	assert.Equal(t, java_util.ToLocale(java_util.ENGLISH), java_util.GetLocaleFromHandler(english.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.FRENCH), java_util.GetLocaleFromHandler(french.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.GERMANY), java_util.GetLocaleFromHandler(germany.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.ITALIAN), java_util.GetLocaleFromHandler(italian.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.JAPANESE), java_util.GetLocaleFromHandler(japanese.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.KOREAN), java_util.GetLocaleFromHandler(korean.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.CHINESE), java_util.GetLocaleFromHandler(chinese.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.SIMPLIFIED_CHINESE), java_util.GetLocaleFromHandler(simplified_chinese.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.TRADITIONAL_CHINESE), java_util.GetLocaleFromHandler(traditional_chinese.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.FRANCE), java_util.GetLocaleFromHandler(france.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.GERMANY), java_util.GetLocaleFromHandler(germany.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.JAPAN), java_util.GetLocaleFromHandler(japan.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.KOREA), java_util.GetLocaleFromHandler(korea.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.CHINA), java_util.GetLocaleFromHandler(china.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.PRC), java_util.GetLocaleFromHandler(prc.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.TAIWAN), java_util.GetLocaleFromHandler(taiwan.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.UK), java_util.GetLocaleFromHandler(uk.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.US), java_util.GetLocaleFromHandler(us.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.CANADA), java_util.GetLocaleFromHandler(canada.(*java_util.LocaleHandle)))
	assert.Equal(t, java_util.ToLocale(java_util.ROOT), java_util.GetLocaleFromHandler(root.(*java_util.LocaleHandle)))
}
